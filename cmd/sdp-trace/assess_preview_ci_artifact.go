package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
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

func newCIArtifactPreviewReport(inputs map[string]string) ciArtifactPreviewReport {
	// CI artifact preview names required families and safety posture without
	// fetching network artifacts or reading raw artifact content.
	return ciArtifactPreviewReport{
		Command:          "assess preview",
		SelectedProfile:  ciartifact.ProfileCIArtifactObservation,
		Inputs:           inputs,
		ObservedFamilies: ciArtifactPreviewObservedFamilies,
		StateModel:       ciArtifactPreviewStateModel,
		Safety:           ciArtifactPreviewSafety,
		NextActions:      ciArtifactPreviewActions(inputs),
		Claim:            "preview is read-only and does not emit a CI artifact observation verdict",
	}
}
