package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (db database) UpdateChat(chat *Chat) (Chat, error) {
	if chat.ID == "" {
		return Chat{}, errors.New("chat ID is required")
	}

	var existingChat Chat
	if err := db.db.First(&existingChat, "id = ?", chat.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Chat{}, fmt.Errorf("chat not found")
		}
		return Chat{}, fmt.Errorf("failed to fetch chat: %w", err)
	}

	if chat.Title != "" {
		existingChat.Title = chat.Title
	}
	if chat.Status != "" {
		existingChat.Status = chat.Status
	}
	existingChat.UpdatedAt = time.Now()

	if err := db.db.Save(&existingChat).Error; err != nil {
		return Chat{}, fmt.Errorf("failed to update chat: %w", err)
	}

	return existingChat, nil
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

func (db database) GetChatsForWorkspace(workspaceID string, chatStatus string) ([]Chat, error) {

	if workspaceID == "" {
		return []Chat{}, errors.New("workspace ID is required")
	}

	var chats []Chat

	if chatStatus == "" {
		chatStatus = string(ActiveStatus)
	}

	result := db.db.Where("workspace_id = ? AND status = ?", workspaceID, chatStatus).
		Order("updated_at DESC").
		Find(&chats)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch chats: %w", result.Error)
	}

	return chats, nil
}

func (db database) GetAllChatsForWorkspace(workspaceID string) ([]Chat, error) {
	var chats []Chat
	result := db.db.Where("workspace_id = ?", workspaceID).
		Order("updated_at DESC").
		Find(&chats)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch chats: %w", result.Error)
	}
	return chats, nil
}

func (db database) CreateFileAsset(asset *FileAsset) (*FileAsset, error) {
	now := time.Now()
	asset.CreatedAt = now
	asset.UpdatedAt = now
	asset.UploadTime = now
	asset.LastReferenced = now

	if err := db.db.Create(asset).Error; err != nil {
		return nil, fmt.Errorf("failed to create file asset: %w", err)
	}
	return asset, nil
}

func (db database) GetFileAssetByHash(fileHash string) (*FileAsset, error) {
	var asset FileAsset
	if err := db.db.Where("file_hash = ? AND status != ?", fileHash, DeletedFileStatus).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (db database) GetFileAssetByID(id uint) (*FileAsset, error) {
	var asset FileAsset
	if err := db.db.First(&asset, id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (db database) UpdateFileAssetReference(id uint) error {
	result := db.db.Model(&FileAsset{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_referenced": time.Now(),
			"status":          ActiveFileStatus,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no file asset found with id %d", id)
	}
	return nil
}

func (db database) ListFileAssets(params ListFileAssetsParams) ([]FileAsset, int64, error) {
	var assets []FileAsset
	var total int64

	query := db.db.Model(&FileAsset{})

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.MimeType != nil {
		query = query.Where("mime_type = ?", *params.MimeType)
	}
	if params.BeforeDate != nil {
		query = query.Where("upload_time <= ?", *params.BeforeDate)
	}
	if params.AfterDate != nil {
		query = query.Where("upload_time >= ?", *params.AfterDate)
	}
	if params.LastAccessedBefore != nil {
		query = query.Where("last_referenced <= ?", *params.LastAccessedBefore)
	}
	if params.WorkspaceID != nil {
		query = query.Where("workspace_id = ?", *params.WorkspaceID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).
		Limit(params.PageSize).
		Order("upload_time DESC").
		Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

func (db database) UpdateFileAsset(asset *FileAsset) error {
	asset.UpdatedAt = time.Now()
	return db.db.Save(asset).Error
}

func (db database) DeleteFileAsset(id uint) error {

	var asset FileAsset
	if err := db.db.First(&asset, id).Error; err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	now := time.Now()
	result := db.db.Model(&FileAsset{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     DeletedFileStatus,
			"deleted_at": &now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no file asset found with id %d", id)
	}
	return nil
}

func (db database) CreateOrEditChatWorkflow(workflow *ChatWorkflow) (*ChatWorkflow, error) {
	var existing ChatWorkflow
	result := db.db.Where("workspace_id = ?", workflow.WorkspaceID).First(&existing)

	now := time.Now()
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		workflow.CreatedAt = now
		workflow.UpdatedAt = now
		if err := db.db.Create(workflow).Error; err != nil {
			return nil, fmt.Errorf("failed to create chat workflow: %w", err)
		}
		return workflow, nil
	} else if result.Error != nil {
		return nil, fmt.Errorf("failed to check existing workflow: %w", result.Error)
	}

	existing.URL = workflow.URL
	existing.StackworkID = workflow.StackworkID
	existing.UpdatedAt = now

	if err := db.db.Save(&existing).Error; err != nil {
		return nil, fmt.Errorf("failed to update chat workflow: %w", err)
	}
	return &existing, nil
}

func (db database) GetChatWorkflowByWorkspaceID(workspaceID string) (*ChatWorkflow, error) {
	var workflow ChatWorkflow
	if err := db.db.Where("workspace_id = ?", workspaceID).First(&workflow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat workflow not found for workspace: %s", workspaceID)
		}
		return nil, fmt.Errorf("failed to fetch chat workflow: %w", err)
	}
	return &workflow, nil
}

func (db database) DeleteChatWorkflow(workspaceID string) error {
	result := db.db.Where("workspace_id = ?", workspaceID).Delete(&ChatWorkflow{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete chat workflow: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("chat workflow not found for workspace: %s", workspaceID)
	}
	return nil
}

func (db database) AddChatStatus(status *ChatWorkflowStatus) (ChatWorkflowStatus, error) {
	if status.ChatID == "" {
		return ChatWorkflowStatus{}, errors.New("chat ID is required")
	}

	if status.Status == "success" {
		status.Message = "Agent workflow completed"
	} else if status.Status == "error" && status.Message == "" {
		status.Message = "An error occurred during the workflow"
	} else if status.Message == "" {
		status.Message = status.Status
	}

	now := time.Now()
	status.CreatedAt = now
	status.UpdatedAt = now

	if err := db.db.Create(&status).Error; err != nil {
		return ChatWorkflowStatus{}, fmt.Errorf("failed to create chat status: %w", err)
	}

	return *status, nil
}

func (db database) UpdateChatStatus(status *ChatWorkflowStatus) (ChatWorkflowStatus, error) {
	if status.UUID == uuid.Nil {
		return ChatWorkflowStatus{}, errors.New("chat status UUID is required")
	}

	var existingStatus ChatWorkflowStatus
	if err := db.db.First(&existingStatus, "uuid = ?", status.UUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ChatWorkflowStatus{}, fmt.Errorf("chat status not found")
		}
		return ChatWorkflowStatus{}, fmt.Errorf("failed to fetch chat status: %w", err)
	}

	if status.Status != "" {
		existingStatus.Status = status.Status
		
		if status.Status == "success" {
			existingStatus.Message = "Agent workflow completed"
		} else if status.Status == "error" && status.Message == "" {
			existingStatus.Message = "An error occurred during the workflow"
		} else if status.Message == "" {
			existingStatus.Message = status.Status
		}
	}
	
	if status.Message != "" {
		existingStatus.Message = status.Message
	}
	
	existingStatus.UpdatedAt = time.Now()

	if err := db.db.Save(&existingStatus).Error; err != nil {
		return ChatWorkflowStatus{}, fmt.Errorf("failed to update chat status: %w", err)
	}

	return existingStatus, nil
}

func (db database) GetChatStatusByChatID(chatID string) ([]ChatWorkflowStatus, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	var statuses []ChatWorkflowStatus
	result := db.db.Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Find(&statuses)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch chat statuses: %w", result.Error)
	}

	return statuses, nil
}

func (db database) GetLatestChatStatusByChatID(chatID string) (ChatWorkflowStatus, error) {
	if chatID == "" {
		return ChatWorkflowStatus{}, errors.New("chat ID is required")
	}

	var status ChatWorkflowStatus
	result := db.db.Where("chat_id = ?", chatID).
		Order("created_at DESC").
		First(&status)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ChatWorkflowStatus{}, fmt.Errorf("no chat status found for chat ID: %s", chatID)
		}
		return ChatWorkflowStatus{}, fmt.Errorf("failed to fetch chat status: %w", result.Error)
	}

	return status, nil
}

func (db database) DeleteChatStatus(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("valid UUID is required")
	}

	result := db.db.Delete(&ChatWorkflowStatus{}, "uuid = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete chat status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("chat status not found with UUID: %s", id.String())
	}

	return nil
}
