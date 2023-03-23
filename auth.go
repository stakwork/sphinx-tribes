package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/form3tech-oss/jwt-go"
)

var (
	// signedMsgPrefix is a special prefix that we'll prepend to any
	// messages we sign/verify. We do this to ensure that we don't
	// accidentally sign a sighash, or other sensitive material. By
	// prepending this fragment, we mind message signing to our particular
	// context.
	signedMsgPrefix = []byte("Lightning Signed Message:")
)

type contextKey string

// ContextKey ...
var ContextKey = contextKey("key")

// PubKeyContext parses pukey from signed timestamp
func PubKeyContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("x-jwt")
		}

		if token == "" {
			fmt.Println("[auth] no token")
			http.Error(w, http.StatusText(401), 401)
			return
		}

		isJwt := strings.Contains(token, ".")

		if isJwt {
			claims, err := DecodeToken(token)

			if err != nil {
				fmt.Println("Failed to parse JWT")
				http.Error(w, http.StatusText(401), 401)
				return
			}

			if claims.VerifyExpiresAt(time.Now().UnixNano(), true) {
				fmt.Println("Token has expired")
				http.Error(w, http.StatusText(401), 401)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKey, claims["pubkey"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			pubkey, err := VerifyTribeUUID(token, true)

			if pubkey == "" || err != nil {
				fmt.Println("[auth] no pubkey || err != nil")
				if err != nil {
					fmt.Println(err)
				}
				http.Error(w, http.StatusText(401), 401)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKey, pubkey)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// VerifyTribeUUID takes base64 uuid and returns hex pubkey
func VerifyTribeUUID(uuid string, checkTimestamp bool) (string, error) {
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

	if checkTimestamp {
		// 5 MINUTE MAX
		ts := int64(binary.BigEndian.Uint32(timeBuf))
		now := time.Now().Unix()
		if ts < now-300 {
			fmt.Println("TOO LATE!")
			return "", errors.New("too late")
		}
	}

	return pubkey, nil
}

// VerifyArbitrary takes base64 sig and msg and returns hex pubkey
func VerifyArbitrary(sig string, msg string) (string, error) {
	sigByes, err := base64.URLEncoding.DecodeString(sig)
	if err != nil {
		return "", err
	}
	pubkey, valid, err := VerifyAndExtract([]byte(msg), sigByes)
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
	pubKey, valid, err := btcecdsa.RecoverCompact(sig, digest)
	if err != nil {
		fmt.Printf("ERR: %+v\n", err)
		return "", false, err
	}
	pubKeyHex := hex.EncodeToString(pubKey.SerializeCompressed())

	return pubKeyHex, valid, nil
}

func DecodeToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		key := jwtKey
		return []byte(key), nil
	})

	return claims, err
}

func EncodeToken(pubkey string) (string, error) {
	exp := ExpireInHours(24 * 7)

	claims := jwt.MapClaims{
		"pubkey": pubkey,
		"exp":    exp,
	}

	_, tokenString, err := TokenAuth.Encode(claims)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
