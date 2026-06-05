package prreview

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestPrreviewSchemaConstantsAndPatternsPreserveContracts(t *testing.T) {
	for name, tt := range map[string]struct{ got, want string }{
		"packet schema":                        {SchemaVersionPacket, "block30-pr-review-packet-v1"},
		"profile schema":                       {SchemaVersionProfile, "block30-pr-review-profile-v1"},
		"runset schema":                        {SchemaVersionRunSet, "block30-pr-review-runs-v1"},
		"ledger schema":                        {SchemaVersionLedger, "block30-pr-review-ledger-v1"},
		"validation schema":                    {SchemaVersionValidation, "block30-pr-review-validation-v1"},
		"pass state":                           {StatePass, "pass"},
		"fail state":                           {StateFail, "fail"},
		"pending state":                        {StatePending, "pending"},
		"not assessed state":                   {StateNotAssessed, "not_assessed"},
		"cannot verify state":                  {StateCannotVerify, "cannot_verify"},
		"diff ref kind":                        {RefKindDiff, "diff"},
		"metadata ref kind":                    {RefKindMetadata, "metadata"},
		"spec ref kind":                        {RefKindSpec, "spec"},
		"plan ref kind":                        {RefKindPlan, "plan"},
		"task ref kind":                        {RefKindTask, "task"},
		"doc ref kind":                         {RefKindDoc, "doc"},
		"schema ref kind":                      {RefKindSchema, "schema"},
		"source excerpt ref kind":              {RefKindSourceExcerpt, "source_excerpt"},
		"verification ref kind":                {RefKindVerification, "verification"},
		"prompt ref kind":                      {RefKindPrompt, "prompt"},
		"raw output ref kind":                  {RefKindRawOutput, "raw_output"},
		"sanitized output ref kind":            {RefKindSanitizedOutput, "sanitized_output"},
		"external ref kind":                    {RefKindExternal, "external"},
		"unified diff content":                 {ContentUnifiedDiff, "unified_diff"},
		"markdown content":                     {ContentMarkdown, "markdown"},
		"json content":                         {ContentJSON, "json"},
		"text content":                         {ContentText, "text"},
		"none redaction":                       {RedactionNone, "none"},
		"redacted redaction":                   {RedactionRedacted, "redacted"},
		"digest-only redaction":                {RedactionDigestOnly, "digest_only"},
		"encrypted redaction":                  {RedactionEncrypted, "encrypted_ref"},
		"external ref redaction":               {RedactionExternalRef, "external_ref"},
		"withheld redaction":                   {RedactionWithheld, "withheld"},
		"not assessed redaction":               {RedactionNotAssessed, "not_assessed"},
		"code plane":                           {PlaneCodeCorrectness, "code_correctness"},
		"trace plane":                          {PlaneTraceEvidence, "trace_evidence_provenance"},
		"requirements plane":                   {PlaneRequirements, "requirements_vs_implementation"},
		"security plane":                       {PlaneSecurity, "security_forgery_overclaim"},
		"dx plane":                             {PlaneDXReplayability, "dx_replayability"},
		"privacy plane":                        {PlanePrivacySafety, "privacy_output_safety"},
		"pi runner":                            {RunnerPI, "pi"},
		"opencode runner":                      {RunnerOpenCode, "opencode"},
		"manual runner":                        {RunnerManualExternal, "manual_external"},
		"findings reported status":             {StatusFindingsReported, "findings_reported"},
		"no findings status":                   {StatusNoFindings, "no_findings"},
		"status not assessed":                  {StatusNotAssessed, "not_assessed"},
		"failed status":                        {StatusFailed, "failed"},
		"timed out status":                     {StatusTimedOut, "timed_out"},
		"empty output status":                  {StatusEmptyOutput, "empty_output"},
		"off task status":                      {StatusOffTask, "off_task"},
		"parse failed status":                  {StatusParseFailed, "parse_failed"},
		"status cannot verify":                 {StatusCannotVerify, "cannot_verify"},
		"critical severity":                    {SeverityCritical, "critical"},
		"major severity":                       {SeverityMajor, "major"},
		"minor severity":                       {SeverityMinor, "minor"},
		"informational severity":               {SeverityInformational, "informational"},
		"accepted fixed disposition":           {DispositionAcceptedFixed, "accepted_fixed"},
		"accepted review blocking disposition": {DispositionAcceptedReviewBlocking, "accepted_review_blocking"},
		"accepted narrower disposition":        {DispositionAcceptedNarrower, "accepted_narrower"},
		"rejected false positive disposition":  {DispositionRejectedFalsePositive, "rejected_false_positive"},
		"deferred not assessed disposition":    {DispositionDeferredNotAssessed, "deferred_not_assessed"},
		"unresolved disposition":               {DispositionUnresolvedReviewBlocker, "unresolved_review_blocker"},
		"coverage satisfied":                   {CoverageSatisfied, "coverage_satisfied"},
		"coverage partial":                     {CoveragePartial, "coverage_partial"},
		"coverage unresolved":                  {CoverageUnresolved, "coverage_unresolved"},
		"coverage not assessed":                {CoverageNotAssessed, "not_assessed"},
		"coverage cannot verify":               {CoverageCannotVerify, "cannot_verify"},
		"authority review record only":         {AuthorityReviewRecordOnly, "review_record_only"},
		"decision not authorized":              {DecisionNotAuthorized, "not_authorized_by_sdp_trace"},
	} {
		if tt.got != tt.want {
			t.Fatalf("%s = %q, want %q", name, tt.got, tt.want)
		}
	}
	if errPromptEvidenceCannotVerify.Error() != "prompt_evidence_cannot_verify" || errPromptTemplateCannotVerify.Error() != "prompt_template_cannot_verify" {
		t.Fatalf("prompt errors drifted: %v / %v", errPromptEvidenceCannotVerify, errPromptTemplateCannotVerify)
	}
	if repoIDPattern.String() != `^[a-z0-9][a-z0-9._-]{0,62}[a-z0-9]$` || changeRefPattern.String() != `^(pr|mr|change)-[A-Za-z0-9._-]{1,64}$` || sha40Pattern.String() != `^[0-9a-f]{40}$` {
		t.Fatalf("regex pattern strings drifted: %q / %q / %q", repoIDPattern.String(), changeRefPattern.String(), sha40Pattern.String())
	}
	for _, id := range []string{"demo_repo", "demo.repo-1"} {
		if !repoIDPattern.MatchString(id) {
			t.Fatalf("repoIDPattern rejected %q", id)
		}
	}
	for _, id := range []string{"Demo", "-demo", strings.Repeat("a", 65)} {
		if repoIDPattern.MatchString(id) {
			t.Fatalf("repoIDPattern accepted %q", id)
		}
	}
	if !changeRefPattern.MatchString("pr-123") || !changeRefPattern.MatchString("mr-feature.1") || changeRefPattern.MatchString("issue-123") {
		t.Fatalf("changeRefPattern contract drifted")
	}
	if !sha40Pattern.MatchString(forty("a")) || sha40Pattern.MatchString(forty("g")) || sha40Pattern.MatchString(sixtyFour("a")) {
		t.Fatalf("sha40Pattern contract drifted")
	}
}

func TestPrreviewPortableTypesPreserveJSONShape(t *testing.T) {
	packetKeys := jsonKeys(t, Packet{
		SchemaVersion:    SchemaVersionPacket,
		PacketID:         "packet-1",
		PacketDigest:     "sha256:" + sixtyFour("1"),
		RepoID:           "demo_repo",
		ChangeRef:        "pr-123",
		BaseCommit:       forty("a"),
		HeadCommit:       forty("b"),
		DiffRef:          SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone},
		ContextRefs:      []SafeRef{},
		VerificationRefs: []SafeRef{},
		CIState:          StateNotAssessed,
		CreatedAt:        "2026-05-09T12:00:00Z",
		CreatedBy:        "test",
		RedactionState:   RedactionNone,
	})
	assertJSONKeys(t, packetKeys, []string{"schema_version", "packet_id", "packet_digest", "repo_id", "change_ref", "base_commit", "head_commit", "diff_ref", "context_refs", "verification_refs", "ci_state", "created_at", "created_by", "redaction_state"}, []string{"metadata_ref", "unavailable_fields"})

	roleKeys := jsonKeys(t, ReviewRole{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "manual"})
	assertJSONKeys(t, roleKeys, []string{"role_id", "plane", "runner", "requested_model"}, []string{"command", "timeout_seconds", "prompt_template_ref", "required_output_schema", "raw_output_retention", "read_only_enforced", "working_tree_mode"})

	resultKeys := jsonKeys(t, ReviewerResult{ReviewRunID: "run-1", PacketDigest: "sha256:" + sixtyFour("1"), Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "manual", ObservedModel: "manual", ModelFamily: "manual", ModelVersion: "v1", Status: StatusNoFindings, Findings: []Finding{}})
	assertJSONKeys(t, resultKeys, []string{"review_run_id", "packet_digest", "plane", "role_id", "runner", "requested_model", "observed_model", "model_family", "model_version", "status", "findings"}, []string{"fallback_for_model", "fallback_reason", "command_digest", "raw_output_ref", "prompt_ref", "context_refs", "started_at", "ended_at", "reason"})

	validationKeys := jsonKeys(t, Validation{SchemaVersion: SchemaVersionValidation, PacketDigest: "sha256:" + sixtyFour("1"), ReviewCoverageState: CoverageNotAssessed, CIState: StateNotAssessed, AuthorityScope: AuthorityReviewRecordOnly, MergeDecision: DecisionNotAuthorized, ReleaseDecision: DecisionNotAuthorized, RiskAcceptance: DecisionNotAuthorized, PlaneResults: []PlaneResult{}, Findings: []LedgerFinding{}, Reasons: []string{}, NextActions: []string{}})
	assertJSONKeys(t, validationKeys, []string{"schema_version", "packet_digest", "review_coverage_state", "ci_state", "authority_scope", "merge_decision", "release_decision", "risk_acceptance", "plane_results", "findings", "reasons", "next_actions"}, nil)

	planeKeys := jsonKeys(t, PlaneResult{Plane: PlaneCodeCorrectness, Status: StateNotAssessed})
	assertJSONKeys(t, planeKeys, []string{"plane", "status", "usable"}, []string{"review_run_id", "reason", "next_action"})

	assertJSONKeys(t, jsonKeys(t, SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}), []string{"id", "kind", "ref", "digest_sha256", "content_type", "redaction_state"}, nil)
	assertJSONKeys(t, jsonKeys(t, UnavailableField{Field: "metadata_ref", State: StateNotAssessed, Reason: "missing"}), []string{"field", "state", "reason"}, nil)
	assertJSONKeys(t, jsonKeys(t, ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "default", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{}}), []string{"schema_version", "profile_id", "required_planes", "roles"}, nil)
	assertJSONKeys(t, jsonKeys(t, RunPreview{SchemaVersion: SchemaVersionRunSet, PacketDigest: "sha256:" + sixtyFour("1"), Roles: []PreviewRole{}}), []string{"schema_version", "packet_digest", "roles"}, nil)
	assertJSONKeys(t, jsonKeys(t, PreviewRole{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "manual"}), []string{"role_id", "plane", "runner", "requested_model", "timeout_seconds", "command_digest"}, []string{"prompt_template_ref", "prompt_digest"})
	assertJSONKeys(t, jsonKeys(t, RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: "sha256:" + sixtyFour("1"), Results: []ReviewerResult{}}), []string{"schema_version", "packet_digest", "results"}, nil)
	assertJSONKeys(t, jsonKeys(t, Finding{ID: "F1", Severity: SeverityMajor, Citation: Citation{ContextRefID: "diff"}, Summary: "summary"}), []string{"id", "severity", "citation", "summary"}, []string{"rationale", "suggested_fix", "question", "evidence_refs"})
	assertJSONKeys(t, jsonKeys(t, Citation{}), nil, []string{"context_ref_id", "diff_hunk_id", "source_digest", "line_start", "line_end"})
	assertJSONKeys(t, jsonKeys(t, Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: "sha256:" + sixtyFour("1"), Findings: []LedgerFinding{}}), []string{"schema_version", "packet_digest", "findings"}, nil)
	assertJSONKeys(t, jsonKeys(t, LedgerFinding{ID: "F1", ReviewRunID: "run-1", Plane: PlaneCodeCorrectness, RoleID: "code", Severity: SeverityMajor, Summary: "summary", Citation: Citation{ContextRefID: "diff"}, Disposition: DispositionUnresolvedReviewBlocker}), []string{"id", "review_run_id", "plane", "role_id", "severity", "summary", "citation", "disposition"}, []string{"evidence_refs", "disposition_evidence"})
}

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

	unsafeOutDir := filepath.Join(root, "unsafe")
	_, err = BuildPacket(PacketOptions{
		OutDir:     unsafeOutDir,
		RepoID:     "/tmp/private/repo",
		ChangeRef:  "pr-123",
		BaseCommit: forty("a"),
		HeadCommit: forty("b"),
		DiffPath:   diffPath,
	})
	if err == nil || !strings.Contains(err.Error(), "unsafe_repo_id") {
		t.Fatalf("expected unsafe repo id rejection, got %v", err)
	}
	if _, err := os.Stat(unsafeOutDir); !os.IsNotExist(err) {
		t.Fatalf("unsafe output directory should not be created before identity validation, got %v", err)
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
		"context-ref-with-hunk": {
			citation: Citation{ContextRefID: "spec", DiffHunkID: "hunk-1"},
			want:     true,
		},
		"context-ref-with-digest": {
			citation: Citation{ContextRefID: "spec", SourceDigest: "sha256:abc"},
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
		"verification-ref-with-digest": {
			citation: Citation{ContextRefID: "verify", SourceDigest: "sha256:abc"},
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
		"unknown-ref-without-digest": {
			citation: Citation{ContextRefID: "unknown"},
			want:     false,
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

func TestBuildPacketCopiesInputsAndComputesStableDigests(t *testing.T) {
	root := t.TempDir()
	diffData := "diff --git a/a.go b/a.go\n+package main\n"
	contextData := "# Spec\n"
	verificationData := "go test ./...\n"
	diffPath := writeText(t, root, "change.diff", diffData)
	contextPath := writeText(t, root, "spec.md", contextData)
	verificationPath := writeText(t, root, "verify.log", verificationData)
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
	assertCopiedInput(t, outDir, packet.ContextRefs[0], "context-1.md", "context-1", RefKindDoc, ContentMarkdown, contextData)
	assertCopiedInput(t, outDir, packet.VerificationRefs[0], "verification-1.txt", "verification-1", RefKindVerification, ContentText, verificationData)

	canonical := packet
	canonical.PacketDigest = "sha256:" + sixtyFour("f")
	replayed, err := packetDigest(canonical)
	if err != nil {
		t.Fatal(err)
	}
	if "sha256:"+replayed != packet.PacketDigest {
		t.Fatalf("packet digest should be replayable with packet_digest cleared: got sha256:%s want %s", replayed, packet.PacketDigest)
	}
}

func assertCopiedInput(t *testing.T, outDir string, ref SafeRef, name, id, kind, contentType, wantData string) {
	t.Helper()
	if ref.ID != id || ref.Kind != kind || ref.Ref != filepath.ToSlash(filepath.Join("inputs", name)) ||
		ref.ContentType != contentType || ref.RedactionState != RedactionNone {
		t.Fatalf("copied input ref drifted: %+v", ref)
	}
	wantDigest := sha256.Sum256([]byte(wantData))
	if ref.DigestSHA256 != hex.EncodeToString(wantDigest[:]) {
		t.Fatalf("copied input digest = %s want %s", ref.DigestSHA256, hex.EncodeToString(wantDigest[:]))
	}
	copiedPath := filepath.Join(outDir, "inputs", name)
	data, err := os.ReadFile(copiedPath)
	if err != nil {
		t.Fatalf("read copied input: %v", err)
	}
	if string(data) != wantData {
		t.Fatalf("copied input data = %q want %q", string(data), wantData)
	}
	info, err := os.Stat(copiedPath)
	if err != nil {
		t.Fatalf("stat copied input: %v", err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Fatalf("copied input mode = %v want 0644", info.Mode().Perm())
	}
}

func TestPrreviewPacketBuildHelpersPreserveDefaultsRefsAndUnavailableFields(t *testing.T) {
	root := t.TempDir()
	diffPath := writeText(t, root, "change.diff", "diff --git a/a.go b/a.go\n+package main\n")
	metadataPath := writeText(t, root, "metadata.json", `{"pr":123}`)
	specPath := writeText(t, root, "spec.md", "# Spec\n")
	planPath := writeText(t, root, "plan.txt", "plain plan\n")
	verifyPath := writeText(t, root, "verify.log", "go test ./...\n")
	outDir := filepath.Join(root, "packet")

	packet, err := BuildPacket(PacketOptions{
		OutDir:            outDir,
		RepoID:            "demo_repo",
		ChangeRef:         "pr-123",
		BaseCommit:        forty("a"),
		HeadCommit:        forty("b"),
		DiffPath:          diffPath,
		MetadataPath:      metadataPath,
		ContextPaths:      []string{specPath, planPath},
		VerificationPaths: []string{verifyPath},
	})
	if err != nil {
		t.Fatal(err)
	}
	if packet.PacketID != "demo_repo-pr-123-"+forty("b")[:12] || packet.CreatedBy != "sdp-trace" || packet.CIState != StateNotAssessed {
		t.Fatalf("packet defaults drifted: %+v", packet)
	}
	if packet.DiffRef.Kind != RefKindDiff || packet.DiffRef.ContentType != ContentUnifiedDiff {
		t.Fatalf("diff ref drifted: %+v", packet.DiffRef)
	}
	if packet.MetadataRef == nil || packet.MetadataRef.Kind != RefKindMetadata || packet.MetadataRef.ContentType != ContentJSON {
		t.Fatalf("metadata ref drifted: %+v", packet.MetadataRef)
	}
	if len(packet.ContextRefs) != 2 || packet.ContextRefs[0].Kind != RefKindDoc || packet.ContextRefs[1].Kind != RefKindSourceExcerpt {
		t.Fatalf("context refs drifted: %+v", packet.ContextRefs)
	}
	if len(packet.VerificationRefs) != 1 || packet.VerificationRefs[0].Kind != RefKindVerification {
		t.Fatalf("verification refs drifted: %+v", packet.VerificationRefs)
	}
	if len(packet.UnavailableFields) != 0 {
		t.Fatalf("available optional inputs should not be unavailable: %+v", packet.UnavailableFields)
	}

	sentinel := filepath.Join(outDir, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = BuildPacket(PacketOptions{
		OutDir:     outDir,
		RepoID:     "demo_repo",
		ChangeRef:  "pr-123",
		BaseCommit: forty("a"),
		HeadCommit: forty("b"),
		DiffPath:   diffPath,
	})
	if err == nil {
		t.Fatalf("BuildPacket accepted existing output directory")
	}
	if got, readErr := os.ReadFile(sentinel); readErr != nil || string(got) != "keep\n" {
		t.Fatalf("existing output content overwritten: data=%q err=%v", got, readErr)
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

func TestValidationAndLedgerLifecyclePreserveTrustSemantics(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("4")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{SchemaVersion: SchemaVersionProfile, ProfileID: "ledger-lifecycle", RequiredPlanes: []string{PlaneCodeCorrectness}, Roles: []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"}}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{
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
		Findings:       []Finding{{ID: "F-lifecycle", Severity: SeverityMajor, Citation: Citation{ContextRefID: "diff", DiffHunkID: "hunk-1"}, Summary: "lifecycle issue", EvidenceRefs: []string{"review-output"}}},
	}}}

	ledger := SynthesizeLedger(packet, runs, nil)
	if ledger.SchemaVersion != SchemaVersionLedger || ledger.PacketDigest != packetDigest || len(ledger.Findings) != 1 {
		t.Fatalf("ledger lifecycle shape drifted: %+v", ledger)
	}
	if ledger.Findings[0].ReviewRunID != "run-code" || ledger.Findings[0].EvidenceRefs[0] != "review-output" {
		t.Fatalf("ledger should preserve reviewer evidence binding: %+v", ledger.Findings[0])
	}
	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageUnresolved || validation.MergeDecision != DecisionNotAuthorized || validation.AuthorityScope != AuthorityReviewRecordOnly {
		t.Fatalf("validation trust boundary drifted: %+v", validation)
	}
}

func TestLedgerDispositionCarryForward(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("5")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketDigest: packetDigest}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{{
		ReviewRunID: "run-code",
		Plane:       PlaneCodeCorrectness,
		RoleID:      "code",
		Findings: []Finding{
			{ID: "F-carried", Severity: SeverityCritical, Summary: "kept"},
			{ID: "F-default", Severity: SeverityMajor, Summary: "new blocker"},
		},
	}}}
	existing := &Ledger{Findings: []LedgerFinding{{ID: "F-carried", Disposition: DispositionAcceptedFixed}}}

	ledger := SynthesizeLedger(packet, runs, existing)
	byID := map[string]LedgerFinding{}
	for _, finding := range ledger.Findings {
		byID[finding.ID] = finding
	}
	if byID["F-carried"].Disposition != DispositionAcceptedFixed {
		t.Fatalf("prior disposition was not carried forward: %+v", byID["F-carried"])
	}
	if byID["F-default"].Disposition != DispositionUnresolvedReviewBlocker {
		t.Fatalf("new major finding should default to unresolved blocker: %+v", byID["F-default"])
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

func TestPrreviewStatusDispositionHelpersPreserveContracts(t *testing.T) {
	for _, state := range []string{StatePass, StateFail, StatePending, StateNotAssessed, StateCannotVerify} {
		if !validCIState(state) {
			t.Fatalf("validCIState(%q) = false", state)
		}
	}
	for _, state := range []string{"", "green", StatusNoFindings, StatusFailed} {
		if validCIState(state) {
			t.Fatalf("validCIState(%q) = true", state)
		}
	}
	for _, runner := range []string{RunnerPI, RunnerOpenCode, RunnerManualExternal} {
		if !validRunner(runner) {
			t.Fatalf("validRunner(%q) = false", runner)
		}
	}
	for _, runner := range []string{"", "codex", "manual"} {
		if validRunner(runner) {
			t.Fatalf("validRunner(%q) = true", runner)
		}
	}

	for _, status := range []string{StatusFindingsReported, StatusNoFindings} {
		if !reviewerStatusUsable(status) {
			t.Fatalf("reviewerStatusUsable(%q) = false", status)
		}
		withoutRaw := planeResult(ReviewerResult{ReviewRunID: "run-" + status, Plane: PlaneCodeCorrectness, Status: status})
		if withoutRaw.Status != StatusCannotVerify || withoutRaw.Reason != "reviewer_output_not_retained" || withoutRaw.Usable {
			t.Fatalf("usable status without raw output should degrade: %+v", withoutRaw)
		}
		withRaw := planeResult(ReviewerResult{ReviewRunID: "run-" + status, Plane: PlaneCodeCorrectness, Status: status, RawOutputRef: retainedRawRef("run-" + status)})
		if withRaw.Status != status || !withRaw.Usable || withRaw.Reason != "" || withRaw.NextAction != "" {
			t.Fatalf("usable retained result drifted: %+v", withRaw)
		}
	}
	for _, status := range []string{StatusNotAssessed, StatusTimedOut, StatusEmptyOutput, StatusOffTask, StatusParseFailed, StatusCannotVerify} {
		if reviewerStatusUsable(status) {
			t.Fatalf("reviewerStatusUsable(%q) = true", status)
		}
	}

	for status, want := range map[string][2]string{
		StatusNotAssessed: {"reviewer_not_assessed", "Run a configured reviewer or import a usable result for this plane."},
		StatusTimedOut:    {"reviewer_timed_out", "Increase timeout or replace the reviewer for this plane."},
		StatusEmptyOutput: {"reviewer_empty_output", "Retry with a shorter bounded prompt or replace the reviewer."},
		StatusOffTask:     {"reviewer_off_task", "Rerun with the frozen packet and required output schema."},
		StatusParseFailed: {"reviewer_parse_failed", "Rerun with JSON-only output matching the required schema."},
		StatusNoFindings:  {"reviewer_output_not_retained", "Attach digest-bound reviewer output before counting this plane."},
		"unknown":         {"reviewer_cannot_verify", "Replace or rerun the reviewer."},
	} {
		reason, action := reviewerStatusAction(status)
		if reason != want[0] || action != want[1] {
			t.Fatalf("reviewerStatusAction(%q) = %q/%q want %q/%q", status, reason, action, want[0], want[1])
		}
	}

	for _, severity := range []string{SeverityCritical, SeverityMajor} {
		if defaultDisposition(severity) != DispositionUnresolvedReviewBlocker {
			t.Fatalf("%s should default to unresolved blocker", severity)
		}
	}
	for _, severity := range []string{SeverityMinor, SeverityInformational, "unknown", ""} {
		if defaultDisposition(severity) != DispositionDeferredNotAssessed {
			t.Fatalf("%s should default to deferred_not_assessed", severity)
		}
	}
	if safeSeverity("unknown") != SeverityInformational {
		t.Fatalf("unknown severity should fall back to informational")
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

func TestPrreviewValidationOrchestrationPreservesDigestRequiredPlaneAndAuthorityContracts(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("c")
	staleDigest := "sha256:" + sixtyFour("d")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, CIState: StatePending}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "validation-orchestration",
		RequiredPlanes: []string{PlaneTraceEvidence, PlaneCodeCorrectness, PlaneSecurity, PlanePrivacySafety, PlaneCodeCorrectness, ""},
		Roles: []ReviewRole{
			{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
			{RoleID: "trace", Plane: PlaneTraceEvidence, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
			{RoleID: "privacy", Plane: PlanePrivacySafety, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
			{RoleID: "security", Plane: PlaneSecurity, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
		},
	}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: staleDigest, Results: []ReviewerResult{
		{ReviewRunID: "run-code-failed", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusParseFailed},
		{ReviewRunID: "run-code-good", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-code-good")},
		{ReviewRunID: "run-privacy-failed", PacketDigest: packetDigest, Plane: PlanePrivacySafety, RoleID: "privacy", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusParseFailed},
		{ReviewRunID: "run-trace-stale", PacketDigest: staleDigest, Plane: PlaneTraceEvidence, RoleID: "trace", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-trace-stale")},
	}}
	ledger := Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: staleDigest}

	validation := Validate(packet, profile, runs, ledger)
	if validation.SchemaVersion != SchemaVersionValidation || validation.PacketDigest != packetDigest || validation.CIState != StatePending {
		t.Fatalf("validation identity/state drifted: %+v", validation)
	}
	if validation.AuthorityScope != AuthorityReviewRecordOnly || validation.MergeDecision != DecisionNotAuthorized || validation.ReleaseDecision != DecisionNotAuthorized || validation.RiskAcceptance != DecisionNotAuthorized {
		t.Fatalf("authority defaults drifted: %+v", validation)
	}
	if len(validation.PlaneResults) != 4 || validation.PlaneResults[0].Plane != PlaneCodeCorrectness || validation.PlaneResults[1].Plane != PlanePrivacySafety || validation.PlaneResults[2].Plane != PlaneSecurity || validation.PlaneResults[3].Plane != PlaneTraceEvidence {
		t.Fatalf("required plane sorting/dedup drifted: %+v", validation.PlaneResults)
	}
	if validation.PlaneResults[0].RunID != "run-code-good" || !validation.PlaneResults[0].Usable {
		t.Fatalf("best usable code result not selected: %+v", validation.PlaneResults[0])
	}
	if validation.PlaneResults[1].RunID != "run-privacy-failed" || validation.PlaneResults[1].Status != StatusParseFailed || validation.PlaneResults[1].Usable {
		t.Fatalf("non-usable privacy result not selected: %+v", validation.PlaneResults[1])
	}
	if validation.PlaneResults[2].RunID != "" || validation.PlaneResults[2].Status != StateNotAssessed || validation.PlaneResults[2].Reason != "required_plane_not_assessed" || validation.PlaneResults[2].Usable || validation.PlaneResults[2].NextAction != "Run or import a reviewer result for this plane." {
		t.Fatalf("missing required plane fallback drifted: %+v", validation.PlaneResults[2])
	}
	if validation.PlaneResults[3].RunID != "run-trace-stale" || !validation.PlaneResults[3].Usable {
		t.Fatalf("trace plane result not selected: %+v", validation.PlaneResults[3])
	}
	reasons := strings.Join(validation.Reasons, ",")
	if !strings.Contains(reasons, "packet_digest_mismatch") || !strings.Contains(reasons, "result_packet_digest_mismatch:run-trace-stale") || !strings.Contains(reasons, PlanePrivacySafety+":"+StatusParseFailed) || !strings.Contains(reasons, PlaneSecurity+":"+StateNotAssessed) {
		t.Fatalf("digest mismatch reasons missing: %+v", validation.Reasons)
	}
	packetDigestReasonCount := 0
	for _, reason := range validation.Reasons {
		if reason == "packet_digest_mismatch" {
			packetDigestReasonCount++
		}
	}
	if packetDigestReasonCount != 1 {
		t.Fatalf("reasons should be unique: %+v", validation.Reasons)
	}
	actions := strings.Join(validation.NextActions, "\n")
	if !strings.Contains(actions, "Create a new packet and rerun review") || !strings.Contains(actions, "Discard stale reviewer results") || !strings.Contains(actions, "Rerun with JSON-only output matching the required schema.") || !strings.Contains(actions, "Run or import a reviewer result for this plane.") {
		t.Fatalf("validation next actions missing: %+v", validation.NextActions)
	}
}

func TestPrreviewValidationRankingModelAndLedgerFindingContracts(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("e")
	packet := Packet{
		SchemaVersion: SchemaVersionPacket,
		PacketID:      "packet-1",
		PacketDigest:  packetDigest,
		DiffRef:       SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone},
		CIState:       StateNotAssessed,
	}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "validation-ranking",
		RequiredPlanes: []string{PlaneCodeCorrectness, PlaneSecurity, PlanePrivacySafety},
		Roles: []ReviewRole{
			{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
			{RoleID: "security", Plane: PlaneSecurity, Runner: RunnerManualExternal, RequestedModel: "model-a"},
			{RoleID: "privacy", Plane: PlanePrivacySafety, Runner: RunnerManualExternal, RequestedModel: "not_assessed"},
		},
	}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{
		{ReviewRunID: "run-code-clean", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-code-clean")},
		{ReviewRunID: "run-code-findings", PacketDigest: packetDigest, Plane: PlaneCodeCorrectness, RoleID: "code", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusFindingsReported, RawOutputRef: retainedRawRef("run-code-findings")},
		{ReviewRunID: "run-privacy-not-assessed", PacketDigest: packetDigest, Plane: PlanePrivacySafety, RoleID: "privacy", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StateNotAssessed},
		{ReviewRunID: "run-privacy-timeout", PacketDigest: packetDigest, Plane: PlanePrivacySafety, RoleID: "privacy", Runner: RunnerManualExternal, RequestedModel: "not_assessed", ObservedModel: "not_assessed", ModelFamily: "not_assessed", ModelVersion: "not_assessed", Status: StatusTimedOut},
		{ReviewRunID: "run-security-mismatch", PacketDigest: packetDigest, Plane: PlaneSecurity, RoleID: "security", Runner: RunnerManualExternal, RequestedModel: "model-a", ObservedModel: "model-b", ModelFamily: "family", ModelVersion: "v1", Status: StatusNoFindings, RawOutputRef: retainedRawRef("run-security-mismatch")},
	}}
	ledger := Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: packetDigest, Findings: []LedgerFinding{
		{ID: "F-valid", Severity: SeverityMajor, Summary: "SYNTHETIC_TOKEN_SECRET_RANKING", Citation: Citation{ContextRefID: "diff", DiffHunkID: "hunk-1"}, Disposition: DispositionUnresolvedReviewBlocker},
		{ID: "F-invalid", Severity: SeverityMinor, Summary: "citation missing", Disposition: DispositionDeferredNotAssessed},
	}}

	validation := Validate(packet, profile, runs, ledger)
	if validation.ReviewCoverageState != CoverageCannotVerify {
		t.Fatalf("coverage = %s want cannot_verify validation=%+v", validation.ReviewCoverageState, validation)
	}
	byPlane := map[string]PlaneResult{}
	for _, result := range validation.PlaneResults {
		byPlane[result.Plane] = result
	}
	if byPlane[PlaneCodeCorrectness].RunID != "run-code-findings" || byPlane[PlaneCodeCorrectness].Status != StatusFindingsReported || !byPlane[PlaneCodeCorrectness].Usable {
		t.Fatalf("usable findings result should outrank clean usable retry: %+v", byPlane[PlaneCodeCorrectness])
	}
	if byPlane[PlanePrivacySafety].RunID != "run-privacy-timeout" || byPlane[PlanePrivacySafety].Status != StatusTimedOut || byPlane[PlanePrivacySafety].Usable {
		t.Fatalf("cannot-verify non-usable result should outrank not_assessed: %+v", byPlane[PlanePrivacySafety])
	}
	if byPlane[PlaneSecurity].RunID != "run-security-mismatch" || byPlane[PlaneSecurity].Status != StatusCannotVerify || byPlane[PlaneSecurity].Reason != "model_identity_mismatch" || byPlane[PlaneSecurity].Usable {
		t.Fatalf("model identity mismatch projection drifted: %+v", byPlane[PlaneSecurity])
	}
	if !containsString(validation.NextActions, "Rerun the reviewer or record fallback provenance for the observed model.") {
		t.Fatalf("model mismatch next action missing: %+v", validation.NextActions)
	}
	if !containsString(validation.Reasons, "finding_citation_cannot_verify") {
		t.Fatalf("unresolvable finding citation reason missing: %+v", validation.Reasons)
	}
	if validation.Findings[0].Summary != redactedUnsafeReviewerText {
		t.Fatalf("unsafe ledger summary was not sanitized: %+v", validation.Findings[0])
	}
	if validation.Findings[0].Disposition != DispositionUnresolvedReviewBlocker {
		t.Fatalf("unresolved blocker disposition drifted: %+v", validation.Findings[0])
	}
	for _, severity := range []string{SeverityCritical, SeverityMajor} {
		finding := LedgerFinding{Severity: severity, Disposition: DispositionUnresolvedReviewBlocker}
		if !ledgerFindingUnresolved(finding) {
			t.Fatalf("%s unresolved blocker should affect validation coverage", severity)
		}
	}
}

func TestPrreviewCoverageModelAndSummaryRenderingContracts(t *testing.T) {
	required := map[string]bool{PlaneCodeCorrectness: true, PlaneSecurity: true}
	if got := reviewCoverageState(required, 1, false, false); got != CoveragePartial {
		t.Fatalf("partial coverage = %s want %s", got, CoveragePartial)
	}
	if got := reviewCoverageState(required, 2, false, true); got != CoverageUnresolved {
		t.Fatalf("unresolved coverage = %s want %s", got, CoverageUnresolved)
	}
	if got := reviewCoverageState(required, 2, false, false); got != CoverageSatisfied {
		t.Fatalf("satisfied coverage = %s want %s", got, CoverageSatisfied)
	}
	if got := reviewCoverageState(required, 2, true, true); got != CoverageCannotVerify {
		t.Fatalf("cannot_verify should dominate coverage, got %s", got)
	}
	if got := reviewCoverageState(required, 0, false, false); got != CoverageNotAssessed {
		t.Fatalf("zero usable coverage = %s want %s", got, CoverageNotAssessed)
	}
	if got := reviewCoverageState(map[string]bool{}, 0, false, false); got != CoverageNotAssessed {
		t.Fatalf("empty required coverage = %s want %s", got, CoverageNotAssessed)
	}
	if modelMismatchWithoutFallback(
		ReviewRole{RequestedModel: StateNotAssessed},
		ReviewerResult{RequestedModel: "model-a", ObservedModel: "model-b"},
	) {
		t.Fatal("role not_assessed requested model should not create a mismatch")
	}
	if !modelMismatchWithoutFallback(
		ReviewRole{RequestedModel: "model-a"},
		ReviewerResult{RequestedModel: "model-a", ObservedModel: "model-b"},
	) {
		t.Fatal("unexplained model mismatch should require cannot_verify projection")
	}
	if modelMismatchWithoutFallback(
		ReviewRole{RequestedModel: "model-a"},
		ReviewerResult{RequestedModel: "model-a", ObservedModel: "model-b", FallbackForModel: "model-a", FallbackReason: "primary_unavailable"},
	) {
		t.Fatal("fallback metadata should explain model mismatch")
	}

	validation := Validation{
		ReviewCoverageState: CoveragePartial,
		CIState:             StateNotAssessed,
		AuthorityScope:      AuthorityReviewRecordOnly,
		MergeDecision:       DecisionNotAuthorized,
		ReleaseDecision:     DecisionNotAuthorized,
		RiskAcceptance:      DecisionNotAuthorized,
		PlaneResults: []PlaneResult{{
			Plane:      PlaneSecurity,
			Status:     StatusCannotVerify,
			NextAction: "Investigate SYNTHETIC_TOKEN_SECRET_SUMMARY before merge.",
		}},
	}
	ledger := Ledger{Findings: []LedgerFinding{{
		ID:          "F-summary",
		Severity:    SeverityMajor,
		Summary:     "SYNTHETIC_PRIVATE_PATH_SUMMARY",
		Disposition: DispositionUnresolvedReviewBlocker,
	}}}
	summary := Summarize(validation, ledger)
	if !strings.Contains(summary, "Review coverage: "+CoveragePartial) || !strings.Contains(summary, "Merge decision: "+DecisionNotAuthorized) {
		t.Fatalf("summary boundary fields missing: %s", summary)
	}
	if !strings.Contains(summary, "This is review-record evidence only") {
		t.Fatalf("summary authority boundary missing: %s", summary)
	}
	if strings.Contains(summary, "SYNTHETIC_TOKEN_SECRET_SUMMARY") || strings.Contains(summary, "SYNTHETIC_PRIVATE_PATH_SUMMARY") {
		t.Fatalf("summary leaked unsafe text: %s", summary)
	}
	if !strings.Contains(summary, redactedUnsafeReviewerText) {
		t.Fatalf("summary missing redaction marker: %s", summary)
	}
	for _, forbidden := range []string{"safe to merge", "approved", "ready", "policy passed"} {
		if strings.Contains(strings.ToLower(summary), forbidden) {
			t.Fatalf("summary contains forbidden phrase %q: %s", forbidden, summary)
		}
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

func TestPrreviewLedgerSynthesisPreservesOrderingCarryForwardAndSanitization(t *testing.T) {
	packetDigest := "sha256:" + sixtyFour("9")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketDigest: packetDigest}
	existing := &Ledger{Findings: []LedgerFinding{
		{ID: "F-1", Disposition: DispositionUnresolvedReviewBlocker},
		{ID: "F-1", Disposition: DispositionAcceptedNarrower},
	}}
	runs := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest, Results: []ReviewerResult{
		{
			ReviewRunID: "run-z",
			Plane:       PlaneSecurity,
			RoleID:      "security",
			Findings: []Finding{{
				ID:           "F-2",
				Severity:     "unknown",
				Summary:      "SYNTHETIC_TOKEN_SECRET_LEDGER",
				Citation:     Citation{ContextRefID: "diff", DiffHunkID: "hunk-2"},
				EvidenceRefs: []string{"evidence-2"},
			}},
		},
		{
			ReviewRunID: "run-a",
			Plane:       PlaneCodeCorrectness,
			RoleID:      "code",
			Findings: []Finding{
				{
					ID:           "F-1",
					Severity:     SeverityCritical,
					Summary:      "blocking issue",
					Citation:     Citation{ContextRefID: "diff", DiffHunkID: "hunk-1"},
					EvidenceRefs: []string{"evidence-1"},
				},
				{
					Severity: SeverityMajor,
					Summary:  "fallback id issue",
				},
			},
		},
	}}

	ledger := SynthesizeLedger(packet, runs, existing)
	if ledger.SchemaVersion != SchemaVersionLedger || ledger.PacketDigest != packetDigest {
		t.Fatalf("ledger identity drifted: %+v", ledger)
	}
	if len(ledger.Findings) != 3 {
		t.Fatalf("ledger finding count = %d want 3: %+v", len(ledger.Findings), ledger.Findings)
	}
	gotIDs := []string{ledger.Findings[0].ID, ledger.Findings[1].ID, ledger.Findings[2].ID}
	wantIDs := []string{"F-1", "F-2", "run-a-finding"}
	if strings.Join(gotIDs, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("ledger findings not sorted by ID: got %v want %v", gotIDs, wantIDs)
	}
	byID := map[string]LedgerFinding{}
	for _, finding := range ledger.Findings {
		byID[finding.ID] = finding
	}
	if byID["F-1"].Disposition != DispositionAcceptedNarrower {
		t.Fatalf("existing duplicate last disposition should carry forward: %+v", byID["F-1"])
	}
	if byID["F-2"].Severity != SeverityInformational || byID["F-2"].Summary != redactedUnsafeReviewerText || byID["F-2"].Disposition != DispositionDeferredNotAssessed {
		t.Fatalf("severity/sanitization/default disposition drifted: %+v", byID["F-2"])
	}
	if byID["F-2"].ReviewRunID != "run-z" || byID["F-2"].Plane != PlaneSecurity || byID["F-2"].RoleID != "security" {
		t.Fatalf("reviewer result metadata not preserved: %+v", byID["F-2"])
	}
	if len(byID["F-2"].EvidenceRefs) != 1 || byID["F-2"].EvidenceRefs[0] != "evidence-2" {
		t.Fatalf("evidence refs not preserved: %+v", byID["F-2"])
	}
	if byID["run-a-finding"].Disposition != DispositionUnresolvedReviewBlocker {
		t.Fatalf("major fallback ID finding should default to unresolved blocker: %+v", byID["run-a-finding"])
	}
	empty := SynthesizeLedger(packet, RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packetDigest}, nil)
	if empty.SchemaVersion != SchemaVersionLedger || empty.PacketDigest != packetDigest || len(empty.Findings) != 0 {
		t.Fatalf("empty ledger behavior drifted: %+v", empty)
	}
}

func TestRunReviewArtifactPipelineRedactsUnsafeReviewerText(t *testing.T) {
	root := t.TempDir()
	packetDigest := "sha256:" + sixtyFour("9")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-1", PacketDigest: packetDigest, RepoID: "demo_repo", ChangeRef: "pr-123", BaseCommit: forty("a"), HeadCommit: forty("b"), DiffRef: SafeRef{ID: "diff", Kind: RefKindDiff, Ref: "inputs/diff.patch", DigestSHA256: sixtyFour("2"), ContentType: ContentUnifiedDiff, RedactionState: RedactionNone}, CIState: StateNotAssessed}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "artifact-safety",
		RequiredPlanes: []string{PlanePrivacySafety},
		Roles: []ReviewRole{{
			RoleID:             "privacy",
			Plane:              PlanePrivacySafety,
			Runner:             RunnerPI,
			RequestedModel:     "fake-pi",
			Command:            []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "unsafe-structured-output"},
			RawOutputRetention: RedactionDigestOnly,
		}},
	}
	t.Setenv("GO_WANT_PR_REVIEW_HELPER_PROCESS", "1")
	runs, _, err := RunReview(packet, profile, RunOptions{
		OutDir:         filepath.Join(root, "runs"),
		AllowedRunners: map[string]bool{RunnerPI: true},
		Now:            time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	ledger := SynthesizeLedger(packet, runs, nil)
	validation := Validate(packet, profile, runs, ledger)
	summary := Summarize(validation, ledger)
	artifactPaths := []string{
		filepath.Join(root, "runs", "results.json"),
		filepath.Join(root, "ledger.json"),
		filepath.Join(root, "validation.json"),
		filepath.Join(root, "summary.md"),
	}
	if err := WriteJSON(artifactPaths[1], ledger); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(artifactPaths[2], validation); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifactPaths[3], []byte(summary), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range artifactPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, marker := range []string{
			"SYNTHETIC_TOKEN_SECRET_PIPELINE",
			"SYNTHETIC_PROMPT_SECRET_PIPELINE",
			"https://access_token=secret@example.invalid/review",
			"/Users/private/repo",
		} {
			if strings.Contains(string(data), marker) {
				t.Fatalf("%s leaked marker %q:\n%s", path, marker, string(data))
			}
		}
		if !strings.Contains(string(data), "[redacted unsafe reviewer text]") && path != filepath.Join(root, "summary.md") {
			t.Fatalf("%s missing redaction marker:\n%s", path, string(data))
		}
		if strings.Contains(string(data), "SYNTHETIC_RATIONALE_SECRET_PIPELINE") {
			t.Fatalf("%s leaked rationale marker:\n%s", path, string(data))
		}
	}
	if _, err := os.Stat(filepath.Join(root, "runs", "raw", "run-privacy.out")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("digest-only raw output should not persist raw bytes, err=%v", err)
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
			{RoleID: "unknown-field", Plane: PlaneTraceEvidence, Runner: RunnerManualExternal, RequestedModel: "fake-unknown", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "unknown-field"}},
			{RoleID: "offtask", Plane: PlaneRequirements, Runner: RunnerManualExternal, RequestedModel: "fake-offtask", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "offtask"}},
			{RoleID: "minimal", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "fake-minimal", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "minimal-success"}, PromptTemplateRef: promptPath},
			{RoleID: "readonly", Plane: PlanePrivacySafety, Runner: RunnerOpenCode, RequestedModel: "fake-opencode", ReadOnlyEnforced: false, Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "success"}},
			{RoleID: "pi-success", Plane: PlaneSecurity, Runner: RunnerPI, RequestedModel: "fake-pi", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "pi-success"}, PromptTemplateRef: promptPath, RawOutputRetention: RedactionDigestOnly},
			{RoleID: "findings-default", Plane: PlaneSecurity, Runner: RunnerManualExternal, RequestedModel: "fake-findings", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "findings-no-status"}},
			{RoleID: "opencode-mutation", Plane: PlaneDXReplayability, Runner: RunnerOpenCode, RequestedModel: "fake-opencode", ReadOnlyEnforced: true, WorkingTreeMode: "clean_required", Command: []string{helper, "-test.run=TestPRReviewFakeRunnerHelper", "--", "opencode-mutation"}},
		},
	}
	t.Setenv("GO_WANT_PR_REVIEW_HELPER_PROCESS", "1")
	runs, preview, err := RunReview(packet, profile, RunOptions{
		OutDir:         filepath.Join(root, "runs"),
		AllowedRunners: map[string]bool{RunnerManualExternal: true, RunnerOpenCode: true, RunnerPI: true},
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
		"unknown-field":     StatusParseFailed,
		"offtask":           StatusOffTask,
		"minimal":           StatusNoFindings,
		"readonly":          StatusNotAssessed,
		"pi-success":        StatusNoFindings,
		"findings-default":  StatusFindingsReported,
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
	malformedRawPath := filepath.Join(root, "runs", "raw", "run-malformed.out")
	malformedRaw, err := os.ReadFile(malformedRawPath)
	if err != nil {
		t.Fatalf("retained raw output should be readable: %v", err)
	}
	malformedSum := sha256.Sum256(malformedRaw)
	if statuses["malformed"].RawOutputRef.Ref != "raw/run-malformed.out" || statuses["malformed"].RawOutputRef.DigestSHA256 != hex.EncodeToString(malformedSum[:]) {
		t.Fatalf("retained raw output ref drifted: %+v", statuses["malformed"].RawOutputRef)
	}
	info, err := os.Stat(malformedRawPath)
	if err != nil {
		t.Fatalf("retained raw output stat failed: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("retained raw output mode = %v, want 0600", info.Mode().Perm())
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
	if statuses["unknown-field"].Reason != "runner_output_parse_failed" {
		t.Fatalf("unknown reviewer output field should parse-fail: %+v", statuses["unknown-field"])
	}
	minimal := statuses["minimal"]
	if minimal.ReviewRunID != "run-minimal" || minimal.Runner != RunnerManualExternal || minimal.RequestedModel != "fake-minimal" || minimal.ObservedModel != StateNotAssessed || minimal.ModelFamily != StateNotAssessed || minimal.ModelVersion != StateNotAssessed || minimal.Status != StatusNoFindings {
		t.Fatalf("parsed default propagation drifted: %+v", minimal)
	}
	if minimal.StartedAt == "" || minimal.EndedAt == "" || minimal.CommandDigest == "" {
		t.Fatalf("parsed execution metadata missing: %+v", minimal)
	}
	if minimal.PromptRef == nil || minimal.PromptRef.RedactionState != RedactionDigestOnly {
		t.Fatalf("parsed prompt ref was not retained: %+v", minimal)
	}
	if len(statuses["findings-default"].Findings) != 1 || statuses["findings-default"].Status != StatusFindingsReported {
		t.Fatalf("parsed findings status default drifted: %+v", statuses["findings-default"])
	}
}

func TestRunReviewPreviewReturnsPreviewOnly(t *testing.T) {
	root := t.TempDir()
	promptPath := writeText(t, root, "prompt.md", "review {{packet_digest}}\n")
	packet := Packet{PacketDigest: "sha256:" + sixtyFour("v"), SchemaVersion: SchemaVersionPacket}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "preview",
		RequiredPlanes: []string{PlaneCodeCorrectness, PlaneTraceEvidence},
		Roles: []ReviewRole{
			{
				RoleID:            "code",
				Plane:             PlaneCodeCorrectness,
				Runner:            RunnerManualExternal,
				RequestedModel:    "not_assessed",
				TimeoutSeconds:    120,
				Command:           []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper"},
				PromptTemplateRef: promptPath,
			},
			{
				RoleID:         "trace",
				Plane:          PlaneTraceEvidence,
				Runner:         RunnerManualExternal,
				RequestedModel: "",
				TimeoutSeconds: 240,
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
	if len(preview.Roles) != 2 {
		t.Fatalf("preview roles = %d want 2", len(preview.Roles))
	}
	if preview.Roles[0].RoleID != "code" || preview.Roles[1].RoleID != "trace" {
		t.Fatalf("preview role order drifted: %+v", preview.Roles)
	}
	if preview.Roles[0].CommandDigest == "" || preview.Roles[0].PromptDigest == "" || preview.Roles[0].PromptRef != promptPath {
		t.Fatalf("preview command/prompt digest drifted: %+v", preview.Roles[0])
	}
	if preview.Roles[1].RequestedModel != StateNotAssessed || preview.Roles[1].PromptRef != "" || preview.Roles[1].PromptDigest != "" {
		t.Fatalf("preview empty prompt/default model drifted: %+v", preview.Roles[1])
	}
	if _, err := os.Stat(outDir); !os.IsNotExist(err) {
		t.Fatalf("preview should not create output directory, got stat err=%v", err)
	}
}

func TestRunReviewPreservesValidationDefaultsAndOutputContracts(t *testing.T) {
	root := t.TempDir()
	packetDigest := "sha256:" + sixtyFour("5")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketDigest: packetDigest}
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "contracts",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles: []ReviewRole{{
			RoleID:         "code",
			Plane:          PlaneCodeCorrectness,
			Runner:         RunnerManualExternal,
			RequestedModel: "manual",
		}},
	}

	defaults := normalizeRunOptions(RunOptions{})
	if defaults.Now.IsZero() || defaults.WorkDir != "." {
		t.Fatalf("run option defaults drifted: %+v", defaults)
	}

	invalidOutDir := filepath.Join(root, "invalid-profile-runs")
	invalidProfile := cloneReviewProfile(profile)
	invalidProfile.SchemaVersion = "bad"
	_, _, err := RunReview(packet, invalidProfile, RunOptions{OutDir: invalidOutDir})
	if err == nil || err.Error() != "invalid_profile_schema_version: bad" {
		t.Fatalf("expected profile validation error before run preparation, got %v", err)
	}
	if _, err := os.Stat(invalidOutDir); !os.IsNotExist(err) {
		t.Fatalf("invalid profile should not create output directory, got %v", err)
	}

	_, _, err = RunReview(packet, profile, RunOptions{OutDir: " \t\n"})
	if err == nil || err.Error() != "missing_output_path" {
		t.Fatalf("expected missing output path contract, got %v", err)
	}

	existingOutDir := filepath.Join(root, "existing-runs")
	if err := os.MkdirAll(existingOutDir, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(existingOutDir, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err = RunReview(packet, profile, RunOptions{OutDir: existingOutDir, NotAssessedReason: "configured_elsewhere"})
	if err == nil || err.Error() != "output_exists: existing-runs" {
		t.Fatalf("expected existing output directory rejection contract, got %v", err)
	}
	if data, readErr := os.ReadFile(sentinel); readErr != nil || string(data) != "keep\n" {
		t.Fatalf("existing output directory contents should remain untouched, data=%q err=%v", string(data), readErr)
	}

	outDir := filepath.Join(root, "runs")
	runs, preview, err := RunReview(packet, profile, RunOptions{OutDir: outDir, NotAssessedReason: "configured_elsewhere"})
	if err != nil {
		t.Fatal(err)
	}
	if preview != nil {
		t.Fatalf("unexpected preview: %+v", preview)
	}
	if runs.SchemaVersion != SchemaVersionRunSet || runs.PacketDigest != packetDigest || len(runs.Results) != 1 {
		t.Fatalf("run-set shape drifted: %+v", runs)
	}
	if runs.Results[0].RoleID != "code" || runs.Results[0].Status != StatusNotAssessed || runs.Results[0].Reason != "configured_elsewhere" {
		t.Fatalf("run result shape drifted: %+v", runs.Results[0])
	}
	if _, err := os.Stat(filepath.Join(outDir, "raw")); err != nil {
		t.Fatalf("raw directory missing: %v", err)
	}
	written, err := ReadRunSet(filepath.Join(outDir, "results.json"))
	if err != nil {
		t.Fatalf("results.json should be readable: %v", err)
	}
	if written.PacketDigest != packetDigest || len(written.Results) != 1 || written.Results[0].RoleID != "code" {
		t.Fatalf("written results shape drifted: %+v", written)
	}

	orderedProfile := cloneReviewProfile(profile)
	orderedProfile.RequiredPlanes = []string{PlaneCodeCorrectness, PlaneTraceEvidence}
	orderedProfile.Roles = []ReviewRole{
		{RoleID: "first", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal, RequestedModel: "manual"},
		{RoleID: "second", Plane: PlaneTraceEvidence, Runner: RunnerManualExternal, RequestedModel: "manual"},
	}
	orderedRuns, _, err := RunReview(packet, orderedProfile, RunOptions{OutDir: filepath.Join(root, "ordered-runs"), NotAssessedReason: "configured_elsewhere"})
	if err != nil {
		t.Fatal(err)
	}
	if len(orderedRuns.Results) != 2 || orderedRuns.Results[0].RoleID != "first" || orderedRuns.Results[1].RoleID != "second" {
		t.Fatalf("profile role order drifted: %+v", orderedRuns.Results)
	}

	noCommandProfile := cloneReviewProfile(profile)
	noCommandRuns, _, err := RunReview(packet, noCommandProfile, RunOptions{OutDir: filepath.Join(root, "no-command-runs")})
	if err != nil {
		t.Fatal(err)
	}
	if len(noCommandRuns.Results) != 1 || noCommandRuns.Results[0].Status != StatusNotAssessed || noCommandRuns.Results[0].Reason != "runner_command_not_configured" {
		t.Fatalf("no-command role state drifted: %+v", noCommandRuns.Results)
	}
	if _, err := roleCommand(context.Background(), noCommandProfile.Roles[0], root); err == nil || !strings.Contains(err.Error(), "runner_command_not_configured") {
		t.Fatalf("empty command guard drifted: %v", err)
	}

	dirtyWorkDir := filepath.Join(root, "dirty-work")
	if err := os.MkdirAll(dirtyWorkDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runGitInit(dirtyWorkDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirtyWorkDir, "dirty.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirtyProfile := cloneReviewProfile(profile)
	dirtyProfile.Roles[0] = ReviewRole{
		RoleID:           "dirty",
		Plane:            PlaneCodeCorrectness,
		Runner:           RunnerOpenCode,
		RequestedModel:   "fake-opencode",
		ReadOnlyEnforced: true,
		WorkingTreeMode:  "clean_required",
		Command:          []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "should-not-run"},
	}
	dirtyRuns, _, err := RunReview(packet, dirtyProfile, RunOptions{
		OutDir:         filepath.Join(root, "dirty-runs"),
		AllowedRunners: map[string]bool{RunnerOpenCode: true},
		WorkDir:        dirtyWorkDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirtyRuns.Results) != 1 || dirtyRuns.Results[0].Status != StatusNotAssessed || dirtyRuns.Results[0].Reason != "working_tree_dirty" {
		t.Fatalf("dirty OpenCode baseline state drifted: %+v", dirtyRuns.Results)
	}

	errorOutDir := filepath.Join(root, "runner-error")
	disallowed := cloneReviewProfile(profile)
	disallowed.Roles[0].Runner = RunnerPI
	_, _, err = RunReview(packet, disallowed, RunOptions{OutDir: errorOutDir})
	if err == nil || !strings.Contains(err.Error(), "runner_not_allowed: pi") {
		t.Fatalf("expected role execution error propagation, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(errorOutDir, "results.json")); !os.IsNotExist(err) {
		t.Fatalf("role execution error should not write successful results.json, got %v", err)
	}

	manualCommandOutDir := filepath.Join(root, "manual-command-disallowed")
	manualCommandProfile := cloneReviewProfile(profile)
	manualCommandProfile.Roles = []ReviewRole{{
		RoleID:         "manual-command",
		Plane:          PlaneCodeCorrectness,
		Runner:         RunnerManualExternal,
		RequestedModel: "manual",
		Command:        []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "should-not-run"},
	}}
	_, _, err = RunReview(packet, manualCommandProfile, RunOptions{OutDir: manualCommandOutDir})
	if err == nil || !strings.Contains(err.Error(), "runner_not_allowed: manual_external") {
		t.Fatalf("manual command should require explicit runner allowance, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(manualCommandOutDir, "results.json")); !os.IsNotExist(err) {
		t.Fatalf("disallowed manual command should not write successful results.json, got %v", err)
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

func TestPrreviewArtifactIOReadWriteContracts(t *testing.T) {
	root := t.TempDir()
	outDir := filepath.Join(root, "nested", "artifacts")
	packet := Packet{SchemaVersion: SchemaVersionPacket, PacketID: "packet-io", PacketDigest: "sha256:" + sixtyFour("1")}
	packetPath := filepath.Join(outDir, "packet.json")
	oldUmask := syscall.Umask(0)
	if err := WriteJSON(packetPath, packet); err != nil {
		syscall.Umask(oldUmask)
		t.Fatalf("WriteJSON(packet) error = %v", err)
	}
	syscall.Umask(oldUmask)
	dirInfo, err := os.Stat(outDir)
	if err != nil {
		t.Fatalf("stat output dir: %v", err)
	}
	if got := dirInfo.Mode().Perm(); got != 0o755 {
		t.Fatalf("output dir mode = %#o want 0755", got)
	}
	fileInfo, err := os.Stat(packetPath)
	if err != nil {
		t.Fatalf("stat packet: %v", err)
	}
	if got := fileInfo.Mode().Perm(); got != 0o644 {
		t.Fatalf("packet file mode = %#o want 0644", got)
	}
	data, err := os.ReadFile(packetPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(data), "\n") || !strings.Contains(string(data), "\n  \"packet_id\": \"packet-io\"") {
		t.Fatalf("written JSON should be indented and newline terminated:\n%s", string(data))
	}
	readPacket, err := ReadPacket(outDir)
	if err != nil {
		t.Fatalf("ReadPacket(dir) error = %v", err)
	}
	if readPacket.PacketID != "packet-io" {
		t.Fatalf("packet read drifted: %+v", readPacket)
	}

	ledger := Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: packet.PacketDigest}
	ledgerPath := filepath.Join(outDir, "ledger.json")
	if err := WriteJSON(ledgerPath, ledger); err != nil {
		t.Fatalf("write ledger: %v", err)
	}
	readLedger, err := ReadLedger(ledgerPath)
	if err != nil {
		t.Fatalf("ReadLedger() error = %v", err)
	}
	if readLedger.SchemaVersion != SchemaVersionLedger || readLedger.PacketDigest != packet.PacketDigest {
		t.Fatalf("ledger read drifted: %+v", readLedger)
	}

	validation := Validation{SchemaVersion: SchemaVersionValidation, PacketDigest: packet.PacketDigest, ReviewCoverageState: CoverageNotAssessed}
	validationPath := filepath.Join(outDir, "validation.json")
	if err := WriteJSON(validationPath, validation); err != nil {
		t.Fatalf("write validation: %v", err)
	}
	readValidation, err := ReadValidation(validationPath)
	if err != nil {
		t.Fatalf("ReadValidation() error = %v", err)
	}
	if readValidation.SchemaVersion != SchemaVersionValidation || readValidation.ReviewCoverageState != CoverageNotAssessed {
		t.Fatalf("validation read drifted: %+v", readValidation)
	}

	badPath := filepath.Join(outDir, "bad.json")
	if err := os.WriteFile(badPath, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadLedger(badPath); err == nil {
		t.Fatal("ReadLedger should return invalid JSON error")
	}
	if _, err := ReadValidation(filepath.Join(outDir, "missing.json")); err == nil {
		t.Fatal("ReadValidation should return missing file error")
	}
}

func TestPrreviewOptionSchemaValidationContracts(t *testing.T) {
	validPacket := PacketOptions{
		OutDir:     "out",
		RepoID:     "demo_repo",
		ChangeRef:  "pr-123",
		BaseCommit: forty("a"),
		HeadCommit: forty("b"),
		DiffPath:   "diff.patch",
		CIState:    StatePass,
	}
	for name, tc := range map[string]struct {
		mutate  func(*PacketOptions)
		wantErr string
	}{
		"missing-out": {
			mutate:  func(opts *PacketOptions) { opts.OutDir = " " },
			wantErr: "pr_review_packet_requires_out",
		},
		"unsafe-repo": {
			mutate:  func(opts *PacketOptions) { opts.RepoID = "/tmp/private/repo" },
			wantErr: "unsafe_repo_id: repo_id must match " + repoIDPattern.String(),
		},
		"unsafe-change": {
			mutate:  func(opts *PacketOptions) { opts.ChangeRef = "branch/main" },
			wantErr: "unsafe_change_ref: change_ref must match " + changeRefPattern.String(),
		},
		"bad-commit": {
			mutate:  func(opts *PacketOptions) { opts.BaseCommit = forty("A") },
			wantErr: "invalid_commit_sha: base and head must be 40 lowercase hex characters",
		},
		"missing-diff": {
			mutate:  func(opts *PacketOptions) { opts.DiffPath = " " },
			wantErr: "pr_review_packet_requires_diff",
		},
		"bad-ci": {
			mutate:  func(opts *PacketOptions) { opts.CIState = "green" },
			wantErr: "invalid_ci_state: green",
		},
	} {
		t.Run("packet-"+name, func(t *testing.T) {
			opts := validPacket
			tc.mutate(&opts)
			if err := validatePacketOptions(opts); err == nil || err.Error() != tc.wantErr {
				t.Fatalf("validatePacketOptions() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
	if err := validatePacketOptions(validPacket); err != nil {
		t.Fatalf("valid packet options rejected: %v", err)
	}

	if err := validateRunSet(RunSet{Results: []ReviewerResult{{ReviewRunID: " "}}}); err == nil || err.Error() != "review_result_requires_review_run_id" {
		t.Fatalf("missing run id error = %v", err)
	}
	if err := validateRunSet(RunSet{Results: []ReviewerResult{{ReviewRunID: "run-1"}, {ReviewRunID: "run-1"}}}); err == nil || err.Error() != "duplicate_review_run_id: run-1" {
		t.Fatalf("duplicate run id error = %v", err)
	}

	validProfile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "profile",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles:          []ReviewRole{{RoleID: "code", Plane: PlaneCodeCorrectness, Runner: RunnerManualExternal}},
	}
	profileCases := map[string]struct {
		mutate  func(*ReviewProfile)
		wantErr string
	}{
		"bad-schema": {
			mutate:  func(profile *ReviewProfile) { profile.SchemaVersion = "bad" },
			wantErr: "invalid_profile_schema_version: bad",
		},
		"missing-profile-id": {
			mutate:  func(profile *ReviewProfile) { profile.ProfileID = " " },
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
		"missing-role-field": {
			mutate:  func(profile *ReviewProfile) { profile.Roles[0].RoleID = "" },
			wantErr: "profile_role_requires_id_plane_runner",
		},
		"invalid-runner": {
			mutate:  func(profile *ReviewProfile) { profile.Roles[0].Runner = "unknown" },
			wantErr: "profile_role_invalid_runner: unknown",
		},
		"required-plane-without-role": {
			mutate:  func(profile *ReviewProfile) { profile.RequiredPlanes = []string{PlaneCodeCorrectness, PlaneSecurity} },
			wantErr: "profile_required_plane_without_role: security_forgery_overclaim",
		},
	}
	for name, tc := range profileCases {
		t.Run("profile-"+name, func(t *testing.T) {
			profile := cloneReviewProfile(validProfile)
			tc.mutate(&profile)
			if err := validateProfile(profile); err == nil || err.Error() != tc.wantErr {
				t.Fatalf("validateProfile() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
	withoutSchema := cloneReviewProfile(validProfile)
	withoutSchema.SchemaVersion = ""
	if err := validateProfile(withoutSchema); err != nil {
		t.Fatalf("empty profile schema version should be accepted: %v", err)
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
	tests := []struct {
		name       string
		err        error
		wantStatus string
		wantReason string
	}{
		{name: "prompt template", err: errPromptTemplateCannotVerify, wantStatus: StatusCannotVerify, wantReason: "prompt_ref_cannot_verify"},
		{name: "prompt evidence", err: errPromptEvidenceCannotVerify, wantStatus: StatusCannotVerify, wantReason: "prompt_evidence_cannot_verify"},
		{name: "runner unavailable", err: exec.ErrNotFound, wantStatus: StatusNotAssessed, wantReason: "runner_unavailable"},
		{name: "runner failed", err: fmt.Errorf("boom"), wantStatus: StatusFailed, wantReason: "runner_failed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ReviewerResult{}
			if err := applyRunnerError(&result, test.err); err == nil {
				t.Fatalf("expected error returned")
			}
			if result.Status != test.wantStatus || result.Reason != test.wantReason {
				t.Fatalf("result = %+v, want status=%s reason=%s", result, test.wantStatus, test.wantReason)
			}
		})
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
	runs, _, err := RunReview(packet, profile, RunOptions{
		OutDir:         filepath.Join(root, "runs"),
		AllowedRunners: map[string]bool{RunnerManualExternal: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(runs.Results) != 1 || runs.Results[0].Status != StatusCannotVerify || runs.Results[0].Reason != "prompt_ref_cannot_verify" {
		t.Fatalf("missing prompt should be cannot_verify without runner execution: %+v", runs.Results)
	}
}

func TestRunReviewPromptIncludesPacketEvidence(t *testing.T) {
	root := t.TempDir()
	diffPath := writeText(t, root, "change.diff", "diff --git a/a.go b/a.go\n+package main\n")
	metadataPath := writeText(t, root, "metadata.txt", "PR #123\n")
	verificationPath := writeText(t, root, "verification.txt", "verify: pass\n")
	packetDir := filepath.Join(root, "packet")
	packet, err := BuildPacket(PacketOptions{
		OutDir:            packetDir,
		RepoID:            "demo_repo",
		ChangeRef:         "pr-123",
		BaseCommit:        forty("a"),
		HeadCommit:        forty("b"),
		DiffPath:          diffPath,
		MetadataPath:      metadataPath,
		VerificationPaths: []string{verificationPath},
	})
	if err != nil {
		t.Fatal(err)
	}
	promptPath := writeText(t, root, "prompt.md", "review {{packet_digest}}\n")
	profile := ReviewProfile{
		SchemaVersion:  SchemaVersionProfile,
		ProfileID:      "prompt-evidence",
		RequiredPlanes: []string{PlaneCodeCorrectness},
		Roles: []ReviewRole{{
			RoleID:            "code",
			Plane:             PlaneCodeCorrectness,
			Runner:            RunnerManualExternal,
			RequestedModel:    "fake",
			Command:           []string{os.Args[0], "-test.run=TestPRReviewFakeRunnerHelper", "--", "prompt-includes-evidence"},
			PromptTemplateRef: promptPath,
		}},
	}
	t.Setenv("GO_WANT_PR_REVIEW_HELPER_PROCESS", "1")
	t.Setenv("PR_REVIEW_EXPECTED_DIGEST", packet.PacketDigest)
	runs, _, err := RunReview(packet, profile, RunOptions{
		OutDir:         filepath.Join(root, "runs"),
		PacketDir:      packetDir,
		AllowedRunners: map[string]bool{RunnerManualExternal: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(runs.Results) != 1 || runs.Results[0].Status != StatusNoFindings {
		t.Fatalf("prompt evidence run failed: %+v", runs.Results)
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
	runs, _, err := RunReview(packet, profile, RunOptions{
		OutDir:         filepath.Join(root, "runs"),
		AllowedRunners: map[string]bool{RunnerManualExternal: true},
	})
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

func TestPrreviewSmallHelpersPreserveContracts(t *testing.T) {
	if commandDigest(nil) != "" {
		t.Fatalf("empty command digest should be empty")
	}
	wantCommandDigest := sha256.Sum256([]byte("go\x00test\x00./internal/prreview"))
	if got := commandDigest([]string{"go", "test", "./internal/prreview"}); got != "sha256:"+hex.EncodeToString(wantCommandDigest[:]) {
		t.Fatalf("commandDigest() = %q", got)
	}
	if got := defaultString(" \t\n", "fallback"); got != "fallback" {
		t.Fatalf("defaultString whitespace fallback = %q", got)
	}
	if got := defaultString(" value ", "fallback"); got != " value " {
		t.Fatalf("defaultString should preserve non-empty input, got %q", got)
	}

	for path, want := range map[string]string{
		"a.JSON":   ".json",
		"a.md":     ".md",
		"a.txt":    ".txt",
		"a.diff":   ".diff",
		"a.patch":  ".patch",
		"a.log":    ".txt",
		"Makefile": ".txt",
	} {
		if got := normalizedExt(path); got != want {
			t.Fatalf("normalizedExt(%q) = %q want %q", path, got, want)
		}
	}
	for path, want := range map[string]string{
		"payload.json":   ContentJSON,
		"notes.md":       ContentMarkdown,
		"notes.markdown": ContentMarkdown,
		"run.log":        ContentText,
	} {
		if got := contentType(path); got != want {
			t.Fatalf("contentType(%q) = %q want %q", path, got, want)
		}
	}
	for path, want := range map[string]string{
		"task-review.md": RefKindTask,
		"notes.md":       RefKindDoc,
		"schema.json":    RefKindSchema,
		"diff.patch":     RefKindSourceExcerpt,
	} {
		if got := contextKind(path); got != want {
			t.Fatalf("contextKind(%q) = %q want %q", path, got, want)
		}
	}
}

func TestSafeEvidenceRefPreservesRefsAndRedactsUnsafeText(t *testing.T) {
	if got := safeEvidenceRef("diff"); got != "diff" {
		t.Fatalf("safe ref = %q", got)
	}
	if got := safeEvidenceRef("context.ref-1"); got != "context.ref-1" {
		t.Fatalf("safe dotted ref = %q", got)
	}
	if got := safeEvidenceRef("https://access_token=secret@example.invalid/review"); got != "[redacted unsafe reviewer text]" {
		t.Fatalf("unsafe ref = %q", got)
	}
}

func TestPrreviewPromptSanitizerAndCopyHelpersPreserveContracts(t *testing.T) {
	var copied strings.Builder
	if err := Copy(strings.NewReader("review bytes"), &copied); err != nil {
		t.Fatalf("Copy returned error: %v", err)
	}
	if copied.String() != "review bytes" {
		t.Fatalf("Copy wrote %q", copied.String())
	}

	packet := Packet{
		PacketDigest: "sha256:" + sixtyFour("p"),
		RepoID:       "repo",
		ChangeRef:    "change",
		BaseCommit:   forty("a"),
		HeadCommit:   forty("b"),
		CIState:      StatePass,
	}
	role := ReviewRole{RoleID: "code", Plane: PlaneCodeCorrectness}
	rendered, err := renderPromptTemplate(packet, role)
	if err != nil || rendered != "" {
		t.Fatalf("empty prompt template rendered=%q err=%v", rendered, err)
	}
	role.PromptTemplateRef = filepath.Join(t.TempDir(), "missing.md")
	if _, err := renderPromptTemplate(packet, role); !errors.Is(err, errPromptTemplateCannotVerify) {
		t.Fatalf("missing prompt template err=%v", err)
	}

	rendered = applyPromptReplacements("{{packet_digest}} {{repo_id}} {{role_id}} {{unknown}}", packet, ReviewRole{RoleID: "code"})
	if rendered != packet.PacketDigest+" repo code {{unknown}}" {
		t.Fatalf("prompt replacements drifted: %q", rendered)
	}
	if got := replacePromptToken("x {{role_id}} x", promptReplacement{key: "role_id", value: "security"}); got != "x security x" {
		t.Fatalf("replacePromptToken = %q", got)
	}

	var packetJSON strings.Builder
	if err := appendPromptPacketJSON(&packetJSON, packet); err != nil {
		t.Fatalf("appendPromptPacketJSON: %v", err)
	}
	if !strings.Contains(packetJSON.String(), "Review packet JSON:") || !strings.Contains(packetJSON.String(), "```json") || !strings.Contains(packetJSON.String(), packet.PacketDigest) {
		t.Fatalf("packet JSON prompt block drifted:\n%s", packetJSON.String())
	}

	packetDir := t.TempDir()
	writePacketRef := func(name, data string) SafeRef {
		t.Helper()
		path := filepath.Join(packetDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256([]byte(data))
		return SafeRef{ID: safeID(name), Ref: name, DigestSHA256: hex.EncodeToString(sum[:]), ContentType: contentType(name)}
	}
	diffRef := writePacketRef("inputs/diff.patch", "diff --git a/a.go b/a.go\n+package main\n")
	diffRef.ID = "diff"
	diffRef.ContentType = ContentUnifiedDiff
	metadataRef := writePacketRef("metadata.json", `{"pr":123}`)
	contextRef := writePacketRef("context.md", "# Context\n")
	verificationRef := writePacketRef("verify.txt", "verify: pass\n")
	packet.DiffRef = diffRef
	packet.MetadataRef = &metadataRef
	packet.ContextRefs = []SafeRef{contextRef}
	packet.VerificationRefs = []SafeRef{verificationRef}

	var evidence strings.Builder
	if err := appendPromptEvidenceRefs(&evidence, packetDir, promptEvidenceRefs(packet)); err != nil {
		t.Fatalf("appendPromptEvidenceRefs: %v", err)
	}
	evidenceText := evidence.String()
	ordered := []string{"diff ref diff", "metadata ref metadata.json", "context ref context.md", "verification ref verify.txt"}
	last := -1
	for _, marker := range ordered {
		idx := strings.Index(evidenceText, marker)
		if idx <= last {
			t.Fatalf("prompt evidence order drifted for %q in:\n%s", marker, evidenceText)
		}
		last = idx
	}
	for _, fence := range []string{"```diff", "```json", "```text"} {
		if !strings.Contains(evidenceText, fence) {
			t.Fatalf("missing fence %q in:\n%s", fence, evidenceText)
		}
	}

	if _, err := readPacketRef(packetDir, SafeRef{ID: "unsafe", Ref: "../secret.txt", DigestSHA256: sixtyFour("0")}); !errors.Is(err, errPromptEvidenceCannotVerify) {
		t.Fatalf("unsafe packet ref err=%v", err)
	}
	if _, err := readPacketRef(packetDir, SafeRef{ID: "bad-digest", Ref: diffRef.Ref, DigestSHA256: sixtyFour("0")}); !errors.Is(err, errPromptEvidenceCannotVerify) {
		t.Fatalf("digest mismatch err=%v", err)
	}

	result := sanitizeReviewerResult(ReviewerResult{
		Reason: "Bearer SYNTHETIC_SECRET",
		Findings: []Finding{{
			Summary:      "ok",
			Rationale:    "/Users/private/path",
			SuggestedFix: "token=secret",
			Question:     "https://access_token=secret@example.invalid/review",
			EvidenceRefs: []string{"diff", "https://access_token=secret@example.invalid/review"},
		}},
	})
	if result.Reason != redactedUnsafeReviewerText || result.Findings[0].Rationale != redactedUnsafeReviewerText || result.Findings[0].SuggestedFix != redactedUnsafeReviewerText || result.Findings[0].Question != redactedUnsafeReviewerText {
		t.Fatalf("unsafe reviewer text was not redacted: %+v", result)
	}
	if result.Findings[0].Summary != "ok" || result.Findings[0].EvidenceRefs[0] != "diff" || result.Findings[0].EvidenceRefs[1] != redactedUnsafeReviewerText {
		t.Fatalf("safe/unsafe reviewer fields drifted: %+v", result.Findings[0])
	}
	if redactedUnsafeReviewerText != "[redacted unsafe reviewer text]" {
		t.Fatalf("redaction spelling drifted: %q", redactedUnsafeReviewerText)
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
	case "unknown-field":
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("4") + `","plane":"trace_evidence","role_id":"unknown-field","runner":"manual_external","requested_model":"fake","observed_model":"fake","model_family":"fake","model_version":"fake","status":"no_findings","findings":[],"extra_field":"reject-me"}`)
		os.Exit(0)
	case "offtask":
		fmt.Print(`{"packet_digest":"sha256:0000000000000000000000000000000000000000000000000000000000000000","plane":"requirements_vs_implementation","role_id":"offtask","runner":"manual_external","requested_model":"fake","observed_model":"fake","model_family":"fake","model_version":"fake","status":"no_findings","findings":[]}`)
		os.Exit(0)
	case "minimal-success":
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("4") + `","plane":"code_correctness","role_id":"minimal","findings":[]}`)
		os.Exit(0)
	case "findings-no-status":
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("4") + `","plane":"security_forgery_overclaim","role_id":"findings-default","findings":[{"id":"F1","severity":"major","citation":{"context_ref_id":"diff","diff_hunk_id":"hunk-1"},"summary":"finding"}]}`)
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
	case "unsafe-structured-output":
		fmt.Print(`{"packet_digest":"sha256:` + sixtyFour("9") + `","plane":"privacy_output_safety","role_id":"privacy","runner":"pi","requested_model":"fake-pi","observed_model":"fake-pi","model_family":"fake","model_version":"v1","status":"findings_reported","findings":[{"id":"F1","severity":"minor","citation":{"context_ref_id":"diff","diff_hunk_id":"hunk-1"},"summary":"SYNTHETIC_TOKEN_SECRET_PIPELINE","rationale":"SYNTHETIC_RATIONALE_SECRET_PIPELINE","suggested_fix":"remove SYNTHETIC_PROMPT_SECRET_PIPELINE","question":"is /Users/private/repo visible?","evidence_refs":["https://access_token=secret@example.invalid/review"]}]}`)
		os.Exit(0)
	case "prompt-includes-evidence":
		stdin, err := io.ReadAll(os.Stdin)
		expectedDigest := os.Getenv("PR_REVIEW_EXPECTED_DIGEST")
		prompt := string(stdin)
		if err != nil || expectedDigest == "" || !strings.Contains(prompt, expectedDigest) || !strings.Contains(prompt, "diff --git a/a.go b/a.go") || !strings.Contains(prompt, "PR #123") || !strings.Contains(prompt, "verify: pass") {
			os.Exit(3)
		}
		fmt.Print(`{"packet_digest":"` + expectedDigest + `","plane":"code_correctness","role_id":"code","runner":"manual_external","requested_model":"fake","observed_model":"fake","model_family":"fake","model_version":"v1","status":"no_findings","findings":[]}`)
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func jsonKeys(t *testing.T, value any) map[string]bool {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	keys := map[string]bool{}
	for key := range decoded {
		keys[key] = true
	}
	return keys
}

func assertJSONKeys(t *testing.T, keys map[string]bool, present, absent []string) {
	t.Helper()
	for _, key := range present {
		if !keys[key] {
			t.Fatalf("json keys missing %q in %#v", key, keys)
		}
	}
	for _, key := range absent {
		if keys[key] {
			t.Fatalf("json keys unexpectedly contain %q in %#v", key, keys)
		}
	}
}
