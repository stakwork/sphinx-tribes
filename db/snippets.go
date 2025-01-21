package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (db database) CreateSnippet(snippet *TextSnippet) (*TextSnippet, error) {
	if snippet.WorkspaceUUID == "" {
		return nil, errors.New("workspace UUID is required")
	}

	if snippet.Title == "" {
		return nil, errors.New("title is required")
	}

	if snippet.Snippet == "" {
		return nil, errors.New("snippet content is required")
	}

	now := time.Now()
	snippet.DateCreated = now
	snippet.LastEdited = now

	if err := db.db.Create(snippet).Error; err != nil {
		return nil, fmt.Errorf("failed to create snippet: %w", err)
	}

	return snippet, nil
}

func (db database) GetSnippetsByWorkspace(workspaceUUID string) ([]TextSnippet, error) {
	var snippets []TextSnippet

	if err := db.db.Where("workspace_uuid = ?", workspaceUUID).Find(&snippets).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch snippets: %w", err)
	}

	return snippets, nil
}

func (db database) GetSnippetByID(id uint) (*TextSnippet, error) {
	var snippet TextSnippet

	if err := db.db.First(&snippet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to fetch snippet: %w", err)
	}

	return &snippet, nil
}

func (db database) UpdateSnippet(snippet *TextSnippet) (*TextSnippet, error) {
	if snippet.ID == 0 {
		return nil, errors.New("snippet ID is required")
	}

	snippet.LastEdited = time.Now()

	if err := db.db.Model(&TextSnippet{}).Where("id = ?", snippet.ID).Updates(snippet).Error; err != nil {
		return nil, fmt.Errorf("failed to update snippet: %w", err)
	}

	return db.GetSnippetByID(snippet.ID)
}

func (db database) DeleteSnippet(id uint) error {
	result := db.db.Delete(&TextSnippet{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete snippet: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("snippet not found")
	}

	return nil
}
