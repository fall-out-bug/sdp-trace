package repoobserver

import "sort"

func sdpTraceConfigPaths() []string {
	// Installed file paths are an index of generated surfaces, not proof that
	// those surfaces have executed.
	paths := make([]string, 0)
	for _, target := range installTargetsForManifest() {
		paths = append(paths, target.path)
	}
	paths = append(paths, ".gitignore:# sdp-trace begin")
	sort.Strings(paths)
	return paths
}

func sdpTraceInstallMetadata() map[string]string {
	// Local core.hooksPath is called out explicitly because it is checkout-local
	// configuration and not committed repository evidence.
	return map[string]string{
		"generated_by":      "sdp-trace install repo-observer",
		"template_version":  SchemaVersion,
		"local_config_note": "core.hooksPath is local checkout configuration",
	}
}

func installTargetsForManifest() []targetFile {
	// Manifest targets omit content so the config stays a compact structural
	// index.
	return []targetFile{
		{path: ".sdp-trace/README.md"},
		{path: ".sdp-trace/config.json"},
		{path: ".githooks/pre-commit"},
		{path: ".githooks/post-commit"},
		{path: ".githooks/pre-push"},
		{path: ".github/workflows/sdp-trace-observe.yml"},
	}
}
