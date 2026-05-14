package main

import (
	"io"
)

func requireAdapterCaptureAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Adapter capture has only a run source and a durable result path; both are
	// required before evaluation can produce citeable assessment JSON.
	return requireNamedValues(map[string]string{
		"--out": opts.stringValue("out"),
		"--run": opts.stringValue("run"),
	}, stderr, "adapter capture assess")
}
