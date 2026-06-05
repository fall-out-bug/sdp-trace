package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func previewInputCannotVerify(state string) bool {
	return state == "present_unreadable" || state == "present_malformed"
}

func managedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		return "absent"
	}
	info, err := os.Stat(path)
	if err != nil {
		// Preview status intentionally avoids leaking filesystem error details.
		return "present_unreadable"
	}
	if info.IsDir() {
		// Run-directory inputs are assessed through their normalized run.json.
		return jsonReadableStatus(filepath.Join(path, "run.json"))
	}
	return jsonReadableStatus(path)
}

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
