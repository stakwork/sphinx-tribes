package db

import (
	"fmt"
	"log"
	"os"

	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/logger"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configHandler struct {
	db Database
}

func NewConfigHandler(database Database) *configHandler {
	return &configHandler{
		db: database,
	}
}

type database struct {
	db                 *gorm.DB
	getWorkspaceByUuid func(uuid string) Workspace
	getUserRoles       func(uuid string, pubkey string) []WorkspaceUserRoles
}

func NewDatabaseConfig(db *gorm.DB) *database {
	return &database{
		db:                 db,
		getWorkspaceByUuid: DB.GetWorkspaceByUuid,
		getUserRoles:       DB.GetUserRoles,
	}
}

// DB is the object
var DB database

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	logger.Log.Info("db url : %v", dbURL)

	if dbURL == "" {
		rdsHost := os.Getenv("RDS_HOSTNAME")
		rdsPort := os.Getenv("RDS_PORT")
		rdsDbName := os.Getenv("RDS_DB_NAME")
		rdsUsername := os.Getenv("RDS_USERNAME")
		rdsPassword := os.Getenv("RDS_PASSWORD")
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", rdsUsername, rdsPassword, rdsHost, rdsPort, rdsDbName)
	}

	if dbURL == "" {
		panic("DB env vars not found")
	}

	var err error

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.db = db
	logger.Log.Info("db connected")

	// migrate table changes
	db.AutoMigrate(&Tribe{})
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Channel{})
	db.AutoMigrate(&LeaderBoard{})
	db.AutoMigrate(&ConnectionCodes{})
	db.AutoMigrate(&BountyRoles{})
	db.AutoMigrate(&UserInvoiceData{})
	db.AutoMigrate(&WorkspaceRepositories{})
	db.AutoMigrate(&WorkspaceCodeGraph{})
	db.AutoMigrate(&WorkspaceFeatures{})
	db.AutoMigrate(&FeaturePhase{})
	db.AutoMigrate(&FeatureStory{})
	db.AutoMigrate(&WfRequest{})
	db.AutoMigrate(&WfProcessingMap{})
	db.AutoMigrate(&Tickets{})
	db.AutoMigrate(&ChatMessage{})
	db.AutoMigrate(&Chat{})
	db.AutoMigrate(&ProofOfWork{})
	db.AutoMigrate(&BountyTiming{})
	db.AutoMigrate(&FeatureFlag{})
	db.AutoMigrate(&Endpoint{})
	db.AutoMigrate(&FeaturedBounty{})
	db.AutoMigrate(&Notification{})
	db.AutoMigrate(&TextSnippet{})
	db.AutoMigrate(&BountyTiming{})
	db.AutoMigrate(&FileAsset{})
	db.AutoMigrate(&TicketPlan{})

	DB.MigrateTablesWithOrgUuid()
	DB.MigrateOrganizationToWorkspace()

	people := DB.GetAllPeople()
	for _, p := range people {
		if p.Uuid == "" {
			DB.AddUuidToPerson(p.ID, xid.New().String())
		}
	}

}

const (
	EditOrg        = "EDIT ORGANIZATION"
	AddBounty      = "ADD BOUNTY"
	UpdateBounty   = "UPDATE BOUNTY"
	DeleteBounty   = "DELETE BOUNTY"
	PayBounty      = "PAY BOUNTY"
	AddUser        = "ADD USER"
	UpdateUser     = "UPDATE USER"
	DeleteUser     = "DELETE USER"
	AddRoles       = "ADD ROLES"
	AddBudget      = "ADD BUDGET"
	WithdrawBudget = "WITHDRAW BUDGET"
	ViewReport     = "VIEW REPORT"
)

var ConfigBountyRoles []BountyRoles = []BountyRoles{
	{
		Name: EditOrg,
	},
	{
		Name: AddBounty,
	},
	{
		Name: UpdateBounty,
	},
	{
		Name: DeleteBounty,
	},
	{
		Name: PayBounty,
	},
	{
		Name: AddUser,
	},
	{
		Name: UpdateUser,
	},
	{
		Name: DeleteUser,
	},
	{
		Name: AddRoles,
	},
	{
		Name: AddBudget,
	},
	{
		Name: WithdrawBudget,
	},
	{
		Name: ViewReport,
	},
}

var ManageBountiesGroup = []string{AddBounty, UpdateBounty, DeleteBounty, PayBounty}

var Updatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_to_join", "price_per_message",
	"escrow_amount", "escrow_millis",
	"unlisted", "private", "deleted",
	"app_url", "bots", "feed_url", "feed_type",
	"owner_route_hint", "updated", "pin",
	"profile_filters", "second_brain_url",
}
var Botupdatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_per_use",
	"unlisted", "deleted",
	"owner_route_hint", "updated",
}
var Peopleupdatables = []string{
	"description", "tags", "img",
	"owner_alias",
	"unlisted", "deleted",
	"owner_route_hint",
	"price_to_meet", "updated",
	"extras",
}

var Validate *validator.Validate = validator.New()

var Channelupdatables = []string{
	"name", "deleted"}

func (db database) GetRolesCount() int64 {
	var count int64
	query := db.db.Model(&BountyRoles{})

	query.Count(&count)
	return count
}

func (db database) MigrateTablesWithOrgUuid() {
	if db.db.Migrator().HasTable("bounty") {
		if !db.db.Migrator().HasColumn(Bounty{}, "workspace_uuid") {
			db.db.AutoMigrate(&Bounty{})
		} else {
			db.db.AutoMigrate(&NewBounty{})
		}
	} else {
		db.db.AutoMigrate(&NewBounty{})
	}
	if !db.db.Migrator().HasTable("budget_histories") {
		if !db.db.Migrator().HasColumn(BudgetHistory{}, "workspace_uuid") {
			db.db.AutoMigrate(&BudgetHistory{})
		}
	}
	if !db.db.Migrator().HasTable("payment_histories") {
		if !db.db.Migrator().HasColumn(PaymentHistory{}, "workspace_uuid") {
			db.db.AutoMigrate(&PaymentHistory{})
		} else {
			db.db.AutoMigrate(&NewPaymentHistory{})
		}
	} else {
		db.db.AutoMigrate(&NewPaymentHistory{})
	}
	if !db.db.Migrator().HasTable("invoice_list") {
		if !db.db.Migrator().HasColumn(InvoiceList{}, "workspace_uuid") {
			db.db.AutoMigrate(&InvoiceList{})
		} else {
			db.db.AutoMigrate(&NewInvoiceList{})
		}
	}
	if !db.db.Migrator().HasTable("bounty_budgets") {
		if !db.db.Migrator().HasColumn(BountyBudget{}, "workspace_uuid") {
			db.db.AutoMigrate(&BountyBudget{})
		} else {
			db.db.AutoMigrate(&NewBountyBudget{})
		}
	}
	if !db.db.Migrator().HasTable("workspace_user_roles") {
		if !db.db.Migrator().HasColumn(UserRoles{}, "workspace_uuid") {
			db.db.AutoMigrate(&UserRoles{})
		}
	} else {
		db.db.AutoMigrate(&WorkspaceUserRoles{})
	}
	if !db.db.Migrator().HasTable("workspaces") {
		db.db.AutoMigrate(&Organization{})
	} else {
		db.db.AutoMigrate(&Workspace{})
	}
	if !db.db.Migrator().HasTable("workspace_users") {
		db.db.AutoMigrate(&OrganizationUsers{})
	} else {
		db.db.AutoMigrate(&WorkspaceUsers{})
	}
}

func (db database) MigrateOrganizationToWorkspace() {
	if (db.db.Migrator().HasTable(&Organization{}) && !db.db.Migrator().HasTable("workspaces")) {
		db.db.Migrator().RenameTable(&Organization{}, "workspaces")
	}

	if (db.db.Migrator().HasTable(&OrganizationUsers{}) && !db.db.Migrator().HasTable("workspace_users")) {
		if db.db.Migrator().HasColumn(&OrganizationUsers{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&OrganizationUsers{}, "org_uuid", "workspace_uuid")
		}
		db.db.Migrator().RenameTable(&OrganizationUsers{}, "workspace_users")
	}

	if (db.db.Migrator().HasTable(&UserRoles{}) && !db.db.Migrator().HasTable("workspace_user_roles")) {
		if db.db.Migrator().HasColumn(&UserRoles{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&UserRoles{}, "org_uuid", "workspace_uuid")
		}

		db.db.Migrator().RenameTable(&UserRoles{}, "workspace_user_roles")
	}

	if (db.db.Migrator().HasTable(&Bounty{})) {
		if db.db.Migrator().HasColumn(&Bounty{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&Bounty{}, "org_uuid", "workspace_uuid")
		}
	}

	if (db.db.Migrator().HasTable(&BountyBudget{})) {
		if db.db.Migrator().HasColumn(&BountyBudget{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&BountyBudget{}, "org_uuid", "workspace_uuid")
		}
	}

	if (db.db.Migrator().HasTable(&BudgetHistory{})) {
		if db.db.Migrator().HasColumn(&BudgetHistory{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&BudgetHistory{}, "org_uuid", "workspace_uuid")
		}
	}

	if (db.db.Migrator().HasTable(&PaymentHistory{})) {
		if db.db.Migrator().HasColumn(&PaymentHistory{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&PaymentHistory{}, "org_uuid", "workspace_uuid")
		}
	}

	if (db.db.Migrator().HasTable(&InvoiceList{})) {
		if db.db.Migrator().HasColumn(&InvoiceList{}, "org_uuid") {
			db.db.Migrator().RenameColumn(&InvoiceList{}, "org_uuid", "workspace_uuid")
		}
	}
}

func (db database) CreateRoles() {
	db.db.Create(&ConfigBountyRoles)
}

func (db database) DeleteRoles() {
	db.db.Exec("DELETE FROM bounty_roles")
}

func InitRoles() {
	count := DB.GetRolesCount()
	if count != int64(len(ConfigBountyRoles)) {
		// delete all the roles and insert again
		if count != 0 {
			DB.DeleteRoles()
		}
		DB.CreateRoles()
	}
}

func GetRolesMap() map[string]string {
	roles := map[string]string{}
	for _, v := range ConfigBountyRoles {
		roles[v.Name] = v.Name
	}
	return roles
}

func GetUserRolesMap(userRoles []WorkspaceUserRoles) map[string]string {
	roles := map[string]string{}
	for _, v := range userRoles {
		roles[v.Role] = v.Role
	}
	return roles
}

func (db database) ConvertMetricsBountiesToMap(metricsCsv []MetricsBountyCsv) []map[string]interface{} {
	var metricsMap []map[string]interface{}
	for _, m := range metricsCsv {
		metricMap := make(map[string]interface{})

		metricMap["DatePosted"] = m.DatePosted
		metricMap["Workspace"] = m.Organization
		metricMap["BountyAmount"] = m.BountyAmount
		metricMap["Provider"] = m.Provider
		metricMap["Hunter"] = m.Hunter
		metricMap["BountyTitle"] = m.BountyTitle
		metricMap["BountyLink"] = m.BountyLink
		metricMap["BountyStatus"] = m.BountyStatus
		metricMap["DateAssigned"] = m.DateAssigned
		metricMap["DatePaid"] = m.DatePaid

		metricsMap = append(metricsMap, metricMap)
	}

	return metricsMap
}

func RolesCheck(userRoles []WorkspaceUserRoles, check string) bool {
	rolesMap := GetRolesMap()
	userRolesMap := GetUserRolesMap(userRoles)

	// check if roles exist in config
	_, ok := rolesMap[check]
	_, ok1 := userRolesMap[check]

	// if any of the roles does not exist, return false
	// if any of the roles does not exist, user roles return false
	if !ok {
		return false
	} else if !ok1 {
		return false
	}
	return true
}

func CheckUser(userRoles []WorkspaceUserRoles, pubkey string) bool {
	for _, role := range userRoles {
		if role.OwnerPubKey == pubkey {
			return true
		}
	}
	return false
}

func UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool {
	org := DB.GetWorkspaceByUuid(uuid)
	var hasRole bool = false
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := DB.GetUserRoles(uuid, pubKeyFromAuth)
		hasRole = RolesCheck(userRoles, role)
		return hasRole
	}
	return true
}

func (db database) UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool {
	org := db.getWorkspaceByUuid(uuid)
	var hasRole bool = false
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := db.getUserRoles(uuid, pubKeyFromAuth)
		hasRole = RolesCheck(userRoles, role)
		return hasRole
	}
	return true
}

func (ch configHandler) UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool {
	org := ch.db.GetWorkspaceByUuid(uuid)
	var hasRole bool = false
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := ch.db.GetUserRoles(uuid, pubKeyFromAuth)
		hasRole = RolesCheck(userRoles, role)
		return hasRole
	}
	return true
}

func (ch configHandler) UserHasManageBountyRoles(pubKeyFromAuth string, uuid string) bool {
	var manageRolesCount = len(ManageBountiesGroup)
	org := ch.db.GetWorkspaceByUuid(uuid)
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := ch.db.GetUserRoles(uuid, pubKeyFromAuth)

		for _, role := range ManageBountiesGroup {
			// check for the manage bounty roles
			hasRole := RolesCheck(userRoles, role)
			if hasRole {
				manageRolesCount--
			}
		}

		if manageRolesCount != 0 {
			return false
		}
	}
	return true
}

func (db database) UserHasManageBountyRoles(pubKeyFromAuth string, uuid string) bool {
	var manageRolesCount = len(ManageBountiesGroup)
	org := db.getWorkspaceByUuid(uuid)
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := db.getUserRoles(uuid, pubKeyFromAuth)

		for _, role := range ManageBountiesGroup {
			// check for the manage bounty roles
			hasRole := RolesCheck(userRoles, role)
			if hasRole {
				manageRolesCount--
			}
		}

		if manageRolesCount != 0 {
			return false
		}
	}
	return true
}

func (db database) ProcessUpdateTicketsWithoutGroup() {
	// get all tickets without group
	tickets, err := db.GetTicketsWithoutGroup()

	if err != nil {
		log.Printf("Error getting tickets without group: %v", err)
		return
	}

	// update each ticket with group uuid
	for _, ticket := range tickets {
		logger.Log.Info("ticket from process: %v", ticket)
		err := db.UpdateTicketsWithoutGroup(ticket)
		if err != nil {
			log.Printf("Error updating ticket: %v", err)
		}
	}
}
