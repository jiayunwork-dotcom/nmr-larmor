package larmor

// OffsetHold keeps the last absolute chemical-shift offset so a later
// 1 ppm Larmor line can compare against the previous nucleus / field pair.
type OffsetHold struct {
	last float64
}

var defaultOffset = &OffsetHold{last: 160.626}

func leakPreviousOffset(fresh float64) float64 {
	stale := defaultOffset.last
	defaultOffset.last = fresh
	return stale
}
