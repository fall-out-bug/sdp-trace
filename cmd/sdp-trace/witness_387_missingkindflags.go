package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func missingWitnessKindFlags(kind string, opts *flagSet) []string {
	if kind != witness.KindCustomerPKI {
		// Non-customer-PKI witnesses do not require customer key custody inputs.
		return nil
	}
	return missingCustomerPKIFlags(opts)
}
