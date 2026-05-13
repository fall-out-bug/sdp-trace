package adaptercapture

func nonExecutedTestProvenanceCondition(event AdapterEvent) Condition {
	// nonExecutedTestProvenanceCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	switch event.TestProvenance {
	case "agent_reported":
		return reportedTestCondition(event, "agent_reported_test_not_executed", "agent-reported tests are claimed as executed evidence", "agent-reported test evidence is visible but non-executed")
	case "harness_observed":
		return reportedTestCondition(event, "harness_observed_test_not_executed", "harness-observed test intent is claimed as executed evidence", "harness-observed test evidence is correlation-only")
	default:
		return cannotVerify("test_provenance_not_overclaimed", "test_provenance_missing", "test provenance is missing or unverifiable", "Record ci_executed or wrapper_executed test provenance.")
	}
}

func testProvenanceExecuted(provenance string) bool {
	return provenance == "ci_executed" || provenance == "wrapper_executed"
}

func reportedTestCondition(event AdapterEvent, failCode, failReason, cannotReason string) Condition {
	// reportedTestCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if event.ExecutedEvidenceClaimed {

		return fail("test_provenance_not_overclaimed", failCode, failReason, "Bind test evidence to CI or wrapper execution.")
	}
	return cannotVerify("test_provenance_not_overclaimed", "test_execution_unverified", cannotReason, "Capture CI or wrapper-executed test evidence.")
}
