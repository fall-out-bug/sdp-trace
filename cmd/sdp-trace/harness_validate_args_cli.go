package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

var documentedHarnessCannotVerifyStates = map[string]bool{
	harnessobs.StateNotAssessed:  true,
	harnessobs.StateCannotVerify: true,
}

func parseHarnessValidateArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "harness validate"}
	// Validation is bound to one profile and one retained run; out is optional
	// persistence for the validation artifact.
	opts.setString("profile", "")
	opts.setString("run", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Validation reads retained observation evidence and optionally persists a
	// validation artifact through harnessobs.
	if !requireOnlyFlags(opts, stderr, "harness validate accepts only flags", harnessValidateRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func harnessValidationExitCode(validation harnessobs.Validation, stderr io.Writer) int {
	// Exit codes mirror validation state while preserving unknown states as a
	// diagnostic cannot_verify path.
	code := harnessStateExitCode(validation.ValidationState)
	if code == exitCannotVerify && !documentedHarnessCannotVerifyStates[validation.ValidationState] {
		fmt.Fprintf(stderr, "unknown harness validation state: %s\n", validation.ValidationState)
	}
	return code
}
