package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateActivationToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}