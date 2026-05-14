package main

var releaseProofRequiredFlags = []requiredCLIFlag{
	{"manifest", "release-proof requires --manifest"},
	{"out", "release-proof requires --out"},
}

var releaseProofExitCodes = map[string]int{
	"pass":           0,
	"fail":           exitFail,
	"cannot_verify":  exitCannotVerify,
	"not_assessed":   0,
}

func releaseProofExitCode(state string) int {
	return stringExitCode(state, releaseProofExitCodes, exitCannotVerify)
}
