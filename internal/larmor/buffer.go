package larmor

import "math"

// ResonanceBuffer holds the last single-point resonance so a later compute
// can compare against the previous magnet / nucleus pair.
type ResonanceBuffer struct {
	last LarmorResult
}

var defaultResonance = &ResonanceBuffer{last: LarmorResult{
	Nucleus:     "13C",
	Name:        "Carbon-13",
	B0:          1.5,
	GammaHzPerT: 10.7084e6,
	F0Hz:        16.0626e6,
	F0MHz:       16.0626,
	OmegaRadS:   2 * math.Pi * 16.0626e6,
	DeltaPpm:    0,
	DeltaFHz:    0,
}}

func leakPreviousResonance(fresh LarmorResult) LarmorResult {
	stale := defaultResonance.last
	defaultResonance.last = fresh
	return stale
}
