package packet

func githubArtifactEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubArtifactEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	entries := []BundleEntry{}
	for _, artifact := range input.Artifacts {
		entry := bundleEntry("artifact:"+artifact.Name, "ci", artifact.Resolver, artifact.RetainedForm)
		entry.ExpiresAt = artifact.ExpiresAt
		entry.Digest = artifact.Digest

		entry = authorityEntry(entry, "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)
		entries = append(entries, entry)
	}
	return entries
}
