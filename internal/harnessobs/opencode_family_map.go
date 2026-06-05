package harnessobs

func openCodeFamilies(raw map[string]any, signals []string) map[string]bool {
	families := map[string]bool{}
	setFamily(families, "harness", openCodeHarnessFamily(signals))
	setFamily(families, "model", hasKey(raw, "model", "model_id", "modelid", "provider"))
	setFamily(families, "interaction", openCodeInteractionFamily(raw, signals))
	setFamily(families, "tool", openCodeToolFamily(raw, signals))
	setFamily(families, "mutation", openCodeMutationFamily(raw, signals))
	setFamily(families, "test", openCodeTestFamily(signals))
	setFamily(families, "phase", openCodePhaseFamily(raw, signals))
	setFamily(families, "review", openCodeReviewFamily(signals))
	setFamily(families, "pr", openCodePRFamily(signals))
	setFamily(families, "merge", openCodeMergeFamily(signals))
	return families
}
