package main

import (
	"log"
	"time"

	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

// Init auth
func initJwt() {
	if jwtKey == "" {
		log.Fatal("No JWT key")
	}
	TokenAuth = jwtauth.New("HS256", []byte(jwtKey), nil)
}

// ExpireInHours for jwt
func ExpireInHours(hours int) int64 {
	return jwtauth.ExpireIn(time.Duration(hours) * time.Hour)
}
