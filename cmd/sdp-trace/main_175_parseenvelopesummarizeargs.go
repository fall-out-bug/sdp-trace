package main

import (
	"fmt"
	"io"
)

func parseEnvelopeSummarizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	if len(args) == 0 || args[0] != "summarize" {
		// The envelope namespace currently has one explicit verb, keeping room for
		// future envelope operations without ambiguous positional parsing.
		fmt.Fprintln(stderr, "envelope requires summarize")
		return nil, exitUsage, false
	}
	opts := &flagSet{name: "envelope summarize"}
	// Envelope summaries require a concrete envelope path; optional output is a
	// second copy of the same derived view.
	opts.setString("envelope", "")
	opts.setString("out", "")
	if err := opts.parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	code, ok := requireOnlyFlagsCode(opts, stderr, "envelope summarize accepts only flags", envelopeSummarizeRequiredFlags)
	return opts, code, ok
}
