package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func repoObserverExitCode(status repoobserver.Status) int {
	if status.InstallState == repoobserver.StateCannotVerify || status.ProofState == repoobserver.StateCannotVerify {
		// Cannot-verify install/proof state stays distinct from failed install.
		return exitCannotVerify
	}
	if status.InstallState == repoobserver.StateFail {
		return 1
	}
	return 0
}
