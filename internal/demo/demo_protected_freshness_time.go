package demo

import (
	"time"
)

func protectedWitnessFreshnessAt(generatedAtText string, now time.Time) ProtectedCondition {
	// protectedWitnessFreshnessAt keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	generatedAt, err := time.Parse(time.RFC3339, generatedAtText)
	if err != nil {

		return witnessFreshnessCannotVerify("invalid_witness_freshness", "CI witness generated_at cannot be parsed", "Regenerate CI witness evidence with an RFC3339 generated_at timestamp.")
	}
	if now.IsZero() {

		now = time.Now().UTC()
	}
	if condition, ok := invalidWitnessFreshnessCondition(generatedAt, now); ok {
		return condition
	}
	return ProtectedCondition{
		ID:         "witness_freshness_valid",
		State:      GatePass,
		ReasonCode: "witness_fresh",
		Reason:     "CI witness freshness is within the protected profile window",
	}
}

func invalidWitnessFreshnessCondition(generatedAt, now time.Time) (ProtectedCondition, bool) {
	// invalidWitnessFreshnessCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if generatedAt.After(now.Add(5 * time.Minute)) {
		return witnessFreshnessFail("witness_from_future", "CI witness generated_at is after the verifier time window", "Regenerate CI witness evidence in the selected CI run."), true
	}
	if now.Sub(generatedAt) > 24*time.Hour {
		return witnessFreshnessFail("stale_witness", "CI witness generated_at is outside the protected freshness window", "Regenerate CI witness evidence for the selected run."), true
	}
	return ProtectedCondition{}, false
}
