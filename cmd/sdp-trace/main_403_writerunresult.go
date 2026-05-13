package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/recorder"
)

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
