package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateOrEditTicket(t *testing.T) {
	// test create or edit tickers
	InitTestDB()

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

	TestDB.CreateOrEditPerson(person)

	// create workspace
	workspace := Workspace{
		Uuid:    uuid.New().String(),
		Name:    "Test tickets space",
		Created: &now,
		Updated: &now,
	}

	TestDB.CreateOrEditWorkspace(workspace)

	// create WorkspaceFeatures
	workspaceFeatures := WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test",
		Brief:         "test breieft",
		Requirements:  "Test requirements",
		Architecture:  "Test architecture",
		Url:           "Test url",
		Priority:      1,
		Created:       &now,
		Updated:       &now,
		CreatedBy:     "test",
		UpdatedBy:     "test",
	}

	TestDB.CreateOrEditFeature(workspaceFeatures)

	// create FeaturePhase
	featurePhase := FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: workspaceFeatures.Uuid,
		Name:        "test feature phase",
		Priority:    1,
		Created:     &now,
		Updated:     &now,
	}

	TestDB.CreateOrEditFeaturePhase(featurePhase)

	// test that an error is returned if the required fields are missing
	t.Run("test that an error is returned if the required fields are missing", func(t *testing.T) {
		ticket := Tickets{
			UUID:        uuid.New().String(),
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
			UUID:        uuid.New().String(),
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
	t.Run("should create a ticket if all fields are provided", func(t *testing.T) {
		ticket := Tickets{
			UUID:        uuid.New().String(),
			FeatureUUID: workspaceFeatures.Uuid,
			PhaseUUID:   featurePhase.Uuid,
			Name:        "test ticket",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		_, err := TestDB.CreateOrEditTicket(&ticket)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	})
}
