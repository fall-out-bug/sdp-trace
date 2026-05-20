package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/capturedepth"
)

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
