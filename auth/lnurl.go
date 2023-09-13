package auth

import (
	"crypto/rand"
	"fmt"
	"strings"

	lnurl "github.com/fiatjaf/go-lnurl"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
	"github.com/stakwork/sphinx-tribes/config"
)

type LnEncodeData struct {
	Encode string
	K1     string
}

func EncodeLNURL(host string) (LnEncodeData, error) {
	fmt.Println("Host ===", host)
	hostUrl := config.Host
	if !strings.Contains(host, "localhost") {
		hostUrl = "https://" + host
	}
	k1 := generate32Bytes()
	url := hostUrl + "/" + "lnauth_login?tag=login&k1=" + k1 + "&action=login"

	encode, err := lnurl.Encode(url)

	if err != nil {
		return LnEncodeData{}, err
	}

	return LnEncodeData{Encode: encode, K1: k1}, nil
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
