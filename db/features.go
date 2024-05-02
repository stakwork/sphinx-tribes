package db

import (
	"strings"
	"time"
)

func (db database) GetFeaturesByWorkspaceUuid(uuid string) []WorkspaceFeatures {
	ms := []WorkspaceFeatures{}

	db.db.Model(&WorkspaceFeatures{}).Where("workspace_uuid = ?", uuid).Order("Created").Find(&ms)

	return ms
}

func (db database) GetFeatureByUuid(uuid string) WorkspaceFeatures {
	ms := WorkspaceFeatures{}

	db.db.Model(&WorkspaceFeatures{}).Where("uuid = ?", uuid).Find(&ms)

	return ms
}

func (db database) CreateOrEditFeature(m WorkspaceFeatures) (WorkspaceFeatures, error) {
	m.Name = strings.TrimSpace(m.Name)
	m.Brief = strings.TrimSpace(m.Brief)
	m.Requirements = strings.TrimSpace(m.Requirements)
	m.Architecture = strings.TrimSpace(m.Architecture)

	now := time.Now()
	m.Updated = &now

	if db.db.Model(&m).Where("uuid = ?", m.Uuid).Updates(&m).RowsAffected == 0 {
		m.Created = &now
		db.db.Create(&m)
	}

	db.db.Model(&WorkspaceFeatures{}).Where("uuid = ?", m.Uuid).Find(&m)

	return m, nil
}
