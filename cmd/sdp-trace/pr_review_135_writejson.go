package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func writePRReviewJSON(path string, value any, stderr io.Writer) bool {
	if err := prreview.WriteJSON(path, value); err != nil {
		// Artifact write failure means the review evidence cannot be cited later.
		fmt.Fprintln(stderr, err)
		return false
	}
	return true
}
