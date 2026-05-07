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
	selection, err := readSelection(selectionPath)
	if err != nil {
		return ExportResult{}, err
	}
	if err := validateSelection(selection); err != nil {
		return ExportResult{}, err
	}
	activeKeys := groupingKeys(selection.GroupingSetID)
	if len(activeKeys) == 0 {
		return ExportResult{}, fmt.Errorf("unsupported grouping set")
	}
	cutoff, hasCutoff, err := parseFreshnessBoundary(selection.FreshnessBoundary, now)
	if err != nil {
		return ExportResult{}, err
	}

	var inputs []InputSelection
	var refusals []RefusalRow
	groups := map[string]*aggregateGroup{}
	refusalCounter := 0
	handoff := selection.Handoff
	if handoff == nil {
		handoff = map[string]string{}
	}
	if !safeHandoff(handoff) {
		return ExportResult{}, fmt.Errorf("unsafe handoff")
	}
	for _, repo := range selection.Repositories {
		if err := validateRepoLabels(repo); err != nil {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, "unsafe_label", "cannot_verify_input"))
			continue
		}
		if err := validateInputPaths(repo); err != nil {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, "malformed_input", "cannot_verify_input"))
			continue
		}
		if hasCutoff && isStale(repo.InputObservedAt, cutoff) {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, "stale_input", "stale_input"))
			inputs = append(inputs, inputSelection(repo, "", "stale_input"))
			continue
		}
		digest, err := verifyDigestManifest(repo.ArtifactDigestManifest, repo.QueryPackResult)
		if err != nil {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, reasonForDigestErr(err), trustForDigestErr(err)))
			inputs = append(inputs, inputSelection(repo, "", trustForDigestErr(err)))
			continue
		}
		result, err := readQueryPack(repo.QueryPackResult)
		if err != nil || result.SchemaVersion != query.QueryPackSchemaVersion || result.QueryPackID != query.QueryPackForensicsBasic {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, "malformed_input", "cannot_verify_input"))
			inputs = append(inputs, inputSelection(repo, digest, "cannot_verify_input"))
			continue
		}
		signals, err := readSignals(repo.PostureSignalManifest)
		if err != nil {
			refusalCounter++
			refusals = append(refusals, refusal(refusalCounter, repo, "malformed_input", "cannot_verify_input"))
			inputs = append(inputs, inputSelection(repo, digest, "cannot_verify_input"))
			continue
		}
		inputs = append(inputs, inputSelection(repo, digest, "trusted_input"))
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

	metricRows := buildMetrics(groups)
	movementRows, summary := buildMovements(metricRows, selection.CurrentWindow, selection.PreviousWindow)
	return ExportResult{
		SchemaVersion:        SchemaVersion,
		ExportProfileID:      ProfileID,
		ExportProfileVersion: ProfileVer,
		ExportID:             deterministicExportID(selectionPath, metricRows, refusals),
		Producer:             "sdp-trace",
		GeneratedAt:          now.UTC().Format(time.RFC3339),
		GroupingSetID:        selection.GroupingSetID,
		ActiveGroupingKeys:   activeKeys,
		InputSelection:       inputs,
		MetricRows:           metricRows,
		MovementRows:         movementRows,
		MovementSummary:      summary,
		RefusalRows:          refusals,
		Handoff:              handoff,
		OutputSafety: OutputSafety{
			VerifiedAbsentSensitiveClasses: SensitiveClasses(),
		},
	}, nil
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
	var manifest DigestManifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}
	if manifest.SchemaVersion != DigestManifestSchemaVersion {
		return "", fmt.Errorf("unsupported digest manifest schema")
	}
	for _, artifact := range manifest.Artifacts {
		if artifact.Role != "query_pack_result" {
			continue
		}
		if unsafeSelectionPath(artifact.Path) || artifact.Path != filepathBase(queryPackPath) {
			return "", errUnsafePath
		}
		payload, err := os.ReadFile(queryPackPath)
		if err != nil {
			return "", err
		}
		sum := sha256.Sum256(payload)
		actual := hex.EncodeToString(sum[:])
		if artifact.SHA256 != actual {
			return "", errDigestMismatch
		}
		return actual, nil
	}
	return "", errMissingRequired
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
	switch metricID {
	case "missing_telemetry_rows":
		return row.EvidenceState == query.RowStateMissingTelemetry
	case "not_assessed_rows":
		return row.EvidenceState == query.RowStateNotAssessed
	case "cannot_verify_rows":
		return row.EvidenceState == query.RowStateCannotVerify
	case "unsupported_observer_rows":
		return row.EvidenceState == query.RowStateUnsupported || (hasSignal && signal.ObserverState == "unsupported")
	case "not_integrated_rows":
		return row.EvidenceState == query.RowStateNotIntegrated
	case "retention_limited_rows":
		return row.EvidenceState == query.RowStateRetentionLimited
	case "local_only_evidence_rows":
		return hasSignal && signal.WitnessScope == "local_only"
	case "ci_witnessed_evidence_rows":
		return hasSignal && signal.WitnessScope == "ci_witnessed"
	case "external_witnessed_evidence_rows":
		return hasSignal && signal.WitnessScope == "external_witnessed"
	case "issue_observed_rows":
		return row.EvidenceState == query.RowStateIssueObserved
	case "override_rows":
		return hasSignal && signal.OverrideMarker == "override_present"
	case "late_attach_rows":
		return hasSignal && signal.LateAttachMarker == "late_attach_observed"
	case "contract_change_rows":
		return hasSignal && signal.ContractChangeMarker == "contract_change_observed"
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

func unsafeOutput(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "http://") ||
		strings.Contains(lower, "https://") ||
		strings.Contains(lower, "secret") ||
		(strings.Contains(lower, "token") && !strings.Contains(lower, "credential_or_token")) ||
		(strings.Contains(lower, "credential") && !strings.Contains(lower, "credential_or_token")) ||
		strings.Contains(lower, "@") ||
		strings.Contains(value, "/") ||
		strings.Contains(value, "\\")
}

func unsafeLabel(value string) bool {
	lower := strings.ToLower(value)
	return unsafeOutput(value) ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "credential")
}

func unsafePath(value string) bool {
	return strings.Contains(value, "://") || strings.Contains(value, "..") || strings.HasPrefix(value, "/")
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
