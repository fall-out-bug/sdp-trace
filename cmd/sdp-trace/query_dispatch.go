package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/capturedepth"
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func runNamedQuery(queryName, runDir string, stderr io.Writer) ([]byte, int, bool) {
	if queryName == capturedepth.QueryName {
		// Capture-depth is a read-only diagnostic query over retained evidence.
		return captureDepthQuery(runDir, stderr)
	}
	if queryName != query.QueryMissingEvidence {
		// Unsupported query names are usage errors, not empty findings.
		fmt.Fprintf(stderr, "unsupported query: %s\n", queryName)
		return nil, exitUsage, false
	}
	return missingEvidenceQuery(runDir, stderr)
}

func captureDepthQuery(runDir string, stderr io.Writer) ([]byte, int, bool) {
	payload, err := capturedepth.CaptureDepth(runDir)
	if err != nil {
		// Query load/replay failures mean the retained evidence cannot be
		// verified for this diagnostic.
		fmt.Fprintln(stderr, err)
		return nil, exitCannotVerify, false
	}
	return payload, 0, true
}

func missingEvidenceQuery(runDir string, stderr io.Writer) ([]byte, int, bool) {
	payload, err := query.MissingEvidence(runDir)
	if err != nil {
		// Missing-evidence query failures are cannot_verify for the query
		// result, not an empty missing-evidence list.
		fmt.Fprintln(stderr, err)
		return nil, exitCannotVerify, false
	}
	return payload, 0, true
}
