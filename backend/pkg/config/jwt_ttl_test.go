package config

import (
	"testing"
	"time"
)

func TestJWTTTLDefault(t *testing.T) {
	jwtTTL = 30 * time.Minute
	t.Setenv("JWT_TTL_MINUTES", "")
	loadJWTTTL()
	if JWTTTL() != 30*time.Minute {
		t.Errorf("default TTL: expected 30m, got %v", JWTTTL())
	}
}

func TestJWTTTLOverride(t *testing.T) {
	jwtTTL = 30 * time.Minute
	t.Setenv("JWT_TTL_MINUTES", "15")
	loadJWTTTL()
	if JWTTTL() != 15*time.Minute {
		t.Errorf("expected 15m, got %v", JWTTTL())
	}
}

func TestJWTTTLClamps(t *testing.T) {
	cases := []struct {
		env  string
		want time.Duration
	}{
		{"1", 5 * time.Minute},        // below floor
		{"99999", 1440 * time.Minute}, // above ceiling
	}
	for _, tc := range cases {
		jwtTTL = 30 * time.Minute
		t.Setenv("JWT_TTL_MINUTES", tc.env)
		loadJWTTTL()
		if JWTTTL() != tc.want {
			t.Errorf("env %q: expected %v, got %v", tc.env, tc.want, JWTTTL())
		}
	}
}

func TestJWTTTLInvalidKeepsDefault(t *testing.T) {
	for _, v := range []string{"abc", "-5", "0"} {
		jwtTTL = 30 * time.Minute
		t.Setenv("JWT_TTL_MINUTES", v)
		loadJWTTTL()
		if JWTTTL() != 30*time.Minute {
			t.Errorf("env %q: expected default kept, got %v", v, JWTTTL())
		}
	}
}
