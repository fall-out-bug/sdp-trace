package main

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"io"
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

func runAuthorityAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireNamedValues(map[string]string{
		"--out":               opts.stringValue("out"),
		"--authority-package": opts.stringValue("authority-package"),
	}, stderr, "authority-envelope assess") {
		// Authority assessment requires an explicit package because raw prompts
		// or model outputs are outside this CLI's accepted evidence shape.
		return exitUsage
	}
	// Authority evaluation is package-bound; unreadable packages are
	// cannot_verify because the caller supplied an authority artifact path.
	pkg, err := authority.ReadPackage(opts.stringValue("authority-package"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	result := authority.Evaluate(pkg)
	return writeAuthorityAssessment(opts, result, stdout, stderr)
}
