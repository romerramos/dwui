package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var (
	passwordHash string
)

func ValidatePassword(password string) bool {
	hash := sha256.Sum256([]byte(password))
	providedHash := hex.EncodeToString(hash[:])
	return providedHash == passwordHash
}

func SetPassword(password string) {
	hash := sha256.Sum256([]byte(password))
	passwordHash = hex.EncodeToString(hash[:])
}

func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	password := fmt.Sprintf("%x", bytes)[:12]
	return password, nil
}
