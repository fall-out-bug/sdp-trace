package posture

import (
	"encoding/json"
	"os"
)

func readJSONFile[T any](path string) (T, error) {
	// JSON replay inputs must parse structurally before profile-specific checks.
	var result T
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}
	return result, json.Unmarshal(data, &result)
}
