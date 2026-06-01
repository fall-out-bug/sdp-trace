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

func requireQueryPackRequiredInputs(runPath, outPath string) error {
	if runPath == "" {
		// The run path is the replayable source evidence for this pack.
		return fmt.Errorf("query-pack requires --run")
	}
	if outPath == "" {
		// The output path is the durable artifact reviewed by later commands.
		return fmt.Errorf("query-pack requires --out")
	}
	return nil
}

func readQueryPackResult(path string) (query.QueryPackResult, error) {
	var result query.QueryPackResult
	if err := readJSONFile(path, &result); err != nil {
		// Artifact read failures are verification failures, not empty results.
		return query.QueryPackResult{}, err
	}
	return result, nil
}

func validateQueryPackExplainResult(result query.QueryPackResult) error {
	if result.SchemaVersion != query.QueryPackSchemaVersion || result.QueryPackID != query.QueryPackForensicsBasic {
		// Explain only understands the current forensics-basic result contract.
		return fmt.Errorf("unsupported query-pack result")
	}
	// The query package owns detailed result validation; CLI validation only
	// gates the schema/profile pair before rendering.
	return nil
}
