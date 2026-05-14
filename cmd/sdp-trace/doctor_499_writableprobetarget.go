package main

import (
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
