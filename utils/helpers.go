package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
)

func GetRandomToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Random token erorr ==", err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
