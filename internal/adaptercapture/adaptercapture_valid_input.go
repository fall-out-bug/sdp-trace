package adaptercapture

// ValidTestInput is used by CLI and fixture tests to avoid duplicating a large
// representative adapter capture run outside this package.
func ValidTestInput() Input {
	return validInput()
}

func validInput() Input {
	// validInput preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	runID := "adapter-run-1"
	nonce := "nonce-1"
	source := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	policy := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	return Input{Run: validRunEvidence(runID, nonce, source, policy)}
}

func validRunEvidence(runID, nonce, source, policy string) RunEvidence {
	// validRunEvidence preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	run := validRunHeader(runID, nonce, source, policy)
	run.AdapterEvents = validAdapterEvents(runID, nonce, source, policy)
	run.ProviderRefs = validProviderRefs(source)
	run.EventFamilySummaries = validEventFamilySummaries()
	return run
}

func validRunHeader(runID, nonce, source, policy string) RunEvidence {
	// validRunHeader preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	return RunEvidence{
		RunID:                 runID,
		RunNonce:              nonce,
		SourceBaseline:        source,
		RunClosedSequence:     20,
		RequiredEventTypes:    validRequiredEventTypes(),
		RedactionPolicyDigest: policy,
		GatewayIntegrated:     true,
		GatewayEvidenceBound:  true,
		TaskDriftAssessed:     true,
	}
}

func validRequiredEventTypes() []string {
	return []string{"run_started", "task_locked", "tool_call", "command_started", "file_mutation", "model_call_observed", "test_observed", "run_closed"}
}

func validAdapterEvents(runID, nonce, source, policy string) []AdapterEvent {
	// validAdapterEvents preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	out := make([]AdapterEvent, 0, len(validEventSpecs))
	for _, spec := range validEventSpecs {

		out = append(out, validEvent(spec.id, spec.eventType, spec.sequence, runID, nonce, source, policy))
	}
	return out
}
