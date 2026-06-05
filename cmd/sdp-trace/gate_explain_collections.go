package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainGateCollections(result demo.GateResult, stdout io.Writer) {
	// Collection explainers preserve the original evidence categories and
	// remediation lists from the persisted gate result.
	explainRequiredRuns(result.RequiredRuns, stdout)
	explainWitnessBindings(result.WitnessBindings, stdout)
	explainMissingAuditEvidence(result.MissingAuditEvidence, stdout)
	explainOverrideRequests(result.OverrideRequests, stdout)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
}

func explainRequiredRuns(requiredRuns []demo.RequiredRunResult, stdout io.Writer) {
	for _, requiredRun := range requiredRuns {
		// One stable line per required run keeps the human explanation auditable
		// without inventing a separate summary verdict.
		fmt.Fprintf(stdout, "Required run %s: %s\n", requiredRun.ID, requiredRun.State)
	}
}

func explainWitnessBindings(bindings []demo.WitnessBinding, stdout io.Writer) {
	for _, binding := range bindings {
		// Binding lines expose the witness-to-run link directly instead of
		// hiding provenance under a combined health score.
		fmt.Fprintf(stdout, "Witness binding %s: %s\n", binding.ID, binding.State)
	}
}

func explainMissingAuditEvidence(missingEvidence []string, stdout io.Writer) {
	for _, missing := range missingEvidence {
		// Missing audit evidence stays visible as a concrete gap; explanation
		// output must not collapse it into a green summary.
		fmt.Fprintf(stdout, "Missing audit evidence: %s\n", missing)
	}
}

func explainOverrideRequests(overrides []demo.OverrideRequest, stdout io.Writer) {
	for _, override := range overrides {
		// Override requests remain separate records because each one needs its
		// own evidence-backed state.
		fmt.Fprintf(stdout, "Override %s: %s\n", override.OverrideID, override.State)
	}
}
