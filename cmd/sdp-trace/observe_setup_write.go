package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

var writeObserveSetupPayload = func(stdout, stderr io.Writer, session harnessobs.SessionRun) int {
	if !writeJSONPayload(stdout, stderr, session, "marshal observe setup") {
		return exitCannotVerify
	}
	return 0
}
