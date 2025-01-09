package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTicketReviewRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   *TicketReviewRequest
		wantErr string
	}{
		{
			name: "valid request",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "test-uuid",
					TicketDescription: "test description",
				},
			},
			wantErr: "",
		},
		{
			name: "missing ticket UUID",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketDescription: "test description",
				},
			},
			wantErr: "ticketUUID is required",
		},
		{
			name: "missing ticket description",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID: "test-uuid",
				},
			},
			wantErr: "ticketDescription is required",
		},
		{
			name: "Both TicketUUID and TicketDescription Empty",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "",
					TicketDescription: "",
				},
			},
			wantErr: "ticketUUID is required",
		},
		{
			name:    "Nil Request",
			input:   nil,
			wantErr: "nil request",
		},
		{
			name: "Whitespace TicketUUID",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "   ",
					TicketDescription: "This is a valid ticket description",
				},
			},
			wantErr: "ticketUUID is required",
		},
		{
			name: "Whitespace TicketDescription",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "123e4567-e89b-12d3-a456-426614174000",
					TicketDescription: "   ",
				},
			},
			wantErr: "ticketDescription is required",
		},
		{
			name: "Large TicketDescription",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "123e4567-e89b-12d3-a456-426614174000",
					TicketDescription: strings.Repeat("a", 10000),
				},
			},
			wantErr: "",
		},
		{
			name: "Non-UUID TicketUUID",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "not-a-uuid",
					TicketDescription: "This is a valid ticket description",
				},
			},
			wantErr: "",
		},
		{
			name: "With Optional Fields",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					FeatureUUID:       "feature-123",
					PhaseUUID:         "phase-456",
					TicketUUID:        "ticket-789",
					TicketDescription: "test description",
				},
			},
			wantErr: "",
		},
		{
			name: "With RequestUUID and SourceWebsocket",
			input: &TicketReviewRequest{
				Value: struct {
					FeatureUUID       string `json:"featureUUID"`
					PhaseUUID         string `json:"phaseUUID"`
					TicketUUID        string `json:"ticketUUID" validate:"required"`
					TicketDescription string `json:"ticketDescription" validate:"required"`
				}{
					TicketUUID:        "ticket-789",
					TicketDescription: "test description",
				},
				RequestUUID:     "req-123",
				SourceWebsocket: "ws://example.com",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTicketReviewRequest(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
