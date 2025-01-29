package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/logger"
	"gorm.io/gorm"
)

func (db database) CreateOrEditTicket(ticket *Tickets) (Tickets, error) {

	if ticket.UUID == uuid.Nil {
		return Tickets{}, errors.New("ticket UUID is required")
	}

	if ticket.FeatureUUID == "" {
		return Tickets{}, errors.New("feature UUID is required")
	}

	if ticket.Status != "" && !IsValidTicketStatus(ticket.Status) {
		return Tickets{}, errors.New("invalid ticket status")
	}

	var existingTicket Tickets
	result := db.db.Where("uuid = ?", ticket.UUID).First(&existingTicket)

	now := time.Now()
	ticket.UpdatedAt = now

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ticket.CreatedAt = now

		if ticket.Status == "" {
			ticket.Status = DraftTicket
		}

		if err := db.db.Create(&ticket).Error; err != nil {
			return Tickets{}, fmt.Errorf("failed to create ticket: %w", err)
		}
		return *ticket, nil
	}

	if result.Error != nil {
		return Tickets{}, fmt.Errorf("database error: %w", result.Error)
	}

	if err := db.db.Model(&existingTicket).Updates(ticket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to update ticket: %w", err)
	}

	var updatedTicket Tickets
	if err := db.db.Where("uuid = ?", ticket.UUID).First(&updatedTicket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to fetch updated ticket: %w", err)
	}

	return updatedTicket, nil
}

func (db database) GetTicket(uuid string) (Tickets, error) {
	ticket := Tickets{}

	results := db.db.Model(&Tickets{}).Where("uuid = ?", uuid).Find(&ticket)

	if results.Error != nil {
		return Tickets{}, fmt.Errorf("failed to get ticket: %w", results.Error)
	}

	if results.RowsAffected == 0 {
		return Tickets{}, fmt.Errorf("ticket not found")
	}

	return ticket, nil
}

func IsValidTicketStatus(status TicketStatus) bool {
	switch status {
	case DraftTicket, ReadyTicket, InProgressTicket, TestTicket, DeployTicket, PayTicket, CompletedTicket:
		return true
	default:
		return false
	}
}

func (db database) UpdateTicket(ticket Tickets) (Tickets, error) {
	if ticket.UUID == uuid.Nil {
		return Tickets{}, errors.New("ticket UUID is required")
	}

	if ticket.Status != "" && !IsValidTicketStatus(ticket.Status) {
		return Tickets{}, errors.New("invalid ticket status")
	}

	var existingTicket Tickets
	result := db.db.Where("uuid = ?", ticket.UUID).First(&existingTicket)

	now := time.Now()
	ticket.UpdatedAt = now

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ticket.CreatedAt = now

			if ticket.Status == "" {
				ticket.Status = DraftTicket
			}
			if err := db.db.Create(&ticket).Error; err != nil {
				return Tickets{}, fmt.Errorf("failed to create ticket: %w", err)
			}
			return ticket, nil
		}
		return Tickets{}, fmt.Errorf("database error: %w", result.Error)
	}

	if err := db.db.Model(&existingTicket).Updates(ticket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to update ticket: %w", err)
	}

	var updatedTicket Tickets
	if err := db.db.Where("uuid = ?", ticket.UUID).First(&updatedTicket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to fetch updated ticket: %w", err)
	}

	return updatedTicket, nil
}

func (db database) GetTicketsByGroup(ticketGroupUUID string) ([]Tickets, error) {
	var tickets []Tickets

	result := db.db.Model(&Tickets{}).Where("ticket_group = ?", ticketGroupUUID).Find(&tickets)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch tickets by group: %w", result.Error)
	}

	return tickets, nil
}

func (db database) DeleteTicket(uuid string) error {
	result := db.db.Where("uuid = ?", uuid).Delete(&Tickets{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete ticket: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}
	return nil
}

func (db database) GetTicketsByPhaseUUID(featureUUID string, phaseUUID string) ([]Tickets, error) {
	var tickets []Tickets

	result := db.db.
		Where("feature_uuid = ? AND phase_uuid = ?", featureUUID, phaseUUID).
		Order("sequence ASC").
		Find(&tickets)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch tickets: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return []Tickets{}, nil
	}

	return tickets, nil
}

func (db database) GetTicketsWithoutGroup() ([]Tickets, error) {
	var tickets []Tickets

	result := db.db.
		Where("ticket_group IS NULL OR ticket_group = ?", uuid.Nil).Find(&tickets)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch tickets: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return []Tickets{}, nil
	}

	return tickets, nil
}

func (db database) UpdateTicketsWithoutGroup(ticket Tickets) error {
	data := map[string]interface{}{}

	data["ticket_group"] = ticket.UUID

	if ticket.AuthorID == nil {
		data["author_id"] = "12345"
	}

	if ticket.Author == nil {
		data["author"] = "HUMAN"
	}

	logger.Log.Info("data === %v", data)

	result := db.db.Model(&Tickets{}).Where("uuid = ?", ticket.UUID).Updates(data)

	if result.Error != nil {
		return fmt.Errorf("failed to update ticket: %w", result.Error)
	}

	return nil
}

func (db database) CreateBountyFromTicket(ticket Tickets, pubkey string) (*NewBounty, error) {
	now := time.Now()

	feature := db.GetFeatureByUuid(ticket.FeatureUUID)

	bounty := &NewBounty{
		Title:           ticket.Name,
		Description:     ticket.Description,
		PhaseUuid:       ticket.PhaseUUID,
		FeatureUuid:     ticket.FeatureUUID,
		WorkspaceUuid:   feature.WorkspaceUuid,
		OwnerID:         pubkey,
		Type:            "freelance_job_request",
		WantedType:      "Other",
		Price:           21,
		Created:         now.Unix(),
		Updated:         &now,
		Show:            true,
		CodingLanguages: pq.StringArray{},
	}

	if err := db.db.Create(bounty).Error; err != nil {
		logger.Log.Error("failed to create bounty", "error", err, "ticket_id", ticket.UUID)
		return nil, fmt.Errorf("failed to create bounty: %w", err)
	}
	//tst
	return bounty, nil
}

func (db database) GetLatestTicketByGroup(ticketGroup uuid.UUID) (Tickets, error) {
	var ticket Tickets
	result := db.db.Where("ticket_group = ?", ticketGroup).
		Order("version DESC").
		Limit(1).
		First(&ticket)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Tickets{}, fmt.Errorf("no tickets found for group %s", ticketGroup)
		}
		return Tickets{}, fmt.Errorf("failed to fetch latest ticket: %w", result.Error)
	}

	return ticket, nil
}

func (db database) GetAllTicketGroups(workspaceUuid string) ([]uuid.UUID, error) {
	var groups []uuid.UUID
	result := db.db.Model(&Tickets{}).
		Joins("JOIN workspace_features ON tickets.feature_uuid = workspace_features.uuid").
		Where("workspace_features.workspace_uuid = ? AND tickets.ticket_group IS NOT NULL", workspaceUuid).
		Select("DISTINCT tickets.ticket_group").
		Find(&groups)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch ticket groups: %w", result.Error)
	}

	return groups, nil
}

func (db database) GetWorkspaceDraftTicket(workspaceUuid string, uuid string) (Tickets, error) {
	var ticket Tickets

	result := db.db.
		Where("workspace_uuid = ? AND uuid = ? AND feature_uuid IS NULL AND phase_uuid IS NULL",
			workspaceUuid, uuid).
		First(&ticket)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Tickets{}, fmt.Errorf("draft ticket not found")
		}
		return Tickets{}, fmt.Errorf("failed to fetch draft ticket: %w", result.Error)
	}

	return ticket, nil
}

func (db database) CreateWorkspaceDraftTicket(ticket *Tickets) (Tickets, error) {
	if ticket.UUID == uuid.Nil {
		return Tickets{}, errors.New("ticket UUID is required")
	}

	if ticket.WorkspaceUuid == "" {
		return Tickets{}, errors.New("workspace UUID is required")
	}

	now := time.Now()
	ticket.CreatedAt = now
	ticket.UpdatedAt = now
	ticket.Status = DraftTicket
	ticket.Version = 1

	if err := db.db.Omit("Features", "FeaturePhase").Create(ticket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to create draft ticket: %w", err)
	}

	var createdTicket Tickets
	if err := db.db.Where("uuid = ?", ticket.UUID).First(&createdTicket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to fetch created ticket: %w", err)
	}

	return createdTicket, nil
}

func (db database) UpdateWorkspaceDraftTicket(ticket *Tickets) (Tickets, error) {
	var existingTicket Tickets
	result := db.db.Where("uuid = ? AND workspace_uuid = ?",
		ticket.UUID, ticket.WorkspaceUuid).First(&existingTicket)

	if result.Error != nil {
		return Tickets{}, fmt.Errorf("failed to find draft ticket: %w", result.Error)
	}

	ticket.UpdatedAt = time.Now()
	ticket.Version = existingTicket.Version + 1

	if err := db.db.Model(&existingTicket).
		Omit("Features", "FeaturePhase").
		Updates(map[string]interface{}{
			"name":        ticket.Name,
			"description": ticket.Description,
			"status":      ticket.Status,
			"updated_at":  ticket.UpdatedAt,
			"version":     ticket.Version,
		}).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to update draft ticket: %w", err)
	}

	var updatedTicket Tickets
	if err := db.db.Where("uuid = ?", ticket.UUID).First(&updatedTicket).Error; err != nil {
		return Tickets{}, fmt.Errorf("failed to fetch updated ticket: %w", err)
	}

	return updatedTicket, nil
}

func (db database) DeleteWorkspaceDraftTicket(workspaceUuid string, uuid string) error {
	result := db.db.Where("workspace_uuid = ? AND uuid = ? AND feature_uuid IS NULL AND phase_uuid IS NULL",
		workspaceUuid, uuid).Delete(&Tickets{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete draft ticket: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("draft ticket not found")
	}

	return nil
}
