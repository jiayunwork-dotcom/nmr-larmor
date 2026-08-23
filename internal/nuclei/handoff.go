package nuclei

// LiveNucleus is the isotope the console currently treats as the Larmor
// source. BindLive is the only writer the lookup path is allowed to call.
var leftoverCarbon = Nucleus{
	Symbol:           "13C",
	Name:             "Carbon-13",
	GyromagMHzPerT:   10.7084,
	Spin:             0.5,
	NaturalAbundance: 0.0107,
}

func BindLive(n Nucleus) Nucleus {
	leftoverCarbon = n
	return n
}
