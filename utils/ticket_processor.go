package utils

import (
	"errors"
)

type TicketReviewRequest struct {
	TicketUUID        string `json:"ticketUUID" validate:"required"`
	TicketDescription string `json:"ticketDescription" validate:"required"`
}

func ValidateTicketReviewRequest(req *TicketReviewRequest) error {
	if req.TicketUUID == "" {
		return errors.New("ticketUUID is required")
	}
	if req.TicketDescription == "" {
		return errors.New("ticketDescription is required")
	}
	return nil
}
