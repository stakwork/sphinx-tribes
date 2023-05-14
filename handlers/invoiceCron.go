package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
)

func InitInvoiceCron() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Seconds().Do(func() {
		invoiceCount, _ := db.Store.GetInvoiceCount(config.InvoiceCount)

		if invoiceCount > 0 {
			url := fmt.Sprintf("%s/invoices", config.RelayUrl)

			client := &http.Client{}
			req, err := http.NewRequest(http.MethodGet, url, nil)

			req.Header.Set("x-user-token", config.RelayAuthKey)
			req.Header.Set("Content-Type", "application/json")
			res, _ := client.Do(req)

			if err != nil {
				log.Printf("Request Failed: %s", err)
				return
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)

			// Unmarshal result
			invoiceRes := db.InvoiceList{}
			err = json.Unmarshal(body, &invoiceRes)

			if err != nil {
				log.Printf("Reading body failed: %s", err)
				return
			}

			for _, v := range invoiceRes.Invoices {
				if v.Settled {
					storeInvoice, _ := db.Store.GetInvoiceCache(v.Payment_request)
					if storeInvoice.Invoice == v.Payment_request {
						/**
						If the invoice is settled and still in store
						make keysend payment
						*/

						url := fmt.Sprintf("%s/payment", config.RelayUrl)

						bodyData := fmt.Sprintf(`{"amount": %s, "destination_key": "%s"}`, storeInvoice.Amount, storeInvoice.User_pubkey)

						jsonBody := []byte(bodyData)

						client := &http.Client{}
						req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

						req.Header.Set("x-user-token", config.RelayAuthKey)
						req.Header.Set("Content-Type", "application/json")
						res, _ := client.Do(req)

						if err != nil {
							log.Printf("Request Failed: %s", err)
							return
						}

						defer res.Body.Close()

						body, err = ioutil.ReadAll(res.Body)

						if res.StatusCode == 200 {
							// Unmarshal result
							keysendRes := db.KeysendSuccess{}
							err = json.Unmarshal(body, &keysendRes)

							var p = db.DB.GetPersonByPubkey(storeInvoice.Owner_pubkey)

							wanteds, _ := p.Extras["wanted"].([]interface{})

							for _, wanted := range wanteds {
								w, ok2 := wanted.(map[string]interface{})
								if !ok2 {
									continue // next wanted
								}

								created, ok3 := w["created"].(float64)
								createdArr := strings.Split(fmt.Sprintf("%f", created), ".")
								createdString := createdArr[0]
								createdInt, _ := strconv.ParseInt(createdString, 10, 32)

								dateInt, _ := strconv.ParseInt(storeInvoice.Created, 10, 32)

								if !ok3 {
									continue
								}

								if createdInt == dateInt {
									w["paid"] = true
								}
							}

							p.Extras["wanted"] = wanteds
							b := new(bytes.Buffer)
							decodeErr := json.NewEncoder(b).Encode(p.Extras)

							if decodeErr != nil {
								log.Printf("Could not encode extras json data")
							} else {
								db.DB.UpdatePerson(p.ID, map[string]interface{}{
									"extras": b,
								})

								// Delete the invoice from store
								db.Store.DeleteCache(storeInvoice.Invoice)

								invoiceCount, _ := db.Store.GetInvoiceCount(config.InvoiceCount)

								if invoiceCount > 0 {
									// reduce the invoice count
									db.Store.SetInvoiceCount(config.InvoiceCount, invoiceCount-1)
								}
							}
						} else {
							// Unmarshal result
							keysendError := db.KeysendError{}
							err = json.Unmarshal(body, &keysendError)
							log.Printf("Keysend Payment to %s Failed, with Error: %s", storeInvoice.User_pubkey, keysendError.Error)
						}

						if err != nil {
							log.Printf("Reading body failed: %s", err)
							return
						}
					}
				}
			}

		}
	})

	s.StartAsync()
}
