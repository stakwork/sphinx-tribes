package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) CreateSSEMessageLog(event map[string]interface{}, chatID, from, to string) (*SSEMessageLog, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}
	if from == "" {
		return nil, errors.New("source URL is required")
	}
	if to == "" {
		return nil, errors.New("target URL is required")
	}

	now := time.Now()
	messageLog := &SSEMessageLog{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Event:     event,
		ChatID:    chatID,
		From:      from,
		To:        to,
		Status:    SSEStatusNew,
	}

	if err := db.db.Create(messageLog).Error; err != nil {
		return nil, fmt.Errorf("failed to create SSE message log: %w", err)
	}

	return messageLog, nil
}

func (db database) DeleteSSEMessageLog(id uuid.UUID) error {
	result := db.db.Delete(&SSEMessageLog{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete SSE message log: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("SSE message log with ID %s not found", id)
	}
	return nil
}

func (db database) UpdateSSEMessageLogStatusBatch(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided for batch update")
	}

	result := db.db.Model(&SSEMessageLog{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     SSEStatusSent,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update SSE message logs batch: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no SSE message logs found with the provided IDs")
	}

	return nil
}

func (db database) UpdateSSEMessageLog(id uuid.UUID, updates map[string]interface{}) (*SSEMessageLog, error) {
	var messageLog SSEMessageLog
	if err := db.db.First(&messageLog, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("SSE message log with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find SSE message log: %w", err)
	}

	eventUpdate, hasEvent := updates["event"]
	if hasEvent {
		delete(updates, "event")
		
		if eventMap, ok := eventUpdate.(map[string]interface{}); ok {
			messageLog.Event = eventMap
		}
	}

	updates["updated_at"] = time.Now()

	if err := db.db.Model(&messageLog).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update SSE message log: %w", err)
	}

	if hasEvent {
		if err := db.db.Model(&messageLog).Update("event", messageLog.Event).Error; err != nil {
			return nil, fmt.Errorf("failed to update event field: %w", err)
		}
	}

	if err := db.db.First(&messageLog, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve updated SSE message log: %w", err)
	}

	return &messageLog, nil
}

func (db database) GetSSEMessageLogByID(id uuid.UUID) (*SSEMessageLog, error) {
	var messageLog SSEMessageLog
	if err := db.db.First(&messageLog, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("SSE message log with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve SSE message log: %w", err)
	}
	return &messageLog, nil
}

func (db database) GetSSEMessageLogsByChatID(chatID string) ([]SSEMessageLog, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	var messageLogs []SSEMessageLog
	if err := db.db.Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Find(&messageLogs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve SSE message logs for chat %s: %w", chatID, err)
	}

	return messageLogs, nil
}

func (db database) GetNewSSEMessageLogsByChatID(chatID string) ([]SSEMessageLog, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	var messageLogs []SSEMessageLog
	if err := db.db.Where("chat_id = ? AND status = ?", chatID, SSEStatusNew).
		Order("created_at DESC").
		Find(&messageLogs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve new SSE message logs for chat %s: %w", chatID, err)
	}

	return messageLogs, nil
} 