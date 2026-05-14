package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func handleRepoObserverInstallError(status repoobserver.Status, err error, stdout, stderr io.Writer) (int, bool) {
	if err == nil {
		return 0, false
	}
	if status.SchemaVersion != "" {
		// Partial status with a schema is still useful human evidence.
		fmt.Fprint(stdout, repoobserver.HumanTable(status))
	}
	fmt.Fprintln(stderr, err)
	return exitCannotVerify, true
}
