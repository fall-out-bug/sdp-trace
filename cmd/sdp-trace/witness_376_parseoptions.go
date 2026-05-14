package main

import (
	"fmt"
	"io"
)

func parseWitnessOptions(args []string, stderr io.Writer) (witnessOptions, bool) {
	opts, ok := parseWitnessFlagSet(args, stderr)
	if !ok {
		return witnessOptions{}, false
	}
	// Validation is separated from flag parsing so missing required values can
	// return product-specific messages.
	options, message, ok := witnessOptionsFromFlags(opts)
	if !ok {
		fmt.Fprintln(stderr, message)
		return witnessOptions{}, false
	}
	return options, true
}
