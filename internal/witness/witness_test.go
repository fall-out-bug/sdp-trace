package witness

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestFetchGitHubOIDCTokenSuccess(t *testing.T) {
	env := map[string]string{
		"ACTIONS_ID_TOKEN_REQUEST_URL":   "https://token.actions.githubusercontent.com/token?audience=old",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN": "request-token",
	}
	oldClient := http.DefaultClient
	requested := false
	http.DefaultClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requested = true
			if req.Method != http.MethodGet {
				t.Fatalf("method = %s", req.Method)
			}
			if got := req.URL.Query().Get("audience"); got != "sdp-trace" {
				t.Fatalf("audience = %s", got)
			}
			if got := req.Header.Get("Authorization"); got != "bearer request-token" {
				t.Fatalf("authorization = %q", got)
			}
			if got := req.Header.Get("Accept"); got != "application/json" {
				t.Fatalf("accept = %q", got)
			}
			return &http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"value":"github-oidc-token"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	t.Cleanup(func() { http.DefaultClient = oldClient })

	token, err := FetchGitHubOIDCToken(env)
	if err != nil {
		t.Fatalf("FetchGitHubOIDCToken: %v", err)
	}
	if !requested {
		t.Fatalf("token endpoint was not requested")
	}
	if token != "github-oidc-token" {
		t.Fatalf("token = %q", token)
	}
}

func TestFetchGitHubOIDCTokenInvalidRequestHost(t *testing.T) {
	_, err := FetchGitHubOIDCToken(map[string]string{
		"ACTIONS_ID_TOKEN_REQUEST_URL": "https://malicious.example/token",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "unexpected oidc request host: malicious.example") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchGitHubOIDCTokenMalformedRequestURL(t *testing.T) {
	_, err := FetchGitHubOIDCToken(map[string]string{
		"ACTIONS_ID_TOKEN_REQUEST_URL": "://bad-url",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestFetchGitHubOIDCTokenNonSuccessResponse(t *testing.T) {
	env := map[string]string{
		"ACTIONS_ID_TOKEN_REQUEST_URL":   "https://token.actions.githubusercontent.com/token",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN": "request-token",
	}
	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "500 Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	t.Cleanup(func() { http.DefaultClient = oldClient })

	_, err := FetchGitHubOIDCToken(env)
	if err == nil || !strings.Contains(err.Error(), "oidc token request returned 500") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchGitHubOIDCTokenEmptyValue(t *testing.T) {
	env := map[string]string{
		"ACTIONS_ID_TOKEN_REQUEST_URL":   "https://token.actions.githubusercontent.com/token",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN": "request-token",
	}
	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	t.Cleanup(func() { http.DefaultClient = oldClient })

	_, err := FetchGitHubOIDCToken(env)
	if err == nil || !strings.Contains(err.Error(), "oidc token response missing value") {
		t.Fatalf("unexpected error: %v", err)
	}
}
