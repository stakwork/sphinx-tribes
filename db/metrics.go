package db

import "math"

func (db database) TotalPaymentsByDateRange(r PaymentDateRange) uint {
	var sum uint

	db.db.Model(&PaymentHistory{}).Where("payment_type = ?", r.PaymentType).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(amount)").Row().Scan(&sum)

	return sum
}

func (db database) TotalSatsPaid(r PaymentDateRange) uint {
	var sum uint

	db.db.Model(&Bounty{}).Where("paid_date >= ?", r.StartDate).Where("paid_date <= ?", r.EndDate).Select("SUM(amount)").Row().Scan(&sum)

	return sum
}

func (db database) TotalSatsPosted(r PaymentDateRange) uint {
	var sum uint

	db.db.Model(&Bounty{}).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(amount)").Row().Scan(&sum)

	return sum
}

func (db database) SatsPaidPercentage(r PaymentDateRange) uint {
	satsPosted := DB.TotalSatsPosted(r)
	satsPaid := DB.TotalSatsPaid(r)

	value := satsPosted * 100 / satsPaid
	paidPercentage := math.Round(float64(value))
	return uint(paidPercentage)
}

func (db database) TotalPaidBounties(r PaymentDateRange) int64 {
	var count int64

	db.db.Model(&Bounty{}).Where("paid_date >= ?", r.StartDate).Where("paid_date <= ?", r.EndDate).Count(&count)

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

	value := bountiesPaid * 100 / bountiesPosted
	paidPercentage := math.Round(float64(value))
	return uint(paidPercentage)
}

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
