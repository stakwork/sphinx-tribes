package main

import (
	"crypto/rand"
	"os"

	lnurl "github.com/fiatjaf/go-lnurl"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
)

type LnEncodeData struct {
	encode string
	k1     string
}

func encodeLNURL() (LnEncodeData, error) {
	var host string
	host = os.Getenv("LN_SERVER_BASE_URL")

	if host == "" {
		host = "people.sphinx.chat"
	}

	k1 := generate32Bytes()
	url := host + "lnurl_login?tag=login&k1=" + k1 + "&action=login"

	encode, err := lnurl.Encode(url)

	if err != nil {
		return LnEncodeData{}, err
	}

	return LnEncodeData{encode: encode, k1: k1}, nil
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
