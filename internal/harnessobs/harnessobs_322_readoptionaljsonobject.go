package harnessobs

import (
	"os"
)

func readOptionalJSONObject(path string) (map[string]any, error) {
	// readOptionalJSONObject keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	data, err := os.ReadFile(path)
	if err != nil {
		return optionalJSONObjectReadError(err)
	}

	return parseOptionalJSONObject(data)
}
