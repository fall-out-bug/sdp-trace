package posture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/query"
	"os"
	"regexp"

	"sort"
	"strings"
	"time"
)

const (
	// SchemaVersion identifies the only posture export contract this package
	// can build or validate.
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

var signalMetricPredicates = map[string]func(PostureSignal) bool{
	"local_only_evidence_rows":         func(signal PostureSignal) bool { return signal.WitnessScope == "local_only" },
	"ci_witnessed_evidence_rows":       func(signal PostureSignal) bool { return signal.WitnessScope == "ci_witnessed" },
	"external_witnessed_evidence_rows": func(signal PostureSignal) bool { return signal.WitnessScope == "external_witnessed" },
	"override_rows":                    func(signal PostureSignal) bool { return signal.OverrideMarker == "override_present" },
	"late_attach_rows":                 func(signal PostureSignal) bool { return signal.LateAttachMarker == "late_attach_observed" },
	"contract_change_rows":             func(signal PostureSignal) bool { return signal.ContractChangeMarker == "contract_change_observed" },
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

// RepositoryWindow binds one selected repository window to its replayable
// query, digest, and optional posture-signal inputs.
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

// DigestManifest records artifact digests that must match before selected
// query-pack evidence can enter posture aggregation.
type DigestManifest struct {
	SchemaVersion string           `json:"schema_version"`
	Artifacts     []DigestArtifact `json:"artifacts"`
}

type DigestArtifact struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// SignalManifest carries optional posture signals that may refine metric counts
// only after signal payload safety checks pass.
type SignalManifest struct {
	SchemaVersion string          `json:"schema_version"`
	Signals       []PostureSignal `json:"signals"`
}

// PostureSignal contains row-bound posture evidence that is never trusted as
// free text; every field is screened before aggregation.
type PostureSignal struct {
	RowRef               string `json:"row_ref"`
	WitnessScope         string `json:"witness_scope,omitempty"`
	ObserverState        string `json:"observer_state,omitempty"`
	OverrideMarker       string `json:"override_marker,omitempty"`
	LateAttachMarker     string `json:"late_attach_marker,omitempty"`
	ContractChangeMarker string `json:"contract_change_marker,omitempty"`
}

// ExportResult is the portable posture export surface produced from trusted
// selections, refusals, metrics, movement rows, and output-safety claims.
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

// InputSelection records selected input identity without exposing raw artifact
// paths in downstream posture output.
type InputSelection struct {
	InputID         string `json:"input_id"`
	Repository      string `json:"repo"`
	TimeWindow      string `json:"time_window"`
	PathRedactedID  string `json:"path_redacted_id"`
	SHA256          string `json:"sha256,omitempty"`
	InputTrustState string `json:"input_trust_state"`
}

// MetricRow is an evidence-backed aggregate for one metric, window, and
// dimension key.
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

// MovementRow compares matching metric rows across the selected current and
// previous windows without fabricating missing-window evidence.
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

// MovementSummary totals comparable and non-comparable movement rows using the
// same reason vocabulary as individual rows.
type MovementSummary struct {
	ComparableCount     int            `json:"comparable_count"`
	NonComparableCount  int            `json:"non_comparable_count"`
	NonComparableReason map[string]int `json:"non_comparable_reason_counts"`
}

// RefusalRow preserves rejected inputs as explicit posture evidence rather than
// silently dropping them from the export.
type RefusalRow struct {
	ID              string `json:"id"`
	InputID         string `json:"input_id"`
	TimeWindow      string `json:"time_window,omitempty"`
	RefusalReason   string `json:"refusal_reason"`
	InputTrustState string `json:"input_trust_state"`
}

// OutputSafety lists sensitive classes this package checked as absent from
// rendered human-readable posture output.
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

	selection, err := readSelection(selectionPath)
	if err != nil {
		return ExportResult{}, err
	}
	input, err := prepareBuildSelectionInput(now, selection)
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

func prepareBuildSelectionInput(now time.Time, selection SelectionManifest) (buildInput, error) {

	activeKeys, err := validateBuildSelection(selection)
	if err != nil {
		return buildInput{}, err
	}
	cutoff, hasCutoff, err := parseFreshnessBoundary(selection.FreshnessBoundary, now)
	if err != nil {
		return buildInput{}, err
	}

	handoff, err := validatedHandoff(selection.Handoff)
	if err != nil {
		return buildInput{}, err
	}
	return buildInput{selection: selection, activeKeys: activeKeys, cutoff: cutoff, hasCutoff: hasCutoff, handoff: handoff}, nil
}

func validateBuildSelection(selection SelectionManifest) ([]string, error) {
	if err := validateSelection(selection); err != nil {
		return nil, err
	}

	activeKeys := groupingKeys(selection.GroupingSetID)
	if len(activeKeys) == 0 {
		return nil, fmt.Errorf("unsupported grouping set")
	}
	return activeKeys, nil
}

func validatedHandoff(handoff map[string]string) (map[string]string, error) {
	if handoff == nil {

		handoff = map[string]string{}
	}
	if !safeHandoff(handoff) {
		return nil, fmt.Errorf("unsafe handoff")
	}
	return handoff, nil
}
func ingestRepositories(selection SelectionManifest, activeKeys []string, cutoff time.Time, hasCutoff bool) ([]InputSelection, []RefusalRow, map[string]*aggregateGroup) {
	inputs := []InputSelection{}
	refusals := []RefusalRow{}
	groups := map[string]*aggregateGroup{}
	refusalCounter := 0

	for _, repo := range selection.Repositories {
		ingested := ingestRepository(repo, cutoff, hasCutoff)
		if !ingested.trusted {

			inputs, refusals, refusalCounter = recordRefusedRepository(inputs, refusals, refusalCounter, repo, ingested)
			continue
		}

		inputs = append(inputs, inputSelection(repo, ingested.digest, "trusted_input"))
		addTrustedRepositoryGroup(groups, repo, activeKeys, ingested.result, ingested.signals, ingested.digest)
	}

	return inputs, refusals, groups
}

func recordRefusedRepository(inputs []InputSelection, refusals []RefusalRow, counter int, repo RepositoryWindow, ingested repositoryIngest) ([]InputSelection, []RefusalRow, int) {

	counter++
	refusals = append(refusals, refusal(counter, repo, ingested.refusalReason, ingested.inputTrustState))
	if ingested.recordSelection {
		inputs = append(inputs, inputSelection(repo, ingested.digest, ingested.inputTrustState))
	}
	return inputs, refusals, counter
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

	if refusal, ok := repositoryPreflightRefusal(repo, cutoff, hasCutoff); ok {
		return refusal
	}
	return ingestRepositoryArtifacts(repo)
}

func repositoryPreflightRefusal(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) (repositoryIngest, bool) {

	if invalidRepositoryLabels(repo) {
		return refusedInput("unsafe_label", "cannot_verify_input", "", false), true
	}
	if invalidRepositoryInputPaths(repo) {

		return refusedInput("malformed_input", "cannot_verify_input", "", false), true
	}
	if staleRepositoryInput(repo, cutoff, hasCutoff) {
		return refusedInput("stale_input", "stale_input", "", true), true
	}
	return repositoryIngest{}, false
}

func staleRepositoryInput(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) bool {
	return hasCutoff && isStale(repo.InputObservedAt, cutoff)
}

func invalidRepositoryLabels(repo RepositoryWindow) bool {
	return validateRepoLabels(repo) != nil
}

func invalidRepositoryInputPaths(repo RepositoryWindow) bool {
	return validateInputPaths(repo) != nil
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
	return ingestSupportedRepositoryArtifacts(repo, digest, result)
}

func ingestSupportedRepositoryArtifacts(repo RepositoryWindow, digest string, result query.QueryPackResult) repositoryIngest {
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

// refusedInput crosses the trust boundary from ingestion to refusal record.
// Refused inputs retain a trust state but do not contribute to metric evidence.
func refusedInput(reason, trustState, digest string, recordSelection bool) repositoryIngest {
	return repositoryIngest{refusalReason: reason, inputTrustState: trustState, digest: digest, recordSelection: recordSelection}
}

func addTrustedRepositoryGroup(groups map[string]*aggregateGroup, repo RepositoryWindow, activeKeys []string, result query.QueryPackResult, signals map[string]PostureSignal, digest string) {
	group := trustedAggregateGroup(groups, repo, activeKeys)
	group.rows = append(group.rows, flattenRows(result)...)
	group.inputRefs = append(group.inputRefs, repo.InputID)
	group.digests = append(group.digests, digest)
	group.trustStates["trusted_input"]++

	for rowRef, signal := range signals {
		group.signals[rowRef] = signal
	}
}

func trustedAggregateGroup(groups map[string]*aggregateGroup, repo RepositoryWindow, activeKeys []string) *aggregateGroup {

	key, dims := dimensionKey(repo, activeKeys)
	groupKey := repo.TimeWindow + "|" + key
	group := groups[groupKey]
	if group == nil {

		group = newAggregateGroup(repo, key, dims)
		groups[groupKey] = group
	}
	return group
}

func newAggregateGroup(repo RepositoryWindow, key string, dims map[string]string) *aggregateGroup {

	return &aggregateGroup{
		dimensions:   dims,
		dimensionKey: key,
		window:       repo.TimeWindow,
		trustStates:  map[string]int{},
		signals:      map[string]PostureSignal{},
	}
}
func buildExportResult(selectionPath string, now time.Time, input buildInput, inputs []InputSelection, metricRows []MetricRow, movementRows []MovementRow, summary MovementSummary, refusals []RefusalRow) ExportResult {

	result := exportResultHeader(selectionPath, now, input, metricRows, refusals)
	result.InputSelection = inputs
	result.MetricRows = metricRows
	result.MovementRows = movementRows
	result.MovementSummary = summary
	result.RefusalRows = refusals
	return result
}

func exportResultHeader(selectionPath string, now time.Time, input buildInput, metricRows []MetricRow, refusals []RefusalRow) ExportResult {

	return ExportResult{

		SchemaVersion:        SchemaVersion,
		ExportProfileID:      ProfileID,
		ExportProfileVersion: ProfileVer,
		ExportID:             deterministicExportID(selectionPath, metricRows, refusals),
		Producer:             "sdp-trace",
		GeneratedAt:          now.UTC().Format(time.RFC3339),

		GroupingSetID:      input.selection.GroupingSetID,
		ActiveGroupingKeys: input.activeKeys,
		Handoff:            input.handoff,
		OutputSafety:       exportOutputSafety(),
	}
}

func exportOutputSafety() OutputSafety {

	return OutputSafety{VerifiedAbsentSensitiveClasses: SensitiveClasses()}
}

func Explain(result ExportResult) (string, error) {
	// Explanation output is derived from structured rows only, then checked
	// again so renderer text cannot leak unsafe payload classes.
	var lines []string
	lines = append(lines, explainHeaderLines(result)...)
	lines = append(lines, fmt.Sprintf("movement_summary comparable=%d non_comparable=%d", result.MovementSummary.ComparableCount, result.MovementSummary.NonComparableCount))
	lines = append(lines, explainRefusalLines(result.RefusalRows)...)
	lines = append(lines, explainMetricLines(result.MetricRows)...)
	lines = append(lines, explainMovementLines(result.MovementRows)...)
	lines = append(lines, explainOutputSafetyLines(result.OutputSafety.VerifiedAbsentSensitiveClasses)...)

	rendered := strings.Join(lines, "\n") + "\n"
	if unsafeOutput(rendered) {
		return "", fmt.Errorf("output_safety_violation")
	}
	return rendered, nil
}

func explainHeaderLines(result ExportResult) []string {

	return []string{
		"schema_version=" + result.SchemaVersion,
		"export_profile_id=" + result.ExportProfileID,
		"grouping_set_id=" + result.GroupingSetID,
	}
}

func explainRefusalLines(rows []RefusalRow) []string {
	return formattedLines(rows, explainRefusalLine)
}

func explainRefusalLine(row RefusalRow) string {
	return fmt.Sprintf("refusal %s input=%s reason=%s state=%s", row.ID, row.InputID, row.RefusalReason, row.InputTrustState)
}

func explainMetricLines(rows []MetricRow) []string {
	return formattedLines(rows, explainMetricLine)
}

func explainMetricLine(row MetricRow) string {
	return fmt.Sprintf("metric %s %s numerator=%d denominator=%d window=%s dimension_key=%s not_assessed_count=%d", row.ID, row.MetricID, row.Numerator, row.Denominator, row.TimeWindow, row.DimensionKey, row.NotAssessedCount)
}

func explainMovementLines(rows []MovementRow) []string {
	return formattedLines(rows, explainMovementLine)
}

func explainMovementLine(row MovementRow) string {
	return fmt.Sprintf("movement %s %s current=%d previous=%d delta=%d comparable=%t reason=%s", row.ID, row.MetricID, row.CurrentValue, row.PreviousValue, row.Delta, row.Comparable, row.NonComparableReason)
}

func explainOutputSafetyLines(classes []string) []string {
	return formattedLines(classes, explainOutputSafetyLine)
}

func explainOutputSafetyLine(class string) string {
	return "output_safety absent=" + class
}

func formattedLines[T any](rows []T, format func(T) string) []string {

	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		lines = append(lines, format(row))
	}
	return lines
}
func ValidateExportResult(result ExportResult) error {

	for _, validate := range exportResultValidators {
		if err := validate(result); err != nil {
			return err
		}
	}
	return nil
}

var exportResultValidators = []func(ExportResult) error{
	validateExportHeader,
	validateExportCollections,
	validateExportInputSelections,
	validateExportMetrics,
	validateExportMovements,
	validateExportMovementSummary,
	validateExportRefusals,
}

func validateExportHeader(result ExportResult) error {

	if err := validateExportSchemaHeader(result); err != nil {
		return err
	}
	return validateExportGeneratedAt(result.GeneratedAt)
}

func validateExportSchemaHeader(result ExportResult) error {

	if unsupportedExportHeader(result) {
		return errors.New("unsupported posture export")
	}
	if malformedExportHeader(result) {
		return errors.New("malformed posture export")
	}
	return nil
}

func validateExportGeneratedAt(generatedAt string) error {

	if _, err := time.Parse(time.RFC3339, generatedAt); err != nil {
		return errors.New("malformed posture export generated_at")
	}
	return nil
}

func validateExportInputSelections(result ExportResult) error {
	return validateInputSelectionRows(result.InputSelection)
}

func validateExportMetrics(result ExportResult) error {
	return validateMetricRows(result.MetricRows)
}

func validateExportMovements(result ExportResult) error {
	return validateMovementRows(result.MovementRows)
}

func validateExportMovementSummary(result ExportResult) error {
	return validateMovementSummary(result.MovementSummary)
}

func validateExportRefusals(result ExportResult) error {
	return validateRefusalRows(result.RefusalRows)
}

func unsupportedExportHeader(result ExportResult) bool {
	return result.SchemaVersion != SchemaVersion ||
		result.ExportProfileID != ProfileID ||
		result.ExportProfileVersion != ProfileVer
}

func malformedExportHeader(result ExportResult) bool {
	return result.ExportID == "" || result.Producer != "sdp-trace" || result.GeneratedAt == ""
}

func validateInputSelectionRows(rows []InputSelection) error {
	return validateRows(rows, validateInputSelectionRow)
}

func validateMetricRows(rows []MetricRow) error {
	return validateRows(rows, validateMetricRow)
}

func validateMovementRows(rows []MovementRow) error {
	return validateRows(rows, validateMovementRow)
}

func validateRefusalRows(rows []RefusalRow) error {
	return validateRows(rows, validateRefusalRow)
}

func validateRows[T any](rows []T, validate func(T) error) error {

	for _, row := range rows {
		if err := validate(row); err != nil {
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
	return malformedRowError(malformedInputSelectionRow(item), "malformed posture export input_selection")
}

func malformedInputSelectionRow(item InputSelection) bool {
	return !safeLabel(item.InputID) || !safeLabel(item.Repository) || !safeLabel(item.TimeWindow) || missingInputSelectionFields(item)
}

func missingInputSelectionFields(item InputSelection) bool {
	return item.PathRedactedID == "" || !validInputTrustState(item.InputTrustState)
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
	return malformedRowError(malformedMovementRow(row), "malformed posture export movement_row")
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
	if malformedMovementSummaryReasons(summary.NonComparableReason) {
		return fmt.Errorf("malformed posture export movement_summary")
	}
	return nil
}

func malformedMovementSummaryReasons(reasons map[string]int) bool {
	for reason, count := range reasons {

		if reason != "non_comparable_missing_window" || count < 0 {
			return true
		}
	}
	return false
}

func validateRefusalRow(row RefusalRow) error {
	return malformedRowError(malformedRefusalRow(row), "malformed posture export refusal_row")
}

func malformedRowError(malformed bool, message string) error {

	if malformed {
		return errors.New(message)
	}
	return nil
}

func malformedRefusalRow(row RefusalRow) bool {
	return missingRefusalIdentity(row) || malformedRefusalState(row) || malformedOptionalRefusalWindow(row)
}

func missingRefusalIdentity(row RefusalRow) bool {
	return row.ID == "" || !safeLabel(row.InputID)
}

func malformedRefusalState(row RefusalRow) bool {
	return !validRefusalReason(row.RefusalReason) || !validInputTrustState(row.InputTrustState)
}

func malformedOptionalRefusalWindow(row RefusalRow) bool {
	return row.TimeWindow != "" && !safeLabel(row.TimeWindow)
}

func SensitiveClasses() []string {

	return append([]string(nil), sensitiveClasses...)
}

var sensitiveClasses = []string{
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
	return readJSONFile[SelectionManifest](path)
}

func validateSelection(selection SelectionManifest) error {

	if selection.SchemaVersion != SelectionSchemaVersion {
		return fmt.Errorf("unsupported selection schema")
	}
	if err := validateSelectionProfile(selection); err != nil {
		return err
	}
	if err := validateSelectionGrouping(selection); err != nil {
		return err
	}
	return validateSelectionRepositories(selection)
}

func validateSelectionProfile(selection SelectionManifest) error {

	switch {
	case selection.ProfileID != ProfileID:
		return fmt.Errorf("unsupported profile")
	case unsupportedSelectionProfileVersion(selection.ProfileVersion):
		return fmt.Errorf("unsupported profile version")
	}
	return nil
}

func unsupportedSelectionProfileVersion(version string) bool {
	return version != "" && version != ProfileVer
}

func validateSelectionGrouping(selection SelectionManifest) error {
	if len(groupingKeys(selection.GroupingSetID)) == 0 {
		return fmt.Errorf("unsupported grouping set")
	}

	if !groupingAllowedByExposure(selection.GroupingSetID, selection.DimensionExposurePolicy) {
		return fmt.Errorf("dimension exposure policy excludes grouping key")
	}
	return nil
}

func validateSelectionRepositories(selection SelectionManifest) error {

	if len(selection.Repositories) == 0 {
		return fmt.Errorf("empty selection")
	}
	return nil
}

func readQueryPack(path string) (query.QueryPackResult, error) {
	return readJSONFile[query.QueryPackResult](path)
}

func readJSONFile[T any](path string) (T, error) {
	// JSON replay inputs must parse structurally before profile-specific checks.
	var result T
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
	manifest, err := readSignalManifest(path)
	if err != nil {
		return nil, err
	}
	return validatedSignals(manifest.Signals)
}

func readSignalManifest(path string) (SignalManifest, error) {
	manifest, err := readJSONFile[SignalManifest](path)
	if err != nil {
		return manifest, err
	}

	if manifest.SchemaVersion != SignalManifestSchemaVersion {
		return manifest, fmt.Errorf("unsupported signal manifest schema")
	}
	return manifest, nil
}

func validatedSignals(signals []PostureSignal) (map[string]PostureSignal, error) {

	out := map[string]PostureSignal{}
	for _, signal := range signals {
		if unsafeSignal(signal) {
			return nil, fmt.Errorf("unsafe signal")
		}
		out[signal.RowRef] = signal
	}
	return out, nil
}

// unsafeSignal checks whether a posture signal crosses the output safety boundary
// before incorporation into metric evidence. Unsafe keywords block signal incorporation.
func unsafeSignal(signal PostureSignal) bool {
	return unsafeOutput(signal.RowRef + signal.WitnessScope + signal.ObserverState + signal.OverrideMarker + signal.LateAttachMarker + signal.ContractChangeMarker)
}

func verifyDigestManifest(manifestPath, queryPackPath string) (string, error) {

	expected, err := expectedDigestForQueryPack(manifestPath, queryPackPath)
	if err != nil {
		return "", err
	}

	actual, err := fileSHA256Hex(queryPackPath)
	if err != nil {
		return "", err
	}
	return checkedDigest(actual, expected)
}

func expectedDigestForQueryPack(manifestPath, queryPackPath string) (string, error) {

	manifest, err := readDigestManifest(manifestPath)
	if err != nil {
		return "", err
	}
	return digestForQueryPackFromManifest(manifest, queryPackPath)
}

// checkedDigest enforces the evidence boundary where manifest expectation meets artifact reality.
// A digest mismatch crosses from unverifiable into untrusted, distinct from missing or malformed.
func checkedDigest(actual, expected string) (string, error) {

	if expected != actual {
		return "", errDigestMismatch
	}
	return actual, nil
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

type digestErrorReason struct {
	err    error
	reason string
}

var digestErrorReasons = []digestErrorReason{
	{err: errDigestMismatch, reason: "untrusted_input_digest_mismatch"},
	{err: errMissingRequired, reason: "missing_required_input"},
	{err: errUnsafePath, reason: "malformed_input"},
}

// reasonForDigestErr maps digest errors to refusal reasons at the trust boundary.
// Mismatch is distinct to avoid blurring tamper evidence with cannot-verify failures.
func reasonForDigestErr(err error) string {

	for _, item := range digestErrorReasons {
		if errors.Is(err, item.err) {
			return item.reason
		}
	}
	return "malformed_input"
}

// trustForDigestErr maps digest verification outcomes to trust states at the evidence boundary.
// Distinguishes tamper evidence (untrusted) from verification failures (cannot_verify).
func trustForDigestErr(err error) string {

	if errors.Is(err, errDigestMismatch) {
		return "untrusted_input"
	}
	return "cannot_verify_input"
}

func flattenRows(result query.QueryPackResult) []query.QueryPackRow {
	var rows []query.QueryPackRow

	for _, key := range sortedMapKeys(result.QueryRows) {
		rows = append(rows, result.QueryRows[key]...)
	}
	return rows
}

func buildMetrics(groups map[string]*aggregateGroup) []MetricRow {
	// Metric rows are emitted in stable group/catalog order because row ids are
	// part of downstream movement and explanation references.
	var rows []MetricRow
	counter := 0
	for _, groupKey := range sortedMapKeys(groups) {
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
	counts := metricCounts(def, group)
	sourceRefs, digestHash := metricEvidenceRefs(group)

	row := metricRowHeader(counter, def, group)
	row.Numerator = counts.numerator
	row.Denominator = len(group.rows)
	row.SourceInputRefs = sourceRefs
	row.SourceArtifactDigestSet = digestHash
	row.SourceFieldState = counts.sourceFieldState
	row.NotAssessedCount = counts.notAssessed
	row.InputTrustStateSummary = copyTrust(group.trustStates)
	return row
}

func metricRowHeader(counter int, def metricDef, group *aggregateGroup) MetricRow {

	return MetricRow{
		ID:            fmt.Sprintf("metric.%04d", counter),
		MetricID:      def.id,
		MetricVersion: def.version,
		Unit:          "rows",
		TimeWindow:    group.window,
		Dimensions:    group.dimensions,
		DimensionKey:  group.dimensionKey,
	}
}

func metricEvidenceRefs(group *aggregateGroup) ([]string, string) {

	sourceRefs := sortedStringsCopy(group.inputRefs)
	digests := sortedStringsCopy(group.digests)
	return sourceRefs, digestSetHash(digests)
}

func sortedStringsCopy(values []string) []string {

	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}

type metricCount struct {
	numerator        int
	notAssessed      int
	sourceFieldState string
}

func metricCounts(def metricDef, group *aggregateGroup) metricCount {

	counts := metricCount{sourceFieldState: "present"}
	for _, row := range group.rows {
		signal, hasSignal := group.signals[row.ID]
		applyMetricCount(&counts, def, row, signal, hasSignal)
	}
	return counts
}

func applyMetricCount(counts *metricCount, def metricDef, row query.QueryPackRow, signal PostureSignal, hasSignal bool) {

	if metricMatches(def.id, row, signal, hasSignal) {
		counts.numerator++
	}
	if metricNotAssessed(def, row, hasSignal) {
		counts.notAssessed++
	}
	if missingPostureSignalMetric(def, hasSignal) {
		counts.sourceFieldState = "not_assessed"
	}
}

// missingPostureSignalMetric gates the evidence boundary for signal-sourced metrics.
// Absence of posture signal marks evidence-absent, downgrading sourceFieldState.
func missingPostureSignalMetric(def metricDef, hasSignal bool) bool {
	return def.source == "posture_signal" && !hasSignal
}

func metricMatches(metricID string, row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {

	if expectedState, ok := rowStateMetrics[metricID]; ok {
		return row.EvidenceState == expectedState
	}
	return nonRowStateMetricMatches(metricID, row, signal, hasSignal)
}

func nonRowStateMetricMatches(metricID string, row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {

	if metricID == "unsupported_observer_rows" {
		return unsupportedObserverMetricMatches(row, signal, hasSignal)
	}
	return hasSignal && signalMetricMatches(metricID, signal)
}

func unsupportedObserverMetricMatches(row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {
	return row.EvidenceState == query.RowStateUnsupported || (hasSignal && signal.ObserverState == "unsupported")
}

func signalMetricMatches(metricID string, signal PostureSignal) bool {
	matches, ok := signalMetricPredicates[metricID]
	return ok && matches(signal)
}

// metricNotAssessed applies the evidence-boundary rule for not_assessed counts.
// For row_state metrics, follows row.EvidenceState. For signal metrics, signal absence is not_assessed.
func metricNotAssessed(def metricDef, row query.QueryPackRow, hasSignal bool) bool {
	if def.source == "posture_signal" {

		return !hasSignal
	}
	return row.EvidenceState == query.RowStateNotAssessed
}

func buildMovements(metrics []MetricRow, currentWindow, previousWindow string) ([]MovementRow, MovementSummary) {
	byKey := metricsByMovementKey(metrics)
	keys := sortedMovementKeys(byKey)
	var rows []MovementRow
	summary := MovementSummary{NonComparableReason: map[string]int{}}
	for i, key := range keys {

		row := movementRowForKey(i+1, key, byKey[key], currentWindow, previousWindow)
		summarizeMovement(&summary, &row)
		rows = append(rows, row)
	}
	return rows, summary
}

func sortedMovementKeys(metrics map[string]map[string]MetricRow) []string {
	return sortedMapKeys(metrics)
}

func sortedMapKeys[T any](items map[string]T) []string {

	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func movementRowForKey(index int, key string, rowsByWindow map[string]MetricRow, currentWindow, previousWindow string) MovementRow {

	parts := strings.Split(key, "|")
	current, hasCurrent := rowsByWindow[currentWindow]
	previous, hasPrevious := rowsByWindow[previousWindow]

	row := MovementRow{
		ID:              fmt.Sprintf("movement.%04d", index),
		MetricID:        parts[0],
		MetricVersion:   parts[1],
		DimensionKey:    parts[2],
		ComparisonBasis: "same_profile_metric_dimension_window",
		Comparable:      hasCurrent && hasPrevious,
	}
	applyMovementWindowValues(&row, current, hasCurrent, previous, hasPrevious)
	return row
}

func applyMovementWindowValues(row *MovementRow, current MetricRow, hasCurrent bool, previous MetricRow, hasPrevious bool) {
	if hasCurrent {
		row.CurrentMetricRowRef = current.ID
		row.CurrentValue = current.Numerator
	}
	if hasPrevious {
		row.PreviousMetricRowRef = previous.ID
		row.PreviousValue = previous.Numerator
	}

	row.Delta = row.CurrentValue - row.PreviousValue
}

func metricsByMovementKey(metrics []MetricRow) map[string]map[string]MetricRow {
	byKey := map[string]map[string]MetricRow{}
	for _, row := range metrics {

		key := row.MetricID + "|" + row.MetricVersion + "|" + row.DimensionKey
		if byKey[key] == nil {
			byKey[key] = map[string]MetricRow{}
		}
		byKey[key][row.TimeWindow] = row
	}
	return byKey
}

// summarizeMovement applies the threshold rule for movement comparability.
// Presence of both current and previous window rows is the evidence boundary.
func summarizeMovement(summary *MovementSummary, row *MovementRow) {
	if row.Comparable {
		summary.ComparableCount++
		return
	}

	row.ComparisonBasis = "non_comparable_missing_window"
	row.NonComparableReason = "non_comparable_missing_window"
	summary.NonComparableCount++
	summary.NonComparableReason[row.NonComparableReason]++
}

var groupingKeysBySet = map[string][]string{
	GroupingRepoWindow:          {"repo", "time_window"},
	GroupingTeamServiceWindow:   {"team", "service", "time_window"},
	GroupingHarnessChangeWindow: {"harness", "change_type", "time_window"},
}

func groupingKeys(groupingSet string) []string {
	return groupingKeysBySet[groupingSet]
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

	values := dimensionValues(repo)
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

func dimensionValues(repo RepositoryWindow) map[string]string {

	return map[string]string{
		"repo":        repo.Repo,
		"team":        repo.Team,
		"service":     repo.Service,
		"harness":     repo.Harness,
		"change_type": repo.ChangeType,
		"time_window": repo.TimeWindow,
	}
}

func validateRepoLabels(repo RepositoryWindow) error {
	if unsafeInputLabel(repo) {
		return fmt.Errorf("unsafe label")
	}

	for _, value := range repositoryOutputLabels(repo) {

		if !safeLabel(value) {
			return fmt.Errorf("unsafe label")
		}
	}
	return nil
}

func repositoryOutputLabels(repo RepositoryWindow) []string {

	return []string{repo.Repo, repo.Team, repo.Service, repo.Harness, repo.ChangeType}
}

func unsafeInputLabel(repo RepositoryWindow) bool {
	return !safeLabel(repo.InputID) || !safeLabel(repo.TimeWindow)
}

// safeLabel enforces the input trust boundary for identifier syntax.
// Unsafe output keywords or credential/token terminology crosses the injection boundary.
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
	return hasUnsafeSelectionPathPrefix(clean) || strings.Contains(clean, "../") || strings.Contains(clean, "/..")
}

func hasUnsafeSelectionPathPrefix(clean string) bool {

	return strings.Contains(clean, "://") ||
		hasWindowsVolume(clean) ||
		strings.HasPrefix(clean, "/") ||
		strings.HasPrefix(clean, "../")
}

func hasWindowsVolume(value string) bool {
	return len(value) >= 3 && isASCIIAlpha(value[0]) && value[1] == ':' && value[2] == '/'
}

func isASCIIAlpha(value byte) bool {
	return (value >= 'A' && value <= 'Z') || (value >= 'a' && value <= 'z')
}

// safeHandoff enforces the export handoff boundary. Keys must be safe labels;
// values must not contain unsafe output. Crossing this threshold prevents injection.
func safeHandoff(values map[string]string) bool {
	for key, value := range values {

		if !safeLabel(key) || unsafeOutput(value) {
			return false
		}
	}
	return true
}

// filepathBase extracts the basename for manifest-artifact matching.
// Slash normalization ensures Windows-style separators cannot escape the comparison.
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

// isStale applies the temporal evidence boundary. Parse failures default to stale (safe).
// InputObservedAt before cutoff crosses the freshness threshold.
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

// deterministicExportID produces identifiers stable across runs. These are identifiers,
// not proof digests; they do not carry evidence authority.
func deterministicExportID(selectionPath string, metrics []MetricRow, refusals []RefusalRow) string {

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%d", selectionPath, len(metrics), len(refusals))))
	return "export:" + hex.EncodeToString(sum[:8])
}

// shortDigest produces a stable redacted identifier. When crossing from verified digest
// to missing/empty, the evidence boundary requires a placeholder distinct from a real hash.
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
