package nuclei

import (
	"math"
	"testing"
)

// TestNucleiLookupProton checks the proton gyromagnetic ratio is the pinned
// CODATA value and that a 7 T resonance lands near 300 MHz.
func TestNucleiLookupProton(t *testing.T) {
	n, err := LookupStrict("1H")
	if err != nil {
		t.Fatalf("LookupStrict(1H) returned error: %v", err)
	}
	if math.Abs(n.GyromagMHzPerT-ProtonGyromagMHzPerT) > 1e-9 {
		t.Errorf("proton gamma/2pi = %v, want %v", n.GyromagMHzPerT, ProtonGyromagMHzPerT)
	}
	f0 := n.LarmorMHz(7.0)
	if f0 < 290 || f0 > 310 {
		t.Errorf("1H @ 7T f0 = %v MHz, want ~300 MHz", f0)
	}
}

// TestNucleiUnknownReturnsError checks that an unknown isotope is rejected with
// an error rather than a zero-value nucleus.
func TestNucleiUnknownReturnsError(t *testing.T) {
	_, err := LookupStrict("999X")
	if err == nil {
		t.Fatal("LookupStrict(999X) expected an error, got nil")
	}
	if IsKnown("999X") {
		t.Error("IsKnown(999X) = true, want false")
	}
}

// TestNucleiC13QuarterProton checks the carbon-13 ratio is approximately one
// quarter of the proton ratio (the cross-rule in the spec).
func TestNucleiC13QuarterProton(t *testing.T) {
	h, _ := LookupStrict("1H")
	c, _ := LookupStrict("13C")
	ratio := h.GyromagMHzPerT / c.GyromagMHzPerT
	if ratio < 3.9 || ratio > 4.1 {
		t.Errorf("1H/13C gamma ratio = %v, want ~4", ratio)
	}
}

// TestNucleiAliasResolution checks loose spellings normalize to catalogued
// symbols (so "proton", "C", "P31" all resolve).
func TestNucleiAliasResolution(t *testing.T) {
	cases := map[string]string{"proton": "1H", "C": "13C", "P31": "31P", "F": "19F", "h": "1H"}
	for in, want := range cases {
		n, err := LookupStrict(in)
		if err != nil {
			t.Errorf("LookupStrict(%q) error: %v", in, err)
			continue
		}
		if n.Symbol != want {
			t.Errorf("LookupStrict(%q) = %q, want %q", in, n.Symbol, want)
		}
	}
}

// TestNucleiValidateB0 checks the positivity requirement on the field.
func TestNucleiValidateB0(t *testing.T) {
	if err := ValidateB0(7.0); err != nil {
		t.Errorf("ValidateB0(7.0) unexpected error: %v", err)
	}
	if err := ValidateB0(0); err == nil {
		t.Error("ValidateB0(0) expected error (no resonance at zero field)")
	}
	if err := ValidateB0(-3); err == nil {
		t.Error("ValidateB0(-3) expected error")
	}
}
