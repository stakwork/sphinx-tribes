package main

import (
	"crypto/rand"
	"os"

	lnurl "github.com/fiatjaf/go-lnurl"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
)

func encodeLNURL() (string, error) {
	host := os.Getenv("LN_SERVER_BASE_URL")

	k1 := generate32Bytes()
	url := host + "lnurl_login?tag=login&k1=" + k1 + "&action=login"

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
