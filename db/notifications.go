package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrMissingEvent   = errors.New("event is required")
	ErrMissingPubKey  = errors.New("pub_key is required")
	ErrMissingContent = errors.New("content is required")
	ErrMissingUUID    = errors.New("uuid is required")
)

func (db database) CreateNotification(n *Notification) error {

	if n.Event == "" {
		return ErrMissingEvent
	}
	if n.PubKey == "" {
		return ErrMissingPubKey
	}
	if n.Content == "" {
		return ErrMissingContent
	}

	if n.UUID == "" {
		n.UUID = uuid.New().String()
	}

	now := time.Now()
	n.CreatedAt = &now
	n.UpdatedAt = &now

	if n.Status == "" {
		n.Status = NotificationStatusPending
	}

	return db.db.Create(n).Error
}

func (db database) GetNotification(uuid string) (*Notification, error) {
	if uuid == "" {
		return nil, ErrMissingUUID
	}

	var notification Notification
	err := db.db.Where("uuid = ?", uuid).First(&notification).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &notification, err
}

func (db database) UpdateNotification(uuid string, updates map[string]interface{}) error {
	if uuid == "" {
		return ErrMissingUUID
	}

	if event, ok := updates["event"].(string); ok && event == "" {
		return ErrMissingEvent
	}
	if pubKey, ok := updates["pub_key"].(string); ok && pubKey == "" {
		return ErrMissingPubKey
	}
	if content, ok := updates["content"].(string); ok && content == "" {
		return ErrMissingContent
	}

	updates["updated_at"] = time.Now()
	result := db.db.Model(&Notification{}).Where("uuid = ?", uuid).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (db database) DeleteNotification(uuid string) error {
	if uuid == "" {
		return ErrMissingUUID
	}

	result := db.db.Where("uuid = ?", uuid).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (db database) GetPendingNotifications() ([]Notification, error) {
	var notifications []Notification
	err := db.db.Where("status = ?", NotificationStatusPending).Find(&notifications).Error
	return notifications, err
}

func (db database) GetFailedNotifications(maxRetries int) ([]Notification, error) {
	if maxRetries < 0 {
		maxRetries = 0
	}

	var notifications []Notification
	err := db.db.Where("status = ? AND retries < ?", NotificationStatusFailed, maxRetries).Find(&notifications).Error
	return notifications, err
}

func (db database) GetNotificationsByPubKey(pubKey string, limit, offset int) ([]Notification, error) {
	if pubKey == "" {
		return nil, ErrMissingPubKey
	}

	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}

	var notifications []Notification
	query := db.db.Where("pub_key = ?", pubKey)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (db database) IncrementRetryCount(uuid string) error {
	if uuid == "" {
		return ErrMissingUUID
	}

	result := db.db.Model(&Notification{}).
		Where("uuid = ?", uuid).
		UpdateColumn("retries", gorm.Expr("retries + 1"))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (db database) GetNotificationCount(pubKey string) (int64, error) {
	if pubKey == "" {
		return 0, ErrMissingPubKey
	}

	var count int64
	err := db.db.Model(&Notification{}).Where("pub_key = ?", pubKey).Count(&count).Error
	return count, err
}
