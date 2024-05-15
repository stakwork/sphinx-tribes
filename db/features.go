package db

import (
	"errors"
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
