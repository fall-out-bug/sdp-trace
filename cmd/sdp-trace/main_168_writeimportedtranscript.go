package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

func writeImportedTranscript(trace interaction.Trace, err error, stdout, stderr io.Writer) int {
	if err != nil {
		// Import failures mean the transcript source cannot be trusted as trace
		// evidence.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writeJSONPayloadUnchecked(stdout, trace)
	return 0
}
