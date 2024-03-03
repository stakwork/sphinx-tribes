package db

import (
	"fmt"
	"os"

	"github.com/rs/xid"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

// DB is the object
var DB database

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("db url : %v", dbURL)

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

	fmt.Println("db connected")

	// migrate table changes
	db.AutoMigrate(&Tribe{})
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Channel{})
	db.AutoMigrate(&LeaderBoard{})
	db.AutoMigrate(&ConnectionCodes{})
	db.AutoMigrate(&Bounty{})
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&OrganizationUsers{})
	db.AutoMigrate(&BountyRoles{})
	db.AutoMigrate(&UserRoles{})
	db.AutoMigrate(&BountyBudget{})
	db.AutoMigrate(&BudgetHistory{})
	db.AutoMigrate(&PaymentHistory{})
	db.AutoMigrate(&InvoiceList{})
	db.AutoMigrate(&UserInvoiceData{})

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
	"profile_filters",
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

func GetUserRolesMap(userRoles []UserRoles) map[string]string {
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
		metricMap["Organization"] = m.Organization
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

func RolesCheck(rolesMap map[string]string, userRoles []UserRoles, check string) bool {

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

func CheckUser(userRoles []UserRoles, pubkey string) bool {
	for _, role := range userRoles {
		if role.OwnerPubKey == pubkey {
			return true
		}
	}
	return false
}

func UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool {
	org := DB.GetOrganizationByUuid(uuid)
	var hasRole bool = false
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := DB.GetUserRoles(uuid, pubKeyFromAuth)
		rolesMap := GetRolesMap()
		hasRole = RolesCheck(rolesMap, userRoles, role)
		return hasRole
	}
	return true
}

func (db database) UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool {
	org := DB.GetOrganizationByUuid(uuid)
	var hasRole bool = false
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := DB.GetUserRoles(uuid, pubKeyFromAuth)
		rolesMap := GetRolesMap()
		hasRole = RolesCheck(rolesMap, userRoles, role)
		return hasRole
	}
	return true
}

func (db database) UserHasManageBountyRoles(pubKeyFromAuth string, uuid string) bool {
	var manageRolesCount = len(ManageBountiesGroup)
	org := DB.GetOrganizationByUuid(uuid)
	if pubKeyFromAuth != org.OwnerPubKey {
		userRoles := DB.GetUserRoles(uuid, pubKeyFromAuth)
		rolesMap := GetRolesMap()

		for _, role := range ManageBountiesGroup {
			// check for the manage bounty roles
			hasRole := RolesCheck(rolesMap, userRoles, role)
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
