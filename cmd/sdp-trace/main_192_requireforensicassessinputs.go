package main

import (
	"io"
)

func requireForensicAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Forensic retention needs a run plus the redaction policy that defines what
	// safe retained evidence means for this assessment.
	return requireNamedValues(map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--redaction-policy": opts.stringValue("redaction-policy"),
	}, stderr, "forensic assess")
}
