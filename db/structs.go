package db

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Tribe struct
type Tribe struct {
	UUID            string         `json:"uuid"`
	OwnerPubKey     string         `json:"owner_pubkey"`
	OwnerAlias      string         `json:"owner_alias"`
	GroupKey        string         `json:"group_key"`
	Name            string         `json:"name"`
	UniqueName      string         `json:"unique_name"`
	Description     string         `json:"description"`
	Tags            pq.StringArray `gorm:"type:text[]" json:"tags"`
	Img             string         `json:"img"`
	PriceToJoin     int64          `json:"price_to_join"`
	PricePerMessage int64          `json:"price_per_message"`
	EscrowAmount    int64          `json:"escrow_amount"`
	EscrowMillis    int64          `json:"escrow_millis"`
	Created         *time.Time     `json:"created"`
	Updated         *time.Time     `json:"updated"`
	MemberCount     uint64         `json:"member_count"`
	Unlisted        bool           `json:"unlisted"`
	Private         bool           `json:"private"`
	Deleted         bool           `json:"deleted"`
	AppURL          string         `json:"app_url"`
	FeedURL         string         `json:"feed_url"`
	SecondBrainUrl  string         `json:"second_brain_url"`
	FeedType        uint64         `json:"feed_type"`
	LastActive      int64          `json:"last_active"`
	Bots            string         `json:"bots"`
	OwnerRouteHint  string         `json:"owner_route_hint"`
	Pin             string         `json:"pin"`
	Preview         string         `json:"preview"`
	ProfileFilters  string         `json:"profile_filters"` // "twitter,github"
	Badges          pq.StringArray `gorm:"type:text[]" json:"badges"`
}

// Bot struct
type Bot struct {
	UUID           string         `json:"uuid"`
	OwnerPubKey    string         `json:"owner_pubkey"`
	OwnerAlias     string         `json:"owner_alias"`
	Name           string         `json:"name"`
	UniqueName     string         `json:"unique_name"`
	Description    string         `json:"description"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags"`
	Img            string         `json:"img"`
	PricePerUse    int64          `json:"price_per_use"`
	Created        *time.Time     `json:"created"`
	Updated        *time.Time     `json:"updated"`
	Unlisted       bool           `json:"unlisted"`
	Deleted        bool           `json:"deleted"`
	MemberCount    uint64         `json:"member_count"`
	OwnerRouteHint string         `json:"owner_route_hint"`
	Tsv            string         `gorm:"type:tsvector"`
}

// Bot struct
type BotRes struct {
	UUID        string         `json:"uuid"`
	OwnerPubKey string         `json:"owner_pubkey"`
	Name        string         `json:"name"`
	UniqueName  string         `json:"unique_name"`
	Description string         `json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Img         string         `json:"img"`
	PricePerUse int64          `json:"price_per_use"`
}

// for bot pricing info
type BotInfo struct {
	Commands *[]BotCommand `json:"commands"`
	Prefix   string        `json:"prefix"`
	Price    int64         `json:"price"`
}
type BotCommand struct {
	Command   string `json:"command"`
	Price     int64  `json:"price"`
	MinPrice  int64  `json:"min_price"`
	MaxPrice  int64  `json:"max_price"`
	WordIndex uint   `json:"word_index"`
	AdminOnly bool   `json:"admin_only"`
}

type Tabler interface {
	TableName() string
}

// Person struct
type Person struct {
	ID               uint           `json:"id"`
	Uuid             string         `json:"uuid"`
	OwnerPubKey      string         `gorm:"uniqueIndex,unique" json:"owner_pubkey"`
	OwnerAlias       string         `json:"owner_alias"`
	UniqueName       string         `json:"unique_name"`
	Description      string         `json:"description"`
	Tags             pq.StringArray `gorm:"type:text[]" json:"tags" null`
	Img              string         `json:"img"`
	Created          *time.Time     `json:"created"`
	Updated          *time.Time     `json:"updated"`
	Unlisted         bool           `json:"unlisted"`
	Deleted          bool           `json:"deleted"`
	LastLogin        int64          `json:"last_login"`
	OwnerRouteHint   string         `json:"owner_route_hint"`
	OwnerContactKey  string         `json:"owner_contact_key"`
	PriceToMeet      int64          `json:"price_to_meet"`
	NewTicketTime    int64          `json:"new_ticket_time", gorm: "-:all"`
	TwitterConfirmed bool           `json:"twitter_confirmed"`
	ReferredBy       uint           `json:"referred_by"`
	Extras           PropertyMap    `json:"extras", type: jsonb not null default '{}'::jsonb`
	GithubIssues     PropertyMap    `json:"github_issues", type: jsonb not null default '{}'::jsonb`
}

type GormDataTypeInterface interface {
	GormDataType() string
}

type GormDBDataTypeInterface interface {
	GormDBDataType(*gorm.DB, *schema.Field) string
}

type StringArray pq.StringArray

func (StringArray) GormDataType() string {
	return `gorm:"type:text[]"`
}

func (p StringArray) Value() (driver.Value, error) {
	b := pq.StringArray(p)

	return b, nil
}

type PersonInShort struct {
	ID          uint   `json:"id"`
	Uuid        string `json:"uuid"`
	OwnerPubKey string `json:"owner_pubkey"`
	OwnerAlias  string `json:"owner_alias"`
	UniqueName  string `json:"unique_name"`
	Img         string `json:"img"`
}

// Github struct
type GithubIssue struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
	Description string `json:"description"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

type Channel struct {
	ID        uint       `json:"id"`
	TribeUUID string     `json:"tribe_uuid"`
	Name      string     `json:"name"`
	Created   *time.Time `json:"created"`
	Deleted   bool       `json:"deleted"`
}

type AssetTx struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	AssetId  uint   `json:"asset_id"`
	Amount   uint   `json:"amount"`
	Metadata string `json:"metadata"`
	Txid     string `json:"metadata"`
	Onchain  bool   `json:"onchain"`
}

type AssetResponse struct {
	Balances []AssetBalanceData `json:"balances"`
	Txs      []AssetTx          `json:"txs"`
}

type AssetBalanceData struct {
	OwnerPubkey string `json:"owner_pubkey"`
	AssetId     uint   `json:"asset_id"`
	Balance     uint   `json:"balance"`
}

type LeaderBoard struct {
	TribeUuid  string `json:"tribe_uuid"`
	Alias      string `json:"alias"`
	Spent      int64  `json:"spent"`
	Earned     int64  `json:"earned"`
	Reputation int64  `json:"reputation"`
}

type AssetListData struct {
	ID      uint   `json:"id"`
	Icon    string `json:"icon"`
	Name    string `json:"name"`
	Asset   string `json:"asset"`
	Token   string `json:"token"`
	Amount  uint   `json:"amount"`
	Creator string `json:"creator"`
	Balance uint   `json:"balance"`
}

type BadgeCreationData struct {
	Badge     string `json:"badge"`
	TribeUUID string `json:"tribeId"`
	Action    string `json:"action"`
}

type ConnectionCodes struct {
	ID               uint       `json:"id"`
	ConnectionString string     `json:"connection_string"`
	IsUsed           bool       `json:"is_used"`
	DateCreated      *time.Time `json:"date_created"`
}

type ConnectionCodesShort struct {
	ConnectionString string     `json:"connection_string"`
	DateCreated      *time.Time `json:"date_created"`
}

type ConnectionCodesList struct {
	ID               uint       `json:"id"`
	ConnectionString string     `json:"connection_string"`
	Pubkey           string     `json:"pubkey"`
	RouteHint        string     `json:"route_hint"`
	SatsAmount       uint64     `json:"sats_amount"`
	DateCreated      *time.Time `json:"date_created"`
	IsUsed           bool       `json:"is_used"`
}

type InvoiceRequest struct {
	Amount          string `json:"amount"`
	Memo            string `json:"memo"`
	Owner_pubkey    string `json:"owner_pubkey"`
	User_pubkey     string `json:"user_pubkey"`
	Created         string `json:"created"`
	Type            string `json:"type"`
	Assigned_hours  uint   `json:"assigned_hours,omitempty"`
	Commitment_fee  uint   `json:"commitment_fee,omitempty"`
	Bounty_expires  string `json:"bounty_expires,omitempty"`
	Websocket_token string `json:"websocket_token,omitempty"`
	Route_hint      string `json:"route_hint,omitempty"`
}

type Invoice struct {
	Invoice string `json:"invoice"`
}

type InvoiceResponse struct {
	Succcess bool    `json:"success"`
	Response Invoice `json:"response"`
}

type InvoiceStoreData struct {
	Invoice        string `json:"invoice"`
	Owner_pubkey   string `json:"owner_pubkey"`
	User_pubkey    string `json:"user_pubkey"`
	Amount         string `json:"amount"`
	Created        string `json:"created"`
	Host           string `json:"host,omitempty"`
	Type           string `json:"type"`
	Assigned_hours uint   `json:"assigned_hours,omitempty"`
	Commitment_fee uint   `json:"commitment_fee,omitempty"`
	Bounty_expires string `json:"bounty_expires,omitempty"`
	Route_hint     string `json:"route_hint,omitempty"`
}

type InvoiceStatus struct {
	Payment_request string `json:"payment_request"`
	Status          bool   `json:"Status"`
}

type InvoiceResult struct {
	Success  bool                 `json:"success"`
	Response InvoiceCheckResponse `json:"response"`
}

type InvoiceError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// TODO change amount back to string
type InvoiceCheckResponse struct {
	Settled         bool   `json:"settled"`
	Payment_request string `json:"payment_request"`
	Payment_hash    string `json:"payment_hash"`
	Preimage        string `json:"preimage"`
	Amount          string `json:"amount"`
}

type InvoicePaySuccess struct {
	Success  bool                 `json:"success"`
	Response InvoiceCheckResponse `json:"response"`
}

type InvoicePayError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type InvoiceSuccessResponse struct {
	Success  bool                     `json:"success"`
	Response InvoiceSuccessPaymentReq `json:"response"`
}

type InvoiceSuccessPaymentReq struct {
	Payment_request string `json:"payment_request"`
}

type DeleteBountyAssignee struct {
	Owner_pubkey string `json:"owner_pubkey"`
	Created      string `json:"created"`
}

type KeysendPayment struct {
	Amount          string `json:"amount"`
	Destination_key string `json:"destination_key"`
}

type KeysendSuccess struct {
	Success  bool        `json:"success"`
	Response PropertyMap `json:"response"`
}

type KeysendError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type LnHost struct {
	Msg  string `json:"msg"`
	Host string `json:"host"`
	K1   string `json:"k1"`
}

type LnEncode struct {
	Host string `json:"host"`
}

type BountyLeaderboard struct {
	Owner_pubkey             string `json:"owner_pubkey"`
	Total_bounties_completed uint   `json:"total_bounties_completed"`
	Total_sats_earned        uint   `json:"total_sats_earned"`
}

type YoutubeDownload struct {
	YoutubeUrls []string `json:"youtube_urls"`
}

type Client struct {
	Host string
	Conn *websocket.Conn
}

type AccessRestrictionType string

const (
	BlankAccess     AccessRestrictionType = ""
	WorkspaceAccess AccessRestrictionType = "workspace"
	OwnerAccess     AccessRestrictionType = "owner"
	AssignedAccess  AccessRestrictionType = "assigned"
	AdminAccess     AccessRestrictionType = "admins"
)

type Bounty struct {
	ID                      uint                   `json:"id"`
	OwnerID                 string                 `json:"owner_id"`
	Paid                    bool                   `json:"paid"`
	Show                    bool                   `gorm:"default:false" json:"show"`
	Completed               bool                   `gorm:"default:false" json:"completed"`
	Type                    string                 `json:"type"`
	Award                   string                 `json:"award"`
	AssignedHours           uint8                  `json:"assigned_hours"`
	BountyExpires           string                 `json:"bounty_expires"`
	CommitmentFee           uint64                 `json:"commitment_fee"`
	Price                   uint                   `json:"price"`
	Title                   string                 `json:"title"`
	Tribe                   string                 `json:"tribe"`
	Assignee                string                 `json:"assignee"`
	TicketUrl               string                 `json:"ticket_url"`
	OrgUuid                 string                 `json:"org_uuid"`
	Description             string                 `json:"description"`
	WantedType              string                 `json:"wanted_type"`
	Deliverables            string                 `json:"deliverables"`
	GithubDescription       bool                   `json:"github_description"`
	OneSentenceSummary      string                 `json:"one_sentence_summary"`
	EstimatedSessionLength  string                 `json:"estimated_session_length"`
	EstimatedCompletionDate string                 `json:"estimated_completion_date"`
	Created                 int64                  `json:"created"`
	Updated                 *time.Time             `json:"updated"`
	AssignedDate            *time.Time             `json:"assigned_date,omitempty"`
	CompletionDate          *time.Time             `json:"completion_date,omitempty"`
	MarkAsPaidDate          *time.Time             `json:"mark_as_paid_date,omitempty"`
	PaidDate                *time.Time             `json:"paid_date,omitempty"`
	CodingLanguages         pq.StringArray         `gorm:"type:text[];not null default:'[]'" json:"coding_languages"`
	PhaseUuid               *string                `json:"phase_uuid"`
	PhasePriority           *int                   `json:"phase_priority"`
	PaymentPending          bool                   `gorm:"default:false" json:"payment_pending"`
	PaymentFailed           bool                   `gorm:"default:false" json:"payment_failed"`
	AccessRestriction       *AccessRestrictionType `gorm:"type:varchar(20);default:null" json:"access_restriction,omitempty"`
	UnlockCode              *string                `gorm:"type:varchar(6);default:null;index" json:"unlock_code,omitempty"`
	IsStakable              bool                   `gorm:"default:false" json:"is_stakable"`
	StakeMin                uint                   `gorm:"default:0" json:"stake_min"`
	MaxStakers              int                    `gorm:"default:1" json:"max_stakers"`
	CurrentStakers          int                    `gorm:"default:0" json:"current_stakers"`
	Stakes                  []BountyStake          `gorm:"foreignKey:BountyID" json:"stakes,omitempty"`
}

// Todo: Change back to Bounty
type NewBounty struct {
	ID                      uint                   `json:"id"`
	OwnerID                 string                 `json:"owner_id"`
	Paid                    bool                   `json:"paid"`
	Show                    bool                   `gorm:"default:false" json:"show"`
	Completed               bool                   `gorm:"default:false" json:"completed"`
	Type                    string                 `json:"type"`
	Award                   string                 `json:"award"`
	AssignedHours           uint8                  `json:"assigned_hours"`
	BountyExpires           string                 `json:"bounty_expires"`
	CommitmentFee           uint64                 `json:"commitment_fee"`
	Price                   uint                   `json:"price"`
	Title                   string                 `json:"title"`
	Tribe                   string                 `json:"tribe"`
	Assignee                string                 `json:"assignee"`
	TicketUrl               string                 `json:"ticket_url"`
	OrgUuid                 string                 `gorm:"-" json:"org_uuid"`
	WorkspaceUuid           string                 `json:"workspace_uuid"`
	FeatureUuid             string                 `json:"feature_uuid"`
	Description             string                 `json:"description"`
	WantedType              string                 `json:"wanted_type"`
	Deliverables            string                 `json:"deliverables"`
	GithubDescription       bool                   `json:"github_description"`
	OneSentenceSummary      string                 `json:"one_sentence_summary"`
	EstimatedSessionLength  string                 `json:"estimated_session_length"`
	EstimatedCompletionDate string                 `json:"estimated_completion_date"`
	Created                 int64                  `json:"created"`
	Updated                 *time.Time             `json:"updated"`
	AssignedDate            *time.Time             `json:"assigned_date,omitempty"`
	CompletionDate          *time.Time             `json:"completion_date,omitempty"`
	MarkAsPaidDate          *time.Time             `json:"mark_as_paid_date,omitempty"`
	PaidDate                *time.Time             `json:"paid_date,omitempty"`
	CodingLanguages         pq.StringArray         `gorm:"type:text[];not null default:'[]'" json:"coding_languages"`
	PhaseUuid               string                 `json:"phase_uuid"`
	PhasePriority           int                    `json:"phase_priority"`
	PaymentPending          bool                   `gorm:"default:false" json:"payment_pending"`
	PaymentFailed           bool                   `gorm:"default:false" json:"payment_failed"`
	ProofOfWorkCount        int                    `gorm:"type:integer;default:0;not null" json:"pow"`
	AccessRestriction       *AccessRestrictionType `gorm:"type:varchar(20);default:null" json:"access_restriction,omitempty"`
	UnlockCode              *string                `gorm:"type:varchar(6);default:null;index" json:"unlock_code,omitempty"`
	IsStakable              bool                   `gorm:"default:false" json:"is_stakable"`
	StakeMin                uint                   `gorm:"default:0" json:"stake_min"`
	MaxStakers              int                    `gorm:"default:1" json:"max_stakers"`
	CurrentStakers          int                    `gorm:"default:0" json:"current_stakers"`
	Stakes                  []BountyStake          `gorm:"foreignKey:BountyID" json:"stakes,omitempty"`
}

type BountyOwners struct {
	OwnerID string `json:"owner_id"`
}

type BountyData struct {
	NewBounty
	BountyId          uint       `json:"bounty_id"`
	BountyCreated     int64      `json:"bounty_created"`
	BountyUpdated     *time.Time `json:"bounty_updated"`
	BountyDescription string     `json:"bounty_description"`
	Person
	AssigneeAlias           string         `json:"assignee_alias"`
	AssigneeId              uint           `json:"assignee_id"`
	AssigneeImg             string         `json:"assignee_img"`
	AssigneeCreated         *time.Time     `json:"assignee_created"`
	AssigneeUpdated         *time.Time     `json:"assignee_updated"`
	AssigneeDescription     string         `json:"assignee_description"`
	AssigneeRouteHint       string         `json:"assignee_route_hint"`
	BountyOwnerId           uint           `json:"bounty_owner_id"`
	OwnerUuid               string         `json:"owner_uuid"`
	OwnerKey                string         `json:"owner_key"`
	OwnerAlias              string         `json:"owner_alias"`
	OwnerUniqueName         string         `json:"owner_unique_name"`
	OwnerDescription        string         `json:"owner_description"`
	OwnerTags               pq.StringArray `gorm:"type:text[]" json:"owner_tags" null`
	OwnerImg                string         `json:"owner_img"`
	OwnerCreated            *time.Time     `json:"owner_created"`
	OwnerUpdated            *time.Time     `json:"owner_updated"`
	OwnerLastLogin          int64          `json:"owner_last_login"`
	OwnerRouteHint          string         `json:"owner_route_hint"`
	OwnerContactKey         string         `json:"owner_contact_key"`
	OwnerPriceToMeet        int64          `json:"owner_price_to_meet"`
	OwnerTwitterConfirmed   bool           `json:"owner_twitter_confirmed"`
	OrganizationName        string         `json:"organization_name"`
	OrganizationImg         string         `json:"organization_img"`
	OrganizationUuid        string         `json:"organization_uuid"`
	OrganizationDescription string         `json:"description"`
	WorkspaceName           string         `json:"workspace_name"`
	WorkspaceImg            string         `json:"workspace_img"`
	WorkspaceUuid           string         `json:"workspace_uuid"`
	WorkspaceDescription    string         `json:"workspace_description"`
}

type BountyResponse struct {
	Bounty       NewBounty      `json:"bounty"`
	Assignee     Person         `json:"assignee"`
	Owner        Person         `json:"owner"`
	Organization WorkspaceShort `json:"organization"`
	Workspace    WorkspaceShort `json:"workspace"`
	Proofs       []ProofOfWork  `json:"proofs,omitempty"`
	Pow          int            `json:"pow"`
}

type BountyCountResponse struct {
	OpenCount     int64 `json:"open_count"`
	AssignedCount int64 `json:"assigned_count"`
	PaidCount     int64 `json:"paid_count"`
}

type Organization struct {
	ID           uint       `json:"id"`
	Uuid         string     `json:"uuid"`
	Name         string     `gorm:"unique;not null" json:"name"`
	OwnerPubKey  string     `json:"owner_pubkey"`
	Img          string     `json:"img"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
	Show         bool       `json:"show"`
	Deleted      bool       `gorm:"default:false" json:"deleted"`
	BountyCount  int64      `json:"bounty_count,omitempty"`
	Budget       uint       `json:"budget,omitempty"`
	Website      string     `json:"website" validate:"omitempty,uri"`
	Github       string     `json:"github" validate:"omitempty,uri"`
	Description  string     `json:"description" validate:"omitempty,lte=120"`
	Mission      string     `json:"mission"`
	Tactics      string     `json:"tactics"`
	SchematicUrl string     `json:"schematic_url"`
	SchematicImg string     `json:"schematic_img"`
}

type Workspace struct {
	ID           uint       `json:"id"`
	Uuid         string     `json:"uuid"`
	Name         string     `gorm:"unique;not null" json:"name"`
	OwnerPubKey  string     `json:"owner_pubkey"`
	Img          string     `json:"img"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
	Show         bool       `json:"show"`
	Deleted      bool       `gorm:"default:false" json:"deleted"`
	BountyCount  int64      `json:"bounty_count,omitempty"`
	Budget       uint       `json:"budget,omitempty"`
	Website      string     `json:"website" validate:"omitempty,uri"`
	Github       string     `json:"github" validate:"omitempty,uri"`
	Description  string     `json:"description" validate:"omitempty,lte=120"`
	Mission      string     `json:"mission"`
	Tactics      string     `json:"tactics"`
	SchematicUrl string     `json:"schematic_url"`
	SchematicImg string     `json:"schematic_img"`
}

type WorkspaceShort struct {
	Uuid string `json:"uuid"`
	Name string `gorm:"unique;not null" json:"name"`
	Img  string `json:"img"`
}

type OrganizationUsers struct {
	ID          uint       `json:"id"`
	OwnerPubKey string     `json:"owner_pubkey"`
	OrgUuid     string     `json:"org_uuid"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
}

type WorkspaceUsers struct {
	ID            uint       `json:"id"`
	OwnerPubKey   string     `json:"owner_pubkey"`
	OrgUuid       string     `gorm:"-" json:"org_uuid"`
	WorkspaceUuid string     `json:"workspace_uuid,omitempty"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
}

type WorkspaceUsersData struct {
	OrgUuid       string     `gorm:"-" json:"org_uuid"`
	WorkspaceUuid string     `json:"workspace_uuid,omitempty"`
	UserCreated   *time.Time `json:"user_created"`
	Person
}

type WorkspaceRepositories struct {
	ID            uint       `json:"id"`
	Uuid          string     `gorm:"not null" json:"uuid"`
	WorkspaceUuid string     `gorm:"not null" json:"workspace_uuid"`
	Name          string     `gorm:"not null" json:"name"`
	Url           string     `json:"url"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
	CreatedBy     string     `json:"created_by"`
	UpdatedBy     string     `json:"updated_by"`
}

type WorkspaceCodeGraph struct {
	ID            uint       `json:"id"`
	Uuid          string     `gorm:"not null" json:"uuid"`
	WorkspaceUuid string     `gorm:"not null" json:"workspace_uuid"`
	Name          string     `gorm:"not null" json:"name"`
	Url           string     `json:"url"`
	SecretAlias   string     `json:"secret_alias"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
	CreatedBy     string     `json:"created_by"`
	UpdatedBy     string     `json:"updated_by"`
}

type FeatureStatus string

const (
	ActiveFeature    FeatureStatus = "active"
	ArchivedFeature  FeatureStatus = "archived"
	CompletedFeature FeatureStatus = "completed"
	BacklogFeature   FeatureStatus = "backlog"
)

type WorkspaceFeatures struct {
	ID                     uint          `json:"id"`
	Uuid                   string        `gorm:"unique;not null" json:"uuid"`
	WorkspaceUuid          string        `gorm:"not null" json:"workspace_uuid"`
	Name                   string        `gorm:"not null" json:"name"`
	Brief                  string        `json:"brief"`
	Requirements           string        `json:"requirements"`
	Architecture           string        `json:"architecture"`
	Url                    string        `json:"url"`
	Priority               int           `json:"priority"`
	Created                *time.Time    `json:"created"`
	Updated                *time.Time    `json:"updated"`
	CreatedBy              string        `json:"created_by"`
	UpdatedBy              string        `json:"updated_by"`
	BountiesCountCompleted int           `gorm:"-" json:"bounties_count_completed"`
	BountiesCountAssigned  int           `gorm:"-" json:"bounties_count_assigned"`
	BountiesCountOpen      int           `gorm:"-" json:"bounties_count_open"`
	FeatStatus             FeatureStatus `gorm:"type:varchar(20);default:'active';not null" json:"feat_status"`
}

type FeaturePhase struct {
	Uuid         string     `json:"uuid" gorm:"primary_key"`
	FeatureUuid  string     `json:"feature_uuid"`
	Name         string     `json:"name"`
	Priority     int        `json:"priority"`
	PhasePurpose string     `json:"phase_purpose,omitempty" gorm:"default:null"`
	PhaseOutcome string     `json:"phase_outcome,omitempty" gorm:"default:null"`
	PhaseScope   string     `json:"phase_scope,omitempty" gorm:"default:null"`
	PhaseDesign  string     `json:"phase_design,omitempty" gorm:"default:null"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
	CreatedBy    string     `json:"created_by"`
	UpdatedBy    string     `json:"updated_by"`
}

type BountyRoles struct {
	Name string `json:"name"`
}

type UserRoles struct {
	Role        string     `json:"role"`
	OwnerPubKey string     `json:"owner_pubkey"`
	OrgUuid     string     `json:"org_uuid"`
	Created     *time.Time `json:"created"`
}

// change back to UserRoles after migration
type WorkspaceUserRoles struct {
	Role          string     `json:"role"`
	OwnerPubKey   string     `json:"owner_pubkey"`
	OrgUuid       string     `gorm:"-" json:"org_uuid"`
	WorkspaceUuid string     `json:"workspace_uuid,omitempty"`
	Created       *time.Time `json:"created"`
}

type BountyBudget struct {
	ID            uint       `json:"id"`
	OrgUuid       string     `json:"org_uuid"`
	WorkspaceUuid string     `gorm:"-" json:"workspace_uuid,omitempty"`
	TotalBudget   uint       `json:"total_budget"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
}

// Rename back to BountyBudget
type NewBountyBudget struct {
	ID            uint       `json:"id"`
	OrgUuid       string     `gorm:"-" json:"org_uuid"`
	WorkspaceUuid string     `json:"workspace_uuid"`
	TotalBudget   uint       `json:"total_budget"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
}

type StatusBudget struct {
	OrgUuid             string `json:"org_uuid"`
	WorkspaceUuid       string `json:"workspace_uuid"`
	CurrentBudget       uint   `json:"current_budget"`
	OpenBudget          uint   `json:"open_budget"`
	OpenCount           int64  `json:"open_count"`
	OpenDifference      int    `json:"open_difference"`
	AssignedBudget      uint   `json:"assigned_budget"`
	AssignedCount       int64  `json:"assigned_count"`
	AssignedDifference  int    `json:"assigned_difference"`
	CompletedBudget     uint   `json:"completed_budget"`
	CompletedCount      int64  `json:"completed_count"`
	CompletedDifference int    `json:"completed_difference"`
}

type BudgetInvoiceRequest struct {
	Amount          uint        `json:"amount"`
	SenderPubKey    string      `json:"sender_pubkey"`
	OrgUuid         string      `json:"org_uuid,omitempty"`
	WorkspaceUuid   string      `json:"workspace_uuid,omitempty"`
	BountyID        uint        `json:"bounty_id,omitempty"`
	PaymentType     PaymentType `json:"payment_type,omitempty"`
	Websocket_token string      `json:"websocket_token,omitempty"`
}

type StakeInvoiceRequest struct {
	Amount         uint        `json:"amount"`
	SenderPubKey   string      `json:"sender_pubkey"`
	WorkspaceUuid  string      `json:"workspace_uuid"`
	PaymentType    PaymentType `json:"payment_type"`
	BountyID       uint        `json:"bounty_id,omitempty"`
	StakeOperation bool        `json:"is_stake,omitempty"`
}

type BudgetStoreData struct {
	Amount       uint       `json:"amount"`
	SenderPubKey string     `json:"sender_pubkey"`
	OrgUuid      string     `json:"org_uuid"`
	Invoice      string     `json:"invoice"`
	Host         string     `json:"host,omitempty"`
	Created      *time.Time `json:"created"`
}

type PaymentType string

const (
	Deposit  PaymentType = "deposit"
	Withdraw PaymentType = "withdraw"
	Payment  PaymentType = "payment"
	Reversal PaymentType = "reversal"
	Stake    PaymentType = "stake"
)

type BudgetHistory struct {
	ID           uint        `json:"id"`
	OrgUuid      string      `json:"org_uuid"`
	Amount       uint        `json:"amount"`
	SenderPubKey string      `json:"sender_pubkey"`
	Created      *time.Time  `json:"created"`
	Updated      *time.Time  `json:"updated"`
	Status       bool        `json:"status"`
	PaymentType  PaymentType `json:"payment_type"`
}

type FeatureStory struct {
	ID          uint       `json:"id"`
	Uuid        string     `json:"uuid"`
	FeatureUuid string     `json:"feature_uuid"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by"`
}

type BudgetHistoryData struct {
	BudgetHistory
	SenderName string `json:"sender_name"`
}

type PaymentHistory struct {
	ID             uint        `json:"id"`
	Amount         uint        `json:"amount"`
	BountyId       uint        `json:"bounty_id"`
	PaymentType    PaymentType `json:"payment_type"`
	OrgUuid        string      `json:"org_uuid"`
	SenderPubKey   string      `json:"sender_pubkey"`
	ReceiverPubKey string      `json:"receiver_pubkey"`
	Tag            string      `json:"tag,omitempty"`
	PaymentStatus  string      `json:"payment_status,omitempty"`
	Error          string      `json:"error,omitempty"`
	Created        *time.Time  `json:"created"`
	Updated        *time.Time  `json:"updated"`
	Status         bool        `json:"status"`
}

type NewPaymentHistory struct {
	ID             uint        `json:"id"`
	Amount         uint        `json:"amount"`
	BountyId       uint        `json:"bounty_id"`
	PaymentType    PaymentType `json:"payment_type"`
	OrgUuid        string      `gorm:"-" json:"org_uuid"`
	WorkspaceUuid  string      `json:"workspace_uuid,omitempty"`
	SenderPubKey   string      `json:"sender_pubkey"`
	ReceiverPubKey string      `json:"receiver_pubkey"`
	Tag            string      `json:"tag,omitempty"`
	PaymentStatus  string      `json:"payment_status,omitempty"`
	Error          string      `json:"error,omitempty"`
	Created        *time.Time  `json:"created"`
	Updated        *time.Time  `json:"updated"`
	Status         bool        `json:"status"`
}

type PaymentHistoryData struct {
	NewPaymentHistory
	SenderName   string `json:"sender_name"`
	ReceiverName string `json:"receiver_name"`
	SenderImg    string `json:"sender_img"`
	ReceiverImg  string `json:"receiver_img"`
}

type PaymentData struct {
	ID             uint        `json:"id"`
	OrgUuid        string      `json:"org_uuid"`
	PaymentType    PaymentType `json:"payment_type"`
	SenderName     string      `json:"sender_name"`
	SenderPubKey   string      `json:"sender_pubkey"`
	ReceiverName   string      `json:"receiver_name"`
	ReceiverPubKey string      `json:"receiver_pubkey"`
	Amount         uint        `json:"amount"`
	BountyId       uint        `json:"bounty_id"`
	Created        *time.Time  `json:"created"`
}

type BountyPayRequest struct {
	Websocket_token string `json:"websocket_token,omitempty"`
}

type InvoiceType string

const (
	Keysend    InvoiceType = "KEYSEND"
	Budget     InvoiceType = "BUDGET"
	PayInvoice InvoiceType = "ASSIGN"
)

type InvoiceList struct {
	ID             uint        `json:"id"`
	PaymentRequest string      `json:"payment_request"`
	Status         bool        `json:"status"`
	Type           InvoiceType `json:"type"`
	OwnerPubkey    string      `json:"owner_pubkey"`
	OrgUuid        string      `json:"org_uuid,omitempty"`
	Created        *time.Time  `json:"created"`
	Updated        *time.Time  `json:"updated"`
}

// Todo: Rename back to InvoiceList
type NewInvoiceList struct {
	ID             uint        `json:"id"`
	PaymentRequest string      `json:"payment_request"`
	Status         bool        `json:"status"`
	Type           InvoiceType `json:"type"`
	OwnerPubkey    string      `json:"owner_pubkey"`
	OrgUuid        string      `gorm:"-" json:"org_uuid"`
	WorkspaceUuid  string      `json:"workspace_uuid"`
	Created        *time.Time  `json:"created"`
	Updated        *time.Time  `json:"updated"`
}

type UserInvoiceData struct {
	ID             uint   `json:"id"`
	Amount         uint   `json:"amount"`
	PaymentRequest string `json:"payment_request"`
	Created        int    `json:"created"`
	UserPubkey     string `json:"user_pubkey"`
	AssignedHours  uint   `json:"assigned_hours,omitempty"`
	CommitmentFee  uint   `json:"commitment_fee,omitempty"`
	BountyExpires  string `json:"bounty_expires,omitempty"`
	RouteHint      string `json:"route_hint,omitempty"`
}

type WithdrawBudgetRequest struct {
	PaymentRequest  string `json:"payment_request"`
	Websocket_token string `json:"websocket_token,omitempty"`
	OrgUuid         string `json:"org_uuid"`
}

// change back to WithdrawBudgetReques
type NewWithdrawBudgetRequest struct {
	PaymentRequest  string `json:"payment_request"`
	Websocket_token string `json:"websocket_token,omitempty"`
	WorkspaceUuid   string `json:"workspace_uuid"`
}

type PaymentDateRange struct {
	StartDate   string      `json:"start_date"`
	EndDate     string      `json:"end_date"`
	PaymentType PaymentType `json:"payment_type,omitempty"`
}

type MemeChallenge struct {
	Id        string `json:"id"`
	Challenge string `json:"challenge"`
}

type SignerResponse struct {
	Sig string `json:"sig"`
}

type RelaySignerResponse struct {
	Success  bool           `json:"success"`
	Response SignerResponse `json:"response"`
}

type MemeTokenSuccess struct {
	Token string `json:"token"`
}

type Meme struct {
	Muid        string      `json:"muid"`
	OwnerPubKey string      `json:"owner_pub_key"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       int64       `json:"price"`
	Tags        StringArray `json:"tags"`
	Filename    string      `json:"filename"`
	Ttl         int64       `json:"ttl"`
	Size        int64       `json:"size"`
	Mime        string      `json:"mime"`
	Created     *time.Time  `json:"created"`
	Updated     *time.Time  `json:"updates"`
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	Template    bool        `json:"template"`
	Expiry      *time.Time  `json:"expiry"`
}

type DateDifference struct {
	Diff float64 `json:"diff"`
}

type BountyMetrics struct {
	BountiesPosted         int64 `json:"bounties_posted"`
	BountiesPaid           int64 `json:"bounties_paid"`
	BountiesAssigned       int64 `json:"bounties_assigned"`
	BountiesPaidPercentage uint  `json:"bounties_paid_average"`
	SatsPosted             uint  `json:"sats_posted"`
	SatsPaid               uint  `json:"sats_paid"`
	SatsPaidPercentage     uint  `json:"sats_paid_percentage"`
	AveragePaid            uint  `json:"average_paid"`
	AverageCompleted       uint  `json:"average_completed"`
	UniqueHuntersPaid      int64 `json:"unique_hunters_paid"`
	NewHuntersPaid         int64 `json:"new_hunters_paid"`
	NewHunters             int64 `json:"new_hunters"`
	NewHuntersByPeriod     int64 `json:"new_hunters_by_period"`
}

type MetricsBountyCsv struct {
	DatePosted   *time.Time `json:"date_posted"`
	Organization string     `json:"organization"`
	BountyAmount uint       `json:"bounty_amount"`
	Provider     string     `json:"provider"`
	Hunter       string     `json:"hunter"`
	BountyTitle  string     `json:"bounty_title"`
	BountyLink   string     `json:"bounty_link"`
	BountyStatus string     `json:"bounty_status"`
	DatePaid     *time.Time `json:"date_paid"`
	DateAssigned *time.Time `json:"date_assigned"`
}

type FilterStatusCount struct {
	Open      int64 `json:"open"`
	Assigned  int64 `json:"assigned"`
	Completed int64 `json:"completed"`
	Paid      int64 `json:"paid"`
	Pending   int64 `json:"pending"`
	Failed    int64 `json:"failed"`
}

type BountyStatus string

const (
	StatusTodo       BountyStatus = "TODO"
	StatusInProgress BountyStatus = "IN_PROGRESS"
	StatusInReview   BountyStatus = "IN_REVIEW"
	StatusComplete   BountyStatus = "COMPLETED"
	StatusPaid       BountyStatus = "PAID"
	StatusDraft      BountyStatus = "DRAFT"
)

type BountyCard struct {
	BountyID     uint              `json:"id"`
	TicketUUID   *uuid.UUID        `json:"ticket_uuid,omitempty"`
	TicketGroup  *uuid.UUID        `json:"ticket_group,omitempty"`
	Title        string            `json:"title"`
	AssigneePic  string            `json:"assignee_img,omitempty"`
	Assignee     string            `json:"assignee"`
	AssigneeName string            `json:"assignee_name"`
	Features     WorkspaceFeatures `json:"features"`
	Phase        FeaturePhase      `json:"phase"`
	Workspace    Workspace         `json:"workspace"`
	Status       BountyStatus      `json:"status"`
}

type WfRequestStatus string

const (
	StatusNew       WfRequestStatus = "NEW"
	StatusPending   WfRequestStatus = "PENDING"
	StatusCompleted WfRequestStatus = "COMPLETED"
	StatusFailed    WfRequestStatus = "FAILED"
)

type WfProcessingMap struct {
	ID                 uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Type               string      `gorm:"index;not null" json:"type"`
	ProcessKey         string      `gorm:"index;not null" json:"process_key"`
	RequiresProcessing bool        `gorm:"default:false" json:"requires_processing"`
	HandlerFunc        string      `json:"handler_func,omitempty"`
	Config             PropertyMap `gorm:"type:jsonb" json:"config,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

type WfRequest struct {
	ID           uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID    string          `gorm:"unique;not null" json:"request_id"`
	WorkflowID   string          `gorm:"index" json:"workflow_id"`
	Source       string          `gorm:"index" json:"source"`
	Action       string          `gorm:"index" json:"action"`
	Status       WfRequestStatus `json:"status"`
	ProjectID    string          `json:"project_id,omitempty"`
	RequestData  PropertyMap     `gorm:"type:jsonb" json:"request_data"`
	ResponseData PropertyMap     `gorm:"type:jsonb" json:"response_data,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type TicketStatus string

type Author string

const (
	DraftTicket      TicketStatus = "DRAFT"
	ReadyTicket      TicketStatus = "READY"
	InProgressTicket TicketStatus = "IN_PROGRESS"
	TestTicket       TicketStatus = "TEST"
	DeployTicket     TicketStatus = "DEPLOY"
	PayTicket        TicketStatus = "PAY"
	CompletedTicket  TicketStatus = "COMPLETED"
)

const (
	HumanAuthor Author = "HUMAN"
	AgentAuthor Author = "AGENT"
)

type Category string

const (
	WebDevelopment     Category = "Web development"
	BackendDevelopment Category = "Backend development"
	Design             Category = "Design"
	Other              Category = "Other"
)

type Tickets struct {
	UUID          uuid.UUID         `gorm:"primaryKey;type:uuid"`
	TicketGroup   *uuid.UUID        `gorm:"type:uuid;index:group_index" json:"ticket_group,omitempty"`
	WorkspaceUuid string            `gorm:"type:varchar(255);index:workspace_index" json:"workspace_uuid"`
	FeatureUUID   string            `gorm:"type:varchar(255);index:composite_index;default:null" json:"feature_uuid"`
	Features      WorkspaceFeatures `gorm:"foreignKey:FeatureUUID;references:Uuid;constraint:OnDelete:SET NULL"`
	PhaseUUID     string            `gorm:"type:varchar(255);index:phase_index;default:null" json:"phase_uuid"`
	FeaturePhase  FeaturePhase      `gorm:"foreignKey:PhaseUUID;references:Uuid;constraint:OnDelete:SET NULL"`
	Name          string            `gorm:"type:varchar(255)" json:"name"`
	Sequence      int               `gorm:"type:integer;index:composite_index;default:0" json:"sequence"`
	Dependency    []int             `gorm:"type:integer[]" json:"dependency"`
	Description   string            `gorm:"type:text" json:"description"`
	Status        TicketStatus      `gorm:"type:varchar(50);default:'DRAFT'" json:"status"`
	Version       int               `gorm:"type:integer;default:0" json:"version"`
	Author        *Author           `gorm:"type:varchar(50)" json:"author,omitempty"`
	AuthorID      *string           `gorm:"type:varchar(255)" json:"author_id,omitempty"`
	Amount        *int64            `gorm:"type:bigint;default:null" json:"amount,omitempty"`
	Category      *Category         `gorm:"type:varchar(50);default:null" json:"category,omitempty"`
	Mode          string            `json:"mode,omitempty"`
	CreatedAt     time.Time         `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

type TicketArrayItem struct {
	TicketName        string `json:"ticket_name"`
	TicketDescription string `json:"ticket_description"`
}

type BroadcastType string

const (
	PoolBroadcast BroadcastType = "pool"

	DirectBroadcast BroadcastType = "direct"
)

type ActionType string

const (
	ProcessAction ActionType = "process"

	MessageAction ActionType = "message"

	RunLinkAction ActionType = "run-link"
)

type TicketMessage struct {
	BroadcastType   BroadcastType `json:"broadcastType"`
	SourceSessionID string        `json:"sourceSessionID"`
	Message         string        `json:"message"`
	Action          ActionType    `json:"action"`
	TicketDetails   Tickets       `json:"ticketDetails"`
}

type ContextTagType string

const (
	ProductBriefContext ContextTagType = "productBrief"
	FeatureBriefContext ContextTagType = "featureBrief"
	SchematicContext    ContextTagType = "schematic"
)

type ContextTag struct {
	Type ContextTagType `json:"type"`
	ID   string         `json:"id"`
}

type ChatRole string

const (
	UserRole      ChatRole = "user"
	AssistantRole ChatRole = "assistant"
)

type ChatMessageStatus string

const (
	SendingStatus ChatMessageStatus = "sending"
	SentStatus    ChatMessageStatus = "sent"
	ErrorStatus   ChatMessageStatus = "error"
)

type ChatSource string

const (
	UserSource  ChatSource = "user"
	AgentSource ChatSource = "agent"
)

type ChatMessage struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	ChatID      string            `json:"chatId" gorm:"index"`
	Message     string            `json:"message"`
	PDFURL      string            `json:"pdf_url,omitempty"`
	Role        ChatRole          `json:"role"`
	Timestamp   time.Time         `json:"timestamp"`
	ContextTags []ContextTag      `json:"contextTags" gorm:"type:jsonb"`
	Status      ChatMessageStatus `json:"status"`
	Source      ChatSource        `json:"source"`
}

type ChatStatus string

const (
	ActiveStatus  ChatStatus = "active"
	ArchiveStatus ChatStatus = "archived"
)

type Chat struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	WorkspaceID string     `json:"workspaceId" gorm:"index"`
	Title       string     `json:"title"`
	Status      ChatStatus `json:"status" gorm:"default:active"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ChatWorkflowStatus struct {
	UUID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"uuid"`
	ChatID    string    `gorm:"index;not null" json:"chat_id"`
	Status    string    `gorm:"type:varchar(255);not null" json:"status"`
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

type ProofOfWorkStatus string

const (
	NewStatus             ProofOfWorkStatus = "New"
	AcceptedStatus        ProofOfWorkStatus = "Accepted"
	RejectedStatus        ProofOfWorkStatus = "Rejected"
	ChangeRequestedStatus ProofOfWorkStatus = "Change Requested"
)

type ProofOfWork struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey"`
	BountyID    uint              `json:"bounty_id"`
	Description string            `json:"description" gorm:"type:text;not null"`
	Status      ProofOfWorkStatus `json:"status" gorm:"type:varchar(20);default:'New'"`
	CreatedAt   time.Time         `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	SubmittedAt time.Time         `json:"submitted_at" gorm:"type:timestamp;default:current_timestamp"`
}

type BountyTiming struct {
	ID                      uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	BountyID                uint       `json:"bounty_id" gorm:"not null"`
	TotalWorkTimeSeconds    int        `json:"total_work_time_seconds" gorm:"default:0"`
	TotalDurationSeconds    int        `json:"total_duration_seconds" gorm:"default:0"`
	TotalAttempts           int        `json:"total_attempts" gorm:"default:0"`
	FirstAssignedAt         *time.Time `json:"first_assigned_at"`
	LastPoWAt               *time.Time `json:"last_pow_at"`
	ClosedAt                *time.Time `json:"closed_at"`
	IsPaused                bool       `json:"is_paused" gorm:"default:false"`
	LastPausedAt            *time.Time `json:"last_paused_at"`
	AccumulatedPauseSeconds int        `json:"accumulated_pause_seconds" gorm:"default:0"`
	CreatedAt               time.Time  `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt               time.Time  `json:"updated_at" gorm:"default:current_timestamp"`
}

type FeatureFlag struct {
	UUID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"uuid"`
	Name        string     `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Enabled     bool       `gorm:"type:boolean;default:false" json:"enabled"`
	Endpoints   []Endpoint `gorm:"foreignKey:FeatureFlagUUID" json:"endpoints,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"-"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"-"`
}

type Endpoint struct {
	UUID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"uuid"`
	Path            string    `gorm:"type:varchar(255);not null" json:"path"`
	FeatureFlagUUID uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"-"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:current_timestamp" json:"-"`
}

type FeaturedBounty struct {
	BountyID  string    `json:"bountyId" gorm:"uniqueIndex; unique; not null"`
	URL       string    `gorm:"not null" json:"url"`
	AddedAt   int64     `json:"addedAt"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}

type NotificationStatus string

const (
	NotificationStatusPending            NotificationStatus = "PENDING"
	NotificationStatusComplete           NotificationStatus = "COMPLETE"
	NotificationStatusFailed             NotificationStatus = "FAILED"
	NotificationStatusWaitingKeyExchange NotificationStatus = "WAITING_KEY_EXCHANGE"
)

type Notification struct {
	ID        uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID      string             `json:"uuid" gorm:"type:uuid;uniqueIndex;not null"`
	Event     string             `json:"event" gorm:"type:varchar(50);not null"`
	PubKey    string             `json:"pub_key" gorm:"type:varchar(100);not null;index"`
	Content   string             `json:"content" gorm:"type:text;not null"`
	Retries   int                `json:"retries" gorm:"default:0"`
	Status    NotificationStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt *time.Time         `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt *time.Time         `json:"updated_at" gorm:"default:current_timestamp"`
}

type TextSnippet struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	WorkspaceUUID string    `json:"workspace_uuid" gorm:"type:varchar(255);not null;index"`
	Title         string    `json:"title" gorm:"type:varchar(255);not null"`
	Snippet       string    `json:"snippet" gorm:"type:text;not null"`
	DateCreated   time.Time `json:"date_created" gorm:"autoCreateTime"`
	LastEdited    time.Time `json:"last_edited" gorm:"autoUpdateTime"`
}

type FileStatus string

const (
	ActiveFileStatus   FileStatus = "active"
	ArchivedFileStatus FileStatus = "archived"
	DeletedFileStatus  FileStatus = "deleted"
)

type FileAsset struct {
	ID             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	OriginFilename string     `json:"originFilename"`
	FileHash       string     `json:"fileHash" gorm:"index"`
	UploadFilename string     `json:"uploadFilename" gorm:"uniqueIndex"`
	UploadTime     time.Time  `json:"uploadTime"`
	LastReferenced time.Time  `json:"lastReferenced"`
	FileSize       int64      `json:"fileSize"`
	MimeType       string     `json:"mimeType"`
	Status         FileStatus `json:"status" gorm:"type:varchar(20);default:'active'"`
	UploadedBy     string     `json:"uploadedBy"`
	StoragePath    string     `json:"storagePath"`
	WorkspaceID    string     `json:"workspaceId" gorm:"index"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

type ListFileAssetsParams struct {
	Status             *FileStatus `form:"status"`
	MimeType           *string     `form:"mimeType"`
	BeforeDate         *time.Time  `form:"beforeDate"`
	AfterDate          *time.Time  `form:"afterDate"`
	LastAccessedBefore *time.Time  `form:"lastAccessedBefore"`
	WorkspaceID        *string     `form:"workspaceId"`
	Page               int         `form:"page,default=1"`
	PageSize           int         `form:"pageSize,default=50"`
}

type PlanStatus string

const (
	DraftPlan    PlanStatus = "DRAFT"
	ApprovedPlan PlanStatus = "APPROVED"
)

type TicketPlan struct {
	UUID          uuid.UUID         `gorm:"primaryKey;type:uuid"`
	WorkspaceUuid string            `gorm:"type:varchar(255);index:workspace_index" json:"workspace_uuid"`
	FeatureUUID   string            `gorm:"type:varchar(255);index:composite_index;default:null" json:"feature_uuid"`
	Features      WorkspaceFeatures `gorm:"foreignKey:FeatureUUID;references:Uuid;constraint:OnDelete:SET NULL"`
	PhaseUUID     string            `gorm:"type:varchar(255);index:phase_index;default:null" json:"phase_uuid"`
	FeaturePhase  FeaturePhase      `gorm:"foreignKey:PhaseUUID;references:Uuid;constraint:OnDelete:SET NULL"`
	Name          string            `gorm:"type:varchar(255);not null" json:"name"`
	Description   string            `gorm:"type:text" json:"description"`
	TicketGroups  pq.StringArray    `gorm:"type:uuid[];not null;default:'{}'" json:"ticket_groups"`
	Status        PlanStatus        `gorm:"type:varchar(50);default:'DRAFT'" json:"status"`
	Version       int               `gorm:"type:integer;default:0" json:"version"`
	CreatedBy     string            `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedBy     string            `gorm:"type:varchar(255)" json:"updated_by"`
	CreatedAt     time.Time         `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

type StubTicket struct {
	TicketName        string `json:"ticketName"`
	TicketDescription string `json:"ticketDescription"`
	Reasoning         string `json:"reasoning"`
}

type PhasePlan struct {
	PhaseApproach string       `json:"phaseApproach"`
	StubTickets   []StubTicket `json:"stubTickets"`
}

type TicketPlanReviewRequest struct {
	Value struct {
		FeatureUUID string    `json:"featureUUID"`
		PhaseUUID   string    `json:"phaseUUID"`
		PhasePlan   PhasePlan `json:"phasePlan"`
	} `json:"value"`
	RequestUUID     string `json:"requestUUID"`
	SourceWebsocket string `json:"sourceWebsocket"`
}

type TicketPlanReviewResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

type AuthorType string

const (
	HumansAuthor AuthorType = "human"
	HiveAuthor   AuthorType = "hive"
)

type ContentType string

const (
	FeatureCreation   ContentType = "feature_creation"
	StoryUpdate       ContentType = "story_update"
	RequirementChange ContentType = "requirement_change"
	GeneralUpdate     ContentType = "general_update"
)

type Activity struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid"`
	ThreadID    uuid.UUID      `gorm:"type:uuid;index:thread_index" json:"thread_id"`
	Title       string         `gorm:"type:varchar(200)" json:"title,omitempty"`
	Sequence    int            `gorm:"type:integer;not null" json:"sequence"`
	ContentType ContentType    `gorm:"type:varchar(50);not null" json:"content_type"`
	Content     string         `gorm:"type:text;not null;check:content,length(content) <= 10000" json:"content"`
	Workspace   string         `gorm:"type:varchar(255);index:workspace_index" json:"workspace"`
	FeatureUUID string         `gorm:"type:varchar(255);index:feature_index" json:"feature_uuid"`
	PhaseUUID   string         `gorm:"type:varchar(255);index:phase_index" json:"phase_uuid"`
	Feedback    string         `gorm:"type:text" json:"feedback"`
	Actions     pq.StringArray `gorm:"type:text[];default:'{}'" json:"actions"`
	Questions   pq.StringArray `gorm:"type:text[];default:'{}'" json:"questions"`
	TimeCreated time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"time_created"`
	TimeUpdated time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"time_updated"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	Author      AuthorType     `gorm:"type:varchar(10);not null" json:"author"`
	AuthorRef   string         `gorm:"type:varchar(255);not null" json:"author_ref"`
}

type NodeData struct {
	BountyID    uint   `json:"bounty_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type QuickBountyItem struct {
	BountyID      uint         `json:"bountyID"`
	BountyTitle   string       `json:"bountyTitle"`
	Status        BountyStatus `json:"status"`
	AssignedAlias *string      `json:"assignedAlias"`
	PhaseID       *string      `json:"phaseID"`
}

type QuickBountiesResponse struct {
	FeatureID string                       `json:"featureID"`
	Phases    map[string][]QuickBountyItem `json:"phases"`
	Unphased  []QuickBountyItem            `json:"unphased"`
}

type QuickTicketItem struct {
	TicketUUID    uuid.UUID    `json:"ticketUUID"`
	TicketTitle   string       `json:"ticketTitle"`
	Status        BountyStatus `json:"status"`
	AssignedAlias *string      `json:"assignedAlias"`
	PhaseID       *string      `json:"phaseID"`
}

type QuickTicketsResponse struct {
	FeatureID string                       `json:"featureID"`
	Phases    map[string][]QuickTicketItem `json:"phases"`
	Unphased  []QuickTicketItem            `json:"unphased"`
}

type Node struct {
	NodeType string   `json:"node_type"`
	NodeData NodeData `json:"node_data"`
}

type NodeListResponse struct {
	NodeList []Node `json:"node_list"`
}

type ArtifactType string

const (
	TextArtifact   ArtifactType = "text"
	VisualArtifact ArtifactType = "visual"
	ActionArtifact ArtifactType = "action"
	SSEArtifact    ArtifactType = "sse_connection"
)

type Artifact struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	MessageID string       `json:"message_id" gorm:"index"`
	Type      ArtifactType `json:"type" gorm:"type:varchar(20);not null"`
	Content   PropertyMap  `json:"content" gorm:"type:jsonb;not null;default:'{}'::jsonb"`
	CreatedAt time.Time    `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"type:timestamp;default:current_timestamp"`
}

type TextContent struct {
	TextType string `json:"text_type"`
	Content  string `json:"content"`
}

type VisualContent struct {
	TextType string    `json:"text_type,omitempty"`
	URL      string    `json:"url,omitempty"`
	Examples []Example `json:"examples,omitempty"`
}

type Example struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ActionContent struct {
	ActionText string   `json:"action_text"`
	Options    []Option `json:"options"`
}

type Option struct {
	ActionType     string `json:"action_type"`
	OptionLabel    string `json:"option_label"`
	OptionResponse string `json:"option_response"`
	Webhook        string `json:"webhook"`
}

type FeatureCall struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkspaceID string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"workspace_id"`
	URL         string         `gorm:"type:text" json:"url"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ChatWorkflow struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WorkspaceID string    `json:"workspaceId" gorm:"index;not null"`
	URL         string    `json:"url" gorm:"type:text;not null"`
	StackworkID string    `json:"stackworkId" gorm:"column:stackwork_id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ChargeModel string

const (
	FreeChargeModel ChargeModel = "Free"
	PAYGChargeModel ChargeModel = "PAYG"
)

type SkillStatus string

const (
	ApprovedSkillStatus SkillStatus = "Approved"
	DraftSkillStatus    SkillStatus = "Draft"
	ArchivedSkillStatus SkillStatus = "Archived"
)

type ClientType string

const (
	ClaudeDesktopClient ClientType = "Claude Desktop"
	CursorClient        ClientType = "Cursor"
	ClineClient         ClientType = "Cline"
	GooseClient         ClientType = "Goose"
)

type Skill struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Tagline     string         `gorm:"type:varchar(255)" json:"tagline"`
	Description string         `gorm:"type:text" json:"description"`
	IconURL     string         `gorm:"type:varchar(255)" json:"icon_url"`
	OwnerPubkey string         `gorm:"type:varchar(255);not null" json:"owner_pubkey"`
	ChargeModel ChargeModel    `gorm:"type:varchar(50);not null" json:"charge_model"`
	Labels      pq.StringArray `gorm:"type:text[]" json:"labels"`
	Status      SkillStatus    `gorm:"type:varchar(50);not null;default:'Draft'" json:"status"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

type SkillInstall struct {
	ID                 uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SkillID            uuid.UUID  `gorm:"type:uuid;not null;index:skill_index" json:"skill_id"`
	Skill              Skill      `gorm:"foreignKey:SkillID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Client             ClientType `gorm:"type:varchar(50);not null" json:"client"`
	InstallDescription string     `gorm:"type:text" json:"install_description"`
	InstallFile        string     `gorm:"type:varchar(255)" json:"install_file"`
	CreatedAt          time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

type SSEMessageStatus string

const (
	SSEStatusNew  SSEMessageStatus = "new"
	SSEStatusSent SSEMessageStatus = "sent"
)

type SSEMessageLog struct {
	ID        uuid.UUID        `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt time.Time        `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time        `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	Event     PropertyMap      `gorm:"type:jsonb;not null;default:'{}'::jsonb" json:"event"`
	ChatID    string           `gorm:"index;not null" json:"chat_id"`
	From      string           `gorm:"not null" json:"from"`
	To        string           `gorm:"not null" json:"to"`
	Status    SSEMessageStatus `gorm:"type:varchar(10);default:'new'" json:"status"`
}

type CodeSpaceMap struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	WorkspaceID  string    `json:"workspaceID" gorm:"index"`
	CodeSpaceURL string    `json:"codeSpaceURL"`
	UserPubkey   string    `json:"userPubkey" gorm:"index"`
	Username     string    `json:"username,omitempty"`
	// TODO: Encrypt this field
	GithubPat    string    `json:"githubPat,omitempty" gorm:"column:github_pat"`
	BaseBranch string `json:"baseBranch"`
}

type StakeStatus string

const (
	StakeStatusNew       StakeStatus = "NEW"
	StakeStatusPending   StakeStatus = "PENDING"
	StakeStatusActive    StakeStatus = "ACTIVE"
	StakeStatusCompleted StakeStatus = "COMPLETED"
	StakeStatusReturned  StakeStatus = "RETURNED"
	StakeStatusFailed    StakeStatus = "FAILED"
)

type BountyStake struct {
	ID           uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BountyID     uint        `json:"bounty_id" gorm:"index;not null"`
	HunterPubKey string      `json:"hunter_pubkey" gorm:"type:varchar(255);not null"`
	Amount       uint        `json:"amount" gorm:"not null"`
	Status       StakeStatus `json:"status" gorm:"type:varchar(20);default:'NEW'"`
	Invoice      string      `json:"invoice" gorm:"type:text"`
	StakeReceipt string      `json:"stake_receipt" gorm:"type:text"`
	StakeReturn  string      `json:"stake_return" gorm:"type:text"`
	Note         string      `json:"note" gorm:"type:text"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	StakedAt     *time.Time  `json:"staked_at"`
	ReturnedAt   *time.Time  `json:"returned_at"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type BountyStakeProcess struct {
	ID           uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BountyID     uint               `json:"bounty_id" gorm:"index;not null"`
	HunterPubKey string             `json:"hunter_pubkey" gorm:"type:varchar(255);not null;index"`
	Amount       uint               `json:"amount" gorm:"not null"`
	Status       StakeProcessStatus `json:"status" gorm:"type:varchar(20);default:'NEW'"`
	Invoice      string             `json:"invoice" gorm:"type:text"` 
	StakeReceipt string             `json:"stake_receipt" gorm:"type:text"`
	StakeReturn  string             `json:"stake_return" gorm:"type:text"`
	StakedAt     *time.Time         `json:"staked_at"`
	ReturnedAt   *time.Time         `json:"returned_at"`
	CreatedAt    time.Time          `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt    time.Time          `json:"updated_at" gorm:"type:timestamp;default:current_timestamp"`
}

type StakeProcessStatus string

const (
	StakeProcessStatusNew     StakeProcessStatus = "NEW"
	StakeProcessStatusPending StakeProcessStatus = "PENDING"
	StakeProcessStatusPaid    StakeProcessStatus = "PAID"
	StakeProcessStatusFailed  StakeProcessStatus = "FAILED"
	StakeProcessStatusReturned StakeProcessStatus = "RETURNED"
)

func (Person) TableName() string {
	return "people"
}

func (PersonInShort) TableName() string {
	return "people"
}

func (Bounty) TableName() string {
	return "bounty"
}

func (NewBounty) TableName() string {
	return "bounty"
}

func (NewBountyBudget) TableName() string {
	return "bounty_budgets"
}

func (NewInvoiceList) TableName() string {
	return "invoice_lists"
}

func (NewPaymentHistory) TableName() string {
	return "payment_histories"
}

func (ConnectionCodes) TableName() string {
	return "connectioncodes"
}

func (ConnectionCodesShort) TableName() string {
	return "connectioncodes"
}

func (WfProcessingMap) TableName() string {
	return "wf_processing_maps"
}

func (WfRequest) TableName() string {
	return "wf_requests"
}

// PropertyMap ...
type PropertyMap map[string]interface{}

// Value ...
func (p PropertyMap) Value() (driver.Value, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(p)
	return b, err
}

// Scan ...
func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}
	return nil
}

type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
