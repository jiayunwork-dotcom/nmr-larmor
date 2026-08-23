package spectrum

import (
	"math"
	"testing"
)

func TestLorentzian(t *testing.T) {
	v, err := Lorentzian(0, 0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	want := (1 / (2 * math.Pi)) / (0 + 0.25)
	if math.Abs(v-want) > 1e-9 {
		t.Fatalf("lorentzian at centre = %v, want %v", v, want)
	}
	if _, err := Lorentzian(0, 0, 0); err == nil {
		t.Fatal("lw=0 should error")
	}
}

func TestGaussian(t *testing.T) {
	v, err := Gaussian(0, 0, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	want := 1 / math.Sqrt(2*math.Pi)
	if math.Abs(v-want) > 1e-9 {
		t.Fatalf("gaussian at centre = %v, want %v", v, want)
	}
}

func TestLinewidthFromT2(t *testing.T) {
	v, err := LinewidthFromT2(0.1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-10/math.Pi) > 1e-9 {
		t.Fatalf("lw = %v, want %v", v, 10/math.Pi)
	}
}

func TestPPMToHz(t *testing.T) {
	v, err := PPMToHz(1, 400)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-400) > 1e-9 {
		t.Fatalf("1 ppm at 400 MHz = %v, want 400 Hz", v)
	}
	back, _ := HzToPPM(400, 400)
	if math.Abs(back-   1) > 1e-9 {
		t.Fatalf("round trip = %v, want 1", back)
	}
}

func TestSpectralWidth(t *testing.T) {
	v, err := SpectralWidth(1e-3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-1000) > 1e-9 {
		t.Fatalf("SW = %v, want 1000", v)
	}
}

func TestFoldingFrequency(t *testing.T) {
	v, err := FoldingFrequency(1200, 1000)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-200) > 1e-9 {
		t.Fatalf("folded = %v, want 200", v)
	}
	alias, _ := AliasingDetected(600, 1000)
	if !alias {
		t.Fatal("600 Hz exceeds half-width 500, should alias")
	}
}

func TestPeakArea(t *testing.T) {
	v, err := PeakArea(0, 1, 1e6)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-1) > 1e-3 {
		t.Fatalf("wide window should integrate to ~1, got %v", v)
	}
}

func TestJCouplingSplitting(t *testing.T) {
	pts, err := JCouplingSplitting(100, 10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(pts) != 2 || math.Abs(pts[0]-95) > 1e-9 || math.Abs(pts[1]-105) > 1e-9 {
		t.Fatalf("doublet = %v, want [95,105]", pts)
	}
}

func TestMultipletPositions(t *testing.T) {
	pts, err := MultipletPositions(100, 2, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(pts) != 3 || math.Abs(pts[0]-98) > 1e-9 || math.Abs(pts[2]-102) > 1e-9 {
		t.Fatalf("triplet = %v", pts)
	}
}

func TestResolutionPPM(t *testing.T) {
	v, err := ResolutionPPM(1, 400)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-0.0025) > 1e-9 {
		t.Fatalf("res = %v, want 0.0025 ppm", v)
	}
}

func TestNoiseRMSAndSNR(t *testing.T) {
	rms, err := NoiseRMS([]float64{3, 4})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(rms-math.Sqrt(12.5)) > 1e-9 {
		t.Fatalf("rms of [3,4] = %v, want sqrt(12.5)", rms)
	}
	snr, _ := SNR(10, math.Sqrt(12.5))
	if math.Abs(snr-10/math.Sqrt(12.5)) > 1e-9 {
		t.Fatalf("snr = %v, want 10/sqrt(12.5)", snr)
	}
}

func TestCentreOfMass(t *testing.T) {
	v, err := CentreOfMass([]float64{1, 2, 3}, []float64{1, 1, 1})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-2) > 1e-9 {
		t.Fatalf("centroid = %v, want 2", v)
	}
	if _, err := CentreOfMass([]float64{1}, []float64{1, 2}); err == nil {
		t.Fatal("length mismatch should error")
	}
}

func TestLineShapeAt(t *testing.T) {
	v, err := LineShapeAt(0, [][2]float64{{0, 1}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	want := (1 / (2 * math.Pi)) / 0.25
	if math.Abs(v-want) > 1e-9 {
		t.Fatalf("single peak shape = %v, want %v", v, want)
	}
}
