package main

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"io"
)

func writeAssessmentArtifact[T any](path string, result T, stdout, stderr io.Writer, exitCode func(T) int) int {
	if err := writeJSONFile(path, result); err != nil {
		// Persisted JSON is the assessment authority for downstream gates.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return exitCode(result)
}

func writeCIArtifactAssessment(opts *flagSet, result ciartifact.ObservationResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		// Observation JSON is the durable artifact used by later review gates.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return ciArtifactExitCode(result)
}

func writeAuthorityAssessment(opts *flagSet, result authority.Result, stdout, stderr io.Writer) int {
	if err := authority.Write(opts.stringValue("out"), result); err != nil {
		// Authority results need the package-specific writer to preserve schema
		// shape and deterministic formatting.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return authorityExitCode(result)
}
