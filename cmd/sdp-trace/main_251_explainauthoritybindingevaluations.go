package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

func explainAuthorityBindingEvaluations(bindings []authority.EvidenceBinding, stdout io.Writer) {
	for _, binding := range bindings {
		// Binding lines keep evidence-binding state separate from action state
		// so provenance gaps remain visible in CLI output.
		fmt.Fprintf(stdout, "Binding %s: %s (%s)\n", binding.BindingID, binding.BindingState, binding.ReasonCode)
	}
}
