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
