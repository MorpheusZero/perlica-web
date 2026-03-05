package util

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	time       = 4
	memory     = 64 * 1024
	threads    = 4
	keyLength  = 32
)

// HashPassword hashes a password using argon2id and returns the hashed password as a string
func HashPassword(password string) string {
	salt := make([]byte, saltLength)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLength)

	return base64.StdEncoding.EncodeToString(salt) + "$" + base64.StdEncoding.EncodeToString(hash)
}

// VerifyPassword verifies a plaintext password against a hashed password
func VerifyPassword(plainTextPassword, hashedPassword string) bool {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(plainTextPassword), salt, time, memory, threads, keyLength)
	hashedInput := base64.StdEncoding.EncodeToString(hash)

	return hashedInput == parts[1]
}

// GenerateRandomPassword generates a random 10-character password with letters, numbers, and symbols
func GenerateRandomPassword() (string, error) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	password := make([]byte, 10)

	for i := range password {
		randomIndex := make([]byte, 1)
		_, err := rand.Read(randomIndex)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex[0]%byte(len(charset))]
	}

	return string(password), nil
}
