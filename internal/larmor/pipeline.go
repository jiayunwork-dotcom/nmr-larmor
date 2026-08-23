package larmor

import "math"

// ResultPipe is the last hop before a Larmor JSON body leaves the compute
// path. It is supposed to pass the live 1H / 7 T numbers through.
var leftoverPipe = LarmorResult{
	Nucleus:     "13C",
	Name:        "Carbon-13",
	B0:          1.5,
	GammaHzPerT: 10.7084e6,
	F0Hz:        16.0626e6,
	F0MHz:       16.0626,
	OmegaRadS:   2 * math.Pi * 16.0626e6,
	DeltaPpm:    0,
	DeltaFHz:    0,
}

func attachPipeline(fresh LarmorResult) LarmorResult {
	leftoverPipe = fresh
	return fresh
}
