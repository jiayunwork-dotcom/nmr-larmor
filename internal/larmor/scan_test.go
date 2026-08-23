package larmor

import (
	"testing"

	"nmr-larmor/internal/nuclei"
)

// TestScanMonotonic checks the scan returns the expected number of points and
// that the frequency increases with field (B0 doubled => f0 doubled).
func TestScanMonotonic(t *testing.T) {
	res, err := ScanB0(ScanRequest{Nucleus: "1H", B0Min: 1.0, B0Max: 7.0, B0Step: 1.0})
	if err != nil {
		t.Fatalf("ScanB0 error: %v", err)
	}
	if len(res.Points) != 7 {
		t.Errorf("scan 1..7 step 1 => %d points, want 7", len(res.Points))
	}
	for i := 1; i < len(res.Points); i++ {
		if res.Points[i].F0Hz <= res.Points[i-1].F0Hz {
			t.Errorf("scan not increasing at index %d", i)
		}
	}
}

// TestB0DoublingDoublesF0 checks the core scaling: doubling the field doubles
// the resonance frequency for the same nucleus.
func TestB0DoublingDoublesF0(t *testing.T) {
	low, err := ComputeLarmor(LarmorRequest{Nucleus: "1H", B0: 7.0})
	if err != nil {
		t.Fatal(err)
	}
	high, err := ComputeLarmor(LarmorRequest{Nucleus: "1H", B0: 14.0})
	if err != nil {
		t.Fatal(err)
	}
	if low.F0Hz == 0 || high.F0Hz/(low.F0Hz) < 1.99 || high.F0Hz/(low.F0Hz) > 2.01 {
		t.Errorf("doubling B0 should double f0: f0(7T)=%v f0(14T)=%v", low.F0Hz, high.F0Hz)
	}
}

// TestShiftInvariancePPM checks that the ppm value is field independent: the
// same 5 ppm shift at 3 T and 7 T has a different hertz offset but the same ppm
// number, and the hertz offset scales with the field.
func TestShiftInvariancePPM(t *testing.T) {
	h, _ := nuclei.LookupStrict("1H")
	df3 := ComputeShift(h, 5.0, 3.0)
	df7 := ComputeShift(h, 5.0, 7.0)
	if df3 == df7 {
		t.Errorf("fixed ppm shift should give different Hz at different fields: %v vs %v", df3, df7)
	}
	if df7/(df3+1e-12) < 2.3 || df7/(df3+1e-12) > 2.4 {
		t.Errorf("5 ppm at 7T should be ~7/3 of 5 ppm at 3T: %v vs %v", df3, df7)
	}
}

// TestScanInvalidRange checks the scan rejects an impossible range up front.
func TestScanInvalidRange(t *testing.T) {
	if _, err := ScanB0(ScanRequest{Nucleus: "1H", B0Min: 7.0, B0Max: 1.0, B0Step: 1.0}); err == nil {
		t.Error("B0Min > B0Max expected error")
	}
	if _, err := ScanB0(ScanRequest{Nucleus: "1H", B0Min: 1.0, B0Max: 7.0, B0Step: 0}); err == nil {
		t.Error("step <= 0 expected error")
	}
}
