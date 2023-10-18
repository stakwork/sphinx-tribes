package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strconv"
)

func GetRandomToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Random token erorr ==", err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func ConvertStringToUint(number string) (uint, error) {
	numberParse, err := strconv.ParseUint(number, 10, 32)

	if err != nil {
		fmt.Println("could not parse string to uint")
		return 0, err
	}

	return uint(numberParse), nil
}

func ConvertStringToInt(number string) (int, error) {
	numberParse, err := strconv.ParseInt(number, 10, 32)

	if err != nil {
		fmt.Println("could not parse string to uint")
		return 0, err
	}

	return int(numberParse), nil
}
