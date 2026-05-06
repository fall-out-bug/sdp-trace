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
	manifestPath := filepath.Join(runDir, "run.json")
	var manifest trace.RunManifest
	if err := trace.ReadJSON(context.Background(), manifestPath, &manifest); err != nil {
		return trace.VerifierResult{
				RunDir:       runDir,
				Result:       trace.VerdictCannotVerify,
				TrustScope:   trace.TrustScopeLocalObserved,
				Completeness: trace.CompletenessUnknown,
				Reason:       "run manifest missing",
			}, trace.MissingEvidenceTable{}, &trace.IntegrityAudit{
				RunID:   "",
				Issue:   "run_manifest_missing",
				Reason:  "run.json is required for verification",
				Details: map[string]string{"run_dir": runDir},
			}, nil
	}
	if err := manifest.Validate(); err != nil {
		return trace.VerifierResult{
				RunDir:       runDir,
				Result:       trace.VerdictCannotVerify,
				TrustScope:   trace.TrustScopeLocalObserved,
				Completeness: trace.CompletenessUnknown,
				Reason:       err.Error(),
			}, trace.MissingEvidenceTable{}, &trace.IntegrityAudit{
				RunID:   "",
				Issue:   "run_manifest_invalid",
				Reason:  err.Error(),
				Details: map[string]string{"run_dir": runDir},
			}, err
	}

	events, err := loadRunEvents(runDir)
	if err != nil {
		return trace.VerifierResult{
				RunDir:       runDir,
				Result:       trace.VerdictCannotVerify,
				TrustScope:   trace.TrustScopeLocalObserved,
				Completeness: trace.CompletenessUnknown,
				Reason:       err.Error(),
			}, trace.MissingEvidenceTable{}, &trace.IntegrityAudit{
				RunID:   manifest.RunID,
				Issue:   "event_load_failed",
				Reason:  err.Error(),
				Details: map[string]string{"run_dir": runDir},
			}, nil
	}
	if len(events) == 0 {
		return trace.VerifierResult{
				RunDir:       runDir,
				Result:       trace.VerdictCannotVerify,
				TrustScope:   trace.TrustScopeLocalObserved,
				Completeness: trace.CompletenessMissingTelemetry,
				Reason:       "no events",
			}, trace.MissingEvidenceTable{}, &trace.IntegrityAudit{
				RunID:   manifest.RunID,
				Issue:   "empty_events",
				Reason:  "run contains no events",
				Details: map[string]string{"run_dir": runDir},
			}, nil
	}

	result := trace.VerifierResult{
		RunID:         manifest.RunID,
		RunDir:        runDir,
		TrustScope:    trace.TrustScopeLocalObserved,
		Result:        trace.VerdictObserved,
		Completeness:  trace.CompletenessComplete,
		Replayability: trace.ReplayabilityPartial,
	}

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
			}, trace.MissingEvidenceTable{}, &trace.IntegrityAudit{
				RunID:   manifest.RunID,
				Issue:   "tampered_chain",
				Reason:  chainIssue,
				Details: map[string]string{"run_dir": runDir},
			}, nil
	}

	contract := trace.DefaultContract
	if manifest.ContractPath != "" {
		resolvedContract, err := trace.LoadContract(filepath.Clean(manifest.ContractPath))
		if err != nil {
			result.Result = trace.VerdictCannotVerify
			result.Completeness = trace.CompletenessUnknown
			result.Replayability = trace.ReplayabilityNone
			result.Reason = err.Error()
			return result, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, &trace.IntegrityAudit{
				RunID:   manifest.RunID,
				Issue:   "contract_unreadable",
				Reason:  err.Error(),
				Details: map[string]string{"contract_path": manifest.ContractPath},
			}, err
		}
		contractDigest := trace.SHA256Hex(string(mustMarshalJSON(resolvedContract)))
		if manifest.ContractDigest != "" && manifest.ContractDigest != contractDigest {
			result.Result = trace.VerdictCannotVerify
			result.Completeness = trace.CompletenessUnknown
			result.Replayability = trace.ReplayabilityNone
			result.Reason = "contract digest mismatch"
			return result, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, &trace.IntegrityAudit{
				RunID:   manifest.RunID,
				Issue:   "contract_digest_mismatch",
				Reason:  "contract digest changed after run",
				Details: map[string]string{"contract_path": manifest.ContractPath},
			}, nil
		}
		contract = resolvedContract
	} else if manifest.ContractDigest != "" {
		defaultDigest := trace.SHA256Hex(string(mustMarshalJSON(trace.DefaultContract)))
		if manifest.ContractDigest != defaultDigest {
			result.Result = trace.VerdictCannotVerify
			result.Completeness = trace.CompletenessUnknown
			result.Replayability = trace.ReplayabilityNone
			result.Reason = "default contract digest mismatch"
			return result, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, &trace.IntegrityAudit{
				RunID:  manifest.RunID,
				Issue:  "contract_digest_mismatch",
				Reason: "default contract digest changed after run",
			}, nil
		}
	}

	observedEvents := map[string]bool{}
	for _, event := range events {
		observedEvents[string(event.EventType)] = true
	}
	missingEvidence := trace.GenerateMissingEvidenceTable(contract, observedEvents)
	if len(missingEvidence.Rows) > 0 {
		result.Result = trace.VerdictNotAssessed
		result.Completeness = trace.CompletenessMissingTelemetry
		result.Replayability = trace.ReplayabilityPartial
	}
	return result, missingEvidence, nil, nil
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
