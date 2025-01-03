package db

import (
	"testing"

	"github.com/google/uuid"
)

func TestProcessUpdateBudget(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// create a user
	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test_user_update_budget",
		OwnerAlias:  "test_user_update_budget",
		Description: "test_user_update_budget_description",
	}

	TestDB.db.Create(&person)

	// Create a new workspace
	uuid := uuid.New()
	workspace := Workspace{
		OwnerPubKey: person.OwnerPubKey,
		Uuid:        uuid.String(),
		Name:        "Test Workspace",
	}

	TestDB.db.Create(&workspace)

	// create invoice
	invoice := NewInvoiceList{
		WorkspaceUuid:  workspace.Uuid,
		PaymentRequest: "test_payment_request",
		Status:         true,
		OwnerPubkey:    person.OwnerPubKey,
	}

	TestDB.db.Create(&invoice)

	// create paymentHistory

	amount := 50000

	paymentHistory := NewPaymentHistory{
		WorkspaceUuid: workspace.Uuid,
		Amount:        uint(amount),
		PaymentStatus: PaymentComplete,
		PaymentType:   Deposit,
		SenderPubKey:  person.OwnerPubKey,
	}

	TestDB.db.Create(&paymentHistory)

	// Process the update budget
	err := TestDB.ProcessUpdateBudget(invoice)
	if err != nil {
		t.Fatalf("Failed to process update budget: %v", err)
	}
}
