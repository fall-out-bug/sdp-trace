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

func parseWitnessFlagSet(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := newWitnessFlagSet()
	if err := opts.parse(args); err != nil {
		// Malformed witness flags stop before any CI or Customer PKI material is
		// read from disk.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	// The witness command has no positional target; target comes from flags so
	// generated records have explicit provenance fields.
	return opts, true
}

func witnessOptionsFromFlags(opts *flagSet) (witnessOptions, string, bool) {
	// Required fields are normalized before optional witness-specific material
	// is copied into the final options struct.
	fields, message, ok := witnessRequiredFieldsFromFlags(opts)
	if !ok {
		return witnessOptions{}, message, false
	}
	return witnessOptionsFromRequiredFields(fields, opts), "", true
}
