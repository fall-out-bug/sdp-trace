package export

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func buildAuditBundle(runDir string, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit, events []trace.Event) (AuditBundle, error) {
	// Re-open the run artifact so the exported manifest and default event list
	// come from the run directory, not from caller-provided memory.
	runArtifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return AuditBundle{}, err
	}
	events = auditBundleEvents(events, runArtifact.Events)

	return newAuditBundle(runArtifact.Manifest, result, table, audit, events), nil
}

func auditBundleEvents(events, artifactEvents []trace.Event) []trace.Event {
	if len(events) > 0 {
		// Explicit events let callers export a selected or redacted slice while
		// keeping the manifest sourced from the run artifact.
		return events
	}
	return artifactEvents
}

func newAuditBundle(run trace.RunManifest, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit, events []trace.Event) AuditBundle {
	// Field order mirrors review order: run identity, replay events, verifier
	// result, missing evidence, then optional integrity diagnostics.
	return AuditBundle{
		Run:       run,
		Events:    events,
		Result:    result,
		Missing:   table,
		Integrity: audit,
	}
}
