package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
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

	randomWorkspaceName := fmt.Sprintf("Test Workspace Budget %d", rand.Intn(1000))
	workspace := Workspace{
		OwnerPubKey: person.OwnerPubKey,
		Uuid:        uuid.String(),
		Name:        randomWorkspaceName,
	}

	TestDB.db.Create(&workspace)

	amount := 50000

	t.Run("Should test that the budget is updated", func(t *testing.T) {
		now := time.Now()

		randomPaymentRequest := fmt.Sprintf("test_update_budget_payment_request_%d", rand.Intn(1000))
		// create invoice
		invoice := NewInvoiceList{
			WorkspaceUuid:  workspace.Uuid,
			PaymentRequest: randomPaymentRequest,
			Status:         false,
			OwnerPubkey:    person.OwnerPubKey,
			Created:        &now,
		}

		TestDB.db.Create(&invoice)

		// create paymentHistory

		paymentHistory := NewPaymentHistory{
			WorkspaceUuid: workspace.Uuid,
			Amount:        uint(amount),
			PaymentStatus: PaymentComplete,
			PaymentType:   Deposit,
			SenderPubKey:  person.OwnerPubKey,
			Created:       &now,
		}

		TestDB.db.Create(&paymentHistory)

		// Process the update budget
		err := TestDB.ProcessUpdateBudget(invoice)
		if err != nil {
			t.Fatalf("Failed to process update budget: %v", err)
		}

		// get workspace budget
		workspaceBudget := TestDB.GetWorkspaceBudget(workspace.Uuid)

		if workspaceBudget.TotalBudget != uint(amount) {
			t.Fatalf("Total budget is not correct: %v", workspaceBudget.TotalBudget)
		}

		// assert that balance is updated
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount))
	})

	t.Run("Budget should not be updated if the invoice is already paid	", func(t *testing.T) {

		now := time.Now()

		randomPaymentRequest := fmt.Sprintf("test_payment_request_%d", rand.Intn(1000))
		// create invoice
		invoice := NewInvoiceList{
			WorkspaceUuid:  workspace.Uuid,
			PaymentRequest: randomPaymentRequest,
			Status:         true,
			Created:        &now,
		}

		TestDB.db.Create(&invoice)

		// Process the update budget
		err := TestDB.ProcessUpdateBudget(invoice)

		// assert that error is returned
		assert.Error(t, err, "Expected error to be returned")

		// get workspace budget
		workspaceBudget := TestDB.GetWorkspaceBudget(workspace.Uuid)

		// assert that balance is not updated
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount))
	})

	t.Run("Assert that budget is accurate after multiple updates", func(t *testing.T) {

		totalAmount := amount

		for i := 0; i < 10; i++ {
			now := time.Now()

			randomPaymentRequest := fmt.Sprintf("test_update_budget_payment_request_%d", rand.Intn(1000))
			// create invoice
			invoice := NewInvoiceList{
				WorkspaceUuid:  workspace.Uuid,
				PaymentRequest: randomPaymentRequest,
				Status:         false,
				OwnerPubkey:    person.OwnerPubKey,
				Created:        &now,
			}

			TestDB.ProcessUpdateBudget(invoice)

			totalAmount += amount

			paymentHistory := NewPaymentHistory{
				WorkspaceUuid: workspace.Uuid,
				Amount:        uint(amount),
				PaymentStatus: PaymentComplete,
				PaymentType:   Deposit,
				SenderPubKey:  person.OwnerPubKey,
				Created:       &now,
			}

			TestDB.db.Create(&paymentHistory)

			// Process the update budget
			err := TestDB.ProcessUpdateBudget(invoice)
			if err != nil {
				t.Fatalf("Failed to process update budget: %v", err)
			}
		}

		// get workspace budget
		workspaceBudget := TestDB.GetWorkspaceBudget(workspace.Uuid)

		if workspaceBudget.TotalBudget != uint(totalAmount) {
			t.Fatalf("Total budget is not correct: %v", workspaceBudget.TotalBudget)
		}

		// assert that balance is updated
		assert.Equal(t, workspaceBudget.TotalBudget, uint(totalAmount))
	})
}
