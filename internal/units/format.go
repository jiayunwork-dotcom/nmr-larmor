package units

import "fmt"

// siPrefixes lists SI prefixes from hertz upward used by FormatHz.
var siPrefixes = []struct {
	multiplier float64
	unit       string
}{
	{1e9, "GHz"},
	{1e6, "MHz"},
	{1e3, "kHz"},
	{1.0, "Hz"},
}

// FormatHz renders a frequency in hertz using the largest SI prefix that keeps
// the numeric value at or above one, so 2.98e8 Hz becomes "298.04 MHz".
func FormatHz(hz float64) string {
	v, unit := SplitHz(hz)
	return fmt.Sprintf("%.4g %s", v, unit)
}

// SplitHz returns the value and unit string for the best SI representation of
// hz. It is exposed separately so callers can compose their own formatting.
func SplitHz(hz float64) (float64, string) {
	a := hz
	if a < 0 {
		a = -a
	}
	for _, p := range siPrefixes {
		if a >= p.multiplier || p.multiplier == 1.0 {
			return hz / p.multiplier, p.unit
		}
	}
	return hz, "Hz"
}

// FormatPPM renders a chemical shift with its unit.
func FormatPPM(ppm float64) string {
	return fmt.Sprintf("%.4g ppm", ppm)
}
