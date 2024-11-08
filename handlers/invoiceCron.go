package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

// Keep for future reference
func InitInvoiceCron() {
	s := gocron.NewScheduler(time.UTC)
	msg := make(map[string]interface{})

	s.Every(5).Seconds().Do(func() {
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

				body, err := io.ReadAll(res.Body)

				// Unmarshal result
				invoiceRes := db.InvoiceResult{}

				err = json.Unmarshal(body, &invoiceRes)

				if err != nil {
					log.Printf("Reading Invoice body failed: %s", err)
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

							amount, _ := utils.ConvertStringToUint(inv.Amount)

							bodyData := utils.BuildKeysendBodyData(amount, inv.User_pubkey, inv.Route_hint, "")

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

							body, err = io.ReadAll(res.Body)

							if res.StatusCode == 200 {
								// Unmarshal result
								keysendRes := db.KeysendSuccess{}
								err = json.Unmarshal(body, &keysendRes)

								dateInt, _ := strconv.ParseInt(inv.Created, 10, 32)
								bounty, err := db.DB.GetBountyByCreated(uint(dateInt))

								if err == nil {
									bounty.Paid = true
								}

								db.DB.UpdateBounty(bounty)

								// Delete the index from the store array list and reset the store
								updateInvoiceCache(invoiceList, index)

								msg["msg"] = "keysend_success"
								msg["invoice"] = inv.Invoice

								socket, err := db.Store.GetSocketConnections(inv.Host)
								if err == nil {
									socket.Conn.WriteJSON(msg)
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
							dateInt, _ := strconv.ParseInt(inv.Created, 10, 32)
							bounty, err := db.DB.GetBountyByCreated(uint(dateInt))

							if err == nil {
								bounty.Assignee = inv.User_pubkey
								bounty.CommitmentFee = uint64(inv.Commitment_fee)
								bounty.AssignedHours = uint8(inv.Assigned_hours)
								bounty.BountyExpires = inv.Bounty_expires
							}

							db.DB.UpdateBounty(bounty)

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
	})

	s.Every(5).Seconds().Do(func() {
		invoiceList, _ := db.Store.GetBudgetInvoiceCache()
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

				body, err := io.ReadAll(res.Body)

				// Unmarshal result
				invoiceRes := db.InvoiceResult{}

				err = json.Unmarshal(body, &invoiceRes)

				if err != nil {
					log.Printf("Reading Workspace Invoice body failed: %s", err)
					return
				}

				if invoiceRes.Response.Settled {
					if inv.Invoice == invoiceRes.Response.Payment_request {
						/**
						  If the invoice is settled and still in store
						  make keysend payment
						*/
						msg["msg"] = "budget_success"
						msg["invoice"] = inv.Invoice

						socket, err := db.Store.GetSocketConnections(inv.Host)

						if err == nil {
							socket.Conn.WriteJSON(msg)
						}

						// db.DB.AddAndUpdateBudget(inv)
						updateBudgetInvoiceCache(invoiceList, index)
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

func updateBudgetInvoiceCache(invoiceList []db.BudgetStoreData, index int) {
	newInvoiceList := append(invoiceList[:index], invoiceList[index+1:]...)
	db.Store.SetBudgetInvoiceCache(newInvoiceList)
}
