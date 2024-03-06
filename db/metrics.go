package db

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/stakwork/sphinx-tribes/utils"
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
	query := fmt.Sprintf(`SELECT COUNT(DISTINCT assignee) FROM bounty WHERE assignee !='' AND paid=true AND created >= %s AND created <= %s`, r.StartDate, r.EndDate)

	db.db.Raw(query).Count(&count)
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
	providers := keys.Get("provider")

	orderQuery := ""
	limitQuery := ""

	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = " ORDER BY " + sortBy + "" + "DESC"
	}
	if limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	providerCondition := ""
	if len(providers) > 0 {
		providerSlice := strings.Split(providers, ",")
		providerCondition = " AND owner_id IN ('" + strings.Join(providerSlice, "','") + "')"
	}

	query := `SELECT * FROM public.bounty WHERE created >= '` + r.StartDate + `'  AND created <= '` + r.EndDate + `'` + providerCondition
	allQuery := query + " " + statusQuery + " " + orderQuery + " " + limitQuery

	b := []Bounty{}
	db.db.Raw(allQuery).Find(&b)
	return b
}

func (db database) GetBountiesByDateRangeCount(r PaymentDateRange, re *http.Request) int64 {
	keys := re.URL.Query()
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	providers := keys.Get("provider")

	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	providerCondition := ""
	if len(providers) > 0 {
		providerSlice := strings.Split(providers, ",")
		providerCondition = " AND owner_id IN ('" + strings.Join(providerSlice, "','") + "')"
	}

	var count int64

	query := `SELECT COUNT(*) FROM public.bounty WHERE created >= '` + r.StartDate + `'  AND created <= '` + r.EndDate + `'` + providerCondition
	allQuery := query + " " + statusQuery
	db.db.Raw(allQuery).Scan(&count)
	return count
}

func (db database) GetBountiesProviders(r PaymentDateRange, re *http.Request) []Person {
	offset, limit, _, _, _ := utils.GetPaginationParams(re)
	keys := re.URL.Query()
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	providers := keys.Get("provider")

	var statusConditions []string

	limitQuery := ""

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	providerCondition := ""
	if len(providers) > 0 {
		providerSlice := strings.Split(providers, ",")
		providerCondition = " AND owner_id IN ('" + strings.Join(providerSlice, "','") + "')"
	}

	if limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	bountyOwners := []BountyOwners{}
	bountyProviders := []Person{}

	query := `SELECT DISTINCT owner_id FROM public.bounty WHERE created >= '` + r.StartDate + `'  AND created <= '` + r.EndDate + `'` + providerCondition
	allQuery := query + " " + statusQuery + " " + limitQuery
	db.db.Raw(allQuery).Scan(&bountyOwners)

	for _, owner := range bountyOwners {
		person := db.GetPersonByPubkey(owner.OwnerID)
		bountyProviders = append(bountyProviders, person)
	}
	return bountyProviders
}
