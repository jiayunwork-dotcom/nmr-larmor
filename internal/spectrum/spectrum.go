// Package spectrum models the frequency-domain appearance of an NMR signal: the
// line shape produced by a damped coherence (Lorentzian and Gaussian), the
// relation between linewidth and T2, the conversion between ppm and hertz across
// a spectrometer frequency, and practical bookkeeping for peak lists. All widths
// and offsets are handled in a single, consistent convention so that peak
// integration and assignment do not mix frames.
package spectrum

import (
	"errors"
	"math"
)

// Lorentzian returns the normalized Lorentzian line shape value at frequency f
// relative to a centre f0, with full-width-at-half-maximum lw (in the same
// frequency units as f and f0).
//
//	L(f) = (lw/2pi) / ((f-f0)^2 + (lw/2)^2)
//
// The integral over all f is 1 for positive lw.
func Lorentzian(f, f0, lw float64) (float64, error) {
	if lw <= 0 {
		return 0, errors.New("linewidth must be > 0")
	}
	d := f - f0
	return (lw / (2 * math.Pi)) / (d*d + (lw/2)*(lw/2)), nil
}

// Gaussian returns the normalized Gaussian line shape at f around f0 with
// standard deviation sigma (same units).
//
//	G(f) = exp(-(f-f0)^2 / (2 sigma^2)) / (sigma sqrt(2 pi))
func Gaussian(f, f0, sigma float64) (float64, error) {
	if sigma <= 0 {
		return 0, errors.New("sigma must be > 0")
	}
	d := f - f0
	return math.Exp(-d*d/(2*sigma*sigma)) / (sigma*math.Sqrt(2*math.Pi)), nil
}

// LinewidthFromT2 returns the full-width-at-half-maximum (Hz) of a Lorentzian
// whose decay is governed by T2:
//
//	lw = 1 / (pi * T2)
//
// This is the natural linewidth of an ideal homogeneous line; extra broadening
// only adds linearly to 1/T2, not to the width directly.
func LinewidthFromT2(T2 float64) (float64, error) {
	if T2 <= 0 {
		return 0, errors.New("T2 must be > 0")
	}
	return 1 / (math.Pi * T2), nil
}

// PPMToHz converts a chemical shift in ppm to an absolute frequency offset (Hz)
// at spectrometer frequency f0MHz (MHz).
//
//	deltaHz = deltaPpm * f0MHz * 1e6 * 1e-6 = deltaPpm * f0MHz
//
// The 1e6 (ppm) and the MHz-to-Hz factor cancel to a clean multiply by f0MHz,
// which is exactly why ppm is such a convenient dimensionless shift unit.
func PPMToHz(deltaPpm, f0MHz float64) (float64, error) {
	if f0MHz <= 0 {
		return 0, errors.New("spectrometer frequency must be > 0")
	}
	return deltaPpm * f0MHz, nil
}

// HzToPPM converts an absolute frequency offset (Hz) back to ppm at f0MHz (MHz).
func HzToPPM(deltaHz, f0MHz float64) (float64, error) {
	if f0MHz <= 0 {
		return 0, errors.New("spectrometer frequency must be > 0")
	}
	return deltaHz / f0MHz, nil
}

// SpectralWidth returns the total swept width (Hz) for dwell time dt (seconds):
//
//	SW = 1 / dt
//
// It is the Nyquist span and sets the maximum unambiguous frequency.
func SpectralWidth(dwellTime float64) (float64, error) {
	if dwellTime <= 0 {
		return 0, errors.New("dwell time must be > 0")
	}
	return 1 / dwellTime, nil
}

// FoldingFrequency returns the alias (folded) frequency when an NMR signal at
// true frequency f appears in a spectrum of width SW (centred at 0).
//
//	folded = f - round(f / SW) * SW
func FoldingFrequency(f, SW float64) (float64, error) {
	if SW <= 0 {
		return 0, errors.New("spectral width must be > 0")
	}
	return f - math.Round(f/SW)*SW, nil
}

// AliasingDetected reports whether f folds (appears at a different location than
// its true value) within the window of width SW.
func AliasingDetected(f, SW float64) (bool, error) {
	if SW <= 0 {
		return false, errors.New("spectral width must be > 0")
	}
	return math.Abs(f) > SW/2, nil
}

// PeakArea integrates a Lorentzian peak analytically over a symmetric window of
// half-width w around f0. For a normalized Lorentzian the area is the fraction
// of total intensity inside the window.
func PeakArea(f0, lw, halfWindow float64) (float64, error) {
	if lw <= 0 || halfWindow < 0 {
		return 0, errors.New("linewidth must be > 0 and window >= 0")
	}
	// integrate L from f0-w to f0+w: (2/pi)*arctan(2w/lw)
	return (2.0 / math.Pi) * math.Atan(2*halfWindow/lw), nil
}

// TotalArea sums the integrated area of several peaks (already-normalized
// contributions in [0,1]); returns the cumulative intensity.
func TotalArea(areas []float64) (float64, error) {
	if len(areas) == 0 {
		return 0, errors.New("no peaks supplied")
	}
	sum := 0.0
	for _, a := range areas {
		if a < 0 {
			return 0, errors.New("negative peak area")
		}
		sum += a
	}
	return sum, nil
}

// PeakSymmetry measures left/right asymmetry of a peak from two heights sampled
// at equal offsets +/- d from the centre. A perfectly symmetric peak yields 0.
func PeakSymmetry(heightLeft, heightRight float64) float64 {
	return (heightRight - heightLeft) / (heightLeft + heightRight + 1e-12)
}

// JCouplingSplitting returns the two satellite frequencies for a doublet split
// by J (Hz) around a centre f0: [f0 - J/2, f0 + J/2].
func JCouplingSplitting(f0, J float64) ([]float64, error) {
	if J < 0 {
		return nil, errors.New("J coupling must be >= 0")
	}
	return []float64{f0 - J/2, f0 + J/2}, nil
}

// MultipletPositions returns equidistant multiplet lines for n equivalent
// couplings of size J around f0. n=1 -> singlet (just f0); n=2 -> doublet; etc.
func MultipletPositions(f0, J float64, n int) ([]float64, error) {
	if n < 1 {
		return nil, errors.New("number of lines must be >= 1")
	}
	if J < 0 {
		return nil, errors.New("J coupling must be >= 0")
	}
	out := make([]float64, n)
	start := -float64(n-1) / 2.0
	for i := 0; i < n; i++ {
		out[i] = f0 + (start+float64(i))*J
	}
	return out, nil
}

// ResolutionPPM reports the achievable chemical-shift resolution (ppm) given the
// linewidth in Hz and the spectrometer frequency (MHz): res = lwHz / f0MHz.
func ResolutionPPM(lwHz, f0MHz float64) (float64, error) {
	if f0MHz <= 0 {
		return 0, errors.New("spectrometer frequency must be > 0")
	}
	return lwHz / f0MHz, nil
}

// BaselineFlat returns True when the endpoints of a spectrum are within tol of
// each other, indicating a flat (un-drifted) baseline.
func BaselineFlat(values []float64, tol float64) (bool, error) {
	if len(values) < 2 {
		return false, errors.New("need at least two points")
	}
	return math.Abs(values[0]-values[len(values)-1]) <= tol, nil
}

// NoiseRMS returns the root-mean-square of a noise trace, the standard measure
// of spectral noise used to set a detection threshold.
func NoiseRMS(noise []float64) (float64, error) {
	if len(noise) == 0 {
		return 0, errors.New("empty noise trace")
	}
	var sum float64
	for _, v := range noise {
		sum += v * v
	}
	return math.Sqrt(sum / float64(len(noise))), nil
}

// SNR returns the signal-to-noise ratio of a peak height above a given noise RMS.
func SNR(peakHeight, noiseRMS float64) (float64, error) {
	if noiseRMS <= 0 {
		return 0, errors.New("noise RMS must be > 0")
	}
	return peakHeight / noiseRMS, nil
}

// LineShapeAt returns the composite line shape value at f by summing individual
// Lorentzian peaks described by (centre, linewidth) pairs. Used to render a
// synthetic spectrum from a peak list.
func LineShapeAt(f float64, peaks [][2]float64) (float64, error) {
	total := 0.0
	for _, p := range peaks {
		v, err := Lorentzian(f, p[0], p[1])
		if err != nil {
			return 0, err
		}
		total += v
	}
	return total, nil
}

// CentreOfMass returns the ppm-weighted centre of a set of (position, height)
// pairs -- the "centroid" of a peak cluster, robust to slight asymmetry.
func CentreOfMass(positions, heights []float64) (float64, error) {
	if len(positions) != len(heights) {
		return 0, errors.New("positions and heights length mismatch")
	}
	var num, den float64
	for i, p := range positions {
		num += p * heights[i]
		den += heights[i]
	}
	if den == 0 {
		return 0, errors.New("total height is zero")
	}
	return num / den, nil
}
