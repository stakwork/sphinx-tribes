package db

import (
	"errors"
	"time"
)

func (db database) GetCodeGraphByUUID(uuid string) (WorkspaceCodeGraph, error) {
	var codeGraph WorkspaceCodeGraph
	result := db.db.Where("uuid = ?", uuid).First(&codeGraph)

	if result.Error != nil {
		return WorkspaceCodeGraph{}, result.Error
	}

	return codeGraph, nil
}

func (db database) GetCodeGraphsByWorkspaceUuid(workspace_uuid string) ([]WorkspaceCodeGraph, error) {
	var codeGraphs []WorkspaceCodeGraph
	result := db.db.Where("workspace_uuid = ?", workspace_uuid).Find(&codeGraphs)

	if result.Error != nil {
		return nil, result.Error
	}

	return codeGraphs, nil
}

func (db database) CreateOrEditCodeGraph(m WorkspaceCodeGraph) (WorkspaceCodeGraph, error) {
	if m.Uuid == "" {
		return WorkspaceCodeGraph{}, errors.New("uuid is required")
	}

	var existing WorkspaceCodeGraph
	result := db.db.Where("uuid = ?", m.Uuid).First(&existing)

	now := time.Now()
	if result.Error != nil {

		m.Created = &now
		m.Updated = &now
		if err := db.db.Create(&m).Error; err != nil {
			return WorkspaceCodeGraph{}, err
		}
		return m, nil
	}

	m.Created = existing.Created
	m.Updated = &now
	if err := db.db.Model(&existing).Updates(m).Error; err != nil {
		return WorkspaceCodeGraph{}, err
	}

	var updated WorkspaceCodeGraph
	if err := db.db.Where("uuid = ?", m.Uuid).First(&updated).Error; err != nil {
		return WorkspaceCodeGraph{}, err
	}

	return updated, nil
}

func (db database) DeleteCodeGraph(workspace_uuid string, uuid string) error {
	if uuid == "" || workspace_uuid == "" {
		return errors.New("workspace_uuid and uuid are required")
	}

	result := db.db.Where("workspace_uuid = ? AND uuid = ?", workspace_uuid, uuid).Delete(&WorkspaceCodeGraph{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("code graph not found")
	}

	return nil
}
