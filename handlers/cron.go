package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

func InitNotificationCron() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minute().Do(func() {
		processPendingNotifications()
	})

	s.Every(5).Minutes().Do(func() {
		recheckContactKey()
	})

	s.StartAsync()
}

func recheckContactKey() {
	notifications, err := db.DB.GetNotificationsByStatus(NotificationWaitingKeyExchange)
	if err != nil {
		logger.Log.Error("Failed to get notifications with status WAITING_KEY_EXCHANGE: %v", err)
		return
	}

	for _, notification := range notifications {
		contactKey, err := verifyUserOnV2Bot(notification.Pubkey)
		if err != nil {
			logger.Log.Error("Failed to verify user on v2 bot: %v", err)
			continue
		}

		if contactKey != "" {
			notification.Status = NotificationPending
			err = db.DB.UpdateNotification(&notification)
			if err != nil {
				logger.Log.Error("Failed to update notification status: %v", err)
			}
		}
	}
}
