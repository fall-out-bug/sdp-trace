package witness

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGitHubActionsWitnessMissingIdentityCannotVerify(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "run.json"), []byte(`{"id":"run"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	record, err := BuildGitHubActions(root, "", map[string]string{})
	if err != nil {
		t.Fatalf("BuildGitHubActions: %v", err)
	}
	if record.Status != StatusCannotVerify {
		t.Fatalf("status = %s", record.Status)
	}
	if record.TrustScope != TrustScopeLocalObserved {
		t.Fatalf("trust_scope = %s", record.TrustScope)
	}
	if record.Reason != ReasonMissingCIIdentity {
		t.Fatalf("reason = %s", record.Reason)
	}
	if len(record.MissingIdentityFields) == 0 {
		t.Fatalf("missing identity fields not recorded")
	}
	payload, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	if string(payload) == "" || !json.Valid(payload) {
		t.Fatalf("invalid witness json")
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("decode witness json: %v", err)
	}
	if _, ok := decoded["report_artifacts"].([]any); !ok {
		t.Fatalf("report_artifacts must serialize as array: %s", string(payload))
	}
}

func TestGitHubActionsWitnessPassesWithCompleteIdentity(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "run.json"), []byte(`{"id":"run"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}
	reportDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(reportDir, "summary.json"), []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatalf("write summary: %v", err)
	}

	record, err := BuildGitHubActionsWithFetcher(root, reportDir, completeGitHubEnv(), func(map[string]string) (string, error) {
		return fakeOIDCToken(t, map[string]any{
			"iss":        githubOIDCIssuer,
			"sub":        "repo:org/repo:ref:refs/heads/main",
			"aud":        "sdp-trace",
			"repository": "org/repo",
			"ref":        "refs/heads/main",
			"sha":        "abc123",
		}), nil
	})
	if err != nil {
		t.Fatalf("BuildGitHubActions: %v", err)
	}
	if record.Status != StatusPass {
		t.Fatalf("status = %s", record.Status)
	}
	if record.TrustScope != TrustScopeCIWitnessed {
		t.Fatalf("trust_scope = %s", record.TrustScope)
	}
	if record.Source.CommitSHA != "abc123" || record.Source.Repository != "org/repo" {
		t.Fatalf("source = %+v", record.Source)
	}
	if len(record.RunArtifacts) == 0 {
		t.Fatalf("run artifact digest missing")
	}
	if len(record.ReportArtifacts) == 0 {
		t.Fatalf("report artifact digest missing")
	}
}

func TestGitHubActionsWitnessRequiresOIDCClaims(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "run.json"), []byte(`{"id":"run"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	env := completeGitHubEnv()
	delete(env, "ACTIONS_ID_TOKEN_REQUEST_URL")
	delete(env, "ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	record, err := BuildGitHubActions(root, "", env)
	if err != nil {
		t.Fatalf("BuildGitHubActions: %v", err)
	}
	if record.Status != StatusCannotVerify {
		t.Fatalf("status = %s", record.Status)
	}
	if record.Reason != ReasonMissingCIOIDC {
		t.Fatalf("reason = %s", record.Reason)
	}
	if len(record.MissingIdentityFields) == 0 {
		t.Fatalf("missing oidc fields not recorded")
	}
}

func completeGitHubEnv() map[string]string {
	return map[string]string{
		"GITHUB_ACTIONS":                 "true",
		"GITHUB_SHA":                     "abc123",
		"GITHUB_RUN_ID":                  "42",
		"GITHUB_RUN_ATTEMPT":             "1",
		"GITHUB_WORKFLOW":                "sdp-trace",
		"GITHUB_JOB":                     "test",
		"GITHUB_ACTOR":                   "octocat",
		"GITHUB_REPOSITORY":              "org/repo",
		"GITHUB_REF":                     "refs/heads/main",
		"GITHUB_SERVER_URL":              "https://github.com",
		"ACTIONS_ID_TOKEN_REQUEST_URL":   "https://pipelines.actions.githubusercontent.com/token",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN": "request-token",
	}
}

func fakeOIDCToken(t *testing.T, claims map[string]any) string {
	t.Helper()
	header, err := json.Marshal(map[string]any{"alg": "none", "typ": "JWT"})
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(header) + "." +
		base64.RawURLEncoding.EncodeToString(payload) + "."
}
