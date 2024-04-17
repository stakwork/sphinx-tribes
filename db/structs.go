package db

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

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
	Tags           pq.StringArray ` `
	Img            string         `json:"img"`
	PricePerUse    int64          `json:"price_per_use"`
	Created        *time.Time     `json:"created"`
	Updated        *time.Time     `json:"updated"`
	Unlisted       bool           `json:"unlisted"`
	Deleted        bool           `json:"deleted"`
	MemberCount    uint64         `json:"member_count"`
	OwnerRouteHint string         `json:"owner_route_hint"`
}

// Bot struct
type BotRes struct {
	UUID        string         `json:"uuid"`
	OwnerPubKey string         `json:"owner_pubkey"`
	Name        string         `json:"name"`
	UniqueName  string         `json:"unique_name"`
	Description string         `json:"description"`
	Tags        pq.StringArray `json:"tags"`
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

type Bounty struct {
	ID                      uint           `json:"id"`
	OwnerID                 string         `json:"owner_id"`
	Paid                    bool           `json:"paid"`
	Show                    bool           `gorm:"default:false" json:"show"`
	Completed               bool           `gorm:"default:false" json:"completed"`
	Type                    string         `json:"type"`
	Award                   string         `json:"award"`
	AssignedHours           uint8          `json:"assigned_hours"`
	BountyExpires           string         `json:"bounty_expires"`
	CommitmentFee           uint64         `json:"commitment_fee"`
	Price                   uint           `json:"price"`
	Title                   string         `json:"title"`
	Tribe                   string         `json:"tribe"`
	Assignee                string         `json:"assignee"`
	TicketUrl               string         `json:"ticket_url"`
	OrgUuid                 string         `json:"org_uuid"`
	Description             string         `json:"description"`
	WantedType              string         `json:"wanted_type"`
	Deliverables            string         `json:"deliverables"`
	GithubDescription       bool           `json:"github_description"`
	OneSentenceSummary      string         `json:"one_sentence_summary"`
	EstimatedSessionLength  string         `json:"estimated_session_length"`
	EstimatedCompletionDate string         `json:"estimated_completion_date"`
	Created                 int64          `json:"created"`
	Updated                 *time.Time     `json:"updated"`
	AssignedDate            *time.Time     `json:"assigned_date,omitempty"`
	CompletionDate          *time.Time     `json:"completion_date,omitempty"`
	MarkAsPaidDate          *time.Time     `json:"mark_as_paid_date,omitempty"`
	PaidDate                *time.Time     `json:"paid_date,omitempty"`
	CodingLanguages         pq.StringArray `gorm:"type:text[];not null default:'[]'" json:"coding_languages"`
}

type BountyOwners struct {
	OwnerID string `json:"owner_id"`
}

type BountyData struct {
	Bounty
	BountyId          uint       `json:"bounty_id"`
	BountyCreated     int64      `json:"bounty_created"`
	BountyUpdated     *time.Time `json:"bounty_updated"`
	BountyDescription string     `json:"bounty_description"`
	Person
	AssigneeAlias         string         `json:"assignee_alias"`
	AssigneeId            uint           `json:"assignee_id"`
	AssigneeImg           string         `json:"assignee_img"`
	AssigneeCreated       *time.Time     `json:"assignee_created"`
	AssigneeUpdated       *time.Time     `json:"assignee_updated"`
	AssigneeDescription   string         `json:"assignee_description"`
	AssigneeRouteHint     string         `json:"assignee_route_hint"`
	BountyOwnerId         uint           `json:"bounty_owner_id"`
	OwnerUuid             string         `json:"owner_uuid"`
	OwnerKey              string         `json:"owner_key"`
	OwnerAlias            string         `json:"owner_alias"`
	OwnerUniqueName       string         `json:"owner_unique_name"`
	OwnerDescription      string         `json:"owner_description"`
	OwnerTags             pq.StringArray `gorm:"type:text[]" json:"owner_tags" null`
	OwnerImg              string         `json:"owner_img"`
	OwnerCreated          *time.Time     `json:"owner_created"`
	OwnerUpdated          *time.Time     `json:"owner_updated"`
	OwnerLastLogin        int64          `json:"owner_last_login"`
	OwnerRouteHint        string         `json:"owner_route_hint"`
	OwnerContactKey       string         `json:"owner_contact_key"`
	OwnerPriceToMeet      int64          `json:"owner_price_to_meet"`
	OwnerTwitterConfirmed bool           `json:"owner_twitter_confirmed"`
	OrganizationName      string         `json:"organization_name"`
	OrganizationImg       string         `json:"organization_img"`
	WorkspaceUuid         string         `json:"organization_uuid"`
	WorkspaceDescription  string         `json:"description"`
}

type BountyResponse struct {
	Bounty       Bounty            `json:"bounty"`
	Assignee     Person            `json:"assignee"`
	Owner        Person            `json:"owner"`
	Organization OrganizationShort `json:"organization"`
}

type BountyCountResponse struct {
	OpenCount     int64 `json:"open_count"`
	AssignedCount int64 `json:"assigned_count"`
	PaidCount     int64 `json:"paid_count"`
}

type Organization struct {
	ID          uint       `json:"id"`
	Uuid        string     `json:"uuid"`
	Name        string     `gorm:"unique;not null" json:"name"`
	OwnerPubKey string     `json:"owner_pubkey"`
	Img         string     `json:"img"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
	Show        bool       `json:"show"`
	Deleted     bool       `gorm:"default:false" json:"deleted"`
	BountyCount int64      `json:"bounty_count,omitempty"`
	Budget      uint       `json:"budget,omitempty"`
	Website     string     `json:"website" validate:"omitempty,uri"`
	Github      string     `json:"github" validate:"omitempty,uri"`
	Description string     `json:"description" validate:"omitempty,lte=120"`
}

type OrganizationShort struct {
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

type OrganizationUsersData struct {
	OrgUuid     string     `json:"org_uuid"`
	UserCreated *time.Time `json:"user_created"`
	Person
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

type BountyBudget struct {
	ID          uint       `json:"id"`
	OrgUuid     string     `json:"org_uuid"`
	TotalBudget uint       `json:"total_budget"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
}

type StatusBudget struct {
	OrgUuid         string `json:"org_uuid"`
	CurrentBudget   uint   `json:"current_budget"`
	OpenBudget      uint   `json:"open_budget"`
	OpenCount       int64  `json:"open_count"`
	AssignedBudget  uint   `json:"assigned_budget"`
	AssignedCount   int64  `json:"assigned_count"`
	CompletedBudget uint   `json:"completed_budget"`
	CompletedCount  int64  `json:"completed_count"`
}

type BudgetInvoiceRequest struct {
	Amount          uint        `json:"amount"`
	SenderPubKey    string      `json:"sender_pubkey"`
	OrgUuid         string      `json:"org_uuid"`
	PaymentType     PaymentType `json:"payment_type,omitempty"`
	Websocket_token string      `json:"websocket_token,omitempty"`
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
	Created        *time.Time  `json:"created"`
	Updated        *time.Time  `json:"updated"`
	Status         bool        `json:"status"`
}

type PaymentHistoryData struct {
	PaymentHistory
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
	BountiesPaidPercentage uint  `json:"bounties_paid_average"`
	SatsPosted             uint  `json:"sats_posted"`
	SatsPaid               uint  `json:"sats_paid"`
	SatsPaidPercentage     uint  `json:"sats_paid_percentage"`
	AveragePaid            uint  `json:"average_paid"`
	AverageCompleted       uint  `json:"average_completed"`
	UniqueHuntersPaid      int64 `json:"unique_hunters_paid"`
	NewHuntersPaid         int64 `json:"new_hunters_paid"`
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

type FilterStattuCount struct {
	Open     int64 `json:"open"`
	Assigned int64 `json:"assigned"`
	Paid     int64 `json:"paid"`
}

func (Person) TableName() string {
	return "people"
}

func (PersonInShort) TableName() string {
	return "people"
}

func (Bounty) TableName() string {
	return "bounty"
}

func (ConnectionCodes) TableName() string {
	return "connectioncodes"
}

func (ConnectionCodesShort) TableName() string {
	return "connectioncodes"
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
