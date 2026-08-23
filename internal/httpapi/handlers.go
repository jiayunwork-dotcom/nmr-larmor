package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"nmr-larmor/internal/larmor"
)

// handleLarmor implements POST /api/larmor.
//
// Request body (JSON):
//
//	{"nucleus":"1H","B0":7.0,"delta_ppm":4.5,"delta_unit":"ppm"}
//
// On success it returns a LarmorResult. On any validation failure (unknown
// nucleus, B0 <= 0, bad unit) it returns a 400 ErrorResponse whose message
// comes from the computation layer.
func handleLarmor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"POST is required for /api/larmor")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "could not read request body")
		return
	}
	var req larmor.LarmorRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request",
			"invalid JSON: "+err.Error())
		return
	}
	res, err := larmor.ComputeLarmor(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// handleScan implements POST /api/scan.
//
// Request body (JSON):
//
//	{"nucleus":"1H","B0_min":1.0,"B0_max":14.0,"B0_step":1.0}
//
// It returns a ScanResult with the Larmor frequency sampled across the field
// window. The same validation as /api/larmor applies to the nucleus, and the
// range/step are checked so an impossible scan is rejected up front.
func handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed",
			"POST is required for /api/scan")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "could not read request body")
		return
	}
	var req larmor.ScanRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request",
			"invalid JSON: "+err.Error())
		return
	}
	res, err := larmor.ScanB0(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}
