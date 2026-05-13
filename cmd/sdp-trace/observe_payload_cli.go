package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Observe payload helpers keep the combined session/run response identical
// across collect and session modes.
// The helper writes structured replay data first; exit-code mapping happens
// only after reviewers have a JSON artifact to inspect.

func writeObserveRunPayload(stdout, stderr io.Writer, session harnessobs.SessionRun, observed harnessobs.Run, message string) bool {
	// Emit session metadata beside observed run evidence for replay/debugging.
	payload := struct {
		Session harnessobs.SessionRun `json:"session"`
		Run     harnessobs.Run        `json:"run"`
	}{Session: session, Run: observed}
	return writeJSONPayload(stdout, stderr, payload, message)
}

func observeCollectExitCode(session harnessobs.SessionRun) int {
	if session.CollectionState == harnessobs.StateCannotVerify {
		// Collection can write a session payload while still failing as evidence.
		return exitCannotVerify
	}
	return 0
}
