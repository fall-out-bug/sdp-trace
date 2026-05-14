package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func buildCustomerPKIWitness(opts witnessOptions) (witness.Record, error) {
	// Customer-PKI witnesses require explicit authority, public credential,
	// payload digest, and freshness evidence paths.
	return witness.WriteProfile(witness.KindCustomerPKI, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		CustomerPKIAuthorityPolicy: opts.customerPKIAuthorityPath,
		CustomerPKIPublicCert:      opts.customerPKIPublicCertPath,
		CustomerPKIPublicKey:       opts.customerPKIPublicKeyPath,
		CustomerPKIPayloadDigest:   opts.customerPKIPayloadDigest,
		CustomerPKIFreshness:       opts.customerPKIFreshnessPath,
	})
}
