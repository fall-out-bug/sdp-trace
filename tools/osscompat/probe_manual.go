package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// runCUEImport tests CUE JSON Schema import.
func runCUEImport() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	args := []string{
		"import",
		"--package", "sdptrace",
		"-o", "-",
		filepath.Join(repoRoot(), "schema/flight-recorder-run.schema.json"),
	}
	if out, err := runExternalTool(ctx, "cue", args...); err != nil {
		return stateFail, fmt.Sprintf("cue import failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "cue can import flight-recorder JSON Schema to stdout"
}
