package units

import "math"

// RoundToPlaces rounds x to the given number of decimal places.
func RoundToPlaces(x float64, places int) float64 {
	if places < 0 {
		places = 0
	}
	factor := math.Pow(10, float64(places))
	return math.Round(x*factor) / factor
}

// RoundHz rounds a frequency in hertz to the nearest integer hertz, which is
// appropriate when presenting spectrometer readouts.
func RoundHz(hz float64) float64 {
	return math.Round(hz)
}

// NearestPPM rounds a computed ppm value to a sensible precision for reporting
// (four decimal places), since chemical shifts are rarely quoted beyond that.
func NearestPPM(ppm float64) float64 {
	return RoundToPlaces(ppm, 4)
}

// Significant reports whether two frequencies are equal within a relative
// tolerance, used when asserting scaling invariants in tests and validations.
func Significant(a, b, relTol float64) bool {
	if a == b {
		return true
	}
	denom := math.Max(math.Abs(a), math.Abs(b))
	if denom == 0 {
		return true
	}
	return math.Abs(a-b)/denom <= relTol
}
