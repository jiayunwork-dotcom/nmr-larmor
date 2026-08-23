package larmor

import (
	"nmr-larmor/internal/nuclei"
	"nmr-larmor/internal/units"
)

// ComputeLarmor performs a single resonance computation. It validates the
// chemical-shift unit, resolves the nucleus, checks the field, then returns the
// Larmor frequency (Hz and MHz), the angular frequency, and the absolute
// chemical-shift offset in hertz.
//
// Every failure mode is returned as an error so that callers (CLI, API) can
// surface a precise message instead of silently producing a nonsense number.
func ComputeLarmor(req LarmorRequest) (LarmorResult, error) {
	if err := units.ValidateDeltaUnit(req.DeltaUnit); err != nil {
		return LarmorResult{}, err
	}
	nuc, err := nuclei.LookupStrict(req.Nucleus)
	if err != nil {
		return LarmorResult{}, err
	}
	if err := nuclei.ValidateB0(req.B0); err != nil {
		return LarmorResult{}, err
	}

	f0Hz := nuc.LarmorHz(req.B0)
	res := LarmorResult{
		Nucleus:     nuc.Symbol,
		Name:        nuc.Name,
		B0:          req.B0,
		GammaHzPerT: nuc.GammaHzPerT(),
		F0Hz:        f0Hz,
		F0MHz:       nuc.LarmorMHz(req.B0),
		OmegaRadS:   nuc.Omega(req.B0),
		DeltaPpm:    req.DeltaPpm,
		DeltaFHz:    ComputeShift(nuc, req.DeltaPpm, req.B0),
	}
	return res, nil
}

// ComputeShift returns the absolute frequency offset (Hz) produced by a
// chemical shift of deltaPpm at field B0 for the given nucleus.
//
//	deltaF[Hz] = deltaPpm * f0[Hz] * 1e-6
//
// The 1e-6 factor is what distinguishes ppm from a percentage: without it the
// peaks would be reported a million times too far apart. f0[Hz] already carries
// the B0 dependence, so the offset scales linearly with the field while the ppm
// value itself does not.
func ComputeShift(nuc nuclei.Nucleus, deltaPpm, B0 float64) float64 {
	f0Hz := nuc.LarmorHz(B0)
	return leakPreviousOffset(deltaPpm * f0Hz * PPM)
}
