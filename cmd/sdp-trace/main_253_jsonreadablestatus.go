package main

import (
	"encoding/json"
	"os"
)

func jsonReadableStatus(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		// Readability is enough for preview; the real assessment reports the
		// concrete parse/load error.
		return "present_unreadable"
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		// Malformed JSON blocks setup without interpreting partial contents.
		return "present_malformed"
	}
	return "present_readable"
}
