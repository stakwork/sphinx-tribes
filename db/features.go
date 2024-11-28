package db

import (
	"errors"
	"fmt"
	"net/http"
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
		orderQuery = "ORDER BY priority ASC"
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
	db.db.Model(&FeaturePhase{}).Where("feature_uuid = ?", featureUuid).Order("created ASC").Find(&phases)
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
	result := db.db.Where("feature_uuid = ?", featureUuid).Order("priority ASC").Find(&stories)
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

func (db database) GetBountiesByFeatureAndPhaseUuid(featureUuid string, phaseUuid string, r *http.Request) ([]NewBounty, error) {
	keys := r.URL.Query()
	tags := keys.Get("tags")
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)
	open := keys.Get("Open")
	assigned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")

	var bounties []NewBounty

	// Initialize the query with the necessary joins and initial filters
	query := db.db.Model(&Bounty{}).
		Select("bounty.*").
		Joins(`INNER JOIN "feature_phases" ON "feature_phases"."uuid" = "bounty"."phase_uuid"`).
		Where(`"feature_phases"."feature_uuid" = ? AND "feature_phases"."uuid" = ?`, featureUuid, phaseUuid)

	// Add pagination if applicable
	if limit > 1 {
		query = query.Limit(limit).Offset(offset)
	}

	// Add sorting if applicable
	if sortBy != "" && direction != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, direction))
	} else {
		query = query.Order("created_at DESC")
	}

	// Add search filter
	if search != "" {
		searchQuery := fmt.Sprintf("LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
		query = query.Where(searchQuery)
	}

	// Add language filter
	if len(languageArray) > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				query = query.Where("coding_languages && ARRAY[" + langs + "]")
			}
		}
	}

	// Add status filters
	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true AND completed != true")
	}
	if assigned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false AND completed = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}

	if len(statusConditions) > 0 {
		query = query.Where(strings.Join(statusConditions, " OR "))
	}

	// Execute the query
	result := query.Find(&bounties)

	if result.RowsAffected == 0 {
		return bounties, errors.New("no bounty found")
	}

	// Handle tags if any
	if tags != "" {
		// pull out the tags and add them in here
		t := strings.Split(tags, ",")
		for _, s := range t {
			query = query.Where("'" + s + "'" + " = any (tags)")
		}
		query.Scan(&bounties)
	}

	return bounties, nil
}

func (db database) GetBountiesCountByFeatureAndPhaseUuid(featureUuid string, phaseUuid string, r *http.Request) int64 {
	keys := r.URL.Query()
	open := keys.Get("Open")
	assigned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")

	// Initialize the query with the necessary joins and initial filters
	query := db.db.Model(&Bounty{}).
		Select("COUNT(*)").
		Joins(`INNER JOIN "feature_phases" ON "feature_phases"."uuid" = "bounty"."phase_uuid"`).
		Where(`"feature_phases"."feature_uuid" = ? AND "feature_phases"."uuid" = ?`, featureUuid, phaseUuid)

	// Add status filters
	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true AND completed != true")
	}
	if assigned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false AND completed = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}

	if len(statusConditions) > 0 {
		query = query.Where(strings.Join(statusConditions, " OR "))
	}

	var count int64

	query.Count(&count)

	return count
}

func (db database) GetPhaseByUuid(phaseUuid string) (FeaturePhase, error) {
	phase := FeaturePhase{}
	result := db.db.Model(&FeaturePhase{}).Where("uuid = ?", phaseUuid).First(&phase)
	if result.RowsAffected == 0 {
		return phase, errors.New("no phase found")
	}
	return phase, nil
}

func (db database) GetBountiesByPhaseUuid(phaseUuid string) []Bounty {
	bounties := []Bounty{}
	db.db.Model(&Bounty{}).Where("phase_uuid = ?", phaseUuid).Find(&bounties)
	return bounties
}

func (db database) GetTicketsByPhase(featureUuid string, phaseUuid string) ([]Tickets, error) {
	var tickets []Tickets

	result := db.db.Where("feature_uuid = ? AND phase_uuid = ?", featureUuid, phaseUuid).
		Order("sequence ASC").
		Find(&tickets)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch tickets for phase: %w", result.Error)
	}

	return tickets, nil
}

func (db database) GetProductBrief(workspaceUuid string) (string, error) {
	workspace := Workspace{}
	result := db.db.Model(&Workspace{}).Where("uuid = ?", workspaceUuid).First(&workspace)
	if result.Error != nil {
		return "", fmt.Errorf("error getting workspace: %v", result.Error)
	}

	productBrief := fmt.Sprintf("Product: %s. Product Brief:\n Mission: %s.\n\n Objectives: %s",
		workspace.Name,
		workspace.Mission,
		workspace.Tactics)

	return productBrief, nil
}

func (db database) GetFeatureBrief(featureUuid string) (string, error) {
	feature := WorkspaceFeatures{}
	result := db.db.Model(&WorkspaceFeatures{}).Where("uuid = ?", featureUuid).First(&feature)
	if result.Error != nil {
		return "", fmt.Errorf("error getting feature: %v", result.Error)
	}

	featureBrief := fmt.Sprintf("Feature: %s. Brief: %s",
		feature.Name,
		feature.Brief)

	return featureBrief, nil
}
