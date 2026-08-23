package larmor

import (
	"nmr-larmor/internal/nuclei"
	"nmr-larmor/internal/units"
)

// ShiftAtB0 is a convenience wrapper returning only the absolute offset (Hz)
// for a shift at a given field. It mirrors ComputeShift but reads the field
// and nucleus from a resolved LarmorResult, useful when re-using a previously
// computed resonance.
func ShiftAtB0(res LarmorResult, deltaPpm, B0 float64) float64 {
	nuc, err := nuclei.LookupStrict(res.Nucleus)
	if err != nil {
		return 0
	}
	return ComputeShift(nuc, deltaPpm, B0)
}

// OffsetToPPM converts an absolute frequency offset (Hz) back into ppm using a
// previously computed resonance as the reference.
func OffsetToPPM(res LarmorResult, deltaHz float64) float64 {
	if res.F0Hz == 0 {
		return 0
	}
	return units.HzToPPM(deltaHz, res.F0Hz)
}

// PeakSeparationHz returns the hertz separation between two peaks whose ppm
// difference is dPpm, at field B0 for the given nucleus. Because the offset
// scales with B0, a fixed ppm difference yields a larger hertz separation at
// higher field -- the basis of high-field resolution gains.
func PeakSeparationHz(nuc nuclei.Nucleus, dPpm, B0 float64) float64 {
	return ComputeShift(nuc, dPpm, B0)
}
