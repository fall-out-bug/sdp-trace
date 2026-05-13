package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Harness validate consumes a retained observation run and emits validation
// rows before mapping state to process exit code. The CLI does not recompute
// harness facts; it delegates that to harnessobs and only preserves artifact
// shape, required inputs, and diagnostic exit behavior.

func runHarnessValidate(args []string, stdout, stderr io.Writer) int {
	// Parse-time validation keeps profile and retained run inputs explicit.
	opts, code, ok := parseHarnessValidateArgs(args, stderr)
	if !ok {
		return code
	}
	// Validation evidence is produced before any CLI exit-code mapping.
	validation, err := validateHarnessRun(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeHarnessValidation(stdout, stderr, validation)
}

func validateHarnessRun(opts *flagSet) (harnessobs.Validation, error) {
	// Validation rows are package-owned; the CLI only supplies the profile,
	// retained run, and optional artifact path.
	return harnessobs.Validate(harnessobs.ValidateOptions{
		ProfilePath: opts.stringValue("profile"),
		RunDir:      opts.stringValue("run"),
		OutPath:     opts.stringValue("out"),
	})
}

func writeHarnessValidation(stdout, stderr io.Writer, validation harnessobs.Validation) int {
	// Emit validation rows before mapping their state to the process exit code.
	if !writeJSONPayload(stdout, stderr, validation, "marshal harness validation") {
		return exitCannotVerify
	}
	return harnessValidationExitCode(validation, stderr)
}
