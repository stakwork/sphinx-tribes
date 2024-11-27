package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) CreateOrEditTicket(ticket *Tickets) (Tickets, error) {

	if ticket.UUID == uuid.Nil || ticket.FeatureUUID == "" || ticket.PhaseUUID == "" || ticket.Name == "" {
		return Tickets{}, errors.New("required fields are missing")
	}

	// check if ticket exists and update it
	if db.db.Model(&Tickets{}).Where("uuid = ?", ticket.UUID).First(&ticket).RowsAffected != 0 {
		now := time.Now()
		ticket.UpdatedAt = now

		// update ticket
		if db.db.Model(&ticket).Where("uuid = ?", ticket.UUID).Updates(&ticket).RowsAffected == 0 {
			return Tickets{}, errors.New("failed to update ticket")
		}

		return *ticket, nil
	}

	// create ticket and return error if it fails
	if db.db.Create(&ticket).Error != nil {
		return Tickets{}, db.db.Create(&ticket).Error
	}

	return *ticket, nil
}

func (db database) GetTicket(uuid string) (Tickets, error) {
	ticket := Tickets{}

	results := db.db.Model(&Tickets{}).Where("uuid = ?", uuid).Find(&ticket)

	if results.Error != nil {
		return Tickets{}, fmt.Errorf("failed to get ticket: %w", results.Error)
	}

	if results.RowsAffected == 0 {
		return Tickets{}, fmt.Errorf("failed to get ticket: %w", results.Error)
	}

	return ticket, nil
}

func (db database) UpdateTicket(ticket Tickets) (Tickets, error) {
	if ticket.UUID == uuid.Nil {
		return Tickets{}, errors.New("ticket UUID is required")
	}

	if ticket.FeatureUUID == "" || ticket.PhaseUUID == "" || ticket.Name == "" {
		return Tickets{}, errors.New("feature_uuid, phase_uuid, and name are required")
	}

	var existingTicket Tickets
	result := db.db.Where("uuid = ?", ticket.UUID).First(&existingTicket)

	now := time.Now()
	ticket.UpdatedAt = now

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ticket.CreatedAt = now
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
