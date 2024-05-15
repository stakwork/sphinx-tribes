package db

import (
<<<<<<< HEAD
	"fmt"
	"net/http"
=======
	"errors"
>>>>>>> b77bffb1 (tests for Add Phases to Features)
	"strings"
	"time"

	"github.com/stakwork/sphinx-tribes/utils"
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

	if db.db.Model(&m).Where("uuid = ?", m.Uuid).Updates(&m).RowsAffected == 0 {
		m.Created = &now
		db.db.Create(&m)
	}

	return m, nil
}

func (db database) CreateOrEditFeaturePhase(phase FeaturePhase) (FeaturePhase, error) {
	phase.Name = strings.TrimSpace(phase.Name)

	now := time.Now()
	phase.Updated = &now

	existingPhase := FeaturePhase{}
	result := db.db.Model(&FeaturePhase{}).Where("uuid = ?", phase.Uuid).First(&existingPhase)

	if result.RowsAffected == 0 {

		phase.Created = &now
		db.db.Create(&phase)
	} else {

		db.db.Model(&FeaturePhase{}).Where("uuid = ?", phase.Uuid).Updates(phase)
	}

	db.db.Model(&FeaturePhase{}).Where("uuid = ?", phase.Uuid).Find(&phase)

	return phase, nil
}

func (db database) GetFeaturePhasesByFeatureUuid(featureUuid string) []FeaturePhase {
	phases := []FeaturePhase{}
	db.db.Model(&FeaturePhase{}).Where("feature_uuid = ?", featureUuid).Order("Created ASC").Find(&phases)
	return phases
}

func (db database) GetFeaturePhaseByUuid(featureUuid, phaseUuid string) (FeaturePhase, error) {
	phase := FeaturePhase{}
	result := db.db.Model(&FeaturePhase{}).Where("feature_uuid = ? AND uuid = ?", featureUuid, phaseUuid).First(&phase)
	if result.RowsAffected == 0 {
		return phase, errors.New("no phase found")
	}
	return phase, nil
}

func (db database) DeleteFeaturePhase(featureUuid, phaseUuid string) error {
	result := db.db.Where("feature_uuid = ? AND uuid = ?", featureUuid, phaseUuid).Delete(&FeaturePhase{})
	if result.RowsAffected == 0 {
		return errors.New("no phase found to delete")
	}
	return nil
}
