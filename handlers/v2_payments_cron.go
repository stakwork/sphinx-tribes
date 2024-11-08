package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
)

func InitV2PaymentsCron() {
	paymentHistories := db.DB.GetPendingPaymentHistory()
	for _, payment := range paymentHistories {
		tag := payment.Tag
		tagResult := GetInvoiceStatusByTag(tag)

		if tagResult.Status == db.PaymentComplete {
			db.DB.SetPaymentAsComplete(tag)

			bounty := db.DB.GetBounty(payment.ID)

			if bounty.ID > 0 {
				now := time.Now()

				bounty.Paid = true
				bounty.PaymentPending = false
				bounty.PaymentFailed = false
				bounty.PaidDate = &now
				bounty.Completed = true
				bounty.CompletionDate = &now

				db.DB.UpdateBounty(bounty)
			}
		} else if tagResult.Status == db.PaymentFailed {
			// Handle failed payments
			bounty := db.DB.GetBounty(payment.ID)

			if bounty.ID > 0 {
				db.DB.SetPaymentStatusByBountyId(bounty.ID, tagResult)

				bounty.Paid = false
				bounty.PaymentPending = false
				bounty.PaymentFailed = true

				db.DB.UpdateBounty(bounty)
			}
		}
	}
}

func GetInvoiceStatusByTag(tag string) db.V2TagRes {
	url := fmt.Sprintf("%s/sends/%s", config.V2BotUrl, tag)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Printf("Error paying invoice: %s", err)
	}

	req.Header.Set("x-admin-token", config.V2BotToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Printf("[Get Tag] Request Failed: %s", err)
		return db.V2TagRes{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Could not read body: %s", err)
	}

	tagRes := []db.V2TagRes{}
	err = json.Unmarshal(body, &tagRes)

	if err != nil {
		log.Printf("Could not unmarshall get tag result: %s", err)
	}

	resultLength := len(tagRes)

	if resultLength > 0 {
		return tagRes[0]
	}

	return db.V2TagRes{}
}
