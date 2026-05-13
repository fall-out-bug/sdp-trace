package harnessobs

import (
	"time"
)

func unixMillisTimestamp(value float64) string {
	// unixMillisTimestamp keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if value <= 0 || value <= 1_000_000_000 {

		return ""
	}
	if value > 1_000_000_000_000 {
		return time.UnixMilli(int64(value)).UTC().Format(time.RFC3339)
	}
	return time.Unix(int64(value), 0).UTC().Format(time.RFC3339)
}
