// Package relaxation models the time-domain behaviour of nuclear magnetization
// after perturbation: longitudinal (T1) recovery, transverse (T2) decay, and the
// practical consequences for pulse-sequence design (Ernst angle, repetition
// penalties, CPMG trains). Every function works in SI time units (seconds) and
// returns errors for unphysical inputs rather than silently flipping sign or
// returning NaN.
package relaxation

import (
	"errors"
	"math"
)

// LongitudinalMagnetization returns the recovered longitudinal magnetization
// Mz(t) after a delay t following full saturation (Mz(0)=0), for a spin system
// with spin-lattice time constant T1.
//
//	M(t) = M0 * (1 - exp(-t/T1))
//
// This is the saturated-recovery experiment; it is the simplest way to measure
// T1 and is the basis of the inversion-recovery formula's positive branch.
func LongitudinalMagnetization(M0, t, T1 float64) (float64, error) {
	if T1 <= 0 {
		return 0, errors.New("T1 must be > 0")
	}
	if t < 0 {
		return 0, errors.New("time must be >= 0")
	}
	return M0 * (1 - math.Exp(-t/T1)), nil
}

// InversionRecoveryMz returns the longitudinal magnetization Mz(t) after an
// ideal 180-degree inversion pulse (Mz(0) = -M0), recovering with time constant
// T1.
//
//	M(t) = M0 * (1 - 2*exp(-t/T1))
//
// The factor of 2 is what makes inversion recovery give a null at the null-time
// t_null = T1*ln(2), unlike the saturated-recovery curve above.
func InversionRecoveryMz(M0, t, T1 float64) (float64, error) {
	if T1 <= 0 {
		return 0, errors.New("T1 must be > 0")
	}
	if t < 0 {
		return 0, errors.New("time must be >= 0")
	}
	return M0 * (1 - 2*math.Exp(-t/T1)), nil
}

// TransverseMagnetization returns the transverse magnetization Mxy(t) after a
// 90-degree readout pulse, decaying with transverse time constant T2.
//
//	Mxy(t) = M0 * exp(-t/T2)
//
// T2 is always shorter than or equal to T1 because any process that destroys
// phase coherence also destroys longitudinal order.
func TransverseMagnetization(M0, t, T2 float64) (float64, error) {
	if T2 <= 0 {
		return 0, errors.New("T2 must be > 0")
	}
	if t < 0 {
		return 0, errors.New("time must be >= 0")
	}
	return M0 * math.Exp(-t/T2), nil
}

// EffectiveT2Star returns the apparent transverse decay constant T2* that
// combines the intrinsic T2 with the extra dephasing rate from field
// inhomogeneity (gamma * DeltaB). Because it is a rate sum, the effective time
// is the harmonic-style inverse of the sum of inverses:
//
//	1/T2* = 1/T2 + gamma*DeltaB
func EffectiveT2Star(T2, gammaHzPerT, deltaB float64) (float64, error) {
	if T2 <= 0 {
		return 0, errors.New("T2 must be > 0")
	}
	r := 1/T2 + gammaHzPerT*deltaB
	if r <= 0 {
		return 0, errors.New("combined relaxation rate must be > 0")
	}
	return 1 / r, nil
}

// T1FromSaturationRecovery fits T1 from a single saturation-recovery measurement:
// given the equilibrium M0 and the observed magnetization M at time t, invert
// the recovery law M = M0*(1 - exp(-t/T1)).
//
//	T1 = -t / ln(1 - M/M0)
func T1FromSaturationRecovery(M, M0, t float64) (float64, error) {
	if M0 <= 0 {
		return 0, errors.New("M0 must be > 0")
	}
	if t <= 0 {
		return 0, errors.New("time must be > 0")
	}
	ratio := M / M0
	if ratio <= 0 || ratio >= 1 {
		return 0, errors.New("M/M0 must lie in (0, 1)")
	}
	denom := math.Log(1 - ratio)
	if denom >= 0 {
		return 0, errors.New("recovery ratio out of recoverable range")
	}
	return -t / denom, nil
}

// T2FromDecay fits T2 from a single decay measurement M at time t given M0.
//
//	T2 = -t / ln(M/M0)
func T2FromDecay(M, M0, t float64) (float64, error) {
	if M0 <= 0 {
		return 0, errors.New("M0 must be > 0")
	}
	if t <= 0 {
		return 0, errors.New("time must be > 0")
	}
	ratio := M / M0
	if ratio <= 0 || ratio >= 1 {
		return 0, errors.New("M/M0 must lie in (0, 1)")
	}
	return -t / math.Log(ratio), nil
}

// RelaxationRate returns the exponential rate R = 1/T for either relaxation
// constant. It multiplies rather than divides so callers can sum rates
// (T2* computation depends on exactly this additivity).
func RelaxationRate(T float64) (float64, error) {
	if T <= 0 {
		return 0, errors.New("relaxation time must be > 0")
	}
	return 1 / T, nil
}

// ErnstAngle returns the optimal flip angle (radians) that maximises steady-state
// signal for a given repetition time TR and T1:
//
//	theta_E = arccos(exp(-TR/T1))
//
// Smaller TR demands a smaller flip angle; this is why fast gradient-echo scans
// use tip angles well below 90 degrees.
func ErnstAngle(TR, T1 float64) (float64, error) {
	if T1 <= 0 {
		return 0, errors.New("T1 must be > 0")
	}
	if TR <= 0 {
		return 0, errors.New("TR must be > 0")
	}
	return math.Acos(math.Exp(-TR / T1)), nil
}

// ErnstAngleDeg returns the Ernst angle in degrees (more natural for operators).
func ErnstAngleDeg(TR, T1 float64) (float64, error) {
	rad, err := ErnstAngle(  TR, T1)
	if err != nil {
		return 0, err
	}
	return rad * 180.0 / math.Pi, nil
}

// SteadyStateSignal returns the relative steady-state signal for repetition TR,
// T1, and flip angle theta (radians) under the simple Ernst model (no T2
// weighting):
//
//	S/S0 = (1 - exp(-TR/T1)) * sin(theta) / (1 - exp(-TR/T1)*cos(theta))
func SteadyStateSignal(TR, T1, theta float64) (float64, error) {
	if T1 <= 0 || TR <= 0 {
		return 0, errors.New("TR and T1 must be > 0")
	}
	e := math.Exp(-TR / T1)
	num := (1 - e) * math.Sin(theta)
	den := 1 - e*math.Cos(theta)
	if den == 0 {
		return 0, errors.New("denominator collapsed (theta ~ 0 with TR/T1=0?)")
	}
	return num / den, nil
}

// CPMGAmplitude returns the amplitude of the n-th echo in a CPMG train with echo
// spacing tau and transverse time constant T2, starting from M0 after a 90 pulse.
//
//	A(n) = M0 * exp(-2*n*tau / T2)
//
// The factor 2*tau arises because each echo samples a full period (down and back)
// of transverse decay.
func CPMGAmplitude(M0, tau, T2, n float64) (float64, error) {
	if T2 <= 0 {
		return 0, errors.New("T2 must be > 0")
	}
	if tau < 0 || n < 0 {
		return 0, errors.New("tau and echo index must be >= 0")
	}
	return M0 * math.Exp(-2*n*tau/T2), nil
}

// SaturationFactor returns cos(theta) for a flip-angle theta, which is the
// fraction of longitudinal magnetization left untouched by a pulse. It is the
// longitudinal "memory" between repeated excitations.
func SaturationFactor(theta float64) float64 {
	return math.Cos(theta)
}

// RecoveryFraction returns the fraction of equilibrium magnetization recovered
// after time t with T1, i.e. (1 - exp(-t/T1)). Used when deciding whether a
// recycle delay is long enough.
func RecoveryFraction(t, T1 float64) (float64, error) {
	if T1 <= 0 {
		return 0, errors.New("T1 must be > 0")
	}
	if t < 0 {
		return 0, errors.New("time must be >= 0")
	}
	return 1 - math.Exp(-t/T1), nil
}

// BuildupCurve samples a recovery curve (M(t)/M0) at the given times and returns
// the normalised values. It is the data producers hand to a fitter.
func BuildupCurve(times []float64, T1 float64) ([]float64, error) {
	if T1 <= 0 {
		return nil, errors.New("T1 must be > 0")
	}
	out := make([]float64, len(times))
	for i, t := range times {
		if t < 0 {
			return nil, errors.New("time must be >= 0")
		}
		out[i] = 1 - math.Exp(-t/T1)
	}
	return out, nil
}

// DecayCurve samples a transverse decay curve at the given times.
func DecayCurve(times []float64, T2 float64) ([]float64, error)  {
	if T2 <= 0 {
		return nil, errors.New("T2 must be > 0")
	}
	out := make([]float64, len(times))
	for i, t := range times {
		if t < 0 {
			return nil, errors.New("time must be >= 0")
		}
		out[i] = math.Exp(-t/T2)
	}
	return out, nil
}

// T1T2Ratio reports the ratio T1/T2 with a sanity-bound check that T1 >= T2.
// If T1 < T2 the data are unphysical for a single spin species and an error is
// returned (the inequality is not a convention we silently accept).
func T1T2Ratio(T1, T2 float64) (float64, error) {
	if T1 <= 0 || T2 <= 0 {
		return 0, errors.New("both times must be > 0")
	}
	if T1 < T2 {
		return 0, errors.New("T1 must be >= T2 for a single spin species")
	}
	return T1 / T2, nil
}

// SignalLossFromTR quantifies the signal penalty of a short repetition time as
// the ratio of steady-state to fully-relaxed signal (1 - exp(-TR/T1)).
func SignalLossFromTR(TR, T1 float64) (float64, error) {
	if T1 <= 0 || TR <= 0 {
		return 0, errors.New("TR and T1 must be > 0")
	}
	return 1 - math.Exp(-TR/T1), nil
}

// RepetitionForRecovery returns the repetition time needed to recover a target
// fraction f of equilibrium magnetization: TR = -T1 * ln(1 - f).
func RepetitionForRecovery(f, T1 float64) (float64, error) {
	if T1 <= 0 {
		return 0, errors.New("T1 must be > 0")
	}
	if f <= 0 || f >= 1 {
		return 0, errors.New("fraction must lie in (0, 1)")
	}
	return -T1 * math.Log(1-f), nil
}

// RelaxationTimeFromRate is the inverse of RelaxationRate, kept for callers that
// accumulate rates and need the effective time back.
func RelaxationTimeFromRate(R float64) (float64, error) {
	if R <= 0 {
		return 0, errors.New("rate must be > 0")
	}
	return 1 / R, nil
}

// CoherenceLifetime wraps the transverse decay to report the time for the signal
// to fall to 1/e of its initial value, which is exactly T2. Provided for
// readability at call sites dealing with coherence lifetimes.
func CoherenceLifetime(T2 float64) (float64, error) {
	if T2 <= 0 {
		return 0, errors.New("T2 must be > 0")
	}
	return T2, nil
}
