package main

import (
	"context"
	"fmt"
	"io"
)

func runQuery(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "query"}
	// Query mode is selected explicitly; the run directory is the only
	// positional evidence source.
	opts.setString("query", "")
	if err := opts.parse(args); err != nil {
		// Parser errors are usage failures before any retained run is inspected.
		return exitUsage
	}
	queryName := opts.stringValue("query")
	runDirs := opts.rest()
	if len(runDirs) == 0 {
		// Query reads one retained run directory as its evidence source.
		fmt.Fprintln(stderr, "query requires <run-dir>")
		return exitUsage
	}
	// Only the first retained run is accepted by the current query contract;
	// extra positional arguments remain outside the stable command surface.
	payload, code, ok := runNamedQuery(queryName, runDirs[0], stderr)
	if !ok {
		return code
	}
	// Query payloads are emitted as raw JSON bytes from the query package so the
	// CLI cannot alter diagnostic shape.
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}
