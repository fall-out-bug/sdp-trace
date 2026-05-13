package main

import (
	"fmt"
	"io"
)

func gateOutputPath(opts *flagSet, stderr io.Writer) (string, bool) {
	// The output path is validated after target arity so diagnostics first
	// establish which evidence source the gate would evaluate.
	outPath := opts.stringValue("out")
	if outPath != "" {
		return outPath, true
	}
	// Persisted gate JSON is the artifact later explain/preview commands can
	// inspect; stdout is only a rendered copy.
	fmt.Fprintln(stderr, "gate requires --out <file>")
	return "", false
}
