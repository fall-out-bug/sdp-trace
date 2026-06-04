package packet

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidateAndRenderHappyPath(t *testing.T) {
	bundle := validBundle()
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StatePass {
		t.Fatalf("state = %s errors=%v", result.State, result.Errors)
	}
	md, err := RenderMarkdown(bundle)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	for _, want := range []string{
		"# Change Evidence Packet v0",
		"## Packet Metadata",
		"| PC-THEATER | pass | No P0 theater finding triggered.",
		"## What This Packet Does Not Prove",
		"not merge, release, compliance, production trust, or quality approval",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("rendered packet missing %q:\n%s", want, md)
		}
	}
}

func TestLoadBundleAndGitHubInput(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "bundle.json")
	inputPath := filepath.Join(dir, "github-input.json")
	writeJSONForTest(t, bundlePath, validBundle())
	writeJSONForTest(t, inputPath, validGitHubInput())

	bundle, err := LoadBundle(bundlePath)
	if err != nil {
		t.Fatalf("LoadBundle: %v", err)
	}
	if bundle.Packet.PacketID == "" {
		t.Fatalf("loaded bundle missing packet id")
	}
	input, err := LoadGitHubInput(inputPath)
	if err != nil {
		t.Fatalf("LoadGitHubInput: %v", err)
	}
	if input.PR.Number != 5 {
		t.Fatalf("loaded PR number = %d", input.PR.Number)
	}
}

func TestLoadBundleAndGitHubInputRejectMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "bundle.json")
	inputPath := filepath.Join(dir, "github-input.json")
	if err := os.WriteFile(bundlePath, []byte("{"), 0o644); err != nil {
		t.Fatalf("write bundle: %v", err)
	}
	if err := os.WriteFile(inputPath, []byte("{"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if _, err := LoadBundle(bundlePath); err == nil {
		t.Fatalf("LoadBundle accepted malformed JSON")
	}
	if _, err := LoadGitHubInput(inputPath); err == nil {
		t.Fatalf("LoadGitHubInput accepted malformed JSON")
	}
}

func TestPacketContractCatalogsPreserveTrustSurface(t *testing.T) {
	if PacketSchemaVersion != "change-evidence-packet.v0" || BundleSchemaVersion != "evidence-bundle-manifest.v0" {
		t.Fatalf("schema versions changed packet=%q bundle=%q", PacketSchemaVersion, BundleSchemaVersion)
	}
	wantRows := []string{
		"PC-CHANGE",
		"PC-INITIATOR",
		"PC-AGENT-ROUTE",
		"PC-MUTATION",
		"PC-VERIFICATION",
		"PC-REVIEW",
		"PC-AUTHORITY",
		"PC-THEATER",
		"PC-ATTESTATION",
		"PC-DECISION",
		"PC-RESIDUAL-GAPS",
	}
	if strings.Join(RequiredRows, "|") != strings.Join(wantRows, "|") {
		t.Fatalf("required rows = %#v", RequiredRows)
	}
	wantDecisions := []string{"merge", "release", "risk_acceptance", "security_review"}
	if strings.Join(requiredDecisions, "|") != strings.Join(wantDecisions, "|") {
		t.Fatalf("required decisions = %#v", requiredDecisions)
	}
	assertBoolCatalog(t, "states", states, []string{StatePass, StatePartial, StateFail, StateCannotVerify, StateNotAssessed, StateNotInScope})
	assertBoolCatalog(t, "missingReasonStates", missingReasonStates, []string{StatePartial, StateFail, StateCannotVerify, StateNotAssessed, StateNotInScope})
	assertBoolCatalog(t, "packetStates", packetStates, []string{"draft", "review_ready", "reviewed", "superseded"})
	assertBoolCatalog(t, "authoringMethods", authoringMethods, []string{AuthoringToolGenerated, "hand_authored_before_tooling"})
	assertBoolCatalog(t, "retainedForms", retainedForms, []string{"raw", "redacted", "digest_only", "external_ref", "not_retained"})
	assertBoolCatalog(t, "redactionStatuses", redactionStatuses, []string{"not_needed", "redacted", "digest_only", "withheld", StateCannotVerify})
	assertBoolCatalog(t, "theaterReasonCodes", theaterReasonCodes, []string{"agent_claimed_verification", "unbound_intent", "ci_theater", "scope_theater", "prompt_contamination"})
}

func TestPacketCoreTypesPreserveJSONShape(t *testing.T) {
	packet := Packet{
		PacketVersion:   PacketSchemaVersion,
		PacketID:        "packet-1",
		SourceChange:    SourceChange{Repository: "repo", ChangeID: "PR-1"},
		GeneratedAt:     "2026-06-04T00:00:00Z",
		AuthoringMethod: AuthoringToolGenerated,
		SelectedProfile: "default",
		RedactionPolicy: "standard",
		BundleRef:       "bundle-1",
		PacketState:     "review_ready",
		Projection:      Projection{Kind: ProjectionCanonical, Canonical: true},
		Rows:            []Row{{ID: "PC-CHANGE", State: StatePass, Summary: "ok", EvidenceRefs: []string{"git:change"}, Owner: "maintainer"}},
		NonApproval:     "not approval",
	}
	assertJSONHasKeys(t, packet, []string{"packet_version", "packet_id", "source_change", "generated_at", "authoring_method", "selected_profile", "redaction_policy", "bundle_ref", "packet_state", "projection", "rows", "non_approval"}, []string{"theater_findings", "residual_gaps", "decision_owners", "extensions"})
	assertJSONHasKeys(t, SourceChange{Repository: "repo", ChangeID: "PR-1"}, []string{"repository", "change_id"}, []string{"url", "base_ref", "head_ref", "commit_range", "head_sha"})
	assertJSONHasKeys(t, Projection{Kind: ProjectionCanonical, Canonical: true}, []string{"kind", "canonical"}, []string{"artifact_ref"})
	assertJSONHasKeys(t, Row{ID: "PC-CHANGE", State: StatePass, Summary: "ok", EvidenceRefs: []string{"git:change"}, Owner: "maintainer"}, []string{"id", "state", "summary", "evidence_refs", "owner"}, []string{"reason"})
	assertJSONHasKeys(t, TheaterFinding{ReasonCode: "ci_theater", State: StateFail, Finding: "overclaim", TriggerEvidenceRefs: []string{"ci:run"}}, []string{"reason_code", "state", "finding", "trigger_evidence_refs"}, []string{"severity", "required_closure_evidence"})
	assertJSONHasKeys(t, ResidualGap{RowID: "PC-REVIEW", State: StateNotAssessed, Reason: "missing"}, []string{"row_id", "state", "reason"}, []string{"evidence_refs", "closure_evidence"})
	assertJSONHasKeys(t, DecisionOwner{Decision: "merge", Owner: "cto", State: StateNotAssessed}, []string{"decision", "owner", "state"}, []string{"reason"})
	assertJSONTags(t, Packet{}, requiredJSONFields("packet_version", "packet_id", "source_change", "generated_at", "authoring_method", "selected_profile", "redaction_policy", "bundle_ref", "packet_state", "projection", "rows", "non_approval"), optionalJSONFields("theater_findings", "residual_gaps", "decision_owners", "extensions"))
	assertJSONTags(t, SourceChange{}, requiredJSONFields(), optionalJSONFields("repository", "change_id", "url", "base_ref", "head_ref", "commit_range", "head_sha"))
	assertJSONTags(t, Projection{}, requiredJSONFields("kind", "canonical"), optionalJSONFields("artifact_ref"))
	assertJSONTags(t, Row{}, requiredJSONFields("id", "state", "summary", "evidence_refs", "owner"), optionalJSONFields("reason"))
	assertJSONTags(t, TheaterFinding{}, requiredJSONFields("reason_code", "state", "finding", "trigger_evidence_refs"), optionalJSONFields("severity", "required_closure_evidence"))
	assertJSONTags(t, ResidualGap{}, requiredJSONFields("row_id", "state", "reason"), optionalJSONFields("evidence_refs", "closure_evidence"))
	assertJSONTags(t, DecisionOwner{}, requiredJSONFields("decision", "owner", "state"), optionalJSONFields("reason"))
}

func TestPacketBundleTypesPreserveJSONShape(t *testing.T) {
	manifest := BundleManifest{
		SchemaVersion: BundleSchemaVersion,
		BundleID:      "bundle-1",
		Entries:       []BundleEntry{{Ref: "git:change", SourceClass: "git", RetainedForm: "raw", RedactionStatus: "not_needed"}},
	}
	assertJSONHasKeys(t, manifest, []string{"schema_version", "bundle_id", "entries"}, []string{"packet_digest", "resolvers"})
	entry := BundleEntry{
		Ref:             "git:change",
		SourceClass:     "git",
		RetainedForm:    "raw",
		RedactionStatus: "not_needed",
	}
	assertJSONHasKeys(t, entry, []string{"ref", "source_class", "retained_form", "redaction_status"}, []string{"digest", "resolver", "expires_at", "artifact_access", "projection_role", "evidence_kind", "observed_components", "contradicts_ref", "contradicts_row_id", "actor", "write_authority", "generated_by", "source_commit_state", "source_ref"})
	assertJSONHasKeys(t, ResolverEntry{Ref: "git:change", Resolver: "artifact"}, []string{"ref", "resolver"}, nil)
	assertJSONHasKeys(t, Bundle{Packet: Packet{PacketID: "packet-1"}, Manifest: manifest}, []string{"packet", "manifest"}, nil)
	assertJSONTags(t, BundleManifest{}, requiredJSONFields("schema_version", "bundle_id", "entries"), optionalJSONFields("packet_digest", "resolvers"))
	assertJSONTags(t, BundleEntry{}, requiredJSONFields("ref", "source_class", "retained_form", "redaction_status"), optionalJSONFields("digest", "resolver", "expires_at", "artifact_access", "projection_role", "evidence_kind", "observed_components", "contradicts_ref", "contradicts_row_id", "actor", "write_authority", "generated_by", "source_commit_state", "source_ref"))
	assertJSONTags(t, ResolverEntry{}, requiredJSONFields("ref", "resolver"), optionalJSONFields())
	assertJSONTags(t, Bundle{}, requiredJSONFields("packet", "manifest"), optionalJSONFields())
}

func TestGitHubEvidenceInputTypesPreserveJSONShape(t *testing.T) {
	input := GitHubPREvidenceInput{
		SchemaVersion: "github-pr-evidence.v0",
		PR:            GitHubPR{Number: 5, URL: "https://example.invalid/pr/5", Title: "title", Author: "dev", BaseRef: "main", HeadRef: "feature", HeadSHA: strings.Repeat("a", 40)},
		CommitRange:   GitHubCommitRange{Base: strings.Repeat("b", 40), Head: strings.Repeat("a", 40)},
	}
	assertJSONHasKeys(t, input, []string{"schema_version", "pr", "commit_range", "prompt_boundary"}, []string{"checks", "artifacts", "reviews", "workflow_run_id", "require_prompt_boundary", "agent_route_refs", "agent_route_components", "agent_route_digest", "agent_route_evidence_kind", "integration_actions"})
	assertJSONHasKeys(t, input.PR, []string{"number", "url", "title", "author", "base_ref", "head_ref", "head_sha"}, []string{"body_ref"})
	assertJSONHasKeys(t, input.CommitRange, []string{"base", "head"}, []string{"changed_files_ref"})
	assertJSONHasKeys(t, GitHubCheck{Name: "verify", URL: "https://example.invalid/check", Conclusion: StatePass}, []string{"name", "url", "conclusion"}, []string{"artifact_refs"})
	assertJSONHasKeys(t, GitHubArtifact{Name: "report", Resolver: "artifact", RetainedForm: "raw"}, []string{"name", "resolver", "retained_form"}, []string{"expires_at", "digest"})
	assertJSONHasKeys(t, GitHubReview{Reviewer: "reviewer", Resolver: "review", State: StateNotAssessed}, []string{"reviewer", "resolver", "state"}, nil)
	assertJSONHasKeys(t, PromptBoundary{Text: "prompt"}, []string{"text"}, []string{"digest", "capture_actor", "captured_at", "capture_method"})
	assertJSONHasKeys(t, PromptBoundaryClassification{Verdict: StateCannotVerify, RouteProofEffect: "cannot_support_route"}, []string{"verdict", "route_proof_effect"}, []string{"reasons"})
	assertJSONHasKeys(t, IntegrationAction{Kind: "merge", Actor: "bot", Resolver: "action"}, []string{"kind", "actor", "resolver"}, nil)
	assertJSONHasKeys(t, BuildPRResult{State: StatePass}, []string{"state"}, []string{"bundle_path", "packet_path", "result_path", "errors"})
	assertJSONHasKeys(t, Validation{State: StatePass}, []string{"state"}, []string{"errors"})
	assertJSONTags(t, GitHubPREvidenceInput{}, requiredJSONFields("schema_version", "pr", "commit_range"), optionalJSONFields("checks", "artifacts", "reviews", "workflow_run_id", "require_prompt_boundary", "agent_route_refs", "agent_route_components", "agent_route_digest", "agent_route_evidence_kind", "prompt_boundary", "integration_actions"))
	assertJSONTags(t, GitHubPR{}, requiredJSONFields("number", "url", "title", "author", "base_ref", "head_ref", "head_sha"), optionalJSONFields("body_ref"))
	assertJSONTags(t, GitHubCommitRange{}, requiredJSONFields("base", "head"), optionalJSONFields("changed_files_ref"))
	assertJSONTags(t, GitHubCheck{}, requiredJSONFields("name", "url", "conclusion"), optionalJSONFields("artifact_refs"))
	assertJSONTags(t, GitHubArtifact{}, requiredJSONFields("name", "resolver", "retained_form"), optionalJSONFields("expires_at", "digest"))
	assertJSONTags(t, GitHubReview{}, requiredJSONFields("reviewer", "resolver", "state"), optionalJSONFields())
	assertJSONTags(t, PromptBoundary{}, requiredJSONFields(), optionalJSONFields("text", "digest", "capture_actor", "captured_at", "capture_method"))
	assertJSONTags(t, PromptBoundaryClassification{}, requiredJSONFields("verdict", "route_proof_effect"), optionalJSONFields("reasons"))
	assertJSONTags(t, IntegrationAction{}, requiredJSONFields("kind", "actor", "resolver"), optionalJSONFields())
	assertJSONTags(t, BuildPRResult{}, requiredJSONFields("state"), optionalJSONFields("bundle_path", "packet_path", "result_path", "errors"))
	assertJSONTags(t, Validation{}, requiredJSONFields("state"), optionalJSONFields("errors"))
}

func TestValidateRejectsMissingVerificationPass(t *testing.T) {
	bundle := validBundle()
	setRow(&bundle, "PC-VERIFICATION", Row{
		ID:           "PC-VERIFICATION",
		State:        StatePass,
		Summary:      "Agent said tests passed.",
		EvidenceRefs: nil,
		Owner:        "maintainer",
	})
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "PC-VERIFICATION pass requires retained evidence refs") {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateRejectsExpiredArtifactPass(t *testing.T) {
	bundle := validBundle()
	bundle.Manifest.Entries[0].ExpiresAt = "2026-05-10T12:00:00Z"
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "pass cites expired artifact ref") {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateRejectsUnverifiableArtifactPass(t *testing.T) {
	for _, access := range []string{"inaccessible", "malformed", "cannot_verify", "not_assessed"} {
		t.Run(access, func(t *testing.T) {
			bundle := validBundle()
			for i := range bundle.Manifest.Entries {
				if bundle.Manifest.Entries[i].Ref == "ci:run" {
					bundle.Manifest.Entries[i].ArtifactAccess = access
				}
			}
			result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StateFail || !hasError(result.Errors, "PC-VERIFICATION pass cites unverifiable artifact ref") {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestValidateRejectsBundleMismatch(t *testing.T) {
	bundle := validBundle()
	bundle.Manifest.BundleID = "other-bundle"
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "packet.bundle_ref") {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateRejectsMissingResidualGapForNonPassRow(t *testing.T) {
	bundle := validBundle()
	bundle.Packet.ResidualGaps = nil
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "PC-INITIATOR non-pass row requires residual gap explanation") {
		t.Fatalf("result = %+v", result)
	}
}

func TestGapForRowWithClosure(t *testing.T) {
	gaps := []ResidualGap{
		{RowID: "PC-REVIEW", Reason: "missing review"},
		{RowID: "PC-AGENT-ROUTE", Reason: "partial route", ClosureEvidence: "retained route observation"},
	}
	if !gapForRowWithClosure(gaps, "PC-AGENT-ROUTE") {
		t.Fatalf("expected closure gap for PC-AGENT-ROUTE")
	}
	if gapForRowWithClosure(gaps, "PC-REVIEW") {
		t.Fatalf("gap without closure evidence should not close")
	}
	if gapForRowWithClosure(gaps, "PC-VERIFICATION") {
		t.Fatalf("missing row should not close")
	}
}

func TestValidateRejectsTheaterPassWithFindings(t *testing.T) {
	bundle := validBundle()
	bundle.Packet.TheaterFindings = []TheaterFinding{{
		ReasonCode:          "agent_claimed_verification",
		State:               StatePartial,
		Finding:             "Agent claimed tests without independent artifact.",
		TriggerEvidenceRefs: []string{"harness:route"},
	}}
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "PC-THEATER cannot be pass") {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateRowDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name string
		edit func(*Bundle)
		want []string
	}{
		{
			name: "missing required row",
			edit: func(bundle *Bundle) {
				removeRow(bundle, "PC-REVIEW")
			},
			want: []string{"missing required row \"PC-REVIEW\""},
		},
		{
			name: "unknown row id",
			edit: func(bundle *Bundle) {
				bundle.Packet.Rows = append(bundle.Packet.Rows, Row{ID: "PC-UNKNOWN", State: StatePass, Summary: "unknown", EvidenceRefs: []string{"ci:run"}, Owner: "maintainer"})
			},
			want: []string{"unknown row id \"PC-UNKNOWN\""},
		},
		{
			name: "duplicate row id",
			edit: func(bundle *Bundle) {
				bundle.Packet.Rows = append(bundle.Packet.Rows, row("PC-REVIEW", StatePartial, "duplicate review", []string{"review:packet"}))
			},
			want: []string{"duplicate row id \"PC-REVIEW\""},
		},
		{
			name: "row fields and reason",
			edit: func(bundle *Bundle) {
				setRow(bundle, "PC-REVIEW", Row{ID: "PC-REVIEW", State: "mystery", Summary: "", EvidenceRefs: []string{"review:packet"}, Owner: ""})
			},
			want: []string{
				"PC-REVIEW has unknown state \"mystery\"",
				"PC-REVIEW requires summary",
				"PC-REVIEW requires owner",
			},
		},
		{
			name: "missing non-pass reason",
			edit: func(bundle *Bundle) {
				setRow(bundle, "PC-REVIEW", Row{ID: "PC-REVIEW", State: StatePartial, Summary: "partial review", EvidenceRefs: []string{"review:packet"}, Owner: "maintainer"})
			},
			want: []string{"PC-REVIEW state partial requires reason"},
		},
		{
			name: "pass row evidence refs",
			edit: func(bundle *Bundle) {
				setRow(bundle, "PC-VERIFICATION", Row{ID: "PC-VERIFICATION", State: StatePass, Summary: "tests passed", Owner: "maintainer"})
			},
			want: []string{"PC-VERIFICATION pass requires retained evidence refs"},
		},
		{
			name: "absent evidence ref",
			edit: func(bundle *Bundle) {
				setRow(bundle, "PC-VERIFICATION", Row{ID: "PC-VERIFICATION", State: StatePass, Summary: "tests passed", EvidenceRefs: []string{"ci:missing"}, Owner: "maintainer"})
			},
			want: []string{"PC-VERIFICATION evidence ref \"ci:missing\" is absent from manifest"},
		},
		{
			name: "no resolver",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.Resolver = ""
				})
			},
			want: []string{"PC-VERIFICATION evidence ref \"ci:run\" has no resolver entry"},
		},
		{
			name: "expired pass evidence",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.ExpiresAt = "2026-05-10T12:00:00Z"
				})
			},
			want: []string{"PC-VERIFICATION pass cites expired artifact ref \"ci:run\""},
		},
		{
			name: "unverifiable pass evidence",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.RetainedForm = "not_retained"
				})
			},
			want: []string{"PC-VERIFICATION pass cites unverifiable artifact ref \"ci:run\""},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validBundle()
			tt.edit(&bundle)
			refreshPacketDigest(&bundle)
			result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StateFail {
				t.Fatalf("state = %s errors=%v, want fail", result.State, result.Errors)
			}
			for _, want := range tt.want {
				if !hasError(result.Errors, want) {
					t.Fatalf("errors = %v, want %q", result.Errors, want)
				}
			}
		})
	}
}

func TestRenderCleanTheaterUsesRowState(t *testing.T) {
	bundle := validBundle()
	setRowState(&bundle, "PC-THEATER", StateNotAssessed, "theater assessment not run")
	bundle.Packet.ResidualGaps = append(bundle.Packet.ResidualGaps, ResidualGap{RowID: "PC-THEATER", State: StateNotAssessed, Reason: "theater assessment not run", ClosureEvidence: "run theater assessment"})
	refreshPacketDigest(&bundle)
	md, err := RenderMarkdown(bundle)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if strings.Contains(md, "| none | pass | none | No P0 theater finding triggered.") ||
		!strings.Contains(md, "| none | not_assessed | none | No P0 theater finding triggered.") {
		t.Fatalf("rendered theater row overclaimed:\n%s", md)
	}
}

func TestPacketRenderingHelpersPreserveTables(t *testing.T) {
	var out bytes.Buffer
	renderTheaterFinding(&out, TheaterFinding{
		ReasonCode:              "prompt_contamination",
		State:                   StateFail,
		Severity:                "P0",
		Finding:                 "a|b\nc",
		TriggerEvidenceRefs:     []string{"prompt:boundary", "review:note"},
		RequiredClosureEvidence: "fresh|review",
	})
	if got, want := out.String(), "| prompt_contamination | fail | P0 | a\\|b c | prompt:boundary, review:note | fresh\\|review |\n"; got != want {
		t.Fatalf("theater finding render = %q want %q", got, want)
	}

	out.Reset()
	renderDecisions(&out, []DecisionOwner{{Decision: "merge|release", Owner: "owner\nname", State: StateNotAssessed}})
	if got := out.String(); !strings.Contains(got, "## Decision Ownership\n\n| decision | owner | state | reason |\n| --- | --- | --- | --- |\n") ||
		!strings.Contains(got, "| merge\\|release | owner name | not_assessed | none |") {
		t.Fatalf("decision render drifted:\n%s", got)
	}

	out.Reset()
	renderEvidence(&out, BundleManifest{
		BundleID:  "bundle|one",
		Entries:   []BundleEntry{{Ref: "artifact:one", SourceClass: "ci", RetainedForm: "raw", RedactionStatus: "not_needed"}},
		Resolvers: []ResolverEntry{{Ref: "artifact:one", Resolver: "resolver|fallback"}},
	})
	if got := out.String(); !strings.Contains(got, "Manifest: `bundle\\|one`") ||
		!strings.Contains(got, "| ref | source class | retained form | redaction status | resolver |") ||
		!strings.Contains(got, "| artifact:one | ci | raw | not_needed | resolver\\|fallback |") {
		t.Fatalf("evidence render drifted:\n%s", got)
	}
}

func TestPacketResidualAndNonProofRendering(t *testing.T) {
	var out bytes.Buffer
	renderResidualGaps(&out, nil)
	if got, want := out.String(), "## Residual Gaps\n\nNo residual gaps recorded beyond row states.\n\n"; got != want {
		t.Fatalf("empty residual render = %q want %q", got, want)
	}

	out.Reset()
	renderResidualGaps(&out, []ResidualGap{{RowID: "PC-REVIEW", State: StatePartial, Reason: "needs|review", ClosureEvidence: "review\nagain"}})
	if got := out.String(); !strings.Contains(got, "| row id | state | reason | closure evidence |") ||
		!strings.Contains(got, "| PC-REVIEW | partial | needs\\|review | review again |") {
		t.Fatalf("residual render drifted:\n%s", got)
	}

	out.Reset()
	renderNonProof(&out, Packet{})
	if got := out.String(); !strings.Contains(got, "## What This Packet Does Not Prove") ||
		!strings.Contains(got, "does not approve merge, release, compliance, production trust") {
		t.Fatalf("non-proof fallback drifted:\n%s", got)
	}
}

func TestPacketRenderLookupAndDigestHelpers(t *testing.T) {
	if got := rowByID(nil, "PC-MISSING"); got.ID != "PC-MISSING" || got.State != StateCannotVerify || got.Summary != "row missing" || got.Reason != "row missing" {
		t.Fatalf("missing row fallback = %+v", got)
	}
	if got := requiredRowIndex("PC-CHANGE"); got != 0 {
		t.Fatalf("PC-CHANGE index = %d", got)
	}
	if got := requiredRowIndex("EXTENSION"); got != len(RequiredRows) {
		t.Fatalf("extension index = %d want %d", got, len(RequiredRows))
	}
	if got := resolverFromList([]ResolverEntry{{Ref: "a", Resolver: "resolver-a"}}, "missing"); got != "" {
		t.Fatalf("missing resolver = %q", got)
	}
	if got := md("a|b\nc"); got != "a\\|b c" {
		t.Fatalf("md escape = %q", got)
	}
	if got := md(" \t "); got != "none" {
		t.Fatalf("blank md = %q", got)
	}
	packet := validBundle().Packet
	first := PacketDigest(packet)
	second := PacketDigest(packet)
	if !strings.HasPrefix(first, "sha256:") || len(first) != len("sha256:")+64 || first != second {
		t.Fatalf("packet digest not deterministic sha256: first=%q second=%q", first, second)
	}
}

func TestRenderMarkdownPreservesTopLevelOrderAndHeaders(t *testing.T) {
	rendered, err := RenderMarkdown(validBundle())
	if err != nil {
		t.Fatalf("RenderMarkdown: %v", err)
	}
	wantOrder := []string{
		"# Change Evidence Packet v0",
		"## Executive Summary",
		"## Packet Metadata",
		"## Required Rows",
		"## Theater Findings",
		"## Decision Ownership",
		"## Evidence Bundle",
		"## Residual Gaps",
		"## What This Packet Does Not Prove",
	}
	last := -1
	for _, marker := range wantOrder {
		pos := strings.Index(rendered, marker)
		if pos < 0 {
			t.Fatalf("rendered markdown missing %q:\n%s", marker, rendered)
		}
		if pos <= last {
			t.Fatalf("rendered markdown order drifted at %q:\n%s", marker, rendered)
		}
		last = pos
	}
	for _, header := range []string{
		"| row id | state | answer | evidence refs | gap / next evidence | owner |",
		"| reason code | state | severity | finding | trigger evidence | required closure evidence |",
		"| row id | state | reason | closure evidence |",
		"| decision | owner | state | reason |",
		"| ref | source class | retained form | redaction status | resolver |",
	} {
		if !strings.Contains(rendered, header) {
			t.Fatalf("rendered markdown missing table header %q:\n%s", header, rendered)
		}
	}
}

func TestValidateContradictionRequiresPartialAndGap(t *testing.T) {
	bundle := validBundle()
	bundle.Manifest.Entries = append(bundle.Manifest.Entries, BundleEntry{
		Ref:              "review:contradiction",
		SourceClass:      "review",
		RetainedForm:     "raw",
		RedactionStatus:  "not_needed",
		Resolver:         "examples/change-evidence-packet/contradiction.md",
		ContradictsRef:   "ci:run",
		ContradictsRowID: "PC-VERIFICATION",
	})
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "PC-VERIFICATION has contradictory evidence") {
		t.Fatalf("result = %+v", result)
	}
	setRowState(&bundle, "PC-VERIFICATION", StatePartial, "review evidence contradicts CI claim")
	bundle.Packet.ResidualGaps = append(bundle.Packet.ResidualGaps, ResidualGap{
		RowID:           "PC-VERIFICATION",
		State:           StatePartial,
		Reason:          "review and CI evidence disagree",
		ClosureEvidence: "fresh rerun and reviewer acknowledgement",
	})
	refreshPacketDigest(&bundle)
	result = Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StatePass {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateContradictionTargetSelection(t *testing.T) {
	bundle := validBundle()
	bundle.Manifest.Entries = append(bundle.Manifest.Entries, BundleEntry{
		Ref:              "review:contradiction",
		SourceClass:      "review",
		RetainedForm:     "raw",
		RedactionStatus:  "not_needed",
		Resolver:         "examples/change-evidence-packet/contradiction.md",
		ContradictsRef:   "git:change",
		ContradictsRowID: "PC-VERIFICATION",
	})
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if !hasError(result.Errors, "PC-VERIFICATION has contradictory evidence") ||
		hasError(result.Errors, "PC-CHANGE has contradictory evidence") {
		t.Fatalf("explicit row target did not override ref fallback: %+v", result)
	}
}

func TestRowIDForRefUsesRequiredRowOrder(t *testing.T) {
	bundle := validBundle()
	rows := rowsByID(bundle.Packet.Rows)
	if got := rowIDForRef(rows, "git:change"); got != "PC-CHANGE" {
		t.Fatalf("rowIDForRef(shared ref) = %q, want PC-CHANGE", got)
	}
}

func TestRowIDForRefUsesExtensionFallbackOrderAndExactRefs(t *testing.T) {
	rows := map[string]Row{
		"PC-VERIFICATION": {ID: "PC-VERIFICATION", EvidenceRefs: []string{"custom:ref"}},
		"ZZ-EXT":          {ID: "ZZ-EXT", EvidenceRefs: []string{"custom:ref", "custom:prefix-extra"}},
		"AA-EXT":          {ID: "AA-EXT", EvidenceRefs: []string{"custom:ref"}},
	}
	if got := rowIDForRef(rows, "custom:ref"); got != "PC-VERIFICATION" {
		t.Fatalf("required row precedence = %q, want PC-VERIFICATION", got)
	}

	delete(rows, "PC-VERIFICATION")
	if got := rowIDForRef(rows, "custom:ref"); got != "AA-EXT" {
		t.Fatalf("extension fallback row = %q, want AA-EXT", got)
	}
	if got := rowIDForRef(rows, "custom:prefix"); got != "" {
		t.Fatalf("rowIDForRef used prefix match = %q", got)
	}
}

func TestGapForRowRequiresReasonForResidualCoverage(t *testing.T) {
	bundle := validBundle()
	setRowState(&bundle, "PC-REVIEW", StatePartial, "review evidence is pending")
	for i := range bundle.Packet.ResidualGaps {
		if bundle.Packet.ResidualGaps[i].RowID == "PC-REVIEW" {
			bundle.Packet.ResidualGaps[i].Reason = ""
		}
	}
	refreshPacketDigest(&bundle)
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "PC-REVIEW non-pass row requires residual gap explanation") {
		t.Fatalf("blank residual reason should not close coverage: %+v", result)
	}

	for i := range bundle.Packet.ResidualGaps {
		if bundle.Packet.ResidualGaps[i].RowID == "PC-REVIEW" {
			bundle.Packet.ResidualGaps[i].Reason = "review evidence is retained externally"
		}
	}
	refreshPacketDigest(&bundle)
	result = Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if hasError(result.Errors, "PC-REVIEW non-pass row requires residual gap explanation") {
		t.Fatalf("non-empty residual reason should close coverage: %+v", result)
	}
}

func TestValidateRejectsProjectionMarkedCanonicalOverArtifact(t *testing.T) {
	bundle := validBundle()
	bundle.Packet.Projection = Projection{Kind: "github_pr_comment", Canonical: true}
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "canonical projection must be") {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateMetadataAndManifestDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name string
		edit func(*Bundle)
		want string
	}{
		{
			name: "packet schema",
			edit: func(bundle *Bundle) {
				bundle.Packet.PacketVersion = "other-packet-schema"
				refreshPacketDigest(bundle)
			},
			want: "packet.packet_version must be",
		},
		{
			name: "manifest schema",
			edit: func(bundle *Bundle) {
				bundle.Manifest.SchemaVersion = "other-manifest-schema"
			},
			want: "manifest.schema_version must be",
		},
		{
			name: "packet id",
			edit: func(bundle *Bundle) {
				bundle.Packet.PacketID = ""
				refreshPacketDigest(bundle)
			},
			want: "packet.packet_id is required",
		},
		{
			name: "bundle ref",
			edit: func(bundle *Bundle) {
				bundle.Packet.BundleRef = ""
				refreshPacketDigest(bundle)
			},
			want: "packet.bundle_ref is required",
		},
		{
			name: "manifest bundle id",
			edit: func(bundle *Bundle) {
				bundle.Manifest.BundleID = ""
			},
			want: "manifest.bundle_id is required",
		},
		{
			name: "missing packet digest",
			edit: func(bundle *Bundle) {
				bundle.Manifest.PacketDigest = ""
			},
			want: "manifest.packet_digest is required",
		},
		{
			name: "mismatched packet digest",
			edit: func(bundle *Bundle) {
				bundle.Manifest.PacketDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			},
			want: "manifest.packet_digest does not match packet content",
		},
		{
			name: "missing non approval",
			edit: func(bundle *Bundle) {
				bundle.Packet.NonApproval = ""
				refreshPacketDigest(bundle)
			},
			want: "packet.non_approval is required",
		},
		{
			name: "packet state",
			edit: func(bundle *Bundle) {
				bundle.Packet.PacketState = "approved"
				refreshPacketDigest(bundle)
			},
			want: "packet.packet_state has unknown value",
		},
		{
			name: "authoring method",
			edit: func(bundle *Bundle) {
				bundle.Packet.AuthoringMethod = "self_attested"
				refreshPacketDigest(bundle)
			},
			want: "packet.authoring_method has unknown value",
		},
		{
			name: "non canonical artifact",
			edit: func(bundle *Bundle) {
				bundle.Packet.Projection = Projection{Kind: "github_pr_comment", Canonical: false}
				refreshPacketDigest(bundle)
			},
			want: "non-canonical packet projection requires artifact_ref",
		},
		{
			name: "empty manifest ref",
			edit: func(bundle *Bundle) {
				bundle.Manifest.Entries = append(bundle.Manifest.Entries, entry("", "ci"))
			},
			want: "manifest entry has empty ref",
		},
		{
			name: "retained form",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.RetainedForm = "imaginary"
				})
			},
			want: "unknown retained_form",
		},
		{
			name: "redaction status",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.RedactionStatus = "imaginary"
				})
			},
			want: "unknown redaction_status",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validBundle()
			tt.edit(&bundle)
			result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StateFail || !hasError(result.Errors, tt.want) {
				t.Fatalf("result = %+v, want error containing %q", result, tt.want)
			}
		})
	}
}

func TestValidatePreservesPhaseOrderAndAccumulation(t *testing.T) {
	bundle := validBundle()
	bundle.Packet.PacketVersion = "other-packet-schema"
	bundle.Manifest.Entries = append(bundle.Manifest.Entries, entry("", "ci"))
	setRow(&bundle, "PC-VERIFICATION", Row{
		ID:           "PC-VERIFICATION",
		State:        StatePass,
		Summary:      "Agent said tests passed.",
		EvidenceRefs: []string{"missing:ci"},
		Owner:        "maintainer",
	})
	bundle.Packet.TheaterFindings = []TheaterFinding{{
		ReasonCode:          "agent_claimed_verification",
		State:               StatePartial,
		Finding:             "Agent claimed verification without retained evidence.",
		TriggerEvidenceRefs: []string{"missing:finding"},
	}}
	refreshPacketDigest(&bundle)
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	assertErrorsInOrder(t, result.Errors, []string{
		"packet.packet_version must be",
		"manifest entry has empty ref",
		"PC-VERIFICATION evidence ref \"missing:ci\" is absent from manifest",
		"PC-THEATER cannot be pass when theater findings are present",
		"theater finding agent_claimed_verification evidence ref \"missing:finding\" is absent from manifest",
	})
}

func TestValidateAcceptsResolverEntryForManifestEvidence(t *testing.T) {
	bundle := validBundle()
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == "ci:run" {
			bundle.Manifest.Entries[i].Resolver = ""
		}
	}
	bundle.Manifest.Resolvers = []ResolverEntry{
		{Ref: "", Resolver: "ignored"},
		{Ref: "ci:run", Resolver: "examples/change-evidence-packet/ci-run"},
	}
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StatePass {
		t.Fatalf("resolver entry did not satisfy evidence resolver: %+v", result)
	}
}

func TestValidateManifestResolverIndexingFeedsRowEvidenceRefs(t *testing.T) {
	for _, tt := range []struct {
		name      string
		edit      func(*Bundle)
		wantState string
		wantError string
	}{
		{
			name: "manifest resolver fallback",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.Resolver = "manifest-resolver"
				})
				bundle.Manifest.Resolvers = nil
			},
			wantState: StatePass,
		},
		{
			name: "resolver entry override",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.Resolver = "manifest-resolver"
				})
				bundle.Manifest.Resolvers = []ResolverEntry{{Ref: "ci:run", Resolver: " "}}
			},
			wantState: StateFail,
			wantError: "PC-VERIFICATION evidence ref \"ci:run\" has no resolver entry",
		},
		{
			name: "empty resolver ref ignored",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "ci:run", func(entry *BundleEntry) {
					entry.Resolver = ""
				})
				bundle.Manifest.Resolvers = []ResolverEntry{{Ref: "", Resolver: "ignored"}}
			},
			wantState: StateFail,
			wantError: "PC-VERIFICATION evidence ref \"ci:run\" has no resolver entry",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validBundle()
			tt.edit(&bundle)
			result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != tt.wantState {
				t.Fatalf("state = %s errors=%v, want %s", result.State, result.Errors, tt.wantState)
			}
			if tt.wantError != "" && !hasError(result.Errors, tt.wantError) {
				t.Fatalf("errors = %v, want %q", result.Errors, tt.wantError)
			}
		})
	}
}

func TestBuildGitHubDowngradesIncompleteCommitRange(t *testing.T) {
	input := validGitHubInput()
	input.CommitRange.Base = ""
	bundle := BuildFromGitHubInput(input, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if rowByID(bundle.Packet.Rows, "PC-CHANGE").State != StateCannotVerify ||
		rowByID(bundle.Packet.Rows, "PC-MUTATION").State != StateCannotVerify {
		t.Fatalf("incomplete commit range was not downgraded: %+v", bundle.Packet.Rows)
	}
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StatePass {
		t.Fatalf("generated downgraded bundle should validate: %+v", result)
	}
}

func TestBuildGitHubVerificationPassCitesArtifactRefs(t *testing.T) {
	input := validGitHubInput()
	bundle := BuildFromGitHubInput(input, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	row := rowByID(bundle.Packet.Rows, "PC-VERIFICATION")
	if row.State != StatePass || !containsString(row.EvidenceRefs, "artifact:test-report") {
		t.Fatalf("verification row = %+v", row)
	}
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == "artifact:test-report" {
			bundle.Manifest.Entries[i].ArtifactAccess = "expired"
		}
	}
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "PC-VERIFICATION pass cites unverifiable artifact ref") {
		t.Fatalf("result = %+v", result)
	}
}

func TestCheckDemoRejectsSelfDeclaredRouteEvidence(t *testing.T) {
	input := validGitHubInput()
	input.AgentRouteEvidenceKind = "harness_route_observation"
	input.AgentRouteComponents = []string{"opencode", "gsd-redux", "minimax-m2.5"}
	bundle := BuildFromGitHubInput(input, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "retained structured OpenCode/GSD/MiniMax") {
		t.Fatalf("self-declared route evidence passed: %+v", result)
	}

	input.AgentRouteDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	bundle = BuildFromGitHubInput(input, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	result = CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StatePass {
		t.Fatalf("digest-bound route evidence failed: %+v", result)
	}
}

func TestCheckDemoAcceptsLegacyAndReduxGSDRouteEvidence(t *testing.T) {
	for _, gsdComponent := range []string{"gsd", "gsd-redux"} {
		t.Run(gsdComponent, func(t *testing.T) {
			input := validGitHubInput()
			input.AgentRouteEvidenceKind = "harness_route_observation"
			input.AgentRouteDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			input.AgentRouteComponents = []string{"opencode", gsdComponent, "minimax-m2.5"}
			bundle := BuildFromGitHubInput(input, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StatePass {
				t.Fatalf("route evidence failed: %+v", result)
			}
		})
	}
}

func TestCheckDemoAccumulatesBaseValidationAndDemoErrors(t *testing.T) {
	bundle := demoGateBundle()
	bundle.Manifest.BundleID = "other-bundle"
	bundle.Packet.AuthoringMethod = "hand_authored_before_tooling"
	result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail ||
		!hasError(result.Errors, "packet.bundle_ref") ||
		!hasError(result.Errors, "requires tool_generated authoring_method") {
		t.Fatalf("result did not accumulate base and demo errors: %+v", result)
	}
}

func TestCheckDemoRequiresMinimumPassOrPartialRows(t *testing.T) {
	bundle := demoGateBundle()
	setRowState(&bundle, "PC-CHANGE", StateNotAssessed, "not assessed in fixture")
	setRowState(&bundle, "PC-INITIATOR", StateNotAssessed, "not assessed in fixture")
	setRowState(&bundle, "PC-AGENT-ROUTE", StateNotAssessed, "not assessed in fixture")
	setRowState(&bundle, "PC-MUTATION", StateNotAssessed, "not assessed in fixture")
	setRowState(&bundle, "PC-VERIFICATION", StateNotAssessed, "not assessed in fixture")
	setRowState(&bundle, "PC-DECISION", StateNotAssessed, "not assessed in fixture")
	setRowState(&bundle, "PC-RESIDUAL-GAPS", StateNotAssessed, "not assessed in fixture")
	refreshPacketDigest(&bundle)
	result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "requires at least 4 pass or partial rows") {
		t.Fatalf("result = %+v", result)
	}
}

func TestCheckDemoRowEvidenceMustBeUsable(t *testing.T) {
	for _, tt := range []struct {
		name string
		edit func(*Bundle)
		want string
	}{
		{
			name: "missing PC-CHANGE refs",
			edit: func(bundle *Bundle) {
				setRow(bundle, "PC-CHANGE", row("PC-CHANGE", StatePass, "change evidence missing", nil))
			},
			want: "requires PC-CHANGE retained evidence refs",
		},
		{
			name: "expired PC-MUTATION ref",
			edit: func(bundle *Bundle) {
				for i := range bundle.Manifest.Entries {
					if bundle.Manifest.Entries[i].Ref == "git:change" {
						bundle.Manifest.Entries[i].ExpiresAt = "2026-05-10T12:00:00Z"
					}
				}
			},
			want: "requires PC-MUTATION evidence ref \"git:change\" to be retained and usable",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bundle := demoGateBundle()
			tt.edit(&bundle)
			result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StateFail || !hasError(result.Errors, tt.want) {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestCheckDemoRouteEvidenceRequirements(t *testing.T) {
	for _, tt := range []struct {
		name string
		edit func(*Bundle)
	}{
		{
			name: "non harness source",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "harness:route", func(entry *BundleEntry) {
					entry.SourceClass = "review"
				})
			},
		},
		{
			name: "wrong evidence kind",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "harness:route", func(entry *BundleEntry) {
					entry.EvidenceKind = "review_note"
				})
			},
		},
		{
			name: "missing opencode component",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "harness:route", func(entry *BundleEntry) {
					entry.ObservedComponents = []string{"gsd-redux", "minimax-m2.5"}
				})
			},
		},
		{
			name: "synthetic digest",
			edit: func(bundle *Bundle) {
				setManifestEntry(bundle, "harness:route", func(entry *BundleEntry) {
					entry.Digest = digestPlaceholder(entry.Ref + entry.Resolver)
				})
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bundle := demoGateBundle()
			tt.edit(&bundle)
			result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StateFail || !hasError(result.Errors, "retained structured OpenCode/GSD/MiniMax") {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestCheckDemoAcceptsMiniMaxAliasesAndNormalizesComponents(t *testing.T) {
	for _, minimaxComponent := range []string{"minimax", "minimax-m2.5", "minimax-m2"} {
		t.Run(minimaxComponent, func(t *testing.T) {
			bundle := demoGateBundle()
			setManifestEntry(&bundle, "harness:route", func(entry *BundleEntry) {
				entry.ObservedComponents = []string{" OpenCode ", "GSD-Redux", minimaxComponent}
			})
			result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
			if result.State != StatePass {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestCheckDemoRequiresVerificationOrReviewAssessed(t *testing.T) {
	bundle := demoGateBundle()
	setRowState(&bundle, "PC-VERIFICATION", StateNotAssessed, "verification not assessed")
	setRowState(&bundle, "PC-REVIEW", StateNotAssessed, "review not assessed")
	result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "requires PC-VERIFICATION or PC-REVIEW") {
		t.Fatalf("result = %+v", result)
	}
}

func TestCheckDemoCannotVerifyClosureCap(t *testing.T) {
	bundle := demoGateBundle()
	setRowState(&bundle, "PC-AUTHORITY", StateCannotVerify, "authority could not be verified")
	setRowState(&bundle, "PC-ATTESTATION", StateCannotVerify, "attestation could not be verified")
	for i := range bundle.Packet.ResidualGaps {
		if bundle.Packet.ResidualGaps[i].RowID == "PC-AUTHORITY" || bundle.Packet.ResidualGaps[i].RowID == "PC-ATTESTATION" {
			bundle.Packet.ResidualGaps[i].State = StateCannotVerify
			bundle.Packet.ResidualGaps[i].ClosureEvidence = ""
		}
	}
	refreshPacketDigest(&bundle)
	result := CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "allows at most one cannot_verify row without closure path") {
		t.Fatalf("result = %+v", result)
	}

	for i := range bundle.Packet.ResidualGaps {
		if bundle.Packet.ResidualGaps[i].RowID == "PC-ATTESTATION" {
			bundle.Packet.ResidualGaps[i].ClosureEvidence = "signed packet and witness evidence"
		}
	}
	refreshPacketDigest(&bundle)
	result = CheckDemoFirstPacket(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StatePass {
		t.Fatalf("closure evidence should exempt one cannot_verify row: %+v", result)
	}
}

func validBundle() Bundle {
	rows := []Row{
		row("PC-CHANGE", StatePass, "PR 38 change is bound to source metadata.", []string{"git:change"}),
		row("PC-INITIATOR", StatePartial, "Task source is present but initiator identity is partial.", []string{"task:pr-body"}),
		row("PC-AGENT-ROUTE", StatePartial, "OpenCode/GSD route evidence is retained.", []string{"harness:route"}),
		row("PC-MUTATION", StatePass, "Changed files and commit range are retained.", []string{"git:change"}),
		row("PC-VERIFICATION", StatePass, "CI run evidence is retained.", []string{"ci:run"}),
		row("PC-REVIEW", StatePartial, "External review artifact is retained.", []string{"review:packet"}),
		row("PC-AUTHORITY", StateNotAssessed, "Authority was not assessed for this fixture.", nil),
		row("PC-THEATER", StatePass, "No P0 theater finding triggered.", []string{"theater:assessment"}),
		row("PC-ATTESTATION", StateNotAssessed, "Signed or external attestation was not assessed.", nil),
		row("PC-DECISION", StatePass, "Next decision owner is identified.", []string{"decision:owner"}),
		row("PC-RESIDUAL-GAPS", StatePartial, "Some rows remain partial or not assessed.", []string{"gap:residual"}),
	}
	bundle := Bundle{
		Packet: Packet{
			PacketVersion:   PacketSchemaVersion,
			PacketID:        "packet-pr-38",
			SourceChange:    SourceChange{Repository: "fall-out-bug/sdp-trace", ChangeID: "PR-38", URL: "https://github.com/fall-out-bug/sdp-trace/pull/38"},
			GeneratedAt:     "2026-05-11T12:00:00Z",
			AuthoringMethod: AuthoringToolGenerated,
			SelectedProfile: "change-host-rich-v0",
			RedactionPolicy: "not_needed",
			BundleRef:       "bundle-pr-38",
			PacketState:     "review_ready",
			Projection:      Projection{Kind: ProjectionCanonical, Canonical: true, ArtifactRef: "packet:markdown"},
			Rows:            rows,
			ResidualGaps:    validResidualGaps(),
			DecisionOwners:  validDecisionOwners(),
			NonApproval:     "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.",
		},
		Manifest: BundleManifest{
			SchemaVersion: BundleSchemaVersion,
			BundleID:      "bundle-pr-38",
			Entries: []BundleEntry{
				entry("git:change", "git"),
				entry("task:pr-body", "change_host"),
				entry("harness:route", "harness"),
				entry("ci:run", "ci"),
				entry("review:packet", "review"),
				entry("theater:assessment", "witness"),
				entry("gap:residual", "manual"),
				entry("decision:owner", "manual"),
			},
		},
	}
	bundle.Manifest.PacketDigest = PacketDigest(bundle.Packet)
	return bundle
}

func demoGateBundle() Bundle {
	bundle := validBundle()
	setManifestEntry(&bundle, "harness:route", func(entry *BundleEntry) {
		entry.EvidenceKind = "harness_route_observation"
		entry.ObservedComponents = []string{"opencode", "gsd-redux", "minimax-m2.5"}
	})
	return bundle
}

func validResidualGaps() []ResidualGap {
	return []ResidualGap{
		{RowID: "PC-INITIATOR", State: StatePartial, Reason: "initiator is bound to PR body only", ClosureEvidence: "GitHub issue or retained prompt boundary"},
		{RowID: "PC-AGENT-ROUTE", State: StatePartial, Reason: "route evidence is retained observation output only", ClosureEvidence: "complete delegation chain"},
		{RowID: "PC-REVIEW", State: StatePartial, Reason: "review is retained externally", ClosureEvidence: "GitHub review or signed retained review"},
		{RowID: "PC-AUTHORITY", State: StateNotAssessed, Reason: "authority profile is outside fixture", ClosureEvidence: "authority profile evidence"},
		{RowID: "PC-ATTESTATION", State: StateNotAssessed, Reason: "signed trust is outside fixture", ClosureEvidence: "signed packet and witness evidence"},
	}
}

func validDecisionOwners() []DecisionOwner {
	return []DecisionOwner{
		{Decision: "merge", Owner: "maintainer", State: StateNotAssessed, Reason: "packet is not approval"},
		{Decision: "release", Owner: "release owner", State: StateNotAssessed, Reason: "packet is not release approval"},
		{Decision: "risk_acceptance", Owner: "risk owner", State: StateNotAssessed, Reason: "packet is not risk acceptance"},
		{Decision: "security_review", Owner: "security owner", State: StateNotAssessed, Reason: "packet is not security review"},
	}
}

func row(id, state, summary string, refs []string) Row {
	r := Row{ID: id, State: state, Summary: summary, EvidenceRefs: refs, Owner: "maintainer"}
	if missingReasonStates[state] {
		r.Reason = "fixture records non-pass state explicitly"
	}
	return r
}

func entry(ref, sourceClass string) BundleEntry {
	return BundleEntry{
		Ref:             ref,
		SourceClass:     sourceClass,
		Digest:          "sha256:1111111111111111111111111111111111111111111111111111111111111111",
		RetainedForm:    "raw",
		RedactionStatus: "not_needed",
		Resolver:        "examples/change-evidence-packet/" + ref,
	}
}

func setRow(bundle *Bundle, id string, row Row) {
	for i := range bundle.Packet.Rows {
		if bundle.Packet.Rows[i].ID == id {
			bundle.Packet.Rows[i] = row
			return
		}
	}
}

func removeRow(bundle *Bundle, id string) {
	for i := range bundle.Packet.Rows {
		if bundle.Packet.Rows[i].ID == id {
			bundle.Packet.Rows = append(bundle.Packet.Rows[:i], bundle.Packet.Rows[i+1:]...)
			return
		}
	}
}

func setRowState(bundle *Bundle, id, state, reason string) {
	for i := range bundle.Packet.Rows {
		if bundle.Packet.Rows[i].ID == id {
			bundle.Packet.Rows[i].State = state
			bundle.Packet.Rows[i].Reason = reason
			return
		}
	}
}

func setManifestEntry(bundle *Bundle, ref string, edit func(*BundleEntry)) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			edit(&bundle.Manifest.Entries[i])
			return
		}
	}
}

func hasError(errors []string, want string) bool {
	for _, err := range errors {
		if strings.Contains(err, want) {
			return true
		}
	}
	return false
}

func assertErrorsInOrder(t *testing.T, errors []string, wants []string) {
	t.Helper()
	next := 0
	for _, err := range errors {
		if next < len(wants) && strings.Contains(err, wants[next]) {
			next++
		}
	}
	if next != len(wants) {
		t.Fatalf("errors = %v, want ordered subsequence %v", errors, wants)
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

func validGitHubInput() GitHubPREvidenceInput {
	return GitHubPREvidenceInput{
		SchemaVersion: "github-pr-evidence-input.v0",
		PR: GitHubPR{
			Number:  5,
			URL:     "https://github.com/fall-out-bug/sdp-trace-demo-jvm-gsd/pull/5",
			Title:   "Add stats endpoint",
			BodyRef: "https://github.com/fall-out-bug/sdp-trace-demo-jvm-gsd/pull/5",
			Author:  "demo-agent",
			BaseRef: "main",
			HeadRef: "feature/stats",
			HeadSHA: "1111111111111111111111111111111111111111",
		},
		CommitRange: GitHubCommitRange{
			Base:            "0000000000000000000000000000000000000000",
			Head:            "1111111111111111111111111111111111111111",
			ChangedFilesRef: "git diff --name-only base..head",
		},
		Checks:         []GitHubCheck{{Name: "bazel-test", URL: "https://github.com/fall-out-bug/sdp-trace-demo-jvm-gsd/actions/runs/1", Conclusion: "success", ArtifactRefs: []string{"test-report"}}},
		Artifacts:      []GitHubArtifact{{Name: "test-report", Resolver: "https://github.com/fall-out-bug/sdp-trace-demo-jvm-gsd/actions/runs/1/artifacts", RetainedForm: "external_ref"}},
		Reviews:        []GitHubReview{{Reviewer: "external-reviewer", Resolver: "REVIEW_FEATURE5_GLM47.md", State: StatePartial}},
		AgentRouteRefs: []string{".opencode/session-observation.json"},
	}
}

func refreshPacketDigest(bundle *Bundle) {
	bundle.Manifest.PacketDigest = PacketDigest(bundle.Packet)
}

func rowsByID(rows []Row) map[string]Row {
	byID := map[string]Row{}
	for _, row := range rows {
		byID[row.ID] = row
	}
	return byID
}

func writeJSONForTest(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write json: %v", err)
	}
}

func assertJSONHasKeys(t *testing.T, value any, present, absent []string) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("unmarshal json object %s: %v", data, err)
	}
	for _, key := range present {
		if _, ok := fields[key]; !ok {
			t.Fatalf("json missing key %q in %s", key, data)
		}
	}
	for _, key := range absent {
		if _, ok := fields[key]; ok {
			t.Fatalf("json unexpectedly included key %q in %s", key, data)
		}
	}
}

type jsonFieldExpectation struct {
	name      string
	omitempty bool
}

func requiredJSONFields(names ...string) []jsonFieldExpectation {
	fields := make([]jsonFieldExpectation, 0, len(names))
	for _, name := range names {
		fields = append(fields, jsonFieldExpectation{name: name})
	}
	return fields
}

func optionalJSONFields(names ...string) []jsonFieldExpectation {
	fields := make([]jsonFieldExpectation, 0, len(names))
	for _, name := range names {
		fields = append(fields, jsonFieldExpectation{name: name, omitempty: true})
	}
	return fields
}

func assertBoolCatalog(t *testing.T, name string, got map[string]bool, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s length = %d want %d: %#v", name, len(got), len(want), got)
	}
	for _, key := range want {
		if !got[key] {
			t.Fatalf("%s missing %q: %#v", name, key, got)
		}
	}
	for key, value := range got {
		if !value {
			t.Fatalf("%s[%q] = false", name, key)
		}
		if !containsString(want, key) {
			t.Fatalf("%s has unexpected key %q: %#v", name, key, got)
		}
	}
}

func assertJSONTags(t *testing.T, value any, required, optional []jsonFieldExpectation) {
	t.Helper()
	want := map[string]bool{}
	for _, field := range append(required, optional...) {
		want[field.name] = field.omitempty
	}
	typ := reflect.TypeOf(value)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		t.Fatalf("assertJSONTags requires struct, got %s", typ.Kind())
	}
	if typ.NumField() != len(want) {
		t.Fatalf("%s JSON field count = %d want %d", typ.Name(), typ.NumField(), len(want))
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		name, omitempty := splitJSONTag(tag)
		wantOmitEmpty, ok := want[name]
		if !ok {
			t.Fatalf("%s.%s has unexpected json tag %q", typ.Name(), field.Name, tag)
		}
		if omitempty != wantOmitEmpty {
			t.Fatalf("%s.%s json tag = %q, omitempty = %t want %t", typ.Name(), field.Name, tag, omitempty, wantOmitEmpty)
		}
	}
}

func splitJSONTag(tag string) (string, bool) {
	parts := strings.Split(tag, ",")
	for _, option := range parts[1:] {
		if option == "omitempty" {
			return parts[0], true
		}
	}
	return parts[0], false
}
