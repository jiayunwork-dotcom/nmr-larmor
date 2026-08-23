// Package pulse models the action of radio-frequency pulses and simple pulse
// sequences on a classical magnetization vector. It implements the rotation of
// the Bloch vector by a flip angle about a chosen axis, the relation between
// pulse power, duration, and flip angle, and the repeating structure of pulse
// trains used in imaging and multidimensional NMR. All angles are in radians and
// time in seconds unless stated otherwise.
package pulse

import (
	"errors"
	"math"
)

// FlipAngle returns the rotation applied by a pulse of amplitude omega1 (rad/s)
// applied for duration t (seconds):
//
//	theta = omega1 * t
//
// A 90-degree pulse is reached when omega1*t = pi/2; this is the definition used
// to calibrate pulse lengths on a given probe.
func FlipAngle(omega1, t float64) (float64, error) {
	if t < 0 {
		return 0, errors.New("duration must be >= 0")
	}
	return omega1 * t, nil
}

// PulseDuration returns the time needed to reach a target flip angle at a given
// RF field strength omega1 (rad/s): t = theta / omega1.
func PulseDuration(theta, omega1 float64) (float64, error) {
	if omega1 <= 0 {
		return 0, errors.New("RF field strength must be > 0")
	}
	return theta / omega1, nil
}

// RotateZ applies a rotation by angle theta about the z-axis to a 3D vector and
// returns the new components. The z component is unchanged; x and y rotate.
func RotateZ(x, y, z, theta float64) (float64, float64, float64) {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return x*c - y*s, x*s + y*c, z
}

// RotateX applies a rotation by angle theta about the x-axis. The x component is
// unchanged; y and z rotate. A 90-degree x-pulse tips longitudinal Mz into the
// transverse plane (Mz -> My).
func RotateX(x, y, z, theta float64) (float64, float64, float64) {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return x, y*c - z*s, y*s + z*c
}

// RotateY applies a rotation by angle theta about the y-axis. The y component is
// unchanged; x and z rotate.
func RotateY(x, y, z, theta float64) (float64, float64, float64) {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return x*c + z*s, y, -x*s + z*c
}

// HardPulse90X applies an ideal 90-degree pulse about x to a magnetization unit
// vector (Mx, My, Mz), tipping it from equilibrium (0,0,1) into the transverse
// plane. Returns the rotated vector.
func HardPulse90X(mx, my, mz float64) (float64, float64, float64) {
	return RotateX(mx, my, mz, math.Pi/2)
}

// HardPulse180Y applies an ideal 180-degree pulse about y, inverting My and Mz
// while leaving Mx untouched. Used as a refocusing pulse in spin echo.
func HardPulse180Y(mx, my, mz float64) (float64, float64, float64) {
	return RotateY(mx, my, mz, math.Pi)
}

// NutationFrequency returns the observed nutation (Rabi) frequency of a
// magnetization vector under continuous RF: omega_nut = gamma * B1, expressed
// here through the supplied omega1 directly. Provided for completeness; equals
// the RF amplitude itself.
func NutationFrequency(omega1 float64) float64 {
	return omega1
}

// PulseTrainFlip computes the cumulative flip angle after n identical pulses of
// angle theta, which is simply n*theta (mod 2*pi for the effective orientation).
func PulseTrainFlip(theta float64, n int) (float64, error) {
	if n < 0 {
		return 0, errors.New("pulse count must be >= 0")
	}
	return float64(n) * theta, nil
}

// SoftPulseExcitation returns the excitation profile width (Hz) for a soft pulse
// of duration t and flip angle theta: the bandwidth scales as theta/t. A longer,
// weaker pulse is more frequency-selective.
func SoftPulseBandwidth(theta, t float64) (float64, error) {
	if t <= 0 {
		return 0, errors.New("duration must be > 0")
	}
	return theta / t, nil
}

// GradientMoment returns the gradient moment (area under a gradient lobe) for a
// gradient GI (T/m) applied over time t at a fixed position r (m): k = gamma *
// G * t along the encoded direction is folded into the phase; here we return the
// geometric moment G*t used to choose slice thickness.
func GradientMoment(gradient, t float64) (float64, error) {
	if t < 0 {
		return 0, errors.New("duration must be >= 0")
	}
	return gradient * t, nil
}

// SliceThickness returns the slice thickness (m) selected by a gradient GI and
// RF bandwidth BW (Hz) at gamma (Hz/T): thickness = BW / (gamma * GI).
func SliceThickness(gammaHzPerT, gradient, bandwidth float64) (float64, error) {
	if gradient <= 0 || gammaHzPerT <= 0 {
		return 0, errors.New("gradient and gamma must be > 0")
	}
	if bandwidth < 0 {
		return 0, errors.New("bandwidth must be >= 0")
	}
	return bandwidth / (gammaHzPerT * gradient), nil
}

// PhaseCycleCount returns how many phase steps are needed to cancel an n-th
// order coherent artifact: 2^n steps. Phase cycling is the time-domain analogue
// of averaging with sign patterns.
func PhaseCycleCount(order int) (int, error) {
	if order < 0 {
		return 0, errors.New("order must be >= 0")
	}
	return 1 << uint(order), nil
}

// AcquisitionTime returns the total acquisition time for np complex points at a
// given dwell time dt: T_acq = np * dt. This bounds the observable evolution.
func AcquisitionTime(np, dt float64) (float64, error) {
	if np < 0 || dt <= 0 {
		return 0, errors.New("np must be >= 0 and dt > 0")
	}
	return np * dt, nil
}

// DwellTime returns the sampling interval for a given spectral width SW: dt =
// 1/SW. Faster sampling (smaller dt) widens the observable window.
func DwellTime(spectralWidth float64) (float64, error) {
	if spectralWidth <= 0 {
		return 0, errors.New("spectral width must be > 0")
	}
	return 1 / spectralWidth, nil
}

// PulseArea returns the integrated area of a rectangular pulse of amplitude a and
// duration t; area is proportional to flip angle for a hard pulse.
func PulseArea(amplitude, t float64) (float64, error) {
	if t < 0 {
		return 0, errors.New("duration must be >= 0")
	}
	return amplitude * t, nil
}

// RepetitionTrainSignal models the accumulated transverse dephasing across a
// train of n pulses spaced tau apart, with per-pulse flip theta, under T2 decay.
// Returns the surviving transverse magnitude relative to unity.
func RepetitionTrainSignal(theta, tau, T2 float64, n int) (float64, error) {
	if T2 <= 0 || n < 0 {
		return 0, errors.New("T2 must be > 0 and n >= 0")
	}
	survive := math.Exp(-tau / T2)
	// each pulse multiplies transverse coherence by cos(theta) (partial flip)
	return math.Pow(cosRetain(theta)*survive, float64(n)), nil
}

// cosRetain wraps the longitudinal retention factor cos(theta) for reuse.
func cosRetain(theta float64) float64 {
	return math.Cos(theta)
}

// SpinEchoAmplitude returns the amplitude of a spin echo at time 2*tau after a
// 90-x then 180-y pair, corrected only by T2 (dephasing from static field is
// refocused, so T2* does not enter).
func SpinEchoAmplitude(M0, tau, T2 float64) (float64, error) {
	if T2 <= 0 {
		return 0, errors.New("T2 must be > 0")
	}
	if tau < 0 {
		return 0, errors.New("tau must be >= 0")
	}
	return M0 * math.Exp(-2*tau/T2), nil
}
