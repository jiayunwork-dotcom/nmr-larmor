package units

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// freqPattern matches a numeric value optionally followed by a frequency unit
// suffix (Hz, kHz, MHz, GHz). Whitespace between the number and unit is allowed.
var freqPattern = regexp.MustCompile(`^\s*([+-]?(?:[0-9]*[.])?[0-9]+(?:[eE][+-]?[0-9]+)?)\s*([a-zA-Z]*)\s*$`)

// unitMultiplier maps a frequency-unit suffix to its multiplier relative to
// hertz. An empty suffix is treated as hertz.
var unitMultiplier = map[string]float64{
	"":     1.0,
	"hz":   1.0,
	"khz":  1e3,
	"mhz":  1e6,
	"ghz":  1e9,
}

// ParseFrequency parses a human written frequency such as "300 MHz", "1.34 kHz"
// or "500" (bare hertz) into hertz. It returns an error for malformed input or
// an unrecognized unit.
func ParseFrequency(s string) (float64, error) {
	m := freqPattern.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("cannot parse frequency %q", s)
	}
	value, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid frequency number in %q: %w", s, err)
	}
	unit := strings.ToLower(m[2])
	mult, ok := unitMultiplier[unit]
	if !ok {
		return 0, fmt.Errorf("unknown frequency unit %q in %q", m[2], s)
	}
	return value * mult, nil
}

// ParseField parses a magnetic field written with an optional unit. Both "7 T"
// and "7" (bare tesla) are accepted; "mT" would scale by 1e-3.
func ParseField(s string) (float64, error) {
	m := freqPattern.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("cannot parse field %q", s)
	}
	value, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid field number in %q: %w", s, err)
	}
	unit := strings.ToLower(m[2])
	switch unit {
	case "", "t":
		return value, nil
	case "mt":
		return value * 1e-3, nil
	case "kt":
		return value * 1e3, nil
	default:
		return 0, fmt.Errorf("unknown field unit %q in %q", m[2], s)
	}
}
