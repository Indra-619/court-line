package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// exchangeCodeTTL is how long a one-time exchange code stays valid.
// It is a variable so tests can exercise the expiry logic.
var exchangeCodeTTL = 5 * time.Minute

type exchangeEntry struct {
	userID    string
	expiresAt time.Time
}

var (
	exchangeMu    sync.RWMutex
	exchangeStore = make(map[string]exchangeEntry)
)

func init() {
	go cleanupExpiredExchangeCodes()
}

// createExchangeCode issues a single-use, time-limited exchange code
// for the given user ID.
func createExchangeCode(userID string) string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		// crypto/rand failing means the runtime is broken; panic loudly
		// rather than issuing a predictable code.
		panic("crypto/rand unavailable: " + err.Error())
	}
	code := hex.EncodeToString(buf)

	exchangeMu.Lock()
	exchangeStore[code] = exchangeEntry{
		userID:    userID,
		expiresAt: time.Now().Add(exchangeCodeTTL),
	}
	exchangeMu.Unlock()

	return code
}

// consumeExchangeCode redeems a code exactly once. It returns ok=false
// when the code is unknown, expired, or already used.
func consumeExchangeCode(code string) (string, bool) {
	if code == "" {
		return "", false
	}

	exchangeMu.Lock()
	defer exchangeMu.Unlock()

	entry, found := exchangeStore[code]
	if !found {
		return "", false
	}
	delete(exchangeStore, code)

	if time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.userID, true
}

// cleanupExpiredExchangeCodes periodically drops expired entries so the
// store cannot grow without bound even if codes are never consumed.
func cleanupExpiredExchangeCodes() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		exchangeMu.Lock()
		for code, entry := range exchangeStore {
			if now.After(entry.expiresAt) {
				delete(exchangeStore, code)
			}
		}
		exchangeMu.Unlock()
	}
}
