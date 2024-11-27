package db

import (
	"errors"
	"fmt"
	"time"
)

func (db database) CreateOrEditTicket(ticket *Tickets) (Tickets, error) {

	if ticket.UUID == "" || ticket.FeatureUUID == "" || ticket.PhaseUUID == "" || ticket.Name == "" {
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
