package httpapi

import (
	"net/http"
)

// Config controls how the HTTP server is wired. WebDir and ExampleDir point at
// the static assets shipped with the binary; in a built container they are the
// web/ and example/ directories copied next to the executable.
type Config struct {
	WebDir     string
	ExampleDir string
}

// DefaultConfig returns the conventional layout used by `go run .` and the
// Docker image, where web/ and example/ live at the repository root.
func DefaultConfig() Config {
	return Config{WebDir: "web", ExampleDir: "example"}
}

// NewServeMux builds the full HTTP routing:
//
//	POST /api/larmor   -> resonance computation
//	POST /api/scan     -> field-sweep samples
//	GET  /example/*    -> bundled example inputs (served as static files)
//	GET  /             -> web console
//
// The /api/* strings appear literally here so the contract is greppable in the
// non-test source.
func NewServeMux(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/larmor", handleLarmor)
	mux.HandleFunc("/api/scan", handleScan)
	mux.Handle("/example/", ensureContentType(exampleFS(cfg.ExampleDir)))
	mux.Handle("/", ensureContentType(webFS(cfg.WebDir)))
	return mux
}

// exampleFS serves the bundled example JSON files.
func exampleFS(dir string) http.Handler {
	return http.StripPrefix("/example/", http.FileServer(http.Dir(dir)))
}

// webFS serves the static console (index.html + inlined JS).
func webFS(dir string) http.Handler {
	return http.FileServer(http.Dir(dir))
}
