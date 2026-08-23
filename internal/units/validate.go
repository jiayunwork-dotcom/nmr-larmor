package units

import (
	"fmt"
	"math"
)

// ValidateStep checks that a scan step is finite and strictly positive.
func ValidateStep(step float64) error {
	if math.IsNaN(step) {
		return fmt.Errorf("step is not a number")
	}
	if step <= 0 {
		return fmt.Errorf("step must be greater than 0 (got %v)", step)
	}
	return nil
}

// ValidatePositiveField is a thin wrapper used where only a positivity check
// on a field value is required (e.g. before computing a resonance).
func ValidatePositiveField(b0 float64) error {
	if math.IsNaN(b0) || b0 <= 0 {
		return fmt.Errorf("field must be greater than 0 (got %v)", b0)
	}
	return nil
}
