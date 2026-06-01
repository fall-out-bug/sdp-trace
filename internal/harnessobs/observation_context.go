package harnessobs

func newObservationContext(opts ObserveOptions, outDir, sourcePath, sourceDigest string, profile Profile, events []Event) observationContext {
	return observationContext{
		outDir:       outDir,
		sourcePath:   sourcePath,
		sourceDigest: sourceDigest,
		now:          observationTime(opts.Now),

		profile: profile,
		events:  events,
	}
}
