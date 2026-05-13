package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func readQueryPackResult(path string) (query.QueryPackResult, error) {
	var result query.QueryPackResult
	if err := readJSONFile(path, &result); err != nil {
		// Artifact read failures are verification failures, not empty results.
		return query.QueryPackResult{}, err
	}
	return result, nil
}
