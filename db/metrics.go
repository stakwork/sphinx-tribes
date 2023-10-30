package db

func (db database) TotalPaymentsByDateRange(r PaymentDateRange) uint {
	var sum uint

	db.db.Model(&PaymentHistory{}).Where("payment_type = ?", r.PaymentType).Where("created >= ?", r.StartDate).Where("created <= ?", r.EndDate).Select("SUM(amount)").Row().Scan(&sum)

	return sum
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
