package demo

import (
	"encoding/json"
	"os"
)

func loadWitnessSummary(path string) (WitnessSummary, error) {
	// Witness files are parsed into the portable demo summary shape before any
	// binding decision is made.
	var record WitnessSummary
	data, err := os.ReadFile(path)
	if err != nil {
		return WitnessSummary{}, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return WitnessSummary{}, err
	}
	return record, nil
}
