package main

import (
	"encoding/json"
	"os"
)

func readJSONFile(path string, dst any) error {
	// Shared JSON reads are local artifact loads; callers decide whether a
	// failure is usage, cannot_verify, or ordinary command failure.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
