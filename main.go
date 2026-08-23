// Command nmr-larmor is the entry point for the NMR Larmor-frequency account.
// It wires two front-ends from the same core packages:
//
//   - a web console (go run . -http :8080) that serves /api/* and the bundled
//     web/ page, and
//   - a one-shot CLI compute (go run . -example example/h1-7t.json) that prints
//     the resonance as JSON.
//
// All physics lives in internal/{nuclei,larmor,units,httpapi}; this file only
// parses flags and connects inputs to those packages.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"nmr-larmor/internal/httpapi"
	"nmr-larmor/internal/larmor"
)

func main() {
	httpAddr := flag.String("http", "", "if set, start the web console on this address, e.g. :8080")
	examplePath := flag.String("example", "", "path to an example JSON file (LarmorRequest) to compute and print")
	nucleus := flag.String("nucleus", "", "nucleus symbol for a direct compute, e.g. 1H")
	b0 := flag.Float64("B0", 0, "main magnetic field in tesla for a direct compute")
	delta := flag.Float64("delta", 0, "chemical shift in ppm for a direct compute")
	flag.Parse()

	if *httpAddr != "" {
		runServer(*httpAddr)
		return
	}

	req, err := buildRequest(*examplePath, *nucleus, *b0, *delta)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	res, err := larmor.ComputeLarmor(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

// buildRequest resolves the user's intent into a LarmorRequest. An explicit
// -example file wins; otherwise -nucleus/-B0/-delta are assembled.
func buildRequest(examplePath, nucleus string, b0, delta float64) (larmor.LarmorRequest, error) {
	if examplePath != "" {
		data, err := os.ReadFile(examplePath)
		if err != nil {
			return larmor.LarmorRequest{}, fmt.Errorf("reading example %q: %w", examplePath, err)
		}
		var req larmor.LarmorRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return larmor.LarmorRequest{}, fmt.Errorf("parsing example %q: %w", examplePath, err)
		}
		return req, nil
	}
	if nucleus == "" {
		return larmor.LarmorRequest{}, fmt.Errorf("provide -example <file> or -nucleus <symbol> -B0 <tesla>")
	}
	return larmor.LarmorRequest{
		Nucleus:   nucleus,
		B0:        b0,
		DeltaPpm:  delta,
		DeltaUnit: "ppm",
	}, nil
}

// runServer starts the HTTP console and blocks.
func runServer(addr string) {
	mux := httpapi.NewServeMux(httpapi.DefaultConfig())
	fmt.Printf("nmr-larmor console listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, "server error:", err)
		os.Exit(1)
	}
}
