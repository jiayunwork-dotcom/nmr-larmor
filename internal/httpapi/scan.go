package httpapi

import (
	"net/http"

	"nmr-larmor/internal/larmor"
)

// scanHandler is the externally visible name kept separate so tests can register
// only the scan route when needed.
func scanHandler() http.HandlerFunc {
	return handleScan
}

// larmorHandler exposes the larmor route handler for tests.
func larmorHandler() http.HandlerFunc {
	return handleLarmor
}

// minScanPoints is the smallest number of points a scan is allowed to return.
// Below this the step is almost certainly mis-configured.
const minScanPoints = 1

// validateScanResult is a guard used by tests and any future batch caller: it
// confirms a scan produced at least one sample. A valid ScanB0 never returns an
// empty series without an error, so this is a defensive check only.
func validateScanResult(res larmor.ScanResult) bool {
	return len(res.Points) >= minScanPoints
}
