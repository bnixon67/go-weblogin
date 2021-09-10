package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateRandomString returns n bytes encoded in URL friendly base64.
func GenerateRandomString(n int) (string, error) {
	// buffer to store n bytes
	b := make([]byte, n)

	// get b random bytes
	_, err := rand.Read(b)
	if err != nil {
		log.Panic(err)
		return "", err
	}

	// convert to URL friendly base64
	return base64.URLEncoding.EncodeToString(b), err
}
