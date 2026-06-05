package packet

// Artifact entries carry expiry and digest through to the manifest because pass
// evidence must remain replayable at validation time.
func githubArtifactEntries(input GitHubPREvidenceInput) []BundleEntry {
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

// Integration entries are manual evidence hooks, not CI-owned proof; actor and
// authority fields preserve that distinction in the manifest.
func githubIntegrationEntries(input GitHubPREvidenceInput) []BundleEntry {
	entries := []BundleEntry{}
	for _, action := range input.IntegrationActions {
		entry := bundleEntry("integration:"+action.Kind, "manual", action.Resolver, "external_ref")
		entry = authorityEntry(entry, "integration", "integration_authored", action.Actor, "github_workflow_run", input.WorkflowRunID)
		entries = append(entries, entry)
	}
	return entries
}
