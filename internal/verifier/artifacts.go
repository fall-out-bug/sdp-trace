package verifier

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// WriteVerifierArtifacts writes verifier result and missing-evidence table for later query.
func WriteVerifierArtifacts(runDir string, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit) error {
	verifierDir := filepath.Join(runDir, "verifier")
	if err := os.MkdirAll(verifierDir, 0o755); err != nil {
		return err
	}
	// Write result, missing evidence, and optional audit as separate artifacts so
	// downstream query packs can consume them independently.
	if err := writeJSON(filepath.Join(verifierDir, "verifier-result.json"), result); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(verifierDir, "missing-evidence-table.json"), table); err != nil {
		return err
	}
	return writeIntegrityAudit(verifierDir, audit)
}

func writeIntegrityAudit(verifierDir string, audit *trace.IntegrityAudit) error {
	if audit == nil {
		// Integrity audit is emitted only when there is an assessed structural
		// issue; absence is not a green proof.
		return nil
	}
	return writeJSON(filepath.Join(verifierDir, "integrity-audit.json"), audit)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	// Verifier artifacts are human-reviewable JSON; event hash authority remains
	// in the replayed run chain.
	return os.WriteFile(path, data, 0o644)
}
