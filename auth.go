package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type contextKey string

// ContextKey ...
var ContextKey = contextKey("key")

// PubKeyContext parses pukey from signed timestamp
func PubKeyContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		pubkey, err := VerifyTribeUUID(token)
		if pubkey == "" || err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey, pubkey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifyTribeUUID takes base64 uuid and returns hex pubkey
func VerifyTribeUUID(uuid string) (string, error) {
	sigByes, err := base64.URLEncoding.DecodeString(uuid)
	if err != nil {
		return "", err
	}

	timeBuf := sigByes[:4] // unix timestamp is 4 bytes, or uint32
	sigBuf := sigByes[4:]
	pubkey, valid, err := VerifyAndExtract(timeBuf, sigBuf)
	if err != nil || !valid || pubkey == "" {
		return "", err
	}
	return pubkey, nil
}

// VerifyAndExtract ... pubkey comes out hex encoded
func VerifyAndExtract(msg, sig []byte) (string, bool, error) {

	if sig == nil || msg == nil {
		return "", false, errors.New("bad")
	}
	msg = append(signedMsgPrefix, msg...)
	digest := chainhash.DoubleHashB(msg)

	// RecoverCompact both recovers the pubkey and validates the signature.
	pubKey, valid, err := btcec.RecoverCompact(btcec.S256(), sig, digest)
	if err != nil {
		fmt.Printf("ERR: %+v\n", err)
		return "", false, err
	}
	pubKeyHex := hex.EncodeToString(pubKey.SerializeCompressed())

	return pubKeyHex, valid, nil
}
