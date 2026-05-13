package demo

import (
	"os"
)

func persistReportArtifacts(outDir string, artifacts ReportArtifacts) error {
	// persistReportArtifacts keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	return writeReportArtifacts(outDir, artifacts)
}
