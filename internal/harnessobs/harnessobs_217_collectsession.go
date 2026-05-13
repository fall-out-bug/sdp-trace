package harnessobs

import (
	"errors"
)

func CollectSession(opts SessionCollectOptions) (SessionRun, Run, error) {
	// CollectSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	ctx, err := prepareSessionCollection(opts)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	sourcePath, err := resolveSessionEventSource(&ctx)
	if err != nil {
		if !errors.Is(err, errSessionSourceUnavailable) {
			return SessionRun{}, Run{}, err
		}
		return markSessionSourceUnavailable(ctx)
	}
	return collectSessionSource(ctx, sourcePath)
}
