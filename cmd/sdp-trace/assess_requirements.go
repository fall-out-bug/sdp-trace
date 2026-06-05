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

func requireManagedAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Managed assessment combines five independent evidence inputs. Keeping the
	// named flag list here preserves explicit provenance for each missing input.
	return requireNamedValues(map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--adapter-registry": opts.stringValue("adapter-registry"),
		"--managed-policy":   opts.stringValue("managed-policy"),
		"--managed-witness":  opts.stringValue("managed-witness"),
	}, stderr, "managed assess")
}

func requireForensicAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Forensic retention needs a run plus the redaction policy that defines what
	// safe retained evidence means for this assessment.
	return requireNamedValues(map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--redaction-policy": opts.stringValue("redaction-policy"),
	}, stderr, "forensic assess")
}
