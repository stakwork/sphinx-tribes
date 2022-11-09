package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
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
	Tags            pq.StringArray `json:"tags"`
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
}

// Bot struct
type Bot struct {
	UUID           string         `json:"uuid"`
	OwnerPubKey    string         `json:"owner_pubkey"`
	OwnerAlias     string         `json:"owner_alias"`
	Name           string         `json:"name"`
	UniqueName     string         `json:"unique_name"`
	Description    string         `json:"description"`
	Tags           pq.StringArray `json:"tags"`
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

func (Person) TableName() string {
	return "people"
}

// Person struct
type Person struct {
	ID               uint           `json:"id"`
	Uuid             string         `json:"uuid"`
	OwnerPubKey      string         `json:"owner_pubkey"`
	OwnerAlias       string         `json:"owner_alias"`
	UniqueName       string         `json:"unique_name"`
	Description      string         `json:"description"`
	Tags             pq.StringArray `json:"tags"`
	Img              string         `json:"img"`
	Created          *time.Time     `json:"created"`
	Updated          *time.Time     `json:"updated"`
	Unlisted         bool           `json:"unlisted"`
	Deleted          bool           `json:"deleted"`
	LastLogin        int64          `json:"last_login"`
	OwnerRouteHint   string         `json:"owner_route_hint"`
	OwnerContactKey  string         `json:"owner_contact_key"`
	PriceToMeet      int64          `json:"price_to_meet"`
	Extras           PropertyMap    `json:"extras", type: jsonb not null default '{}'::jsonb`
	TwitterConfirmed bool           `json:"twitter_confirmed"`
	GithubIssues     PropertyMap    `json:"github_issues", type: jsonb not null default '{}'::jsonb`
	NewTicketTime    int64          `json:"new_ticket_time", gorm: "-:all"`
}

// Github struct
type GithubIssue struct {
	// ID          uint `json:"id"`
	// PersonID    uint `json:"person_id"`
	// Person      Person
	// URL         string `json:"url"` // this will function as id
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

type AssetResponse struct {
	Balances []AssetBalanceData `json:"balances"`
	Txs      []string           `json:"txs"`
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

/*
GithubIssues
stakwork/sphinx-relay/229: {
	status: open
	assignee: Evanfeenstra
}
*/

/* loopbot:
{
	prefix:'/loopout',
	price:0,
	commands: [{
		command: '*',
		price: 0,
		min_price: 250000,
		max_price: 16777215,
		price_index: 2,
		admin_only: false
	}]
}
*/

/* btc bot:
{
	prefix:'/btc',
	price:10,
	commands: null
}
*/

// PropertyMap ...
type PropertyMap map[string]interface{}

// Value ...
func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
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
