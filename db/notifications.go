package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) CreateNotification(notification *Notification) error {
	if notification.UUID == "" {
		notification.UUID = uuid.New().String()
	}

	if notification.Event == "" {
		return errors.New("event is required")
	}
	if notification.PubKey == "" {
		return errors.New("public key is required")
	}
	if notification.Content == "" {
		return errors.New("content is required")
	}

	now := time.Now()
	notification.CreatedAt = &now
	notification.UpdatedAt = &now
	notification.Status = NotificationStatusPending

	if err := db.db.Create(&notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func (db database) UpdateNotification(uuid string, updates map[string]interface{}) error {
	if uuid == "" {
		return errors.New("notification UUID is required")
	}

	if len(updates) == 0 {
		return errors.New("no updates provided")
	}

	var existingNotification Notification
	if err := db.db.First(&existingNotification, "uuid = ?", uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("notification not found")
		}
		return fmt.Errorf("failed to fetch notification: %w", err)
	}

	now := time.Now()
	updates["updated_at"] = &now

	if status, ok := updates["status"].(NotificationStatus); ok {
		if status == "" {
			return errors.New("status cannot be empty")
		}
	}

	if err := db.db.Model(&existingNotification).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

func (db database) GetNotification(uuid string) (*Notification, error) {
	if uuid == "" {
		return nil, errors.New("uuid is required")
	}

	var notification Notification
	if err := db.db.Where("uuid = ?", uuid).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch notification: %w", err)
	}

	return &notification, nil
}

func (db database) DeleteNotification(uuid string) error {
	if uuid == "" {
		return errors.New("uuid is required")
	}

	result := db.db.Where("uuid = ?", uuid).Delete(&Notification{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (db database) GetPendingNotifications() ([]Notification, error) {
	var notifications []Notification
	if err := db.db.Where("status = ?", NotificationStatusPending).
		Order("created_at ASC").
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch pending notifications: %w", err)
	}

	return notifications, nil
}

func (db database) GetFailedNotifications(maxRetries int) ([]Notification, error) {
	if maxRetries < 0 {
		return nil, errors.New("maxRetries must be non-negative")
	}

	var notifications []Notification
	if err := db.db.Where("status = ? AND retries < ?", NotificationStatusFailed, maxRetries).
		Order("created_at ASC").
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch failed notifications: %w", err)
	}

	return notifications, nil
}

func (db database) GetNotificationsByPubKey(pubKey string, limit, offset int) ([]Notification, error) {
	if pubKey == "" {
		return nil, errors.New("public key is required")
	}

	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset must be non-negative")
	}

	var notifications []Notification
	if err := db.db.Where("pub_key = ?", pubKey).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	return notifications, nil
}

func (db database) IncrementRetryCount(uuid string) error {
	if uuid == "" {
		return errors.New("uuid is required")
	}

	result := db.db.Model(&Notification{}).
		Where("uuid = ?", uuid).
		Updates(map[string]interface{}{
			"retries":    gorm.Expr("retries + 1"),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to increment retry count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (db database) GetNotificationCount(pubKey string) (int64, error) {
	if pubKey == "" {
		return 0, errors.New("public key is required")
	}

	var count int64
	if err := db.db.Model(&Notification{}).
		Where("pub_key = ?", pubKey).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get notification count: %w", err)
	}

	return count, nil
}
