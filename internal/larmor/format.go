package larmor

import (
	"fmt"

	"nmr-larmor/internal/units"
)

// Describe returns a single-line human readable summary of a resonance, e.g.
// "1H @ 7.00 T: f0 = 298.0 MHz, delta 4.50 ppm -> 1341 Hz".
func (r LarmorResult) Describe() string {
	return fmt.Sprintf("%s @ %.2f T: f0 = %s, delta %.3g ppm -> %s",
		r.Nucleus, r.B0, units.FormatHz(r.F0Hz), r.DeltaPpm, units.FormatHz(r.DeltaFHz))
}

// DescribeScan returns a compact summary of a scan suitable for CLI output.
func (s ScanResult) DescribeScan() string {
	if len(s.Points) == 0 {
		return fmt.Sprintf("%s: empty scan", s.Nucleus)
	}
	first := s.Points[0]
	last := s.Points[len(s.Points)-1]
	return fmt.Sprintf("%s: %d points from %.2f T (%s) to %.2f T (%s)",
		s.Nucleus, len(s.Points),
		first.B0, units.FormatHz(first.F0Hz),
		last.B0, units.FormatHz(last.F0Hz))
}
