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
	SHA256           string `json:"sha256,omitempty"`
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
	return forensicsBasicPack(loadPackInputs(runDir))
}

func forensicsBasicPack(inputs packInputs, err error) (QueryPackResult, error) {
	if err != nil {

		return QueryPackResult{}, err
	}
	return buildForensicsBasicPack(inputs), nil
}

func buildForensicsBasicPack(inputs packInputs) QueryPackResult {

	builder := newPackBuilder(inputs)
	if inputs.runErr != nil {
		builder.addMalformedRequiredInputRows()
		return builder.result()
	}
	builder.addTimelineRows()
	builder.addRedactionRows()
	builder.addCaptureRows()
	builder.addGapRows()
	builder.addUnverifiedClaimRows()
	builder.addSummaryRows()
	return builder.result()
}

func ExplainForensicsPack(result QueryPackResult) string {
	var lines []string
	for _, queryName := range queryOrder {

		rows := sortedQueryRows(result.QueryRows[queryName])
		for _, row := range rows {
			lines = append(lines, explainQueryRow(queryName, row))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func loadPackInputs(runDir string) (packInputs, error) {
	var run runArtifact
	runArtifact, err := readPackArtifact(filepath.Join(runDir, "run.json"), "run", "run", true, &run)
	if err != nil && runArtifact.Role == "" {

		return packInputs{}, err
	}
	inputs := packInputs{run: run, runArtifact: runArtifact, runErr: err}
	return loadOptionalPackInputs(runDir, inputs)
}

func sortedQueryRows(rows []QueryPackRow) []QueryPackRow {
	sorted := append([]QueryPackRow(nil), rows...)

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	return sorted
}

func explainQueryRow(queryName string, row QueryPackRow) string {

	parts := []string{queryName, row.ID, row.EvidenceState, row.EvidenceFamily}
	parts = append(parts, "source_ref="+row.SourceRef)
	parts = appendOptionalPart(parts, "source_condition_id", row.SourceConditionID)
	parts = appendOptionalPart(parts, "source_condition_state", row.SourceConditionState)
	if row.Reconstructable != nil {
		parts = append(parts, fmt.Sprintf("reconstructable=%t", *row.Reconstructable))
	}
	parts = appendOptionalPart(parts, "gap", row.EvidenceGap)
	return strings.Join(parts, " ")
}

func appendOptionalPart(parts []string, key, value string) []string {
	if value == "" {

		return parts
	}
	return append(parts, key+"="+value)
}

func loadOptionalPackInputs(runDir string, inputs packInputs) (packInputs, error) {
	var err error

	inputs, err = loadForensicInput(runDir, inputs)
	if err != nil {
		return packInputs{}, err
	}
	return loadAdapterInput(runDir, inputs)
}

func loadForensicInput(runDir string, inputs packInputs) (packInputs, error) {
	var forensic assessmentEnvelope
	artifact, present, err := readOptionalPackArtifact(filepath.Join(runDir, "forensic-retention.assessment-result.json"), "forensic_retention", "forensic_retention", false, &forensic)
	if err != nil && artifact.Role == "" {

		return packInputs{}, err
	}
	if present {

		inputs.forensicPresent = true
		inputs.forensicArtifact = &artifact
		inputs.forensic = forensic
		inputs.forensicErr = err
	}
	return inputs, nil
}

func loadAdapterInput(runDir string, inputs packInputs) (packInputs, error) {
	var adapter assessmentEnvelope
	artifact, present, err := readOptionalPackArtifact(filepath.Join(runDir, "adapter-capture.assessment-result.json"), "adapter_capture", "adapter_capture", false, &adapter)
	if err != nil && artifact.Role == "" {

		return packInputs{}, err
	}
	if present {

		inputs.adapterPresent = true
		inputs.adapterArtifact = &artifact
		inputs.adapter = adapter
		inputs.adapterErr = err
	}
	return inputs, nil
}

func readOptionalPackArtifact(path, role, redactedID string, required bool, target any) (QueryPackInputArtifact, bool, error) {
	if _, err := os.Stat(path); err != nil {

		return optionalArtifactStatResult(err, role, redactedID, required)
	}

	artifact, err := readPackArtifact(path, role, redactedID, required, target)
	return artifact, true, err
}

func optionalArtifactStatResult(err error, role, redactedID string, required bool) (QueryPackInputArtifact, bool, error) {
	if errors.Is(err, os.ErrNotExist) {

		return QueryPackInputArtifact{}, false, nil
	}
	return QueryPackInputArtifact{
		Role:             role,
		PathRedactedID:   redactedID,
		ArtifactRequired: required,
	}, true, err
}

func readPackArtifact(path, role, redactedID string, required bool, target any) (QueryPackInputArtifact, error) {
	artifact := QueryPackInputArtifact{
		Role:             role,
		PathRedactedID:   redactedID,
		ArtifactRequired: required,
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		return artifact, err
	}

	artifact.SHA256 = sha256Hex(payload)
	artifact.SchemaVersion = readArtifactSchemaVersion(payload)

	if err := json.Unmarshal(payload, target); err != nil {
		return artifact, err
	}
	return artifact, nil
}

func sha256Hex(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func readArtifactSchemaVersion(payload []byte) string {
	var envelope struct {
		SchemaVersion string `json:"schema_version"`
	}
	if err := json.Unmarshal(payload, &envelope); err == nil {

		return envelope.SchemaVersion
	}
	return ""
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

	return QueryPackResult{
		SchemaVersion:    QueryPackSchemaVersion,
		QueryPackID:      QueryPackForensicsBasic,
		QueryPackVersion: "v1",
		RunID:            b.inputs.run.RunID,
		RunNonce:         b.inputs.run.RunNonce,
		SourceBaseline:   b.inputs.run.SourceBaseline,
		InputArtifacts:   b.inputArtifacts(),
		QueryRows:        b.rows,
		OutputSafety:     b.outputSafety(),
	}
}

func (b *packBuilder) inputArtifacts() []QueryPackInputArtifact {
	artifacts := []QueryPackInputArtifact{b.inputs.runArtifact}
	artifacts = appendObservedArtifact(artifacts, b.inputs.forensicArtifact)
	artifacts = appendObservedArtifact(artifacts, b.inputs.adapterArtifact)

	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Role < artifacts[j].Role })
	return artifacts
}

func appendObservedArtifact(artifacts []QueryPackInputArtifact, artifact *QueryPackInputArtifact) []QueryPackInputArtifact {
	if artifact == nil {
		return artifacts
	}

	return append(artifacts, *artifact)
}

func (b *packBuilder) outputSafety() QueryPackOutputSafety {

	return QueryPackOutputSafety{
		RedactionPolicyDigest:          b.inputs.run.RedactionDigest,
		VerifiedAbsentSensitiveClasses: sensitiveClasses(),
	}
}

func (b *packBuilder) addTimelineRows() {
	b.addRunTimelineRows()
	b.addOptionalTimelineRows()
}

func (b *packBuilder) addRunTimelineRows() {
	if len(b.inputs.run.EventRefs) == 0 {

		b.addRow(QueryForensicsTimeline, RowStatePresent, EvidenceFamilyRunChain, "block_09.run.run_id", "", "", "run_timeline_available", "")
		return
	}
	for i, event := range b.inputs.run.EventRefs {
		family := familyForEvent(event.EventType)

		sourceRef := fmt.Sprintf("block_09.event.%s.e%04d", family, i+1)
		b.addRow(QueryForensicsTimeline, RowStatePresent, family, sourceRef, "", "", "timeline_event_present", "")
	}
}

func (b *packBuilder) addOptionalTimelineRows() {
	b.addOptionalTimelineRow(b.inputs.forensicPresent, b.inputs.forensicErr, EvidenceFamilyRetention, "block_18", "missing_optional_block_18_forensic_retention_result")
	b.addOptionalTimelineRow(b.inputs.adapterPresent, b.inputs.adapterErr, EvidenceFamilyAdapterCapture, "block_19", "missing_optional_block_19_adapter_capture_result")
}

func (b *packBuilder) addOptionalTimelineRow(present bool, inputErr error, family, block, missingReason string) {
	if !present {

		b.addRow(QueryForensicsTimeline, RowStateNotAssessed, family, block+".condition.missing", "", "", missingReason, family)
		return
	}
	if inputErr != nil {

		b.addRow(QueryForensicsTimeline, RowStateCannotVerify, EvidenceFamilyInputArtifact, block+".condition.malformed", "", "", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
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

		b.addMissingAdapterCaptureRow()
		return
	}
	if b.inputs.adapterErr != nil {

		b.addMalformedAdapterCaptureRow()
		return
	}
	for _, condition := range b.inputs.adapter.AdapterCaptureConditions {
		b.addAdapterCaptureConditionRow(condition)
	}
}

func (b *packBuilder) addMissingAdapterCaptureRow() {
	b.addRow(QueryForensicsCaptureDepth, RowStateCannotVerify, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "", "", "missing_block_19_adapter_capture_result", "adapter_capture")
}

func (b *packBuilder) addMalformedAdapterCaptureRow() {
	b.addRow(QueryForensicsCaptureDepth, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_19.condition.malformed", "", "", "unreadable_or_malformed_input_artifact", "input_artifact")
}

func (b *packBuilder) addAdapterCaptureConditionRow(condition assessmentCondition) {
	family := familyForAdapterCondition(condition.ID)

	row := b.rowFromCondition(QueryForensicsCaptureDepth, family, "block_19.condition."+safeToken(condition.ID), condition)
	b.rows[QueryForensicsCaptureDepth] = append(b.rows[QueryForensicsCaptureDepth], row)
}
func (b *packBuilder) addGapRows() {
	b.addVerifierGapRows()
	b.addForensicGapRows()
	b.addAdapterGapRows()
}

func (b *packBuilder) addVerifierGapRows() {
	for _, key := range sortedVerifierStateKeys(b.inputs.run.VerifierStates) {

		b.addVerifierGapRow(key, b.inputs.run.VerifierStates[key])
	}
}

func sortedVerifierStateKeys(states map[string]verifierState) []string {
	keys := make([]string, 0, len(states))
	for key := range states {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (b *packBuilder) addVerifierGapRow(key string, state verifierState) {
	rowState := mapSourceState(state.State)
	if rowState == RowStatePresent {

		return
	}
	family := familyForVerifierState(key)
	b.addRow(QueryForensicsGaps, rowState, family, "block_09.run."+safeToken(key), "", "", safeToken(key), family)
}

func (b *packBuilder) addForensicGapRows() {
	if !b.inputs.forensicPresent {

		b.addRow(QueryForensicsGaps, RowStateNotAssessed, EvidenceFamilyRetention, "block_18.condition.missing", "", "", "missing_optional_block_18_forensic_retention_result", "retention")
	} else if b.inputs.forensicErr != nil {

		b.addRow(QueryForensicsGaps, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_18.condition.malformed", "", "", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}

func (b *packBuilder) addAdapterGapRows() {
	if !b.inputs.adapterPresent {

		b.addRow(QueryForensicsGaps, RowStateNotAssessed, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "", "", "missing_optional_block_19_adapter_capture_result", "adapter_capture")
	} else if b.inputs.adapterErr != nil {

		b.addRow(QueryForensicsGaps, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_19.condition.malformed", "", "", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}

func (b *packBuilder) addUnverifiedClaimRows() {
	b.addUnverifiedClaimsFor(QueryForensicsRedactions)
	b.addUnverifiedClaimsFor(QueryForensicsCaptureDepth)
}

func (b *packBuilder) addUnverifiedClaimsFor(queryName string) {
	for _, row := range append([]QueryPackRow{}, b.rows[queryName]...) {
		if row.EvidenceState != RowStatePresent {

			b.addReferencedClaim(row)
		}
	}
}

func (b *packBuilder) addSummaryRows() {
	for _, queryName := range queryOrder {
		if queryName != QueryForensicsSummary {

			b.addSummaryRow(queryName)
		}
	}
}

func (b *packBuilder) addSummaryRow(queryName string) {
	related := b.relatedRows(queryName)
	if len(related) == 0 {

		return
	}
	row := b.newRow(QueryForensicsSummary, RowStatePresent, EvidenceFamilyClaim, "block_09.run.run_id", "", "", "query_group_index", "")
	row.RelatedRows = related
	b.rows[QueryForensicsSummary] = append(b.rows[QueryForensicsSummary], row)
}

func (b *packBuilder) relatedRows(queryName string) []string {
	var related []string
	for _, row := range b.rows[queryName] {

		related = append(related, row.ID)
	}
	return related
}

func (b *packBuilder) addReferencedClaim(source QueryPackRow) {

	row := b.newRow(QueryForensicsUnverifiedClaims, source.EvidenceState, EvidenceFamilyClaim, source.SourceRef, source.SourceConditionID, source.SourceConditionState, source.ReasonCode, source.EvidenceGap)
	row.RelatedRows = []string{source.ID}
	row.Reconstructable = source.Reconstructable
	b.rows[QueryForensicsUnverifiedClaims] = append(b.rows[QueryForensicsUnverifiedClaims], row)
}

func (b *packBuilder) rowFromCondition(queryName, family, sourceRef string, condition assessmentCondition) QueryPackRow {
	if criticalEvidenceRetentionCap(condition) {

		state := RowStateRetentionLimited
		row := b.newRow(queryName, state, family, sourceRef, condition.ID, condition.State, "digest_only_not_reconstructable", EvidenceFamilyRetention)
		row.Reconstructable = falsePointer()
		return row
	}
	state := mapSourceState(condition.State)
	gap := gapForConditionState(state, family)
	row := b.newRow(queryName, state, family, sourceRef, condition.ID, condition.State, condition.ReasonCode, gap)
	row.Reconstructable = reconstructableForCondition(condition)
	return row
}

func criticalEvidenceRetentionCap(condition assessmentCondition) bool {

	return condition.ID == "critical_evidence_reconstructable" &&
		(condition.CappedToRetentionMode != "" || condition.ReasonCode == "critical_evidence_digest_only")
}

func gapForConditionState(state, family string) string {
	if state == RowStatePresent || state == RowStateIssueObserved {

		return ""
	}
	return family
}

func reconstructableForCondition(condition assessmentCondition) *bool {
	if condition.State == RowStateRetentionLimited || condition.CappedToRetentionMode != "" {

		return falsePointer()
	}
	return nil
}

func falsePointer() *bool {
	falseValue := false
	return &falseValue
}
func (b *packBuilder) addRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap string) {

	row := b.newRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap)
	b.rows[queryName] = append(b.rows[queryName], row)
}

func (b *packBuilder) newRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap string) QueryPackRow {

	id := b.nextRowID(queryName)
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

func (b *packBuilder) nextRowID(queryName string) string {
	b.counters[queryName]++
	return fmt.Sprintf("%s.%04d", queryShortName(queryName), b.counters[queryName])
}

func queryShortName(queryName string) string {
	return strings.TrimPrefix(queryName, "forensics-")
}

var sourceStateRows = map[string]string{
	"pass":                   RowStatePresent,
	"":                       RowStateCannotVerify,
	"fail":                   RowStateIssueObserved,
	RowStateCannotVerify:     RowStateCannotVerify,
	RowStateNotAssessed:      RowStateNotAssessed,
	RowStateMissingTelemetry: RowStateMissingTelemetry,
	RowStateUnsupported:      RowStateUnsupported,
	RowStateNotIntegrated:    RowStateNotIntegrated,
	RowStateRetentionLimited: RowStateRetentionLimited,
}

func mapSourceState(state string) string {

	if mapped, ok := sourceStateRows[state]; ok {
		return mapped
	}
	return RowStateCannotVerify
}

func familyForEvent(eventType string) string {
	return firstMatchingFamily(eventType, eventFamilyRules, EvidenceFamilyRunChain)
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
	if strings.Contains(id, "task") {

		return EvidenceFamilyTask
	}
	return nonTaskAdapterFamily(id)
}

func nonTaskAdapterFamily(id string) string {
	switch {
	case strings.Contains(id, "file"):

		return EvidenceFamilyFileMutations
	case strings.Contains(id, "test"):

		return EvidenceFamilyTest
	default:
		return EvidenceFamilyAdapterCapture
	}
}

func familyForVerifierState(id string) string {
	return firstMatchingFamily(id, verifierFamilyRules, EvidenceFamilyRunChain)
}

type familyRule struct {
	token  string
	family string
}

var eventFamilyRules = []familyRule{
	{"supersed", EvidenceFamilySupersession},
	{"task", EvidenceFamilyTask},
	{"command", EvidenceFamilyCommand},
	{"file", EvidenceFamilyFileMutations},
	{"test", EvidenceFamilyTest},
	{"redaction", EvidenceFamilyRedaction},
}

var verifierFamilyRules = []familyRule{
	{"witness", EvidenceFamilyWitness},
	{"supersed", EvidenceFamilySupersession},
	{"task", EvidenceFamilyTask},
	{"command", EvidenceFamilyCommand},
	{"file", EvidenceFamilyFileMutations},
	{"test", EvidenceFamilyTest},
	{"redaction", EvidenceFamilyRedaction},
}

func firstMatchingFamily(value string, rules []familyRule, fallback string) string {
	for _, rule := range rules {
		if strings.Contains(value, rule.token) {

			return rule.family
		}
	}
	return fallback
}

func safeToken(value string) string {
	sanitized := sanitizeToken(value)
	if sanitized == "" {

		return "unknown"
	}
	return sanitized
}

const safeTokenAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func sanitizeToken(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if isSafeTokenChar(r) {

			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func isSafeTokenChar(r rune) bool {
	return strings.ContainsRune(safeTokenAlphabet, r)
}

var verifiedAbsentSensitiveClasses = []string{
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

func sensitiveClasses() []string {

	return append([]string(nil), verifiedAbsentSensitiveClasses...)
}
