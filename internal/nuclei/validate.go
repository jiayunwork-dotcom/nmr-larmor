package nuclei

import (
	"fmt"
	"math"
)

// MaxB0T is a sanity ceiling for the main magnetic field. Real NMR
// spectrometers top out well below this; values above it almost certainly
// indicate a unit or transcription error rather than a real magnet.
const MaxB0T = 60.0

// ValidateB0 checks that a magnetic field is physically usable for a Larmor
// resonance computation. The field must be strictly positive (B0 == 0 means
// there is no resonance to report) and below the sanity ceiling.
func ValidateB0(b0 float64) error {
	if math.IsNaN(b0) {
		return fmt.Errorf("B0 is not a number")
	}
	if b0 <= 0 {
		return fmt.Errorf("B0 must be greater than 0 (got %v)", b0)
	}
	if b0 > MaxB0T {
		return fmt.Errorf("B0 %.3f T exceeds the supported range (max %.0f T)", b0, MaxB0T)
	}
	return nil
}

// ValidateB0Range checks a [min,max] field window used by scans.
func ValidateB0Range(min, max float64) error {
	if math.IsNaN(min) || math.IsNaN(max) {
		return fmt.Errorf("B0 range contains a non-number")
	}
	if min <= 0 {
		return fmt.Errorf("B0_min must be greater than 0 (got %v)", min)
	}
	if max < min {
		return fmt.Errorf("B0_max (%.3f) must be >= B0_min (%.3f)", max, min)
	}
	return nil
}
