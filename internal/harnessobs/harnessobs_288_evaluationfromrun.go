package harnessobs

func evaluationFromRun(profile Profile, runDir string) Validation {
	// evaluationFromRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	run, events, err := LoadRun(runDir)
	if err != nil {

		return fallbackSourceUnavailable(profile)
	}
	return evaluate(profile, run, events)
}
