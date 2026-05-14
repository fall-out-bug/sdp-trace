package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func writeQueryPackArtifact(opts *queryPackOptions) (int, error) {
	result, err := query.ForensicsBasicPack(opts.runPath)
	if err != nil {
		// Pack generation depends on replayable run evidence.
		return exitCannotVerify, err
	}
	if err := writeJSONFile(opts.outPath, result); err != nil {
		// A generated pack that cannot be persisted is not review evidence.
		return 1, err
	}
	return 0, nil
}
