package main

import (
	"context"
	"fmt"
	"io"
)

func runOverride(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseOverrideRequestArgs(args, stderr)
	if !ok {
		return code
	}
	// Appending is the only state-changing step; parsing alone never creates an
	// override artifact.
	event, err := appendOverrideRequestEvent(opts)
	if err != nil {
		// Append failure means no override request was durably recorded.
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "override_event: %s\n", event.EventID)
	return 0
}
