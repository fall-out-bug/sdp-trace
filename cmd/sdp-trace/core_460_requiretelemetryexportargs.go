package main

import (
	"fmt"
	"io"
)

func requireTelemetryExportArgs(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if err := requireTelemetryExportInputs(opts); err != nil {
		// Required input checks keep unsupported profiles and missing artifacts
		// as usage errors before any metric bytes are emitted.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}
