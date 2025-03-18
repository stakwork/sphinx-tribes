package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) CreateArtifact(artifact *Artifact) (*Artifact, error) {

	if artifact.MessageID == "" {
		return nil, fmt.Errorf("message ID cannot be empty")
	}

	var count int64
	if err := db.db.Model(&ChatMessage{}).Where("id = ?", artifact.MessageID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check message existence: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("message with ID %s does not exist", artifact.MessageID)
	}

	validTypes := map[ArtifactType]bool{
		TextArtifact:   true,
		VisualArtifact: true,
		ActionArtifact: true,
		SSEArtifact:    true,
	}
	if !validTypes[artifact.Type] {
		return nil, fmt.Errorf("invalid artifact type: %s", artifact.Type)
	}

	if artifact.ID == uuid.Nil {
		artifact.ID = uuid.New()
	}

	now := time.Now()
	artifact.CreatedAt = now
	artifact.UpdatedAt = now

	if err := db.db.Create(artifact).Error; err != nil {
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}

	return artifact, nil
}

func (db database) GetArtifactByID(id uuid.UUID) (*Artifact, error) {
	var artifact Artifact
	if err := db.db.First(&artifact, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch artifact: %w", err)
	}

	return &artifact, nil
}

func (db database) GetArtifactsByMessageID(messageID string) ([]Artifact, error) {
	var artifacts []Artifact

	if err := db.db.Where("message_id = ?", messageID).
		Order("created_at DESC").
		Find(&artifacts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch artifacts: %w", err)
	}

	return artifacts, nil
}

func (db database) GetAllArtifactsByChatID(chatID string) ([]Artifact, error) {
	var artifacts []Artifact

	if err := db.db.
		Joins("JOIN chat_messages ON artifacts.message_id = chat_messages.id").
		Where("chat_messages.chat_id = ?", chatID).
		Order("artifacts.created_at DESC").
		Find(&artifacts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch artifacts: %w", err)
	}

	return artifacts, nil
}

func (db database) UpdateArtifact(artifact *Artifact) (*Artifact, error) {
	if artifact.ID == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	var existingArtifact Artifact
	if err := db.db.First(&existingArtifact, "id = ?", artifact.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("artifact not found")
		}
		return nil, fmt.Errorf("failed to fetch artifact: %w", err)
	}

	artifact.MessageID = existingArtifact.MessageID
	artifact.Type = existingArtifact.Type
	artifact.UpdatedAt = time.Now()

	if err := db.db.Save(artifact).Error; err != nil {
		return nil, fmt.Errorf("failed to update artifact: %w", err)
	}

	return artifact, nil
}

func (db database) DeleteArtifactByID(id uuid.UUID) error {
	result := db.db.Delete(&Artifact{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete artifact: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("artifact not found")
	}

	return nil
}

func (db database) DeleteAllArtifactsByChatID(chatID string) error {
	result := db.db.Exec(`
		DELETE FROM artifacts 
		WHERE message_id IN (
			SELECT id FROM chat_messages WHERE chat_id = ?
		)
	`, chatID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete artifacts: %w", result.Error)
	}

	return nil
}
