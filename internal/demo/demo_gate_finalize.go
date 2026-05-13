package demo

import (
	"sort"
)

func finalizeGateResult(result *GateResult) {
	// finalizeGateResult keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.GateConditions = gateConditions(*result)
	if len(result.Reasons) == 0 {
		result.Reasons = append(result.Reasons, "local contract evidence is complete for the local gate")
	}
	result.Reasons = append(result.Reasons, "audit-grade release gate cannot verify without CI/OIDC witness and external witness checkpoint")
	sort.Strings(result.Reasons)
	sort.Strings(result.NextActions)
}
