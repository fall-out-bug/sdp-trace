package main

import (
	"context"
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
