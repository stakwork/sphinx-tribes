package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

func (db database) CreateOrUpdateFeatureCall(workspaceID string, url string) (*FeatureCall, error) {
	if workspaceID == "" {
		return nil, errors.New("workspace_id is required")
	}
	if url == "" {
		return nil, errors.New("url is required")
	}

	workspace := Workspace{}
	if result := db.db.Where("uuid = ?", workspaceID).First(&workspace); result.Error != nil {
		return nil, errors.New("workspace not found")
	}

	featureCall := &FeatureCall{}
	result := db.db.Where("workspace_id = ?", workspaceID).First(featureCall)

	if result.Error != nil {
		featureCall = &FeatureCall{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			URL:         url,
		}
		if err := db.db.Create(featureCall).Error; err != nil {
			return nil, err
		}
	} else {
		featureCall.URL = url
		featureCall.UpdatedAt = time.Now()
		if err := db.db.Save(featureCall).Error; err != nil {
			return nil, err
		}
	}

	return featureCall, nil
}

func (db database) GetFeatureCallByWorkspaceID(workspaceID string) (*FeatureCall, error) {
	if workspaceID == "" {
		return nil, errors.New("workspace_id is required")
	}

	featureCall := &FeatureCall{}
	result := db.db.Where("workspace_id = ?", workspaceID).First(featureCall)
	if result.Error != nil {
		return nil, result.Error
	}

	return featureCall, nil
}

func (db database) DeleteFeatureCall(workspaceID string) error {
	if workspaceID == "" {
		return errors.New("workspace_id is required")
	}

	result := db.db.Where("workspace_id = ?", workspaceID).Delete(&FeatureCall{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("feature call not found")
	}

	return nil
} 