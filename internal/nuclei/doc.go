// Package nuclei provides a small catalog of NMR-active nuclei together with
// their gyromagnetic ratios and helpers to look them up by symbol.
//
// The gyromagnetic ratio is stored as gamma/(2*pi) expressed in MHz per tesla,
// which is the form most directly useful for Larmor frequency calculations:
//
//	f0(Hz) = (gamma/2pi)[Hz/T] * B0
//
// For the proton (1H) the CODATA-pinned value 42.577 MHz/T is used so that a
// 7 T magnet lands in the familiar ~300 MHz spectrometer range.
package nuclei
