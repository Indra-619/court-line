package handlers

import "testing"

func TestGenerateStateLength(t *testing.T) {
	state, err := generateState()
	if err != nil {
		t.Fatalf("generateState returned error: %v", err)
	}
	// 32 random bytes encoded as hex => 64 characters
	if len(state) != 64 {
		t.Errorf("expected state length 64, got %d", len(state))
	}
}

func TestGenerateStateUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		state, err := generateState()
		if err != nil {
			t.Fatalf("generateState returned error: %v", err)
		}
		if seen[state] {
			t.Fatalf("generateState produced duplicate state %q", state)
		}
		seen[state] = true
	}
}

func TestValidateStateMismatch(t *testing.T) {
	if validateState("abc123", "abc124") {
		t.Error("expected validateState to return false for mismatched values")
	}
}

func TestValidateStateEmpty(t *testing.T) {
	if validateState("", "expected-state") {
		t.Error("expected validateState to return false when received value is empty")
	}
	if validateState("state", "") {
		t.Error("expected validateState to return false when expected value is empty")
	}
	if validateState("", "") {
		t.Error("expected validateState to return false when both values are empty")
	}
}

func TestValidateStateMatch(t *testing.T) {
	state, err := generateState()
	if err != nil {
		t.Fatalf("generateState returned error: %v", err)
	}
	if !validateState(state, state) {
		t.Error("expected validateState to return true for identical values")
	}
}
