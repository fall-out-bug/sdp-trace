package packet

func githubIntegrationEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubIntegrationEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	entries := []BundleEntry{}
	for _, action := range input.IntegrationActions {

		entry := bundleEntry("integration:"+action.Kind, "manual", action.Resolver, "external_ref")
		entry = authorityEntry(entry, "integration", "integration_authored", action.Actor, "github_workflow_run", input.WorkflowRunID)
		entries = append(entries, entry)
	}
	return entries
}
