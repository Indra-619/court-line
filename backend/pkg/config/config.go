package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"
)

// forbiddenSecrets are known placeholder values that must never be used
// in a running system.
var forbiddenSecrets = []string{
	"your-super-secret-jwt-key-change-this",
	"default-secret-change-in-production",
	"your-super-secret-jwt-key",
}

var jwtSecret string

// cookieSecure controls the Secure flag on session cookies. It is read
// once at load time from COOKIE_SECURE so deployments behind HTTPS can
// opt in without a code change; local dev keeps the default false.
var cookieSecure bool

// validateSecret reports whether a candidate JWT secret is usable. It
// rejects empty values and known insecure defaults.
func validateSecret(s string) error {
	if s == "" {
		return errors.New("JWT_SECRET environment variable is required but not set")
	}
	if len(s) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}
	for _, forbidden := range forbiddenSecrets {
		if s == forbidden {
			return errors.New("JWT_SECRET is set to a known insecure default; generate a strong random secret")
		}
	}
	return nil
}

// jwtTTL is the access-token lifetime. Short-lived access tokens limit
// the replay window of a leaked JWT; the rotating refresh token (30d)
// keeps sessions alive, and the frontend auto-refreshes on 401.
// Override with JWT_TTL_MINUTES (clamped 5..1440).
var jwtTTL = 30 * time.Minute

// loadJWTTTL reads JWT_TTL_MINUTES, keeping the default when unset or
// invalid and clamping to a sane range.
func loadJWTTTL() {
	v := os.Getenv("JWT_TTL_MINUTES")
	if v == "" {
		return
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		log.Printf("warning: invalid JWT_TTL_MINUTES=%q, keeping default %v", v, jwtTTL)
		return
	}
	if n < 5 {
		n = 5
	}
	if n > 1440 {
		n = 1440
	}
	jwtTTL = time.Duration(n) * time.Minute
}

// JWTTTL returns the configured access-token lifetime.
func JWTTTL() time.Duration {
	return jwtTTL
}

// MustLoad reads required configuration from the environment and exits
// the process when a value is missing or left at an insecure default.
func MustLoad() {
	secret := os.Getenv("JWT_SECRET")
	if err := validateSecret(secret); err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	jwtSecret = secret

	loadJWTTTL()

	switch v := os.Getenv("COOKIE_SECURE"); v {
	case "true", "1":
		cookieSecure = true
	default:
		cookieSecure = false
	}
}

// CookieSecure reports whether session cookies must carry the Secure
// flag (set via COOKIE_SECURE=true in production).
func CookieSecure() bool {
	return cookieSecure
}

// JWTSecret returns the loaded JWT signing secret.
func JWTSecret() string {
	return jwtSecret
}
