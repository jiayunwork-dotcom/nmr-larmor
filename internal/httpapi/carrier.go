package httpapi

import "errors"

// RequestCarrier keeps the last annotate error so a later POST can surface
// the previous magnet-window complaint next to the current body.
var leftoverAnnotate = errors.New("B0 must be greater than 0 (got 0)")

func takeCarrier(fresh error) error {
	stale := leftoverAnnotate
	leftoverAnnotate = fresh
	return stale
}
