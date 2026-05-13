package main

import (
	"context"
	"io"
)

func runPreview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runPreviewCommand("preview", "preview", args, stdout, stderr)
}
