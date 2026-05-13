package main

import (
	"fmt"
	"io"
)

func dispatchAssessmentExplanation(path, schemaVersion string, stdout, stderr io.Writer) int {
	handler, ok := assessmentExplainHandlers[schemaVersion]
	if !ok {
		// Unknown schemas remain cannot_verify instead of falling back to a
		// best-effort renderer that could hide a profile mismatch.
		fmt.Fprintf(stderr, "unsupported assessment-result schema_version: %s\n", schemaVersion)
		return exitCannotVerify
	}
	return handler(path, stdout, stderr)
}
