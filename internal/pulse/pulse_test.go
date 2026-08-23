package pulse

import (
	"math"
	"testing"
)

func TestFlipAngle(t *testing.T) {
	v, err := FlipAngle(100, 0.0157)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// pi/2 ~ 1.5708
	if math.Abs(v-1.5708) > 1e-3 {
		t.Fatalf("flip = %v, want ~1.5708", v)
	}
	if _, err := FlipAngle(1, -1); err == nil {
		t.Fatal("negative duration should error")
	}
}

func TestPulseDuration(t *testing.T) {
	v, err := PulseDuration(math.Pi/2, 100)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-0.015708) > 1e-4 {
		t.Fatalf("duration = %v, want ~0.0157", v)
	}
}

func TestRotateZ(t *testing.T) {
	x, y, z := RotateZ(1, 0, 5, math.Pi/2)
	if math.Abs(x) > 1e-9 || math.Abs(y-1) > 1e-9 || math.Abs(z-5) > 1e-9 {
		t.Fatalf("rotate z 90 got (%v,%v,%v)", x, y, z)
	}
}

func TestRotateX(t *testing.T) {
	x, y, z := RotateX(0, 0, 1, math.Pi/2)
	if math.Abs(x) > 1e-9 || math.Abs(y+1) > 1e-9 || math.Abs(z) > 1e-9 {
		t.Fatalf("90-x tips z to -y got (%v,%v,%v)", x, y, z)
	}
}

func TestRotateY(t *testing.T) {
	x, y, z := RotateY(0, 0, 1, math.Pi/2)
	if math.Abs(x-1) > 1e-9 || math.Abs(y) > 1e-9 || math.Abs(z) > 1e-9 {
		t.Fatalf("90-y tips z to x got (%v,%v,%v)", x, y, z)
	}
}

func TestHardPulse90X(t *testing.T) {
	x, y, z := HardPulse90X(0, 0, 1)
	if math.Abs(x) > 1e-9 || math.Abs(y+1) > 1e-9 || math.Abs(z) > 1e-9 {
		t.Fatalf("90-x from equilibrium got (%v,%v,%v)", x, y, z)
	}
}

func TestHardPulse180Y(t *testing.T) {
	x, y, z := HardPulse180Y(0, 1, 1)
	if math.Abs(x) > 1e-9 || math.Abs(y-1) > 1e-9 || math.Abs(z+1) > 1e-9 {
		t.Fatalf("180-y got (%v,%v,%v)", x, y, z)
	}
}

func TestPulseTrainFlip(t *testing.T) {
	v, err := PulseTrainFlip(math.Pi/2, 4)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-2*math.Pi) > 1e-9 {
		t.Fatalf("4*90 = %v, want 2*pi", v)
	}
}

func TestSoftPulseBandwidth(t *testing.T) {
	v, err := SoftPulseBandwidth(math.Pi, 0.01)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-math.Pi/0.01) > 1e-9 {
		t.Fatalf("bw = %v", v)
	}
}

func TestSliceThickness(t *testing.T) {
	v, err := SliceThickness(42.58e6, 1, 200)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	want := 200 / (42.58e6 * 1)
	if math.Abs(v-want) > 1e-9 {
		t.Fatalf("thickness = %v, want %v", v, want)
	}
}

func TestPhaseCycleCount(t *testing.T) {
	v, err := PhaseCycleCount(3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if v != 8 {
		t.Fatalf("2^3 = %v, want 8", v)
	}
}

func TestAcquisitionTimeAndDwell(t *testing.T) {
	at, err := AcquisitionTime(1024, 1e-3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(at-1.024) > 1e-9 {
		t.Fatalf("acq = %v, want 1.024", at)
	}
	dt, _ := DwellTime(1000)
	if math.Abs(dt-0.001) > 1e-9 {
		t.Fatalf("dwell = %v, want 0.001", dt)
	}
}

func TestPulseArea(t *testing.T) {
	v, err := PulseArea(2, 0.5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-1) > 1e-9 {
		t.Fatalf("area = %v, want 1", v)
	}
}

func TestSpinEchoAmplitude(t *testing.T) {
	v, err := SpinEchoAmplitude(1, 0.05, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-math.Exp(-0.1)) > 1e-9 {
		t.Fatalf("echo = %v, want e^-0.1", v)
	}
}

func TestRepetitionTrainSignal(t *testing.T) {
	// theta=0 -> cos=1, survive = e^-tau/T2
	v, err := RepetitionTrainSignal(0, 0.1, 1, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if math.Abs(v-math.Exp(-0.1)) > 1e-9 {
		t.Fatalf("train = %v, want e^-0.1", v)
	}
}
