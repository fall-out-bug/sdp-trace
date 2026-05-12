package packet

import (
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
