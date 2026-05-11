package packet

import (
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

func TestValidateRejectsProjectionMarkedCanonicalOverArtifact(t *testing.T) {
	bundle := validBundle()
	bundle.Packet.Projection = Projection{Kind: "github_pr_comment", Canonical: true}
	result := Validate(bundle, time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC))
	if result.State != StateFail || !hasError(result.Errors, "canonical projection must be") {
		t.Fatalf("result = %+v", result)
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
	input.AgentRouteComponents = []string{"opencode", "gsd", "minimax-m2.5"}
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
