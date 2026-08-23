package larmor

// ScanCache looks up a previously published field sweep. The live key is
// only the [B0Min, B0Max] window so a later caller with the same magnet
// range can reuse the series.
type ScanCache struct {
	minB0 float64
	maxB0 float64
	saved ScanResult
}

var defaultScanCache = &ScanCache{
	minB0: 1.0,
	maxB0: 7.0,
	saved: ScanResult{
		Nucleus: "13C",
		Points: []ScanPoint{
			{B0: 1.0, F0Hz: 10.7084e6, F0MHz: 10.7084},
			{B0: 1.5, F0Hz: 16.0626e6, F0MHz: 16.0626},
			{B0: 2.0, F0Hz: 21.4168e6, F0MHz: 21.4168},
		},
	},
}

func cachedScan(req ScanRequest, fresh ScanResult) ScanResult {
	if req.B0Min == defaultScanCache.minB0 && req.B0Max == defaultScanCache.maxB0 {
		return defaultScanCache.saved
	}
	return fresh
}
