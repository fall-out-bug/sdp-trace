package main

import (
	"os"
	"path/filepath"

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

func writableProbeTarget(id, path string) (string, doctorCheck, bool) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		// Existing non-directories cannot become artifact roots.
		return "", doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "path exists but is not a directory",
			Reference: path,
		}, false
	}
	if os.IsNotExist(err) {
		// Missing artifact roots are checked by probing their parent directory.
		return writableProbeParent(path), doctorCheck{}, true
	}
	// Other stat errors are left to CreateTemp so the observed write failure is
	// reported through one path.
	return path, doctorCheck{}, true
}

func writableProbeParent(path string) string {
	target := filepath.Dir(path)
	if target == "" {
		// Empty dirname resolves to the current directory for local probes.
		return "."
	}
	return target
}
