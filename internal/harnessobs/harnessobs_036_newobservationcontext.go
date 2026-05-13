package harnessobs

func newObservationContext(opts ObserveOptions, outDir, sourcePath, sourceDigest string, profile Profile, events []Event) observationContext {
	// newObservationContext keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	return observationContext{
		outDir:       outDir,
		sourcePath:   sourcePath,
		sourceDigest: sourceDigest,
		now:          observationTime(opts.Now),

		profile: profile,
		events:  events,
	}
}
