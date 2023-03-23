package main

import (
	"math/rand"
	"os"
	"time"
)

var host string
var jwtKey string

func initConfig() {
	host = os.Getenv("LN_SERVER_BASE_URL")
	jwtKey = os.Getenv("LN_JWT_KEY")

	if host == "" {
		host = "people.sphinx.chat"
	}

	if jwtKey == "" {
		jwtKey = generateRandomString()
	}
}

func generateRandomString() string {
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
