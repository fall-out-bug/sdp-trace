package demo

import (
	"os"
	"path/filepath"
)

func persistGateResult(outPath string, result GateResult) error {
	// persistGateResult keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return writeJSON(outPath, result)
}
