package httpapi

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the JSON body returned for every failed request. Both fields
// are produced by the backend so the web console and API clients can show a
// precise, user-facing message rather than a silent 200 with an empty table.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// writeJSON encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes an ErrorResponse with the supplied status and machine
// readable code. The message is safe to surface directly to users.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{Error: message, Code: code})
}
