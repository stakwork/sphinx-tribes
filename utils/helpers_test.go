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
