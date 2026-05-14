package main

import (
	"context"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func runInstall(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseInstallRepoObserverArgs(args, stdout, stderr)
	if !ok {
		return code
	}
	// Install returns status for both preview and write modes; JSON is written
	// before error handling so failed attempts are still inspectable.
	// The repoobserver package owns preview/write semantics; the CLI guarantees
	// a status artifact is attempted first.
	status, err := repoobserver.Install(repoObserverOptionsFromFlags(opts))
	if writeErr := repoobserver.WriteJSON(opts.stringValue("out"), status); writeErr != nil {
		// Install status JSON is the durable diagnostic surface for both preview
		// and write modes.
		fmt.Fprintln(stderr, writeErr)
		return 1
	}
	if code, handled := handleRepoObserverInstallError(status, err, stdout, stderr); handled {
		return code
	}
	fmt.Fprint(stdout, repoobserver.HumanTable(status))
	return repoObserverInstallExitCode(opts.boolValue("write"), status)
}
