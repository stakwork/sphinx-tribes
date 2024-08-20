package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Host string
var JwtKey string
var RelayUrl string
var MemeUrl string
var RelayAuthKey string
var RelayNodeKey string
var SuperAdmins []string = []string{""}

// these are constants for the store
var InvoiceList = "INVOICELIST"
var BudgetInvoiceList = "BUDGETINVOICELIST"
var S3BucketName string
var S3FolderName string
var S3Url string
var AdminCheck string
var AdminDevFreePass = "FREE_PASS"
var Connection_Auth string
var AdminStrings string

var S3Client *s3.Client
var PresignClient *s3.PresignClient

var V2BotUrl string
var V2BotToken string
var V2ContactKey string

func InitConfig() {
	Host = os.Getenv("LN_SERVER_BASE_URL")
	JwtKey = os.Getenv("LN_JWT_KEY")
	RelayUrl = os.Getenv("RELAY_URL")
	MemeUrl = os.Getenv("MEME_URL")
	RelayAuthKey = os.Getenv("RELAY_AUTH_KEY")
	AdminStrings = os.Getenv("ADMINS")
	AwsSecret := os.Getenv("AWS_SECRET_ACCESS")
	AwsAccess := os.Getenv("AWS_ACCESS_KEY_ID")
	AwsRegion := os.Getenv("AWS_REGION")
	S3BucketName = os.Getenv("S3_BUCKET_NAME")
	S3FolderName = os.Getenv("S3_FOLDER_NAME")
	S3Url = os.Getenv("S3_URL")
	AdminCheck = os.Getenv("ADMIN_CHECK")
	Connection_Auth = os.Getenv("CONNECTION_AUTH")
	V2BotUrl = os.Getenv("V2_BOT_URL")
	V2BotToken = os.Getenv("V2_BOT_TOKEN")

	// Add to super admins
	SuperAdmins = StripSuperAdmins(AdminStrings)

	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AwsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AwsAccess, AwsSecret, "")),
	)

	if err != nil {
		fmt.Println("Could not setup AWS session", err)
	}

	// create a s3 client session
	S3Client = s3.NewFromConfig(awsConfig)
	PresignClient = s3.NewPresignClient(S3Client)

	// only make this call if there is a Relay auth key
	if RelayAuthKey != "" {
		RelayNodeKey = GetNodePubKey()
	} else {
		panic("No relay auth key set")
	}

	if V2BotUrl != "" && V2BotToken != "" {
		contact_key := GetV2ContactKey()
		V2ContactKey = contact_key
	}

	if Host == "" {
		Host = "https://people.sphinx.chat"
	}

	if MemeUrl == "" {
		MemeUrl = "https://memes.sphinx.chat"
	}

	if JwtKey == "" {
		JwtKey = GenerateRandomString()
	}

	if S3BucketName == "" {
		S3BucketName = "sphinx-tribes"
	}

	if S3FolderName == "" {
		S3FolderName = "metrics"
	}

	if S3Url == "" {
		S3Url = "https://sphinx-tribes.s3.amazonaws.com"
	}
}

func StripSuperAdmins(adminStrings string) []string {
	superAdmins := []string{}
	if adminStrings != "" {
		if strings.Contains(adminStrings, ",") {
			splitArray := strings.Split(adminStrings, ",")
			splitLength := len(splitArray)

			for i := 0; i < splitLength; i++ {
				// append indexes, and skip all the commas
				if splitArray[i] == "," {
					continue
				} else {
					superAdmins = append(superAdmins, strings.TrimSpace(splitArray[i]))
				}
			}
		} else {
			superAdmins = append(superAdmins, strings.TrimSpace(adminStrings))
		}
	}
	return superAdmins
}

func GenerateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 24)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

type PropertyMap map[string]interface{}

type Feature map[string]PropertyMap

type NodeGetInfoResponse struct {
	Uris                   []string      `json:"uris"`
	Chains                 []PropertyMap `json:"chains"`
	Features               Feature       `json:"features"`
	IdentityPubkey         string        `json:"identity_pubkey"`
	Alias                  string        `json:"alias"`
	NumPendingChannels     uint          `json:"num_pending_channels"`
	NumActiveChannels      uint          `json:"num_active_channels"`
	NumInactiveChannels    uint          `json:"num_inactive_channels"`
	NumPeers               uint          `json:"num_peers"`
	BlockHeight            uint          `json:"block_height"`
	BlockHash              string        `json:"block_hash"`
	SyncedToChain          bool          `json:"synced_to_chain"`
	Testnet                bool          `json:"testnet"`
	BestHeaderTimestamp    string        `json:"best_header_timestamp"`
	Version                string        `json:"version"`
	Color                  string        `json:"color"`
	SyncedToGraph          bool          `json:"synced_to_graph"`
	CommitHash             string        `json:"commit_hash"`
	RequireHtlcInterceptor bool          `json:"require_htlc_interceptor"`
}

type NodeGetInfo struct {
	Success  bool                `json:"success"`
	Response NodeGetInfoResponse `json:"response"`
}

type V2AccountInfo struct {
	ContactInfo string `json:"contact_info"`
	Alias       string `json:"alias"`
	Img         string `json:"img"`
	Network     string `json:"network"`
}

type Contact struct {
	Id                uint   `json:"id"`
	RouteHint         string `json:"route_hint"`
	PublicKey         string `json:"public_key"`
	NodeAlias         string `json:"node_alias,omitempty"`
	Alias             string `json:"alias"`
	PhotoUrl          string `json:"photo_url,omitempty"`
	PrivatePhoto      int    `json:"private_photo"`
	IsOwner           uint   `json:"is_owner"`
	Deleted           int    `json:"deleted"`
	RemoteId          string `json:"remote_id,omitempty"`
	Status            int    `json:"status,omitempty"`
	ContactKey        string `json:"contact_key"`
	DeviceId          string `json:"device_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	FromGroup         int    `json:"from_group"`
	NotificationSound string `json:"notification_sound,omitempty"`
	LastActive        string `json:"last_active"`
	TipAmount         string `json:"tip_amount,omitempty"`
	Tenant            uint   `json:"tenant"`
	PriceToMeet       string `json:"price_to_meet,omitempty"`
	Unmet             string `json:"unmet"`
	Blocked           string `json:"blocked"`
	HmacKey           string `json:"hmac_key"`
	PersonUuid        string `json:"person_uuid,omitempty"`
	LastTimestamp     uint   `json:"last_timestamp"`
	IsAdmin           int    `json:"is_admin"`
	PushKitToken      string `json:"push_kit_token,omitempty"`
	Prune             string `json:"prune,omitempty"`
	AdminToken        string `json:"admin_token,omitempty"`
}

type Chat struct {
	Id                 uint          `json:"id"`
	Uuid               string        `json:"uuid"`
	Name               string        `json:"name"`
	PhotoUrl           string        `json:"photo_url"`
	Type               int           `json:"type"`
	Status             int           `json:"status"`
	ContactIds         []int         `json:"contact_id"`
	IsMuted            int           `json:"is_muted"`
	CreatedAt          string        `json:"created_at"`
	UpdatedAt          string        `json:"updated_at"`
	Deleted            int           `json:"deleted"`
	GroupKey           string        `json:"group_key"`
	Host               string        `json:"host"`
	PriceToJoin        int           `json:"price_to_join"`
	PricePerMessage    string        `json:"price_per_message,omitempty"`
	EscrowAmount       int           `json:"escrow_amount,omitempty"`
	EscrowMillis       int           `json:"escrow_millis,omitempty"`
	Unlisted           int           `json:"unlisted"`
	Private            int           `json:"private"`
	OwnerPubkey        string        `json:"owner_pubkey"`
	Seen               int           `json:"seen"`
	AppUrl             string        `json:"app_url,omitempty"`
	FeedUrl            string        `json:"feed_url,omitempty"`
	FeedType           string        `json:"feed_type,omitempty"`
	Meta               string        `json:"meta,omitempty"`
	MyPhotoUrl         string        `json:"my_photo_url,omitempty"`
	MyALias            string        `json:"my_alias,omitempty"`
	Tenant             uint          `json:"tenant"`
	SkipBroadcastJoins int           `json:"skip_broadcast_joins"`
	Pin                string        `json:"pin,omitempty"`
	Notify             int           `json:"notify"`
	ProfileFilters     string        `json:"profile_filters,omitempty"`
	CallRecording      string        `json:"call_recording,omitempty"`
	MemeServerLocation string        `json:"meme_server_location,omitempty"`
	JitsiServer        string        `json:"jitsi_server,omitempty"`
	StakworkApiKey     string        `json:"stakwork_api_key"`
	StakworkWebhook    string        `json:"stakwork_webhook"`
	DefaultJoin        int           `json:"default_join"`
	Preview            string        `json:"preview,omitempty"`
	PendingContactIds  []interface{} `json:"pending_contact_ids,omitempty"`
}

type ContactResponse struct {
	Contacts      []Contact     `json:"contacts"`
	Chats         []Chat        `json:"chats"`
	Subscriptions []interface{} `json:"subscriptions,omitempty"`
}

type ProxyContacts struct {
	Success  bool            `json:"success"`
	Response ContactResponse `json:"response"`
}

func GetV2ContactKey() string {
	url := fmt.Sprintf("%s/account", V2BotUrl)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Printf("Get Contact Request Failed: %s", err)
	}

	req.Header.Set("x-admin-token", V2BotToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Printf("Get Contact Request Failed: %s", err)
		return ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Get Contact Body Read Failed: %s", err)
	}

	accountInfo := V2AccountInfo{}

	// Unmarshal result
	err = json.Unmarshal(body, &accountInfo)
	if err != nil {
		log.Printf("Reading Relay Node Info body failed: %s", err)
	}

	contact_key := accountInfo.ContactInfo
	return contact_key
}

func GetNodePubKey() string {
	var pubkey string
	var url string
	var isProxy bool = false

	if strings.Contains(RelayUrl, "swarm") {
		url = fmt.Sprintf("%s/contacts", RelayUrl)
		isProxy = true
	} else {
		url = fmt.Sprintf("%s/getinfo", RelayUrl)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	req.Header.Set("x-user-token", RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	if isProxy {
		proxyContacts := ProxyContacts{}
		err = json.Unmarshal(body, &proxyContacts)
		if err != nil {
			log.Printf("Reading Proxy Contacts Info body failed: %s", err)
		}
		contacts := proxyContacts.Response.Contacts
		if len(contacts) > 0 {
			pubkey = contacts[0].PublicKey
		}
	} else {
		nodeInfo := NodeGetInfo{}
		// Unmarshal result
		err = json.Unmarshal(body, &nodeInfo)
		if err != nil {
			log.Printf("Reading Relay Node Info body failed: %s", err)
		}
		pubkey = nodeInfo.Response.IdentityPubkey
	}
	return pubkey
}
