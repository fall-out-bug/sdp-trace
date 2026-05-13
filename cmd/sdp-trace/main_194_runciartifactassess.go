package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

func runCIArtifactAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireNamedValues(map[string]string{
		"--out":               opts.stringValue("out"),
		"--artifact-manifest": opts.stringValue("artifact-manifest"),
	}, stderr, "ci artifact observation assess") {
		// CI artifact assessment starts from a named manifest, not repository
		// discovery.
		return exitUsage
	}
	var manifest ciartifact.Manifest
	// CI artifact observation starts from a manifest snapshot; unreadable
	// manifests are input failures, not not_assessed observations.
	if err := readJSONFile(opts.stringValue("artifact-manifest"), &manifest); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := ciartifact.Evaluate(manifest)
	return writeCIArtifactAssessment(opts, result, stdout, stderr)
}
