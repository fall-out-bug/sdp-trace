package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func runNamedQuery(queryName, runDir string, stderr io.Writer) ([]byte, int, bool) {
	if queryName == query.QueryCaptureDepth {
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
