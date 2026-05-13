package demo

import (
	"os"
	"path/filepath"
)

func writeReportArtifacts(outDir string, artifacts ReportArtifacts) error {
	// writeReportArtifacts keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	writes := []struct {
		name  string
		value any
	}{
		{name: "summary.json", value: artifacts.Summary},
		{name: "evidence-table.json", value: artifacts.EvidenceTable},
		{name: "missing-telemetry.json", value: artifacts.MissingTelemetry},
	}
	for _, write := range writes {
		if err := writeJSON(filepath.Join(outDir, write.name), write.value); err != nil {

			return err
		}
	}

	return os.WriteFile(filepath.Join(outDir, "timeline.md"), []byte(artifacts.Timeline), 0o644)
}
