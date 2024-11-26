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
	"github.com/stakwork/sphinx-tribes/utils"
)

func InitV2PaymentsCron() {
	paymentHistories := db.DB.GetPendingPaymentHistory()
	for _, payment := range paymentHistories {
		tag := payment.Tag
		tagResult := GetInvoiceStatusByTag(tag)

		bounty := db.DB.GetBounty(payment.BountyId)

		if bounty.ID > 0 {

			if tagResult.Status == db.PaymentComplete {
				db.DB.SetPaymentAsComplete(tag)

				now := time.Now()

				bounty.PaymentPending = false
				bounty.PaymentFailed = false
				bounty.Paid = true

				bounty.PaidDate = &now
				bounty.Completed = true
				bounty.CompletionDate = &now

				db.DB.UpdateBountyPaymentStatuses(bounty)

			} else if tagResult.Status == db.PaymentFailed {
				// Handle failed payments

				err := db.DB.ProcessReversePayments(payment.ID)
				if err != nil {
					log.Printf("Could not reverse bounty payment : Bounty ID - %d, Payment ID - %d, Error - %s", bounty.ID, payment.ID, err)
				}

			} else if tagResult.Status == db.PaymentPending {
				if payment.PaymentStatus == db.PaymentPending {
					created := utils.ConvertTimeToTimestamp(payment.Created.String())

					now := time.Now()
					daysDiff := utils.GetDateDaysDifference(int64(created), &now)

					if daysDiff >= 7 {

						err := db.DB.ProcessReversePayments(payment.ID)
						if err != nil {
							log.Printf("Could not reverse bounty payment after 7 days : Bounty ID - %d, Payment ID - %d, Error - %s", bounty.ID, payment.ID, err)
						}
					}
				}
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
