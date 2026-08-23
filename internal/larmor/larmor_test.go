package larmor

import (
	"math"
	"testing"
)

// TestLarmorB0ZeroError is the simple-level surface: B0 == 0 means there is no
// resonance to report and must be rejected, not returned as 0.
func TestLarmorB0ZeroError(t *testing.T) {
	_, err := ComputeLarmor(LarmorRequest{Nucleus: "1H", B0: 0, DeltaPpm: 0})
	if err == nil {
		t.Fatal("ComputeLarmor with B0=0 expected an error, got nil")
	}
}

// TestLarmor7TProtonMagnitude checks the headline invariant: a proton at 7 T
// resonates at roughly 300 MHz (298 MHz from the pinned ratio).
func TestLarmor7TProtonMagnitude(t *testing.T) {
	res, err := ComputeLarmor(LarmorRequest{Nucleus: "1H", B0: 7.0, DeltaPpm: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.F0MHz < 290 || res.F0MHz > 310 {
		t.Errorf("1H @ 7T f0 = %v MHz, want ~300 MHz", res.F0MHz)
	}
	// angular frequency must equal 2*pi*f0
	if math.Abs(res.OmegaRadS-2*math.Pi*res.F0Hz) > 1e-3 {
		t.Errorf("omega = %v, want 2*pi*f0 = %v", res.OmegaRadS, 2*math.Pi*res.F0Hz)
	}
}

// TestLarmorDeltaUnitMHzRejected checks the delta-unit guard: passing "MHz"
// (or any non-ppm unit) is a hard error, never a silent mis-interpretation.
func TestLarmorDeltaUnitMHzRejected(t *testing.T) {
	_, err := ComputeLarmor(LarmorRequest{Nucleus: "1H", B0: 7.0, DeltaPpm: 1, DeltaUnit: "MHz"})
	if err == nil {
		t.Fatal("DeltaUnit=MHz expected error, got nil")
	}
	// empty unit is allowed and treated as ppm
	if _, err := ComputeLarmor(LarmorRequest{Nucleus: "1H", B0: 7.0, DeltaPpm: 1}); err != nil {
		t.Errorf("empty DeltaUnit should be allowed: %v", err)
	}
}

// TestLarmorUnknownNucleusError checks an unknown nucleus surfaces an error.
func TestLarmorUnknownNucleusError(t *testing.T) {
	_, err := ComputeLarmor(LarmorRequest{Nucleus: "Xe-999", B0: 7.0})
	if err == nil {
		t.Fatal("unknown nucleus expected error")
	}
}
