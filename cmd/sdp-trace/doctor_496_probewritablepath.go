package main

import (
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func probeWritablePath(id, path, target, okReason string) doctorCheck {
	// Missing directories are probed through their parent, while existing
	// directories are probed directly.
	// Probe with a temporary file so doctor validates actual write capability,
	// not just path syntax.
	probe, err := os.CreateTemp(target, ".sdp-trace-doctor-")
	if err != nil {
		// A failed probe is recorded as cannot_verify rather than inferred from
		// permissions text or path shape.
		return doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "directory is not writable",
			Reference: path,
		}
	}
	probeName := probe.Name()
	// Probe cleanup failures are intentionally ignored; the check is about
	// whether a report/run artifact could be written.
	_ = probe.Close()
	_ = os.Remove(probeName)
	return writablePathPassCheck(id, path, okReason)
}
