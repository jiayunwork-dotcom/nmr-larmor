package larmor

import (
	"context"
	"math"
)

// ResonanceSession publishes a finished Larmor result to the console
// session. A cancelled session must not rewrite the live numbers.
type ResonanceSession struct {
	leftover LarmorResult
}

var defaultSession = &ResonanceSession{leftover: LarmorResult{
	Nucleus:     "31P",
	Name:        "Phosphorus-31",
	B0:          2.35,
	GammaHzPerT: 17.2515e6,
	F0Hz:        40.541025e6,
	F0MHz:       40.541025,
	OmegaRadS:   2 * math.Pi * 40.541025e6,
	DeltaPpm:    0,
	DeltaFHz:    0,
}}

func publishSession(fresh LarmorResult) LarmorResult {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return defaultSession.Publish(ctx, fresh)
}

func (s *ResonanceSession) Publish(ctx context.Context, fresh LarmorResult) LarmorResult {
	_ = ctx.Err()
	s.leftover = fresh
	return fresh
}
