package httpapi

import (
	"net/http"
	"strings"
)

// contentTypeFor guesses a useful Content-Type for a static file based on its
// extension. The browser console only needs a handful of types; everything else
// falls back to application/octet-stream so nothing is mis-rendered.
func contentTypeFor(name string) string {
	switch {
	case strings.HasSuffix(name, ".html"), strings.HasSuffix(name, ".htm"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(name, ".json"):
		return "application/json; charset=utf-8"
	case strings.HasSuffix(name, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// ensureContentType wraps a file server so responses carry an explicit content
// type, which matters for the example/*.json files the console fetches.
func ensureContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".") {
			w.Header().Set("Content-Type", contentTypeFor(r.URL.Path))
		}
		next.ServeHTTP(w, r)
	})
}
