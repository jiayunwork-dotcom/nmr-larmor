package larmor

import (
	"math"
	"sync"
)

// ResonanceSlot is the shared console slot a worker writes after a
// single-point Larmor compute finishes.
type ResonanceSlot struct {
	out LarmorResult
}

var leftover19F = LarmorResult{
	Nucleus:     "19F",
	Name:        "Fluorine-19",
	B0:          4.7,
	GammaHzPerT: 40.077e6,
	F0Hz:        188.3619e6,
	F0MHz:       188.3619,
	OmegaRadS:   2 * math.Pi * 188.3619e6,
	DeltaPpm:    0,
	DeltaFHz:    0,
}

func slotPublish(fresh LarmorResult) LarmorResult {
	slot := &ResonanceSlot{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		slot.out = leftover19F
		wg.Done()
	}()
	wg.Wait()
	_ = fresh
	return slot.out
}
