package main

func ciArtifactPreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "artifact_manifest", "Supply artifact manifest before CI artifact observation assessment.", "Fix artifact manifest so it is readable JSON.")
}
