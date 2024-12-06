package auth

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

	btcec "github.com/btcsuite/btcd/btcec/v2"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/form3tech-oss/jwt-go"
	"github.com/stakwork/sphinx-tribes/config"
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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		isJwt := strings.Contains(token, ".") && !strings.HasPrefix(token, ".")

		if isJwt {
			claims, err := DecodeJwt(token)

			if err != nil {
				fmt.Println("Failed to parse JWT")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if claims.VerifyExpiresAt(time.Now().UnixNano(), true) {
				fmt.Println("Token has expired")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKey, pubkey)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// PubKeyContext parses pukey from signed timestamp
func PubKeyContextSuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("x-jwt")
		}

		if token == "" {
			fmt.Println("[auth] no token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		isJwt := strings.Contains(token, ".") && !strings.HasPrefix(token, ".")
		if isJwt {
			claims, err := DecodeJwt(token)

			if err != nil {
				fmt.Println("Failed to parse JWT")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if claims.VerifyExpiresAt(time.Now().UnixNano(), true) {
				fmt.Println("Token has expired")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			pubkey := fmt.Sprintf("%v", claims["pubkey"])
			if !IsFreePass() && !AdminCheck(pubkey) {
				fmt.Println("Not a super admin")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if !IsFreePass() && !AdminCheck(pubkey) {
				fmt.Println("Not a super admin : auth")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKey, pubkey)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// ConnectionContext parses token for connection code
func ConnectionCodeContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r == nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		token := r.Header.Get("token")

		if token == "" {
			fmt.Println("[auth] no token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if token != config.Connection_Auth {
			fmt.Println("Not a super admin : auth")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CypressContext allows testing for cypress
func CypressContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsFreePass() {
			ctx := context.WithValue(r.Context(), ContextKey, "")
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println("Endpoint is for testing only : test endpoint")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	})
}

func AdminCheck(pubkey string) bool {
	for _, val := range config.SuperAdmins {
		if val == pubkey {
			return true
		}
	}
	return false
}

func IsFreePass() bool {
	if len(config.SuperAdmins) == 1 && config.SuperAdmins[0] == config.AdminDevFreePass || config.AdminStrings == "" {
		return true
	}
	return false
}

// VerifyTribeUUID takes base64 uuid and returns hex pubkey
func VerifyTribeUUID(uuid string, checkTimestamp bool) (string, error) {

	ts, timeBuf, sigBuf, err := ParseTokenString(uuid)
	if err != nil {
		return "", err
	}

	pubkey, valid, err := VerifyAndExtract(timeBuf, sigBuf)
	if err != nil || !valid || pubkey == "" {
		return "", err
	}

	if checkTimestamp {
		// 5 MINUTE MAX
		now := time.Now().Unix()
		if int64(ts) < now-300 {
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

// all 3 arguments are hex strings
func VerifyDerSig(sig string, hash string, pubkey string) (bool, error) {
	decoded, err := hex.DecodeString(sig)
	if err != nil {
		return false, err
	}
	signature, err := btcecdsa.ParseDERSignature(decoded)
	if err != nil {
		return false, err
	}
	msg, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	pubkeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return false, err
	}
	publicKey, err := btcec.ParsePubKey(pubkeyBytes)
	if err != nil {
		return false, err
	}
	return signature.Verify(msg, publicKey), nil
}

func DecodeJwt(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		key := config.JwtKey
		return []byte(key), nil
	})

	return claims, err
}

func EncodeJwt(pubkey string) (string, error) {

	if pubkey == "" || strings.ContainsAny(pubkey, "!@#$%^&*()") {
		return "", errors.New("invalid public key")
	}

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

// tribe UUID is a base64 encoded string 69 bytes long
// first 4 bytes is the timestamp
// last 65 bytes is the sign

// it can have two signature methods: signing the straight bytes
// OR base64 encoding then utf8-string encoding than signing again.
// the second method always prefixes the token with a "."
// that way, signers that only support utf8 (CLN) can still make tokens

func ParseTokenString(t string) (uint32, []byte, []byte, error) {
	token := t
	forceUtf8 := false
	// this signifies it's forced utf8 sig (for CLN SignMessage)
	if strings.HasPrefix(t, ".") {
		token = t[1:]
		forceUtf8 = true
	}
	tBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return 0, nil, nil, err
	}
	if len(tBytes) < 5 {
		return 0, nil, nil, errors.New("invalid signature (too short)")
	}
	sig := tBytes[4:]
	timeBuf := tBytes[:4]
	ts := binary.BigEndian.Uint32(timeBuf)
	if forceUtf8 {
		ts64 := base64.URLEncoding.EncodeToString(timeBuf)
		return ts, []byte(ts64), sig, nil
	} else {
		timeBuf := tBytes[:4]
		return ts, timeBuf, sig, nil
	}
}

func Sign(msg []byte, privKey *btcec.PrivateKey) ([]byte, error) {
	if msg == nil {
		//w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New("no msg")
	}

	msg = append(signedMsgPrefix, msg...)
	digest := chainhash.DoubleHashB(msg)
	// btcec.S256(), sig, digest

	sigBytes, err := btcecdsa.SignCompact(privKey, digest, true)
	if err != nil {
		return nil, err
	}

	return sigBytes, nil
}
