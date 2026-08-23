package nuclei

// TwoPi is 2*pi, used when converting between linear frequency and angular
// frequency. It is defined here so the nuclei package does not need to import
// the larmor package (which would create an import cycle).
const TwoPi = 2 * 3.141592653589793

// ProtonGyromagMHzPerT is the CODATA-pinned gyromagnetic ratio of the proton
// expressed as gamma/(2*pi) in MHz/T. It is the anchor for the whole table:
// many other isotopes are referenced against it and a 7 T proton resonance
// therefore lands at 42.577 * 7 ~= 298 MHz.
const ProtonGyromagMHzPerT = 42.577

// table is the built-in catalog of NMR-active nuclei. Values of
// gamma/(2*pi) are CODATA-style figures in MHz/T. The proton is pinned to the
// project constant; the rest are accurate to the quoted precision and cover the
// isotopes most commonly probed in solution-state NMR.
var table = []Nucleus{
	{Symbol: "1H", Name: "Hydrogen", GyromagMHzPerT: ProtonGyromagMHzPerT, Spin: 0.5, NaturalAbundance: 0.999885, QuadrupoleMHz: 0},
	{Symbol: "2H", Name: "Deuterium", GyromagMHzPerT: 6.53590532, Spin: 1, NaturalAbundance: 0.000115, QuadrupoleMHz: 0.00226},
	{Symbol: "3H", Name: "Tritium", GyromagMHzPerT: 45.4148, Spin: 0.5, NaturalAbundance: 0, QuadrupoleMHz: 0},
	{Symbol: "3He", Name: "Helium-3", GyromagMHzPerT: 32.434, Spin: 0.5, NaturalAbundance: 1.34e-6, QuadrupoleMHz: 0},
	{Symbol: "13C", Name: "Carbon-13", GyromagMHzPerT: 10.7084, Spin: 0.5, NaturalAbundance: 0.0107, QuadrupoleMHz: 0},
	{Symbol: "15N", Name: "Nitrogen-15", GyromagMHzPerT: -4.3156, Spin: 0.5, NaturalAbundance: 0.00365, QuadrupoleMHz: 0},
	{Symbol: "19F", Name: "Fluorine-19", GyromagMHzPerT: 40.077, Spin: 0.5, NaturalAbundance: 1.0, QuadrupoleMHz: 0},
	{Symbol: "23Na", Name: "Sodium-23", GyromagMHzPerT: 11.269, Spin: 1.5, NaturalAbundance: 1.0, QuadrupoleMHz: 1.08},
	{Symbol: "27Al", Name: "Aluminium-27", GyromagMHzPerT: 11.103, Spin: 2.5, NaturalAbundance: 1.0, QuadrupoleMHz: 1.49},
	{Symbol: "29Si", Name: "Silicon-29", GyromagMHzPerT: -8.465, Spin: 0.5, NaturalAbundance: 0.0467, QuadrupoleMHz: 0},
	{Symbol: "31P", Name: "Phosphorus-31", GyromagMHzPerT: 17.2515, Spin: 0.5, NaturalAbundance: 1.0, QuadrupoleMHz: 0},
	{Symbol: "7Li", Name: "Lithium-7", GyromagMHzPerT: 16.546, Spin: 1.5, NaturalAbundance: 0.9275, QuadrupoleMHz: -0.04},
	{Symbol: "11B", Name: "Boron-11", GyromagMHzPerT: 13.663, Spin: 1.5, NaturalAbundance: 0.801, QuadrupoleMHz: 2.59},
	{Symbol: "35Cl", Name: "Chlorine-35", GyromagMHzPerT: 4.172, Spin: 1.5, NaturalAbundance: 0.7576, QuadrupoleMHz: -0.0789},
	{Symbol: "79Br", Name: "Bromine-79", GyromagMHzPerT: 10.704, Spin: 1.5, NaturalAbundance: 0.5069, QuadrupoleMHz: 0.313},
	{Symbol: "129Xe", Name: "Xenon-129", GyromagMHzPerT: -11.777, Spin: 0.5, NaturalAbundance: 0.264, QuadrupoleMHz: 0},
	{Symbol: "195Pt", Name: "Platinum-195", GyromagMHzPerT: 9.094, Spin: 0.5, NaturalAbundance: 0.3383, QuadrupoleMHz: 0},
	{Symbol: "199Hg", Name: "Mercury-199", GyromagMHzPerT: 7.612, Spin: 0.5, NaturalAbundance: 0.1696, QuadrupoleMHz: 0},
	{Symbol: "89Y", Name: "Yttrium-89", GyromagMHzPerT: 2.0858, Spin: 0.5, NaturalAbundance: 1.0, QuadrupoleMHz: 0},
	{Symbol: "17O", Name: "Oxygen-17", GyromagMHzPerT: -5.771, Spin: 2.5, NaturalAbundance: 0.00038, QuadrupoleMHz: -1.89},
}

// bySymbol indexes the table by normalized symbol for O(1) lookups.
var bySymbol = buildIndex()

func buildIndex() map[string]Nucleus {
	idx := make(map[string]Nucleus, len(table))
	for _, n := range table {
		idx[n.Symbol] = n
	}
	return idx
}

// aliases maps common loose spellings to canonical symbols.
var aliases = map[string]string{
	"H":   "1H",
	"H1":  "1H",
	"PROTON": "1H",
	"P":   "31P",
	"P31": "31P",
	"C":   "13C",
	"C13": "13C",
	"F":   "19F",
	"F19": "19F",
	"N":   "15N",
	"N15": "15N",
	"NA":  "23Na",
	"NA23": "23Na",
	"AL":  "27Al",
	"AL27": "27Al",
	"LI":  "7Li",
	"LI7": "7Li",
	"SI":  "29Si",
	"SI29": "29Si",
	"CL":  "35Cl",
	"CL35": "35Cl",
	"BR":  "79Br",
	"BR79": "79Br",
}
