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
	log.Println("Pending Invoice Cron Job Started")
	paymentHistories := db.DB.GetPendingPaymentHistory()
	for _, payment := range paymentHistories {
		bounty := db.DB.GetBounty(payment.BountyId)
		log.Println("Bounty ID =========================", bounty.ID, bounty)
		log.Println("Payment ID =========================", payment.ID, payment)

		if bounty.ID > 0 {
			tag := payment.Tag
			tagResult := GetInvoiceStatusByTag(tag)
			log.Println("Tag Result =================", tagResult)

			if tagResult.Status == db.PaymentComplete {
				log.Println("Payment Status From V2 BOT IS Complete =================================", payment)
				db.DB.SetPaymentAsComplete(tag)

				now := time.Now()

				bounty.PaymentPending = false
				bounty.PaymentFailed = false
				bounty.Paid = true

				bounty.PaidDate = &now
				bounty.Completed = true
				bounty.CompletionDate = &now

				db.DB.UpdateBountyPaymentStatuses(bounty)
				log.Println("Bounty Payment Statuses Updated =================================", bounty)
			} else if tagResult.Status == db.PaymentPending {
				log.Println("Payment Status From V2 BOT IS Pending =================================", payment)
				if payment.PaymentStatus == db.PaymentPending {
					log.Println("Payment Status From DB IS Pending =================================	", payment)
					created := utils.ConvertTimeToTimestamp(payment.Created.String())

					now := time.Now()
					daysDiff := utils.GetDateDaysDifference(int64(created), &now)

					log.Println("Payment Date Difference Is ================================================", daysDiff)

					if daysDiff >= 7 {

						log.Println("Payment Date Difference Is Greater Or Equals 7 Days ================================================", payment)

						err := db.DB.ProcessReversePayments(payment.ID)
						if err != nil {
							log.Printf("Could not reverse bounty payment after 7 days : Bounty ID - %d, Payment ID - %d, Error - %s ================================================", bounty.ID, payment.ID, err)
						}

						log.Println("Bounty Payment Statuses Updated After 7 Days ================================================", bounty)
					}
				}
			} else if tagResult.Status == db.PaymentFailed {
				// Handle failed payments
				err := db.DB.ProcessReversePayments(payment.ID)
				if err != nil {
					log.Printf("Could not reverse bounty payment : Bounty ID - %d, Payment ID - %d, Error - %s ================================================", bounty.ID, payment.ID, err)
				}

				log.Println("Bounty Payment Statuses Updated After Failed Payment ================================================", bounty)

			} else {
				log.Println("Payment Status From V2 BOT IS Unknown ================================================", payment, tagResult)
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
		log.Printf("Could not unmarshal get tag result: %s", err)
	}

	resultLength := len(tagRes)

	if resultLength > 0 {
		return tagRes[0]
	}

	return db.V2TagRes{}
}
