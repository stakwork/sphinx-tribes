package db

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB database

func InitTestDB() {
	rdsHost := "localhost"
	rdsPort := fmt.Sprintf("%d", 5532)
	rdsDbName := "test_db"
	rdsUsername := "test_user"
	rdsPassword := "test_password"
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", rdsUsername, rdsPassword, rdsHost, rdsPort, rdsDbName)

	if dbURL == "" {
		panic("TESTDB URL is not set")
	}

	var err error

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	TestDB.db = db

	logger.Log.Info("DB CONNECTED")

	db.Exec(`
		DO $$ BEGIN
			ALTER TABLE tickets 
				ALTER COLUMN feature_uuid DROP NOT NULL,
				ALTER COLUMN phase_uuid DROP NOT NULL;
		EXCEPTION
			WHEN others THEN null;
		END $$;
	`)

	// migrate table changes
	db.AutoMigrate(&Tribe{})
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Channel{})
	db.AutoMigrate(&LeaderBoard{})
	db.AutoMigrate(&ConnectionCodes{})
	db.AutoMigrate(&BountyRoles{})
	db.AutoMigrate(&UserInvoiceData{})
	db.AutoMigrate(&WorkspaceRepositories{})
	db.AutoMigrate(&WorkspaceFeatures{})
	db.AutoMigrate(&FeaturePhase{})
	db.AutoMigrate(&FeatureStory{})
	db.AutoMigrate(&NewBounty{})
	db.AutoMigrate(&BudgetHistory{})
	db.AutoMigrate(&NewPaymentHistory{})
	db.AutoMigrate(&NewInvoiceList{})
	db.AutoMigrate(&NewBountyBudget{})
	db.AutoMigrate(&Workspace{})
	db.AutoMigrate(&WorkspaceUsers{})
	db.AutoMigrate(&WorkspaceUserRoles{})
	db.AutoMigrate(&Bot{})
	db.AutoMigrate(&WfRequest{})
	db.AutoMigrate(&WfProcessingMap{})
	db.AutoMigrate(&Tickets{})
	db.AutoMigrate(&ChatMessage{})
	db.AutoMigrate(&Chat{})
	db.AutoMigrate(&WorkspaceCodeGraph{})
	db.AutoMigrate(&FeatureFlag{})
	db.AutoMigrate(&Bounty{})
	db.AutoMigrate(&Notification{})
	db.AutoMigrate(&BountyTiming{})
	db.AutoMigrate(&FileAsset{})
	db.AutoMigrate(&TextSnippet{})

	people := TestDB.GetAllPeople()
	for _, p := range people {
		if p.Uuid == "" {
			TestDB.AddUuidToPerson(p.ID, xid.New().String())
		}
	}
}

func CleanDB() {
	TestDB.db.Exec("DELETE FROM people")
}

func DeleteAllChats() {
	TestDB.db.Exec("DELETE FROM chats")
}

func CleanTestData() {

	TestDB.DeleteAllBounties()

	TestDB.db.Exec("DELETE FROM workspaces")

	TestDB.db.Exec("DELETE FROM workspace_features")

	TestDB.db.Exec("DELETE FROM feature_phases")

	TestDB.db.Exec("DELETE FROM people")

	TestDB.db.Exec("DELETE FROM tickets")

	TestDB.db.Exec("DELETE FROM chats")

	TestDB.db.Exec("DELETE FROM notifications")
}

func DeleteAllChatMessages() {
	TestDB.db.Exec("DELETE FROM chat_messages")
}

func CloseTestDB() {
	instance, err := TestDB.db.DB()

	if err != nil {
		panic(err)
	}

	instance.Close()
}
