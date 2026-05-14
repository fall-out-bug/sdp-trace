package main

import (
	"fmt"
	"io"
	"strings"
)

func parseCrossRepoPostureExplainArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export cross-repo-posture-explain"}
	// The result flag points at the posture artifact that will be explained.
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "export cross-repo-posture-explain accepts only flags") {
		return nil, exitUsage, false
	}
	if strings.TrimSpace(opts.stringValue("result")) == "" {
		// A persisted export result is required before human explanation.
		fmt.Fprintln(stderr, "export cross-repo-posture-explain requires --result")
		return nil, exitUsage, false
	}
	// Successful parsing only identifies the artifact; schema and output-safety
	// checks happen when the artifact is read and rendered.
	return opts, 0, true
}
