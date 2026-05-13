package witness

import (
	"crypto/ed25519"
)

func applyCustomerPKIPass(record *Record) {
	// Passing Customer PKI records establish external trust only after every
	// policy, freshness, signature, and run-binding gate has returned pass.
	record.Status = StatusPass
	record.TrustScope = TrustScopeExternal
	record.EstablishedTrustScope = TrustScopeExternal
	record.Reason = ReasonProfileVerified
}

func validateCustomerPKIRecord(record *Record, states *ProfileStates, publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence, runsRoot, payloadDigest string) bool {
	// Authority, freshness, and signature are checked as separate gates so the
	// failing profile state can name the exact evidence boundary that broke.
	if !validateCustomerPKIAuthority(record, states, publicKey, policy, freshness) {
		return false
	}
	if !validateCustomerPKIFreshness(record, states, runsRoot, payloadDigest, freshness) {
		return false
	}
	if !verifyFreshnessSignature(publicKey, freshness) {
		customerPKIFail(record, states, "freshness", ReasonSignerMismatch)
		return false
	}
	return true
}
