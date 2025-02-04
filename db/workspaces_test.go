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

func TestAddAndUpdateBudget(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test_user_add_and_update_budget",
		OwnerAlias:  "test_user_add_and_update_budget",
		Description: "test_user_add_and_update_budget_description",
	}

	TestDB.db.Create(&person)

	// Create a new workspace
	uuid := uuid.New()

	randomWorkspaceName := fmt.Sprintf("Test Workspace Add and Update Budget %d", rand.Intn(1000))
	workspace := Workspace{
		OwnerPubKey: person.OwnerPubKey,
		Uuid:        uuid.String(),
		Name:        randomWorkspaceName,
	}

	TestDB.db.Create(&workspace)

	amount := 50000

	t.Run("Should test that the budget is updated", func(t *testing.T) {
		now := time.Now()

		randomPaymentRequest := fmt.Sprintf("test_add_and_update_budget_payment_request_%d", rand.Intn(1000))
		// create invoice
		invoice := NewInvoiceList{
			WorkspaceUuid:  workspace.Uuid,
			PaymentRequest: randomPaymentRequest,
			Status:         false,
			OwnerPubkey:    person.OwnerPubKey,
			Created:        &now,
			Type:           "BUDGET",
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
		TestDB.AddAndUpdateBudget(invoice)

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

		randomPaymentRequest := fmt.Sprintf("test_add_and_update_budget_payment_request_%d", rand.Intn(1000))
		// create invoice
		invoice := NewInvoiceList{
			WorkspaceUuid:  workspace.Uuid,
			PaymentRequest: randomPaymentRequest,
			Status:         true,
			Created:        &now,
		}

		TestDB.db.Create(&invoice)

		// Process the update budget
		TestDB.AddAndUpdateBudget(invoice)

		// get workspace budget
		workspaceBudget := TestDB.GetWorkspaceBudget(workspace.Uuid)

		// assert that balance is not updated
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount))
	})

	t.Run("Assert that budget is accurate after multiple updates", func(t *testing.T) {

		totalAmount := amount

		for i := 0; i < 10; i++ {
			now := time.Now()

			randomPaymentRequest := fmt.Sprintf("test_add_and_update_budget_payment_request_%d", rand.Intn(1000))
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
			TestDB.AddAndUpdateBudget(invoice)
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

func TestGetUserCreatedWorkspaces(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()


	t.Run("Basic Retrieval with Valid Pubkey", func(t *testing.T) {

		person := Person{
			Uuid:        uuid.New().String(),
			OwnerPubKey: "test_user_workspaces",
			OwnerAlias:  "test_user",
		}
		TestDB.db.Create(&person)

		workspace1 := Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.New().String(),
			Name:        "Test Workspace 1",
		}
		workspace2 := Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.New().String(),
			Name:        "Test Workspace 2",
		}
		TestDB.db.Create(&workspace1)
		TestDB.db.Create(&workspace2)

		workspaces := TestDB.GetUserCreatedWorkspaces(person.OwnerPubKey)
		assert.Equal(t, 2, len(workspaces))
		assert.Equal(t, workspace1.Name, workspaces[0].Name)
		assert.Equal(t, workspace2.Name, workspaces[1].Name)
	})

	t.Run("No Workspaces Found", func(t *testing.T) {
		workspaces := TestDB.GetUserCreatedWorkspaces("non_existent_pubkey")
		assert.Equal(t, 0, len(workspaces))
	})

	t.Run("Empty Pubkey", func(t *testing.T) {
		workspaces := TestDB.GetUserCreatedWorkspaces("")
		assert.Equal(t, 0, len(workspaces))
	})

	t.Run("Workspaces with a Mixture of Deleted and Non-Deleted", func(t *testing.T) {
		person := Person{
			Uuid:        uuid.New().String(),
			OwnerPubKey: "test_user_deleted_mix",
			OwnerAlias:  "test_user_deleted",
		}
		TestDB.db.Create(&person)

		workspace1 := Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.New().String(),
			Name:        "Active Workspace",
			Deleted:     false,
		}
		workspace2 := Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.New().String(),
			Name:        "Deleted Workspace",
			Deleted:     true,
		}
		TestDB.db.Create(&workspace1)
		TestDB.db.Create(&workspace2)

		workspaces := TestDB.GetUserCreatedWorkspaces(person.OwnerPubKey)
		assert.Equal(t, 1, len(workspaces))
		assert.Equal(t, workspace1.Name, workspaces[0].Name)
	})

	t.Run("SQL Injection Prevention", func(t *testing.T) {
		maliciousPubkey := "' OR '1'='1"
		workspaces := TestDB.GetUserCreatedWorkspaces(maliciousPubkey)
		assert.Equal(t, 0, len(workspaces))
	})

	t.Run("Performance Test with Large Result Set", func(t *testing.T) {
		person := Person{
			Uuid:        uuid.New().String(),
			OwnerPubKey: "test_user_performance",
			OwnerAlias:  "test_user_perf",
		}
		TestDB.db.Create(&person)

		for i := 0; i < 100; i++ {
			workspace := Workspace{
				OwnerPubKey: person.OwnerPubKey,
				Uuid:        uuid.New().String(),
				Name:        fmt.Sprintf("Performance Workspace %d", i),
			}
			TestDB.db.Create(&workspace)
		}

		start := time.Now()
		workspaces := TestDB.GetUserCreatedWorkspaces(person.OwnerPubKey)
		duration := time.Since(start)

		assert.Equal(t, 100, len(workspaces))
		assert.Less(t, duration.Milliseconds(), int64(1000), "Query should complete within 1 second")
	})
}
