package harnessobs

import (
	"path/filepath"
)

func Observe(opts ObserveOptions) (Run, error) {
	// Observe keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	ctx, err := prepareObservation(opts)
	if err != nil {
		return Run{}, err
	}

	if err := writeObservationEvents(ctx.outDir, ctx.events); err != nil {
		return Run{}, err
	}
	run := newObservedRun(ctx)

	if err := writeJSON(filepath.Join(ctx.outDir, "run.json"), run); err != nil {
		return Run{}, err
	}
	return run, nil
}
