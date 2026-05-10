package verifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// VerifyRun reads a run directory and returns a verifiable result object.
func VerifyRun(runDir string) (trace.VerifierResult, trace.MissingEvidenceTable, *trace.IntegrityAudit, error) {
	manifest, events, result, audit, err, ok := verifiedRunInput(runDir)
	if !ok {
		return result, trace.MissingEvidenceTable{}, audit, err
	}

	result = observedResult(runDir, manifest.RunID)
	chainResult, chainAudit, ok := verifiedChain(runDir, manifest, events)
	if !ok {
		return chainResult, trace.MissingEvidenceTable{}, chainAudit, nil
	}

	contract, table, result, audit, err, ok := verifiedContract(manifest, result)
	if !ok {
		return result, table, audit, err
	}
	result, missingEvidence := resultWithMissingEvidence(result, contract, events)
	return result, missingEvidence, nil, nil
}

func verifiedRunInput(runDir string) (trace.RunManifest, []trace.Event, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
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
		return manifest, cannotVerifyResult(runDir, "", "run manifest missing"), audit("", "run_manifest_missing", "run.json is required for verification", "run_dir", runDir), nil, false
	}
	if err := manifest.Validate(); err != nil {
		return manifest, cannotVerifyResult(runDir, "", err.Error()), audit("", "run_manifest_invalid", err.Error(), "run_dir", runDir), err, false
	}
	return manifest, trace.VerifierResult{}, nil, nil, true
}

func verifiedEvents(runDir string, manifest trace.RunManifest) ([]trace.Event, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	events, err := loadRunEvents(runDir)
	if err != nil {
		return nil, cannotVerifyResult(runDir, "", err.Error()), audit(manifest.RunID, "event_load_failed", err.Error(), "run_dir", runDir), nil, false
	}
	if len(events) == 0 {
		result := cannotVerifyResult(runDir, "", "no events")
		result.Completeness = trace.CompletenessMissingTelemetry
		return nil, result, audit(manifest.RunID, "empty_events", "run contains no events", "run_dir", runDir), nil, false
	}
	return events, trace.VerifierResult{}, nil, nil, true
}

func verifiedChain(runDir string, manifest trace.RunManifest, events []trace.Event) (trace.VerifierResult, *trace.IntegrityAudit, bool) {
	chainHead := manifest.EventChainHead
	if chainHead == "" {
		chainHead = manifest.FinalChainHead
	}
	chainOk, chainIssue := verifyChain(events, chainHead)
	if !chainOk {
		return trace.VerifierResult{
			RunID:         manifest.RunID,
			RunDir:        runDir,
			TrustScope:    trace.TrustScopeLocalObserved,
			Result:        trace.VerdictFail,
			Completeness:  trace.CompletenessMissingTelemetry,
			Replayability: trace.ReplayabilityNone,
			Reason:        chainIssue,
		}, audit(manifest.RunID, "tampered_chain", chainIssue, "run_dir", runDir), false
	}
	return trace.VerifierResult{}, nil, true
}

func verifiedContract(manifest trace.RunManifest, result trace.VerifierResult) (trace.Contract, trace.MissingEvidenceTable, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	contract := trace.DefaultContract
	if manifest.ContractPath != "" {
		resolvedContract, err := trace.LoadContract(filepath.Clean(manifest.ContractPath))
		if err != nil {
			return contract, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, err.Error()), audit(manifest.RunID, "contract_unreadable", err.Error(), "contract_path", manifest.ContractPath), err, false
		}
		contractDigest := trace.SHA256Hex(string(mustMarshalJSON(resolvedContract)))
		if manifest.ContractDigest != "" && manifest.ContractDigest != contractDigest {
			return contract, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, "contract digest mismatch"), audit(manifest.RunID, "contract_digest_mismatch", "contract digest changed after run", "contract_path", manifest.ContractPath), nil, false
		}
		contract = resolvedContract
	} else if manifest.ContractDigest != "" {
		defaultDigest := trace.SHA256Hex(string(mustMarshalJSON(trace.DefaultContract)))
		if manifest.ContractDigest != defaultDigest {
			return contract, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, "default contract digest mismatch"), audit(manifest.RunID, "contract_digest_mismatch", "default contract digest changed after run", "", ""), nil, false
		}
	}
	return contract, trace.MissingEvidenceTable{}, result, nil, nil, true
}

func observedResult(runDir, runID string) trace.VerifierResult {
	return trace.VerifierResult{
		RunID:         runID,
		RunDir:        runDir,
		TrustScope:    trace.TrustScopeLocalObserved,
		Result:        trace.VerdictObserved,
		Completeness:  trace.CompletenessComplete,
		Replayability: trace.ReplayabilityPartial,
	}
}

func cannotVerifyResult(runDir, runID, reason string) trace.VerifierResult {
	return trace.VerifierResult{
		RunID:        runID,
		RunDir:       runDir,
		Result:       trace.VerdictCannotVerify,
		TrustScope:   trace.TrustScopeLocalObserved,
		Completeness: trace.CompletenessUnknown,
		Reason:       reason,
	}
}

func cannotVerifyReplayResult(result trace.VerifierResult, reason string) trace.VerifierResult {
	result.Result = trace.VerdictCannotVerify
	result.Completeness = trace.CompletenessUnknown
	result.Replayability = trace.ReplayabilityNone
	result.Reason = reason
	return result
}

func resultWithMissingEvidence(result trace.VerifierResult, contract trace.Contract, events []trace.Event) (trace.VerifierResult, trace.MissingEvidenceTable) {
	missingEvidence := trace.GenerateMissingEvidenceTable(contract, observedEventSet(events))
	if len(missingEvidence.Rows) > 0 {
		result.Result = trace.VerdictNotAssessed
		result.Completeness = trace.CompletenessMissingTelemetry
		result.Replayability = trace.ReplayabilityPartial
	}
	return result, missingEvidence
}

func observedEventSet(events []trace.Event) map[string]bool {
	observedEvents := map[string]bool{}
	for _, event := range events {
		observedEvents[string(event.EventType)] = true
	}
	return observedEvents
}

func audit(runID, issue, reason, detailKey, detailValue string) *trace.IntegrityAudit {
	result := &trace.IntegrityAudit{
		RunID:  runID,
		Issue:  issue,
		Reason: reason,
	}
	if detailKey != "" {
		result.Details = map[string]string{detailKey: detailValue}
	}
	return result
}

func mustMarshalJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func loadRunEvents(runDir string) ([]trace.Event, error) {
	files, err := os.ReadDir(filepath.Join(runDir, "events"))
	if err != nil {
		return nil, fmt.Errorf("events directory missing: %w", err)
	}
	eventFiles := make([]string, 0, len(files))
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		_, err := trace.EventSeqFromFilename(name)
		if err != nil {
			continue
		}
		eventFiles = append(eventFiles, filepath.Join(runDir, "events", name))
	}
	if len(eventFiles) == 0 {
		return nil, errors.New("no event files")
	}
	sort.Strings(eventFiles)
	events := make([]trace.Event, 0, len(eventFiles))
	for _, path := range eventFiles {
		var event trace.Event
		if err := trace.ReadJSON(context.Background(), path, &event); err != nil {
			return nil, fmt.Errorf("invalid event %s: %w", filepath.Base(path), err)
		}
		events = append(events, event)
	}
	return events, nil
}

func verifyChain(events []trace.Event, expectedHead string) (bool, string) {
	for i, event := range events {
		if event.Sequence != i {
			return false, fmt.Sprintf("sequence mismatch at %s", event.EventID)
		}
		if i == 0 && event.PrevEventHash != trace.NullEventHash {
			return false, "first event has non-empty prev_event_hash"
		}
		if i > 0 && event.PrevEventHash != events[i-1].EventHash {
			return false, fmt.Sprintf("broken chain at %d (%s)", i+1, event.EventID)
		}
		if err := event.VerifyPayloadDigest(); err != nil {
			return false, fmt.Sprintf("invalid payload digest for %s: %s", event.EventID, err)
		}
		recomputed, err := trace.EventHash(event)
		if err != nil {
			return false, fmt.Sprintf("invalid event hash for %s", event.EventID)
		}
		if event.EventHash != recomputed {
			return false, fmt.Sprintf("hash mismatch for %s", event.EventID)
		}
	}
	if expectedHead != "" && events[len(events)-1].EventHash != expectedHead {
		return false, "run head does not match manifest final_chain_head"
	}
	return true, ""
}

// WriteVerifierArtifacts writes verifier result and missing-evidence table for later query.
func WriteVerifierArtifacts(runDir string, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit) error {
	verifierDir := filepath.Join(runDir, "verifier")
	if err := os.MkdirAll(verifierDir, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(verifierDir, "verifier-result.json"), result); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(verifierDir, "missing-evidence-table.json"), table); err != nil {
		return err
	}
	if audit != nil {
		return writeJSON(filepath.Join(verifierDir, "integrity-audit.json"), audit)
	}
	return nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
