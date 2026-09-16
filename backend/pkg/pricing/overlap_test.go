package pricing

import "testing"

func TestOverlaps(t *testing.T) {
	tests := []struct {
		name         string
		startA, endA string
		startB, endB string
		want         bool
		wantErr      bool
	}{
		{"back-to-back slots do not overlap", "10:00", "12:00", "12:00", "14:00", false, false},
		{"back-to-back reversed does not overlap", "12:00", "14:00", "10:00", "12:00", false, false},
		{"partial overlap", "10:00", "12:00", "11:00", "13:00", true, false},
		{"partial overlap reversed", "11:00", "13:00", "10:00", "12:00", true, false},
		{"contained inside", "10:00", "14:00", "11:00", "12:00", true, false},
		{"containing", "11:00", "12:00", "10:00", "14:00", true, false},
		{"identical slots", "10:00", "12:00", "10:00", "12:00", true, false},
		{"same start different end", "10:00", "12:00", "10:00", "11:00", true, false},
		{"disjoint earlier", "14:00", "16:00", "10:00", "12:00", false, false},
		{"disjoint later", "10:00", "12:00", "14:00", "16:00", false, false},
		{"invalid start A", "25:00", "26:00", "10:00", "12:00", false, true},
		{"invalid end B", "10:00", "12:00", "10:00", "99:99", false, true},
		{"garbage", "abc", "12:00", "10:00", "12:00", false, true},
		{"end before start A", "12:00", "10:00", "10:00", "12:00", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Overlaps(tt.startA, tt.endA, tt.startB, tt.endB)
			if (err != nil) != tt.wantErr {
				t.Errorf("Overlaps(%q,%q,%q,%q) error = %v, wantErr %v", tt.startA, tt.endA, tt.startB, tt.endB, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Overlaps(%q,%q,%q,%q) = %v, want %v", tt.startA, tt.endA, tt.startB, tt.endB, got, tt.want)
			}
		})
	}
}
