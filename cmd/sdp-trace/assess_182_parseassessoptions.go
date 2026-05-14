package main

import (
	"fmt"
	"io"
)

func parseAssessOptions(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := newStringFlagSet("assess", assessStringFlags)
	if err := opts.parse(args); err != nil {
		// Parse failures happen before profile-specific evidence loading.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	if len(opts.rest()) != 0 {
		// Assessments are entirely flag-addressed so verdict artifacts can be
		// replayed from named evidence inputs.
		fmt.Fprintln(stderr, "assess accepts only flags")
		return nil, false
	}
	return opts, true
}
