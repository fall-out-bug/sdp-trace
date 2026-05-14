package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func buildEnvelopeWitness(opts witnessOptions) (witness.Record, error) {
	// Envelope-backed profiles all use the supplied envelope as their portable
	// provenance input.
	return witness.WriteProfile(opts.kind, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		EnvelopePath: opts.witnessEnvelope,
	})
}
