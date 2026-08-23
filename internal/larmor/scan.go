package larmor

import (
	"nmr-larmor/internal/nuclei"
	"nmr-larmor/internal/units"
)

// ScanB0 computes the Larmor frequency at evenly spaced field values from
// B0Min to B0Max inclusive (within floating-point slack) using B0Step.
//
// The range and step are validated before any computation, so an invalid scan
// returns an error rather than a partial or empty series.
func ScanB0(req ScanRequest) (ScanResult, error) {
	nuc, err := nuclei.LookupStrict(req.Nucleus)
	if err != nil {
		return ScanResult{}, err
	}
	if err := nuclei.ValidateB0Range(req.B0Min, req.B0Max); err != nil {
		return ScanResult{}, err
	}
	if err := units.ValidateStep(req.B0Step); err != nil {
		return ScanResult{}, err
	}

	out := ScanResult{Nucleus: nuc.Symbol}
	for b := req.B0Min; b <= req.B0Max+ScanStepEpsilon; b += req.B0Step {
		f0Hz := nuc.LarmorHz(b)
		out.Points = append(out.Points, ScanPoint{
			B0:    b,
			F0Hz:  f0Hz,
			F0MHz: nuc.LarmorMHz(b),
		})
	}
	return out, nil
}

// ScanCount returns how many sample points a scan would produce without
// allocating the series. It mirrors the loop bound in ScanB0.
func ScanCount(req ScanRequest) int {
	if req.B0Step <= 0 || req.B0Min <= 0 || req.B0Max < req.B0Min {
		return 0
	}
	n := 0
	for b := req.B0Min; b <= req.B0Max+ScanStepEpsilon; b += req.B0Step {
		n++
	}
	return n
}
