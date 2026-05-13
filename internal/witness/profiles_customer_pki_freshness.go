package witness

import (
	"time"
)

func validateCustomerPKIFreshness(record *Record, states *ProfileStates, runsRoot, payloadDigest string, freshness CustomerPKIFreshnessEvidence) bool {
	if invalidFreshnessPayloadDigest(payloadDigest, freshness.PayloadDigest) {
		// Payload digest mismatch means the signer did not authorize the artifact
		// payload currently being assessed.
		customerPKIFail(record, states, "artifact", ReasonArtifactMismatch)
		return false
	}
	if !freshnessCurrent(freshness, time.Now().UTC()) {
		// Expired or future-dated freshness evidence cannot establish a live
		// external witness even if its signature is valid.
		customerPKIFail(record, states, "freshness", ReasonStaleFreshness)
		return false
	}
	if !runIDMatches(runsRoot, freshness.RunID) {
		// The signed run ID must resolve to a discovered local run so the
		// external freshness evidence binds to this evidence set.
		customerPKIFail(record, states, "run", ReasonRunMismatch)
		return false
	}
	return true
}

func invalidFreshnessPayloadDigest(expected, actual string) bool {
	// A digest must both match the expected payload and have strong hex shape;
	// malformed values are treated as binding failures.
	return expected != actual || !strongDigest(actual)
}

func customerPKIFail(record *Record, states *ProfileStates, field, reason string) {
	// The profile-level fail and the specific state fail are updated together so
	// downstream gates can explain which evidence class contradicted the claim.
	failCustomerPKIState(states, field)
	applyProfileState(record, StatusFail, stateFail, reason)
}

func failCustomerPKIState(states *ProfileStates, field string) {
	setter := customerPKIStateSetters[field]
	if setter == nil {
		// Unknown failure fields fall back to identity failure so the profile
		// cannot accidentally pass an unmapped Customer PKI condition.
		setter = func(states *ProfileStates) { states.IdentityState = stateFail }
	}
	setter(states)
}
