package db

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func SetupSuite(_ *testing.T) func(tb testing.TB) {
	InitTestDB()

	return func(_ testing.TB) {
		defer CloseTestDB()
		log.Println("Teardown test")
	}
}

func TestCreateOrEditTicket(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	// create person
	now := time.Now()

	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "testfeaturepubkey",
		OwnerAlias:  "testfeaturealias",
		Description: "testfeaturedescription",
		Created:     &now,
		Updated:     &now,
		Deleted:     false,
	}

	workspace := Workspace{
		Uuid:    uuid.New().String(),
		Name:    "Test tickets space",
		Created: &now,
		Updated: &now,
	}

	workspaceFeatures := WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test",
		Brief:         "test brief",
		Requirements:  "Test requirements",
		Architecture:  "Test architecture",
		Url:           "Test url",
		Priority:      1,
		Created:       &now,
		Updated:       &now,
		CreatedBy:     "test",
		UpdatedBy:     "test",
	}

	featurePhase := FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: workspaceFeatures.Uuid,
		Name:        "test feature phase",
		Priority:    1,
		Created:     &now,
		Updated:     &now,
	}

	ticket := Tickets{
		UUID:        uuid.New(),
		FeatureUUID: workspaceFeatures.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "test ticket",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// create person
	TestDB.CreateOrEditPerson(person)

	// create workspace
	TestDB.CreateOrEditWorkspace(workspace)

	// create WorkspaceFeatures
	TestDB.CreateOrEditFeature(workspaceFeatures)

	// create FeaturePhase
	TestDB.CreateOrEditFeaturePhase(featurePhase)

	// test that an error is returned if the required fields are missing
	t.Run("test that an error is returned if the required fields are missing", func(t *testing.T) {
		ticket := Tickets{
			UUID:        uuid.New(),
			FeatureUUID: "",
			PhaseUUID:   "",
			Name:        "test ticket",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		_, err := TestDB.CreateOrEditTicket(&ticket)
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	})

	// test that an error is thrown if the FeatureUUID, and PhaseUUID does not exists
	t.Run("test that an error is returned if the required fields are missing", func(t *testing.T) {
		ticket := Tickets{
			UUID:        uuid.New(),
			FeatureUUID: "testfeatureuuid",
			PhaseUUID:   "testphaseuuid",
			Name:        "test ticket",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		_, err := TestDB.CreateOrEditTicket(&ticket)
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	})

	// shpuld create a ticket if all fields are provided
	t.Run("should create a ticket if the all fields are provided", func(t *testing.T) {

		_, err := TestDB.CreateOrEditTicket(&ticket)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	})
}

func TestGetTicket(t *testing.T) {
	InitTestDB()

	defer CloseTestDB()

	// create person
	now := time.Now()

	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "testfeaturepubkey",
		OwnerAlias:  "testfeaturealias",
		Description: "testfeaturedescription",
		Created:     &now,
		Updated:     &now,
		Deleted:     false,
	}

	// create person
	TestDB.CreateOrEditPerson(person)

	workspace := Workspace{
		Uuid:    uuid.New().String(),
		Name:    "Test tickets space",
		Created: &now,
		Updated: &now,
	}

	// create workspace
	TestDB.CreateOrEditWorkspace(workspace)

	workspaceFeatures := WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test",
		Brief:         "test get brief",
		Requirements:  "Test get requirements",
		Architecture:  "Test get architecture",
		Url:           "Test get url",
		Priority:      1,
		Created:       &now,
		Updated:       &now,
		CreatedBy:     "test",
		UpdatedBy:     "test",
	}

	// create WorkspaceFeatures
	TestDB.CreateOrEditFeature(workspaceFeatures)

	featurePhase := FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: workspaceFeatures.Uuid,
		Name:        "test get feature phase",
		Priority:    1,
		Created:     &now,
		Updated:     &now,
	}

	// create FeaturePhase
	TestDB.CreateOrEditFeaturePhase(featurePhase)

	ticket := Tickets{
		UUID:        uuid.New(),
		FeatureUUID: workspaceFeatures.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "test get ticket",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// create ticket
	TestDB.CreateOrEditTicket(&ticket)

	// test that an error is returned if the ticket does not exist
	t.Run("test that an error is returned if the ticket does not exist", func(t *testing.T) {
		_, err := TestDB.GetTicket(uuid.New().String())
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	})

	// should return a ticket if it exists
	t.Run("should return a ticket if it exists", func(t *testing.T) {
		result, err := TestDB.GetTicket(ticket.UUID.String())
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if result.UUID != ticket.UUID {
			t.Errorf("expected %v but got %v", ticket.UUID, result.UUID)
		}

		if result.FeatureUUID != ticket.FeatureUUID {
			t.Errorf("expected %v but got %v", ticket.FeatureUUID, result.FeatureUUID)
		}

		if result.PhaseUUID != ticket.PhaseUUID {
			t.Errorf("expected %v but got %v", ticket.PhaseUUID, result.PhaseUUID)
		}

		if result.Name != ticket.Name {
			t.Errorf("expected %v but got %v", ticket.Name, result.Name)
		}
	})
}
