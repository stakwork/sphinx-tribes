package auth

import (
	"log"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/stakwork/sphinx-tribes/config"
)

var TokenAuth *jwtauth.JWTAuth

// Init auth
func InitJwt() {
	if config.JwtKey == "" {
		log.Fatal("No JWT key")
	}
	TokenAuth = jwtauth.New("HS256", []byte(config.JwtKey), nil)
}

// ExpireInHours for jwt
func ExpireInHours(hours int) int64 {
	return jwtauth.ExpireIn(time.Duration(hours) * time.Hour)
}
