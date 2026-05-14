package main

import (
	"fmt"
	"io"
)

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
