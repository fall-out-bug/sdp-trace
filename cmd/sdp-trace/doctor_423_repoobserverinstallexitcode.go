package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func repoObserverInstallExitCode(write bool, status repoobserver.Status) int {
	if !write {
		// Preview mode reports planned changes but does not fail on an uninstalled
		// repository state.
		return 0
	}
	return repoObserverExitCode(status)
}
