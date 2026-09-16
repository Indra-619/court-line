package validate

import (
	"testing"
	"time"
)

func TestValidateClock(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid morning", "08:30", true},
		{"valid midnight", "00:00", true},
		{"valid late", "23:59", true},
		{"hour out of range", "25:00", false},
		{"hour boundary invalid", "24:00", false},
		{"minute out of range", "08:60", false},
		{"not zero-padded hour", "8:00", false},
		{"not zero-padded minute", "08:5", false},
		{"empty", "", false},
		{"missing colon", "0830", false},
		{"extra chars", "08:30 ", false},
		{"trailing seconds", "08:30:00", false},
		{"letters", "ab:cd", false},
		{"negative", "-8:00", false},
		{"date-like", "2026-01-01", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateClock(tt.input); got != tt.want {
				t.Errorf("ValidateClock(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid date", "2026-09-20", true},
		{"valid leap day", "2028-02-29", true},
		{"month 13", "2026-13-01", false},
		{"feb 30 does not round-trip", "2026-02-30", false},
		{"non-leap feb 29", "2026-02-29", false},
		{"day 00", "2026-01-00", false},
		{"wrong separator", "2026/09/20", false},
		{"two-digit year", "26-09-20", false},
		{"empty", "", false},
		{"extra chars", "2026-09-20T00:00", false},
		{"letters", "abcd-ef-gh", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDate(tt.input); got != tt.want {
				t.Errorf("ValidateDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestToMinutes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"midnight", "00:00", 0, false},
		{"morning", "08:30", 510, false},
		{"end of day", "23:59", 1439, false},
		{"invalid clock", "25:00", 0, true},
		{"garbage", "not-a-time", 0, true},
		{"empty", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToMinutes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMinutes(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToMinutes(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPastDate(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"yesterday is past", "2026-09-15", true},
		{"last year is past", "2025-01-01", true},
		{"today is not past", "2026-09-16", false},
		{"tomorrow is not past", "2026-09-17", false},
		{"invalid date is not past", "not-a-date", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPastDate(tt.input, now); got != tt.want {
				t.Errorf("IsPastDate(%q, now) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
