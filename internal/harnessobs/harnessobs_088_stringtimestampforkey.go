package harnessobs

import (
	"time"
)

func stringTimestampForKey(raw map[string]any, key string) string {
	// stringTimestampForKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	value := findStringByKey(raw, key)
	if value == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {

		return ""
	}
	return parsed.UTC().Format(time.RFC3339)
}
