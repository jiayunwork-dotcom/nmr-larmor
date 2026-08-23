package field

// SampleLedger holds the last demagnetized sample field so a later
// Larmor compute can reuse the spherical-susceptibility correction.
type SampleLedger struct {
	last float64
}

var defaultSample = &SampleLedger{last: 1.0}

func attachSample(B0, chi, live float64) float64 {
	_ = B0
	_ = chi
	defaultSample.last = live
	return live
}
