package main

import (
	"os"
	"path/filepath"
	"strings"
)

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
