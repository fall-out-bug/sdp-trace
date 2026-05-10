package query

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestForensicsBasicPackDerivesRowsWithoutPolicyVerdict(t *testing.T) {
	runDir := writeForensicsPackFixture(t)

	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	if result.SchemaVersion != QueryPackSchemaVersion ||
		result.QueryPackID != QueryPackForensicsBasic ||
		result.QueryPackVersion != "v1" {
		t.Fatalf("identity = %+v", result)
	}
	if result.TopLevelAssessment != "" {
		t.Fatalf("query pack emitted top-level assessment: %q", result.TopLevelAssessment)
	}
	if len(result.InputArtifacts) != 3 {
		t.Fatalf("input artifacts = %+v", result.InputArtifacts)
	}
	assertNoUnsafeMarkers(t, result)

	redactionRows := result.QueryRows[QueryForensicsRedactions]
	digestCap := findRow(t, redactionRows, "critical_evidence_reconstructable")
	if digestCap.EvidenceState != RowStateRetentionLimited ||
		digestCap.EvidenceFamily != EvidenceFamilyRetention ||
		digestCap.Reconstructable == nil ||
		*digestCap.Reconstructable ||
		digestCap.ReasonCode != "digest_only_not_reconstructable" {
		t.Fatalf("digest cap row = %+v", digestCap)
	}
	if digestCap.SourceRef != "block_18.condition.critical_evidence_reconstructable" ||
		digestCap.SourceConditionState != "fail" {
		t.Fatalf("digest cap source = %+v", digestCap)
	}

	captureRows := result.QueryRows[QueryForensicsCaptureDepth]
	missingTelemetry := findRow(t, captureRows, "tool_call_depth_visible")
	if missingTelemetry.EvidenceState != RowStateMissingTelemetry ||
		missingTelemetry.EvidenceFamily != EvidenceFamilyAdapterCapture ||
		missingTelemetry.EvidenceGap == "" {
		t.Fatalf("missing telemetry row = %+v", missingTelemetry)
	}
	if strings.HasPrefix(missingTelemetry.ID, "event_") {
		t.Fatalf("row id leaks event sequence: %s", missingTelemetry.ID)
	}

	summaryRows := result.QueryRows[QueryForensicsSummary]
	if len(summaryRows) == 0 || len(summaryRows[0].RelatedRows) == 0 {
		t.Fatalf("summary rows should reference other query rows: %+v", summaryRows)
	}
}

func TestForensicsBasicPackPreservesMissingRequiredArtifact(t *testing.T) {
	runDir := writeForensicsPackFixture(t)
	if err := os.Remove(filepath.Join(runDir, "adapter-capture.assessment-result.json")); err != nil {
		t.Fatalf("remove adapter result: %v", err)
	}

	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	row := result.QueryRows[QueryForensicsCaptureDepth][0]
	if row.EvidenceState != RowStateCannotVerify ||
		row.EvidenceFamily != EvidenceFamilyAdapterCapture ||
		row.ReasonCode != "missing_block_19_adapter_capture_result" {
		t.Fatalf("missing adapter row = %+v", row)
	}
}

func TestTimelineRowsFallbackWhenRunHasNoEventRefsAndOptionalArtifactsMissing(t *testing.T) {
	runDir := t.TempDir()
	writeQueryPackJSON(t, filepath.Join(runDir, "run.json"), map[string]any{
		"run_id": "run-no-events",
	})

	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	timeline := result.QueryRows[QueryForensicsTimeline]
	if len(timeline) != 3 {
		t.Fatalf("timeline rows = %+v", timeline)
	}
	if timeline[0].ReasonCode != "run_timeline_available" || timeline[0].SourceRef != "block_09.run.run_id" {
		t.Fatalf("fallback run row = %+v", timeline[0])
	}
	if timeline[1].EvidenceState != RowStateNotAssessed || timeline[1].EvidenceGap != EvidenceFamilyRetention {
		t.Fatalf("missing forensic row = %+v", timeline[1])
	}
	if timeline[2].EvidenceState != RowStateNotAssessed || timeline[2].EvidenceGap != EvidenceFamilyAdapterCapture {
		t.Fatalf("missing adapter row = %+v", timeline[2])
	}
}

func TestForensicsBasicPackPreservesUnreadableOptionalArtifactAsCannotVerifyRows(t *testing.T) {
	runDir := writeForensicsPackFixture(t)
	adapterPath := filepath.Join(runDir, "adapter-capture.assessment-result.json")
	if err := os.Remove(adapterPath); err != nil {
		t.Fatalf("remove adapter result: %v", err)
	}
	if err := os.Mkdir(adapterPath, 0o755); err != nil {
		t.Fatalf("replace adapter result with directory: %v", err)
	}

	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	row := findMatrixRow(t, result.QueryRows[QueryForensicsGaps], struct {
		Scenario        string `json:"scenario"`
		Query           string `json:"query"`
		EvidenceState   string `json:"evidence_state"`
		EvidenceFamily  string `json:"evidence_family"`
		ReasonCode      string `json:"reason_code"`
		SourceRefPrefix string `json:"source_ref_prefix"`
		SourceState     string `json:"source_condition_state,omitempty"`
		EvidenceGap     string `json:"evidence_gap,omitempty"`
		Reconstructable *bool  `json:"reconstructable,omitempty"`
	}{
		ReasonCode: "unreadable_or_malformed_input_artifact",
	})
	if row.EvidenceState != RowStateCannotVerify || row.EvidenceFamily != EvidenceFamilyInputArtifact {
		t.Fatalf("unreadable optional artifact row = %+v", row)
	}
	artifact := findInputArtifact(t, result.InputArtifacts, "adapter_capture")
	if artifact.SHA256 != "" {
		t.Fatalf("unreadable optional artifact unexpectedly has digest: %+v", artifact)
	}
}

func TestForensicsBasicPackMapsEmptyUpstreamStateToCannotVerify(t *testing.T) {
	runDir := writeForensicsPackFixture(t)
	writeQueryPackJSON(t, filepath.Join(runDir, "forensic-retention.assessment-result.json"), map[string]any{
		"schema_version":                "block18-forensic-retention-assessment-v1",
		"selected_profile":              "forensic_retention",
		"forensic_retention_assessment": "cannot_verify",
		"trust_scope":                   "local_observed",
		"forensic_conditions": []map[string]any{
			{
				"id":          "retention_mode_declared",
				"state":       "",
				"reason_code": "empty_state_fixture",
				"reason":      "state is empty",
			},
		},
	})

	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	row := findRow(t, result.QueryRows[QueryForensicsRedactions], "retention_mode_declared")
	if row.EvidenceState != RowStateCannotVerify || row.EvidenceGap != EvidenceFamilyRetention {
		t.Fatalf("empty state row = %+v", row)
	}
}

func TestExplainForensicsPackRendersStableSafeRows(t *testing.T) {
	runDir := writeForensicsPackFixture(t)
	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	explanation := ExplainForensicsPack(result)
	if !strings.Contains(explanation, "forensics-redactions redactions.0001 retention_limited retention source_ref=block_18.condition.critical_evidence_reconstructable source_condition_id=critical_evidence_reconstructable source_condition_state=fail reconstructable=false") {
		t.Fatalf("explanation missing digest cap: %s", explanation)
	}
	if strings.Contains(explanation, "secret-token") ||
		strings.Contains(explanation, "rotate-vault-key") ||
		strings.Contains(explanation, "s3://") {
		t.Fatalf("explanation leaked unsafe marker")
	}
}

func TestCommittedBlock20FixtureMatrix(t *testing.T) {
	fixtureRoot := filepath.Join("..", "..", "examples", "block20-forensics-query-pack")
	matrixPath := filepath.Join(fixtureRoot, "fixture-matrix.jsonl")
	file, err := os.Open(matrixPath)
	if err != nil {
		t.Fatalf("open fixture matrix: %v", err)
	}
	defer file.Close()

	scenarioDirs := committedScenarioDirs(t, fixtureRoot)
	seenScenarios := map[string]bool{}
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		rows++
		var expected struct {
			Scenario        string `json:"scenario"`
			Query           string `json:"query"`
			EvidenceState   string `json:"evidence_state"`
			EvidenceFamily  string `json:"evidence_family"`
			ReasonCode      string `json:"reason_code"`
			SourceRefPrefix string `json:"source_ref_prefix"`
			SourceState     string `json:"source_condition_state,omitempty"`
			EvidenceGap     string `json:"evidence_gap,omitempty"`
			Reconstructable *bool  `json:"reconstructable,omitempty"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &expected); err != nil {
			t.Fatalf("decode matrix row %d: %v", rows, err)
		}
		if !scenarioDirs[expected.Scenario] {
			t.Fatalf("matrix row references uncommitted scenario %q", expected.Scenario)
		}
		seenScenarios[expected.Scenario] = true
		runDir := filepath.Join(fixtureRoot, expected.Scenario)
		result, err := ForensicsBasicPack(runDir)
		if err != nil {
			t.Fatalf("%s pack: %v", expected.Scenario, err)
		}
		row := findMatrixRow(t, result.QueryRows[expected.Query], expected)
		if row.EvidenceState != expected.EvidenceState || row.EvidenceFamily != expected.EvidenceFamily {
			t.Fatalf("%s row mismatch = %+v expected=%+v", expected.Scenario, row, expected)
		}
		if expected.SourceRefPrefix == "" || !strings.HasPrefix(row.SourceRef, expected.SourceRefPrefix) {
			t.Fatalf("%s source_ref = %q expected prefix %q", expected.Scenario, row.SourceRef, expected.SourceRefPrefix)
		}
		if strings.HasPrefix(row.SourceRef, "block_09.") && (row.SourceConditionID != "" || row.SourceConditionState != "") {
			t.Fatalf("%s Block 09 row %s must omit source condition fields", expected.Scenario, row.ID)
		}
		if expected.SourceState != "" && row.SourceConditionState != expected.SourceState {
			t.Fatalf("%s source_condition_state = %q expected %q", expected.Scenario, row.SourceConditionState, expected.SourceState)
		}
		if expected.EvidenceGap != "" && row.EvidenceGap != expected.EvidenceGap {
			t.Fatalf("%s evidence_gap = %q expected %q", expected.Scenario, row.EvidenceGap, expected.EvidenceGap)
		}
		if (row.EvidenceState == RowStateCannotVerify || row.EvidenceState == RowStateNotAssessed) && row.EvidenceGap == "" {
			t.Fatalf("%s %s row has no evidence_gap", expected.Scenario, row.EvidenceState)
		}
		if expected.Reconstructable != nil {
			if row.Reconstructable == nil || *row.Reconstructable != *expected.Reconstructable {
				t.Fatalf("%s reconstructable = %+v expected=%t", expected.Scenario, row.Reconstructable, *expected.Reconstructable)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan matrix: %v", err)
	}
	if rows < 10 {
		t.Fatalf("fixture matrix too small: %d", rows)
	}
	for scenario := range scenarioDirs {
		if !seenScenarios[scenario] {
			t.Fatalf("committed scenario %q has no matrix row", scenario)
		}
	}
}

func TestForensicsBasicPackSafetyClassesAreVerifiedAgainstOutput(t *testing.T) {
	runDir := writeForensicsPackFixture(t)
	markers := map[string]string{}
	for i, class := range sensitiveClasses() {
		markers[class] = fmt.Sprintf("synthetic-sensitive-sentinel-%02d", i)
	}
	writeQueryPackJSON(t, filepath.Join(runDir, "run.json"), map[string]any{
		"run_id":          "run-block20",
		"run_nonce":       "nonce-block20",
		"source_baseline": "sha256:source",
		"event_refs": []map[string]any{
			{"event_type": "command_started", "uri": markers["raw_command_arguments"]},
			{"event_type": "test_observed", "uri": markers["unsafe_test_identifiers"]},
		},
		"verifier_states": map[string]any{
			"event_chain_structurally_valid": map[string]any{"state": "pass", "reason": markers["source_snippets"]},
			"event_chain_witnessed":          map[string]any{"state": "not_assessed", "reason": markers["prompts"]},
		},
	})
	writeQueryPackJSON(t, filepath.Join(runDir, "forensic-retention.assessment-result.json"), map[string]any{
		"schema_version":                "block18-forensic-retention-assessment-v1",
		"selected_profile":              "forensic_retention",
		"forensic_retention_assessment": "fail",
		"trust_scope":                   "local_observed",
		"forensic_conditions": []map[string]any{
			{"id": "critical_evidence_reconstructable", "state": "fail", "reason_code": "critical_evidence_digest_only", "reason": markers["stdout_stderr_bodies"], "capped_to_retention_mode": "digest_only"},
			{"id": "redaction_unresolved_visible", "state": "pass", "reason_code": "redaction_resolved", "reason": markers["raw_review_bodies"]},
		},
	})
	writeQueryPackJSON(t, filepath.Join(runDir, "adapter-capture.assessment-result.json"), map[string]any{
		"schema_version":             "block19-adapter-capture-assessment-v1",
		"selected_profile":           "adapter_capture",
		"adapter_capture_assessment": "cannot_verify",
		"trust_scope":                "local_observed",
		"adapter_capture_conditions": []map[string]any{
			{"id": "tool_call_depth_visible", "state": "missing_telemetry", "reason_code": "tool_event_missing", "reason": markers["tool_payloads"]},
			{"id": "provider_refs_portable", "state": "fail", "reason_code": "provider_ref_contains_secret", "reason": markers["authenticated_provider_urls"]},
			{"id": "redaction_metadata_consistent", "state": "fail", "reason_code": "adapter_config_unsafe", "reason": markers["adapter_configuration"]},
			{"id": "model_identity_not_overclaimed", "state": "fail", "reason_code": "raw_model_payload_seen", "reason": markers["raw_model_payloads"]},
			{"id": "file_mutation_correlated", "state": "fail", "reason_code": "file_path_seen", "reason": markers["executable_paths"]},
			{"id": "test_provenance_not_overclaimed", "state": "fail", "reason_code": "script_path_seen", "reason": markers["script_paths"]},
			{"id": "run_binding_established", "state": "fail", "reason_code": "gateway_ref_seen", "reason": markers["gateway_evidence_refs"]},
			{"id": "adapter_identity_visible", "state": "fail", "reason_code": "credential_seen", "reason": markers["credentials"]},
			{"id": "adapter_event_contract_valid", "state": "fail", "reason_code": "token_seen", "reason": markers["tokens"]},
			{"id": "capture_depth_not_overclaimed", "state": "fail", "reason_code": "raw_reference_note_seen", "reason": markers["unsafe_raw_reference_access_notes"]},
			{"id": "task_drift_visible", "state": "fail", "reason_code": "key_material_seen", "reason": markers["key_material"]},
			{"id": "test_provenance_not_overclaimed", "state": "fail", "reason_code": "command_name_seen", "reason": markers["command_names"]},
		},
	})

	result, err := ForensicsBasicPack(runDir)
	if err != nil {
		t.Fatalf("forensics pack: %v", err)
	}
	outputClasses := map[string]bool{}
	for _, class := range result.OutputSafety.VerifiedAbsentSensitiveClasses {
		outputClasses[class] = true
	}
	payload, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	explanation := ExplainForensicsPack(result)
	for class, marker := range markers {
		if !outputClasses[class] {
			t.Fatalf("output safety class %q missing", class)
		}
		assertMarkerAbsent(t, payload, marker)
		assertMarkerAbsent(t, []byte(explanation), marker)
	}
}

func TestSafeTokenPreservesAllowedCharsAndDropsUnsafe(t *testing.T) {
	t.Run("keepsAllowedCharacters", func(t *testing.T) {
		got := safeToken("abc_DEF-012")
		if got != "abc_DEF-012" {
			t.Fatalf("safeToken preserved = %q", got)
		}
	})

	t.Run("dropsUnsafeCharacters", func(t *testing.T) {
		got := safeToken("a b+c:d/e.f:g")
		if got != "abcdefg" {
			t.Fatalf("safeToken filtered = %q", got)
		}
	})

	t.Run("normalizesEmptyResult", func(t *testing.T) {
		got := safeToken(" a+b ")
		if got != "ab" {
			t.Fatalf("safeToken removed all unsafe chars = %q", got)
		}
	})

	t.Run("dropsNonAscii", func(t *testing.T) {
		got := safeToken("a-b_🙂-1")
		if got != "a-b_-1" {
			t.Fatalf("safeToken filtered unicode = %q", got)
		}
	})
}

func TestSafeTokenUnknownForEmptyOrFullyUnsafe(t *testing.T) {
	t.Run("emptyValue", func(t *testing.T) {
		if got := safeToken(""); got != "unknown" {
			t.Fatalf("empty value = %q", got)
		}
	})

	t.Run("preservesSingleSafeCharacter", func(t *testing.T) {
		if got := safeToken("a!@#$"); got != "a" {
			t.Fatalf("fully unsafe with one safe char = %q", got)
		}
	})

	t.Run("allUnsafe", func(t *testing.T) {
		if got := safeToken(" !*"); got != "unknown" {
			t.Fatalf("all-unsafe value = %q", got)
		}
	})
}

func writeForensicsPackFixture(t *testing.T) string {
	t.Helper()
	runDir := t.TempDir()
	writeQueryPackJSON(t, filepath.Join(runDir, "run.json"), map[string]any{
		"run_id":          "run-block20",
		"run_nonce":       "nonce-block20",
		"source_baseline": "sha256:source",
		"event_refs": []map[string]any{
			{"sequence": 7, "event_hash": "eventhash", "event_type": "command_started", "uri": "events/secret-command.json"},
			{"sequence": 8, "event_hash": "eventhash2", "event_type": "test_observed", "uri": "events/rotate-vault-key.json"},
		},
		"verifier_states": map[string]any{
			"event_chain_structurally_valid": map[string]any{"state": "pass"},
			"event_chain_witnessed":          map[string]any{"state": "not_assessed", "reason": "no witness"},
		},
	})
	writeQueryPackJSON(t, filepath.Join(runDir, "forensic-retention.assessment-result.json"), map[string]any{
		"schema_version":                "block18-forensic-retention-assessment-v1",
		"selected_profile":              "forensic_retention",
		"forensic_retention_assessment": "fail",
		"trust_scope":                   "local_observed",
		"forensic_conditions": []map[string]any{
			{
				"id":                       "critical_evidence_reconstructable",
				"state":                    "fail",
				"reason_code":              "critical_evidence_digest_only",
				"reason":                   "critical evidence retained as digest only",
				"capped_to_retention_mode": "digest_only",
			},
			{
				"id":          "redaction_unresolved_visible",
				"state":       "pass",
				"reason_code": "redaction_resolved",
				"reason":      "redaction resolved",
			},
		},
	})
	writeQueryPackJSON(t, filepath.Join(runDir, "adapter-capture.assessment-result.json"), map[string]any{
		"schema_version":             "block19-adapter-capture-assessment-v1",
		"selected_profile":           "adapter_capture",
		"adapter_capture_assessment": "cannot_verify",
		"trust_scope":                "local_observed",
		"adapter_capture_conditions": []map[string]any{
			{
				"id":          "tool_call_depth_visible",
				"state":       "missing_telemetry",
				"reason_code": "tool_event_missing",
				"reason":      "required tool-call adapter event is missing",
			},
			{
				"id":          "provider_refs_portable",
				"state":       "fail",
				"reason_code": "provider_ref_contains_secret",
				"reason":      "provider reference contains unsafe material",
			},
		},
	})
	return runDir
}

func writeQueryPackJSON(t *testing.T, path string, value any) {
	t.Helper()
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}

func findInputArtifact(t *testing.T, artifacts []QueryPackInputArtifact, role string) QueryPackInputArtifact {
	t.Helper()
	for _, artifact := range artifacts {
		if artifact.Role == role {
			return artifact
		}
	}
	t.Fatalf("input artifact role %s not found in %+v", role, artifacts)
	return QueryPackInputArtifact{}
}

func committedScenarioDirs(t *testing.T, fixtureRoot string) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(fixtureRoot)
	if err != nil {
		t.Fatalf("read fixture root: %v", err)
	}
	scenarios := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() {
			scenarios[entry.Name()] = true
		}
	}
	if len(scenarios) == 0 {
		t.Fatalf("no committed fixture scenarios found")
	}
	return scenarios
}

func findRow(t *testing.T, rows []QueryPackRow, sourceConditionID string) QueryPackRow {
	t.Helper()
	for _, row := range rows {
		if row.SourceConditionID == sourceConditionID {
			return row
		}
	}
	t.Fatalf("row with source condition %s not found in %+v", sourceConditionID, rows)
	return QueryPackRow{}
}

func findMatrixRow(t *testing.T, rows []QueryPackRow, expected struct {
	Scenario        string `json:"scenario"`
	Query           string `json:"query"`
	EvidenceState   string `json:"evidence_state"`
	EvidenceFamily  string `json:"evidence_family"`
	ReasonCode      string `json:"reason_code"`
	SourceRefPrefix string `json:"source_ref_prefix"`
	SourceState     string `json:"source_condition_state,omitempty"`
	EvidenceGap     string `json:"evidence_gap,omitempty"`
	Reconstructable *bool  `json:"reconstructable,omitempty"`
}) QueryPackRow {
	t.Helper()
	for _, row := range rows {
		if row.ReasonCode == expected.ReasonCode {
			return row
		}
	}
	t.Fatalf("matrix row reason %s not found in %+v", expected.ReasonCode, rows)
	return QueryPackRow{}
}

func assertMarkerAbsent(t *testing.T, payload []byte, marker string) {
	t.Helper()
	if strings.Contains(string(payload), marker) {
		digest := sha256.Sum256([]byte(marker))
		t.Fatalf("query-pack output leaked synthetic marker digest=%x", digest)
	}
}

func assertNoUnsafeMarkers(t *testing.T, result QueryPackResult) {
	t.Helper()
	payload, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	for _, marker := range []string{"secret-token", "rotate-vault-key", "s3://", "events/secret-command.json"} {
		if strings.Contains(string(payload), marker) {
			t.Fatalf("query pack leaked unsafe marker digest=%x", marker)
		}
	}
}
