package posture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

const (
	SchemaVersion               = "block21-cross-repo-posture-export-v1"
	ProfileID                   = "cross-repo-evidence-posture-v1"
	ProfileVer                  = "v1"
	SelectionSchemaVersion      = "block21-cross-repo-selection-v1"
	DigestManifestSchemaVersion = "block21-artifact-digest-manifest-v1"
	SignalManifestSchemaVersion = "block21-posture-signal-manifest-v1"

	GroupingRepoWindow          = "repo_window_v1"
	GroupingTeamServiceWindow   = "team_service_window_v1"
	GroupingHarnessChangeWindow = "harness_change_window_v1"
)

var metricCatalog = []metricDef{
	{"missing_telemetry_rows", "v1", "row_state"},
	{"not_assessed_rows", "v1", "row_state"},
	{"cannot_verify_rows", "v1", "row_state"},
	{"unsupported_observer_rows", "v1", "row_or_signal"},
	{"not_integrated_rows", "v1", "row_state"},
	{"retention_limited_rows", "v1", "row_state"},
	{"local_only_evidence_rows", "v1", "posture_signal"},
	{"ci_witnessed_evidence_rows", "v1", "posture_signal"},
	{"external_witnessed_evidence_rows", "v1", "posture_signal"},
	{"issue_observed_rows", "v1", "row_state"},
	{"override_rows", "v1", "posture_signal"},
	{"late_attach_rows", "v1", "posture_signal"},
	{"contract_change_rows", "v1", "posture_signal"},
}

var rowStateMetrics = map[string]string{
	"missing_telemetry_rows": query.RowStateMissingTelemetry,
	"not_assessed_rows":      query.RowStateNotAssessed,
	"cannot_verify_rows":     query.RowStateCannotVerify,
	"not_integrated_rows":    query.RowStateNotIntegrated,
	"retention_limited_rows": query.RowStateRetentionLimited,
	"issue_observed_rows":    query.RowStateIssueObserved,
}

var safeLabelPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

type SelectionManifest struct {
	SchemaVersion           string             `json:"schema_version"`
	ProfileID               string             `json:"profile_id"`
	ProfileVersion          string             `json:"profile_version,omitempty"`
	GroupingSetID           string             `json:"grouping_set_id"`
	FreshnessBoundary       string             `json:"freshness_boundary"`
	DimensionExposurePolicy []string           `json:"dimension_exposure_policy"`
	CurrentWindow           string             `json:"current_window"`
	PreviousWindow          string             `json:"previous_window"`
	Repositories            []RepositoryWindow `json:"repositories"`
	Handoff                 map[string]string  `json:"handoff,omitempty"`
}

type RepositoryWindow struct {
	InputID                string `json:"input_id"`
	Repo                   string `json:"repo"`
	Team                   string `json:"team"`
	Service                string `json:"service"`
	Harness                string `json:"harness"`
	ChangeType             string `json:"change_type"`
	TimeWindow             string `json:"time_window"`
	InputObservedAt        string `json:"input_observed_at"`
	QueryPackResult        string `json:"query_pack_result"`
	ArtifactDigestManifest string `json:"artifact_digest_manifest"`
	PostureSignalManifest  string `json:"posture_signal_manifest,omitempty"`
}

type DigestManifest struct {
	SchemaVersion string           `json:"schema_version"`
	Artifacts     []DigestArtifact `json:"artifacts"`
}

type DigestArtifact struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type SignalManifest struct {
	SchemaVersion string          `json:"schema_version"`
	Signals       []PostureSignal `json:"signals"`
}

type PostureSignal struct {
	RowRef               string `json:"row_ref"`
	WitnessScope         string `json:"witness_scope,omitempty"`
	ObserverState        string `json:"observer_state,omitempty"`
	OverrideMarker       string `json:"override_marker,omitempty"`
	LateAttachMarker     string `json:"late_attach_marker,omitempty"`
	ContractChangeMarker string `json:"contract_change_marker,omitempty"`
}

type ExportResult struct {
	SchemaVersion        string            `json:"schema_version"`
	ExportProfileID      string            `json:"export_profile_id"`
	ExportProfileVersion string            `json:"export_profile_version"`
	ExportID             string            `json:"export_id"`
	Producer             string            `json:"producer"`
	GeneratedAt          string            `json:"generated_at"`
	GroupingSetID        string            `json:"grouping_set_id"`
	ActiveGroupingKeys   []string          `json:"active_grouping_keys"`
	InputSelection       []InputSelection  `json:"input_selection"`
	MetricRows           []MetricRow       `json:"metric_rows"`
	MovementRows         []MovementRow     `json:"movement_rows"`
	MovementSummary      MovementSummary   `json:"movement_summary"`
	RefusalRows          []RefusalRow      `json:"refusal_rows"`
	Handoff              map[string]string `json:"handoff"`
	OutputSafety         OutputSafety      `json:"output_safety"`
}

type InputSelection struct {
	InputID         string `json:"input_id"`
	Repository      string `json:"repo"`
	TimeWindow      string `json:"time_window"`
	PathRedactedID  string `json:"path_redacted_id"`
	SHA256          string `json:"sha256,omitempty"`
	InputTrustState string `json:"input_trust_state"`
}

type MetricRow struct {
	ID                      string            `json:"id"`
	MetricID                string            `json:"metric_id"`
	MetricVersion           string            `json:"metric_version"`
	Numerator               int               `json:"numerator"`
	Denominator             int               `json:"denominator"`
	Unit                    string            `json:"unit"`
	TimeWindow              string            `json:"time_window"`
	Dimensions              map[string]string `json:"dimensions"`
	DimensionKey            string            `json:"dimension_key"`
	SourceInputRefs         []string          `json:"source_input_refs"`
	SourceArtifactDigestSet string            `json:"source_artifact_digest_set_hash"`
	SourceFieldState        string            `json:"source_field_state"`
	NotAssessedCount        int               `json:"not_assessed_count"`
	InputTrustStateSummary  map[string]int    `json:"input_trust_state_summary"`
}

type MovementRow struct {
	ID                   string `json:"id"`
	MetricID             string `json:"metric_id"`
	MetricVersion        string `json:"metric_version"`
	DimensionKey         string `json:"dimension_key"`
	CurrentMetricRowRef  string `json:"current_metric_row_ref,omitempty"`
	PreviousMetricRowRef string `json:"previous_metric_row_ref,omitempty"`
	CurrentValue         int    `json:"current_value"`
	PreviousValue        int    `json:"previous_value"`
	Delta                int    `json:"delta"`
	ComparisonBasis      string `json:"comparison_basis"`
	Comparable           bool   `json:"comparable"`
	NonComparableReason  string `json:"non_comparable_reason,omitempty"`
}

type MovementSummary struct {
	ComparableCount     int            `json:"comparable_count"`
	NonComparableCount  int            `json:"non_comparable_count"`
	NonComparableReason map[string]int `json:"non_comparable_reason_counts"`
}

type RefusalRow struct {
	ID              string `json:"id"`
	InputID         string `json:"input_id"`
	TimeWindow      string `json:"time_window,omitempty"`
	RefusalReason   string `json:"refusal_reason"`
	InputTrustState string `json:"input_trust_state"`
}

type OutputSafety struct {
	VerifiedAbsentSensitiveClasses []string `json:"verified_absent_sensitive_classes"`
}

type metricDef struct {
	id      string
	version string
	source  string
}

type aggregateGroup struct {
	dimensions   map[string]string
	dimensionKey string
	window       string
	rows         []query.QueryPackRow
	inputRefs    []string
	digests      []string
	trustStates  map[string]int
	signals      map[string]PostureSignal
}

func Build(selectionPath string, now time.Time) (ExportResult, error) {
	input, err := prepareBuildInput(selectionPath, now)
	if err != nil {
		return ExportResult{}, err
	}
	inputs, refusals, groups := ingestRepositories(input.selection, input.activeKeys, input.cutoff, input.hasCutoff)
	metricRows := buildMetrics(groups)
	movementRows, summary := buildMovements(metricRows, input.selection.CurrentWindow, input.selection.PreviousWindow)
	return buildExportResult(selectionPath, now, input, inputs, metricRows, movementRows, summary, refusals), nil
}

type buildInput struct {
	selection  SelectionManifest
	activeKeys []string
	cutoff     time.Time
	hasCutoff  bool
	handoff    map[string]string
}

func prepareBuildInput(selectionPath string, now time.Time) (buildInput, error) {
	selection, err := readSelection(selectionPath)
	if err != nil {
		return buildInput{}, err
	}
	if err := validateSelection(selection); err != nil {
		return buildInput{}, err
	}
	activeKeys := groupingKeys(selection.GroupingSetID)
	if len(activeKeys) == 0 {
		return buildInput{}, fmt.Errorf("unsupported grouping set")
	}
	cutoff, hasCutoff, err := parseFreshnessBoundary(selection.FreshnessBoundary, now)
	if err != nil {
		return buildInput{}, err
	}
	handoff := selection.Handoff
	if handoff == nil {
		handoff = map[string]string{}
	}
	if !safeHandoff(handoff) {
		return buildInput{}, fmt.Errorf("unsafe handoff")
	}
	return buildInput{selection: selection, activeKeys: activeKeys, cutoff: cutoff, hasCutoff: hasCutoff, handoff: handoff}, nil
}

func ingestRepositories(selection SelectionManifest, activeKeys []string, cutoff time.Time, hasCutoff bool) ([]InputSelection, []RefusalRow, map[string]*aggregateGroup) {
	inputs := []InputSelection{}
	refusals := []RefusalRow{}
	groups := map[string]*aggregateGroup{}
	refusalCounter := 0
	for _, repo := range selection.Repositories {
		ingested := ingestRepository(repo, cutoff, hasCutoff)
		if !ingested.trusted {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, ingested.refusalReason, ingested.inputTrustState))
			if ingested.recordSelection {
				inputs = append(inputs, inputSelection(repo, ingested.digest, ingested.inputTrustState))
			}
			continue
		}
		inputs = append(inputs, inputSelection(repo, ingested.digest, "trusted_input"))
		addTrustedRepositoryGroup(groups, repo, activeKeys, ingested.result, ingested.signals, ingested.digest)
	}
	return inputs, refusals, groups
}

type repositoryIngest struct {
	trusted         bool
	recordSelection bool
	digest          string
	refusalReason   string
	inputTrustState string
	result          query.QueryPackResult
	signals         map[string]PostureSignal
}

func ingestRepository(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) repositoryIngest {
	checked := ingestRepositoryChecks(repo, cutoff, hasCutoff)
	if !checked.trusted {
		return checked
	}
	return ingestRepositoryArtifacts(repo)
}

func ingestRepositoryChecks(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) repositoryIngest {
	if err := validateRepoLabels(repo); err != nil {
		return refusedInput("unsafe_label", "cannot_verify_input", "", false)
	}
	if err := validateInputPaths(repo); err != nil {
		return refusedInput("malformed_input", "cannot_verify_input", "", false)
	}
	if hasCutoff && isStale(repo.InputObservedAt, cutoff) {
		return refusedInput("stale_input", "stale_input", "", true)
	}
	return repositoryIngest{trusted: true}
}

func ingestRepositoryArtifacts(repo RepositoryWindow) repositoryIngest {
	digest, err := verifyDigestManifest(repo.ArtifactDigestManifest, repo.QueryPackResult)
	if err != nil {
		return refusedInput(reasonForDigestErr(err), trustForDigestErr(err), "", true)
	}
	result, err := readQueryPack(repo.QueryPackResult)
	if err != nil {
		return refusedInput("malformed_input", "cannot_verify_input", digest, true)
	}
	if !isSupportedQueryPack(result) {
		return refusedInput("malformed_input", "cannot_verify_input", digest, true)
	}
	signals, err := readSignals(repo.PostureSignalManifest)
	if err != nil {
		return refusedInput("malformed_input", "cannot_verify_input", digest, true)
	}
	return repositoryIngest{trusted: true, digest: digest, result: result, signals: signals}
}

func isSupportedQueryPack(result query.QueryPackResult) bool {
	return result.SchemaVersion == query.QueryPackSchemaVersion && result.QueryPackID == query.QueryPackForensicsBasic
}

func refusedInput(reason, trustState, digest string, recordSelection bool) repositoryIngest {
	return repositoryIngest{refusalReason: reason, inputTrustState: trustState, digest: digest, recordSelection: recordSelection}
}

func addTrustedRepositoryGroup(groups map[string]*aggregateGroup, repo RepositoryWindow, activeKeys []string, result query.QueryPackResult, signals map[string]PostureSignal, digest string) {
	key, dims := dimensionKey(repo, activeKeys)
	groupKey := repo.TimeWindow + "|" + key
	group := groups[groupKey]
	if group == nil {
		group = &aggregateGroup{
			dimensions:   dims,
			dimensionKey: key,
			window:       repo.TimeWindow,
			trustStates:  map[string]int{},
			signals:      map[string]PostureSignal{},
		}
		groups[groupKey] = group
	}
	group.rows = append(group.rows, flattenRows(result)...)
	group.inputRefs = append(group.inputRefs, repo.InputID)
	group.digests = append(group.digests, digest)
	group.trustStates["trusted_input"]++
	for rowRef, signal := range signals {
		group.signals[rowRef] = signal
	}
}

func buildExportResult(selectionPath string, now time.Time, input buildInput, inputs []InputSelection, metricRows []MetricRow, movementRows []MovementRow, summary MovementSummary, refusals []RefusalRow) ExportResult {
	return ExportResult{
		SchemaVersion:        SchemaVersion,
		ExportProfileID:      ProfileID,
		ExportProfileVersion: ProfileVer,
		ExportID:             deterministicExportID(selectionPath, metricRows, refusals),
		Producer:             "sdp-trace",
		GeneratedAt:          now.UTC().Format(time.RFC3339),
		GroupingSetID:        input.selection.GroupingSetID,
		ActiveGroupingKeys:   input.activeKeys,
		InputSelection:       inputs,
		MetricRows:           metricRows,
		MovementRows:         movementRows,
		MovementSummary:      summary,
		RefusalRows:          refusals,
		Handoff:              input.handoff,
		OutputSafety: OutputSafety{
			VerifiedAbsentSensitiveClasses: SensitiveClasses(),
		},
	}
}

func Explain(result ExportResult) (string, error) {
	var lines []string
	lines = append(lines, "schema_version="+result.SchemaVersion)
	lines = append(lines, "export_profile_id="+result.ExportProfileID)
	lines = append(lines, "grouping_set_id="+result.GroupingSetID)
	lines = append(lines, fmt.Sprintf("movement_summary comparable=%d non_comparable=%d", result.MovementSummary.ComparableCount, result.MovementSummary.NonComparableCount))
	for _, row := range result.RefusalRows {
		lines = append(lines, fmt.Sprintf("refusal %s input=%s reason=%s state=%s", row.ID, row.InputID, row.RefusalReason, row.InputTrustState))
	}
	for _, row := range result.MetricRows {
		lines = append(lines, fmt.Sprintf("metric %s %s numerator=%d denominator=%d window=%s dimension_key=%s not_assessed_count=%d", row.ID, row.MetricID, row.Numerator, row.Denominator, row.TimeWindow, row.DimensionKey, row.NotAssessedCount))
	}
	for _, row := range result.MovementRows {
		lines = append(lines, fmt.Sprintf("movement %s %s current=%d previous=%d delta=%d comparable=%t reason=%s", row.ID, row.MetricID, row.CurrentValue, row.PreviousValue, row.Delta, row.Comparable, row.NonComparableReason))
	}
	for _, class := range result.OutputSafety.VerifiedAbsentSensitiveClasses {
		lines = append(lines, "output_safety absent="+class)
	}
	rendered := strings.Join(lines, "\n") + "\n"
	if unsafeOutput(rendered) {
		return "", fmt.Errorf("output_safety_violation")
	}
	return rendered, nil
}

func ValidateExportResult(result ExportResult) error {
	if err := validateExportHeader(result); err != nil {
		return err
	}
	if err := validateExportCollections(result); err != nil {
		return err
	}
	return validateExportRows(result)
}

func validateExportHeader(result ExportResult) error {
	if result.SchemaVersion != SchemaVersion ||
		result.ExportProfileID != ProfileID ||
		result.ExportProfileVersion != ProfileVer {
		return fmt.Errorf("unsupported posture export")
	}
	if result.ExportID == "" || result.Producer != "sdp-trace" || result.GeneratedAt == "" {
		return fmt.Errorf("malformed posture export")
	}
	if _, err := time.Parse(time.RFC3339, result.GeneratedAt); err != nil {
		return fmt.Errorf("malformed posture export generated_at")
	}
	return nil
}

func validateExportRows(result ExportResult) error {
	return firstError(
		validateInputSelectionRows(result.InputSelection),
		validateMetricRows(result.MetricRows),
		validateMovementRows(result.MovementRows),
		validateMovementSummary(result.MovementSummary),
		validateRefusalRows(result.RefusalRows),
	)
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func validateInputSelectionRows(rows []InputSelection) error {
	for _, row := range rows {
		if err := validateInputSelectionRow(row); err != nil {
			return err
		}
	}
	return nil
}

func validateMetricRows(rows []MetricRow) error {
	for _, row := range rows {
		if err := validateMetricRow(row); err != nil {
			return err
		}
	}
	return nil
}

func validateMovementRows(rows []MovementRow) error {
	for _, row := range rows {
		if err := validateMovementRow(row); err != nil {
			return err
		}
	}
	return nil
}

func validateRefusalRows(rows []RefusalRow) error {
	for _, row := range rows {
		if err := validateRefusalRow(row); err != nil {
			return err
		}
	}
	return nil
}

func validateExportCollections(result ExportResult) error {
	if !validExportGrouping(result.GroupingSetID, result.ActiveGroupingKeys) {
		return fmt.Errorf("malformed posture export grouping")
	}
	if !hasRequiredCollections(result) {
		return fmt.Errorf("malformed posture export missing required collection")
	}
	if !hasOutputSafety(result.OutputSafety.VerifiedAbsentSensitiveClasses) {
		return fmt.Errorf("malformed posture export output_safety")
	}
	return nil
}

func validExportGrouping(groupingSet string, keys []string) bool {
	return len(groupingKeys(groupingSet)) > 0 && len(keys) >= 2
}

func hasRequiredCollections(result ExportResult) bool {
	for _, present := range []bool{
		result.InputSelection != nil,
		result.MetricRows != nil,
		result.MovementRows != nil,
		result.RefusalRows != nil,
		result.Handoff != nil,
		result.MovementSummary.NonComparableReason != nil,
	} {
		if !present {
			return false
		}
	}
	return true
}

func hasOutputSafety(classes []string) bool {
	return len(classes) > 0
}

func validateInputSelectionRow(item InputSelection) error {
	if !safeLabel(item.InputID) || !safeLabel(item.Repository) || !safeLabel(item.TimeWindow) ||
		item.PathRedactedID == "" || !validInputTrustState(item.InputTrustState) {
		return fmt.Errorf("malformed posture export input_selection")
	}
	return nil
}

func validateMetricRow(row MetricRow) error {
	if err := validateMetricRowShape(row); err != nil {
		return err
	}
	if err := validateMetricDimensions(row.Dimensions); err != nil {
		return err
	}
	return validateMetricTrustSummary(row.InputTrustStateSummary)
}

func validateMetricRowShape(row MetricRow) error {
	if malformedMetricIdentity(row) || malformedMetricCounts(row) || malformedMetricSource(row) {
		return fmt.Errorf("malformed posture export metric_row")
	}
	return nil
}

func malformedMetricIdentity(row MetricRow) bool {
	return row.ID == "" || !validMetricID(row.MetricID) || row.MetricVersion != ProfileVer || !safeLabel(row.TimeWindow)
}

func malformedMetricCounts(row MetricRow) bool {
	return row.Numerator < 0 || row.Denominator < 0 || row.Unit != "rows" || row.NotAssessedCount < 0
}

func malformedMetricSource(row MetricRow) bool {
	return missingMetricSourceRefs(row) || missingMetricTrustSource(row)
}

func missingMetricSourceRefs(row MetricRow) bool {
	return row.Dimensions == nil || row.DimensionKey == "" || row.SourceInputRefs == nil || row.SourceArtifactDigestSet == ""
}

func missingMetricTrustSource(row MetricRow) bool {
	return !validSourceFieldState(row.SourceFieldState) || row.InputTrustStateSummary == nil
}

func validateMetricDimensions(dimensions map[string]string) error {
	for key, value := range dimensions {
		if !validDimensionName(key) || !safeLabel(value) {
			return fmt.Errorf("malformed posture export metric_row dimensions")
		}
	}
	return nil
}

func validateMetricTrustSummary(summary map[string]int) error {
	for state, count := range summary {
		if !validInputTrustState(state) || count < 0 {
			return fmt.Errorf("malformed posture export input_trust_state_summary")
		}
	}
	return nil
}

func validateMovementRow(row MovementRow) error {
	if malformedMovementRow(row) {
		return fmt.Errorf("malformed posture export movement_row")
	}
	return nil
}

func malformedMovementRow(row MovementRow) bool {
	return malformedMovementIdentity(row) || malformedMovementValues(row) || malformedMovementComparison(row)
}

func malformedMovementIdentity(row MovementRow) bool {
	return row.ID == "" || !validMetricID(row.MetricID) || row.MetricVersion != ProfileVer || row.DimensionKey == ""
}

func malformedMovementValues(row MovementRow) bool {
	return row.CurrentValue < 0 || row.PreviousValue < 0
}

func malformedMovementComparison(row MovementRow) bool {
	return !validComparisonBasis(row.ComparisonBasis) || (!row.Comparable && row.NonComparableReason != "non_comparable_missing_window")
}

func validateMovementSummary(summary MovementSummary) error {
	if summary.ComparableCount < 0 || summary.NonComparableCount < 0 {
		return fmt.Errorf("malformed posture export movement_summary")
	}
	for reason, count := range summary.NonComparableReason {
		if reason != "non_comparable_missing_window" || count < 0 {
			return fmt.Errorf("malformed posture export movement_summary")
		}
	}
	return nil
}

func validateRefusalRow(row RefusalRow) error {
	if row.ID == "" || !safeLabel(row.InputID) || !validRefusalReason(row.RefusalReason) ||
		!validInputTrustState(row.InputTrustState) ||
		(row.TimeWindow != "" && !safeLabel(row.TimeWindow)) {
		return fmt.Errorf("malformed posture export refusal_row")
	}
	return nil
}

func SensitiveClasses() []string {
	return []string{
		"raw_command_args",
		"command_name_or_path",
		"unsafe_test_identifier",
		"stdout_stderr_body",
		"prompt_body",
		"source_snippet",
		"tool_payload",
		"adapter_configuration",
		"gateway_evidence_ref",
		"credential_or_token",
		"authenticated_provider_url",
		"model_request_response_payload",
		"raw_review_body",
		"unsafe_raw_reference_note",
		"private_filesystem_path",
		"unsafe_personal_identifier",
		"unsafe_label",
		"raw_digest_manifest_path",
		"free_text_exception_or_refusal_reason",
	}
}

func validMetricID(value string) bool {
	for _, item := range metricCatalog {
		if item.id == value {
			return true
		}
	}
	return false
}

func validDimensionName(value string) bool {
	switch value {
	case "repo", "team", "service", "harness", "change_type", "time_window":
		return true
	default:
		return false
	}
}

func validInputTrustState(value string) bool {
	switch value {
	case "trusted_input", "stale_input", "untrusted_input", "cannot_verify_input", "not_assessed_input":
		return true
	default:
		return false
	}
}

func validSourceFieldState(value string) bool {
	switch value {
	case "present", "not_assessed", "cannot_verify", "unsupported":
		return true
	default:
		return false
	}
}

func validComparisonBasis(value string) bool {
	switch value {
	case "same_profile_metric_dimension_window", "non_comparable_missing_window":
		return true
	default:
		return false
	}
}

func validRefusalReason(value string) bool {
	switch value {
	case "stale_input", "malformed_input", "untrusted_input_digest_mismatch", "unsafe_label",
		"unsupported_input", "missing_required_input", "missing_optional_input",
		"non_comparable_metric_version", "non_comparable_dimension_key",
		"non_comparable_denominator_basis", "non_comparable_input_trust_rule",
		"non_comparable_missing_window", "output_safety_violation":
		return true
	default:
		return false
	}
}

func readSelection(path string) (SelectionManifest, error) {
	var selection SelectionManifest
	data, err := os.ReadFile(path)
	if err != nil {
		return selection, err
	}
	return selection, json.Unmarshal(data, &selection)
}

func validateSelection(selection SelectionManifest) error {
	if selection.SchemaVersion != SelectionSchemaVersion {
		return fmt.Errorf("unsupported selection schema")
	}
	if selection.ProfileID != ProfileID {
		return fmt.Errorf("unsupported profile")
	}
	if selection.ProfileVersion != "" && selection.ProfileVersion != ProfileVer {
		return fmt.Errorf("unsupported profile version")
	}
	if len(groupingKeys(selection.GroupingSetID)) == 0 {
		return fmt.Errorf("unsupported grouping set")
	}
	if !groupingAllowedByExposure(selection.GroupingSetID, selection.DimensionExposurePolicy) {
		return fmt.Errorf("dimension exposure policy excludes grouping key")
	}
	if len(selection.Repositories) == 0 {
		return fmt.Errorf("empty selection")
	}
	return nil
}

func readQueryPack(path string) (query.QueryPackResult, error) {
	var result query.QueryPackResult
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}
	return result, json.Unmarshal(data, &result)
}

func readSignals(path string) (map[string]PostureSignal, error) {
	if strings.TrimSpace(path) == "" {
		return map[string]PostureSignal{}, nil
	}
	var manifest SignalManifest
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if manifest.SchemaVersion != SignalManifestSchemaVersion {
		return nil, fmt.Errorf("unsupported signal manifest schema")
	}
	out := map[string]PostureSignal{}
	for _, signal := range manifest.Signals {
		if unsafeOutput(signal.RowRef + signal.WitnessScope + signal.ObserverState + signal.OverrideMarker + signal.LateAttachMarker + signal.ContractChangeMarker) {
			return nil, fmt.Errorf("unsafe signal")
		}
		out[signal.RowRef] = signal
	}
	return out, nil
}

func verifyDigestManifest(manifestPath, queryPackPath string) (string, error) {
	manifest, err := readDigestManifest(manifestPath)
	if err != nil {
		return "", err
	}

	expected, err := digestForQueryPackFromManifest(manifest, queryPackPath)
	if err != nil {
		return "", err
	}

	actual, err := fileSHA256Hex(queryPackPath)
	if err != nil {
		return "", err
	}

	if err := checkDigestMatch(expected, actual); err != nil {
		return "", err
	}
	return actual, nil
}

func checkDigestMatch(expected, actual string) error {
	if expected != actual {
		return errDigestMismatch
	}
	return nil
}

func readDigestManifest(path string) (DigestManifest, error) {
	var manifest DigestManifest
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, err
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest, err
	}
	if manifest.SchemaVersion != DigestManifestSchemaVersion {
		return manifest, fmt.Errorf("unsupported digest manifest schema")
	}

	return manifest, nil
}

func digestForQueryPackFromManifest(manifest DigestManifest, queryPackPath string) (string, error) {
	filename := filepathBase(queryPackPath)
	for _, artifact := range manifest.Artifacts {
		if artifact.Role != "query_pack_result" {
			continue
		}
		if !digestArtifactMatchesPath(artifact.Path, filename) {
			return "", errUnsafePath
		}
		return artifact.SHA256, nil
	}
	return "", errMissingRequired
}

func digestArtifactMatchesPath(artifactPath, filename string) bool {
	return !unsafeSelectionPath(artifactPath) && artifactPath == filename
}

func fileSHA256Hex(path string) (string, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

var (
	errDigestMismatch  = errors.New("digest mismatch")
	errMissingRequired = errors.New("missing required input")
	errUnsafePath      = errors.New("unsafe path")
)

func reasonForDigestErr(err error) string {
	switch {
	case errors.Is(err, errDigestMismatch):
		return "untrusted_input_digest_mismatch"
	case errors.Is(err, errMissingRequired):
		return "missing_required_input"
	case errors.Is(err, errUnsafePath):
		return "malformed_input"
	default:
		return "malformed_input"
	}
}

func trustForDigestErr(err error) string {
	if errors.Is(err, errDigestMismatch) {
		return "untrusted_input"
	}
	return "cannot_verify_input"
}

func flattenRows(result query.QueryPackResult) []query.QueryPackRow {
	var rows []query.QueryPackRow
	keys := make([]string, 0, len(result.QueryRows))
	for key := range result.QueryRows {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		rows = append(rows, result.QueryRows[key]...)
	}
	return rows
}

func buildMetrics(groups map[string]*aggregateGroup) []MetricRow {
	groupKeys := make([]string, 0, len(groups))
	for key := range groups {
		groupKeys = append(groupKeys, key)
	}
	sort.Strings(groupKeys)
	var rows []MetricRow
	counter := 0
	for _, groupKey := range groupKeys {
		group := groups[groupKey]
		for _, def := range metricCatalog {
			counter++
			row := metricForGroup(counter, def, group)
			rows = append(rows, row)
		}
	}
	return rows
}

func metricForGroup(counter int, def metricDef, group *aggregateGroup) MetricRow {
	numerator := 0
	notAssessed := 0
	sourceFieldState := "present"
	for _, row := range group.rows {
		signal, hasSignal := group.signals[row.ID]
		if metricMatches(def.id, row, signal, hasSignal) {
			numerator++
		}
		if metricNotAssessed(def, row, hasSignal) {
			notAssessed++
		}
		if def.source == "posture_signal" && !hasSignal {
			sourceFieldState = "not_assessed"
		}
	}
	sourceRefs := append([]string(nil), group.inputRefs...)
	sort.Strings(sourceRefs)
	digests := append([]string(nil), group.digests...)
	sort.Strings(digests)
	return MetricRow{
		ID:                      fmt.Sprintf("metric.%04d", counter),
		MetricID:                def.id,
		MetricVersion:           def.version,
		Numerator:               numerator,
		Denominator:             len(group.rows),
		Unit:                    "rows",
		TimeWindow:              group.window,
		Dimensions:              group.dimensions,
		DimensionKey:            group.dimensionKey,
		SourceInputRefs:         sourceRefs,
		SourceArtifactDigestSet: digestSetHash(digests),
		SourceFieldState:        sourceFieldState,
		NotAssessedCount:        notAssessed,
		InputTrustStateSummary:  copyTrust(group.trustStates),
	}
}

func metricMatches(metricID string, row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {
	if expectedState, ok := rowStateMetrics[metricID]; ok {
		return row.EvidenceState == expectedState
	}
	if metricID == "unsupported_observer_rows" {
		return row.EvidenceState == query.RowStateUnsupported || (hasSignal && signal.ObserverState == "unsupported")
	}
	if !hasSignal {
		return false
	}
	return signalMetricMatches(metricID, signal)
}

func signalMetricMatches(metricID string, signal PostureSignal) bool {
	switch metricID {
	case "local_only_evidence_rows":
		return signal.WitnessScope == "local_only"
	case "ci_witnessed_evidence_rows":
		return signal.WitnessScope == "ci_witnessed"
	case "external_witnessed_evidence_rows":
		return signal.WitnessScope == "external_witnessed"
	case "override_rows":
		return signal.OverrideMarker == "override_present"
	case "late_attach_rows":
		return signal.LateAttachMarker == "late_attach_observed"
	case "contract_change_rows":
		return signal.ContractChangeMarker == "contract_change_observed"
	default:
		return false
	}
}

func metricNotAssessed(def metricDef, row query.QueryPackRow, hasSignal bool) bool {
	if def.source == "posture_signal" {
		return !hasSignal
	}
	return row.EvidenceState == query.RowStateNotAssessed
}

func buildMovements(metrics []MetricRow, currentWindow, previousWindow string) ([]MovementRow, MovementSummary) {
	byKey := map[string]map[string]MetricRow{}
	for _, row := range metrics {
		key := row.MetricID + "|" + row.MetricVersion + "|" + row.DimensionKey
		if byKey[key] == nil {
			byKey[key] = map[string]MetricRow{}
		}
		byKey[key][row.TimeWindow] = row
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var rows []MovementRow
	summary := MovementSummary{NonComparableReason: map[string]int{}}
	for i, key := range keys {
		parts := strings.Split(key, "|")
		current, hasCurrent := byKey[key][currentWindow]
		previous, hasPrevious := byKey[key][previousWindow]
		row := MovementRow{
			ID:              fmt.Sprintf("movement.%04d", i+1),
			MetricID:        parts[0],
			MetricVersion:   parts[1],
			DimensionKey:    parts[2],
			ComparisonBasis: "same_profile_metric_dimension_window",
			Comparable:      hasCurrent && hasPrevious,
		}
		if hasCurrent {
			row.CurrentMetricRowRef = current.ID
			row.CurrentValue = current.Numerator
		}
		if hasPrevious {
			row.PreviousMetricRowRef = previous.ID
			row.PreviousValue = previous.Numerator
		}
		row.Delta = row.CurrentValue - row.PreviousValue
		if row.Comparable {
			summary.ComparableCount++
		} else {
			row.ComparisonBasis = "non_comparable_missing_window"
			row.NonComparableReason = "non_comparable_missing_window"
			summary.NonComparableCount++
			summary.NonComparableReason[row.NonComparableReason]++
		}
		rows = append(rows, row)
	}
	return rows, summary
}

func groupingKeys(groupingSet string) []string {
	switch groupingSet {
	case GroupingRepoWindow:
		return []string{"repo", "time_window"}
	case GroupingTeamServiceWindow:
		return []string{"team", "service", "time_window"}
	case GroupingHarnessChangeWindow:
		return []string{"harness", "change_type", "time_window"}
	default:
		return nil
	}
}

func groupingAllowedByExposure(groupingSet string, exposure []string) bool {
	allowed := map[string]bool{"time_window": true}
	for _, key := range exposure {
		allowed[key] = true
	}
	for _, key := range groupingKeys(groupingSet) {
		if !allowed[key] {
			return false
		}
	}
	return true
}

func dimensionKey(repo RepositoryWindow, keys []string) (string, map[string]string) {
	values := map[string]string{
		"repo":        repo.Repo,
		"team":        repo.Team,
		"service":     repo.Service,
		"harness":     repo.Harness,
		"change_type": repo.ChangeType,
		"time_window": repo.TimeWindow,
	}
	dims := map[string]string{}
	var parts []string
	for _, key := range keys {
		dims[key] = values[key]
		if key == "time_window" {
			continue
		}
		parts = append(parts, key+"="+values[key])
	}
	return strings.Join(parts, "|"), dims
}

func validateRepoLabels(repo RepositoryWindow) error {
	if !safeLabel(repo.InputID) || !safeLabel(repo.TimeWindow) {
		return fmt.Errorf("unsafe label")
	}
	labels := map[string]string{
		"repo":        repo.Repo,
		"team":        repo.Team,
		"service":     repo.Service,
		"harness":     repo.Harness,
		"change_type": repo.ChangeType,
	}
	for _, value := range labels {
		if !safeLabel(value) {
			return fmt.Errorf("unsafe label")
		}
	}
	return nil
}

func safeLabel(value string) bool {
	if !safeLabelPattern.MatchString(value) {
		return false
	}
	return !unsafeLabel(value)
}

var unsafeOutputKeywords = []string{
	"http://",
	"https://",
	"secret",
	"@",
}

func unsafeOutput(value string) bool {
	lower := strings.ToLower(value)
	if hasUnsafeOutputKeyword(lower) {
		return true
	}
	if hasUnsafeTokenOrCredential(lower) {
		return true
	}
	if hasUnsafePath(value) {
		return true
	}
	return false
}

func hasUnsafeOutputKeyword(value string) bool {
	for _, keyword := range unsafeOutputKeywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}
	return false
}

func hasUnsafeTokenOrCredential(value string) bool {
	if strings.Contains(value, "credential_or_token") {
		return false
	}
	return strings.Contains(value, "token") || strings.Contains(value, "credential")
}

func hasUnsafePath(value string) bool {
	return strings.Contains(value, "/") || strings.Contains(value, "\\")
}

func unsafeLabel(value string) bool {
	lower := strings.ToLower(value)
	return unsafeOutput(value) ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "credential")
}

func validateInputPaths(repo RepositoryWindow) error {
	for _, path := range []string{repo.QueryPackResult, repo.ArtifactDigestManifest, repo.PostureSignalManifest} {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if unsafeSelectionPath(path) {
			return fmt.Errorf("unsafe input path")
		}
	}
	return nil
}

func unsafeSelectionPath(value string) bool {
	clean := strings.ReplaceAll(value, "\\", "/")
	return strings.Contains(clean, "://") ||
		hasWindowsVolume(clean) ||
		strings.HasPrefix(clean, "/") ||
		strings.HasPrefix(clean, "../") ||
		strings.Contains(clean, "../") ||
		strings.Contains(clean, "/..")
}

func hasWindowsVolume(value string) bool {
	return len(value) >= 3 &&
		((value[0] >= 'A' && value[0] <= 'Z') || (value[0] >= 'a' && value[0] <= 'z')) &&
		value[1] == ':' &&
		value[2] == '/'
}

func safeHandoff(values map[string]string) bool {
	for key, value := range values {
		if !safeLabel(key) || unsafeOutput(value) {
			return false
		}
	}
	return true
}

func filepathBase(path string) string {
	clean := strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(clean, "/")
	return parts[len(parts)-1]
}

func parseFreshnessBoundary(value string, now time.Time) (time.Time, bool, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, false, nil
	}
	if strings.HasPrefix(value, "P") {
		return time.Time{}, false, fmt.Errorf("duration freshness boundaries are not supported in v1")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false, err
	}
	return parsed, true, nil
}

func isStale(value string, cutoff time.Time) bool {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return true
	}
	return parsed.Before(cutoff)
}

func inputSelection(repo RepositoryWindow, digest, trust string) InputSelection {
	return InputSelection{
		InputID:         repo.InputID,
		Repository:      repo.Repo,
		TimeWindow:      repo.TimeWindow,
		PathRedactedID:  "artifact:query_pack_result:" + shortDigest(digest),
		SHA256:          digest,
		InputTrustState: trust,
	}
}

func refusal(counter int, repo RepositoryWindow, reason, trust string) RefusalRow {
	return RefusalRow{
		ID:              fmt.Sprintf("refusal.%04d", counter),
		InputID:         repo.InputID,
		TimeWindow:      repo.TimeWindow,
		RefusalReason:   reason,
		InputTrustState: trust,
	}
}

func digestSetHash(digests []string) string {
	sum := sha256.Sum256([]byte(strings.Join(digests, "\n")))
	return hex.EncodeToString(sum[:])
}

func deterministicExportID(selectionPath string, metrics []MetricRow, refusals []RefusalRow) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%d", selectionPath, len(metrics), len(refusals))))
	return "export:" + hex.EncodeToString(sum[:8])
}

func shortDigest(digest string) string {
	if len(digest) >= 16 {
		return digest[:16]
	}
	return "not_assessed0000"
}

func copyTrust(in map[string]int) map[string]int {
	out := map[string]int{}
	for key, value := range in {
		out[key] = value
	}
	return out
}
