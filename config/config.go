package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var Host string
var JwtKey string
var RelayUrl string
var MemeUrl string
var RelayAuthKey string
var RelayNodeKey string

// these are constants for the store
var InvoiceList = "INVOICELIST"
var BudgetInvoiceList = "BUDGETINVOICELIST"

func InitConfig() {
	Host = os.Getenv("LN_SERVER_BASE_URL")
	JwtKey = os.Getenv("LN_JWT_KEY")
	RelayUrl = os.Getenv("RELAY_URL")
	MemeUrl = os.Getenv("MEME_URL")
	RelayAuthKey = os.Getenv("RELAY_AUTH_KEY")

	// only make this call if there is a Relay auth key
	if RelayAuthKey != "" {
		RelayNodeKey = GetNodePubKey()
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

func GetNodePubKey() string {
	var pubkey string
	url := fmt.Sprintf("%s/getinfo", RelayUrl)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("x-user-token", RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	nodeInfo := NodeGetInfo{}

	// Unmarshal result
	err = json.Unmarshal(body, &nodeInfo)

	if err != nil {
		log.Printf("Reading Relay Node Info body failed: %s", err)
	}

	pubkey = nodeInfo.Response.IdentityPubkey

	return pubkey
}
