package config

import (
	"errors"
	"log"
	"os"
)

// forbiddenSecrets are known placeholder values that must never be used
// in a running system.
var forbiddenSecrets = []string{
	"your-super-secret-jwt-key-change-this",
	"default-secret-change-in-production",
	"your-super-secret-jwt-key",
}

var jwtSecret string

// validateSecret reports whether a candidate JWT secret is usable. It
// rejects empty values and known insecure defaults.
func validateSecret(s string) error {
	if s == "" {
		return errors.New("JWT_SECRET environment variable is required but not set")
	}
	for _, forbidden := range forbiddenSecrets {
		if s == forbidden {
			return errors.New("JWT_SECRET is set to a known insecure default; generate a strong random secret")
		}
	}
	return nil
}

// MustLoad reads required configuration from the environment and exits
// the process when a value is missing or left at an insecure default.
func MustLoad() {
	secret := os.Getenv("JWT_SECRET")
	if err := validateSecret(secret); err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	jwtSecret = secret
}

// JWTSecret returns the loaded JWT signing secret.
func JWTSecret() string {
	return jwtSecret
}
