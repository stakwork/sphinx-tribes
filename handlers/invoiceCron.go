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
	msg := make(map[string]interface{})

	s.Every(1).Seconds().Do(func() {
		invoiceList, _ := db.Store.GetInvoiceCache()
		invoiceCount := len(invoiceList)

		if invoiceCount > 0 {

			for index, inv := range invoiceList {

				url := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, inv.Invoice)

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
				invoiceRes := db.InvoiceResult{}

				err = json.Unmarshal(body, &invoiceRes)

				if err != nil {
					log.Printf("Reading body failed: %s", err)
					return
				}

				if invoiceRes.Response.Settled {
					if inv.Invoice == invoiceRes.Response.Payment_request {
						/**
						  If the invoice is settled and still in store
						  make keysend payment
						*/
						msg["msg"] = "invoice_success"
						msg["invoice"] = inv.Invoice

						socket, err := db.Store.GetSocketConnections(inv.Host)
						if err == nil {
							socket.Conn.WriteJSON(msg)
						}

						if inv.Type == "KEYSEND" {
							url := fmt.Sprintf("%s/payment", config.RelayUrl)

							bodyData := fmt.Sprintf(`{"amount": %s, "destination_key": "%s"}`, inv.Amount, inv.User_pubkey)

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

								var p = db.DB.GetPersonByPubkey(inv.Owner_pubkey)

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

									dateInt, _ := strconv.ParseInt(inv.Created, 10, 32)

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

									// Delete the index from the store array list and reset the store
									updateInvoiceCache(invoiceList, index)

									msg["msg"] = "keysend_success"
									msg["invoice"] = inv.Invoice

									socket, err := db.Store.GetSocketConnections(inv.Host)
									if err == nil {
										socket.Conn.WriteJSON(msg)
									}
								}
							} else {
								// Unmarshal result
								keysendError := db.KeysendError{}
								err = json.Unmarshal(body, &keysendError)
								log.Printf("Keysend Payment to %s Failed, with Error: %s", inv.User_pubkey, keysendError.Error)

								msg["msg"] = "keysend_error"
								msg["invoice"] = inv.Invoice

								socket, err := db.Store.GetSocketConnections(inv.Host)

								if err == nil {
									socket.Conn.WriteJSON(msg)
								}

								updateInvoiceCache(invoiceList, index)
							}

							if err != nil {
								log.Printf("Reading body failed: %s", err)
								return
							}
						} else {
							var p = db.DB.GetPersonByPubkey(inv.Owner_pubkey)

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

								dateInt, _ := strconv.ParseInt(inv.Created, 10, 32)

								if !ok3 {
									continue
								}

								if createdInt == dateInt {
									var user = db.DB.GetPersonByPubkey(inv.User_pubkey)

									var assignee = make(map[string]interface{})

									assignee["img"] = user.Img
									assignee["label"] = user.OwnerAlias
									assignee["value"] = user.OwnerPubKey
									assignee["owner_pubkey"] = user.OwnerPubKey
									assignee["owner_alias"] = user.OwnerAlias
									assignee["commitment_fee"] = inv.Commitment_fee
									assignee["assigned_hours"] = inv.Assigned_hours
									assignee["bounty_expires"] = inv.Bounty_expires

									w["assignee"] = assignee
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

								// Delete the index from the store array list and reset the store
								updateInvoiceCache(invoiceList, index)

								msg := make(map[string]interface{})
								msg["msg"] = "assign_success"
								msg["invoice"] = inv.Invoice

								socket, err := db.Store.GetSocketConnections(inv.Host)
								if err == nil {
									socket.Conn.WriteJSON(msg)
								}
							}
						}
					}
				}
			}
		}
	})

	s.StartAsync()
}

func updateInvoiceCache(invoiceList []db.InvoiceStoreData, index int) {
	newInvoiceList := append(invoiceList[:index], invoiceList[index+1:]...)
	db.Store.SetInvoiceCache(newInvoiceList)
}
