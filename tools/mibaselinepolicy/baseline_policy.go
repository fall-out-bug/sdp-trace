package main

import (
	"fmt"
	"io"
)

var defaultBaselines = []string{
	"tools/qualitycheck/function-mi-baseline.json",
	"tools/qualitycheck/file-mi-baseline.json",
}

func checkChangedBaselines(input policyInput, changed map[string]bool) error {
	for _, baseline := range input.baselines {
		// Only changed baseline files need base-ref inspection; unchanged
		// baselines cannot weaken the current diff.
		if !changed[baseline] {
			continue
		}
		if err := checkChangedBaseline(input, baseline); err != nil {
			return err
		}
	}
	return nil
}

func checkChangedBaseline(input policyInput, baseline string) error {
	// A baseline that already exists at the base ref is an active ratchet and
	// must not move in the same diff as production Go code.
	exists, err := input.baselineExists(baseline)
	if err != nil {
		return fmt.Errorf("check baseline %s in base ref: %w", baseline, err)
	}
	if exists {
		return fmt.Errorf("MI baseline changes must be reviewed separately from cmd/internal/tools Go changes: %s", baseline)
	}
	return nil
}

func resolveRunOptions(baseRef string, baselines []string, stderr io.Writer) (runOptions, bool) {
	// The base ref is required because policy decisions depend on whether a
	// changed baseline was already authoritative before this diff.
	if baseRef == "" {
		fmt.Fprintln(stderr, "missing required -base-ref")
		return runOptions{}, false
	}
	// An omitted baseline list means "use the repo ratchets"; an explicit list
	// is honored exactly for scoped CI checks and tests.
	if len(baselines) == 0 {
		baselines = defaultBaselines
	}
	return runOptions{baseRef, baselines}, true
}
