package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func readGateExplainResult(path string, stderr io.Writer) (demo.GateResult, int, bool) {
	var result demo.GateResult
	if err := readJSONFile(path, &result); err != nil {
		// Missing gate artifacts are cannot_verify for explanation, not usage.
		fmt.Fprintln(stderr, err)
		return demo.GateResult{}, exitCannotVerify, false
	}
	if result.SchemaVersion != demo.GateSchemaVersion && result.SchemaVersion != demo.GateSchemaVersionBlock16 {
		// Unsupported result schemas remain cannot_verify instead of being
		// rendered with stale field assumptions.
		fmt.Fprintf(stderr, "unsupported gate-result schema_version: %s\n", result.SchemaVersion)
		return demo.GateResult{}, exitCannotVerify, false
	}
	return result, 0, true
}
