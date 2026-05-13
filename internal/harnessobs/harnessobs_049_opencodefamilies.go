package harnessobs

func openCodeFamilies(raw map[string]any, signals []string) map[string]bool {
	// openCodeFamilies keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	families := map[string]bool{}

	setFamily(families, "harness", openCodeHarnessFamily(signals))
	setFamily(families, "model", hasKey(raw, "model", "model_id", "modelid", "provider"))
	setFamily(families, "interaction", openCodeInteractionFamily(raw, signals))
	setFamily(families, "tool", openCodeToolFamily(raw, signals))
	setFamily(families, "mutation", openCodeMutationFamily(raw, signals))
	setFamily(families, "test", openCodeTestFamily(signals))
	setFamily(families, "phase", openCodePhaseFamily(raw, signals))
	setFamily(families, "review", hasSignal(signals, "review") || hasSignalPrefix(signals, "review."))
	setFamily(families, "pr", hasSignal(signals, "pull_request", "pull request") || hasSignalPrefix(signals, "pr.", "pr_"))

	setFamily(families, "merge", hasSignal(signals, "merge") || hasSignalPrefix(signals, "merge."))

	return families
}
