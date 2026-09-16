// Package pricing provides booking duration and overlap helpers built on
// strict HH:MM clock values.
package pricing

import (
	"fmt"

	"github.com/your-username/book-lapangan/backend/pkg/validate"
)

// CalculateHours returns the exact duration in hours between two HH:MM
// clock values. It errors on invalid input or when end <= start.
func CalculateHours(startTime, endTime string) (float64, error) {
	startMinutes, err := validate.ToMinutes(startTime)
	if err != nil {
		return 0, err
	}
	endMinutes, err := validate.ToMinutes(endTime)
	if err != nil {
		return 0, err
	}
	if endMinutes <= startMinutes {
		return 0, fmt.Errorf("end time %q must be after start time %q", endTime, startTime)
	}
	return float64(endMinutes-startMinutes) / 60.0, nil
}
