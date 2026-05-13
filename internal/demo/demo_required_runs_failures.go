package demo

import (
	"fmt"
)

func cannotVerifyRequiredRunEvidence(result RequiredRunResult, requiredID, evidenceID string) RequiredRunResult {
	// cannotVerifyRequiredRunEvidence keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.State = GateCannotVerify
	result.Reasons = []string{fmt.Sprintf("required run %s missing evidence %s", requiredID, evidenceID)}
	return result
}

func cannotVerifyRequiredRun(result RequiredRunResult, requiredID, runName string) RequiredRunResult {
	// cannotVerifyRequiredRun keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.State = GateCannotVerify
	result.Reasons = []string{fmt.Sprintf("required run %s cannot verify from run %s", requiredID, runName)}
	return result
}

func applyProtectedFutureConstraint(result RequiredRunResult, requiredID string) RequiredRunResult {
	// applyProtectedFutureConstraint keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if result.Profile != GateModeProtectedFuture {
		return result
	}

	return cannotVerifyRequiredRunReason(result, requiredID, "requests protected_future profile, which cannot verify before signed checkpoint evidence exists")
}

func cannotVerifyRequiredRunReason(result RequiredRunResult, requiredID, reason string) RequiredRunResult {
	// cannotVerifyRequiredRunReason keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.State = GateCannotVerify
	result.Reasons = []string{fmt.Sprintf("required run %s %s", requiredID, reason)}
	return result
}
