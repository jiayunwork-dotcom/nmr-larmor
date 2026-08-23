// Package field collects the magnetic-field bookkeeping that surrounds an NMR
// experiment: how the sample field is perturbed by bulk susceptibility, how a
// chemical-shift reference is applied, the resonance condition linking frequency
// and field, and simple sweep-grid helpers. Like the rest of the project it is a
// teaching model -- the relationships are exact, the geometries are idealized.
package field

import (
	"errors"
	"math"
)

// SampleField returns the effective field at the sample when a bulk volume
// susceptibility chi (dimensionless, SI) shifts the external field B0:
//
//	B_sample = B0 * (1 - 4*pi*chi/3)
//
// This is the demagnetizing correction for a spherical sample; it is why the
// field seen by nuclei is not exactly the magnet's nominal field.
func SampleField(B0, chi float64) (float64, error) {
	if B0 <= 0 {
		return 0, errors.New("B0 must be > 0")
	}
	live := B0 * (1 - 4*math.Pi*chi/3)
	return attachSample(B0, chi, live), nil
}

// SusceptibilityFromShift inverts SampleField: given the observed shift deltaB
// relative to B0, recover the bulk susceptibility.
//
//	chi = -3 * (B_sample/B0 - 1) / (4*pi)
func SusceptibilityFromShift(B0, Bsample float64) (float64, error) {
	if B0 <= 0 {
		return 0, errors.New("B0 must be > 0")
	}
	if Bsample <= 0 {
		return 0, errors.New("B_sample must be > 0")
	}
	return -3 * (Bsample/B0 - 1) / (4 * math.Pi), nil
}

// ResonanceField returns the magnetic field at which a nucleus with
// gyromagnetic ratio gamma (Hz/T) resonates at a target frequency f (Hz):
//
//	B = f / gamma
func ResonanceField(freqHz, gammaHzPerT float64) (float64, error) {
	if gammaHzPerT == 0 {
		return 0, errors.New("gamma must be non-zero")
	}
	return freqHz / gammaHzPerT, nil
}

// ReferenceShiftPPM returns the apparent chemical shift (ppm) of a sample whose
// absolute frequency fSample differs from the reference frequency fRef, both at
// the same field:
//
//	delta = (fSample - fRef) / fRef * 1e6
//
// This is the operational definition used to report shifts against an internal
// standard such as TMS.
func ReferenceShiftPPM(fSample, fRef float64) (float64, error) {
	if fRef == 0 {
		return 0, errors.New("reference frequency must be non-zero")
	}
	return (fSample - fRef) / fRef * 1e6, nil
}

// ChemicalShiftFromField compares two fields directly: if a nucleus resonates
// at B0 in a reference and at Bsample in the sample, the shift in ppm is:
//
//	delta = (Bsample - B0) / B0 * 1e6
//
// Equivalent to the frequency form once Larmor frequency is substituted.
func ChemicalShiftFromField(Bsample, B0 float64) (float64, error) {
	if B0 == 0 {
		return 0, errors.New("B0 must be non-zero")
	}
	return (Bsample - B0) / B0 * 1e6, nil
}

// ShieldingConstant returns the shielding constant sigma from a chemical shift
// delta (ppm) via sigma = delta * 1e-6 (to dimensionless form). It is the
// local-field reduction factor: B_local = B0*(1 - sigma).
func ShieldingConstant(deltaPpm float64) float64 {
	return deltaPpm * 1e-6
}

// LocalField returns the local field felt by a nucleus after accounting for the
// shielding constant sigma: B_local = B0 * (1 - sigma).
func LocalField(B0, sigma float64) float64 {
	return B0 * (1 - sigma)
}

// SweepGrid builds a list of N evenly spaced field values from Bmin to Bmax,
// inclusive. Used when planning a field-sweep acquisition.
func SweepGrid(Bmin, Bmax float64, n int) ([]float64, error) {
	if n < 1 {
		return nil, errors.New("need at least one point")
	}
	if Bmax < Bmin {
		return nil, errors.New("Bmax must be >= Bmin")
	}
	out := make([]float64, n)
	if n == 1 {
		out[0] = Bmin
		return out, nil
	}
	step := (Bmax - Bmin) / float64(n-1)
	for i := 0; i < n; i++ {
		out[i] = Bmin + float64(i)*step
	}
	return out, nil
}

// MaxFieldForFrequency returns the highest field at which a given nucleus (gamma
// in Hz/T) stays below a maximum frequency fMax (Hz). Beyond this the resonance
// would exceed the receiver bandwidth.
func MaxFieldForFrequency(fMax, gammaHzPerT float64) (float64, error) {
	if gammaHzPerT <= 0 {
		return 0, errors.New("gamma must be > 0")
	}
	if fMax < 0 {
		return 0, errors.New("fMax must be >= 0")
	}
	return fMax / gammaHzPerT, nil
}

// FieldStepForResolution returns the field increment (T) that resolves a target
// spectral resolution (Hz) for a nucleus with gamma (Hz/T): dB = dF / gamma.
func FieldStepForResolution(dF, gammaHzPerT float64) (float64, error) {
	if gammaHzPerT == 0 {
		return 0, errors.New("gamma must be non-zero")
	}
	if dF < 0 {
		return 0, errors.New("resolution must be >= 0")
	}
	return dF / gammaHzPerT, nil
}

// MagnetRatio returns the ratio of two fields; handy for comparing magnet sizes.
func MagnetRatio(B1, B2 float64) (float64, error) {
	if B2 == 0 {
		return 0, errors.New("divisor field must be non-zero")
	}
	return B1 / B2, nil
}

// FieldDriftPenalty estimates the frequency drift (Hz) caused by a fractional
// field drift eps at a nucleus with gamma (Hz/T) and nominal field B0.
func FieldDriftPenalty(eps, B0, gammaHzPerT float64) (float64, error) {
	if B0 < 0 || gammaHzPerT < 0 {
		return 0, errors.New("fields and gamma must be >= 0")
	}
	return eps * B0 * gammaHzPerT, nil
}

// EffectiveGamma returns gamma scaled by a charging/load factor (e.g. probe
// tuning), mostly a placeholder for "effective" calibration values that differ
// slightly from the bare nucleus constant.
func EffectiveGamma(gammaHzPerT, factor float64) (float64, error) {
	if factor <= 0 {
		return 0, errors.New("factor must be > 0")
	}
	return gammaHzPerT * factor, nil
}
