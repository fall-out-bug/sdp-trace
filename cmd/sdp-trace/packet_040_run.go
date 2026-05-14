package main

import (
	"context"
	"io"
)

func runPacket(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "packet <build-pr|build-github|validate|check-demo|render> [flags]", "packet requires build-pr, build-github, validate, check-demo, or render", packetHandlers)
}
