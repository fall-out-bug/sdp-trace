package packet

import (
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

func TestRowIDForRefUsesRequiredRowOrder(t *testing.T) {
	bundle := validBundle()
	rows := rowsByID(bundle.Packet.Rows)
	if got := rowIDForRef(rows, "git:change"); got != "PC-CHANGE" {
		t.Fatalf("rowIDForRef(shared ref) = %q, want PC-CHANGE", got)
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

func setRowState(bundle *Bundle, id, state, reason string) {
	for i := range bundle.Packet.Rows {
		if bundle.Packet.Rows[i].ID == id {
			bundle.Packet.Rows[i].State = state
			bundle.Packet.Rows[i].Reason = reason
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
