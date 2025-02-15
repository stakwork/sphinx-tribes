package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) CreateOrEditTicketPlan(plan *TicketPlan) (*TicketPlan, error) {
	if plan.UUID == uuid.Nil {
		return nil, errors.New("ticket plan UUID is required")
	}

	if plan.WorkspaceUuid == "" {
		return nil, errors.New("workspace UUID is required")
	}

	if plan.Name == "" {
		return nil, errors.New("name is required")
	}

	var existingPlan TicketPlan
	result := db.db.Where("uuid = ?", plan.UUID).First(&existingPlan)

	now := time.Now()
	if result.Error != nil {

		plan.CreatedAt = now
		plan.UpdatedAt = now
		plan.Version = 1
		if err := db.db.Create(plan).Error; err != nil {
			return nil, fmt.Errorf("failed to create ticket plan: %w", err)
		}
	} else {
		plan.UpdatedAt = now
		plan.Version = existingPlan.Version + 1
		if err := db.db.Model(&existingPlan).Updates(plan).Error; err != nil {
			return nil, fmt.Errorf("failed to update ticket plan: %w", err)
		}
	}

	var updatedPlan TicketPlan
	if err := db.db.Where("uuid = ?", plan.UUID).First(&updatedPlan).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated ticket plan: %w", err)
	}

	return &updatedPlan, nil
}

func (db database) GetTicketPlan(uuid string) (*TicketPlan, error) {
	var plan TicketPlan
	result := db.db.Where("uuid = ?", uuid).First(&plan)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ticket plan not found")
		}
		return nil, result.Error
	}
	
	return &plan, nil
}

func (db database) DeleteTicketPlan(uuid string) error {
	result := db.db.Where("uuid = ?", uuid).Delete(&TicketPlan{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete ticket plan: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket plan not found")
	}
	return nil
}

func (db database) GetTicketPlansByFeature(featureUUID string) ([]TicketPlan, error) {
	var plans []TicketPlan
	if err := db.db.Where("feature_uuid = ?", featureUUID).Find(&plans).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ticket plans by feature: %w", err)
	}
	return plans, nil
}

func (db database) GetTicketPlansByPhase(phaseUUID string) ([]TicketPlan, error) {
	var plans []TicketPlan
	if err := db.db.Where("phase_uuid = ?", phaseUUID).Find(&plans).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ticket plans by phase: %w", err)
	}
	return plans, nil
}

func (db database) GetTicketPlansByWorkspace(workspaceUUID string) ([]TicketPlan, error) {
	var plans []TicketPlan
	if err := db.db.Where("workspace_uuid = ?", workspaceUUID).Find(&plans).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ticket plans by workspace: %w", err)
	}
	return plans, nil
} 

func (db database) BuildTicketArray(groupIDs []string) []TicketArrayItem {
    var ticketArray []TicketArrayItem

    for _, groupID := range groupIDs {
        var tickets []Tickets
        result := db.db.Where("ticket_group = ?", groupID).Find(&tickets)
        if result.Error != nil {
            continue
        }

        var latestVersion int
        var latestName, latestDescription string

        for _, ticket := range tickets {
            if ticket.Version > latestVersion {
                latestVersion = ticket.Version
                latestName = ticket.Name
                latestDescription = ticket.Description
            }
        }

        if latestVersion > 0 {
            ticketArray = append(ticketArray, TicketArrayItem{
                TicketName:        latestName,
                TicketDescription: latestDescription,
            })
        }
    }

    return ticketArray
}
