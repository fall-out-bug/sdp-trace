package harnessobs

import (
	"os"
)

func prepareObservation(opts ObserveOptions) (observationContext, error) {
	// prepareObservation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, sourcePath, outDir, err := validateObserveOptions(opts)
	if err != nil {
		return observationContext{}, err
	}

	profile, events, sourceDigest, err := loadObservationSource(profilePath, sourcePath)
	if err != nil {
		return observationContext{}, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return observationContext{}, err
	}
	return newObservationContext(opts, outDir, sourcePath, sourceDigest, profile, events), nil
}
