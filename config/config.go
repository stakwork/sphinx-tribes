package config

import (
	"math/rand"
	"os"
	"time"
)

var Host string
var JwtKey string

func InitConfig() {
	Host = os.Getenv("LN_SERVER_BASE_URL")
	JwtKey = os.Getenv("LN_JWT_KEY")

	if Host == "" {
		Host = "https://people.sphinx.chat"
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
