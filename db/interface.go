package db

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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
	UpdateBountyPaymentStatuses(bounty NewBounty) (NewBounty, error)
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
	GetConnectionCodesList(page, limit int) ([]ConnectionCodesList, int64, error)
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
	GetUserRoles(uuid string, pubkey string) []WorkspaceUserRoles
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
	ProcessBountyPayment(payment NewPaymentHistory, bounty NewBounty) error
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
	GetLastWithdrawal(workspace_uuid string) NewPaymentHistory
	GetSumOfDeposits(workspace_uuid string) uint
	GetSumOfWithdrawal(workspace_uuid string) uint
	DeleteAllUsersFromWorkspace(uuid string) error
	GetFilterStatusCount() FilterStatusCount
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
	GetPendingPaymentHistory() []NewPaymentHistory
	GetPaymentByBountyId(bountyId uint) NewPaymentHistory
	GetDailyEarnings() []DailyEarning
	SetPaymentAsComplete(tag string) bool
	SetPaymentStatusByBountyId(bountyId uint, tagResult V2TagRes) bool
	GetWorkspacePendingPayments(workspace_uuid string) []NewPaymentHistory
	CreateWorkflowRequest(req *WfRequest) error
	UpdateWorkflowRequest(req *WfRequest) error
	GetWorkflowRequestByID(requestID string) (*WfRequest, error)
	GetWorkflowRequestsByStatus(status WfRequestStatus) ([]WfRequest, error)
	GetWorkflowRequest(requestID string) (*WfRequest, error)
	UpdateWorkflowRequestStatusAndResponse(requestID string, status WfRequestStatus, responseData PropertyMap) error
	GetWorkflowRequestsByWorkflowID(workflowID string) ([]WfRequest, error)
	GetPendingWorkflowRequests(limit int) ([]WfRequest, error)
	DeleteWorkflowRequest(requestID string) error
	CreateProcessingMap(pm *WfProcessingMap) error
	UpdateProcessingMap(pm *WfProcessingMap) error
	GetProcessingMapByKey(processType, processKey string) (*WfProcessingMap, error)
	GetProcessingMapsByType(processType string) ([]WfProcessingMap, error)
	DeleteProcessingMapByKey(processType, processKey string) error
	DeleteProcessingMap(id uint) error
	ProcessReversePayments(paymentId uint) error
	CreateOrEditTicket(ticket *Tickets) (Tickets, error)
	GetTicketsByGroup(ticketGroupUUID string) ([]Tickets, error)
	GetTicket(uuid string) (Tickets, error)
	UpdateTicket(ticket Tickets) (Tickets, error)
	DeleteTicket(uuid string) error
	GetProductBrief(workspaceUuid string) (string, error)
	GetFeatureBrief(featureUuid string) (string, error)
	GetFeatureArchitecture(featureUuid string) (string, error)
	GetPhaseDesign(phaseUUID string) (string, error)
	GetTicketsByPhaseUUID(featureUUID string, phaseUUID string) ([]Tickets, error)
	AddChat(chat *Chat) (Chat, error)
	UpdateChat(chat *Chat) (Chat, error)
	GetChatByChatID(chatID string) (Chat, error)
	AddChatMessage(message *ChatMessage) (ChatMessage, error)
	UpdateChatMessage(message *ChatMessage) (ChatMessage, error)
	GetChatMessagesForChatID(chatID string) ([]ChatMessage, error)
	GetChatsForWorkspace(workspaceID string, chatStatus string, limit int, offset int) ([]Chat, int64, error)
	GetCodeGraphByUUID(uuid string) (WorkspaceCodeGraph, error)
	GetCodeGraphByWorkspaceUuid(workspace_uuid string) (WorkspaceCodeGraph, error)
	CreateOrEditCodeGraph(m WorkspaceCodeGraph) (WorkspaceCodeGraph, error)
	DeleteCodeGraph(workspace_uuid string, uuid string) error
	GetTicketsWithoutGroup() ([]Tickets, error)
	UpdateTicketsWithoutGroup(ticket Tickets) error
	ProcessUpdateTicketsWithoutGroup()
	GetNewHunters(r PaymentDateRange) int64
	TotalPeopleByPeriod(r PaymentDateRange) int64
	GetProofsByBountyID(bountyID uint) []ProofOfWork
	CreateProof(proof ProofOfWork) error
	DeleteProof(proofID string) error
	UpdateProofStatus(proofID string, status ProofOfWorkStatus) error
	IncrementProofCount(bountyID uint) error
	DecrementProofCount(bountyID uint) error
	CreateBountyTiming(bountyID uint) (*BountyTiming, error)
	GetBountyTiming(bountyID uint) (*BountyTiming, error)
	UpdateBountyTiming(timing *BountyTiming) error
	StartBountyTiming(bountyID uint) error
	CloseBountyTiming(bountyID uint) error
	UpdateBountyTimingOnProof(bountyID uint) error
	GetWorkspaceBountyCardsData(r *http.Request) []NewBounty
	UpdateFeatureStatus(uuid string, status FeatureStatus) (WorkspaceFeatures, error)
	CreateBountyFromTicket(ticket Tickets, pubkey string) (*NewBounty, error)
	AddFeatureFlag(flag *FeatureFlag) (FeatureFlag, error)
	UpdateFeatureFlag(flag *FeatureFlag) (FeatureFlag, error)
	DeleteFeatureFlag(flagUUID uuid.UUID) error
	GetFeatureFlags() ([]FeatureFlag, error)
	GetFeatureFlagByUUID(flagUUID uuid.UUID) (FeatureFlag, error)
	GetEndpointByUUID(uuid uuid.UUID) (Endpoint, error)
	AddEndpoint(endpoint *Endpoint) (Endpoint, error)
	UpdateEndpoint(endpoint *Endpoint) (Endpoint, error)
	DeleteEndpoint(endpointUUID uuid.UUID) error
	GetEndpointsByFeatureFlag(flagUUID uuid.UUID) ([]Endpoint, error)
	GetEndpointByPath(path string) (Endpoint, error)
	GetAllEndpoints() ([]Endpoint, error)
	GetLatestTicketByGroup(ticketGroup uuid.UUID) (Tickets, error)
	GetAllTicketGroups(workspaceUuid string) ([]uuid.UUID, error)
	GetFeaturedBountyById(id string) (FeaturedBounty, error)
	GetAllFeaturedBounties() ([]FeaturedBounty, error)
	CreateFeaturedBounty(bounty FeaturedBounty) error
	UpdateFeaturedBounty(bountyID string, bounty FeaturedBounty) error
	DeleteFeaturedBounty(bountyID string) error
	CreateNotification(notification *Notification) error
	GetNotification(uuid string) (*Notification, error)
	UpdateNotification(uuid string, updates map[string]interface{}) error
	DeleteNotification(uuid string) error
	GetPendingNotifications() ([]Notification, error)
	GetFailedNotifications(maxRetries int) ([]Notification, error)
	GetNotificationsByPubKey(pubKey string, limit, offset int) ([]Notification, error)
	IncrementRetryCount(uuid string) error
	GetNotificationCount(pubKey string) (int64, error)
	GetWorkspaceDraftTicket(workspaceUuid string, uuid string) (Tickets, error)
	CreateWorkspaceDraftTicket(ticket *Tickets) (Tickets, error)
	UpdateWorkspaceDraftTicket(ticket *Tickets) (Tickets, error)
	DeleteWorkspaceDraftTicket(workspaceUuid string, uuid string) error
	CreateSnippet(snippet *TextSnippet) (*TextSnippet, error)
	GetSnippetsByWorkspace(workspaceUUID string) ([]TextSnippet, error)
	GetSnippetByID(id uint) (*TextSnippet, error)
	UpdateSnippet(snippet *TextSnippet) (*TextSnippet, error)
	DeleteSnippet(id uint) error
	CreateFileAsset(asset *FileAsset) (*FileAsset, error)
	GetFileAssetByHash(fileHash string) (*FileAsset, error)
	GetFileAssetByID(id uint) (*FileAsset, error)
	UpdateFileAssetReference(id uint) error
	ListFileAssets(params ListFileAssetsParams) ([]FileAsset, int64, error)
	UpdateFileAsset(asset *FileAsset) error
	DeleteFileAsset(id uint) error
	DeleteBountyTiming(bountyID uint) error
	DeleteTicketGroup(TicketGroupUUID uuid.UUID) error
	PauseBountyTiming(bountyID uint) error
	ResumeBountyTiming(bountyID uint) error
	SaveNotification(pubkey, event, content, status string) error
	GetNotificationsByStatus(status string) []Notification
	IncrementNotificationRetry(notificationUUID string)
	UpdateNotificationStatus(notificationUUID string, status string)
	CreateOrEditTicketPlan(plan *TicketPlan) (*TicketPlan, error)
	GetTicketPlan(uuid string) (*TicketPlan, error)
	DeleteTicketPlan(uuid string) error
	GetTicketPlansByFeature(featureUUID string) ([]TicketPlan, error)
	GetTicketPlansByPhase(phaseUUID string) ([]TicketPlan, error)
	GetTicketPlansByWorkspace(workspaceUUID string) ([]TicketPlan, error)
	CreateActivity(activity *Activity) (*Activity, error)
	UpdateActivity(activity *Activity) (*Activity, error)
	GetActivity(id string) (*Activity, error)
	GetActivitiesByThread(threadID string) ([]Activity, error)
	GetActivitiesByFeature(featureUUID string) ([]Activity, error)
	GetActivitiesByPhase(phaseUUID string) ([]Activity, error)
	GetActivitiesByWorkspace(workspace string) ([]Activity, error)
	GetLatestActivityByThread(threadID string) (*Activity, error)
	CreateActivityThread(sourceID string, activity *Activity) (*Activity, error)
	DeleteActivity(id string) error
	BuildTicketArray(groupIDs []string) []TicketArrayItem
	GetBountiesByWorkspaceAndTimeRange(workspaceId string, startDate time.Time, endDate time.Time) ([]NewBounty, error)
	GetBountiesByFeatureUuid(featureUuid string) ([]NewBounty, error)
	GetTicketsByFeatureUUID(featureUuid string) ([]Tickets, error)
	CreateArtifact(artifact *Artifact) (*Artifact, error)
	GetArtifactByID(id uuid.UUID) (*Artifact, error)
	GetArtifactsByMessageID(messageID string) ([]Artifact, error)
	GetAllArtifactsByChatID(chatID string) ([]Artifact, error)
	UpdateArtifact(artifact *Artifact) (*Artifact, error)
	DeleteArtifactByID(id uuid.UUID) error
	DeleteAllArtifactsByChatID(chatID string) error
	CreateOrUpdateFeatureCall(workspaceID string, url string) (*FeatureCall, error)
	GetFeatureCallByWorkspaceID(workspaceID string) (*FeatureCall, error)
	DeleteFeatureCall(workspaceID string) error
	CreateOrEditChatWorkflow(workflow *ChatWorkflow) (*ChatWorkflow, error)
	GetChatWorkflowByWorkspaceID(workspaceID string) (*ChatWorkflow, error)
	DeleteChatWorkflow(workspaceID string) error
	CreateSkill(skill *Skill) (*Skill, error)
	GetAllSkills() ([]Skill, error)
	GetSkillByID(id uuid.UUID) (*Skill, error)
	UpdateSkillByID(skill *Skill) (*Skill, error)
	DeleteSkillByID(id uuid.UUID) error
	CreateSkillInstall(install *SkillInstall) (*SkillInstall, error)
	GetSkillInstallBySkillsID(skillID uuid.UUID) ([]SkillInstall, error)
	GetSkillInstallByID(id uuid.UUID) (*SkillInstall, error)
	UpdateSkillInstallByID(install *SkillInstall) (*SkillInstall, error)
	DeleteSkillInstallByID(id uuid.UUID) error
	CreateSSEMessageLog(event map[string]interface{}, chatID, from, to string) (*SSEMessageLog, error)
	DeleteSSEMessageLog(id uuid.UUID) error
	UpdateSSEMessageLogStatusBatch(ids []uuid.UUID) error
	UpdateSSEMessageLog(id uuid.UUID, updates map[string]interface{}) (*SSEMessageLog, error)
	GetSSEMessagesByChatID(chatID string, limit int, offset int, status string) ([]SSEMessageLog, int64, error)
	GetSSEMessageLogByID(id uuid.UUID) (*SSEMessageLog, error)
	GetSSEMessageLogsByChatID(chatID string) ([]SSEMessageLog, error)
	GetNewSSEMessageLogsByChatID(chatID string) ([]SSEMessageLog, error)
	CreateCodeSpaceMap(codeSpace CodeSpaceMap) (CodeSpaceMap, error)
	GetCodeSpaceMaps() ([]CodeSpaceMap, error)
	GetCodeSpaceMapByWorkspace(workspaceID string) ([]CodeSpaceMap, error)
	GetCodeSpaceMapByUser(userPubkey string) ([]CodeSpaceMap, error)
	GetCodeSpaceMapByURL(codeSpaceURL string) ([]CodeSpaceMap, error)
	GetCodeSpaceMapByWorkspaceAndUser(workspaceID, userPubkey string) (CodeSpaceMap, error)
	GetCodeSpaceMapByID(id uuid.UUID) (CodeSpaceMap, error)
	UpdateCodeSpaceMap(id uuid.UUID, updates map[string]interface{}) (CodeSpaceMap, error)
	DeleteCodeSpaceMap(id uuid.UUID) error
	CreateBountyStake(stake BountyStake) (*BountyStake, error)
	GetAllBountyStakes() ([]BountyStake, error)
	GetBountyStakesByBountyID(bountyID uint) ([]BountyStake, error)
	GetBountyStakeByID(stakeID uuid.UUID) (*BountyStake, error)
	GetBountyStakesByHunterPubKey(hunterPubKey string) ([]BountyStake, error)
	UpdateBountyStake(stakeID uuid.UUID, updates map[string]interface{}) (*BountyStake, error)
	DeleteBountyStake(stakeID uuid.UUID) error
	AddChatStatus(status *ChatWorkflowStatus) (ChatWorkflowStatus, error)
	UpdateChatStatus(status *ChatWorkflowStatus) (ChatWorkflowStatus, error)
	GetChatStatusByChatID(chatID string) ([]ChatWorkflowStatus, error)
	GetLatestChatStatusByChatID(chatID string) (ChatWorkflowStatus, error)
	DeleteChatStatus(uuid uuid.UUID) error
	DeleteOldSSEMessageLogs(maxAge time.Duration) (int64, error)
	CreateBountyStakeProcess(process *BountyStakeProcess) (*BountyStakeProcess, error)
	GetBountyStakeProcessByID(id uuid.UUID) (*BountyStakeProcess, error)
	GetBountyStakeProcessesByBountyID(bountyID uint) ([]BountyStakeProcess, error)
	GetBountyStakeProcessesByHunterPubKey(hunterPubKey string) ([]BountyStakeProcess, error)
	GetAllBountyStakeProcesses() ([]BountyStakeProcess, error)
	UpdateBountyStakeProcess(id uuid.UUID, updates map[string]interface{}) (*BountyStakeProcess, error)
	DeleteBountyStakeProcess(id uuid.UUID) error
}
