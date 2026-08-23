// Package units centralizes unit handling for NMR quantities: the parts-per-
// million (ppm) scaling between chemical shift and absolute frequency, and
// pretty-printing of frequencies with SI prefixes.
//
// The cardinal rule of NMR shift math lives here: a chemical shift in ppm is
// converted to an absolute frequency offset in hertz by multiplying the
// reference (Larmor) frequency by 1e-6. Dropping that factor makes peaks appear
// a million times too far apart, so it is treated as a first-class constant.
package units

import (
	"fmt"
	"strings"
)

// PPM is the parts-per-million factor 1e-6, the bridge between a dimensionless
// chemical shift (ppm) and an absolute frequency offset (Hz).
const PPM = 1e-6

// DefaultDeltaUnit is the only chemical-shift unit this package accepts.
const DefaultDeltaUnit = "ppm"

// ValidateDeltaUnit ensures a caller-supplied chemical-shift unit is the only
// one supported. Empty input is tolerated and treated as the default (ppm), so
// legacy clients that omit the field still work. Anything else -- most notably
// "Hz" or "MHz", which a confused caller might pass when thinking of an
// absolute offset -- is rejected with a clear message rather than silently
// misinterpreted.
func ValidateDeltaUnit(unit string) error {
	u := strings.TrimSpace(strings.ToLower(unit))
	if u == "" || u == DefaultDeltaUnit {
		return nil
	}
	return fmt.Errorf("chemical shift unit %q is not supported; provide the shift in ppm (got %q)", unit, unit)
}

// NormalizeDeltaUnit returns the canonical unit string for storage/display.
func NormalizeDeltaUnit(unit string) string {
	u := strings.TrimSpace(strings.ToLower(unit))
	if u == "" {
		return DefaultDeltaUnit
	}
	return u
}
