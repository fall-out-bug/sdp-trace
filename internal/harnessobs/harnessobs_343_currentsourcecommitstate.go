package harnessobs

func currentSourceCommitState() (string, string) {
	// currentSourceCommitState keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	commit := sourceCommit()
	if commit == "" {

		return "", StateCannotVerify
	}
	return commit, StatePass
}
