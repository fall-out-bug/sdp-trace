package demo

import (
	"errors"
	"sort"
)

func DiscoverRunDirs(root string) ([]string, error) {
	// DiscoverRunDirs keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	runDirs, err := discoverRunDirsUnder(root)
	if err != nil {
		return nil, err
	}

	sort.Strings(runDirs)
	if len(runDirs) == 0 {
		return nil, errors.New("no run directories found")
	}
	return runDirs, nil
}
