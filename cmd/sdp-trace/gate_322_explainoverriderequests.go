package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainOverrideRequests(overrides []demo.OverrideRequest, stdout io.Writer) {
	for _, override := range overrides {
		// Override requests remain separate records because each one needs its
		// own evidence-backed state.
		fmt.Fprintf(stdout, "Override %s: %s\n", override.OverrideID, override.State)
	}
}
