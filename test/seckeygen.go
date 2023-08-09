package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateRandomKey(keyLength int) (string, error) {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

func main() {
	// Replace 32 with the desired key length in bytes
	secureKey, err := generateRandomKey(32)
	if err != nil {
		fmt.Println("Error generating secure key:", err)
		return
	}

	fmt.Println("Secure Key:", secureKey)
}
