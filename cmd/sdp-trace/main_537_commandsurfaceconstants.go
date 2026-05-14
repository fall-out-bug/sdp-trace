package main

const commandSurfaceSchemaVersion = "sdp-trace-command-surface-v1"

func knownAssessmentProfiles() []profileMeta {
	return []profileMeta{
		{ID: "adapter-capture", Command: "assess", Description: "Assess adapter-capture evidence and overclaim risk."},
		{ID: "managed-harness", Command: "assess", Description: "Assess managed harness evidence against policy, registry, and witness inputs."},
		{ID: "forensic-retention", Command: "assess", Description: "Assess whether retained evidence supports forensic reconstruction."},
		{ID: "ci-artifact-observation", Command: "assess", Description: "Assess whether selected artifact families are CI-uploaded or lower-authority facts."},
		{ID: "authority-envelope", Command: "assess", Description: "Assess observed actions against a caller-selected authority envelope."},
	}
}

func knownWitnessKinds() []string {
	return []string{"github-actions", "gitlab-ci", "buildkite", "customer-pki"}
}

func knownStates() []stateMeta {
	return []stateMeta{
		{Name: "observed", Description: "Verifier evidence met required checks for the selected local profile."},
		{Name: "pass", Description: "Selected profile concluded successfully where the command contract uses pass/fail states."},
		{Name: "fail", Description: "Verifier evidence conflicted or was insufficient for required checks."},
		{Name: "not_assessed", Description: "State was outside the run scope; it does not imply success or evidence."},
		{Name: "cannot_verify", Description: "Verifier could not execute a required check or lacked required evidence."},
	}
}
