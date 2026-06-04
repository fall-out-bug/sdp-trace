package main

import (
	"fmt"
	"io"
)

func parsePRReviewRunArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review run"}
	// Runner selection defaults to no external allowance; callers must opt in to
	// every non-default runner family they want recorded.
	// Preview mode shares the same required inputs as an executing review run.
	// That keeps dry-run planning bound to the same packet/profile evidence.
	opts.setString("packet", "")
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	opts.setString("not-assessed-reason", "")
	// Preview is the only boolean because it changes publication, not inputs.
	opts.setBool("preview", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review run accepts only flags") {
		return nil, exitUsage, false
	}
	// Required packet/profile values are checked when the runner tries to load
	// them so read errors can carry file-specific context.
	return opts, 0, true
}
