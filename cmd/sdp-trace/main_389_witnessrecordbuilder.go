package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

type witnessRecordBuilder func(witnessOptions) (witness.Record, error)
