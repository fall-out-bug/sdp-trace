package main

import (
	"fmt"
	"io"
)

func writeTelemetryExportOutput(outPath, rendered string, stdout io.Writer) error {
	if outPath == "-" {
		// Dash keeps review-friendly stdout output without changing the rendered
		// payload.
		fmt.Fprint(stdout, rendered)
		return nil
	}
	if err := writeTextFileAtomic(outPath, rendered); err != nil {
		// File output is all-or-nothing; partial metric files are not accepted as
		// evidence.
		return fmt.Errorf("out_unwritable")
	}
	return nil
}
