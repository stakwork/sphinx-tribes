package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (db database) GetTicketsByPhase(phaseUUID string) ([]Tickets, error) {
	var tickets []Tickets

	result := db.db.Where("phase_uuid = ?", phaseUUID).
		Order("sequence ASC").
		Find(&tickets)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch tickets for phase: %w", result.Error)
	}

	// Return empty slice if no tickets found
	if result.RowsAffected == 0 {
		return []Tickets{}, nil
	}

	return tickets, nil
}
