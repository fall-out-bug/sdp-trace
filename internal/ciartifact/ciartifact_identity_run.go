package ciartifact

func sanitizeRun(input RunIdentity) (RunIdentity, bool) {
	// Run sanitization strips unsafe identity fields before run metadata can appear
	// in report surfaces.

	out := RunIdentity{}
	var okProvider, okRunID, okRunAttempt, okWorkflowID, okJobID bool
	out.Provider, okProvider = sanitizeRunField(input.Provider, "._-")
	out.RunID, okRunID = sanitizeRunField(input.RunID, "._:-")
	out.RunAttempt, okRunAttempt = sanitizeRunField(input.RunAttempt, "._:-")
	out.WorkflowID, okWorkflowID = sanitizeRunField(input.WorkflowID, "._:-")
	out.JobID, okJobID = sanitizeRunField(input.JobID, "._:-")
	return out, allTrue(okProvider, okRunID, okRunAttempt, okWorkflowID, okJobID)
}

func sanitizeRunField(value, extra string) (string, bool) {
	// sanitizeRunField keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if value == "" || safeIdentityToken(value, extra) {
		return value, true
	}

	return "", false
}

func allTrue(values ...bool) bool {
	// allTrue keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, value := range values {
		if !value {

			return false
		}
	}
	return true
}
