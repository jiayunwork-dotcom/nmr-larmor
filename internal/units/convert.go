package units

// PPMToHz converts a chemical shift expressed in ppm to an absolute frequency
// offset in hertz, given the reference (Larmor) frequency in hertz.
//
//	deltaHz = deltaPpm * f0Hz * 1e-6
func PPMToHz(deltaPpm, f0Hz float64) float64 {
	return deltaPpm * f0Hz * PPM
}

// HzToPPM is the inverse of PPMToHz:
//
//	deltaPpm = deltaHz / f0Hz / 1e-6
func HzToPPM(deltaHz, f0Hz float64) float64 {
	if f0Hz == 0 {
		return 0
	}
	return deltaHz / f0Hz / PPM
}

// PPMToRadPerT converts a gyromagnetic ratio expressed as gamma/(2*pi) in
// MHz/T into the absolute ratio gamma in rad/(s*T). This is the factor that
// turns a field (tesla) directly into an angular frequency (rad/s):
//
//	omega = gamma[rad/(s*T)] * B0
func PPMToRadPerT(gammaOver2PiMHzPerT float64) float64 {
	// MHz/T -> Hz/T, then Hz = rad/s / (2*pi).
	return gammaOver2PiMHzPerT * 1e6 * (2 * 3.141592653589793)
}
