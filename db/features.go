package db

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/stakwork/sphinx-tribes/utils"
	"gorm.io/gorm"
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

func (db database) GetPhasesByFeatureUuid(featureUuid string) []FeaturePhase {
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

func (db database) CreateOrEditFeatureStory(story FeatureStory) (FeatureStory, error) {
	story.Description = strings.TrimSpace(story.Description)

	now := time.Now()
	story.Updated = &now

	existingStory := FeatureStory{}
	result := db.db.Model(&FeatureStory{}).Where("uuid = ?", story.Uuid).First(&existingStory)

	if result.RowsAffected == 0 {
		story.Created = &now
		db.db.Create(&story)
	} else {
		db.db.Model(&FeatureStory{}).Where("uuid = ?", story.Uuid).Updates(story)
	}

	db.db.Model(&FeatureStory{}).Where("uuid = ?", story.Uuid).Find(&story)

	return story, nil
}

func (db database) GetFeatureStoriesByFeatureUuid(featureUuid string) ([]FeatureStory, error) {
	var stories []FeatureStory
	result := db.db.Where("feature_uuid = ?", featureUuid).Find(&stories)
	if result.Error != nil {
		return nil, result.Error
	}

	for i := range stories {
		stories[i].Description = strings.TrimSpace(stories[i].Description)
	}
	return stories, nil
}

func (db database) GetFeatureStoryByUuid(featureUuid, storyUuid string) (FeatureStory, error) {
	story := FeatureStory{}
	result := db.db.Model(&FeatureStory{}).Where("feature_uuid = ? AND uuid = ?", featureUuid, storyUuid).First(&story)
	if result.RowsAffected == 0 {
		return story, errors.New("no story found")
	}
	return story, nil
}

func (db database) DeleteFeatureStoryByUuid(featureUuid, storyUuid string) error {
	result := db.db.Where("feature_uuid = ? AND uuid = ?", featureUuid, storyUuid).Delete(&FeatureStory{})
	if result.RowsAffected == 0 {
		return errors.New("no story found to delete")
	}
	return nil
}

func (db *database) GetFeaturePhaseBounty(featureUUID string, phaseUUID string) (*Bounty, error) {
	var bounty Bounty
	result := db.db.Where("feature_uuid = ? AND phase_uuid = ?", featureUUID, phaseUUID).First(&bounty)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &bounty, nil
}
