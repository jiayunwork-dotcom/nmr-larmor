package larmor

import (
	"nmr-larmor/internal/nuclei"
)

// Reference describes the zero point of the chemical-shift scale. In practice
// tetramethylsilane (TMS) is the universal reference for 1H/13C NMR; the label
// is carried through so results can be annotated, but the shift arithmetic is
// performed entirely in ppm and never depends on the label.
type Reference struct {
	Label    string
	DeltaPpm float64
}

// TMS returns the standard TMS reference (zero ppm).
func TMS() Reference {
	return Reference{Label: "TMS", DeltaPpm: 0.0}
}

// RelativeTo returns the chemical shift of value (in ppm) measured relative to
// this reference. Because both are already in ppm, the operation is a simple
// subtraction and the result is independent of the magnetic field.
func (r Reference) RelativeTo(valuePpm float64) float64 {
	return valuePpm - r.DeltaPpm
}

// OffsetHz converts a reference-relative ppm value into an absolute hertz
// offset at field B0 for the given nucleus, reusing the core shift formula.
func (r Reference) OffsetHz(nuc nuclei.Nucleus, valuePpm, B0 float64) float64 {
	return ComputeShift(nuc, r.RelativeTo(valuePpm), B0)
}
