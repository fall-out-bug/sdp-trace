package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func writeProtectedGateResult(path string, result demo.GateResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(path, result); err != nil {
		// A protected verdict that cannot be written is not reviewable evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeIndentedPayload(stdout, result)
	return gateExitCode(result)
}
