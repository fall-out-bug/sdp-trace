package main

import (
	"context"
	"io"
)

func runQueryPack(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "explain" {
		// Explain renders an existing query-pack result; build creates one.
		return runQueryPackExplain(args[1:], stdout, stderr)
	}
	return runQueryPackBuild(args, stderr)
}
