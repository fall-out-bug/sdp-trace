package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

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
