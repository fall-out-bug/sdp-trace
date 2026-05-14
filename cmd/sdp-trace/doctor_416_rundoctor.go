package main

import (
	"context"
	"io"
	"strings"
)

func runDoctor(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts, code, ok := parseDoctorArgs(args, stderr)
	if !ok {
		return code
	}
	if strings.TrimSpace(opts.stringValue("profile")) != "" {
		// Profile mode changes the evidence surface from local defaults to
		// repository installation/proof diagnostics.
		// Profile-scoped doctor delegates to repoobserver because those checks
		// inspect repository installation/proof state, not local run defaults.
		return runRepoObserverDoctor(opts, stdout, stderr)
	}
	return runLocalDoctor(opts, stdout, stderr)
}
