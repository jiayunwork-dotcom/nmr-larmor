package field

import (
	"math"
	"testing"
)

func TestSampleField(t *testing.T) {
	v, err := SampleField(10, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-10) > 1e-9 {
		t.Fatalf("chi=0 -> B unchanged, got %v", v)
	}
	v, _ = SampleField(10, 0.003)
	if v >= 10 {
		t.Fatalf("positive susceptibility lowers field under spherical demag, got %v", v)
	}
	if _, err := SampleField(0, 1); err == nil {
		t.Fatal("B0=0 should error")
	}
}

func TestSusceptibilityFromShift(t *testing.T) {
	// B_sample = 10*(1 - 4pi*chi/3); pick chi=0.001 -> B=10*(1-0.004189)
	B0 := 10.0
	bs := B0 * (1 - 4*math.Pi*0.001/3)
	chi, err := SusceptibilityFromShift(B0, bs)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(chi-0.001) > 1e-9 {
		t.Fatalf("recovered chi = %v, want 0.001", chi)
	}
}

func TestResonanceField(t *testing.T) {
	v, err := ResonanceField(400e6, 42.58e6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-9.39) > 1e-2 {
		t.Fatalf("B = %v, want ~9.39 T", v)
	}
}

func TestReferenceShiftPPM(t *testing.T) {
	v, err := ReferenceShiftPPM(400.004, 400)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-10) > 1e-9 {
		t.Fatalf("shift = %v, want 10 ppm", v)
	}
}

func TestChemicalShiftFromField(t *testing.T) {
	v, err := ChemicalShiftFromField(10.01, 10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-1000) > 1e-9 {
	t.Fatalf("shift = %v, want 1000 ppm", v)
	}
}

func TestShieldingAndLocal(t *testing.T) {
	sigma := ShieldingConstant(100)
	if math.Abs(sigma-1e-4) > 1e-12 {
		t.Fatalf("sigma = %v, want 1e-4", sigma)
	}
	bl := LocalField(10, sigma)
	if math.Abs(bl-10*(1-1e-4)) > 1e-9 {
		t.Fatalf("local field = %v", bl)
	}
}

func TestSweepGrid(t *testing.T) {
	g, err := SweepGrid(0, 10, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(g) != 3 || g[0] != 0 || g[2] != 10 {
		t.Fatalf("grid = %v", g)
	}
	if _, err := SweepGrid(5, 1, 3); err == nil {
		t.Fatal("Bmax<Bmin should error")
	}
}

func TestMaxFieldForFrequency(t *testing.T) {
	v, err := MaxFieldForFrequency(400e6, 42.58e6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-9.39) > 1e-2 {
		t.Fatalf("Bmax = %v, want ~9.39", v)
	}
}

func TestFieldStepForResolution(t *testing.T) {
	v, err := FieldStepForResolution(1, 42.58e6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-1/42.58e6) > 1e-12 {
		t.Fatalf("step = %v", v)
	}
}

func TestFieldDriftPenalty(t *testing.T) {
	v, err := FieldDriftPenalty(1e-6, 10, 42.58e6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-425.8) > 1e-6 {
		t.Fatalf("drift = %v, want 425.8 Hz", v)
	}
}

func TestEffectiveGamma(t *testing.T) {
	v, err := EffectiveGamma(42.58, 1.01)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-42.58*1.01) > 1e-9 {
		t.Fatalf("eff gamma = %v", v)
	}
	if _, err := MagnetRatio(10, 0); err == nil {
		t.Fatal("zero divisor should error")
	}
}
