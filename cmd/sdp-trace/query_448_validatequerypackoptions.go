package main

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func validateQueryPackOptions(opts *queryPackOptions) error {
	if opts.pack == "" {
		return fmt.Errorf("error: ambiguous pack selection; --pack is required")
	}
	if opts.pack != query.QueryPackForensicsBasic {
		// The CLI exposes a closed pack vocabulary so unknown pack names fail as
		// usage errors before any evidence is read or written.
		return fmt.Errorf("error: unknown pack %q", opts.pack)
	}
	return requireQueryPackRequiredInputs(opts.runPath, opts.outPath)
}
