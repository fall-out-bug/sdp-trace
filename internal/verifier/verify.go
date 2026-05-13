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
	// Verification is a live replay boundary: every returned verdict below is
	// produced from the manifest, retained events, and contract available in
	// this run directory now.
	// Checked-in proof JSON never upgrades the outcome; missing or contradictory
	// source evidence keeps the result at cannot_verify, fail, or not_assessed.
	manifest, events, result, audit, err, ok := verifiedRunInput(runDir)
	if !ok {
		// Input failures return the verifier result and optional integrity audit
		// that explain why replay could not continue.
		return result, trace.MissingEvidenceTable{}, audit, err
	}

	// Observation is provisional until chain, contract, and missing-evidence
	// checks complete.
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

func verifiedChain(runDir string, manifest trace.RunManifest, events []trace.Event) (trace.VerifierResult, *trace.IntegrityAudit, bool) {
	chainOk, chainIssue := verifyChain(events, manifestChainHead(manifest))
	if !chainOk {
		// Chain contradictions are hard fail evidence rather than cannot_verify.
		return failedChainResult(runDir, manifest.RunID, chainIssue), audit(manifest.RunID, "tampered_chain", chainIssue, "run_dir", runDir), false
	}
	return trace.VerifierResult{}, nil, true
}

func manifestChainHead(manifest trace.RunManifest) string {
	if manifest.EventChainHead != "" {
		return manifest.EventChainHead
	}
	// Older manifests may only carry final_chain_head; use it as the replay
	// binding when event_chain_head is absent.
	return manifest.FinalChainHead
}

func failedChainResult(runDir, runID, issue string) trace.VerifierResult {
	// A broken event chain is fail evidence because replay contradicted retained
	// source artifacts. It is not downgraded to cannot_verify.
	return trace.VerifierResult{
		RunID:         runID,
		RunDir:        runDir,
		TrustScope:    trace.TrustScopeLocalObserved,
		Result:        trace.VerdictFail,
		Completeness:  trace.CompletenessMissingTelemetry,
		Replayability: trace.ReplayabilityNone,
		Reason:        issue,
	}
}

func verifiedContract(manifest trace.RunManifest, result trace.VerifierResult) (trace.Contract, trace.MissingEvidenceTable, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	contract := trace.DefaultContract
	if manifest.ContractPath != "" {
		// Explicit contract paths take precedence because they bind run-specific
		// requirements.
		return verifiedManifestContract(manifest, result, contract)
	} else if manifest.ContractDigest != "" {
		if manifest.ContractDigest != contractDigest(trace.DefaultContract) {
			// A default-contract digest mismatch means the verifier's default has
			// drifted since run capture.
			return contract, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, "default contract digest mismatch"), audit(manifest.RunID, "contract_digest_mismatch", "default contract digest changed after run", "", ""), nil, false
		}
	}
	return contract, trace.MissingEvidenceTable{}, result, nil, nil, true
}

func verifiedManifestContract(manifest trace.RunManifest, result trace.VerifierResult, fallback trace.Contract) (trace.Contract, trace.MissingEvidenceTable, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	resolvedContract, err := trace.LoadContract(filepath.Clean(manifest.ContractPath))
	if err != nil {
		// Keep the manifest contract id in the missing-evidence table even when
		// the contract file cannot be loaded.
		return fallback, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, err.Error()), audit(manifest.RunID, "contract_unreadable", err.Error(), "contract_path", manifest.ContractPath), err, false
	}
	if manifest.ContractDigest != "" && manifest.ContractDigest != contractDigest(resolvedContract) {
		// Contract digest mismatch blocks replay because requirements may have
		// changed after the run.
		return fallback, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, "contract digest mismatch"), audit(manifest.RunID, "contract_digest_mismatch", "contract digest changed after run", "contract_path", manifest.ContractPath), nil, false
	}
	return resolvedContract, trace.MissingEvidenceTable{}, result, nil, nil, true
}

func contractDigest(contract trace.Contract) string {
	return trace.SHA256Hex(string(mustMarshalJSON(contract)))
}

func observedResult(runDir, runID string) trace.VerifierResult {
	// Observed is the optimistic local replay baseline before requirement gaps
	// are applied.
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
	// Cannot-verify results still carry local trust scope because the verifier
	// only assessed local artifacts.
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
	// Preserve run identity while lowering replayability and completeness.
	result.Result = trace.VerdictCannotVerify
	result.Completeness = trace.CompletenessUnknown
	result.Replayability = trace.ReplayabilityNone
	result.Reason = reason
	return result
}

func resultWithMissingEvidence(result trace.VerifierResult, contract trace.Contract, events []trace.Event) (trace.VerifierResult, trace.MissingEvidenceTable) {
	missingEvidence := trace.GenerateMissingEvidenceTable(contract, observedEventSet(events))
	if len(missingEvidence.Rows) > 0 {
		// Missing required evidence is not_assessed: the replay succeeded, but
		// the contract is not fully covered.
		result.Result = trace.VerdictNotAssessed
		result.Completeness = trace.CompletenessMissingTelemetry
		result.Replayability = trace.ReplayabilityPartial
	}
	return result, missingEvidence
}

func observedEventSet(events []trace.Event) map[string]bool {
	observedEvents := map[string]bool{}
	for _, event := range events {
		// Contract coverage is event-type based; individual event payloads remain
		// chain evidence.
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
		// Details stay optional so audit rows without a stable field do not
		// fabricate empty metadata.
		result.Details = map[string]string{detailKey: detailValue}
	}
	return result
}

func mustMarshalJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		// Contract digesting should stay total for impossible marshal failures.
		return []byte("{}")
	}
	return data
}

func loadRunEvents(runDir string) ([]trace.Event, error) {
	eventDir := filepath.Join(runDir, "events")
	eventFiles, err := loadRunEventFiles(eventDir)
	if err != nil {
		// The public error names the events directory as the missing replay
		// component.
		return nil, fmt.Errorf("events directory missing: %w", err)
	}
	sort.Strings(eventFiles)
	return parseRunEventFiles(eventFiles)
}

func loadRunEventFiles(eventDir string) ([]string, error) {
	files, err := os.ReadDir(eventDir)
	if err != nil {
		return nil, err
	}

	// Directory order is not trusted. The caller sorts the accepted filenames
	// before replay so filesystem enumeration cannot influence the chain.
	eventFiles := make([]string, 0, len(files))
	for _, entry := range files {
		eventFiles = appendRunEventFile(eventFiles, eventDir, entry)
	}
	return eventFiles, nil
}

func appendRunEventFile(eventFiles []string, eventDir string, entry os.DirEntry) []string {
	if !isValidEventFile(entry) {
		// Ignore non-event files so auxiliary artifacts cannot affect replay
		// order.
		return eventFiles
	}
	return append(eventFiles, filepath.Join(eventDir, entry.Name()))
}

func isValidEventFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		// Event chains are flat files only.
		return false
	}
	name := entry.Name()
	if !strings.HasSuffix(name, ".json") {
		// Only JSON event files participate in chain replay.
		return false
	}
	_, err := trace.EventSeqFromFilename(name)
	// Filename sequence parsing rejects unrelated JSON files.
	return err == nil
}

func parseRunEventFiles(eventFiles []string) ([]trace.Event, error) {
	if err := requireRunEventFiles(eventFiles); err != nil {
		return nil, err
	}
	// Parsing preserves one event per accepted file. Structural trust checks
	// stay in chain verification so parse errors and replay contradictions have
	// distinct verifier outcomes.
	events := make([]trace.Event, 0, len(eventFiles))
	for _, path := range eventFiles {
		var err error
		events, err = appendParsedRunEvent(events, path)
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

func requireRunEventFiles(eventFiles []string) error {
	if len(eventFiles) == 0 {
		// No valid event files means no replayable telemetry.
		return errors.New("no event files")
	}
	return nil
}

func appendParsedRunEvent(events []trace.Event, path string) ([]trace.Event, error) {
	event, err := parseRunEvent(path)
	if err != nil {
		// Surface only the basename to keep verifier errors portable.
		return nil, fmt.Errorf("invalid event %s: %w", filepath.Base(path), err)
	}
	return append(events, event), nil
}

func parseRunEvent(path string) (trace.Event, error) {
	var event trace.Event
	if err := trace.ReadJSON(context.Background(), path, &event); err != nil {
		return trace.Event{}, err
	}
	// Structural event validation happens during chain verification.
	return event, nil
}

func verifyChain(events []trace.Event, expectedHead string) (bool, string) {
	if issue := verifyChainEvents(events); issue != "" {
		// Event-level chain defects take precedence over manifest head mismatch.
		return false, issue
	}
	if issue := verifyExpectedHead(events, expectedHead); issue != "" {
		return false, issue
	}
	return true, ""
}

func verifyChainEvents(events []trace.Event) string {
	for i, event := range events {
		if issue := verifyChainEvent(events, i, event); issue != "" {
			// Return first defect to keep diagnostics deterministic.
			return issue
		}
	}
	return ""
}

func verifyExpectedHead(events []trace.Event, expectedHead string) string {
	if expectedHead != "" && events[len(events)-1].EventHash != expectedHead {
		// final_chain_head is optional, but when present it binds the retained
		// event files to the manifest head.
		return "run head does not match manifest final_chain_head"
	}
	return ""
}

func verifyChainEvent(events []trace.Event, index int, event trace.Event) string {
	// Checks are ordered from cheapest structural guard to strongest hash proof.
	// firstChainIssue short-circuits so later evidence cannot mask the earliest
	// chain defect.
	checks := []chainCheck{
		func() string { return verifyChainSequence(index, event) },
		func() string { return verifyChainPrevHash(events, index, event) },
		func() string { return verifyChainPayloadDigest(event) },
		func() string { return verifyChainEventHash(event) },
	}
	return firstChainIssue(checks)
}

type chainCheck func() string

func firstChainIssue(checks []chainCheck) string {
	for _, check := range checks {
		if issue := check(); issue != "" {
			// Return first defect to keep diagnostics deterministic.
			return issue
		}
	}
	return ""
}

func verifyChainSequence(index int, event trace.Event) string {
	if event.Sequence != index {
		// Event sequence is zero-based and must match file replay order.
		return fmt.Sprintf("sequence mismatch at %s", event.EventID)
	}
	return ""
}

func verifyChainPrevHash(events []trace.Event, index int, event trace.Event) string {
	if event.PrevEventHash != chainExpectedPrevHash(events, index) {
		if index == 0 {
			// The first event must link to the null sentinel only.
			return "first event has non-empty prev_event_hash"
		}
		return fmt.Sprintf("broken chain at %d (%s)", index+1, event.EventID)
	}
	return ""
}

func chainExpectedPrevHash(events []trace.Event, index int) string {
	if index == 0 {
		// Genesis events link to the canonical null hash sentinel.
		return trace.NullEventHash
	}
	return events[index-1].EventHash
}

func verifyChainPayloadDigest(event trace.Event) string {
	if err := event.VerifyPayloadDigest(); err != nil {
		// Payload-digest failure blocks chain trust even if event_hash matches a
		// malformed payload copy.
		return fmt.Sprintf("invalid payload digest for %s: %s", event.EventID, err)
	}
	return ""
}

func verifyChainEventHash(event trace.Event) string {
	recomputed, err := trace.EventHash(event)
	if err != nil {
		// Hash recomputation can fail when canonical event shape is invalid.
		return fmt.Sprintf("invalid event hash for %s", event.EventID)
	}
	if event.EventHash != recomputed {
		// Event hash mismatch proves retained event contents no longer match the
		// recorded chain.
		return fmt.Sprintf("hash mismatch for %s", event.EventID)
	}
	return ""
}

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
