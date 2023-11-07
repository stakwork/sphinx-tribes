package config

import (
	"math/rand"
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
	RelayNodeKey = os.Getenv("RELAY_NODE_KEY")

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
