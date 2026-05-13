package harnessobs

import (
	"time"
)

func rawNormalizationTime(now time.Time) time.Time {
	// rawNormalizationTime keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if now.IsZero() {

		return time.Now().UTC()
	}
	return now
}
