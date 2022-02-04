package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

func main() {
	key, hash, _ := GenerateApiKey("damga-proxy")
	fmt.Println(key)
	fmt.Println(hash)
}

func GenerateApiKey(name string) (string, string, error) {
	apiKey, err := GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	fmt.Println("-----------")
	fmt.Println("ApiKey:   ", apiKey)
	fmt.Println("HashedKey:", string(hashedKey))
	fmt.Println("Base64:", name)
	fmt.Println("Base64:", base64.StdEncoding.EncodeToString([]byte(name+":"+apiKey)))
	fmt.Println("-----------")
	return apiKey, string(hashedKey), nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
