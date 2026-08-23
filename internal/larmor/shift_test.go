package larmor

import (
	"math"
	"testing"

	"nmr-larmor/internal/nuclei"
)

// TestShiftZeroDelta checks the invariant delta = 0 => absolute offset 0 Hz.
func TestShiftZeroDelta(t *testing.T) {
	h, _ := nuclei.LookupStrict("1H")
	df := ComputeShift(h, 0, 7.0)
	if df != 0 {
		t.Errorf("ComputeShift(1H, 0 ppm, 7T) = %v Hz, want 0", df)
	}
}

// TestShiftPPMtoHzFactor is the high-difficulty surface: a 1 ppm shift at 7 T
// for the proton is about 298 Hz, never 298 MHz. The defining 1e-6 factor
// (ppm vs percent) must be present; dropping it inflates the separation by a
// factor of one million.
func TestShiftPPMtoHzFactor(t *testing.T) {
	h, _ := nuclei.LookupStrict("1H")
	df := ComputeShift(h, 1.0, 7.0)
	want := 42.577 * 7.0 // delta_ppm * (gamma/2pi in MHz/T) * B0
	if math.Abs(df-want) > 1e-6 {
		t.Errorf("1 ppm @ 7T 1H -> %v Hz, want %v Hz", df, want)
	}
	if df >= 1e6 {
		t.Errorf("deltaF = %v Hz is far too large; the 1e-6 ppm factor is likely missing", df)
	}
}

// TestShiftScalesWithField checks that a fixed ppm shift yields a larger hertz
// separation at higher field: doubling B0 doubles the hertz offset.
func TestShiftScalesWithField(t *testing.T) {
	h, _ := nuclei.LookupStrict("1H")
	low := ComputeShift(h, 5.0, 3.0)
	high := ComputeShift(h, 5.0, 6.0)
	if math.Abs(high/(low+1e-12)-2) > 1e-9 {
		t.Errorf("doubling B0 should double the Hz offset: low=%v high=%v", low, high)
	}
}
