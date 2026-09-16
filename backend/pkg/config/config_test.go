package config

import "testing"

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

func TestValidateSecretAcceptsStrongValue(t *testing.T) {
	if err := validateSecret("a-genuinely-random-64-char-secret-value-used-only-in-tests-aaaa"); err != nil {
		t.Fatalf("unexpected error for strong secret: %v", err)
	}
}
