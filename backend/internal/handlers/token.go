package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// generateJTI returns a cryptographically random 16-byte hex string
// used as the unique JWT identifier (jti) claim. The jti lets a token
// be blacklisted individually on logout.
func generateJTI() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// generateRefreshToken returns a cryptographically random 32-byte hex
// string. The raw value is handed to the client exactly once; only its
// SHA-256 hash is persisted.
func generateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// sha256Hex returns the lowercase hex-encoded SHA-256 digest of the
// given token. Refresh tokens are stored and looked up only by hash so
// a database leak cannot be replayed as a session.
func sha256Hex(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
