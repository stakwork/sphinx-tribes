package main

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	lnurl "github.com/fiatjaf/go-lnurl"
)

func encodeLNURL() (string, error) {
	host := os.Getenv("LN_SERVER_BASE_URL")

	k1 := generate32Bytes()
	println("Host ==", host)
	println("K1 ===", k1)
	url := host + "lnurl?tag=login&k1=" + k1 + "action=login"

	return lnurl.Encode(url)
}

func generate32Bytes() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		return ""
	}

	data := hex.EncodeToString(key)

	return data
}
