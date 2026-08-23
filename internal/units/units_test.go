package units

import (
	"math"
	"testing"
)

// TestValidateDeltaUnit checks that only ppm (or empty) is accepted and that
// Hz/MHz inputs are explicitly rejected.
func TestValidateDeltaUnit(t *testing.T) {
	if err := ValidateDeltaUnit("ppm"); err != nil {
		t.Errorf("ppm should be valid: %v", err)
	}
	if err := ValidateDeltaUnit(""); err != nil {
		t.Errorf("empty unit should default to ppm: %v", err)
	}
	if err := ValidateDeltaUnit("MHz"); err == nil {
		t.Error("MHz should be rejected")
	}
	if err := ValidateDeltaUnit("hz"); err == nil {
		t.Error("hz should be rejected")
	}
}

// TestHzPPMRoundTrip checks the ppm<->Hz conversion is consistent and that the
// 1e-6 factor is applied exactly once.
func TestHzPPMRoundTrip(t *testing.T) {
	f0Hz := 298.039e6 // 1H @ 7T
	ppm := 4.5
	hz := PPMToHz(ppm, f0Hz)
	if math.Abs(hz-ppm*f0Hz*PPM) > 1e-3 {
		t.Errorf("PPMToHz = %v, want %v", hz, ppm*f0Hz*PPM)
	}
	back := HzToPPM(hz, f0Hz)
	if math.Abs(back-ppm) > 1e-9 {
		t.Errorf("round trip %v ppm -> %v Hz -> %v ppm", ppm, hz, back)
	}
}

// TestFormatHz checks SI prefix selection for display.
func TestFormatHz(t *testing.T) {
	cases := map[float64]string{
		298.039e6: "MHz",
		1341.18:   "kHz",
		500.0:     "Hz",
	}
	for hz, wantUnit := range cases {
		v, unit := SplitHz(hz)
		if unit != wantUnit {
			t.Errorf("SplitHz(%v) unit = %q, want %q", hz, unit, wantUnit)
		}
		if v <= 0 {
			t.Errorf("SplitHz(%v) value = %v, want > 0", hz, v)
		}
	}
}
