package db

import (
	"net/http"
	"time"
)

type Database interface {
	CreateOrEditTribe(m Tribe) (Tribe, error)
	CreateChannel(c Channel) (Channel, error)
	CreateOrEditBot(b Bot) (Bot, error)
	CreateOrEditPerson(m Person) (Person, error)
	GetUnconfirmedTwitter() []Person
	UpdateTwitterConfirmed(id uint, confirmed bool)
	GetUnconfirmedGithub() []Person
	UpdateGithubConfirmed(id uint, confirmed bool)
	UpdateGithubIssues(id uint, issues map[string]interface{})
	UpdateTribe(uuid string, u map[string]interface{}) bool
	UpdateChannel(id uint, u map[string]interface{}) bool
	UpdateTribeUniqueName(uuid string, u string)
	GetOpenGithubIssues(r *http.Request) (int64, error)
	GetListedTribes(r *http.Request) []Tribe
	GetTribesByOwner(pubkey string) []Tribe
	GetAllTribesByOwner(pubkey string) []Tribe
	GetTribesByAppUrl(aurl string) []Tribe
	GetChannelsByTribe(tribe_uuid string) []Channel
	GetChannel(id uint) Channel
	GetListedBots(r *http.Request) []Bot
	GetListedPeople(r *http.Request) []Person
	GetPeopleBySearch(r *http.Request) []Person
	GetListedPosts(r *http.Request) ([]PeopleExtra, error)
	GetUserBountiesCount(personKey string, tabType string) int64
	GetBountiesCount(r *http.Request) int64
	GetWorkspaceBounties(r *http.Request, org_uuid string) []Bounty
	GetWorkspaceBountiesCount(r *http.Request, org_uuid string) int64
	GetAssignedBounties(r *http.Request) ([]Bounty, error)
	GetCreatedBounties(r *http.Request) ([]Bounty, error)
	GetBountyById(id string) ([]Bounty, error)
	GetNextBountyByCreated(r *http.Request) (uint, error)
	GetPreviousBountyByCreated(r *http.Request) (uint, error)
	GetNextWorkspaceBountyByCreated(r *http.Request) (uint, error)
	GetPreviousWorkspaceBountyByCreated(r *http.Request) (uint, error)
	GetBountyIndexById(id string) int64
	GetBountyDataByCreated(created string) ([]Bounty, error)
	AddBounty(b Bounty) (Bounty, error)
	GetAllBounties(r *http.Request) []Bounty
	CreateOrEditBounty(b Bounty) (Bounty, error)
	UpdateBountyNullColumn(b Bounty, column string) Bounty
	UpdateBountyBoolColumn(b Bounty, column string) Bounty
	DeleteBounty(pubkey string, created string) (Bounty, error)
	GetBountyByCreated(created uint) (Bounty, error)
	GetBounty(id uint) Bounty
	UpdateBounty(b Bounty) (Bounty, error)
	UpdateBountyPayment(b Bounty) (Bounty, error)
	GetListedOffers(r *http.Request) ([]PeopleExtra, error)
	UpdateBot(uuid string, u map[string]interface{}) bool
	GetAllTribes() []Tribe
	GetTribesTotal() int64
	GetTribeByIdAndPubkey(uuid string, pubkey string) Tribe
	GetTribe(uuid string) Tribe
	GetPerson(id uint) Person
	UpdatePerson(id uint, u map[string]interface{}) bool
	GetPersonByUuid(uuid string) Person
	GetPersonByGithubName(github_name string) Person
	GetFirstTribeByFeedURL(feedURL string) Tribe
	GetBot(uuid string) Bot
	GetTribeByUniqueName(un string) Tribe
	GetBotsByOwner(pubkey string) []Bot
	GetBotByUniqueName(un string) Bot
	GetPersonByUniqueName(un string) Person
	SearchTribes(s string) []Tribe
	SearchBots(s string, limit int, offset int) []BotRes
	SearchPeople(s string, limit int, offset int) []Person
	CreateLeaderBoard(uuid string, leaderboards []LeaderBoard) ([]LeaderBoard, error)
	GetLeaderBoard(uuid string) []LeaderBoard
	GetLeaderBoardByUuidAndAlias(uuid string, alias string) LeaderBoard
	UpdateLeaderBoard(uuid string, alias string, u map[string]interface{}) bool
	CountDevelopers() int64
	CountBounties() uint64
	GetPeopleListShort(count uint32) *[]PersonInShort
	GetConnectionCode() ConnectionCodesShort
	CreateConnectionCode(c []ConnectionCodes) ([]ConnectionCodes, error)
	GetLnUser(lnKey string) int64
	CreateLnUser(lnKey string) (Person, error)
	GetBountiesLeaderboard() []LeaderData
	GetWorkspaces(r *http.Request) []Organization
	GetWorkspacesCount() int64
	GetWorkspaceByUuid(uuid string) Organization
	GetWorkspaceByName(name string) Organization
	CreateOrEditWorkspace(m Organization) (Organization, error)
	GetWorkspaceUsers(uuid string) ([]OrganizationUsersData, error)
	GetWorkspaceUsersCount(uuid string) int64
	GetWorkspaceBountyCount(uuid string) int64
	GetWorkspaceUser(pubkey string, org_uuid string) OrganizationUsers
	CreateWorkspaceUser(orgUser OrganizationUsers) OrganizationUsers
	DeleteWorkspaceUser(orgUser OrganizationUsersData, org string) OrganizationUsersData
	GetBountyRoles() []BountyRoles
	CreateUserRoles(roles []UserRoles, uuid string, pubkey string) []UserRoles
	GetUserCreatedWorkspaces(pubkey string) []Organization
	GetUserAssignedWorkspaces(pubkey string) []OrganizationUsers
	AddBudgetHistory(budget BudgetHistory) BudgetHistory
	CreateWorkspaceBudget(budget BountyBudget) BountyBudget
	UpdateWorkspaceBudget(budget BountyBudget) BountyBudget
	GetPaymentHistoryByCreated(created *time.Time, org_uuid string) PaymentHistory
	GetWorkspaceBudget(org_uuid string) BountyBudget
	GetWorkspaceStatusBudget(org_uuid string) StatusBudget
	GetWorkspaceBudgetHistory(org_uuid string) []BudgetHistoryData
	AddAndUpdateBudget(invoice InvoiceList) PaymentHistory
	WithdrawBudget(sender_pubkey string, org_uuid string, amount uint)
	AddPaymentHistory(payment PaymentHistory) PaymentHistory
	GetPaymentHistory(org_uuid string, r *http.Request) []PaymentHistory
	GetInvoice(payment_request string) InvoiceList
	GetWorkspaceInvoices(org_uuid string) []InvoiceList
	GetWorkspaceInvoicesCount(org_uuid string) int64
	UpdateInvoice(payment_request string) InvoiceList
	AddInvoice(invoice InvoiceList) InvoiceList
	AddUserInvoiceData(userData UserInvoiceData) UserInvoiceData
	GetUserInvoiceData(payment_request string) UserInvoiceData
	DeleteUserInvoiceData(payment_request string) UserInvoiceData
	ChangeWorkspaceDeleteStatus(org_uuid string, status bool) Organization
	UpdateWorkspaceForDeletion(uuid string) error
	DeleteAllUsersFromWorkspace(uuid string) error
	GetFilterStatusCount() FilterStattuCount
	UserHasManageBountyRoles(pubKeyFromAuth string, uuid string) bool
	BountiesPaidPercentage(r PaymentDateRange) uint
	TotalSatsPosted(r PaymentDateRange) uint
	TotalSatsPaid(r PaymentDateRange) uint
	SatsPaidPercentage(r PaymentDateRange) uint
	AveragePaidTime(r PaymentDateRange) uint
	AverageCompletedTime(r PaymentDateRange) uint
	TotalBountiesPosted(r PaymentDateRange) int64
	TotalPaidBounties(r PaymentDateRange) int64
	NewHuntersPaid(r PaymentDateRange) int64
	TotalHuntersPaid(r PaymentDateRange) int64
	GetPersonByPubkey(pubkey string) Person
	GetBountiesByDateRange(r PaymentDateRange, re *http.Request) []Bounty
	GetBountiesByDateRangeCount(r PaymentDateRange, re *http.Request) int64
	GetBountiesProviders(r PaymentDateRange, re *http.Request) []Person
	PersonUniqueNameFromName(name string) (string, error)
	ProcessAlerts(p Person)
	UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool
}
