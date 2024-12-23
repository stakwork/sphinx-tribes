package db

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/stakwork/sphinx-tribes/logger"
)

type Action struct {
	Action   string `json:"action"`    // "dm"
	ChatUuid string `json:"chat_uuid"` // tribe uuid
	Pubkey   string `json:"pubkey"`
	Content  string `json:"content"`
	BotId    string `json:"bot_id"`
}

func (db database) ProcessAlerts(p Person) {
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
	botId := os.Getenv("ALERT_BOT_ID")
	if relayUrl == "" || alertSecret == "" || alertTribeUuid == "" || botId == "" {
		logger.Log.Info("Ticket alerts: ENV information not found")
		return
	}

	var action Action
	action.ChatUuid = alertTribeUuid
	action.Action = "dm"
	action.BotId = botId
	action.Content = "A new ticket relevant to your interests has been created on Sphinx Community - https://community.sphinx.chat/p?owner_id="
	action.Content += p.OwnerPubKey
	action.Content += "&created=" + strconv.Itoa(int(p.NewTicketTime))

	// Check that new ticket time exists
	if p.NewTicketTime == 0 {
		logger.Log.Info("Ticket alerts: New ticket time not found")
		return
	}

	var issue PropertyMap = nil
	wanteds, ok := p.Extras["wanted"].([]interface{})
	if !ok {
		logger.Log.Info("Ticket alerts: No tickets found for person")
	}
	for _, wanted := range wanteds {
		w, ok2 := wanted.(map[string]interface{})
		if !ok2 {
			continue
		}
		timeF, ok3 := w["created"].(float64)
		if !ok3 {
			continue
		}
		time := int64(timeF)
		if time == p.NewTicketTime {
			issue = w
			break
		}
	}

	if issue == nil {
		logger.Log.Info("Ticket alerts: No ticket identified with the correct timestamp")
	}

	languages, ok4 := issue["codingLanguage"].([]interface{})
	if !ok4 {
		logger.Log.Info("Ticket alerts: No languages found in ticket")
		return
	}

	var err error
	people, err := db.GetPeopleForNewTicket(languages)
	if err != nil {
		logger.Log.Error("Ticket alerts: DB query to get interested people failed: %v", err)
		return
	}

	client := http.Client{}

	for _, per := range people {
		action.Pubkey = per.OwnerPubKey
		buf, err := json.Marshal(action)
		if err != nil {
			logger.Log.Error("Ticket alerts: Unable to parse message into byte buffer: %v", err)
			return
		}
		request, err := http.NewRequest("POST", relayUrl, bytes.NewReader(buf))
		if err != nil {
			logger.Log.Error("Ticket alerts: Unable to create a request to send to relay: %v", err)
			return
		}

		mac := hmac.New(sha256.New, []byte(alertSecret))
		mac.Write(buf)
		hmac256Byte := mac.Sum(nil)
		hmac256Hex := "sha256=" + hex.EncodeToString(hmac256Byte)
		request.Header.Set("x-hub-signature-256", hmac256Hex)
		request.Header.Set("Content-Type", "application/json")
		_, err = client.Do(request)
		if err != nil {
			logger.Log.Error("Ticket alerts: Unable to communicate request to relay: %v", err)
		}
	}

	return
}
