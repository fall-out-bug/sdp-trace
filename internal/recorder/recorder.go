package recorder

import (
	"context"
	"errors"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

const defaultBaseOutputDir = ".sdp-trace-runs"

// The recorder package is the portable local capture boundary. It records a
// command run, source snapshot, event chain, and manifest without depending on
// any specific harness runtime.
//
// Public types remain in this file while setup, event writing, command waiting,
// and artifact finalization live in narrower files. That keeps the API visible
// without hiding the evidence lifecycle behind one large implementation file.

// RecorderOptions configures local capture behavior.
type RecorderOptions struct {
	Task               string
	ContractPath       string
	UseDefaultContract bool
	WrapperName        string
	OutputDir          string
	Command            []string
	Env                []string
}

// RecorderResult represents run output for CLI callers and tests.
type RecorderResult struct {
	RunDir   string
	ExitCode int
	Contract trace.Contract
}

// Run executes a command and writes a local chain with minimal milestone fields.
func Run(ctx context.Context, options RecorderOptions) (RecorderResult, error) {
	// The public entrypoint keeps setup, command execution, and closure as
	// separate phases so every non-zero outcome still produces terminal evidence.
	if len(options.Command) == 0 {
		return RecorderResult{}, errors.New("missing command")
	}

	prepared, err := prepareRun(options)
	if err != nil {
		return RecorderResult{}, err
	}

	exitCode, signal := runCommand(ctx, options.Command, prepared.options.Env, prepared.writer)
	if err := finishRun(prepared.writer, prepared.commandID, exitCode, signal, prepared.options); err != nil {
		return RecorderResult{}, err
	}
	return recorderResult(prepared, exitCode), nil
}

func recorderResult(prepared preparedRun, exitCode int) RecorderResult {
	// Results expose only the stable run directory, command exit code, and
	// resolved contract; event details stay in the run artifacts.
	return RecorderResult{
		RunDir:   prepared.runDir,
		ExitCode: exitCode,
		Contract: prepared.contract,
	}
}

type preparedRun struct {
	options   RecorderOptions
	runDir    string
	contract  trace.Contract
	writer    *runWriter
	commandID string
}
