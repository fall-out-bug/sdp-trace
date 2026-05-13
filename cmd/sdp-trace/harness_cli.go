package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

var harnessHandlers = map[string]subcommandHandler{
	"observe":   runHarnessObserve,
	"validate":  runHarnessValidate,
	"summarize": runHarnessSummarize,
}

var harnessObserveRequiredFlags = []requiredCLIFlag{
	{"profile", "harness observe requires --profile"},
	{"source", "harness observe requires --source"},
	{"out", "harness observe requires --out"},
}

var harnessValidateRequiredFlags = []requiredCLIFlag{
	{"profile", "harness validate requires --profile"},
	{"run", "harness validate requires --run"},
}

var harnessSummarizeRequiredFlags = []requiredCLIFlag{
	{"validation", "harness summarize requires --validation"},
}

var documentedHarnessCannotVerifyStates = map[string]bool{
	harnessobs.StateNotAssessed:  true,
	harnessobs.StateCannotVerify: true,
}

func runHarness(args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "harness <observe|validate|summarize> [flags]", "harness requires observe, validate, or summarize", harnessHandlers)
}

func runHarnessObserve(args []string, stdout, stderr io.Writer) int {
	// Parse-time validation keeps every evidence input named as a flag.
	opts, code, ok := parseHarnessObserveArgs(args, stderr)
	if !ok {
		return code
	}
	// Observation work stays in harnessobs; the CLI only reports its artifact.
	run, err := observeHarnessRun(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeHarnessRun(stdout, stderr, run)
}

func observeHarnessRun(opts *flagSet) (harnessobs.Run, error) {
	// The CLI preserves the observation boundary: source parsing,
	// normalization, and verdict derivation stay package-owned.
	return harnessobs.Observe(harnessobs.ObserveOptions{
		ProfilePath: opts.stringValue("profile"),
		SourcePath:  opts.stringValue("source"),
		OutDir:      opts.stringValue("out"),
	})
}

func writeHarnessRun(stdout, stderr io.Writer, run harnessobs.Run) int {
	// A marshal failure means the command cannot publish replayable evidence.
	if !writeJSONPayload(stdout, stderr, run, "marshal harness run") {
		return exitCannotVerify
	}
	return 0
}

func parseHarnessObserveArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "harness observe"}
	// Observation requires the profile, raw source, and output run directory to
	// stay named for replayable artifact provenance.
	opts.setString("profile", "")
	opts.setString("source", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Harness observation is flag-only so profile/source/output are auditable.
	if !requireOnlyFlags(opts, stderr, "harness observe accepts only flags", harnessObserveRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

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

func runHarnessSummarize(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "harness summarize"}
	// Summaries are read-only views over a validation artifact.
	opts.setString("validation", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Summarize accepts only a persisted validation artifact, not raw events.
	if !requireOnlyFlags(opts, stderr, "harness summarize accepts only flags", harnessSummarizeRequiredFlags) {
		return exitUsage
	}
	validation, err := harnessobs.LoadValidation(opts.stringValue("validation"))
	if err != nil {
		// Unreadable validation artifacts keep summary in cannot_verify because
		// there is no trusted row set to summarize.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Summaries are derived from persisted validation evidence; missing or
	// malformed validation stays cannot_verify instead of becoming prose truth.
	fmt.Fprint(stdout, harnessobs.Summarize(validation))
	return 0
}
