package main

import (
	"fmt"
	"io"
)

func explainAssessmentResult(path string, stdout, stderr io.Writer) int {
	var envelope struct {
		SchemaVersion   string `json:"schema_version"`
		SelectedProfile string `json:"selected_profile"`
	}
	// Read the minimal schema envelope first; selected_profile is descriptive
	// and must not choose the parser.
	if err := readJSONFile(path, &envelope); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return dispatchAssessmentExplanation(path, envelope.SchemaVersion, stdout, stderr)
}

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

func explainTypedAssessment[T any](explain func(T, io.Writer) int) assessmentExplainHandler {
	return func(path string, stdout, stderr io.Writer) int {
		var result T
		// The typed load is the trust boundary for explanation; renderers only
		// restate fields from that decoded artifact.
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explain(result, stdout)
	}
}
