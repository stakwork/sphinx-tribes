package db

import (
	"fmt"
	"math"
	"net/http"

	"github.com/stakwork/sphinx-tribes/utils"
)

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

func (db database) PaidDifferenceSum(r PaymentDateRange) uint {
	var sum uint
	db.db.Model(&Bounty{}).Where("paid_date_difference != ?",
		"").Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(paid_date_difference)").Row().Scan(&sum)
	return sum
}

func (db database) PaidDifferenceCount(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("paid_date_difference != ?",
		"").Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) AveragePaidTime(r PaymentDateRange) uint {
	paidSum := DB.PaidDifferenceSum(r)
	paidCount := DB.PaidDifferenceCount(r)
	if paidCount != 0 && paidSum != 0 {
		avg := paidSum / uint(paidCount)
		avgDays := math.Round(float64(avg))
		return uint(avgDays)
	}
	return 0
}

func (db database) CompletedDifferenceSum(r PaymentDateRange) uint {
	var sum uint
	db.db.Model(&Bounty{}).Where("completion_date_difference != ?",
		"").Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(completion_date_difference)").Row().Scan(&sum)
	return sum
}

func (db database) CompletedDifferenceCount(r PaymentDateRange) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("completion_date_difference != ?",
		"").Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Count(&count)
	return count
}

func (db database) AverageCompletedTime(r PaymentDateRange) uint {
	paidSum := DB.CompletedDifferenceSum(r)
	paidCount := DB.CompletedDifferenceCount(r)
	if paidCount != 0 && paidSum != 0 {
		avg := paidSum / uint(paidCount)
		avgDays := math.Round(float64(avg))
		return uint(avgDays)
	}
	return 0
}

func (db database) GetBountiesByDateRange(r PaymentDateRange, re *http.Request) []Bounty {
	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(re)

	orderQuery := ""
	limitQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + "body." + sortBy + " " + direction
	} else {
		orderQuery = " ORDER BY " + "body." + sortBy + "" + "DESC"
	}
	if limit != 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	query := `SELECT * public.bounty WHERE created >= '` + r.StartDate + `'  AND created <= '` + r.EndDate + `'`
	allQuery := query + " " + " " + orderQuery + " " + limitQuery

	b := []Bounty{}
	db.db.Raw(allQuery).Scan(&b)
	return b
}
