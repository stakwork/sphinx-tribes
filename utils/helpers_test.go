package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDateDaysDifference(t *testing.T) {
	created := 1700238956
	exactTime := time.Now().Unix()

	testDate := time.Now()
	nextDate := time.Now().AddDate(0, 0, 2)
	days := GetDateDaysDifference(int64(created), &testDate)

	assert.NotEqual(t, 0, days)

	daysEqual := GetDateDaysDifference(exactTime, &testDate)
	assert.Equal(t, int64(0), daysEqual)

	daysNext := GetDateDaysDifference(exactTime, &nextDate)
	assert.Equal(t, int64(2), daysNext)

	wrongDate := GetDateDaysDifference(0, &nextDate)
	assert.Greater(t, wrongDate, int64(365))
}

func TestGetRandomToken(t *testing.T) {
	randomToken32 := GetRandomToken(32)
	assert.GreaterOrEqual(t, len(randomToken32), 32)

	randomToken64 := GetRandomToken(56)
	assert.GreaterOrEqual(t, len(randomToken64), 56)
}

func TestConvertStringToUint(t *testing.T) {
	number := "20"
	result, _ := ConvertStringToUint(number)

	assert.Equal(t, uint(20), result)

	wrongNum := "wrong"
	result2, err := ConvertStringToUint(wrongNum)
	assert.Equal(t, uint(0), result2)

	assert.NotEqual(t, err, nil)
}

func TestConvertStringToInt(t *testing.T) {
	number := "10"
	result, _ := ConvertStringToInt(number)

	assert.Equal(t, int(10), result)

	wrongNum := "wrong"
	result2, err := ConvertStringToInt(wrongNum)
	assert.Equal(t, int(0), result2)

	assert.NotEqual(t, err, nil)
}

func TestGetInvoiceAmount(t *testing.T) {
	invoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"

	amount := GetInvoiceAmount(invoice)
	assert.Equal(t, uint(1500), amount)

	invalidInvoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq"

	amount2 := GetInvoiceAmount(invalidInvoice)
	assert.Equal(t, uint(0), amount2)
}

func TestGetInvoiceExpired(t *testing.T) {
	expiredInvoice := "lnbcrt100u1pnr5gtzpp5r7ew6nzqd9y9w5ktsspftnckxdn3te0y04n9mw7c6hkkrznh4pgsdqhgf6kgem9wssyjmnkda5kxegcqzpgxqyz5vqsp5mc09mpl4l3rllnfl3y902yxa29flke8r4ertqswdcrk766z5nq4q9qyyssq7wteenxtwlxatsd8dqdncqnn6u23jmcpe0d7ne6dcpafwlx9ckr3dp6y4p7sl4j3pq6l93g6vc4w8z04ry9yzwjv6cggm06eecad9psp9dh6u5"
	isInvoiceExpired := GetInvoiceExpired(expiredInvoice)
	assert.Equal(t, true, isInvoiceExpired)
}

func TestConvertTimeToTimestamp(t *testing.T) {
	dateWithPlus := "2024-10-16 09:21:21.743327+00"
	dateWithoutPlus := "2024-10-16 09:21:21.743327"

	dateTimestamp1 := ConvertTimeToTimestamp(dateWithPlus)
	dateTimestamp2 := ConvertTimeToTimestamp(dateWithoutPlus)

	assert.Greater(t, dateTimestamp1, 7000000)
	assert.Greater(t, dateTimestamp2, 7000000)
}

func TestGetHoursDifference(t *testing.T) {
	time1 := time.Now().Unix()
	time2 := time.Now().Add(time.Hour * 1)

	hourDiff := GetHoursDifference(time1, &time2)
	assert.Equal(t, hourDiff, int64(1))
}
