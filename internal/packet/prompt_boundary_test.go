package packet

import (
	"reflect"
	"testing"
	"time"
)

func TestClassifyPromptBoundary(t *testing.T) {
	tests := []struct {
		name   string
		input  PromptBoundary
		want   string
		effect string
	}{
		{
			name:   "clean text",
			input:  PromptBoundary{Text: "Implement the feature and run the existing tests."},
			want:   "clean",
			effect: StatePass,
		},
		{
			name:   "contaminated text",
			input:  PromptBoundary{Text: "Implement the feature and update .evidence plus packet rows."},
			want:   "contaminated",
			effect: StateFail,
		},
		{
			name: "digest only",
			input: PromptBoundary{
				Digest:        "sha256:abc",
				CaptureActor:  "recorder",
				CapturedAt:    "2026-05-12T00:00:00Z",
				CaptureMethod: "external_capture",
			},
			want:   "digest_only",
			effect: StatePartial,
		},
		{
			name:   "missing",
			input:  PromptBoundary{},
			want:   "missing",
			effect: StateCannotVerify,
		},
		{
			name:   "malformed",
			input:  PromptBoundary{Digest: "sha256:abc", CaptureActor: "recorder", CapturedAt: "not-time", CaptureMethod: "external_capture"},
			want:   "malformed",
			effect: StateCannotVerify,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyPromptBoundary(tt.input)
			if got.Verdict != tt.want || got.RouteProofEffect != tt.effect {
				t.Fatalf("classification = %+v, want verdict %s effect %s", got, tt.want, tt.effect)
			}
		})
	}
}

func TestClassifyPromptBoundaryReasons(t *testing.T) {
	tests := []struct {
		name        string
		input       PromptBoundary
		wantVerdict string
		wantReasons []string
	}{
		{
			name:        "clean text has no reasons",
			input:       PromptBoundary{Text: "Implement the feature and run tests."},
			wantVerdict: "clean",
			wantReasons: nil,
		},
		{
			name:        "contaminated text names matching phrase",
			input:       PromptBoundary{Text: "Please update provenance after the run."},
			wantVerdict: "contaminated",
			wantReasons: []string{"developer prompt contains recorder-duty phrase: update provenance"},
		},
		{
			name:        "digest only records retained metadata reason",
			input:       PromptBoundary{Digest: "sha256:abc", CaptureActor: "recorder", CapturedAt: "2026-05-12T00:00:00Z", CaptureMethod: "external_capture"},
			wantVerdict: "digest_only",
			wantReasons: []string{"prompt text unavailable; digest metadata retained"},
		},
		{
			name:        "missing records missing evidence reason",
			input:       PromptBoundary{},
			wantVerdict: "missing",
			wantReasons: []string{"prompt boundary evidence missing"},
		},
		{
			name:        "malformed records malformed metadata reason",
			input:       PromptBoundary{Digest: "sha256:abc", CaptureActor: "recorder", CapturedAt: "not-time", CaptureMethod: "external_capture"},
			wantVerdict: "malformed",
			wantReasons: []string{"prompt boundary metadata malformed"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyPromptBoundary(tt.input)
			if got.Verdict != tt.wantVerdict || !reflect.DeepEqual(got.Reasons, tt.wantReasons) {
				t.Fatalf("classification = %+v, want verdict %s reasons %#v", got, tt.wantVerdict, tt.wantReasons)
			}
		})
	}
}

func TestForbiddenRecorderDutyPhrasesPreserveCatalog(t *testing.T) {
	want := []string{
		"sdp-trace",
		".sdp-trace",
		".evidence",
		"write evidence",
		"update evidence",
		"maintain provenance",
		"update provenance",
		"update packet",
		"update bundle",
		"close gate",
		"claim verification",
	}
	if got := forbiddenRecorderDutyPhrases(); !reflect.DeepEqual(got, want) {
		t.Fatalf("forbidden recorder-duty phrases = %#v, want %#v", got, want)
	}
}

func TestPromptBoundaryEntryMetadata(t *testing.T) {
	text := PromptBoundary{Text: "Implement feature"}
	if got := promptBoundaryResolver(text); got != "prompt:text-retained" {
		t.Fatalf("text resolver = %q", got)
	}
	if got := promptBoundaryRetainedForm(text); got != "redacted" {
		t.Fatalf("text retained form = %q", got)
	}

	digest := PromptBoundary{Digest: "sha256:abc"}
	if got := promptBoundaryResolver(digest); got != "prompt:digest:sha256:abc" {
		t.Fatalf("digest resolver = %q", got)
	}
	if got := promptBoundaryRetainedForm(digest); got != "digest_only" {
		t.Fatalf("digest retained form = %q", got)
	}

	missing := PromptBoundary{}
	if got := promptBoundaryResolver(missing); got != "prompt:missing" {
		t.Fatalf("missing resolver = %q", got)
	}
	if got := promptBoundaryRetainedForm(missing); got != "not_retained" {
		t.Fatalf("missing retained form = %q", got)
	}
}

func TestBuildGitHubPromptBoundaryRequiredBlocksAgentRoute(t *testing.T) {
	tests := []struct {
		name     string
		boundary PromptBoundary
		want     string
	}{
		{
			name:     "contaminated text fails route",
			boundary: PromptBoundary{Text: "Implement feature and update packet rows."},
			want:     StateFail,
		},
		{
			name:     "missing boundary cannot verify route",
			boundary: PromptBoundary{},
			want:     StateCannotVerify,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validGitHubInput()
			input.RequirePromptBoundary = true
			input.PromptBoundary = tt.boundary
			bundle := BuildFromGitHubInput(input, testTime())
			row := rowByID(bundle.Packet.Rows, "PC-AGENT-ROUTE")
			if row.State != tt.want {
				t.Fatalf("agent route state = %s, want %s: %+v", row.State, tt.want, row)
			}
		})
	}
}

func TestBuildFromGitHubInputPreservesPacketShellAndManifest(t *testing.T) {
	input := validGitHubInput()
	input.IntegrationActions = []IntegrationAction{{Kind: "merge", Actor: "bot", Resolver: "action-log"}}
	bundle := BuildFromGitHubInput(input, time.Date(2026, 5, 12, 3, 4, 5, 0, time.FixedZone("MSK", 3*60*60)))

	if bundle.Packet.PacketID != "github-pr-5-change-evidence-packet" ||
		bundle.Packet.BundleRef != "github-pr-5-change-evidence-packet-bundle" ||
		bundle.Manifest.BundleID != "github-pr-5-change-evidence-packet-bundle" {
		t.Fatalf("packet/bundle ids drifted: packet=%q bundle_ref=%q manifest=%q", bundle.Packet.PacketID, bundle.Packet.BundleRef, bundle.Manifest.BundleID)
	}
	if bundle.Packet.GeneratedAt != "2026-05-12T00:04:05Z" {
		t.Fatalf("generated_at = %q", bundle.Packet.GeneratedAt)
	}
	if bundle.Packet.SelectedProfile != "change-host-rich-v0" ||
		bundle.Packet.RedactionPolicy != "not_assessed" ||
		bundle.Packet.PacketState != "draft" ||
		bundle.Packet.Projection.Kind != ProjectionCanonical ||
		!bundle.Packet.Projection.Canonical ||
		bundle.Packet.Projection.ArtifactRef != "packet:markdown" {
		t.Fatalf("packet shell drifted: %+v", bundle.Packet)
	}
	gotActions, ok := bundle.Packet.Extensions["integration_actions"].([]IntegrationAction)
	if !ok || !reflect.DeepEqual(gotActions, input.IntegrationActions) {
		t.Fatalf("integration actions extension = %#v", bundle.Packet.Extensions["integration_actions"])
	}
	if bundle.Manifest.SchemaVersion != BundleSchemaVersion ||
		bundle.Manifest.PacketDigest != PacketDigest(bundle.Packet) ||
		len(bundle.Manifest.Entries) == 0 {
		t.Fatalf("manifest drifted: %+v", bundle.Manifest)
	}
}

func TestBuildFromGitHubInputContaminatedPromptFindingShape(t *testing.T) {
	input := validGitHubInput()
	input.PromptBoundary = PromptBoundary{Text: "Please update packet after tests."}
	bundle := BuildFromGitHubInput(input, testTime())
	if len(bundle.Packet.TheaterFindings) != 1 {
		t.Fatalf("theater findings = %#v", bundle.Packet.TheaterFindings)
	}
	got := bundle.Packet.TheaterFindings[0]
	if got.ReasonCode != "prompt_contamination" ||
		got.State != StateFail ||
		got.Severity != "P0" ||
		got.Finding != "developer prompt contains recorder-duty phrase: update packet" ||
		!reflect.DeepEqual(got.TriggerEvidenceRefs, []string{"prompt:boundary"}) {
		t.Fatalf("prompt contamination finding drifted: %+v", got)
	}
}

func TestResolverFromList(t *testing.T) {
	resolvers := []ResolverEntry{
		{Ref: "artifact:a", Resolver: "a.zip"},
		{Ref: "artifact:b", Resolver: "b.zip"},
	}
	if got := resolverFromList(resolvers, "artifact:b"); got != "b.zip" {
		t.Fatalf("resolver = %q", got)
	}
	if got := resolverFromList(resolvers, "artifact:c"); got != "" {
		t.Fatalf("missing resolver = %q", got)
	}
}

func TestBuildFromGitHubInputRedactsSecretLikeResolvers(t *testing.T) {
	bundle := BuildFromGitHubInput(GitHubPREvidenceInput{
		PR:          GitHubPR{Number: 1, URL: "https://github.com/example/repo/pull/1", BaseRef: "main", HeadRef: "feature", HeadSHA: "head"},
		CommitRange: GitHubCommitRange{Base: "base", Head: "head"},
		Checks: []GitHubCheck{{
			Name:         "ci",
			URL:          "https://example.invalid?token=SECRET_TOKEN_SHOULD_NOT_APPEAR",
			Conclusion:   "success",
			ArtifactRefs: []string{"packet"},
		}},
		Artifacts: []GitHubArtifact{{
			Name:         "packet",
			Resolver:     "https://example.invalid/artifact?token=SECRET_TOKEN_SHOULD_NOT_APPEAR",
			RetainedForm: "external_ref",
			Digest:       "sha256:artifact",
		}},
	}, testTime())
	for _, entry := range bundle.Manifest.Entries {
		if entry.Resolver == "https://example.invalid?token=SECRET_TOKEN_SHOULD_NOT_APPEAR" ||
			entry.Resolver == "https://example.invalid/artifact?token=SECRET_TOKEN_SHOULD_NOT_APPEAR" {
			t.Fatalf("secret-like resolver was not redacted: %+v", entry)
		}
	}
}

func testTime() time.Time {
	return time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
}
