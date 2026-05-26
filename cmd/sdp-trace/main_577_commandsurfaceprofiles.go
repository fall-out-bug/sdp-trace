package main

var commandSurfaceProfiles = []profileMeta{
	{ID: "adapter-capture", Command: "assess", Description: "Assess adapter-capture evidence and overclaim risk."},
	{ID: "managed-harness", Command: "assess", Description: "Assess managed harness evidence against policy, registry, and witness inputs."},
	{ID: "forensic-retention", Command: "assess", Description: "Assess whether retained evidence supports forensic reconstruction."},
	{ID: "ci-artifact-observation", Command: "assess", Description: "Assess whether selected artifact families are CI-uploaded or lower-authority facts."},
	{ID: "authority-envelope", Command: "assess", Description: "Assess observed actions against a caller-selected authority envelope."},
}

var commandSurfaceWitnessKinds = []string{"github-actions", "gitlab-ci", "buildkite", "customer-pki"}
