package utils

import (
	"crypto/rand"
	"encoding/base32"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"strings"
	"time"

	decodepay "github.com/nbd-wtf/ln-decodepay"
	"github.com/stakwork/sphinx-tribes/logger"
)

func GetRandomTokens(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		logger.Log.Error("Random token error: %v", err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func ConvertStringToUint(number string) (uint, error) {
	numberParse, err := strconv.ParseUint(number, 10, 32)

	if err != nil {
		logger.Log.Error("could not parse string to uint: %v", err)
		return 0, err
	}

	return uint(numberParse), nil
}

func ConvertStringToInt(number string) (int, error) {
	numberParse, err := strconv.ParseInt(number, 10, 32)

	if err != nil {
		logger.Log.Error("could not parse string to int: %v", err)
		return 0, err
	}

	return int(numberParse), nil
}

func GetInvoiceAmount(paymentRequest string) uint {
	decodedInvoice, err := decodepay.Decodepay(paymentRequest)

	if err != nil {
		logger.Log.Error("Could not Decode Invoice: %v", err)
		return 0
	}
	amountInt := decodedInvoice.MSatoshi / 1000
	amount := uint(amountInt)

	return amount
}

func GetInvoiceExpired(paymentRequest string) bool {
	decodedInvoice, err := decodepay.Decodepay(paymentRequest)
	if err != nil {
		logger.Log.Error("Could not Decode Invoice: %v", err)
		return false
	}

	timeInUnix := time.Now().Unix()
	expiryDate := decodedInvoice.CreatedAt + decodedInvoice.Expiry

	if timeInUnix > int64(expiryDate) {
		return true
	} else {
		return false
	}
}

func ConvertTimeToTimestamp(date string) int {
	format := "2006-01-02 15:04:05"
	dateTouse := date

	if strings.Contains(date, "+") {
		dateSplit := strings.Split(date, "+")
		dateTouse = strings.Trim(dateSplit[0], " ")
	}

	t, err := time.Parse(format, dateTouse)
	if err != nil {
		logger.Log.Error("Parse string to timestamp: %v", err)
	} else {
		return int(t.Unix())
	}
	return 0
}

func AddHoursToTimestamp(timestamp int, hours int) int {
	tm := time.Unix(int64(timestamp), 0)

	dur := int(time.Hour.Hours()) * hours
	tm = tm.Add(time.Hour * time.Duration(dur))

	return int(tm.Unix())
}

func GetDateDaysDifference(createdDate int64, paidDate *time.Time) int64 {
	firstDate := time.Unix(createdDate, 0)
	difference := paidDate.Sub(firstDate)
	days := int64(difference.Hours() / 24)
	return days
}

func GetHoursDifference(createdDate int64, endDate *time.Time) int64 {
	firstDate := time.Unix(createdDate, 0)
	difference := endDate.Sub(firstDate)

	hours := int64(difference.Hours())
	return hours
}

func isValidUUID(uuid string) bool {
	return len(uuid) == 36
}

func ValidateUUID(r *http.Request) bool {
	uuid := chi.URLParam(r, "uuid")
	return isValidUUID(uuid)
}
