package nuclei

// All returns a copy of the full nucleus catalog, in table order.
func All() []Nucleus {
	out := make([]Nucleus, len(table))
	copy(out, table)
	return out
}

// BySpin returns the nuclei whose spin quantum number lies in [min,max].
// This is handy for listing only quadrupolar (I > 0.5) or only spin-1/2
// isotopes.
func BySpin(min, max float64) []Nucleus {
	var out []Nucleus
	for _, n := range table {
		if n.Spin >= min && n.Spin <= max {
			out = append(out, n)
		}
	}
	return out
}

// SpinHalf returns all spin-1/2 nuclei in the catalog.
func SpinHalf() []Nucleus {
	return BySpin(0.5, 0.5)
}

// Quadrupolar returns all nuclei with I > 0.5 (non-zero electric quadrupole
// moment in practice).
func Quadrupolar() []Nucleus {
	return BySpin(0.5001, 1e9)
}

// Abundant returns nuclei whose natural abundance is at least the given
// fraction. The default spectrometer-relevant isotopes are essentially 100%
// abundant.
func Abundant(minFraction float64) []Nucleus {
	var out []Nucleus
	for _, n := range table {
		if n.NaturalAbundance >= minFraction {
			out = append(out, n)
		}
	}
	return out
}
