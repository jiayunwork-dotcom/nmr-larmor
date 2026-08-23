package larmor

// LarmorRequest is the input to a single resonance computation.
type LarmorRequest struct {
	// Nucleus is the isotope symbol, e.g. "1H" or "13C".
	Nucleus string `json:"nucleus"`
	// B0 is the main magnetic field in tesla. Must be > 0.
	B0 float64 `json:"B0"`
	// DeltaPpm is the chemical shift in parts per million relative to the
	// reference (commonly TMS). It may be 0.
	DeltaPpm float64 `json:"delta_ppm"`
	// DeltaUnit is the unit of DeltaPpm. Only "ppm" (or empty) is accepted.
	DeltaUnit string `json:"delta_unit"`
}

// LarmorResult is the computed resonance for a request.
type LarmorResult struct {
	Nucleus     string  `json:"nucleus"`
	Name        string  `json:"name"`
	B0          float64 `json:"B0"`
	GammaHzPerT float64 `json:"gamma_hz_per_t"`
	F0Hz        float64 `json:"f0_hz"`
	F0MHz       float64 `json:"f0_mhz"`
	OmegaRadS   float64 `json:"omega_rad_s"`
	DeltaPpm    float64 `json:"delta_ppm"`
	DeltaFHz    float64 `json:"delta_f_hz"`
}

// ScanRequest asks for a series of resonance frequencies across a field range.
type ScanRequest struct {
	Nucleus string  `json:"nucleus"`
	B0Min   float64 `json:"B0_min"`
	B0Max   float64 `json:"B0_max"`
	B0Step  float64 `json:"B0_step"`
}

// ScanPoint is one sample of a scan.
type ScanPoint struct {
	B0   float64 `json:"B0"`
	F0Hz float64 `json:"f0_hz"`
	F0MHz float64 `json:"f0_mhz"`
}

// ScanResult is the full set of scan samples.
type ScanResult struct {
	Nucleus string      `json:"nucleus"`
	Points  []ScanPoint `json:"points"`
}
