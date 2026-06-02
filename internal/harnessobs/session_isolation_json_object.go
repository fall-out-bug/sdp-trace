package harnessobs

import (
	"encoding/json"
	"errors"
	"os"
)

// JSON object loading keeps missing settings permissive for setup while
// preserving invalid or unreadable content as explicit installation failures.

// readOptionalJSONObject treats a missing settings file as an empty object and
// leaves invalid or unreadable content as an installation error.
func readOptionalJSONObject(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return optionalJSONObjectReadError(err)
	}

	return parseOptionalJSONObject(data)
}

// optionalJSONObjectReadError separates the intentional missing-file default
// from real filesystem read failures.
func optionalJSONObjectReadError(err error) (map[string]any, error) {
	if errors.Is(err, os.ErrNotExist) {
		return map[string]any{}, nil
	}
	return nil, err
}

// parseOptionalJSONObject accepts blank settings files as empty JSON objects
// but still validates non-blank content with the JSON decoder.
func parseOptionalJSONObject(data []byte) (map[string]any, error) {
	config := map[string]any{}
	if blankJSON(data) {
		return config, nil
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}
