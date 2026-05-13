package verifier

import (
	"context"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifiedRunInput(runDir string) (trace.RunManifest, []trace.Event, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	// Manifest and event loading are separated so each failure can produce a
	// precise audit issue.
	manifest, result, audit, err, ok := verifiedManifest(runDir)
	if !ok {
		return manifest, nil, result, audit, err, false
	}
	events, result, audit, err, ok := verifiedEvents(runDir, manifest)
	return manifest, events, result, audit, err, ok
}

func verifiedManifest(runDir string) (trace.RunManifest, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	manifestPath := filepath.Join(runDir, "run.json")
	var manifest trace.RunManifest
	if err := trace.ReadJSON(context.Background(), manifestPath, &manifest); err != nil {
		// Missing manifest prevents any source-bound replay.
		return manifest, cannotVerifyResult(runDir, "", "run manifest missing"), audit("", "run_manifest_missing", "run.json is required for verification", "run_dir", runDir), nil, false
	}
	if err := manifest.Validate(); err != nil {
		// Invalid manifests are returned as errors because callers can repair the
		// run metadata directly.
		return manifest, cannotVerifyResult(runDir, "", err.Error()), audit("", "run_manifest_invalid", err.Error(), "run_dir", runDir), err, false
	}
	return manifest, trace.VerifierResult{}, nil, nil, true
}

func verifiedEvents(runDir string, manifest trace.RunManifest) ([]trace.Event, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	events, err := loadRunEvents(runDir)
	if err != nil {
		// Event load errors keep the result cannot_verify rather than fail; the
		// verifier cannot prove tampering without event replay.
		return nil, cannotVerifyResult(runDir, "", err.Error()), audit(manifest.RunID, "event_load_failed", err.Error(), "run_dir", runDir), nil, false
	}
	if len(events) == 0 {
		// Empty event sets are missing telemetry, not a successful empty run.
		result := cannotVerifyResult(runDir, "", "no events")
		result.Completeness = trace.CompletenessMissingTelemetry
		return nil, result, audit(manifest.RunID, "empty_events", "run contains no events", "run_dir", runDir), nil, false
	}
	return events, trace.VerifierResult{}, nil, nil, true
}
