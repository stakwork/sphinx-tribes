package utils

import (
	"errors"
)

type TicketReviewRequest struct {
	FeatureUUID       string `json:"featureUUID" validate:"required"`
	PhaseUUID         string `json:"phaseUUID" validate:"required"`
	TicketUUID        string `json:"ticketUUID" validate:"required"`
	TicketDescription string `json:"ticketDescription" validate:"required"`
}

func ValidateTicketReviewRequest(req *TicketReviewRequest) error {
	if req.FeatureUUID == "" {
		return errors.New("featureUUID is required")
	}
	if req.PhaseUUID == "" {
		return errors.New("phaseUUID is required")
	}
	if req.TicketUUID == "" {
		return errors.New("ticketUUID is required")
	}
	if req.TicketDescription == "" {
		return errors.New("ticketDescription is required")
	}
	return nil
}
