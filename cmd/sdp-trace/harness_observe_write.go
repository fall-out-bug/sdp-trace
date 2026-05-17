package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

var writeHarnessRunPayload = func(stdout, stderr io.Writer, run harnessobs.Run) int {
	if !writeJSONPayload(stdout, stderr, run, "marshal harness run") {
		return exitCannotVerify
	}
	return 0
}
