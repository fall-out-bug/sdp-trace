package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

func writeCIArtifactAssessment(opts *flagSet, result ciartifact.ObservationResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		// Observation JSON is the durable artifact used by later review gates.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return ciArtifactExitCode(result)
}
