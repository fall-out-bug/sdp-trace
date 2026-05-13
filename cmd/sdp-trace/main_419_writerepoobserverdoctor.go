package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func writeRepoObserverDoctor(opts *flagSet, status repoobserver.Status, stdout, stderr io.Writer) int {
	if err := repoobserver.WriteJSON(opts.stringValue("out"), status); err != nil {
		// Persisted doctor JSON is the machine-readable diagnostic artifact.
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprint(stdout, repoobserver.HumanTable(status))
	return repoObserverExitCode(status)
}
