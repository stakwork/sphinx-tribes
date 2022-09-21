package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TicketMessage struct {
	Ticket   PropertyMap `json:"ticket"`
	ChatUuid string      `json:"chat_uuid"`
	Status   string      `json:"status"` // In case we want to add more notifications other than just opening ticket
	Action   string      `json:"action"`
	Members  []Person    `json:"members"`
}

func processAlerts(p Person) {
	// get the existing person's github_issues
	// compare to see if there are new ones
	// check the new ones for coding languages like "#Javascript"

	// people need a new "extras"->"alerts" that lets them toggle on and off alerts
	// pull people who have alerts on
	// of those people, who have "Javascript" in their "coding_languages"?
	// if they match, build an Action with their pubkey
	// post all the Actions you have build to relay with the HMAC header

	relayUrl := os.Getenv("ALERT_URL")
	alertSecret := os.Getenv("ALERT_SECRET")
	alertTribeUuid := os.Getenv("ALERT_TRIBE_UUID")
	if relayUrl == "" || alertSecret == "" || alertTribeUuid == "" {
		fmt.Println("Ticket alerts: ENV information not found")
		return
	}

	var message TicketMessage
	message.ChatUuid = alertTribeUuid
	message.Status = "open"
	message.Action = "dm"

	// Check that new ticket time exists
	if p.newTicketTime == 0 {
		fmt.Println("Ticket alerts: New ticket time not found")
		return
	}

	var issue PropertyMap = nil
	wanteds, ok := p.Extras["wanted"].([]interface{})
	if !ok {
		fmt.Println("Ticket alerts: No tickets found for person")
	}
	for _, wanted := range wanteds {
		w, ok2 := wanted.(map[string]interface{})
		if !ok2 {
			continue
		}
		time, ok3 := w["created"].(int64)
		if !ok3 {
			continue
		}
		if time == p.newTicketTime {
			issue = w
			break
		}
	}

	if issue == nil {
		fmt.Println("Ticket alerts: No ticket identified with the correct timestamp")
	}

	message.Ticket = issue
	languages, ok4 := issue["codingLanguage"].([]interface{})
	if !ok4 {
		fmt.Println("Ticket alerts: No languages found in ticket")
		return
	}

	var err error
	message.Members, err = DB.getPeopleForNewTicket(languages)
	if err != nil {
		fmt.Println("Ticket alerts: DB query to get interested people failed", err)
		return
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(message)
	if err != nil {
		fmt.Println("Ticket alerts: Unable to parse message into byte buffer", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", relayUrl, &buf)
	if err != nil {
		fmt.Println("Ticket alerts: Unable to create a request to send to relay", err)
		return
	}

	bytes := buf.Bytes()

	secret, err := hex.DecodeString(alertSecret)
	if err != nil {
		fmt.Println("Ticket alerts: Unable to create a byte array for secret", err)
		return
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write(bytes)
	hmac256Byte := mac.Sum(nil)
	hmac256Hex := hex.EncodeToString(hmac256Byte)

	request.Header.Set("x-hub-signature-256", hmac256Hex)
	_, err = client.Do(request)
	if err != nil {
		fmt.Println("Ticket alerts: Unable to communicate request to relay", err)
		return
	}
	return
}
