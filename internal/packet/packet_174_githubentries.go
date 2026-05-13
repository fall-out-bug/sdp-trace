package packet

func githubEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
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
