package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

// Init auth
func initJwt() {
	jwtKey := os.Getenv("LN_JWT_KEY")
	fmt.Println("key ", jwtKey)
	if jwtKey == "" {
		log.Fatal("No JWT key")
	}
	TokenAuth = jwtauth.New("HS256", []byte(jwtKey), nil)
}

// ExpireInHours for jwt
func ExpireInHours(hours int) int64 {
	return jwtauth.ExpireIn(time.Duration(hours) * time.Hour)
}
