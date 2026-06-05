package main

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func buildWitnessRecord(opts witnessOptions) (witness.Record, error) {
	builder, ok := witnessRecordBuilders()[opts.kind]
	if !ok {
		// This should be unreachable after option validation; keep the error so
		// direct helper misuse cannot silently produce a generic witness.
		return witness.Record{}, fmt.Errorf("unsupported witness kind %q", opts.kind)
	}
	return builder(opts)
}

type witnessRecordBuilder func(witnessOptions) (witness.Record, error)

func witnessRecordBuilders() map[string]witnessRecordBuilder {
	// Builders map closed witness kinds to the package function that owns their
	// evidence interpretation.
	return map[string]witnessRecordBuilder{
		witness.KindGitHubActions: buildGitHubActionsWitness,
		witness.KindGitLabCI:      buildEnvelopeWitness,
		witness.KindBuildkite:     buildEnvelopeWitness,
		witness.KindCustomerPKI:   buildCustomerPKIWitness,
	}
}

func buildGitHubActionsWitness(opts witnessOptions) (witness.Record, error) {
	return witness.WriteGitHubActions(opts.out, opts.target, opts.reportDir, witness.EnvironmentFromOS())
}

func buildEnvelopeWitness(opts witnessOptions) (witness.Record, error) {
	// Envelope-backed profiles all use the supplied envelope as their portable
	// provenance input.
	return witness.WriteProfile(opts.kind, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		EnvelopePath: opts.witnessEnvelope,
	})
}

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
