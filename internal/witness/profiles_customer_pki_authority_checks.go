package witness

import (
	"crypto/ed25519"
)

func orderedCustomerPKIAuthorityChecks(publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) []customerPKIAuthorityCheck {
	// Signer identity fails before key digest so callers see the authority
	// allow-list problem first.
	return []customerPKIAuthorityCheck{
		{field: "signer", reason: ReasonSignerMismatch, matches: func() bool { return customerPKISignerMismatch(policy, freshness) }},
		{field: "signer", reason: ReasonSignerMismatch, matches: func() bool { return customerPKIPublicKeyMismatch(publicKey, policy) }},
		{field: "policy", reason: ReasonPolicyMismatch, matches: func() bool { return customerPKIPolicyDigestMismatch(policy, freshness) }},
		{field: "signer", reason: ReasonRevocationNA, notAssessed: true, matches: func() bool { return customerPKIRevocationAssessmentRequired(policy) }},
		{field: "signer", reason: ReasonCertRevoked, matches: func() bool { return customerPKIRevoked(policy) }},
	}
}

func customerPKISignerMismatch(policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) bool {
	if policy.ProfileID != "customer-pki-v1" {
		// Profile ID binds the policy to this verifier contract; a policy for a
		// different profile cannot authorize customer-pki-v1 evidence.
		return true
	}
	if policy.AllowedSignerID == "" {
		// A blank signer allow-list does not grant universal signer authority.
		return true
	}
	return policy.AllowedSignerID != freshness.SignerID
}

func customerPKIPublicKeyMismatch(publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy) bool {
	// Empty key digests are allowed for early integrations, but a provided
	// digest must bind exactly to the parsed public key.
	return policy.PublicKeySHA256 != "" && policy.PublicKeySHA256 != digestBytes(publicKey)
}

func customerPKIPolicyDigestMismatch(policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) bool {
	// Policy digests are optional compatibility evidence; when present, they
	// prevent freshness evidence from silently switching authority policy.
	return policy.PolicyDigest != "" && policy.PolicyDigest != freshness.PolicyDigest
}

func customerPKIRevocationAssessmentRequired(policy CustomerPKIAuthorityPolicy) bool {
	// Required revocation evidence without a state is an assessment gap. It must
	// block pass without inventing a revoked verdict.
	return policy.RevocationRequired && policy.RevocationState == ""
}

func customerPKIRevoked(policy CustomerPKIAuthorityPolicy) bool {
	// An explicit revoked state is contradictory external authority evidence and
	// therefore fails signer authority.
	return policy.RevocationState == "revoked"
}
