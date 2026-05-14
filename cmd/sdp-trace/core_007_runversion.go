package main

import (
	"context"
	"fmt"
	"io"
)

func runVersion(_ context.Context, _ []string, stdout, _ io.Writer) int {
	fmt.Fprintf(stdout, "sdp-trace %s\n", version)
	return 0
}
