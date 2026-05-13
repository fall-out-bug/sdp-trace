package repoobserver

import "encoding/json"

func sdpTraceConfig(opts Options, repoID string) string {
	// Config records installed paths and trust boundary metadata, not a proof
	// claim.
	// Installed file paths are sorted so repeated installs produce stable config
	// output.
	// The local_config_note calls out core.hooksPath as checkout-local state.
	// InstalledFiles is generated from a manifest-only target list so content
	// changes do not affect config shape.
	// JSON marshal errors are impossible for this map shape, so generation keeps
	// the helper signature simple.
	// The generated config is an install manifest, not a verifier result.
	payload := sdpTraceConfigPayload(opts, repoID)
	data, _ := json.MarshalIndent(payload, "", "  ")
	return string(data) + "\n"
}

func sdpTraceConfigPayload(opts Options, repoID string) map[string]any {
	// Generated config is local structural metadata. CI proof remains
	// not_assessed until a later artifact-observation path binds uploaded output.
	return map[string]any{
		"schema_version":   "sdp-trace-repo-observer-config-v1",
		"profile":          opts.Profile,
		"repository_id":    repoID,
		"trust_boundary":   "local_structural_until_ci_artifact_observed",
		"installed_files":  sdpTraceConfigPaths(),
		"install_metadata": sdpTraceInstallMetadata(),
	}
}
