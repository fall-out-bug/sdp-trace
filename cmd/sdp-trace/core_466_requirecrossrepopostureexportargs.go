package main

import (
	"fmt"
	"io"
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
