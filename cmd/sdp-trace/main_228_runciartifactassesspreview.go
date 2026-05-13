package main

import (
	"io"
)

func runCIArtifactAssessPreview(opts *flagSet, stdout io.Writer) int {
	// CI artifact preview never reads artifact content or calls remote systems.
	inputs := map[string]string{
		"artifact_manifest": managedInputStatus(opts.stringValue("artifact-manifest")),
	}
	writeJSONPayloadUnchecked(stdout, newCIArtifactPreviewReport(inputs))
	return previewInputExitCode(inputs)
}
