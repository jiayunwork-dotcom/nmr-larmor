package nuclei

import (
	"fmt"
	"strings"
)

// NormalizeSymbol canonicalizes a user supplied isotope symbol. It trims
// surrounding whitespace, upper-cases the input and resolves common aliases
// (e.g. "H" -> "1H", "C" -> "13C", "proton" -> "1H").
func NormalizeSymbol(symbol string) string {
	raw := strings.TrimSpace(symbol)
	if raw == "" {
		return ""
	}
	upper := strings.ToUpper(raw)
	if canon, ok := aliases[upper]; ok {
		return canon
	}
	return upper
}

// Lookup returns the nucleus for a symbol together with a boolean indicating
// whether it was found in the catalog.
func Lookup(symbol string) (Nucleus, bool) {
	canon := NormalizeSymbol(symbol)
	n, ok := bySymbol[canon]
	return n, ok
}

// LookupStrict is like Lookup but returns an error when the symbol is unknown
// or empty. The error message is safe to surface to API/CLI users.
func LookupStrict(symbol string) (Nucleus, error) {
	canon := NormalizeSymbol(symbol)
	n, ok := bySymbol[canon]
	if !ok {
		return Nucleus{}, fmt.Errorf("unknown nucleus %q (known examples: 1H, 13C, 31P, 19F)", symbol)
	}
	return n, nil
}

// IsKnown reports whether symbol resolves to a catalogued nucleus.
func IsKnown(symbol string) bool {
	_, ok := Lookup(symbol)
	return ok
}
