package handlers

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Indra-619/court-line/backend/pkg/config"
)

func TestGenerateJTILengthAndUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		jti, err := generateJTI()
		if err != nil {
			t.Fatalf("generateJTI failed: %v", err)
		}
		// 16 random bytes encoded as hex => 32 characters
		if len(jti) != 32 {
			t.Fatalf("expected jti length 32, got %d", len(jti))
		}
		if seen[jti] {
			t.Fatalf("jti collision after %d iterations: %s", i, jti)
		}
		seen[jti] = true
	}
}

func TestGenerateRefreshTokenLengthAndUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		token, err := generateRefreshToken()
		if err != nil {
			t.Fatalf("generateRefreshToken failed: %v", err)
		}
		// 32 random bytes encoded as hex => 64 characters
		if len(token) != 64 {
			t.Fatalf("expected refresh token length 64, got %d", len(token))
		}
		if seen[token] {
			t.Fatalf("refresh token collision after %d iterations", i)
		}
		seen[token] = true
	}
}

func TestSHA256HexDeterministic(t *testing.T) {
	// Known SHA-256 test vector for "abc"
	const want = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got := sha256Hex("abc"); got != want {
		t.Errorf("sha256Hex(\"abc\") = %s, want %s", got, want)
	}
	if sha256Hex("abc") != sha256Hex("abc") {
		t.Error("sha256Hex must be deterministic")
	}
	if sha256Hex("abc") == sha256Hex("abd") {
		t.Error("sha256Hex produced identical hashes for different inputs")
	}
	if got := sha256Hex(""); len(got) != 64 {
		t.Errorf("expected 64 hex chars for empty input, got %d", len(got))
	}
}

func TestGenerateJWTIncludesJTI(t *testing.T) {
	userID := "64b7f1a2c3d4e5f6a7b8c9d0"
	signed, err := generateJWT(userID)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret()), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		t.Fatalf("failed to parse generated JWT: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected MapClaims")
	}

	jti, ok := claims["jti"].(string)
	if !ok || len(jti) != 32 {
		t.Fatalf("expected 32-char jti claim, got %v", claims["jti"])
	}
	if claims["userId"] != userID {
		t.Errorf("expected userId %s, got %v", userID, claims["userId"])
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("expected numeric exp claim")
	}
	remaining := time.Until(time.Unix(int64(exp), 0))
	if remaining < 23*time.Hour || remaining > 24*time.Hour+time.Minute {
		t.Errorf("expected ~24h expiry, got %v remaining", remaining)
	}

	// jti must differ between tokens
	signed2, err := generateJWT(userID)
	if err != nil {
		t.Fatalf("second generateJWT failed: %v", err)
	}
	token2, err := jwt.Parse(signed2, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret()), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		t.Fatalf("failed to parse second JWT: %v", err)
	}
	claims2 := token2.Claims.(jwt.MapClaims)
	if claims2["jti"] == claims["jti"] {
		t.Error("expected unique jti per token")
	}
	if !strings.HasPrefix(signed2, "eyJ") {
		t.Error("expected a signed JWT string")
	}
}
