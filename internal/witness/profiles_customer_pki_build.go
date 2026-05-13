package witness

import (
	"crypto/ed25519"
)

func BuildCustomerPKI(runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	// Customer PKI can only establish external witness trust after local
	// artifacts, authority policy, public key, and freshness evidence all bind.
	// Early cannot_verify/fail records are returned with artifact digests already
	// captured where possible.
	// The happy path sets CI identity from freshness evidence only after the
	// external input files load safely.
	record, err := newCustomerPKIRecord(runsRoot, reportDir)
	if err != nil {
		return Record{}, err
	}
	// Customer PKI starts with local artifact hashes, then imports only the
	// minimal external inputs needed to authorize those artifacts.
	inputs, ok := prepareCustomerPKIValidation(&record, opts)
	if !ok {
		return record, nil
	}
	if !validateCustomerPKIRecord(&record, inputs.states, inputs.publicKey, inputs.policy, inputs.freshness, runsRoot, opts.CustomerPKIPayloadDigest) {
		return record, nil
	}
	applyCustomerPKIPass(&record)
	return record, nil
}

type customerPKIValidationInputs struct {
	publicKey ed25519.PublicKey
	policy    CustomerPKIAuthorityPolicy
	freshness CustomerPKIFreshnessEvidence
	states    *ProfileStates
}

func prepareCustomerPKIValidation(record *Record, opts ProfileOptions) (customerPKIValidationInputs, bool) {
	// External inputs become validation context only after safety checks have
	// accepted the public trust anchor, authority policy, and freshness evidence.
	publicKey, policy, freshness, ok := loadCustomerPKIInputs(record, opts)
	if !ok {
		return customerPKIValidationInputs{}, false
	}
	record.CI = CIIdentity{Provider: KindCustomerPKI, RunID: freshness.RunID}
	states := customerPKIPassStates(policy)
	record.ProfileStates = states
	return customerPKIValidationInputs{publicKey: publicKey, policy: policy, freshness: freshness, states: states}, true
}
