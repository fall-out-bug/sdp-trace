package packet

func githubPRBodyEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubPRBodyEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if input.PR.BodyRef != "" {

		return []BundleEntry{authorityEntry(bundleEntry("github:pr-body", "change_host", input.PR.BodyRef, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)}
	}
	return nil
}
