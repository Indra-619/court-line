package pricing

import (
	"math"
	"testing"
)

func TestCalculateHours(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		want      float64
		tolerance float64
		wantErr   bool
	}{
		{"half hour boundaries", "08:30", "10:00", 1.5, 1e-9, false},
		{"whole hour", "08:00", "09:00", 1.0, 1e-9, false},
		{"full day minus a minute", "00:00", "23:59", 23.983333333333334, 1e-6, false},
		{"quarter hour", "12:00", "12:15", 0.25, 1e-9, false},
		{"end before start", "10:00", "09:00", 0, 0, true},
		{"end equals start", "10:00", "10:00", 0, 0, true},
		{"garbage start", "abc", "10:00", 0, 0, true},
		{"garbage end", "10:00", "xyz", 0, 0, true},
		{"invalid hour", "25:00", "26:00", 0, 0, true},
		{"not zero-padded", "8:00", "10:00", 0, 0, true},
		{"empty strings", "", "", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateHours(tt.start, tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateHours(%q, %q) error = %v, wantErr %v", tt.start, tt.end, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Errorf("CalculateHours(%q, %q) = %v, want %v (tolerance %v)", tt.start, tt.end, got, tt.want, tt.tolerance)
			}
		})
	}
}
