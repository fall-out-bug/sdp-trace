package export

import "github.com/fall_out_bug/sdp-trace/internal/trace"

// AuditBundle is a small reproducible package exported by the command layer.
// It carries verifier output together with run evidence, without making the
// bundle itself a new trust authority.
type AuditBundle struct {
	Run       trace.RunManifest          `json:"run"`
	Events    []trace.Event              `json:"events"`
	Result    trace.VerifierResult       `json:"result"`
	Missing   trace.MissingEvidenceTable `json:"missing_evidence"`
	Integrity *trace.IntegrityAudit      `json:"integrity_audit,omitempty"`
}

// BuildAuditBundle composes run-level artifacts into an exportable structure.
func BuildAuditBundle(runDir string, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit, events []trace.Event) (AuditBundle, error) {
	return buildAuditBundle(runDir, result, table, audit, events)
}
