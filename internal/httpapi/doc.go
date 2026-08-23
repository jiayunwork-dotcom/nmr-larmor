// Package httpapi exposes the NMR Larmor computation over HTTP. It registers the
// two JSON endpoints required by the API contract (/api/larmor and /api/scan),
// serves the static web console from the web/ directory, and serves the bundled
// example inputs from example/ so the page can "load example" without any
// backend-side hard-coding.
package httpapi
