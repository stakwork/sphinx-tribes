package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type ProcessorConfig struct {
	RequiresProcessing bool
	HandlerFunc        string
	Config             json.RawMessage
}

func ProcessWorkflowRequest(requestID string, source string) (string, error) {

	if requestID == "" {
		requestID = uuid.New().String()
	}

	config := lookupProcessingConfig(source)

	if config == nil {
		return requestID, nil
	}

	return processWithHandler(requestID)
}

func lookupProcessingConfig(source string) error {
	return nil
}

func processWithHandler(requestID string) (string, error) {
	fmt.Println("Processing with default handler")
	return fmt.Sprintf("Processed with default handler, RequestID: %s", requestID), nil
}
