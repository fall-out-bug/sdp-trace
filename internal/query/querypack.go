package query

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	QueryPackSchemaVersion  = "block20-forensics-query-pack-result-v1"
	QueryPackForensicsBasic = "forensics-basic-v1"

	QueryForensicsSummary          = "forensics-summary"
	QueryForensicsTimeline         = "forensics-timeline"
	QueryForensicsGaps             = "forensics-gaps"
	QueryForensicsRedactions       = "forensics-redactions"
	QueryForensicsCaptureDepth     = "forensics-capture-depth"
	QueryForensicsUnverifiedClaims = "forensics-unverified-claims"

	RowStatePresent          = "present"
	RowStateIssueObserved    = "issue_observed"
	RowStateNotAssessed      = "not_assessed"
	RowStateCannotVerify     = "cannot_verify"
	RowStateMissingTelemetry = "missing_telemetry"
	RowStateUnsupported      = "unsupported"
	RowStateNotIntegrated    = "not_integrated"
	RowStateRetentionLimited = "retention_limited"

	EvidenceFamilyRunChain       = "run_chain"
	EvidenceFamilyWitness        = "witness"
	EvidenceFamilyRetention      = "retention"
	EvidenceFamilyRedaction      = "redaction"
	EvidenceFamilyAdapterCapture = "adapter_capture"
	EvidenceFamilyTask           = "task"
	EvidenceFamilyCommand        = "command"
	EvidenceFamilyFileMutations  = "file_mutations"
	EvidenceFamilyTest           = "test"
	EvidenceFamilySupersession   = "supersession"
	EvidenceFamilyClaim          = "claim"
	EvidenceFamilyInputArtifact  = "input_artifact"
)

var queryOrder = []string{
	QueryForensicsSummary,
	QueryForensicsTimeline,
	QueryForensicsGaps,
	QueryForensicsRedactions,
	QueryForensicsCaptureDepth,
	QueryForensicsUnverifiedClaims,
}

type QueryPackResult struct {
	SchemaVersion      string                    `json:"schema_version"`
	QueryPackID        string                    `json:"query_pack_id"`
	QueryPackVersion   string                    `json:"query_pack_version"`
	RunID              string                    `json:"run_id,omitempty"`
	RunNonce           string                    `json:"run_nonce,omitempty"`
	SourceBaseline     string                    `json:"source_baseline,omitempty"`
	TopLevelAssessment string                    `json:"top_level_assessment,omitempty"`
	InputArtifacts     []QueryPackInputArtifact  `json:"input_artifacts"`
	QueryRows          map[string][]QueryPackRow `json:"query_rows"`
	OutputSafety       QueryPackOutputSafety     `json:"output_safety"`
}

type QueryPackInputArtifact struct {
	Role             string `json:"role"`
	SHA256           string `json:"sha256"`
	PathRedactedID   string `json:"path_redacted_id"`
	SchemaVersion    string `json:"schema_version,omitempty"`
	ArtifactRequired bool   `json:"artifact_required"`
}

type QueryPackRow struct {
	ID                   string   `json:"id"`
	Query                string   `json:"query"`
	EvidenceState        string   `json:"evidence_state"`
	EvidenceFamily       string   `json:"evidence_family"`
	Reconstructable      *bool    `json:"reconstructable,omitempty"`
	SourceRef            string   `json:"source_ref"`
	SourceConditionID    string   `json:"source_condition_id,omitempty"`
	SourceConditionState string   `json:"source_condition_state,omitempty"`
	ReasonCode           string   `json:"reason_code,omitempty"`
	EvidenceGap          string   `json:"evidence_gap,omitempty"`
	RelatedRows          []string `json:"related_rows,omitempty"`
}

type QueryPackOutputSafety struct {
	VerifiedAbsentSensitiveClasses []string `json:"verified_absent_sensitive_classes"`
	RedactionPolicyDigest          string   `json:"redaction_policy_digest,omitempty"`
}

type packInputs struct {
	run              runArtifact
	runArtifact      QueryPackInputArtifact
	runErr           error
	forensic         assessmentEnvelope
	forensicArtifact *QueryPackInputArtifact
	forensicPresent  bool
	forensicErr      error
	adapter          assessmentEnvelope
	adapterArtifact  *QueryPackInputArtifact
	adapterPresent   bool
	adapterErr       error
}

type runArtifact struct {
	SchemaVersion   string                     `json:"schema_version,omitempty"`
	RunID           string                     `json:"run_id"`
	RunNonce        string                     `json:"run_nonce,omitempty"`
	SourceBaseline  string                     `json:"source_baseline,omitempty"`
	EventRefs       []eventRef                 `json:"event_refs,omitempty"`
	VerifierStates  map[string]verifierState   `json:"verifier_states,omitempty"`
	RedactionDigest string                     `json:"redaction_policy_digest,omitempty"`
	Raw             map[string]json.RawMessage `json:"-"`
}

type eventRef struct {
	EventType string `json:"event_type"`
}

type verifierState struct {
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
}

type assessmentEnvelope struct {
	SchemaVersion            string                `json:"schema_version"`
	ForensicConditions       []assessmentCondition `json:"forensic_conditions,omitempty"`
	AdapterCaptureConditions []assessmentCondition `json:"adapter_capture_conditions,omitempty"`
	ForensicAssessment       string                `json:"forensic_retention_assessment,omitempty"`
	AdapterCaptureAssessment string                `json:"adapter_capture_assessment,omitempty"`
}

type assessmentCondition struct {
	ID                    string `json:"id"`
	State                 string `json:"state"`
	ReasonCode            string `json:"reason_code"`
	Reason                string `json:"reason"`
	CappedToRetentionMode string `json:"capped_to_retention_mode,omitempty"`
}

func ForensicsBasicPack(runDir string) (QueryPackResult, error) {
	inputs, err := loadPackInputs(runDir)
	if err != nil {
		return QueryPackResult{}, err
	}
	builder := newPackBuilder(inputs)
	if inputs.runErr != nil {
		builder.addMalformedRequiredInputRows()
		return builder.result(), nil
	}
	builder.addTimelineRows()
	builder.addRedactionRows()
	builder.addCaptureRows()
	builder.addGapRows()
	builder.addUnverifiedClaimRows()
	builder.addSummaryRows()
	result := builder.result()
	return result, nil
}

func ExplainForensicsPack(result QueryPackResult) string {
	var lines []string
	for _, queryName := range queryOrder {
		rows := append([]QueryPackRow(nil), result.QueryRows[queryName]...)
		sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
		for _, row := range rows {
			parts := []string{queryName, row.ID, row.EvidenceState, row.EvidenceFamily}
			parts = append(parts, "source_ref="+row.SourceRef)
			if row.SourceConditionID != "" {
				parts = append(parts, "source_condition_id="+row.SourceConditionID)
			}
			if row.SourceConditionState != "" {
				parts = append(parts, "source_condition_state="+row.SourceConditionState)
			}
			if row.Reconstructable != nil {
				parts = append(parts, fmt.Sprintf("reconstructable=%t", *row.Reconstructable))
			}
			if row.EvidenceGap != "" {
				parts = append(parts, "gap="+row.EvidenceGap)
			}
			lines = append(lines, strings.Join(parts, " "))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func loadPackInputs(runDir string) (packInputs, error) {
	runPath := filepath.Join(runDir, "run.json")
	var run runArtifact
	runArtifact, err := readPackArtifact(runPath, "run", "run", true, &run)
	if err != nil && runArtifact.Role == "" {
		return packInputs{}, err
	}
	inputs := packInputs{run: run, runArtifact: runArtifact, runErr: err}
	forensicPath := filepath.Join(runDir, "forensic-retention.assessment-result.json")
	var forensic assessmentEnvelope
	if artifact, present, err := readOptionalPackArtifact(forensicPath, "forensic_retention", "forensic_retention", false, &forensic); err != nil && artifact.Role == "" {
		return packInputs{}, err
	} else if present {
		inputs.forensicPresent = true
		inputs.forensicArtifact = &artifact
		inputs.forensic = forensic
		inputs.forensicErr = err
	}
	adapterPath := filepath.Join(runDir, "adapter-capture.assessment-result.json")
	var adapter assessmentEnvelope
	if artifact, present, err := readOptionalPackArtifact(adapterPath, "adapter_capture", "adapter_capture", false, &adapter); err != nil && artifact.Role == "" {
		return packInputs{}, err
	} else if present {
		inputs.adapterPresent = true
		inputs.adapterArtifact = &artifact
		inputs.adapter = adapter
		inputs.adapterErr = err
	}
	return inputs, nil
}

func readOptionalPackArtifact(path, role, redactedID string, required bool, target any) (QueryPackInputArtifact, bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return QueryPackInputArtifact{}, false, nil
		}
		return QueryPackInputArtifact{}, true, err
	}
	artifact, err := readPackArtifact(path, role, redactedID, required, target)
	return artifact, true, err
}

func readPackArtifact(path, role, redactedID string, required bool, target any) (QueryPackInputArtifact, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return QueryPackInputArtifact{}, err
	}
	sum := sha256.Sum256(payload)
	artifact := QueryPackInputArtifact{
		Role:             role,
		SHA256:           hex.EncodeToString(sum[:]),
		PathRedactedID:   redactedID,
		ArtifactRequired: required,
	}
	var envelope struct {
		SchemaVersion string `json:"schema_version"`
	}
	if err := json.Unmarshal(payload, &envelope); err == nil {
		artifact.SchemaVersion = envelope.SchemaVersion
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return artifact, err
	}
	return artifact, nil
}

type packBuilder struct {
	inputs   packInputs
	rows     map[string][]QueryPackRow
	counters map[string]int
}

func newPackBuilder(inputs packInputs) *packBuilder {
	rows := map[string][]QueryPackRow{}
	for _, queryName := range queryOrder {
		rows[queryName] = []QueryPackRow{}
	}
	return &packBuilder{inputs: inputs, rows: rows, counters: map[string]int{}}
}

func (b *packBuilder) result() QueryPackResult {
	artifacts := []QueryPackInputArtifact{b.inputs.runArtifact}
	if b.inputs.forensicArtifact != nil {
		artifacts = append(artifacts, *b.inputs.forensicArtifact)
	}
	if b.inputs.adapterArtifact != nil {
		artifacts = append(artifacts, *b.inputs.adapterArtifact)
	}
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Role < artifacts[j].Role })
	return QueryPackResult{
		SchemaVersion:    QueryPackSchemaVersion,
		QueryPackID:      QueryPackForensicsBasic,
		QueryPackVersion: "v1",
		RunID:            b.inputs.run.RunID,
		RunNonce:         b.inputs.run.RunNonce,
		SourceBaseline:   b.inputs.run.SourceBaseline,
		InputArtifacts:   artifacts,
		QueryRows:        b.rows,
		OutputSafety: QueryPackOutputSafety{
			RedactionPolicyDigest:          b.inputs.run.RedactionDigest,
			VerifiedAbsentSensitiveClasses: sensitiveClasses(),
		},
	}
}

func (b *packBuilder) addTimelineRows() {
	if len(b.inputs.run.EventRefs) == 0 {
		b.addRow(QueryForensicsTimeline, RowStatePresent, EvidenceFamilyRunChain, "block_09.run.run_id", "", "", "run_timeline_available", "")
	} else {
		for i, event := range b.inputs.run.EventRefs {
			family := familyForEvent(event.EventType)
			sourceRef := fmt.Sprintf("block_09.event.%s.e%04d", family, i+1)
			b.addRow(QueryForensicsTimeline, RowStatePresent, family, sourceRef, "", "", "timeline_event_present", "")
		}
	}
	if !b.inputs.forensicPresent {
		b.addRow(QueryForensicsTimeline, RowStateNotAssessed, EvidenceFamilyRetention, "block_18.condition.missing", "", "", "missing_optional_block_18_forensic_retention_result", EvidenceFamilyRetention)
	}
	if !b.inputs.adapterPresent {
		b.addRow(QueryForensicsTimeline, RowStateNotAssessed, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "", "", "missing_optional_block_19_adapter_capture_result", EvidenceFamilyAdapterCapture)
	}
}

func (b *packBuilder) addMalformedRequiredInputRows() {
	for _, queryName := range queryOrder {
		b.addRow(queryName, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_09.run.malformed", "", "", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}

func (b *packBuilder) addRedactionRows() {
	if !b.inputs.forensicPresent {
		b.addRow(QueryForensicsRedactions, RowStateCannotVerify, EvidenceFamilyRedaction, "block_18.condition.missing", "", "", "missing_block_18_forensic_retention_result", "redaction")
		return
	}
	if b.inputs.forensicErr != nil {
		b.addRow(QueryForensicsRedactions, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_18.condition.malformed", "", "", "unreadable_or_malformed_input_artifact", "input_artifact")
		return
	}
	for _, condition := range b.inputs.forensic.ForensicConditions {
		family := familyForForensicCondition(condition.ID)
		row := b.rowFromCondition(QueryForensicsRedactions, family, "block_18.condition."+safeToken(condition.ID), condition)
		b.rows[QueryForensicsRedactions] = append(b.rows[QueryForensicsRedactions], row)
	}
}

func (b *packBuilder) addCaptureRows() {
	if !b.inputs.adapterPresent {
		b.addRow(QueryForensicsCaptureDepth, RowStateCannotVerify, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "", "", "missing_block_19_adapter_capture_result", "adapter_capture")
		return
	}
	if b.inputs.adapterErr != nil {
		b.addRow(QueryForensicsCaptureDepth, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_19.condition.malformed", "", "", "unreadable_or_malformed_input_artifact", "input_artifact")
		return
	}
	for _, condition := range b.inputs.adapter.AdapterCaptureConditions {
		family := familyForAdapterCondition(condition.ID)
		row := b.rowFromCondition(QueryForensicsCaptureDepth, family, "block_19.condition."+safeToken(condition.ID), condition)
		b.rows[QueryForensicsCaptureDepth] = append(b.rows[QueryForensicsCaptureDepth], row)
	}
}

func (b *packBuilder) addGapRows() {
	keys := make([]string, 0, len(b.inputs.run.VerifierStates))
	for key := range b.inputs.run.VerifierStates {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		state := b.inputs.run.VerifierStates[key]
		rowState := mapSourceState(state.State)
		if rowState == RowStatePresent {
			continue
		}
		family := familyForVerifierState(key)
		b.addRow(QueryForensicsGaps, rowState, family, "block_09.run."+safeToken(key), "", state.State, safeToken(key), family)
	}
	if !b.inputs.forensicPresent {
		b.addRow(QueryForensicsGaps, RowStateNotAssessed, EvidenceFamilyRetention, "block_18.condition.missing", "", "", "missing_optional_block_18_forensic_retention_result", "retention")
	}
	if !b.inputs.adapterPresent {
		b.addRow(QueryForensicsGaps, RowStateNotAssessed, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "", "", "missing_optional_block_19_adapter_capture_result", "adapter_capture")
	}
}

func (b *packBuilder) addUnverifiedClaimRows() {
	for _, row := range append([]QueryPackRow{}, b.rows[QueryForensicsRedactions]...) {
		if row.EvidenceState != RowStatePresent {
			b.addReferencedClaim(row)
		}
	}
	for _, row := range append([]QueryPackRow{}, b.rows[QueryForensicsCaptureDepth]...) {
		if row.EvidenceState != RowStatePresent {
			b.addReferencedClaim(row)
		}
	}
}

func (b *packBuilder) addSummaryRows() {
	for _, queryName := range queryOrder {
		if queryName == QueryForensicsSummary {
			continue
		}
		var related []string
		for _, row := range b.rows[queryName] {
			related = append(related, row.ID)
		}
		if len(related) == 0 {
			continue
		}
		row := b.newRow(QueryForensicsSummary, RowStatePresent, EvidenceFamilyClaim, "block_09.run.run_id", "", "", "query_group_index", "")
		row.RelatedRows = related
		b.rows[QueryForensicsSummary] = append(b.rows[QueryForensicsSummary], row)
	}
}

func (b *packBuilder) addReferencedClaim(source QueryPackRow) {
	row := b.newRow(QueryForensicsUnverifiedClaims, source.EvidenceState, EvidenceFamilyClaim, source.SourceRef, source.SourceConditionID, source.SourceConditionState, source.ReasonCode, source.EvidenceGap)
	row.RelatedRows = []string{source.ID}
	row.Reconstructable = source.Reconstructable
	b.rows[QueryForensicsUnverifiedClaims] = append(b.rows[QueryForensicsUnverifiedClaims], row)
}

func (b *packBuilder) rowFromCondition(queryName, family, sourceRef string, condition assessmentCondition) QueryPackRow {
	state := mapSourceState(condition.State)
	reason := condition.ReasonCode
	gap := ""
	if state != RowStatePresent && state != RowStateIssueObserved {
		gap = family
	}
	reconstructable := (*bool)(nil)
	if condition.State == RowStateRetentionLimited || condition.CappedToRetentionMode != "" {
		falseValue := false
		reconstructable = &falseValue
	}
	if condition.ID == "critical_evidence_reconstructable" && (condition.CappedToRetentionMode != "" || condition.ReasonCode == "critical_evidence_digest_only") {
		state = RowStateRetentionLimited
		reason = "digest_only_not_reconstructable"
		falseValue := false
		reconstructable = &falseValue
		gap = EvidenceFamilyRetention
	}
	row := b.newRow(queryName, state, family, sourceRef, condition.ID, condition.State, reason, gap)
	row.Reconstructable = reconstructable
	return row
}

func (b *packBuilder) addRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap string) {
	row := b.newRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap)
	b.rows[queryName] = append(b.rows[queryName], row)
}

func (b *packBuilder) newRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap string) QueryPackRow {
	b.counters[queryName]++
	id := fmt.Sprintf("%s.%04d", queryShortName(queryName), b.counters[queryName])
	return QueryPackRow{
		ID:                   id,
		Query:                queryName,
		EvidenceState:        state,
		EvidenceFamily:       family,
		SourceRef:            sourceRef,
		SourceConditionID:    conditionID,
		SourceConditionState: conditionState,
		ReasonCode:           reasonCode,
		EvidenceGap:          gap,
	}
}

func queryShortName(queryName string) string {
	return strings.TrimPrefix(queryName, "forensics-")
}

func mapSourceState(state string) string {
	switch state {
	case "", "pass":
		return RowStatePresent
	case "fail":
		return RowStateIssueObserved
	case RowStateCannotVerify:
		return RowStateCannotVerify
	case RowStateNotAssessed:
		return RowStateNotAssessed
	case RowStateMissingTelemetry:
		return RowStateMissingTelemetry
	case RowStateUnsupported:
		return RowStateUnsupported
	case RowStateNotIntegrated:
		return RowStateNotIntegrated
	case RowStateRetentionLimited:
		return RowStateRetentionLimited
	default:
		return RowStateCannotVerify
	}
}

func familyForEvent(eventType string) string {
	switch {
	case strings.Contains(eventType, "supersed"):
		return EvidenceFamilySupersession
	case strings.Contains(eventType, "task"):
		return EvidenceFamilyTask
	case strings.Contains(eventType, "command"):
		return EvidenceFamilyCommand
	case strings.Contains(eventType, "file"):
		return EvidenceFamilyFileMutations
	case strings.Contains(eventType, "test"):
		return EvidenceFamilyTest
	case strings.Contains(eventType, "redaction"):
		return EvidenceFamilyRedaction
	default:
		return EvidenceFamilyRunChain
	}
}

func familyForForensicCondition(id string) string {
	switch {
	case strings.Contains(id, "redaction"):
		return EvidenceFamilyRedaction
	case strings.Contains(id, "retention"), strings.Contains(id, "raw_reference"), strings.Contains(id, "critical_evidence"):
		return EvidenceFamilyRetention
	default:
		return EvidenceFamilyClaim
	}
}

func familyForAdapterCondition(id string) string {
	switch {
	case strings.Contains(id, "task"):
		return EvidenceFamilyTask
	case strings.Contains(id, "file"):
		return EvidenceFamilyFileMutations
	case strings.Contains(id, "test"):
		return EvidenceFamilyTest
	default:
		return EvidenceFamilyAdapterCapture
	}
}

func familyForVerifierState(id string) string {
	switch {
	case strings.Contains(id, "witness"):
		return EvidenceFamilyWitness
	case strings.Contains(id, "supersed"):
		return EvidenceFamilySupersession
	case strings.Contains(id, "task"):
		return EvidenceFamilyTask
	case strings.Contains(id, "command"):
		return EvidenceFamilyCommand
	case strings.Contains(id, "file"):
		return EvidenceFamilyFileMutations
	case strings.Contains(id, "test"):
		return EvidenceFamilyTest
	case strings.Contains(id, "redaction"):
		return EvidenceFamilyRedaction
	default:
		return EvidenceFamilyRunChain
	}
}

func safeToken(value string) string {
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "unknown"
	}
	return out
}

func sensitiveClasses() []string {
	return []string{
		"raw_command_arguments",
		"command_names",
		"executable_paths",
		"script_paths",
		"unsafe_test_identifiers",
		"stdout_stderr_bodies",
		"prompts",
		"source_snippets",
		"tool_payloads",
		"adapter_configuration",
		"gateway_evidence_refs",
		"credentials",
		"tokens",
		"authenticated_provider_urls",
		"raw_model_payloads",
		"raw_review_bodies",
		"unsafe_raw_reference_access_notes",
		"key_material",
	}
}
