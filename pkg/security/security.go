package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	saltSize    = 16
	bcryptCost  = bcrypt.DefaultCost
)

// GenerateSalt creates a base64-encoded random salt.
func GenerateSalt() (string, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword combines password + salt and hashes it using bcrypt.
func HashPassword(password, salt string) (string, error) {
	combined := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(combined), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword compares the hash with the password + salt.
func VerifyPassword(password, salt, hash string) bool {
	combined := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(combined))
	return err == nil
}