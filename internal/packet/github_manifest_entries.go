package packet

// GitHub manifest entries are emitted in the same order as packet evidence is
// consumed: base source refs first, then optional route/check/review/artifact
// evidence. Stable order keeps generated packet diffs reviewable.
func githubEntries(input GitHubPREvidenceInput) []BundleEntry {
	entries := githubBaseEntries(input)

	entries = append(entries, githubPromptBoundaryEntries(input)...)
	entries = append(entries, githubPRBodyEntries(input)...)
	entries = append(entries, githubAgentRouteEntries(input)...)
	entries = append(entries, githubCheckEntries(input)...)
	entries = append(entries, githubReviewEntries(input)...)
	entries = append(entries, githubArtifactEntries(input)...)
	entries = append(entries, githubIntegrationEntries(input)...)
	return entries
}

// Base entries are always present so non-pass default rows still cite retained
// generated evidence instead of becoming prose-only claims.
func githubBaseEntries(input GitHubPREvidenceInput) []BundleEntry {
	return []BundleEntry{
		authorityEntry(bundleEntry("github:pr", "change_host", input.PR.URL, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("git:commit-range", "git", input.CommitRange.Base+".."+input.CommitRange.Head, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("theater:builder", "witness", "sdp-trace packet build-pr", "raw"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("decision:owners", "manual", "default generated decision owners", "raw"), "operator", "operator_authored", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("gap:generated", "manual", "generated residual gaps", "raw"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
	}
}
