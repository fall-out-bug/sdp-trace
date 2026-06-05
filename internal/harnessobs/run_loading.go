package harnessobs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Run loading reads a previously captured run and its referenced event files.
// The safe reference check keeps event paths bound to the run directory.
// Schema validation remains separate from JSON decoding so malformed files and
// unsupported artifacts keep their existing error boundaries.
func LoadRun(dir string) (Run, []Event, error) {
	var run Run
	if err := readJSON(filepath.Join(dir, "run.json"), &run); err != nil {
		return Run{}, nil, err
	}
	if err := validateLoadedRun(run); err != nil {
		return Run{}, nil, err
	}
	events, err := loadRunEvents(dir, run.EventRefs)
	if err != nil {
		return Run{}, nil, err
	}
	return run, events, nil
}

func readJSON(path string, target any) error {
	// This narrow reader is used by run artifacts that do not need the stricter
	// unknown-field rejection path used for profile/config loading.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func validateLoadedRun(run Run) error {
	// Run artifacts are replay inputs, so an unsupported schema keeps the whole
	// load in a normal error state instead of being coerced into validation.
	if run.SchemaVersion != RunSchemaVersion {
		return fmt.Errorf("unsupported run schema_version: %s", run.SchemaVersion)
	}
	return nil
}
