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

func (db database) CreateBountyFromTicket(ticket Tickets) (*NewBounty, error) {
	now := time.Now()

	bounty := &NewBounty{
		Title:           ticket.Name,
		Description:     ticket.Description,
		PhaseUuid:       ticket.PhaseUUID,
		FeatureUuid:     ticket.FeatureUUID,
		Type:            "freelance_job_request",
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

	return bounty, nil
}
