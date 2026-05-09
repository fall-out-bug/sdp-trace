package prreview

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildPacketBindsRefsAndRejectsUnsafeIdentity(t *testing.T) {
	root := t.TempDir()
	diffPath := writeText(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	contextPath := writeText(t, root, "spec.md", "# Spec\n")
	verificationPath := writeText(t, root, "verify.txt", "go test ./...\n")
	outDir := filepath.Join(root, "packet")

	packet, err := BuildPacket(PacketOptions{
		OutDir:            outDir,
		RepoID:            "demo_repo",
		ChangeRef:         "pr-123",
		BaseCommit:        forty("a"),
		HeadCommit:        forty("b"),
		DiffPath:          diffPath,
		ContextPaths:      []string{contextPath},
		VerificationPaths: []string{verificationPath},
		CIState:           StateNotAssessed,
		Now:               time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if packet.PacketDigest == "" || packet.DiffRef.DigestSHA256 == "" {
		t.Fatalf("packet missing digest: %+v", packet)
	}
	if len(packet.ContextRefs) != 1 || len(packet.VerificationRefs) != 1 {
		t.Fatalf("refs not captured: %+v", packet)
	}
	if len(packet.UnavailableFields) != 1 || packet.UnavailableFields[0].Field != "metadata_ref" {
		t.Fatalf("missing metadata should be explicit not_assessed, got %+v", packet.UnavailableFields)
	}
	if strings.Contains(packet.ContextRefs[0].Ref, root) || strings.HasPrefix(packet.ContextRefs[0].Ref, "/") {
		t.Fatalf("context ref leaked absolute path: %+v", packet.ContextRefs[0])
	}
	if _, err := os.Stat(filepath.Join(outDir, "packet.json")); err != nil {
		t.Fatalf("packet not written: %v", err)
	}

	_, err = BuildPacket(PacketOptions{
		OutDir:     filepath.Join(root, "unsafe"),
		RepoID:     "/tmp/private/repo",
		ChangeRef:  "pr-123",
		BaseCommit: forty("a"),
		HeadCommit: forty("b"),
		DiffPath:   diffPath,
	})
	if err == nil || !strings.Contains(err.Error(), "unsafe_repo_id") {
		t.Fatalf("expected unsafe repo id rejection, got %v", err)
	}
}

func TestValidateProfileRejectsMalformedProfiles(t *testing.T) {
	valid := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "default",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles:          []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}},
	}
	for name, tc := range map[string]struct {
		mutate  func(*ReviewProfile)
		wantErr string
	}{
		"bad-schema": {
			mutate:  func(profile *ReviewProfile) { profile.SchemaVersion = "unknown" },
			wantErr: "invalid_profile_schema_version: unknown",
		},
		"missing-profile-id": {
			mutate:  func(profile *ReviewProfile) { profile.ProfileID = "" },
			wantErr: "profile_requires_profile_id",
		},
		"missing-required-planes": {
			mutate:  func(profile *ReviewProfile) { profile.RequiredPlanes = nil },
			wantErr: "profile_requires_required_planes",
		},
		"missing-roles": {
			mutate:  func(profile *ReviewProfile) { profile.Roles = nil },
			wantErr: "profile_requires_roles",
		},
		"missing-role-fields": {
			mutate:  func(profile *ReviewProfile) { profile.Roles[0].Runner = "" },
			wantErr: "profile_role_requires_id_plane_runner",
		},
		"invalid-runner": {
			mutate:  func(profile *ReviewProfile) { profile.Roles[0].Runner = "unknown" },
			wantErr: "profile_role_invalid_runner: unknown",
		},
		"required-plane-without-role": {
			mutate: func(profile *ReviewProfile) {
				profile.RequiredPlanes = []string{PlaneCodeCorrectness, PlanePrivacySafety}
			},
			wantErr: "profile_required_plane_without_role: privacy_output_safety",
		},
	} {
		t.Run(name, func(t *testing.T) {
			profile := cloneReviewProfile(valid)
			tc.mutate(&profile)
			if err := validateProfile(profile); err == nil || err.Error() != tc.wantErr {
				t.Fatalf("validateProfile() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
	if err := validateProfile(valid); err != nil {
		t.Fatalf("valid profile rejected: %v", err)
	}
}

func TestCitationResolvableCharacterization(t *testing.T) {
	packet := Packet{
		DiffRef:          SafeRef{ID: "diff-ref"},
		ContextRefs:      []SafeRef{{ID: "spec"}},
		VerificationRefs: []SafeRef{{ID: "verify"}},
	}
	for name, tc := range map[string]struct {
		citation Citation
		want     bool
	}{
		"empty-citation": {
			citation: Citation{},
			want:     false,
		},
		"diff-ref-with-hunk": {
			citation: Citation{ContextRefID: "diff-ref", DiffHunkID: "hunk-1"},
			want:     true,
		},
		"diff-alias-with-digest": {
			citation: Citation{ContextRefID: "diff", SourceDigest: "sha256:abc"},
			want:     true,
		},
		"diff-ref-without-hunk-or-digest": {
			citation: Citation{ContextRefID: "diff-ref"},
			want:     false,
		},
		"context-ref-with-line": {
			citation: Citation{ContextRefID: "spec", LineStart: 12},
			want:     true,
		},
		"context-ref-without-location": {
			citation: Citation{ContextRefID: "spec"},
			want:     false,
		},
		"verification-ref-with-line": {
			citation: Citation{ContextRefID: "verify", LineStart: 4},
			want:     true,
		},
		"verification-ref-with-hunk-only": {
			citation: Citation{ContextRefID: "verify", DiffHunkID: "hunk-1"},
			want:     false,
		},
		"unknown-ref-with-digest": {
			citation: Citation{ContextRefID: "unknown", SourceDigest: "sha256:abc"},
			want:     true,
		},
		"digest-only": {
			citation: Citation{SourceDigest: "sha256:abc"},
			want:     true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := citationResolvable(packet, tc.citation); got != tc.want {
				t.Fatalf("citationResolvable() = %v, want %v for %+v", got, tc.want, tc.citation)
			}
		})
	}
}

func cloneReviewProfile(profile ReviewProfile) ReviewProfile {
	profile.RequiredPlanes = append([]string(nil), profile.RequiredPlanes...)
	profile.Roles = append([]ReviewRole(nil), profile.Roles...)
	return profile
}

func TestBuildPacketRecordsUnavailableInputsAndDigestChangesWithDiff(t *testing.T) {
	root := t.TempDir()
	diffPath := writeText(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	first, err := BuildPacket(PacketOptions{
		OutDir:     filepath.Join(root, "packet-1"),
		RepoID:     "demo_repo",
		ChangeRef:  "pr-123",
		BaseCommit: forty("a"),
		HeadCommit: forty("b"),
		DiffPath:   diffPath,
		CIState:    StateNotAssessed,
		Now:        time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	unavailable := map[string]string{}
	for _, field := range first.UnavailableFields {
		unavailable[field.Field] = field.State
	}
	for _, field := range []string{"metadata_ref", "context_refs", "verification_refs"} {
		if unavailable[field] != StateNotAssessed {
			t.Fatalf("missing %s should be not_assessed, got %+v", field, first.UnavailableFields)
		}
	}

	diffPath = writeText(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+newer\n")
	second, err := BuildPacket(PacketOptions{
		OutDir:     filepath.Join(root, "packet-2"),
		RepoID:     "demo_repo",
		ChangeRef:  "pr-123",
		BaseCommit: forty("a"),
		HeadCommit: forty("b"),
		DiffPath:   diffPath,
		CIState:    StateNotAssessed,
		Now:        time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if first.DiffRef.DigestSHA256 == second.DiffRef.DigestSHA256 || first.PacketDigest == second.PacketDigest {
		t.Fatalf("diff mutation should change diff and packet digests: first=%s/%s second=%s/%s", first.DiffRef.DigestSHA256, first.PacketDigest, second.DiffRef.DigestSHA256, second.PacketDigest)
	}
	canonical := second
	canonical.PacketDigest = "sha256:" + sixtyFour("f")
	replayed, err := packetDigest(canonical)
	if err != nil {
		t.Fatal(err)
	}
	if "sha256:"+replayed != second.PacketDigest {
		t.Fatalf("packet digest should be replayable with packet_digest cleared: got sha256:%s want %s", replayed, second.PacketDigest)
	}
}

func TestValidateReviewStatesAndAuthorityBoundary(t *testing.T) {
	packet := Packet{
		SchemaVersion: SchemaVersionPacket,
		PacketID:      "packet-1",
		PacketDigest:  "sha256:" + sixtyFour("1"),
		RepoID:        "demo_repo",
		ChangeRef:     "pr-123",
		BaseCommit:    forty("a"),
		HeadCommit:    forty("b"),
		DiffRef:       SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone},
		CIState:       StateNotAssessed,
		CreatedAt:     "2026-05-09T12:00:00Z",
		CreatedBy:     "test",
	}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "trust-sensitive-default",
		RequiredPlanes: []string{PlaneCodeCorrectness, PlaneTraceEvidence, PlaneRequirements},
		Roles: []ReviewRole{
			{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
			{RoleID: "trace", Plane: PlaneTraceEvidence, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
			{RoleID: "req", Plane: PlaneRequirements, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
		},
	}
	runs := RunSet{
		SchemaVersion: SchemaVersionRunSet,
		PacketDigest:  packet.PacketDigest,
		Results: []ReviewerResult{
			{ReviewRunID: "run-code", PacketDigest: packet.PacketDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings},
			{ReviewRunID: "run-trace", PacketDigest: packet.PacketDigest, Plane: PlaneTraceEvidence, RoleID: "trace", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusParseFailed},
			{ReviewRunID: "run-req", PacketDigest: packet.PacketDigest, Plane: PlaneRequirements, RoleID: "req", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusFindingsReported, Findings: []Finding{{ID: "F1", Severity: SeverityMajor, Citation: Citation{ContextRefID: "diff", DiffHunkID: "hunk-1"}, Summary: "SYNTHETIC_TOKEN_SECRET_DO_NOT_LEAK"}}},
		},
	}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageCannotVerify {
		t.Fatalf("coverage = %s want cannot_verify", validation.ReviewCoverageState)
	}
	if validation.MergeDecision != DecisionNotAuthorized || validation.AuthorityScope != AuthorityReviewRecordOnly {
		t.Fatalf("authority boundary missing: %+v", validation)
	}
	payload, err := json.Marshal(validation)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), "SYNTHETIC_TOKEN_SECRET_DO_NOT_LEAK") {
		t.Fatalf("validation leaked marker: %s", payload)
	}
	if !strings.Contains(string(payload), "parse_failed") {
		t.Fatalf("validation missing failed run status: %s", payload)
	}
}

func TestSynthesizeAndValidateCoverageSatisfied(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("3")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StatePass}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{ReviewRunID: "run-code", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-code")}}}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageSatisfied {
		t.Fatalf("coverage = %s want %s reasons=%v", validation.ReviewCoverageState, CoverageSatisfied, validation.Reasons)
	}
	summary := Summarize(validation, ledger)
	for _, forbidden := range []string{"safe to merge", "approved", "ready", "policy passed"} {
		if strings.Contains(strings.ToLower(summary), forbidden) {
			t.Fatalf("summary contains forbidden phrase %q: %s", forbidden, summary)
		}
	}
	if !strings.Contains(summary, DecisionNotAuthorized) {
		t.Fatalf("summary missing authority boundary: %s", summary)
	}
}

func TestValidateCannotVerifyUsableStatusWithoutRetainedOutput(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("8")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StatePass}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{ReviewRunID: "run-code", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings}}}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageCannotVerify {
		t.Fatalf("coverage = %s want cannot_verify", validation.ReviewCoverageState)
	}
	if len(validation.PlaneResults) != 1 || validation.PlaneResults[0].Reason != "reviewer_output_not_retained" {
		t.Fatalf("retained-output reason missing: %+v", validation.PlaneResults)
	}
}

func TestValidateUsesBestPlaneResultAcrossRetries(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("8")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StatePass}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{
		{ReviewRunID: "run-code-first", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusParseFailed},
		{ReviewRunID: "run-code-retry", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-code-retry")},
	}}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageSatisfied {
		t.Fatalf("coverage = %s want %s validation=%+v", validation.ReviewCoverageState, CoverageSatisfied, validation)
	}
	if len(validation.PlaneResults) != 1 || validation.PlaneResults[0].RunID != "run-code-retry" {
		t.Fatalf("best plane result not selected: %+v", validation.PlaneResults)
	}
}

func TestValidateCannotVerifyUnexplainedModelMismatch(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("5")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "model-a"}}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{ReviewRunID: "run-code", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "model-a", ObservedModel: "model-b", ModelFamily: "family", ModelVersion: "v1", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-code")}}}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageCannotVerify {
		t.Fatalf("coverage=%s want cannot_verify validation=%+v", validation.ReviewCoverageState, validation)
	}
	if len(validation.PlaneResults) != 1 || validation.PlaneResults[0].Reason != "model_identity_mismatch" {
		t.Fatalf("model mismatch reason missing: %+v", validation.PlaneResults)
	}

	runs.Results[0].FallbackForModel = "model-a"
	runs.Results[0].FallbackReason = "primary_unavailable"
	validation = Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageSatisfied {
		t.Fatalf("fallback provenance should make mismatch verifiable, got %+v", validation)
	}
}

func TestValidateCannotVerifyPerResultPacketDigestMismatch(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("a")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{ReviewRunID: "run-code", PacketDigest: "sha256:" + sixtyFour("b"), Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings}}}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageCannotVerify {
		t.Fatalf("per-result stale digest should be cannot_verify: %+v", validation)
	}
	if !strings.Contains(strings.Join(validation.Reasons, ","), "result_packet_digest_mismatch") {
		t.Fatalf("per-result mismatch reason missing: %+v", validation.Reasons)
	}
}

func TestValidateCoverageStatesForNoReviewersUnresolvedAndStaleDigest(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("6")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}

	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageNotAssessed {
		t.Fatalf("no reviewers coverage=%s want not_assessed validation=%+v", validation.ReviewCoverageState, validation)
	}

	runs.Results = []ReviewerResult{{
		ReviewRunID:    "run-code",
		PacketDigest:   packetDigest,
		Plane:          PlaneCodeCorrectness,
		RoleID:         "code",
		Runner:         RunnerManualExternal,
		RequestedModel: "not_assessed",
		ObservedModel:  "not_assessed",
		ModelFamily:    "not_assessed",
		ModelVersion:   "not_assessed",
		Status:         StatusFindingsReported,
		RawOutputRef:   retainedRawRef("run-code"),
		Findings:       []Finding{{ID: "F1", Severity: SeverityMajor, Citation: Citation{ContextRefID: "diff", DiffHunkID: "hunk-1"}, Summary: "Missing validation."}},
	}}
	ledger = SynthesizeLedger(packet, runs, nil)
	validation = Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageUnresolved {
		t.Fatalf("unresolved finding coverage=%s want coverage_unresolved validation=%+v ledger=%+v", validation.ReviewCoverageState, validation, ledger)
	}

	runs.PacketDigest = "sha256:" + sixtyFour("7")
	validation = Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageCannotVerify {
		t.Fatalf("stale run digest coverage=%s want cannot_verify validation=%+v", validation.ReviewCoverageState, validation)
	}
}

func TestReadRunSetRejectsDuplicateRunIDs(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "results.json")
	data := `{
  "schema_version": "block30-pr-review-runs-v1",
  "packet_digest": "sha256:` + sixtyFour("1") + `",
  "results": [
    {"review_run_id":"dup","packet_digest":"sha256:` + sixtyFour("1") + `","plane":"code_correctness","role_id":"a","runner":"manual_external","requested_model":"not_assessed","observed_model":"not_assessed","model_family":"not_assessed","model_version":"not_assessed","status":"no_findings","findings":[]},
    {"review_run_id":"dup","packet_digest":"sha256:` + sixtyFour("1") + `","plane":"trace_evidence_provenance","role_id":"b","runner":"manual_external","requested_model":"not_assessed","observed_model":"not_assessed","model_family":"not_assessed","model_version":"not_assessed","status":"no_findings","findings":[]}
  ]
}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ReadRunSet(path)
	if err == nil || !strings.Contains(err.Error(), "duplicate_review_run_id") {
		t.Fatalf("expected duplicate review_run_id rejection, got %v", err)
	}
}

func TestValidationAndSummaryRedactUnsafeMarkerClasses(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("8")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlanePrivacySafety}, Roles: []ReviewRole{{RoleID: "privacy", Plane: PlanePrivacySafety, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}
	markers := []string{
		"SYNTHETIC_PROMPT_SECRET_ALPHA",
		"SYNTHETIC_TOKEN_SECRET_ALPHA",
		"SYNTHETIC_PRIVATE_PATH_ALPHA",
		"SYNTHETIC_AUTH_URL_ALPHA",
		"SYNTHETIC_MODEL_RESPONSE_ALPHA",
		"Bearer should-not-render",
		"https://token@example.invalid/path",
		"/Users/private/project",
	}
	findings := make([]Finding, 0, len(markers))
	for i, marker := range markers {
		findings = append(findings, Finding{ID: fmt.Sprintf("F%d", i+1), Severity: SeverityMinor, Citation: Citation{ContextRefID: "diff", DiffHunkID: "hunk-1"}, Summary: marker})
	}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{ReviewRunID: "run-privacy", PacketDigest: packetDigest, Plane: PlanePrivacySafety, RoleID: "privacy", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusFindingsReported, Findings: findings}}}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	summary := Summarize(validation, ledger)
	if !strings.Contains(summary, "CI state: not_assessed") {
		t.Fatalf("summary should render CI state: %s", summary)
	}
	payload, err := json.Marshal(validation)
	if err != nil {
		t.Fatal(err)
	}
	combined := string(payload) + "\n" + summary
	for _, marker := range markers {
		if strings.Contains(combined, marker) {
			t.Fatalf("unsafe marker leaked: %s\n%s", marker, combined)
		}
	}
	if !strings.Contains(combined, "[redacted unsafe reviewer text]") {
		t.Fatalf("redaction marker missing: %s", combined)
	}
}

func TestRunReviewRecordsRunnerFailureStatesAndPromptDigest(t *testing.T) {
	root := t.TempDir()
	workDir := filepath.Join(root, "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runGitInit(workDir); err != nil {
		t.Fatal(err)
	}
	packetDigest := "sha256:" + sixtyFour("4")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	promptPath := writeText(t, root, "prompt.md", "review {{packet_digest}}\n")
	helper := os.Args[0]
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "runner-states",
		RequiredPlanes: []string{PlaneCodeCorrectness, PlaneTraceEvidence, PlaneRequirements, PlanePrivacySafety, PlaneSecurity, PlaneDXReplayability},
		Roles: []ReviewRole{
			{RoleID: "empty", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "fake-empty", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "empty"}, PromptTemplateRef: promptPath},
			{RoleID: "malformed", Plane: PlaneTraceEvidence, Runner: RunnerManualExternal, RequestedModel: "fake-malformed", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "malformed"}},
			{RoleID: "offtask", Plane: PlaneRequirements, Runner: RunnerManualExternal, RequestedModel: "fake-offtask", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "offtask"}},
			{RoleID: "readonly", Plane: PlanePrivacySafety, Runner: RunnerOpenCode, RequestedModel: "fake-opencode", ReadOnlyEnforced: false, Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "success"}},
			{RoleID: "pi-success", Plane: PlaneSecurity, Runner: RunnerPI, RequestedModel: "fake-pi", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "pi-success"}, PromptTemplateRef: promptPath, RawOutputRetention: RedactionDigestOnly},
			{RoleID: "opencode-mutation", Plane: PlaneDXReplayability, Runner: RunnerOpenCode, RequestedModel: "fake-opencode", ReadOnlyEnforced: true, WorkingTreeMode: "clean_required", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "opencode-mutation"}},
		},
	}
	t.Setenv("GO_WANT_PR_REVIEW_HELPER_PROCESS", "1")
	runs, preview, err := RunReview(packet, profile, RunOptions{
		OutDir:         filepath.Join(root, "runs"),
		AllowedRunners: map[string]bool{RunnerOpenCode: true, RunnerPI: true},
		Now:            time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
		WorkDir:        workDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if preview != nil {
		t.Fatalf("unexpected preview: %+v", preview)
	}
	statuses := map[string]ReviewerResult{}
	for _, result := range runs.Results {
		statuses[result.RoleID] = result
		if result.PacketDigest != packetDigest {
			t.Fatalf("result not bound to packet digest: %+v", result)
		}
	}
	expected := map[string]string{
		"empty":             StatusEmptyOutput,
		"malformed":         StatusParseFailed,
		"offtask":           StatusOffTask,
		"readonly":          StatusNotAssessed,
		"pi-success":        StatusNoFindings,
		"opencode-mutation": StatusCannotVerify,
	}
	for roleID, status := range expected {
		if statuses[roleID].Status != status {
			t.Fatalf("%s status=%s want %s result=%+v", roleID, statuses[roleID].Status, status, statuses[roleID])
		}
	}
	if statuses["empty"].CommandDigest == "" || statuses["empty"].PromptRef == nil || statuses["empty"].PromptRef.RedactionState != RedactionDigestOnly {
		t.Fatalf("runner provenance missing command or prompt digest: %+v", statuses["empty"])
	}
	if statuses["empty"].RawOutputRef == nil || statuses["malformed"].RawOutputRef == nil || statuses["offtask"].RawOutputRef == nil {
		t.Fatalf("raw output digest refs missing: %+v", statuses)
	}
	if statuses["pi-success"].RawOutputRef == nil || !strings.HasPrefix(statuses["pi-success"].RawOutputRef.Ref, "digest-only:") {
		t.Fatalf("pi digest-only raw output ref missing: %+v", statuses["pi-success"])
	}
	if _, err := os.Stat(filepath.Join(root, "runs", "raw", "run-pi-success.out")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("digest-only raw output should not persist raw bytes, err=%v", err)
	}
	if statuses["readonly"].Reason != "opencode_read_only_not_enforced" {
		t.Fatalf("opencode read-only reason missing: %+v", statuses["readonly"])
	}
	if statuses["opencode-mutation"].Reason != "mutation_detected" {
		t.Fatalf("opencode mutation reason missing: %+v", statuses["opencode-mutation"])
	}
}

func TestRunReviewPreviewReturnsPreviewOnly(t *testing.T) {
	root := t.TempDir()
	packet := Packet{PacketDigest: "sha256:" + sixtyFour("v"), SchemaVersion: SchemaVersionPacket}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "preview",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles: []ReviewRole{
			{
				RoleID:         "code",
				Plane:          PlaneCodeCorrectness,
				Runner:         RunnerManualExternal,
				RequestedModel: "not_assessed",
				TimeoutSeconds: 120,
			},
		},
	}
	outDir := filepath.Join(root, "unused")
	runSet, preview, err := RunReview(packet, profile, RunOptions{OutDir: outDir, Preview: true})
	if err != nil {
		t.Fatal(err)
	}
	if preview == nil {
		t.Fatal("expected preview response")
	}
	if runSet.SchemaVersion != "" || len(runSet.Results) != 0 {
		t.Fatalf("preview mode should not produce run-set: %+v", runSet)
	}
	if preview.SchemaVersion != SchemaVersionRunSet {
		t.Fatalf("preview schema = %q", preview.SchemaVersion)
	}
	if preview.PacketDigest != packet.PacketDigest {
		t.Fatalf("preview packet digest = %q want %q", preview.PacketDigest, packet.PacketDigest)
	}
	if len(preview.Roles) != 1 {
		t.Fatalf("preview roles = %d want 1", len(preview.Roles))
	}
	if _, err := os.Stat(outDir); !os.IsNotExist(err) {
		t.Fatalf("preview should not create output directory, got stat err=%v", err)
	}
}

func TestWriteJSONAndReadRunSetUseDirectoryContracts(t *testing.T) {
	root := t.TempDir()
	runSet := RunSet{
		SchemaVersion: SchemaVersionRunSet,
		PacketDigest:  "sha256:" + sixtyFour("7"),
		Results: []ReviewerResult{{
			ReviewRunID:    "run-1",
			PacketDigest:   "sha256:" + sixtyFour("7"),
			Plane:          PlaneCodeCorrectness,
			RoleID:         "code",
			Runner:         RunnerManualExternal,
			RequestedModel: "manual",
			ObservedModel:  "manual",
			ModelFamily:    "manual",
			ModelVersion:   "v1",
			Status:         StatusNoFindings,
		}},
	}
	if err := WriteJSON(" ", runSet); err != nil {
		t.Fatalf("blank WriteJSON path should be ignored: %v", err)
	}
	outDir := filepath.Join(root, "runs")
	if err := WriteJSON(filepath.Join(outDir, "results.json"), runSet); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}

	read, err := ReadRunSet(outDir)
	if err != nil {
		t.Fatalf("ReadRunSet(dir) error = %v", err)
	}
	if len(read.Results) != 1 || read.Results[0].ReviewRunID != "run-1" {
		t.Fatalf("read run set = %+v", read)
	}

	runSet.Results = append(runSet.Results, runSet.Results[0])
	if err := WriteJSON(filepath.Join(outDir, "results.json"), runSet); err != nil {
		t.Fatalf("WriteJSON malformed runset: %v", err)
	}
	if _, err := ReadRunSet(filepath.Join(outDir, "results.json")); err == nil {
		t.Fatalf("expected run-set validation error")
	}
}

func TestPacketProfileAndSmallHelpers(t *testing.T) {
	root := t.TempDir()
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1"}
	packetDir := filepath.Join(root, "packet")
	if err := WriteJSON(filepath.Join(packetDir, "packet.json"), packet); err != nil {
		t.Fatalf("write packet: %v", err)
	}
	readPacket, err := ReadPacket(packetDir)
	if err != nil {
		t.Fatalf("ReadPacket(dir) error = %v", err)
	}
	if readPacket.PacketID != "packet-1" {
		t.Fatalf("packet = %+v", readPacket)
	}

	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "profile",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles:          []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "manual"}},
	}
	profilePath := filepath.Join(root, "profile.json")
	if err := WriteJSON(profilePath, profile); err != nil {
		t.Fatalf("write profile: %v", err)
	}
	if _, err := ReadProfile(profilePath); err != nil {
		t.Fatalf("ReadProfile() error = %v", err)
	}

	if got := defaultReviewerStatus([]Finding{{ID: "f1"}}); got != StatusFindingsReported {
		t.Fatalf("default status with finding = %s", got)
	}
	if got := defaultReviewerStatus(nil); got != StatusNoFindings {
		t.Fatalf("default status empty = %s", got)
	}
	if contextKind("task-review.md") != RefKindTask || contextKind("notes.md") != RefKindDoc || contextKind("schema.json") != RefKindSchema || contextKind("diff.patch") != RefKindSourceExcerpt {
		t.Fatalf("contextKind mapping changed")
	}

	metadataPath := writeText(t, root, "metadata.json", `{"ok":true}`)
	inputsDir := filepath.Join(root, "inputs")
	if err := os.MkdirAll(inputsDir, 0o755); err != nil {
		t.Fatalf("mkdir inputs: %v", err)
	}
	ref, err := optionalMetadataRef(inputsDir, metadataPath)
	if err != nil {
		t.Fatalf("optionalMetadataRef() error = %v", err)
	}
	if ref == nil || ref.Kind != RefKindMetadata {
		t.Fatalf("metadata ref = %+v", ref)
	}
	empty, err := optionalMetadataRef(filepath.Join(root, "inputs-empty"), "")
	if err != nil || empty != nil {
		t.Fatalf("empty metadata ref = %+v err=%v", empty, err)
	}
}

func TestApplyRunnerErrorClassifiesUnavailableAndFailure(t *testing.T) {
	result := ReviewerResult{}
	if err := applyRunnerError(&result, exec.ErrNotFound); err == nil {
		t.Fatalf("expected error returned")
	}
	if result.Status != StatusNotAssessed || result.Reason != "runner_unavailable" {
		t.Fatalf("unavailable result = %+v", result)
	}
	result = ReviewerResult{}
	if err := applyRunnerError(&result, fmt.Errorf("boom")); err == nil {
		t.Fatalf("expected error returned")
	}
	if result.Status != StatusFailed || result.Reason != "runner_failed" {
		t.Fatalf("failed result = %+v", result)
	}
}

func TestRunReviewNotAssessedReasonDoesNotInvokeRunner(t *testing.T) {
	root := t.TempDir()
	packetDigest := "sha256:" + sixtyFour("7")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "missing-secrets",
		RequiredPlanes: []string{PlaneCodeCorrectness, PlaneTraceEvidence},
		Roles: []ReviewRole{
			{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerPI, RequestedModel: "minimax/MiniMax-M2.7", Command: []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "should-not-run"}},
			{RoleID: "trace", Plane: PlaneTraceEvidence, Runner: RunnerPI, RequestedModel: "zai/glm-5.1", Command: []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "should-not-run"}},
		},
	}
	runs, _, err := RunReview(packet, profile, RunOptions{OutDir: filepath.Join(root, "runs"), NotAssessedReason: "ci_model_review_not_configured"})
	if err != nil {
		t.Fatal(err)
	}
	if len(runs.Results) != 2 {
		t.Fatalf("unexpected result count: %+v", runs.Results)
	}
	for _, result := range runs.Results {
		if result.Status != StatusNotAssessed || result.Reason != "ci_model_review_not_configured" {
			t.Fatalf("missing-secret result should be not_assessed: %+v", result)
		}
	}
}

func TestRunReviewCannotVerifyUnreadablePromptTemplate(t *testing.T) {
	root := t.TempDir()
	packetDigest := "sha256:" + sixtyFour("8")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "missing-prompt",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles: []ReviewRole{{
			RoleID:            "code",
			Plane:             PlaneCodeCorrectness,
			Runner:            RunnerManualExternal,
			RequestedModel:    "fake",
			Command:           []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "should-not-run"},
			PromptTemplateRef: filepath.Join(root, "missing.md"),
		}},
	}
	runs, _, err := RunReview(packet, profile, RunOptions{OutDir: filepath.Join(root, "runs")})
	if err != nil {
		t.Fatal(err)
	}
	if len(runs.Results) != 1 || runs.Results[0].Status != StatusCannotVerify || runs.Results[0].Reason != "prompt_ref_cannot_verify" {
		t.Fatalf("missing prompt should be cannot_verify without runner execution: %+v", runs.Results)
	}
}

func TestRunReviewMapsTimeoutToTimedOut(t *testing.T) {
	root := t.TempDir()
	packetDigest := "sha256:" + sixtyFour("9")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "timeout",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles: []ReviewRole{{
			RoleID:         "timeout",
			Plane:          PlaneCodeCorrectness,
			Runner:         RunnerManualExternal,
			RequestedModel: "fake-timeout",
			Command:        []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "timeout"},
			TimeoutSeconds: 1,
		}},
	}
	t.Setenv("GO_WANT_PR_REVIEW_HELPER_PROCESS", "1")
	runs, _, err := RunReview(packet, profile, RunOptions{OutDir: filepath.Join(root, "runs")})
	if err != nil {
		t.Fatal(err)
	}
	if len(runs.Results) != 1 || runs.Results[0].Status != StatusTimedOut || runs.Results[0].Reason != "runner_timed_out" {
		t.Fatalf("timeout mapping failed: %+v", runs.Results)
	}
}

func TestSafeID(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{name: "normalizes-case-and-space-trims", in: "  Review_Run-01 ", out: "review_run-01"},
		{name: "replaces-invalid-characters", in: "a@b c", out: "a-b-c"},
		{name: "retains-safe-punctuation", in: "a-b.c_1", out: "a-b.c_1"},
		{name: "trims-unsafe-boundaries", in: "---item.", out: "item"},
		{name: "unicode-to-dash-and-default", in: "π_Т-9", out: "_--9"},
		{name: "all-invalid-becomes-item", in: " !!! ... \n", out: "item"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := safeID(tc.in); got != tc.out {
				t.Fatalf("safeID(%q) = %q, want %q", tc.in, got, tc.out)
			}
		})
	}
}

func TestPRReviewFakeRunnerHelper(t *testing.T) {
	if os.Getenv("GO_WANT_PR_REVIEW_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	mode := args[len(args)-1]
	switch mode {
	case "empty":
		os.Exit(0)
	case "malformed":
		fmt.Print("{")
		os.Exit(0)
	case "offtask":
		fmt.Print(`{"packet_digest":"sha256:0000000000000000000000000000000000000000000000000000000000000000","plane":"requirements_vs_implementation","role_id":"offtask","runner":"manual_external","requested_model":"fake","observed_model":"fake","model_family":"fake","model_version":"fake","status":"no_findings","findings":[]}`)
		os.Exit(0)
	case "success":
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("4") + `","plane":"privacy_output_safety","role_id":"readonly","runner":"opencode","requested_model":"fake","observed_model":"fake","model_family":"fake","model_version":"fake","status":"no_findings","findings":[]}`)
		os.Exit(0)
	case "pi-success":
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil || !strings.Contains(string(stdin), "sha256:"+sixtyFour("4")) {
			os.Exit(3)
		}
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("4") + `","plane":"security_forgery_overclaim","role_id":"pi-success","runner":"pi","requested_model":"fake-pi","observed_model":"fake-pi","model_family":"fake","model_version":"v1","status":"no_findings","findings":[]}`)
		os.Exit(0)
	case "should-not-run":
		os.Exit(4)
	case "opencode-mutation":
		if err := os.WriteFile("mutated-by-helper.txt", []byte("mutation\n"), 0o644); err != nil {
			os.Exit(2)
		}
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("4") + `","plane":"dx_replayability","role_id":"opencode-mutation","runner":"opencode","requested_model":"fake-opencode","observed_model":"fake-opencode","model_family":"fake","model_version":"v1","status":"no_findings","findings":[]}`)
		os.Exit(0)
	case "timeout":
		time.Sleep(2 * time.Second)
		os.Exit(0)
	default:
		os.Exit(2)
	}
}

func runGitInit(workDir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = workDir
	return cmd.Run()
}

func writeText(t *testing.T, root, name, content string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func forty(s string) string {
	return strings.Repeat(s, 40)
}

func sixtyFour(s string) string {
	return strings.Repeat(s, 64)
}

func retainedRawRef(id string) *SafeRef {
	return &SafeRef{
		ID:             "raw-" + id,
		Kind:           RefKindRawOutput,
		Ref:            "raw/" + id + ".out",
		DigestSHA256:   sixtyFour("9"),
		ContentType:    ContentText,
		RedactionState: RedactionDigestOnly,
	}
}
