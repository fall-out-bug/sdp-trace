package witness

import (
	"crypto/ed25519"
)

func loadCustomerPKIInputs(record *Record, opts ProfileOptions) (ed25519.PublicKey, CustomerPKIAuthorityPolicy, CustomerPKIFreshnessEvidence, bool) {
	if missing := missingCustomerPKIInputs(opts); len(missing) > 0 {
		// Missing PKI inputs leave the external authority unresolved; this is
		// cannot_verify rather than fail because no contradictory evidence exists.
		record.MissingIdentityFields = missing
		customerPKICannotVerify(record, ReasonMissingIdentity)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	if privateKeyInput(opts.CustomerPKIPublicCert) || privateKeyInput(opts.CustomerPKIPublicKey) {
		// Private key material in an input slot is a hard safety failure because
		// witness generation must never require or preserve signing secrets.
		customerPKIInputFail(record, ReasonPrivateKeyInput)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	return loadCustomerPKISafeInputs(record, opts)
}

func loadCustomerPKISafeInputs(record *Record, opts ProfileOptions) (ed25519.PublicKey, CustomerPKIAuthorityPolicy, CustomerPKIFreshnessEvidence, bool) {
	// Load the trust anchor before policy JSON so malformed key material cannot
	// be hidden behind otherwise well-formed policy metadata.
	publicKey, err := loadCustomerPublicKey(opts)
	if err != nil {
		customerPKICannotVerify(record, ReasonMalformedInput)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	policy, freshness, ok := loadCustomerPKIJSONInputs(record, opts)
	return publicKey, policy, freshness, ok
}

func loadCustomerPKIJSONInputs(record *Record, opts ProfileOptions) (CustomerPKIAuthorityPolicy, CustomerPKIFreshnessEvidence, bool) {
	var policy CustomerPKIAuthorityPolicy
	if err := readSafeJSON(opts.CustomerPKIAuthorityPolicy, &policy); err != nil {
		// Policy JSON is authority evidence. If it cannot be read safely, no
		// signer can be accepted for the profile.
		customerPKICannotVerify(record, ReasonPolicyMissing)
		return CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	var freshness CustomerPKIFreshnessEvidence
	if err := readSafeJSON(opts.CustomerPKIFreshness, &freshness); err != nil {
		// Freshness JSON links the signed payload to a time window and run ID; a
		// missing value keeps the profile open instead of failing unrelated gates.
		customerPKICannotVerify(record, ReasonMissingFreshness)
		return CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	return policy, freshness, true
}
