package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func requireCrossRepoPostureExportArgs(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if err := requireCrossRepoPostureInputs(opts); err != nil {
		// The profile and selection file are mandatory even in validate-only
		// mode because they define the posture evidence boundary.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func requireCrossRepoPostureInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("profile")) != posture.ProfileID {
		// Keep the CLI bound to the only supported cross-repo posture contract.
		return fmt.Errorf("export cross-repo-posture requires --profile cross-repo-evidence-posture-v1")
	}
	if strings.TrimSpace(opts.stringValue("selection")) == "" {
		// The selection artifact is the auditable repository set for posture.
		return fmt.Errorf("export cross-repo-posture requires --selection")
	}
	return nil
}

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
