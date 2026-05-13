package packet

func githubBaseEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubBaseEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return []BundleEntry{
		authorityEntry(bundleEntry("github:pr", "change_host", input.PR.URL, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("git:commit-range", "git", input.CommitRange.Base+".."+input.CommitRange.Head, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("theater:builder", "witness", "sdp-trace packet build-pr", "raw"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("decision:owners", "manual", "default generated decision owners", "raw"), "operator", "operator_authored", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("gap:generated", "manual", "generated residual gaps", "raw"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
	}
}
