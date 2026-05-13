package main

import (
	"context"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/recorder"
)

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
