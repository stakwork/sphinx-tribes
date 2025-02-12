package db

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidContent     = errors.New("content must not be empty and must be less than 10000 characters")
	ErrInvalidAuthorRef   = errors.New("author reference is required")
	ErrInvalidContentType = errors.New("invalid content type")
	ErrInvalidAuthorType  = errors.New("invalid author type")
	ErrInvalidWorkspace   = errors.New("workspace is required")
	ErrInvalidThreadID    = errors.New("thread ID is required")
)

func validateActivity(activity *Activity) error {

	if strings.TrimSpace(activity.Content) == "" || len(activity.Content) > 10000 {
		return ErrInvalidContent
	}

	if strings.TrimSpace(activity.AuthorRef) == "" {
		return ErrInvalidAuthorRef
	}

	if strings.TrimSpace(activity.Workspace) == "" {
		return ErrInvalidWorkspace
	}

	switch activity.ContentType {
	case FeatureCreation, StoryUpdate, RequirementChange, GeneralUpdate:

	default:
		return ErrInvalidContentType
	}

	switch activity.Author {
	case HumansAuthor, HiveAuthor:

	default:
		return ErrInvalidAuthorType
	}

	if activity.Author == HumansAuthor {
		if len(activity.AuthorRef) < 32 { 
			return errors.New("invalid public key format for human author")
		}
	}

	if activity.Author == HiveAuthor {
		if _, err := uuid.Parse(activity.AuthorRef); err != nil {
			return errors.New("invalid UUID format for hive author")
		}
	}

	return nil
}

func (db database) CreateActivity(activity *Activity) (*Activity, error) {

	if err := validateActivity(activity); err != nil {
		return nil, err
	}

	if activity.ID == uuid.Nil {
		activity.ID = uuid.New()
	}

	if activity.ThreadID == uuid.Nil {
		activity.ThreadID = activity.ID
		activity.Sequence = 1
	}
	
	if activity.Actions == nil {
		activity.Actions = []string{}
	}
	if activity.Questions == nil {
		activity.Questions = []string{}
	}
	
	activity.TimeCreated = time.Now()
	activity.TimeUpdated = time.Now()
	
	if activity.Status == "" {
		activity.Status = "active"
	}
	
	err := db.db.Create(activity).Error
	if err != nil {
		return nil, err
	}
	
	return activity, nil
}

func (db database) UpdateActivity(activity *Activity) (*Activity, error) {

	if err := validateActivity(activity); err != nil {
		return nil, err
	}

	existing := &Activity{}
	if err := db.db.Where("id = ?", activity.ID).First(existing).Error; err != nil {
		return nil, errors.New("activity not found")
	}

	if existing.ThreadID != activity.ThreadID {
		return nil, errors.New("thread_id cannot be modified")
	}

	if existing.Sequence != activity.Sequence {
		return nil, errors.New("sequence cannot be modified")
	}
	
	activity.TimeUpdated = time.Now()
	
	err := db.db.Model(&Activity{}).
		Where("id = ?", activity.ID).
		Updates(activity).Error
	if err != nil {
		return nil, err
	}
	
	return activity, nil
}

func (db database) GetActivity(id string) (*Activity, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid activity ID format")
	}

	var activity Activity
	err := db.db.Where("id = ?", id).First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (db database) GetActivitiesByThread(threadID string) ([]Activity, error) {

	if _, err := uuid.Parse(threadID); err != nil {
		return nil, errors.New("invalid thread ID format")
	}

	var activities []Activity
	err := db.db.Where("thread_id = ?", threadID).
		Order("sequence ASC").
		Find(&activities).Error
	return activities, err
}

func (db database) GetActivitiesByFeature(featureUUID string) ([]Activity, error) {
	if strings.TrimSpace(featureUUID) == "" {
		return nil, errors.New("feature UUID is required")
	}

	var activities []Activity
	err := db.db.Where("feature_uuid = ?", featureUUID).
		Order("time_created DESC").
		Find(&activities).Error
	return activities, err
}

func (db database) GetActivitiesByPhase(phaseUUID string) ([]Activity, error) {
	if strings.TrimSpace(phaseUUID) == "" {
		return nil, errors.New("phase UUID is required")
	}

	var activities []Activity
	err := db.db.Where("phase_uuid = ?", phaseUUID).
		Order("time_created DESC").
		Find(&activities).Error
	return activities, err
}

func (db database) GetActivitiesByWorkspace(workspace string) ([]Activity, error) {
	if strings.TrimSpace(workspace) == "" {
		return nil, errors.New("workspace is required")
	}

	var activities []Activity
	err := db.db.Where("workspace = ?", workspace).
		Order("time_created DESC").
		Find(&activities).Error
	return activities, err
}

func (db database) GetLatestActivityByThread(threadID string) (*Activity, error) {

	if _, err := uuid.Parse(threadID); err != nil {
		return nil, errors.New("invalid thread ID format")
	}

	var activity Activity
	err := db.db.Where("thread_id = ?", threadID).
		Order("sequence DESC").
		First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (db database) CreateActivityThread(sourceID string, activity *Activity) (*Activity, error) {

	if _, err := uuid.Parse(sourceID); err != nil {
		return nil, errors.New("invalid source ID format")
	}

	if err := validateActivity(activity); err != nil {
		return nil, err
	}
	
	existingActivities, err := db.GetActivitiesByThread(sourceID)
	if err != nil {
		return nil, err
	}
	
	activity.ThreadID = uuid.MustParse(sourceID)
	activity.Sequence = len(existingActivities) + 1
	
	return db.CreateActivity(activity)
}

func (db database) DeleteActivity(id string) error {
    if _, err := uuid.Parse(id); err != nil {
        return errors.New("invalid activity ID format")
    }

    existing := &Activity{}
    if err := db.db.Where("id = ?", id).First(existing).Error; err != nil {
        return errors.New("activity not found")
    }

    if err := db.db.Delete(&Activity{}, "id = ?", id).Error; err != nil {
        return err
    }

    return nil
}