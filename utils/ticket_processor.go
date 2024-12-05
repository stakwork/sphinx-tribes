package utils

import (
	"errors"
)

type TicketReviewRequest struct {
	Value struct {
		FeatureUUID       string `json:"featureUUID"`
		PhaseUUID         string `json:"phaseUUID"`
		TicketUUID        string `json:"ticketUUID" validate:"required"`
		TicketDescription string `json:"ticketDescription" validate:"required"`
	} `json:"value"`
	RequestUUID     string `json:"requestUUID"`
	SourceWebsocket string `json:"sourceWebsocket"`
}

func ValidateTicketReviewRequest(req *TicketReviewRequest) error {
	if req.Value.TicketUUID == "" {
		return errors.New("ticketUUID is required")
	}
	if req.Value.TicketDescription == "" {
		return errors.New("ticketDescription is required")
	}
	return nil
}
