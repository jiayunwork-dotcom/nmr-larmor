package relaxation

import (
	"math"
	"testing"
)

func TestLongitudinalMagnetization(t *testing.T) {
	got, err := LongitudinalMagnetization(1, 0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-0) > 1e-9 {
		t.Fatalf("at t=0 should be 0, got %v", got)
	}
	got, _ = LongitudinalMagnetization(1, 1e9, 1)
	if math.Abs(got-1) > 1e-6 {
		t.Fatalf("large t should saturate to 1, got %v", got)
	}
	if _, err := LongitudinalMagnetization(1, 1, 0); err == nil {
		t.Fatal("T1=0 should error")
	}
}

func TestInversionRecoveryMz(t *testing.T) {
	got, err := InversionRecoveryMz(1, 0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got+1) > 1e-9 {
		t.Fatalf("at t=0 should be -1, got %v", got)
	}
	got, _ = InversionRecoveryMz(1, math.Log(2), 1)
	if math.Abs(got) > 1e-9 {
		t.Fatalf("null time should be 0, got %v", got)
	}
}

func TestTransverseMagnetization(t *testing.T) {
	got, err := TransverseMagnetization(2, 0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-2) > 1e-9 {
		t.Fatalf("at t=0 should be 2, got %v", got)
	}
	if _, err := TransverseMagnetization(1, -1, 1); err == nil {
		t.Fatal("negative time should error")
	}
}

func TestEffectiveT2Star(t *testing.T) {
	got, err := EffectiveT2Star(1, 0, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-1) > 1e-9 {
		t.Fatalf("no inhomogeneity -> T2*, got %v", got)
	}
	got, _ = EffectiveT2Star(1, 10, 0.01)
	want := 1 / (1 + 0.1)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("T2* = %v, want %v", got, want)
	}
}

func TestT1FromSaturationRecovery(t *testing.T) {
	// M/M0 = 1 - exp(-t/T1); choose t=1, T1=2 -> M/M0 = 1 - e^-0.5
	M0 := 1.0
	t1 := 2.0
	M := M0 * (1 - math.Exp(-1/t1))
	got, err := T1FromSaturationRecovery(M, M0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-t1) > 1e-9 {
		t.Fatalf("recovered T1 = %v, want %v", got, t1)
	}
	if _, err := T1FromSaturationRecovery( 0, 1, 1); err == nil {
		t.Fatal("M/M0=0 should error")
	}
}

func TestT2FromDecay(t *testing.T) {
	M0 := 1.0
	t2 := 0.5
	M := M0 * math.Exp(-1/t2)
	got, err := T2FromDecay(M, M0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-t2) > 1e-9 {
		t.Fatalf("recovered T2 = %v, want %v", got, t2)
	}
}

func TestErnstAngle(t *testing.T) {
	// TR=T1 -> arccos(e^-1) = 68.4... degrees
	got, err := ErnstAngleDeg(1, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-68.413) > 1e-2 {
		t.Fatalf("Ernst deg = %v, want ~68.4", got)
	}
	if _, err := ErnstAngleDeg(1, 0); err == nil {
		t.Fatal("T1=0 should error")
	}
}

func TestSteadyStateSignal(t *testing.T) {
	// theta=90deg, TR >> T1 -> ~1
	got, err := SteadyStateSignal(1e9, 1, math.Pi/2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-1) > 1e-3 {
		t.Fatalf("fully relaxed steady state should be ~1, got %v", got)
	}
}

func TestCPMGAmplitude(t *testing.T) {
	got, err := CPMGAmplitude(1, 0.01, 1, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-1) > 1e-9 {
		t.Fatalf("n=0 echo should be 1, got %v", got)
	}
	got, _ = CPMGAmplitude(1, 0.01, 1, 50)
	if got <= 0 {
		t.Fatalf("echo should be positive, got %v", got)
	}
}

func TestBuildupDecayCurve(t *testing.T) {
	b, err := BuildupCurve([]float64{0, 1, 10}, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(b) != 3 || b[0] != 0 {
		t.Fatalf("buildup wrong: %v", b)
	}
	d, err := DecayCurve([]float64{0}, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(d[0]-1) > 1e-9 {
		t.Fatalf("decay at 0 should be 1, got %v", d[0])
	}
}

func TestT1T2Ratio(t *testing.T) {
	got, err := T1T2Ratio(2, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-2) > 1e-9 {
		t.Fatalf("ratio = %v, want 2", got)
	}
	if _, err := T1T2Ratio(1, 2); err == nil {
		t.Fatal("T1<T2 should error")
	}
}

func TestRepetitionForRecovery(t *testing.T) {
	// f = 1 - e^-1 -> TR = T1
	got, err := RepetitionForRecovery(1-math.Exp(-1), 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(got-1) > 1e-9 {
		t.Fatalf("TR = %v, want 1", got)
	}
}

func TestSaturationFactorAndRate(t *testing.T) {
	if math.Abs(SaturationFactor(math.Pi/2)-0) > 1e-9 {
		t.Fatal("90-degree pulse leaves no longitudinal magnetization")
	}
	r, err := RelaxationRate(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(r-0.5) > 1e-9 {
		t.Fatalf("rate = %v, want 0.5", r)
	}
}
