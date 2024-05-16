package db

import (
	"errors"
	"fmt"
	"github.com/stakwork/sphinx-tribes/utils"
	"net/http"
	"strings"
	"time"
)

func (db database) GetFeaturesByWorkspaceUuid(uuid string, r *http.Request) []WorkspaceFeatures {
	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(r)

	orderQuery := ""
	limitQuery := ""

	ms := []WorkspaceFeatures{}

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = "ORDER BY created DESC"
	}

	if limit > 1 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	query := `SELECT * FROM public.workspace_features WHERE workspace_uuid = '` + uuid + `'`

	allQuery := query + " " + orderQuery + " " + limitQuery

	theQuery := db.db.Raw(allQuery)

	theQuery.Scan(&ms)

	return ms
}

func (db database) GetWorkspaceFeaturesCount(uuid string) int64 {
	var count int64
	db.db.Model(&WorkspaceFeatures{}).Where("workspace_uuid = ?", uuid).Count(&count)
	return count
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

	var existing WorkspaceFeatures
	result := db.db.Model(&WorkspaceFeatures{}).Where("uuid = ?", m.Uuid).First(&existing)
	if result.RowsAffected == 0 {

		m.Created = &now
		db.db.Create(&m)
	} else {

		db.db.Model(&WorkspaceFeatures{}).Where("uuid = ?", m.Uuid).Updates(m)
	}

	db.db.Model(&WorkspaceFeatures{}).Where("uuid = ?", m.Uuid).First(&m)
	return m, nil
}

func (db database) DeleteFeatureByUuid(uuid string) error {
	result := db.db.Where("uuid = ?", uuid).Delete(&WorkspaceFeatures{})

	if result.RowsAffected == 0 {
		return errors.New("no feature found to delete")
	}
	return nil

}
