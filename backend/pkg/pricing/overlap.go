package pricing

import (
	"fmt"

	"github.com/Indra-619/court-line/backend/pkg/validate"
)

// Overlaps reports whether two half-open HH:MM intervals [start, end)
// intersect. Back-to-back slots (end == next start) do NOT overlap.
func Overlaps(startA, endA, startB, endB string) (bool, error) {
	aStart, err := validate.ToMinutes(startA)
	if err != nil {
		return false, err
	}
	aEnd, err := validate.ToMinutes(endA)
	if err != nil {
		return false, err
	}
	bStart, err := validate.ToMinutes(startB)
	if err != nil {
		return false, err
	}
	bEnd, err := validate.ToMinutes(endB)
	if err != nil {
		return false, err
	}
	if aEnd <= aStart || bEnd <= bStart {
		return false, fmt.Errorf("interval end must be after start (%q-%q, %q-%q)", startA, endA, startB, endB)
	}
	return aStart < bEnd && bStart < aEnd, nil
}
