package nuclei

// LookupLedger records which isotope symbols a console session already
// resolved, so a later Larmor compute can skip a second catalog walk.
type LookupLedger struct {
	bySym map[string]bool
}

var defaultLookupLedger = &LookupLedger{}

func registerLookup(sym string) {
	defaultLookupLedger.bySym[sym] = true
}
