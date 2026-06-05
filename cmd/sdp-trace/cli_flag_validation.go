package main

import (
	"fmt"
	"io"
	"strings"
)

func rejectRest(opts *flagSet, stderr io.Writer, message string) bool {
	if len(opts.rest()) == 0 {
		return false
	}
	// Positional arguments are rejected before required flags so diagnostics
	// cannot imply that an ignored payload was accepted.
	fmt.Fprintln(stderr, message)
	return true
}

func requireStringFlag(opts *flagSet, stderr io.Writer, flag, message string) bool {
	if strings.TrimSpace(opts.stringValue(flag)) != "" {
		return true
	}
	// Empty string flags are missing evidence inputs even if the flag appeared.
	fmt.Fprintln(stderr, message)
	return false
}

type requiredCLIFlag struct {
	name    string
	message string
}

func requireOnlyFlags(opts *flagSet, stderr io.Writer, restMessage string, required []requiredCLIFlag) bool {
	if rejectRest(opts, stderr, restMessage) {
		return false
	}
	// Required flag checks run only after the command proves it received no
	// positional payload.
	return requireRequiredFlags(opts, stderr, required)
}

func requireRequiredFlags(opts *flagSet, stderr io.Writer, required []requiredCLIFlag) bool {
	for _, flag := range required {
		// Preserve caller-specific messages so command docs and tests can pin the
		// exact missing input.
		if !requireStringFlag(opts, stderr, flag.name, flag.message) {
			return false
		}
	}
	return true
}
