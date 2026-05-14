package main

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func validateQueryPackExplainResult(result query.QueryPackResult) error {
	if result.SchemaVersion != query.QueryPackSchemaVersion || result.QueryPackID != query.QueryPackForensicsBasic {
		// Explain only understands the current forensics-basic result contract.
		return fmt.Errorf("unsupported query-pack result")
	}
	// The query package owns detailed result validation; CLI validation only
	// gates the schema/profile pair before rendering.
	return nil
}
