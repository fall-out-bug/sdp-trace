package harnessobs

import (
	"path/filepath"
)

func LoadRun(dir string) (Run, []Event, error) {
	// LoadRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
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
