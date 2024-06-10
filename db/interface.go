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
	GetWorkspaceBounties(r *http.Request, workspace_uuid string) []NewBounty
	GetWorkspaceBountiesCount(r *http.Request, workspace_uuid string) int64
	GetAssignedBounties(r *http.Request) ([]NewBounty, error)
	GetCreatedBounties(r *http.Request) ([]NewBounty, error)
	GetBountyById(id string) ([]NewBounty, error)
	GetNextBountyByCreated(r *http.Request) (uint, error)
	GetPreviousBountyByCreated(r *http.Request) (uint, error)
	GetNextWorkspaceBountyByCreated(r *http.Request) (uint, error)
	GetPreviousWorkspaceBountyByCreated(r *http.Request) (uint, error)
	GetBountyIndexById(id string) int64
	GetBountyDataByCreated(created string) ([]NewBounty, error)
	AddBounty(b Bounty) (Bounty, error)
	GetAllBounties(r *http.Request) []NewBounty
	CreateOrEditBounty(b NewBounty) (NewBounty, error)
	UpdateBountyNullColumn(b NewBounty, column string) NewBounty
	UpdateBountyBoolColumn(b NewBounty, column string) NewBounty
	DeleteBounty(pubkey string, created string) (NewBounty, error)
	GetBountyByCreated(created uint) (NewBounty, error)
	GetBounty(id uint) NewBounty
	UpdateBounty(b NewBounty) (NewBounty, error)
	UpdateBountyPayment(b NewBounty) (NewBounty, error)
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
	GetWorkspaces(r *http.Request) []Workspace
	GetWorkspacesCount() int64
	GetWorkspaceByUuid(uuid string) Workspace
	GetWorkspaceByName(name string) Workspace
	CreateOrEditWorkspace(m Workspace) (Workspace, error)
	GetWorkspaceUsers(uuid string) ([]WorkspaceUsersData, error)
	GetWorkspaceUsersCount(uuid string) int64
	GetWorkspaceBountyCount(uuid string) int64
	GetWorkspaceUser(pubkey string, workspace_uuid string) WorkspaceUsers
	CreateWorkspaceUser(orgUser WorkspaceUsers) WorkspaceUsers
	DeleteWorkspaceUser(orgUser WorkspaceUsersData, org string) WorkspaceUsersData
	GetBountyRoles() []BountyRoles
	CreateUserRoles(roles []WorkspaceUserRoles, uuid string, pubkey string) []WorkspaceUserRoles
	GetUserCreatedWorkspaces(pubkey string) []Workspace
	GetUserAssignedWorkspaces(pubkey string) []WorkspaceUsers
	AddBudgetHistory(budget BudgetHistory) BudgetHistory
	CreateWorkspaceBudget(budget NewBountyBudget) NewBountyBudget
	UpdateWorkspaceBudget(budget NewBountyBudget) NewBountyBudget
	GetPaymentHistoryByCreated(created *time.Time, workspace_uuid string) NewPaymentHistory
	GetWorkspaceBudget(workspace_uuid string) NewBountyBudget
	GetWorkspaceStatusBudget(workspace_uuid string) StatusBudget
	GetWorkspaceBudgetHistory(workspace_uuid string) []BudgetHistoryData
	ProcessUpdateBudget(invoice NewInvoiceList) error
	AddAndUpdateBudget(invoice NewInvoiceList) NewPaymentHistory
	WithdrawBudget(sender_pubkey string, workspace_uuid string, amount uint)
	AddPaymentHistory(payment NewPaymentHistory) NewPaymentHistory
	GetPaymentHistory(workspace_uuid string, r *http.Request) []NewPaymentHistory
	GetInvoice(payment_request string) NewInvoiceList
	GetWorkspaceInvoices(workspace_uuid string) []NewInvoiceList
	GetWorkspaceInvoicesCount(workspace_uuid string) int64
	UpdateInvoice(payment_request string) NewInvoiceList
	AddInvoice(invoice NewInvoiceList) NewInvoiceList
	DeleteInvoice(payment_request string) NewInvoiceList
	AddUserInvoiceData(userData UserInvoiceData) UserInvoiceData
	ProcessAddInvoice(invoice NewInvoiceList, userData UserInvoiceData) error
	ProcessBudgetInvoice(paymentHistory NewPaymentHistory, newInvoice NewInvoiceList) error
	GetUserInvoiceData(payment_request string) UserInvoiceData
	DeleteUserInvoiceData(payment_request string) UserInvoiceData
	ChangeWorkspaceDeleteStatus(workspace_uuid string, status bool) Workspace
	UpdateWorkspaceForDeletion(uuid string) error
	ProcessDeleteWorkspace(workspace_uuid string) error
	DeleteAllUsersFromWorkspace(uuid string) error
	GetFilterStatusCount() FilterStattuCount
	UserHasManageBountyRoles(pubKeyFromAuth string, uuid string) bool
	BountiesPaidPercentage(r PaymentDateRange, workspace string) uint
	TotalSatsPosted(r PaymentDateRange, workspace string) uint
	TotalSatsPaid(r PaymentDateRange, workspace string) uint
	SatsPaidPercentage(r PaymentDateRange, workspace string) uint
	AveragePaidTime(r PaymentDateRange, workspace string) uint
	AverageCompletedTime(r PaymentDateRange, workspace string) uint
	TotalBountiesPosted(r PaymentDateRange, workspace string) int64
	TotalPaidBounties(r PaymentDateRange, workspace string) int64
	TotalAssignedBounties(r PaymentDateRange, workspace string) int64
	NewHuntersPaid(r PaymentDateRange, workspace string) int64
	TotalHuntersPaid(r PaymentDateRange, workspace string) int64
	GetPersonByPubkey(pubkey string) Person
	GetBountiesByDateRange(r PaymentDateRange, re *http.Request) []NewBounty
	GetBountiesByDateRangeCount(r PaymentDateRange, re *http.Request) int64
	GetBountiesProviders(r PaymentDateRange, re *http.Request) []Person
	PersonUniqueNameFromName(name string) (string, error)
	ProcessAlerts(p Person)
	UserHasAccess(pubKeyFromAuth string, uuid string, role string) bool
	CreateOrEditWorkspaceRepository(m WorkspaceRepositories) (WorkspaceRepositories, error)
	GetWorkspaceRepositorByWorkspaceUuid(uuid string) []WorkspaceRepositories
	GetWorkspaceRepoByWorkspaceUuidAndRepoUuid(workspace_uuid string, uuid string) (WorkspaceRepositories, error)
	DeleteWorkspaceRepository(workspace_uuid string, uuid string) bool
	CreateOrEditFeature(m WorkspaceFeatures) (WorkspaceFeatures, error)
	GetFeaturesByWorkspaceUuid(uuid string, r *http.Request) []WorkspaceFeatures
	GetWorkspaceFeaturesCount(uuid string) int64
	GetFeatureByUuid(uuid string) WorkspaceFeatures
	CreateOrEditFeaturePhase(phase FeaturePhase) (FeaturePhase, error)
	GetPhasesByFeatureUuid(featureUuid string) []FeaturePhase
	GetFeaturePhaseByUuid(featureUuid, phaseUuid string) (FeaturePhase, error)
	DeleteFeaturePhase(featureUuid, phaseUuid string) error
	CreateOrEditFeatureStory(story FeatureStory) (FeatureStory, error)
	GetFeatureStoriesByFeatureUuid(featureUuid string) ([]FeatureStory, error)
	GetFeatureStoryByUuid(featureUuid, storyUuid string) (FeatureStory, error)
	DeleteFeatureStoryByUuid(featureUuid, storyUuid string) error
	DeleteFeatureByUuid(uuid string) error
	GetBountiesByFeatureAndPhaseUuid(featureUuid string, phaseUuid string, r *http.Request) ([]NewBounty, error)
	GetBountiesCountByFeatureAndPhaseUuid(featureUuid string, phaseUuid string, r *http.Request) int64
	GetPhaseByUuid(phaseUuid string) (FeaturePhase, error)
	GetBountiesByPhaseUuid(phaseUuid string) []Bounty
	GetFeaturePhasesBountiesCount(bountyType string, phaseUuid string) int64
}
