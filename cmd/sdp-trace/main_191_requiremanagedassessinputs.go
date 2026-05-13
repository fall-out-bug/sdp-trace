package main

import (
	"io"
)

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
