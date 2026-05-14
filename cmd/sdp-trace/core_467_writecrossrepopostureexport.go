package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func writeCrossRepoPostureExport(opts *flagSet, result posture.ExportResult, stderr io.Writer) int {
	if opts.boolValue("validate-only") {
		// Validate-only proves the selection can build without publishing a new
		// posture artifact.
		return 0
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		// Non-preview exports must name the durable posture artifact path.
		fmt.Fprintln(stderr, "export cross-repo-posture requires --out")
		return exitUsage
	}
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		// Failed publication leaves no reviewable posture export.
		fmt.Fprintln(stderr, "out_unwritable")
		return 1
	}
	return 0
}
