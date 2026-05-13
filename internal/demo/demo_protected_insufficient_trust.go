package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func protectedInsufficientTrustScopeCondition(trustScope string) ProtectedCondition {
	// protectedInsufficientTrustScopeCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	code := "protected_trust_scope_not_satisfied"
	if trustScope == checkpoint.TrustScopeLocalSigned {

		code = "local_signed_not_protected"
	}
	return ProtectedCondition{
		ID:         "protected_trust_scope_satisfied",
		State:      GateFail,
		ReasonCode: code,
		Reason:     "observed trust scope does not satisfy protected profile",
		NextAction: "Provide CI signed checkpoint evidence with matching CI witness binding.",
	}
}
