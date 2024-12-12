package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (db database) AddChat(chat *Chat) (Chat, error) {
	if chat.ID == "" {
		return Chat{}, errors.New("chat ID is required")
	}

	now := time.Now()
	chat.CreatedAt = now
	chat.UpdatedAt = now

	if err := db.db.Create(&chat).Error; err != nil {
		return Chat{}, fmt.Errorf("failed to create chat: %w", err)
	}

	return *chat, nil
}

func (db database) GetChatByChatID(chatID string) (Chat, error) {
	var chat Chat
	result := db.db.Where("id = ?", chatID).First(&chat)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Chat{}, fmt.Errorf("chat not found")
		}
		return Chat{}, fmt.Errorf("failed to fetch chat: %w", result.Error)
	}

	return chat, nil
}

func (db database) AddChatMessage(chatMessage *ChatMessage) (ChatMessage, error) {
	if chatMessage.ID == "" {
		return ChatMessage{}, errors.New("message ID is required")
	}

	now := time.Now()
	chatMessage.Timestamp = now

	if err := db.db.Create(&chatMessage).Error; err != nil {
		return ChatMessage{}, fmt.Errorf("failed to create chat message: %w", err)
	}

	return *chatMessage, nil
}

func (db database) UpdateChatMessage(chatMessage *ChatMessage) (ChatMessage, error) {
	if chatMessage.ID == "" {
		return ChatMessage{}, errors.New("message ID is required")
	}

	var existingMessage ChatMessage
	if err := db.db.First(&existingMessage, "id = ?", chatMessage.ID).Error; err != nil {
		return ChatMessage{}, fmt.Errorf("message not found: %w", err)
	}

	if chatMessage.Message != "" {
		existingMessage.Message = chatMessage.Message
	}
	if chatMessage.Status != "" {
		existingMessage.Status = chatMessage.Status
	}
	if chatMessage.Role != "" {
		existingMessage.Role = chatMessage.Role
	}

	existingMessage.Timestamp = time.Now()

	if err := db.db.Save(&existingMessage).Error; err != nil {
		return ChatMessage{}, fmt.Errorf("failed to update chat message: %w", err)
	}

	return existingMessage, nil
}

func (db database) GetChatMessagesForChatID(chatID string) ([]ChatMessage, error) {
	var chatMessages []ChatMessage

	result := db.db.Where("chat_id = ?", chatID).Order("timestamp ASC").Find(&chatMessages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch chat messages: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return []ChatMessage{}, nil
	}

	return chatMessages, nil
}
