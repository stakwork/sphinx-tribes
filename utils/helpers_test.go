package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDateDaysDifference(t *testing.T) {
	created := 1700238956

	testDate := time.Now()
	days := GetDateDaysDifference(int64(created), &testDate)
	assert.NotEqual(t, 0, days)
}
