package main

import (
	"context"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/recorder"
)

func runLegacyWrapRecorder(ctx context.Context, opts *flagSet, command []string, stdout, stderr io.Writer) int {
	// recorder.Run owns artifact creation and event sequencing for wrapped
	// commands.
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: true,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	return writeRunResult(res, err, stdout, stderr)
}

func runTaskRecorder(ctx context.Context, opts *flagSet, command []string, stdout, stderr io.Writer) int {
	useDefault := opts.boolValue("use-default-contract")
	// The recorder is the only layer that writes run manifests and trace events.
	// CLI flags are passed through as explicit recorder options so the retained
	// manifest can explain the task, wrapper, contract, and command sources.
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		Task:               opts.stringValue("task"),
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: useDefault,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	return writeRunResult(res, err, stdout, stderr)
}

func writeRunResult(res recorder.RecorderResult, err error, stdout, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	// The run directory is the durable artifact root; the child exit code is
	// preserved for shell automation.
	fmt.Fprintf(stdout, "run_dir: %s\n", res.RunDir)
	return res.ExitCode
}
