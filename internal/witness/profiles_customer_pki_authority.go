package witness

import (
	"crypto/ed25519"
)

func validateCustomerPKIAuthority(record *Record, states *ProfileStates, publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) bool {
	// Authority validation rejects signer, key, policy, and revocation issues
	// before any freshness signature is allowed to establish trust.
	issue, ok := nextCustomerPKIAuthorityIssue(publicKey, policy, freshness)
	if !ok {
		return true
	}
	if issue.notAssessed {
		states.SignerAuthorityState = stateNotAssessed
		applyProfileState(record, StatusNotAssessed, stateNotAssessed, issue.reason)
		return false
	}
	customerPKIFail(record, states, issue.field, issue.reason)
	return false
}

type customerPKIAuthorityIssue struct {
	field       string
	reason      string
	notAssessed bool
	matches     bool
}

type customerPKIAuthorityCheck struct {
	field       string
	reason      string
	notAssessed bool
	matches     func() bool
}

func nextCustomerPKIAuthorityIssue(publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) (customerPKIAuthorityIssue, bool) {
	// The issue order is intentional: identity/key mismatches are stronger than
	// revocation gaps, and a required-but-absent revocation check stays
	// not_assessed instead of being reported as a false failure.
	// Each row names the profile state that should be lowered if it matches.
	// Matching stops at the first issue so reason-code precedence is stable.
	for _, check := range orderedCustomerPKIAuthorityChecks(publicKey, policy, freshness) {
		if check.matches() {
			return customerPKIAuthorityIssue{
				field:       check.field,
				reason:      check.reason,
				notAssessed: check.notAssessed,
				matches:     true,
			}, true
		}
	}
	return customerPKIAuthorityIssue{}, false
}
