package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

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
