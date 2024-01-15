package db

import (
	"fmt"
	"github.com/stakwork/sphinx-tribes/utils"
	"math"
	"net/http"
)

var SecondsToDateConversion = 60 * 60 * 24

func (db database) TotalPeopleByDateRange(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Person{}).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) TotalOrganizationsByDateRange(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Organization{}).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) TotalPaymentsByDateRange(r PaymentDateRange) uint {
	var sum uint
	db.db.Model(&PaymentHistory{}).Where("payment_type = ?", r.PaymentType).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(amount)").Row().Scan(&sum)
	return sum
}

func (db database) TotalSatsPosted(r PaymentDateRange) uint {
	var sum uint
	db.db.Model(&Bounty{}).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(price)").Row().Scan(&sum)
	return sum
}

func (db database) TotalSatsPaid(r PaymentDateRange) uint {
	var sum uint
	db.db.Model(&Bounty{}).Where("paid = ?", true).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(price)").Row().Scan(&sum)
	return sum
}

func (db database) SatsPaidPercentage(r PaymentDateRange) uint {
	satsPosted := DB.TotalSatsPosted(r)
	satsPaid := DB.TotalSatsPaid(r)
	if satsPaid != 0 && satsPosted != 0 {
		value := (satsPaid * 100) / satsPosted
		paidPercentage := math.Round(float64(value))
		return uint(paidPercentage)
	}
	return 0
}

func (db database) TotalPaidBounties(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("paid = ?", true).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) TotalHuntersPaid(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Bounty{}).Select("DISTINCT assignee").Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) NewHuntersPaid(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Bounty{}).
		Select("DISTINCT assignee").
		Where("paid = true").
		Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).
		Not("assignee IN (?)", db.db.Model(&Bounty{}).
			Select("assignee").
			Where("paid = true").
			Where("created < ?", r.StartDate),
		).Count(&count)
	return count
}

func (db database) TotalBountiesPosted(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) BountiesPaidPercentage(r PaymentDateRange) uint {
	bountiesPosted := DB.TotalBountiesPosted(r)
	bountiesPaid := DB.TotalPaidBounties(r)
	if bountiesPaid != 0 && bountiesPosted != 0 {
		value := bountiesPaid * 100 / bountiesPosted
		paidPercentage := math.Round(float64(value))
		return uint(paidPercentage)
	}
	return 0
}

func (db database) PaidDifference(r PaymentDateRange) []DateDifference {
	ms := []DateDifference{}

	db.db.Raw(`SELECT EXTRACT(EPOCH FROM (paid_date - TO_TIMESTAMP(created))) as diff FROM public.bounty WHERE paid_date IS NOT NULL AND created >= '` + r.StartDate + `' AND created <= '` + r.EndDate + `' `).Find(&ms)
	return ms
}

func (db database) PaidDifferenceCount(r PaymentDateRange) int64 {
	var count int64
	list := db.PaidDifference(r)
	count = int64(len(list))
	return count
}

func (db database) AveragePaidTime(r PaymentDateRange) uint {
	paidList := DB.PaidDifference(r)
	paidCount := DB.PaidDifferenceCount(r)
	var paidSum uint
	for _, diff := range paidList {
		paidSum = uint(math.Round(diff.Diff))
	}
	return CalculateAverageDays(paidCount, paidSum)
}

func (db database) CompletedDifference(r PaymentDateRange) []DateDifference {
	ms := []DateDifference{}

	db.db.Raw(`SELECT EXTRACT(EPOCH FROM (completion_date - TO_TIMESTAMP(created))) as diff FROM public.bounty WHERE completion_date IS NOT NULL AND created >= '` + r.StartDate + `' AND created <= '` + r.EndDate + `' `).Find(&ms)
	return ms
}

func (db database) CompletedDifferenceCount(r PaymentDateRange) int64 {
	var count int64
	list := db.CompletedDifference(r)
	count = int64(len(list))
	return count
}

func (db database) AverageCompletedTime(r PaymentDateRange) uint {
	paidList := DB.CompletedDifference(r)
	paidCount := DB.CompletedDifferenceCount(r)
	var paidSum uint
	for _, diff := range paidList {
		paidSum = uint(math.Round(diff.Diff))
	}
	return CalculateAverageDays(paidCount, paidSum)
}

func CalculateAverageDays(paidCount int64, paidSum uint) uint {
	if paidCount != 0 && paidSum != 0 {
		avg := paidSum / uint(paidCount)
		avgSeconds := math.Round(float64(avg))
		avgDays := math.Round(avgSeconds / float64(SecondsToDateConversion))
		return uint(avgDays)
	}
	return 0
}

func (db database) GetBountiesByDateRange(r PaymentDateRange, re *http.Request) []Bounty {
	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(re)
	keys := re.URL.Query()
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")

	openQuery := ""
	assignedQuery := ""
	paidQuery := ""
	orderQuery := ""
	limitQuery := ""

	if open != "" && open == "true" {
		openQuery = "AND assignee = '' AND paid != true"
	}
	if assingned != "" && assingned == "true" {
		assignedQuery = "AND assignee != '' AND paid = false"
	}
	if paid != "" && paid == "true" {
		paidQuery = "AND paid = true"
	}

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = " ORDER BY " + sortBy + "" + "DESC"
	}
	if limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	query := `SELECT * FROM public.bounty WHERE created >= '` + r.StartDate + `'  AND created <= '` + r.EndDate + `'`
	allQuery := query + " " + openQuery + " " + assignedQuery + " " + paidQuery + " " + orderQuery + " " + limitQuery

	b := []Bounty{}
	db.db.Raw(allQuery).Find(&b)
	return b
}

func (db database) GetBountiesByDateRangeCount(r PaymentDateRange, re *http.Request) int64 {
	keys := re.URL.Query()
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")

	openQuery := ""
	assignedQuery := ""
	paidQuery := ""

	if open != "" && open == "true" {
		openQuery = "AND assignee = '' AND paid != true"
		assignedQuery = ""
	}
	if assingned != "" && assingned == "true" {
		if open != "" && open == "true" {
			assignedQuery = "OR assignee != '' AND paid = false"
		} else {
			assignedQuery = "AND assignee != '' AND paid = false"
		}
	}
	if paid != "" && paid == "true" {
		if open != "" && open == "true" || assingned != "" && assingned == "true" {
			paidQuery = "OR paid = true"
		} else if open != "" && open == "true" && assingned == "" && assingned != "true" {
			assignedQuery = ""
		} else {
			paidQuery = "AND paid = true"
		}
	}

	var count int64

	query := `SELECT COUNT(*) FROM public.bounty WHERE created >= '` + r.StartDate + `'  AND created <= '` + r.EndDate + `'`
	allQuery := query + " " + openQuery + " " + assignedQuery + " " + paidQuery
	db.db.Raw(allQuery).Scan(&count)
	return count
}
