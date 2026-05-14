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
