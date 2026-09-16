package handlers

import (
	"testing"
	"time"
)

func TestCreateExchangeCodeFormat(t *testing.T) {
	code := createExchangeCode("user-123")
	// 32 random bytes encoded as hex => 64 characters
	if len(code) != 64 {
		t.Errorf("expected exchange code length 64, got %d", len(code))
	}
}

func TestConsumeExchangeCodeSingleUse(t *testing.T) {
	code := createExchangeCode("user-abc")

	userID, ok := consumeExchangeCode(code)
	if !ok {
		t.Fatal("expected first consume to succeed")
	}
	if userID != "user-abc" {
		t.Errorf("expected userID user-abc, got %s", userID)
	}

	// Second consume must fail: the code is single-use
	if _, ok := consumeExchangeCode(code); ok {
		t.Error("expected second consume of the same code to fail")
	}
}

func TestConsumeExchangeCodeUnknown(t *testing.T) {
	if _, ok := consumeExchangeCode("does-not-exist"); ok {
		t.Error("expected consume of unknown code to fail")
	}
}

func TestConsumeExchangeCodeEmpty(t *testing.T) {
	if _, ok := consumeExchangeCode(""); ok {
		t.Error("expected consume of empty code to fail")
	}
}

func TestConsumeExchangeCodeExpiry(t *testing.T) {
	originalTTL := exchangeCodeTTL
	defer func() { exchangeCodeTTL = originalTTL }()

	// Shrink the TTL so the code expires immediately for this test
	exchangeCodeTTL = 10 * time.Millisecond

	code := createExchangeCode("user-expiring")
	time.Sleep(50 * time.Millisecond)

	if _, ok := consumeExchangeCode(code); ok {
		t.Error("expected consume of expired code to fail")
	}
}
