package main

import (
	"fmt"
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
