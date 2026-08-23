// Package larmor implements the core NMR Larmor-frequency accounting: given a
// magnetic field B0 and a nucleus (via its gamma/2pi), it computes the resonance
// frequency and the absolute frequency offset produced by a chemical shift.
//
// Two conventions exist in the literature, f = gamma*B0/(2*pi) and
// omega = gamma*B0; this package pins the linear-frequency form and derives the
// angular frequency from it, so the labels are always consistent:
//
//	f0   = (gamma/2pi) * B0
//	omega = 2*pi * f0
//	deltaF = deltaPpm * f0 * 1e-6
//
// Chemical shifts are stored in ppm, which is field independent; the absolute
// offset in hertz scales linearly with B0.
package larmor
