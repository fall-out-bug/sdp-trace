package harnessobs

func RunSession(opts SessionOptions) (SessionRun, Run, error) {
	// RunSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session, err := setupRunnableSession(opts)
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	commandResult, err := runObservedCommand(opts.Command, &session)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	return collectFinishedSession(opts, session, commandResult.waitErr, commandResult.end)
}
