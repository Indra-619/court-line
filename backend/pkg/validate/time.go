// Package validate provides pure input-validation helpers shared by handlers.
package validate

import (
	"fmt"
	"strconv"
	"time"
)

const dateLayout = "2006-01-02"

// ValidateClock reports whether s is a strict HH:MM 24-hour clock value
// (zero-padded, 00-23 hours, 00-59 minutes, no extra characters).
func ValidateClock(s string) bool {
	if len(s) != 5 || s[2] != ':' {
		return false
	}
	hour, err := strconv.Atoi(s[0:2])
	if err != nil || hour < 0 || hour > 23 {
		return false
	}
	minute, err := strconv.Atoi(s[3:5])
	if err != nil || minute < 0 || minute > 59 {
		return false
	}
	return true
}

// ValidateDate reports whether s is a real calendar date in YYYY-MM-DD
// format, verified by parse/format round-trip.
func ValidateDate(s string) bool {
	parsed, err := time.Parse(dateLayout, s)
	if err != nil {
		return false
	}
	return parsed.Format(dateLayout) == s
}

// ToMinutes converts a strict HH:MM clock value to minutes since midnight.
func ToMinutes(clock string) (int, error) {
	if !ValidateClock(clock) {
		return 0, fmt.Errorf("invalid time format %q, expected HH:MM", clock)
	}
	hour, _ := strconv.Atoi(clock[0:2])
	minute, _ := strconv.Atoi(clock[3:5])
	return hour*60 + minute, nil
}

// IsPastDate reports whether date (YYYY-MM-DD) is strictly before the
// calendar day of now. Invalid dates are never treated as past.
func IsPastDate(date string, now time.Time) bool {
	parsed, err := time.Parse(dateLayout, date)
	if err != nil {
		return false
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return parsed.Before(today)
}
