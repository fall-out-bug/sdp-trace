package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func emptyWritablePathCheck(id string) doctorCheck {
	// Empty configured paths cannot support durable run/report artifacts.
	return doctorCheck{
		ID:     id,
		State:  string(trace.VerdictCannotVerify),
		Reason: "path is empty",
	}
}

func writablePathPassCheck(id, path, okReason string) doctorCheck {
	// The temporary probe is removed immediately; doctor reports capability,
	// not an artifact to retain.
	return doctorCheck{
		ID:        id,
		State:     "pass",
		Reason:    okReason,
		Reference: path,
	}
}
