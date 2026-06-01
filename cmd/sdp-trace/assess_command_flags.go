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

func newStringFlagSet(name string, flags []string) *flagSet {
	opts := &flagSet{name: name}
	// Shared assess flags keep the command surface stable while each selected
	// profile validates only the inputs it can actually use.
	for _, flag := range flags {
		opts.setString(flag, "")
	}
	return opts
}

var assessStringFlags = []string{
	"profile",
	"out",
	"contract",
	"run",
	"adapter-registry",
	"managed-policy",
	"managed-witness",
	"redaction-policy",
	"artifact-manifest",
	"authority-package",
}
