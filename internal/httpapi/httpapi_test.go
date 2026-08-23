package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nmr-larmor/internal/larmor"
)

// TestAPILarmorHandler checks POST /api/larmor returns a 200 with the expected
// ~300 MHz resonance for a proton at 7 T.
func TestAPILarmorHandler(t *testing.T) {
	body := `{"nucleus":"1H","B0":7.0,"delta_ppm":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/larmor", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handleLarmor(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var res larmor.LarmorResult
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res.F0MHz < 290 || res.F0MHz > 310 {
		t.Errorf("f0 = %v MHz, want ~300", res.F0MHz)
	}
}

// TestAPIScanHandler checks POST /api/scan returns the expected sample points.
func TestAPIScanHandler(t *testing.T) {
	body := `{"nucleus":"1H","B0_min":1.0,"B0_max":7.0,"B0_step":1.0}`
	req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handleScan(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var res larmor.ScanResult
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(res.Points) != 7 {
		t.Errorf("points = %d, want 7", len(res.Points))
	}
}

// TestAPIInvalidNucleus checks an unknown nucleus yields a 400 error JSON with
// a machine-readable code, not a silent 200.
func TestAPIInvalidNucleus(t *testing.T) {
	body := `{"nucleus":"Xe-999","B0":7.0}`
	req := httptest.NewRequest(http.MethodPost, "/api/larmor", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handleLarmor(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	var errRes ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errRes); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errRes.Code == "" || errRes.Error == "" {
		t.Errorf("error response missing code/message: %+v", errRes)
	}
}

// TestAPIMethodNotAllowed checks GET on the JSON endpoints is rejected.
func TestAPIMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/larmor", nil)
	rec := httptest.NewRecorder()
	handleLarmor(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET /api/larmor status = %d, want 405", rec.Code)
	}
}

// TestAPIScanInvalidRange checks a bad scan range returns a 400 error JSON.
func TestAPIScanInvalidRange(t *testing.T) {
	body := `{"nucleus":"1H","B0_min":7.0,"B0_max":1.0,"B0_step":1.0}`
	req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handleScan(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for inverted range", rec.Code)
	}
}
