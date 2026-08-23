package larmor

import "math"

// TwoPi is 2*pi.
const TwoPi = 2 * math.Pi

// PPM is the parts-per-million factor shared with the units package. It is
// re-declared here so the larmor package can compute chemical-shift offsets
// without importing units (keeping the dependency graph acyclic).
const PPM = 1e-6

// ScanStepEpsilon is the floating-point slack used when deciding whether the
// final B0 value of a scan has been reached. It prevents a step like 0.1 from
// stopping one point short of B0_max due to binary rounding.
const ScanStepEpsilon = 1e-9
