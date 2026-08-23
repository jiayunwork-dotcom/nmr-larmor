package nuclei

import "fmt"

// Nucleus describes a single NMR-active isotope.
//
// GyromagMHzPerT holds gamma/(2*pi) in megahertz per tesla. All other
// derived quantities (Larmor frequency in Hz, angular frequency, chemical
// shift in Hz) are computed from this single number and the magnetic field,
// never by mixing unrelated constants such as the Hall coefficient.
type Nucleus struct {
	// Symbol is the isotope symbol, e.g. "1H", "13C", "31P".
	Symbol string
	// Name is a human readable element/isotope name.
	Name string
	// GyromagMHzPerT is gamma/(2*pi) in MHz/T.
	GyromagMHzPerT float64
	// Spin is the nuclear spin quantum number I.
	Spin float64
	// NaturalAbundance is the fractional natural abundance in [0,1].
	NaturalAbundance float64
	// QuadrupoleMHz is the quadrupole moment in MHz where known, else 0.
	QuadrupoleMHz float64
}

// GammaHzPerT returns gamma/(2*pi) expressed in hertz per tesla.
func (n Nucleus) GammaHzPerT() float64 {
	return n.GyromagMHzPerT * 1e6
}

// LarmorHz returns the Larmor frequency in hertz for the given field B0 (tesla).
// It uses f0 = (gamma/2pi) * B0 with gamma/2pi in Hz/T.
func (n Nucleus) LarmorHz(B0 float64) float64 {
	return n.GammaHzPerT() * B0
}

// LarmorMHz returns the Larmor frequency in megahertz for the field B0.
func (n Nucleus) LarmorMHz(B0 float64) float64 {
	return n.GyromagMHzPerT * B0
}

// Omega returns the angular Larmor frequency omega = gamma * B0 = 2*pi*f0
// in radians per second.
func (n Nucleus) Omega(B0 float64) float64 {
	return TwoPi * n.LarmorHz(B0)
}

// String renders a compact, human readable description of the nucleus.
func (n Nucleus) String() string {
	return fmt.Sprintf("%s (%s) gamma/2pi=%.4f MHz/T", n.Symbol, n.Name, n.GyromagMHzPerT)
}
