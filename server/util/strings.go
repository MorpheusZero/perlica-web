package util

import (
	"encoding/base64"
	"errors"
	"strings"
)

// StartsWith checks if the string s starts with the specified prefix.
func StartsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// DecodeBasicAuth decodes a base64-encoded "username:password" string and returns the username and password.
func DecodeBasicAuth(encoded string) (username, password string, err error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	decoded := string(decodedBytes)
	parts := strings.SplitN(decoded, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid basic auth format")
	}

	return parts[0], parts[1], nil
}
