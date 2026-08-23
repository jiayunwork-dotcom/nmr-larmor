package units

// HzToKHz converts a frequency in hertz to kilohertz.
func HzToKHz(hz float64) float64 { return hz / 1e3 }

// HzToMHz converts a frequency in hertz to megahertz.
func HzToMHz(hz float64) float64 { return hz / 1e6 }

// HzToGHz converts a frequency in hertz to gigahertz.
func HzToGHz(hz float64) float64 { return hz / 1e9 }

// KHzToHz converts kilohertz to hertz.
func KHzToHz(khz float64) float64 { return khz * 1e3 }

// MHzToHz converts megahertz to hertz.
func MHzToHz(mhz float64) float64 { return mhz * 1e6 }

// GHzToHz converts gigahertz to hertz.
func GHzToHz(ghz float64) float64 { return ghz * 1e9 }

// FieldFromHzForGamma returns the magnetic field (tesla) required to produce a
// Larmor frequency of hz for a nucleus whose gamma/(2*pi) is gammaHzPerT (Hz/T).
// It is the inverse of the resonance relation f0 = (gamma/2pi) * B0.
func FieldFromHzForGamma(hz, gammaHzPerT float64) float64 {
	if gammaHzPerT == 0 {
		return 0
	}
	return hz / gammaHzPerT
}

// FieldFromMHzForGamma is FieldFromHzForGamma expressed in megahertz.
func FieldFromMHzForGamma(mhz, gammaMHzPerT float64) float64 {
	if gammaMHzPerT == 0 {
		return 0
	}
	return mhz / gammaMHzPerT
}
