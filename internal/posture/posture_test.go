package posture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestBuildAggregatesMovementAndRefusals(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	current := writeQueryPack(t, ".", "current", "missing_telemetry")
	previous := writeQueryPack(t, ".", "previous", "present")
	bad := writeQueryPack(t, ".", "bad", "issue_observed")
	currentDigest := writeDigest(t, current)
	previousDigest := writeDigest(t, previous)
	badDigest := writeDigest(t, bad)
	rewriteDigestSHA(t, badDigest, previousDigest)
	currentSignals := writeSignals(t, ".", "current-signals.json", "timeline.0001", "ci_witnessed", "override_present")
	selection := SelectionManifest{
		SchemaVersion:           "block21-cross-repo-selection-v1",
		ProfileID:               ProfileID,
		ProfileVersion:          ProfileVer,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories: []RepositoryWindow{
			selectionRepo("current", "2026-w02", current, currentDigest, currentSignals),
			selectionRepo("previous", "2026-w01", previous, previousDigest, ""),
			selectionRepo("bad-digest", "2026-w02", bad, badDigest, ""),
		},
		Handoff: map[string]string{"consumer": "sdp-report"},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)

	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.RefusalRows) != 1 || result.RefusalRows[0].RefusalReason != "untrusted_input_digest_mismatch" {
		t.Fatalf("refusals = %+v", result.RefusalRows)
	}
	missing := findMetric(result, "missing_telemetry_rows", "2026-w02")
	if missing.Numerator != 1 || missing.Denominator != 2 {
		t.Fatalf("missing metric = %+v", missing)
	}
	override := findMetric(result, "override_rows", "2026-w02")
	if override.Numerator != 1 || override.NotAssessedCount != 1 {
		t.Fatalf("override metric = %+v", override)
	}
	if result.MovementSummary.ComparableCount == 0 {
		t.Fatalf("missing comparable movement: %+v", result.MovementSummary)
	}
}

func TestBuildRefusesUnsafeLabels(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	manifest := writeDigest(t, qp)
	selection := SelectionManifest{
		SchemaVersion:           "block21-cross-repo-selection-v1",
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories: []RepositoryWindow{
			{
				InputID:                "unsafe",
				Repo:                   "https://provider.example/private",
				Team:                   "platform",
				Service:                "api",
				Harness:                "generic",
				ChangeType:             "feature",
				TimeWindow:             "2026-w02",
				InputObservedAt:        "2026-01-05T00:00:00Z",
				QueryPackResult:        qp,
				ArtifactDigestManifest: manifest,
			},
		},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.RefusalRows) != 1 || result.RefusalRows[0].RefusalReason != "unsafe_label" {
		t.Fatalf("refusals = %+v", result.RefusalRows)
	}
	rendered, err := Explain(result)
	if err != nil {
		t.Fatalf("explain: %v", err)
	}
	if containsUnsafe(rendered) {
		t.Fatalf("unsafe explain output: %s", rendered)
	}
}

func TestBuildSupportsTeamServiceAndHarnessChangeGrouping(t *testing.T) {
	for name, grouping := range map[string]string{
		"team-service":   GroupingTeamServiceWindow,
		"harness-change": GroupingHarnessChangeWindow,
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			withChdir(t, root)
			current := writeQueryPack(t, ".", "current", "present")
			previous := writeQueryPack(t, ".", "previous", "present")
			selection := SelectionManifest{
				SchemaVersion:           SelectionSchemaVersion,
				ProfileID:               ProfileID,
				GroupingSetID:           grouping,
				FreshnessBoundary:       "2026-01-01T00:00:00Z",
				DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
				CurrentWindow:           "2026-w02",
				PreviousWindow:          "2026-w01",
				Repositories: []RepositoryWindow{
					selectionRepo("current", "2026-w02", current, writeDigest(t, current), ""),
					selectionRepo("previous", "2026-w01", previous, writeDigest(t, previous), ""),
				},
			}
			selectionPath := filepath.Join(root, "selection.json")
			writeJSON(t, selectionPath, selection)
			result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
			if err != nil {
				t.Fatalf("build: %v", err)
			}
			if result.MovementSummary.ComparableCount == 0 {
				t.Fatalf("missing comparable movement for %s", grouping)
			}
		})
	}
}

func TestBuildRefusesStaleInput(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	repo := selectionRepo("current", "2026-w02", qp, writeDigest(t, qp), "")
	repo.InputObservedAt = "2025-12-31T00:00:00Z"
	selection := SelectionManifest{
		SchemaVersion:           SelectionSchemaVersion,
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories:            []RepositoryWindow{repo},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.RefusalRows) != 1 || result.RefusalRows[0].RefusalReason != "stale_input" {
		t.Fatalf("refusals = %+v", result.RefusalRows)
	}
}

func TestBuildRejectsGroupingOutsideExposurePolicy(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	selection := SelectionManifest{
		SchemaVersion:           SelectionSchemaVersion,
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingTeamServiceWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"team"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories:            []RepositoryWindow{selectionRepo("current", "2026-w02", qp, writeDigest(t, qp), "")},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	if _, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)); err == nil {
		t.Fatalf("expected grouping outside exposure policy to be rejected")
	}
}

func TestBuildRejectsDurationFreshnessBoundaryInV1(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	manifest := writeDigest(t, qp)
	selection := SelectionManifest{
		SchemaVersion:           "block21-cross-repo-selection-v1",
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "P7D",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories: []RepositoryWindow{
			selectionRepo("current", "2026-w02", qp, manifest, ""),
		},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	if _, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)); err == nil {
		t.Fatalf("expected P-duration freshness boundary to be rejected")
	}
}

func TestBuildRefusesUnsafeInputIDTimeWindowAndPaths(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	manifest := writeDigest(t, qp)
	for name, mutate := range map[string]func(*RepositoryWindow){
		"input-id": func(repo *RepositoryWindow) { repo.InputID = "credential-token" },
		"window":   func(repo *RepositoryWindow) { repo.TimeWindow = "2026/w02" },
		"path":     func(repo *RepositoryWindow) { repo.QueryPackResult = "../outside.json" },
	} {
		t.Run(name, func(t *testing.T) {
			repo := selectionRepo("current", "2026-w02", qp, manifest, "")
			mutate(&repo)
			selection := SelectionManifest{
				SchemaVersion:           SelectionSchemaVersion,
				ProfileID:               ProfileID,
				GroupingSetID:           GroupingRepoWindow,
				FreshnessBoundary:       "2026-01-01T00:00:00Z",
				DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
				CurrentWindow:           "2026-w02",
				PreviousWindow:          "2026-w01",
				Repositories:            []RepositoryWindow{repo},
			}
			selectionPath := filepath.Join(root, name+"-selection.json")
			writeJSON(t, selectionPath, selection)
			result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
			if err != nil {
				t.Fatalf("build: %v", err)
			}
			if len(result.RefusalRows) != 1 {
				t.Fatalf("refusals = %+v", result.RefusalRows)
			}
		})
	}
}

func TestBuildRejectsMismatchedDigestManifestPath(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	manifest := writeDigest(t, qp)
	var digest DigestManifest
	payload, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(payload, &digest); err != nil {
		t.Fatal(err)
	}
	digest.Artifacts[0].Path = "other-query-pack.json"
	writeJSON(t, manifest, digest)
	selection := SelectionManifest{
		SchemaVersion:           SelectionSchemaVersion,
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories:            []RepositoryWindow{selectionRepo("current", "2026-w02", qp, manifest, "")},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.RefusalRows) != 1 || result.RefusalRows[0].RefusalReason != "malformed_input" {
		t.Fatalf("refusals = %+v", result.RefusalRows)
	}
}

func TestBuildRejectsUnsupportedDigestManifestSchema(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	manifest := writeDigest(t, qp)
	var digest DigestManifest
	readJSONFixture(t, manifest, &digest)
	digest.SchemaVersion = "future-schema"
	writeJSON(t, manifest, digest)
	selection := SelectionManifest{
		SchemaVersion:           SelectionSchemaVersion,
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories:            []RepositoryWindow{selectionRepo("current", "2026-w02", qp, manifest, "")},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.InputSelection) != 1 || result.InputSelection[0].InputTrustState == "trusted_input" {
		t.Fatalf("unsupported digest schema trusted input: %+v", result.InputSelection)
	}
}

func TestBuildNormalizesHandoffAndRejectsUnsafeHandoff(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	manifest := writeDigest(t, qp)
	selection := SelectionManifest{
		SchemaVersion:           SelectionSchemaVersion,
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories:            []RepositoryWindow{selectionRepo("current", "2026-w02", qp, manifest, "")},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if result.Handoff == nil || len(result.Handoff) != 0 {
		t.Fatalf("handoff not normalized: %+v", result.Handoff)
	}
	selection.Handoff = map[string]string{"consumer": "https://provider.example/private"}
	writeJSON(t, selectionPath, selection)
	if _, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)); err == nil {
		t.Fatalf("expected unsafe handoff to be rejected")
	}
}

func TestBuildRefusesWindowsAbsoluteInputPath(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, ".", "current", "present")
	repo := selectionRepo("current", "2026-w02", qp, writeDigest(t, qp), "")
	repo.QueryPackResult = `C:\private\query-pack.json`
	selection := SelectionManifest{
		SchemaVersion:           SelectionSchemaVersion,
		ProfileID:               ProfileID,
		GroupingSetID:           GroupingRepoWindow,
		FreshnessBoundary:       "2026-01-01T00:00:00Z",
		DimensionExposurePolicy: []string{"repo", "team", "service", "harness", "change_type"},
		CurrentWindow:           "2026-w02",
		PreviousWindow:          "2026-w01",
		Repositories:            []RepositoryWindow{repo},
	}
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.RefusalRows) != 1 || result.RefusalRows[0].RefusalReason != "malformed_input" {
		t.Fatalf("refusals = %+v", result.RefusalRows)
	}
}

func TestSchemaMirrorsPostureEnums(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("..", "..", "schema", "cross-repo-posture-export.schema.json"))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var schema map[string]any
	if err := json.Unmarshal(payload, &schema); err != nil {
		t.Fatalf("decode schema: %v", err)
	}
	properties := schemaObjectAt(t, schema, "properties")
	assertSchemaConst(t, properties, "schema_version", SchemaVersion)
	assertSchemaConst(t, properties, "export_profile_id", ProfileID)
	defs := schemaObjectAt(t, schema, "$defs")
	metricID := schemaObjectAt(t, defs, "metricID")
	assertSchemaEnumValue(t, metricID, []string{
		"missing_telemetry_rows",
		"not_assessed_rows",
		"cannot_verify_rows",
		"unsupported_observer_rows",
		"not_integrated_rows",
		"retention_limited_rows",
		"local_only_evidence_rows",
		"ci_witnessed_evidence_rows",
		"external_witnessed_evidence_rows",
		"issue_observed_rows",
		"override_rows",
		"late_attach_rows",
		"contract_change_rows",
	})
	outputSafety := schemaObjectAt(t, defs, "outputSafety")
	outputSafetyProps := schemaObjectAt(t, outputSafety, "properties")
	classes := schemaObjectAt(t, outputSafetyProps, "verified_absent_sensitive_classes")
	items := schemaObjectAt(t, classes, "items")
	assertSchemaEnumValue(t, items, SensitiveClasses())
}

func selectionRepo(inputID, window, qp, digestManifest, signals string) RepositoryWindow {
	return RepositoryWindow{
		InputID:                inputID,
		Repo:                   "repo-a",
		Team:                   "platform",
		Service:                "api",
		Harness:                "generic",
		ChangeType:             "feature",
		TimeWindow:             window,
		InputObservedAt:        "2026-01-05T00:00:00Z",
		QueryPackResult:        qp,
		ArtifactDigestManifest: digestManifest,
		PostureSignalManifest:  signals,
	}
}

func writeQueryPack(t *testing.T, root, name, rowState string) string {
	t.Helper()
	path := filepath.Join(root, name+"-query-pack.json")
	result := map[string]any{
		"schema_version":     "block20-forensics-query-pack-result-v1",
		"query_pack_id":      "forensics-basic-v1",
		"query_pack_version": "v1",
		"input_artifacts": []map[string]any{
			{"role": "run", "path_redacted_id": "run", "artifact_required": true},
		},
		"query_rows": map[string]any{
			"forensics-summary": []map[string]any{},
			"forensics-timeline": []map[string]any{
				{"id": "timeline.0001", "query": "forensics-timeline", "evidence_state": rowState, "evidence_family": "command", "source_ref": "block_09.event.command.e0001"},
				{"id": "timeline.0002", "query": "forensics-timeline", "evidence_state": "not_assessed", "evidence_family": "test", "source_ref": "block_09.event.test.e0002"},
			},
			"forensics-gaps":              []map[string]any{},
			"forensics-redactions":        []map[string]any{},
			"forensics-capture-depth":     []map[string]any{},
			"forensics-unverified-claims": []map[string]any{},
		},
		"output_safety": map[string]any{"verified_absent_sensitive_classes": []string{"tokens"}},
	}
	writeJSON(t, path, result)
	return path
}

func writeDigest(t *testing.T, artifact string) string {
	t.Helper()
	payload, err := os.ReadFile(artifact)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(payload)
	path := artifact + ".digest.json"
	writeJSON(t, path, DigestManifest{
		SchemaVersion: "block21-artifact-digest-manifest-v1",
		Artifacts: []DigestArtifact{{
			Role:   "query_pack_result",
			Path:   filepath.Base(artifact),
			SHA256: hex.EncodeToString(sum[:]),
		}},
	})
	return path
}

func rewriteDigestSHA(t *testing.T, targetManifest, sourceManifest string) {
	t.Helper()
	var target DigestManifest
	var source DigestManifest
	readJSONFixture(t, targetManifest, &target)
	readJSONFixture(t, sourceManifest, &source)
	target.Artifacts[0].SHA256 = source.Artifacts[0].SHA256
	writeJSON(t, targetManifest, target)
}

func readJSONFixture(t *testing.T, path string, target any) {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(payload, target); err != nil {
		t.Fatal(err)
	}
}

func writeSignals(t *testing.T, root, name, rowRef, witnessScope, override string) string {
	t.Helper()
	path := filepath.Join(root, name)
	writeJSON(t, path, SignalManifest{
		SchemaVersion: "block21-posture-signal-manifest-v1",
		Signals: []PostureSignal{{
			RowRef:         rowRef,
			WitnessScope:   witnessScope,
			OverrideMarker: override,
		}},
	})
	return path
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func findMetric(result ExportResult, id, window string) MetricRow {
	for _, row := range result.MetricRows {
		if row.MetricID == id && row.TimeWindow == window {
			return row
		}
	}
	return MetricRow{}
}

func containsUnsafe(value string) bool {
	return strings.Contains(value, "https://") || strings.Contains(value, "/private") || strings.Contains(value, "secret")
}

func withChdir(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
}

func schemaObjectAt(t *testing.T, source map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := source[key].(map[string]any)
	if !ok {
		t.Fatalf("%s is not an object", key)
	}
	return value
}

func assertSchemaConst(t *testing.T, properties map[string]any, key string, expected string) {
	t.Helper()
	property := schemaObjectAt(t, properties, key)
	actual, ok := property["const"].(string)
	if !ok || actual != expected {
		t.Fatalf("%s.const = %v expected %q", key, property["const"], expected)
	}
}

func assertSchemaEnumValue(t *testing.T, property map[string]any, expected []string) {
	t.Helper()
	rawEnum, ok := property["enum"].([]any)
	if !ok {
		t.Fatalf("enum is not an array: %+v", property)
	}
	actual := make([]string, 0, len(rawEnum))
	for _, value := range rawEnum {
		asString, ok := value.(string)
		if !ok {
			t.Fatalf("enum contains non-string: %+v", rawEnum)
		}
		actual = append(actual, asString)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("enum = %+v expected %+v", actual, expected)
	}
}
