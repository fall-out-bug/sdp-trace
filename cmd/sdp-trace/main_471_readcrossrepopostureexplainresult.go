package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func readCrossRepoPostureExplainResult(path string, stderr io.Writer) (posture.ExportResult, int, bool) {
	var result posture.ExportResult
	if err := readJSONFile(path, &result); err != nil {
		// Missing or malformed result artifacts cannot be explained honestly.
		fmt.Fprintln(stderr, "result_unreadable")
		return posture.ExportResult{}, exitCannotVerify, false
	}
	if result.SchemaVersion != posture.SchemaVersion || result.ExportProfileID != posture.ProfileID {
		// Explain only accepts the current posture export schema/profile pair.
		fmt.Fprintln(stderr, "unsupported cross-repo posture export")
		return posture.ExportResult{}, exitCannotVerify, false
	}
	// Shape validation is delegated to posture.Explain, which also enforces
	// output-safety before rendering.
	return result, 0, true
}
