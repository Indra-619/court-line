package config

import (
	"strings"
	"testing"
)

func TestValidateSecretRejectsEmpty(t *testing.T) {
	if err := validateSecret(""); err == nil {
		t.Fatal("expected error for empty secret, got nil")
	}
}

func TestValidateSecretRejectsForbiddenDefaults(t *testing.T) {
	forbidden := []string{
		"your-super-secret-jwt-key-change-this",
		"default-secret-change-in-production",
		"your-super-secret-jwt-key",
	}
	for _, s := range forbidden {
		if err := validateSecret(s); err == nil {
			t.Errorf("expected error for forbidden default %q, got nil", s)
		}
	}
}

func TestValidateSecretRejectsShortSecret(t *testing.T) {
	short := strings.Repeat("a", 31)
	err := validateSecret(short)
	if err == nil {
		t.Fatal("expected error for 31-character secret, got nil")
	}
	if err.Error() != "JWT_SECRET must be at least 32 characters" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateSecretAcceptsMinimumLength(t *testing.T) {
	if err := validateSecret(strings.Repeat("b", 32)); err != nil {
		t.Fatalf("expected 32-character secret to be accepted, got: %v", err)
	}
}

func TestValidateSecretAcceptsStrongValue(t *testing.T) {
	if err := validateSecret("a-genuinely-random-64-char-secret-value-used-only-in-tests-aaaa"); err != nil {
		t.Fatalf("unexpected error for strong secret: %v", err)
	}
}

func TestMustLoadCookieSecureFromEnv(t *testing.T) {
	const strong = "a-genuinely-random-64-char-secret-value-used-only-in-tests-aaaa"

	t.Setenv("JWT_SECRET", strong)

	// Default: not secure (local dev over plain HTTP).
	t.Setenv("COOKIE_SECURE", "")
	MustLoad()
	if CookieSecure() {
		t.Error("expected CookieSecure false when COOKIE_SECURE is unset")
	}

	// "true" and "1" opt in.
	t.Setenv("COOKIE_SECURE", "true")
	MustLoad()
	if !CookieSecure() {
		t.Error("expected CookieSecure true for COOKIE_SECURE=true")
	}

	t.Setenv("COOKIE_SECURE", "1")
	MustLoad()
	if !CookieSecure() {
		t.Error("expected CookieSecure true for COOKIE_SECURE=1")
	}

	// Anything else stays off.
	t.Setenv("COOKIE_SECURE", "yes-please")
	MustLoad()
	if CookieSecure() {
		t.Error("expected CookieSecure false for unrecognized COOKIE_SECURE value")
	}
}
