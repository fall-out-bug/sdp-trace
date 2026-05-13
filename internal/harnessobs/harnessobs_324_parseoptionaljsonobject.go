package harnessobs

import (
	"encoding/json"
)

func parseOptionalJSONObject(data []byte) (map[string]any, error) {
	// parseOptionalJSONObject keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	config := map[string]any{}
	if blankJSON(data) {

		return config, nil
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}
