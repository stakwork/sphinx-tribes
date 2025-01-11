package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
)

const (
	NotificationPending          = "PENDING"
	NotificationComplete         = "COMPLETE"
	NotificationFailed           = "FAILED"
	NotificationWaitingKeyExchange = "WAITING_KEY_EXCHANGE"
)

func sendNotification(pubkey string, event string, content string, retries int) {
	// Generate notification content based on the event
	content = generateNotificationContent(event, content)

	// 1. Verify user on v2 bot:
	contactKey, err := verifyUserOnV2Bot(pubkey)
	if err != nil {
		logger.Log.Error("Failed to verify user on v2 bot: %v", err)
		return
	}

	// 2. Send the notification:
	status, err := sendNotificationToUser(pubkey, content)
	if err != nil {
		logger.Log.Error("Failed to send notification: %v", err)
		return
	}

	// 3. Log the result:
	notification := db.Notification{
		Event:    event,
		Pubkey:   pubkey,
		Content:  content,
		Retries:  retries,
		Status:   status,
		UUID:     utils.GenerateUUID(),
	}
	err = db.DB.CreateNotification(&notification)
	if err != nil {
		logger.Log.Error("Failed to log notification: %v", err)
	}
}

func generateNotificationContent(event string, content string) string {
	switch event {
	case "bounty_assigned":
		return fmt.Sprintf("You have been assigned a new bounty: %s", content)
	case "bounty_paid":
		return fmt.Sprintf("Your bounty has been paid: %s", content)
	default:
		return content
	}
}

func verifyUserOnV2Bot(pubkey string) (string, error) {
	url := fmt.Sprintf("%s/contact/%s", utils.V2BotURL, pubkey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var contact struct {
		ContactKey string `json:"contact_key"`
	}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		return "", err
	}

	if contact.ContactKey == "" {
		// Add user to bot's known contacts
		err = addUserToBotContacts(pubkey)
		if err != nil {
			return "", err
		}

		// Re-verify user
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(body, &contact)
		if err != nil {
			return "", err
		}

		if contact.ContactKey == "" {
			return "", fmt.Errorf("contact key still missing after adding user")
		}
	}

	return contact.ContactKey, nil
}

func addUserToBotContacts(pubkey string) error {
	url := fmt.Sprintf("%s/add_contact", utils.V2BotURL)
	payload := map[string]string{
		"pubkey": pubkey,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add user to bot's contacts")
	}

	return nil
}

func sendNotificationToUser(pubkey string, content string) (string, error) {
	url := fmt.Sprintf("%s/send", utils.V2BotURL)
	payload := map[string]interface{}{
		"dest":      pubkey,
		"amt_msat":  0,
		"content":   content,
		"is_tribe":  false,
		"wait":      true,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}

func processPendingNotifications() {
	notifications, err := db.DB.GetPendingNotifications()
	if err != nil {
		logger.Log.Error("Failed to get pending notifications: %v", err)
		return
	}

	for _, notification := range notifications {
		sendNotification(notification.Pubkey, notification.Event, notification.Content, notification.Retries+1)
	}
}
