package posture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func TestVerifyDigestManifest(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)

	queryPack := writeQueryPack(t, "current", "present")
	digestManifest := writeDigest(t, queryPack)

	actual, err := verifyDigestManifest(digestManifest, queryPack)
	if err != nil {
		t.Fatalf("verifyDigestManifest: %v", err)
	}

	payload, err := os.ReadFile(queryPack)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(payload)
	want := hex.EncodeToString(sum[:])
	if actual != want {
		t.Fatalf("digest = %s, want %s", actual, want)
	}
}

func TestVerifyDigestManifestMissingQueryPackArtifact(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)

	queryPack := writeQueryPack(t, "current", "present")
	digestManifest := writeDigest(t, queryPack)
	var digest DigestManifest
	readJSONFixture(t, digestManifest, &digest)
	digest.Artifacts = nil
	writeJSON(t, digestManifest, digest)

	_, err := verifyDigestManifest(digestManifest, queryPack)
	if !errors.Is(err, errMissingRequired) {
		t.Fatalf("expected errMissingRequired, got %v", err)
	}
}

func TestVerifyDigestManifestRejectsMismatchedQueryPackPath(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)

	queryPack := writeQueryPack(t, "current", "present")
	digestManifest := writeDigest(t, queryPack)
	var digest DigestManifest
	readJSONFixture(t, digestManifest, &digest)
	digest.Artifacts[0].Path = "other-query-pack.json"
	writeJSON(t, digestManifest, digest)

	_, err := verifyDigestManifest(digestManifest, queryPack)
	if !errors.Is(err, errUnsafePath) {
		t.Fatalf("expected errUnsafePath, got %v", err)
	}
}

func TestVerifyDigestManifestMismatchReturnsNoDigest(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)

	queryPack := writeQueryPack(t, "current", "present")
	digestManifest := writeDigest(t, queryPack)
	var digest DigestManifest
	readJSONFixture(t, digestManifest, &digest)
	digest.Artifacts[0].SHA256 = strings.Repeat("0", 64)
	writeJSON(t, digestManifest, digest)

	actual, err := verifyDigestManifest(digestManifest, queryPack)
	if !errors.Is(err, errDigestMismatch) {
		t.Fatalf("expected errDigestMismatch, got %v", err)
	}
	if actual != "" {
		t.Fatalf("expected empty digest on mismatch, got %q", actual)
	}
}

func TestIngestRepository(t *testing.T) {
	t.Run("trusted", func(t *testing.T) {
		root := t.TempDir()
		withChdir(t, root)

		queryPack := writeQueryPack(t, "current", "present")
		digest := writeDigest(t, queryPack)
		signals := writeSignals(t, ".", "current-signals.json", "timeline.0001", "ci_witnessed", "override_present")
		repo := selectionRepo("current", "2026-w02", queryPack, digest, signals)

		result := ingestRepository(repo, time.Time{}, false)
		if !result.trusted {
			t.Fatalf("unexpected refusal: %+v", result)
		}
		if result.digest == "" {
			t.Fatalf("expected digest")
		}
		if result.result.SchemaVersion != query.QueryPackSchemaVersion {
			t.Fatalf("unexpected schema_version: %q", result.result.SchemaVersion)
		}
	})

	t.Run("unsafe label", func(t *testing.T) {
		root := t.TempDir()
		withChdir(t, root)

		queryPack := writeQueryPack(t, "current", "present")
		repo := selectionRepo("current", "2026-w02", queryPack, writeDigest(t, queryPack), "")
		repo.Repo = "https://provider.example/private"

		result := ingestRepository(repo, time.Time{}, false)
		if result.trusted || result.refusalReason != "unsafe_label" || result.inputTrustState != "cannot_verify_input" {
			t.Fatalf("unexpected result: %+v", result)
		}
		if result.recordSelection {
			t.Fatalf("expected recordSelection false")
		}
	})

	t.Run("stale input", func(t *testing.T) {
		root := t.TempDir()
		withChdir(t, root)

		queryPack := writeQueryPack(t, "current", "present")
		repo := selectionRepo("current", "2026-w02", queryPack, writeDigest(t, queryPack), "")
		repo.InputObservedAt = "2025-12-31T00:00:00Z"

		result := ingestRepository(repo, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), true)
		if result.trusted || result.refusalReason != "stale_input" || result.inputTrustState != "stale_input" {
			t.Fatalf("unexpected result: %+v", result)
		}
		if !result.recordSelection {
			t.Fatalf("expected recordSelection true")
		}
	})

	t.Run("unsupported query pack", func(t *testing.T) {
		root := t.TempDir()
		withChdir(t, root)

		queryPack := writeQueryPack(t, "current", "present")
		var payload map[string]any
		readJSONFixture(t, queryPack, &payload)
		payload["schema_version"] = "future"
		writeJSON(t, queryPack, payload)
		digest := writeDigest(t, queryPack)
		repo := selectionRepo("current", "2026-w02", queryPack, digest, "")

		result := ingestRepository(repo, time.Time{}, false)
		if result.trusted || result.refusalReason != "malformed_input" || result.inputTrustState != "cannot_verify_input" {
			t.Fatalf("unexpected result: %+v", result)
		}
		if result.digest == "" {
			t.Fatalf("expected digest on malformed query pack")
		}
	})

	t.Run("malformed signals", func(t *testing.T) {
		root := t.TempDir()
		withChdir(t, root)

		queryPack := writeQueryPack(t, "current", "present")
		repo := selectionRepo("current", "2026-w02", queryPack, writeDigest(t, queryPack), "malformed-signals.json")
		if err := os.WriteFile(repo.PostureSignalManifest, []byte("{invalid json"), 0o644); err != nil {
			t.Fatal(err)
		}

		result := ingestRepository(repo, time.Time{}, false)
		if result.trusted || result.refusalReason != "malformed_input" || result.inputTrustState != "cannot_verify_input" {
			t.Fatalf("unexpected result: %+v", result)
		}
		if result.digest == "" {
			t.Fatalf("expected digest on malformed signals")
		}
	})
}

func TestBuildAggregatesMovementAndRefusals(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	current := writeQueryPack(t, "current", "missing_telemetry")
	previous := writeQueryPack(t, "previous", "present")
	bad := writeQueryPack(t, "bad", "issue_observed")
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

func TestPostureHelperCoverageForMovementAndSafety(t *testing.T) {
	validMovement := MovementRow{
		ID:                  "movement.0001",
		MetricID:            "missing_telemetry_rows",
		MetricVersion:       ProfileVer,
		DimensionKey:        "repo:current",
		CurrentValue:        1,
		PreviousValue:       0,
		ComparisonBasis:     "previous_window",
		Comparable:          true,
		NonComparableReason: "",
	}
	if malformedMovementIdentity(validMovement) {
		t.Fatalf("valid movement identity reported malformed")
	}
	invalidMovement := validMovement
	invalidMovement.MetricVersion = "future"
	if !malformedMovementIdentity(invalidMovement) {
		t.Fatalf("invalid movement identity reported valid")
	}
	if malformedMovementSummaryReasons(map[string]int{"non_comparable_missing_window": 1}) {
		t.Fatalf("valid movement summary reason reported malformed")
	}
	if !malformedMovementSummaryReasons(map[string]int{"future_reason": 1}) {
		t.Fatalf("unknown movement summary reason reported valid")
	}
	if err := malformedRowError(false, "bad"); err != nil {
		t.Fatalf("unexpected row error: %v", err)
	}
	if err := malformedRowError(true, "bad"); err == nil {
		t.Fatalf("expected row error")
	}
}

func TestPostureHelperCoverageForMetricAndPathBoundaries(t *testing.T) {
	row := query.QueryPackRow{EvidenceState: query.RowStateUnsupported}
	if !unsupportedObserverMetricMatches(row, PostureSignal{}, false) {
		t.Fatalf("unsupported row did not match unsupported observer metric")
	}
	signal := PostureSignal{ObserverState: "unsupported"}
	if !unsupportedObserverMetricMatches(query.QueryPackRow{}, signal, true) {
		t.Fatalf("unsupported signal did not match unsupported observer metric")
	}
	if unsupportedObserverMetricMatches(query.QueryPackRow{}, signal, false) {
		t.Fatalf("signal matched without signal presence")
	}
	if unsafeInputLabel(RepositoryWindow{InputID: "input-1", TimeWindow: "2026-w02"}) {
		t.Fatalf("safe input label reported unsafe")
	}
	if !unsafeInputLabel(RepositoryWindow{InputID: "https://example.invalid/input", TimeWindow: "2026-w02"}) {
		t.Fatalf("unsafe input label reported safe")
	}
	for _, value := range []string{"https://example.invalid/query.json", "C:/repo/query.json", "/repo/query.json", "../query.json"} {
		if !hasUnsafeSelectionPathPrefix(value) {
			t.Fatalf("unsafe selection path prefix reported safe: %s", value)
		}
	}
	if hasUnsafeSelectionPathPrefix("relative/query.json") {
		t.Fatalf("relative selection path reported unsafe")
	}
}

func TestPostureHelperCoverageForDigestAndTrustCopies(t *testing.T) {
	if got := shortDigest(strings.Repeat("a", 64)); got != strings.Repeat("a", 16) {
		t.Fatalf("short digest = %q", got)
	}
	if got := shortDigest("abc"); got != "not_assessed0000" {
		t.Fatalf("short missing digest = %q", got)
	}
	source := map[string]int{"trusted": 1}
	copied := copyTrust(source)
	copied["trusted"] = 2
	if source["trusted"] != 1 {
		t.Fatalf("copyTrust aliased source map")
	}
}

func TestBuildRefusesUnsafeLabels(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, "current", "present")
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
			current := writeQueryPack(t, "current", "present")
			previous := writeQueryPack(t, "previous", "present")
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
	assertBuildRefusesSingleRepo(t, "stale_input", func(repo *RepositoryWindow) {
		repo.InputObservedAt = "2025-12-31T00:00:00Z"
	})
}

func TestBuildRejectsGroupingOutsideExposurePolicy(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, "current", "present")
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
	qp := writeQueryPack(t, "current", "present")
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
	qp := writeQueryPack(t, "current", "present")
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
	qp := writeQueryPack(t, "current", "present")
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
	qp := writeQueryPack(t, "current", "present")
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

func TestUnsafeOutput(t *testing.T) {
	cases := []struct {
		name   string
		value  string
		unsafe bool
	}{
		{"safe value", "repository-a", false},
		{"https url", "https://provider.example", true},
		{"http url", "http://provider.example", true},
		{"secret marker", "contains secret", true},
		{"email", "user@example.com", true},
		{"slash path", "path/to/file", true},
		{"windows path", "path\\to\\file", true},
		{"token", "api-token", true},
		{"credential", "api credential", true},
		{"credential and token exact", "credential_or_token", false},
		{"token with exception", "Credential_OR_Token", false},
		{"mixed case token", "A Token Value", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			if got := unsafeOutput(tc.value); got != tc.unsafe {
				t.Fatalf("unsafeOutput(%q)=%t expected %t", tc.value, got, tc.unsafe)
			}
		})
	}
}

func TestBuildNormalizesHandoffAndRejectsUnsafeHandoff(t *testing.T) {
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, "current", "present")
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
	assertBuildRefusesSingleRepo(t, "malformed_input", func(repo *RepositoryWindow) {
		repo.QueryPackResult = `C:\private\query-pack.json`
	})
}

func assertBuildRefusesSingleRepo(t *testing.T, wantReason string, mutate func(*RepositoryWindow)) {
	t.Helper()
	root := t.TempDir()
	withChdir(t, root)
	qp := writeQueryPack(t, "current", "present")
	repo := selectionRepo("current", "2026-w02", qp, writeDigest(t, qp), "")
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
	selectionPath := filepath.Join(root, "selection.json")
	writeJSON(t, selectionPath, selection)
	result, err := Build(selectionPath, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.RefusalRows) != 1 || result.RefusalRows[0].RefusalReason != wantReason {
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

func TestValidateExportResultAcceptsCanonicalResult(t *testing.T) {
	result := validExportResult()
	if err := ValidateExportResult(result); err != nil {
		t.Fatalf("validate export result: %v", err)
	}
}

func TestValidateExportResultRejectsMalformedNestedRows(t *testing.T) {
	for name, mutate := range map[string]func(*ExportResult){
		"metric-dimension": func(result *ExportResult) {
			result.MetricRows[0].Dimensions["repo"] = "https://provider.example/private"
		},
		"metric-trust-summary": func(result *ExportResult) {
			result.MetricRows[0].InputTrustStateSummary["trusted_input"] = -1
		},
		"movement-summary": func(result *ExportResult) {
			result.MovementSummary.NonComparableReason["unknown"] = 1
		},
		"movement-row": func(result *ExportResult) {
			result.MovementRows[0].CurrentValue = -1
		},
		"refusal-reason": func(result *ExportResult) {
			result.RefusalRows[0].RefusalReason = "ok"
		},
	} {
		t.Run(name, func(t *testing.T) {
			result := validExportResult()
			mutate(&result)
			if err := ValidateExportResult(result); err == nil {
				t.Fatalf("expected malformed %s to be rejected", name)
			}
		})
	}
}

func TestValidateExportResultRejectsMalformedCollections(t *testing.T) {
	for name, mutate := range map[string]func(*ExportResult){
		"grouping": func(result *ExportResult) {
			result.GroupingSetID = "unknown"
		},
		"missing-collection": func(result *ExportResult) {
			result.Handoff = nil
		},
		"output-safety": func(result *ExportResult) {
			result.OutputSafety.VerifiedAbsentSensitiveClasses = nil
		},
	} {
		t.Run(name, func(t *testing.T) {
			result := validExportResult()
			mutate(&result)
			if err := ValidateExportResult(result); err == nil {
				t.Fatalf("expected malformed %s to be rejected", name)
			}
		})
	}
}

func TestValidateMovementRowRejectsMalformedRows(t *testing.T) {
	base := validExportResult().MovementRows[0]
	for name, mutate := range map[string]func(*MovementRow){
		"identity": func(row *MovementRow) {
			row.ID = ""
		},
		"metric-id": func(row *MovementRow) {
			row.MetricID = "unknown"
		},
		"version": func(row *MovementRow) {
			row.MetricVersion = "old"
		},
		"dimension-key": func(row *MovementRow) {
			row.DimensionKey = ""
		},
		"current-value": func(row *MovementRow) {
			row.CurrentValue = -1
		},
		"previous-value": func(row *MovementRow) {
			row.PreviousValue = -1
		},
		"comparison-basis": func(row *MovementRow) {
			row.ComparisonBasis = "unsupported"
		},
		"non-comparable-reason": func(row *MovementRow) {
			row.Comparable = false
			row.NonComparableReason = "output_safety_violation"
		},
	} {
		t.Run(name, func(t *testing.T) {
			row := base
			mutate(&row)
			if err := validateMovementRow(row); err == nil || err.Error() != "malformed posture export movement_row" {
				t.Fatalf("error = %v, want malformed posture export movement_row", err)
			}
		})
	}
}

func TestValidateMetricRowShapeRejectsMalformedRows(t *testing.T) {
	base := validExportResult().MetricRows[0]
	for name, mutate := range map[string]func(*MetricRow){
		"identity": func(row *MetricRow) {
			row.ID = ""
		},
		"metric-id": func(row *MetricRow) {
			row.MetricID = "unknown"
		},
		"version": func(row *MetricRow) {
			row.MetricVersion = "old"
		},
		"count": func(row *MetricRow) {
			row.Numerator = -1
		},
		"denominator": func(row *MetricRow) {
			row.Denominator = -1
		},
		"unit": func(row *MetricRow) {
			row.Unit = "bytes"
		},
		"not-assessed-count": func(row *MetricRow) {
			row.NotAssessedCount = -1
		},
		"time-window": func(row *MetricRow) {
			row.TimeWindow = "https://provider.example/private"
		},
		"dimensions": func(row *MetricRow) {
			row.Dimensions = nil
		},
		"source": func(row *MetricRow) {
			row.SourceInputRefs = nil
		},
		"source-state": func(row *MetricRow) {
			row.SourceFieldState = "unknown"
		},
		"trust-summary": func(row *MetricRow) {
			row.InputTrustStateSummary = nil
		},
	} {
		t.Run(name, func(t *testing.T) {
			row := base
			mutate(&row)
			if err := validateMetricRowShape(row); err == nil || err.Error() != "malformed posture export metric_row" {
				t.Fatalf("error = %v, want malformed metric_row", err)
			}
		})
	}
}

func validExportResult() ExportResult {
	return ExportResult{
		SchemaVersion:        SchemaVersion,
		ExportProfileID:      ProfileID,
		ExportProfileVersion: ProfileVer,
		ExportID:             "export.0001",
		Producer:             "sdp-trace",
		GeneratedAt:          "2026-01-10T00:00:00Z",
		GroupingSetID:        GroupingRepoWindow,
		ActiveGroupingKeys:   []string{"repo", "time_window"},
		InputSelection: []InputSelection{{
			InputID:         "input-a",
			Repository:      "repo-a",
			TimeWindow:      "2026-w02",
			PathRedactedID:  "query-pack",
			SHA256:          strings.Repeat("a", 64),
			InputTrustState: "trusted_input",
		}},
		MetricRows: []MetricRow{{
			ID:                      "metric.0001",
			MetricID:                "missing_telemetry_rows",
			MetricVersion:           ProfileVer,
			Numerator:               1,
			Denominator:             2,
			Unit:                    "rows",
			TimeWindow:              "2026-w02",
			Dimensions:              map[string]string{"repo": "repo-a", "time_window": "2026-w02"},
			DimensionKey:            "repo=repo-a|time_window=2026-w02",
			SourceInputRefs:         []string{"input-a"},
			SourceArtifactDigestSet: strings.Repeat("b", 64),
			SourceFieldState:        "present",
			NotAssessedCount:        0,
			InputTrustStateSummary:  map[string]int{"trusted_input": 1},
		}},
		MovementRows: []MovementRow{{
			ID:                   "movement.0001",
			MetricID:             "missing_telemetry_rows",
			MetricVersion:        ProfileVer,
			DimensionKey:         "repo=repo-a|time_window=2026-w02",
			CurrentMetricRowRef:  "metric.0001",
			PreviousMetricRowRef: "metric.0000",
			CurrentValue:         1,
			PreviousValue:        0,
			Delta:                1,
			ComparisonBasis:      "same_profile_metric_dimension_window",
			Comparable:           true,
		}},
		MovementSummary: MovementSummary{
			ComparableCount:     1,
			NonComparableCount:  0,
			NonComparableReason: map[string]int{},
		},
		RefusalRows: []RefusalRow{{
			ID:              "refusal.0001",
			InputID:         "input-b",
			TimeWindow:      "2026-w02",
			RefusalReason:   "stale_input",
			InputTrustState: "stale_input",
		}},
		Handoff: map[string]string{"consumer": "sdp-report"},
		OutputSafety: OutputSafety{
			VerifiedAbsentSensitiveClasses: []string{"tokens"},
		},
	}
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

func writeQueryPack(t *testing.T, name, rowState string) string {
	t.Helper()
	path := name + "-query-pack.json"
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
