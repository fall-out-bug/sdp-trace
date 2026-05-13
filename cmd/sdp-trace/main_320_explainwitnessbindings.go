package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainWitnessBindings(bindings []demo.WitnessBinding, stdout io.Writer) {
	for _, binding := range bindings {
		// Binding lines expose the witness-to-run link directly instead of
		// hiding provenance under a combined health score.
		fmt.Fprintf(stdout, "Witness binding %s: %s\n", binding.ID, binding.State)
	}
}
