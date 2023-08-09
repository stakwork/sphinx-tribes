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
	OwnerPubKey      string         `json:"owner_pubkey"`
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
}

type InvoiceStatus struct {
	Payment_request string `json:"payment_request"`
	Status          bool   `json:"Status"`
}

type InvoiceResult struct {
	Success  bool                 `json:"success"`
	Response InvoiceCheckResponse `json:"response"`
}

type InvoiceCheckResponse struct {
	Settled         bool   `json:"settled"`
	Payment_request string `json:"payment_request"`
	Payment_hash    string `json:"payment_hash"`
	Preimage        string `json:"preimage"`
	Amount          uint   `json:"amount"`
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
	Show                    bool           `json:"show"`
	Type                    string         `json:"type"`
	Award                   string         `json:"award"`
	AssignedHours           uint8          `json:"assigned_hours"`
	BountyExpires           string         `json:"bounty_expires"`
	CommitmentFee           uint64         `json:"commitment_fee"`
	Price                   string         `json:"price"`
	Title                   string         `json:"title"`
	Tribe                   string         `json:"tribe"`
	Created                 int64          `json:"created"`
	Assignee                string         `json:"assignee"`
	TicketUrl               string         `json:"ticket_url"`
	Description             string         `json:"description"`
	WantedType              string         `json:"wanted_type"`
	Deliverables            string         `json:"deliverables"`
	GithubDescription       bool           `json:"github_description"`
	OneSentenceSummary      string         `json:"one_sentence_summary"`
	EstimatedSessionLength  string         `json:"estimated_session_length"`
	EstimatedCompletionDate string         `json:"estimated_completion_date"`
	Updated                 *time.Time     `json:"updated"`
	CodingLanguages         pq.StringArray `gorm:"type:text[];not null default:'[]'" json:"coding_languages"`
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
	AssigneeCreated       *time.Time     `json:"assignee_created"`
	AssigneeUpdated       *time.Time     `json:"assignee_updated"`
	AssigneeDescription   string         `json:"assignee_description"`
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
}

type BountyResponse struct {
	Bounty   Bounty `json:"bounty"`
	Assignee Person `json:"assignee"`
	Owner    Person `json:"owner"`
}

type Organization struct {
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	OwnerPubKey string     `json:"owner_pubkey"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
	Show        bool       `json:"show"`
}

type OrganizationUsers struct {
	OwnerPubKey  string     `json:"owner_pubkey"`
	Organization string     `json:"organization"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
}

type BountyRoles struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserRoles struct {
	Role         string     `json:"role"`
	OwnerPubKey  string     `json:"owner_pubkey"`
	Organization string     `json:"organization"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
}

type BountyBudget struct {
	Organization string     `json:"organization"`
	TotalBudget  uint       `json:"total_budget"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
}

type BudgetHistory struct {
	Organization string     `json:"organization"`
	Amount       uint       `json:"amount"`
	SenderPubKey string     `json:"sender_pubkey"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
}

type PaymentHistory struct {
	Organization   string     `json:"organization"`
	SenderPubKey   string     `json:"sender_pubkey"`
	ReceiverPubKey string     `json:"receiver_pubkey"`
	Amount         uint       `json:"amount"`
	BountyId       uint       `json:"id"`
	Created        *time.Time `json:"created"`
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
