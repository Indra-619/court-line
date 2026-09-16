package config

import (
	"log"
	"os"
)

// forbiddenSecrets are known placeholder values that must never be used
// in a running system.
var forbiddenSecrets = []string{
	"your-super-secret-jwt-key-change-this",
	"default-secret-change-in-production",
}

var jwtSecret string

// MustLoad reads required configuration from the environment and exits
// the process when a value is missing or left at an insecure default.
func MustLoad() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalf("configuration error: JWT_SECRET environment variable is required but not set")
	}
	for _, forbidden := range forbiddenSecrets {
		if secret == forbidden {
			log.Fatalf("configuration error: JWT_SECRET is set to a known insecure default; generate a strong random secret")
		}
	}
	jwtSecret = secret
}

// JWTSecret returns the loaded JWT signing secret.
func JWTSecret() string {
	return jwtSecret
}
