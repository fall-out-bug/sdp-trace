package main

func harnessStateExitCode(state string) int {
	return stringExitCode(state, harnessStateExitCodes, exitCannotVerify)
}
