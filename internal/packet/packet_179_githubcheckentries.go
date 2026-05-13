package packet

func githubCheckEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubCheckEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if len(input.Checks) > 0 {

		return []BundleEntry{authorityEntry(bundleEntry("github:check", "ci", checkResolvers(input.Checks), "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)}
	}
	return nil
}
