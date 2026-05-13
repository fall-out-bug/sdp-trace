package main

func releaseProofExitCode(state string) int {
	return stringExitCode(state, releaseProofExitCodes, exitCannotVerify)
}
