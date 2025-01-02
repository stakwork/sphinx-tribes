// db/feature_flags.go
package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) AddFeatureFlag(flag *FeatureFlag) (FeatureFlag, error) {
	if flag.UUID == uuid.Nil {
		return FeatureFlag{}, errors.New("feature flag UUID is required")
	}

	now := time.Now()
	flag.CreatedAt = now
	flag.UpdatedAt = now

	if err := db.db.Create(&flag).Error; err != nil {
		return FeatureFlag{}, fmt.Errorf("failed to create feature flag: %w", err)
	}

	return *flag, nil
}

func (db database) UpdateFeatureFlag(flag *FeatureFlag) (FeatureFlag, error) {
	if flag.UUID == uuid.Nil {
		return FeatureFlag{}, errors.New("feature flag UUID is required")
	}

	var existingFlag FeatureFlag
	if err := db.db.First(&existingFlag, "uuid = ?", flag.UUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return FeatureFlag{}, fmt.Errorf("feature flag not found")
		}
		return FeatureFlag{}, fmt.Errorf("failed to fetch feature flag: %w", err)
	}

	existingFlag.Name = flag.Name
	existingFlag.Description = flag.Description
	existingFlag.Enabled = flag.Enabled
	existingFlag.UpdatedAt = time.Now()

	if err := db.db.Save(&existingFlag).Error; err != nil {
		return FeatureFlag{}, fmt.Errorf("failed to update feature flag: %w", err)
	}

	return existingFlag, nil
}

func (db database) DeleteFeatureFlag(flagUUID uuid.UUID) error {
	tx := db.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	if err := tx.Delete(&Endpoint{}, "feature_flag_uuid = ?", flagUUID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete associated endpoints: %w", err)
	}

	result := tx.Delete(&FeatureFlag{}, "uuid = ?", flagUUID)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete feature flag: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("feature flag not found")
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db database) GetFeatureFlags() ([]FeatureFlag, error) {
	var flags []FeatureFlag
	if err := db.db.Preload("Endpoints").Find(&flags).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch feature flags: %w", err)
	}
	return flags, nil
}

func (db database) GetFeatureFlagByUUID(flagUUID uuid.UUID) (FeatureFlag, error) {
	var flag FeatureFlag
	if err := db.db.Preload("Endpoints").First(&flag, "uuid = ?", flagUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return FeatureFlag{}, fmt.Errorf("feature flag not found")
		}
		return FeatureFlag{}, fmt.Errorf("failed to fetch feature flag: %w", err)
	}
	return flag, nil
}

func (db database) AddEndpoint(endpoint *Endpoint) (Endpoint, error) {
	if endpoint.UUID == uuid.Nil {
		return Endpoint{}, errors.New("endpoint UUID is required")
	}

	now := time.Now()
	endpoint.CreatedAt = now
	endpoint.UpdatedAt = now

	if err := db.db.Create(&endpoint).Error; err != nil {
		return Endpoint{}, fmt.Errorf("failed to create endpoint: %w", err)
	}

	return *endpoint, nil
}

func (db database) UpdateEndpoint(endpoint *Endpoint) (Endpoint, error) {
	if endpoint.UUID == uuid.Nil {
		return Endpoint{}, errors.New("endpoint UUID is required")
	}

	if endpoint.Path == "" {
		return Endpoint{}, errors.New("endpoint path is required")
	}

	var existingEndpoint Endpoint
	if err := db.db.First(&existingEndpoint, "uuid = ?", endpoint.UUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Endpoint{}, fmt.Errorf("endpoint not found")
		}
		return Endpoint{}, fmt.Errorf("failed to fetch endpoint: %w", err)
	}

	existingEndpoint.Path = endpoint.Path
	existingEndpoint.UpdatedAt = time.Now()

	if err := db.db.Save(&existingEndpoint).Error; err != nil {
		return Endpoint{}, fmt.Errorf("failed to update endpoint: %w", err)
	}

	return existingEndpoint, nil
}

func (db database) DeleteEndpoint(endpointUUID uuid.UUID) error {
	result := db.db.Delete(&Endpoint{}, "uuid = ?", endpointUUID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete endpoint: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("endpoint not found")
	}
	return nil
}

func (db database) GetEndpointsByFeatureFlag(flagUUID uuid.UUID) ([]Endpoint, error) {
	var endpoints []Endpoint
	if err := db.db.Where("feature_flag_uuid = ?", flagUUID).Find(&endpoints).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch endpoints: %w", err)
	}
	return endpoints, nil
}

func (db database) GetEndpointByPath(path string) (Endpoint, error) {
	var endpoint Endpoint
	if err := db.db.First(&endpoint, "path = ?", path).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Endpoint{}, fmt.Errorf("endpoint not found")
		}
		return Endpoint{}, fmt.Errorf("failed to fetch endpoint: %w", err)
	}
	return endpoint, nil
}

func (db database) GetEndpointByUUID(uuid uuid.UUID) (Endpoint, error) {
	var endpoint Endpoint

	if err := db.db.Where("uuid = ?", uuid).First(&endpoint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Endpoint{}, fmt.Errorf("endpoint not found")
		}
		return Endpoint{}, fmt.Errorf("failed to fetch endpoint: %w", err)
	}

	return endpoint, nil
}
