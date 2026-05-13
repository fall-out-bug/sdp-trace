package witness

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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

func TestCRAPHelperEdges(t *testing.T) {
	env := environmentFromEntries([]string{"A=1", "MALFORMED", "B=two=parts"})
	if env["A"] != "1" || env["B"] != "two=parts" {
		t.Fatalf("environmentFromEntries = %v", env)
	}
	if got := stringItems([]any{"sdp", 1, "trace"}); strings.Join(got, ",") != "sdp,trace" {
		t.Fatalf("stringItems = %v", got)
	}
	if got := audienceString([]any{"sdp-trace", "other"}); got != "sdp-trace,other" {
		t.Fatalf("audienceString list = %q", got)
	}
	if got := audienceString("sdp-trace"); got != "sdp-trace" {
		t.Fatalf("audienceString scalar = %q", got)
	}

	root := t.TempDir()
	for name, payload := range map[string]string{
		"one": `{"run_id":"run-1"}`,
		"two": `{"id":"run-2"}`,
	} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "run.json"), []byte(payload), 0o644); err != nil {
			t.Fatalf("write run: %v", err)
		}
	}
	runIDs, err := runIDsFromRoot(root)
	if err != nil {
		t.Fatalf("runIDsFromRoot: %v", err)
	}
	if strings.Join(runIDs, ",") != "run-1,run-2" {
		t.Fatalf("runIDs = %v", runIDs)
	}
	emptyIDDir := filepath.Join(t.TempDir(), "empty-id")
	if err := os.MkdirAll(emptyIDDir, 0o755); err != nil {
		t.Fatalf("mkdir empty-id: %v", err)
	}
	if err := os.WriteFile(filepath.Join(emptyIDDir, "run.json"), []byte(`{"id":""}`), 0o644); err != nil {
		t.Fatalf("write empty run: %v", err)
	}
	runIDs, err = runIDsFromDirs([]string{emptyIDDir})
	if err != nil {
		t.Fatalf("empty run id: %v", err)
	}
	if len(runIDs) != 0 {
		t.Fatalf("empty run id included: %v", runIDs)
	}
	badRunDir := filepath.Join(t.TempDir(), "bad-run")
	if err := os.MkdirAll(badRunDir, 0o755); err != nil {
		t.Fatalf("mkdir bad-run: %v", err)
	}
	if err := os.WriteFile(filepath.Join(badRunDir, "run.json"), []byte(`{`), 0o644); err != nil {
		t.Fatalf("write bad run: %v", err)
	}
	if _, err := runIDsFromDirs([]string{badRunDir}); err == nil {
		t.Fatalf("malformed run.json accepted")
	}
}

func TestMissingCustomerPKIInputsIncludesKeyChoice(t *testing.T) {
	missing := missingCustomerPKIInputs(ProfileOptions{})
	if !testContainsString(missing, "--customer-pki-public-cert|--customer-pki-public-key") {
		t.Fatalf("missing key choice not reported: %v", missing)
	}
	complete := missingCustomerPKIInputs(ProfileOptions{
		CustomerPKIAuthorityPolicy: "policy.json",
		CustomerPKIPayloadDigest:   strings.Repeat("a", 64),
		CustomerPKIFreshness:       "freshness.json",
		CustomerPKIPublicKey:       "key.pem",
	})
	if len(complete) != 0 {
		t.Fatalf("complete inputs missing = %v", complete)
	}
}

func testContainsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
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

func TestWriteGitHubActionsWritesWitnessRecord(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "run.json"), []byte(`{"id":"run"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "ci-witness.json")
	record, err := WriteGitHubActions(outPath, root, "", map[string]string{})
	if err != nil {
		t.Fatalf("WriteGitHubActions: %v", err)
	}
	if record.Status != StatusCannotVerify {
		t.Fatalf("status = %s", record.Status)
	}
	loaded, err := Load(outPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Reason != ReasonMissingCIIdentity {
		t.Fatalf("reason = %s", loaded.Reason)
	}
	if IsPassingCI(loaded) {
		t.Fatalf("missing identity witness must not pass")
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

func TestGitHubActionsWitnessOIDCFetcherErrorCannotVerify(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "run.json"), []byte(`{"id":"run"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	record, err := BuildGitHubActionsWithFetcher(root, "", completeGitHubEnv(), func(map[string]string) (string, error) {
		return "", errors.New("oidc broker failed")
	})
	if err != nil {
		t.Fatalf("BuildGitHubActions: %v", err)
	}
	if record.Status != StatusCannotVerify {
		t.Fatalf("status = %s", record.Status)
	}
	if record.Reason != ReasonInvalidCIOIDC {
		t.Fatalf("reason = %s", record.Reason)
	}
	if record.TrustScope != TrustScopeLocalObserved {
		t.Fatalf("trust_scope = %s", record.TrustScope)
	}
}

func TestGitHubActionsWitnessMismatchedOIDCClaimsCannotVerify(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "run.json"), []byte(`{"id":"run"}`), 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	record, err := BuildGitHubActionsWithFetcher(root, "", completeGitHubEnv(), func(map[string]string) (string, error) {
		return fakeOIDCToken(t, map[string]any{
			"iss":        githubOIDCIssuer,
			"aud":        "sdp-trace",
			"repository": "org/other",
			"ref":        "refs/heads/main",
			"sha":        "abc123",
		}), nil
	})
	if err != nil {
		t.Fatalf("BuildGitHubActions: %v", err)
	}
	if record.Status != StatusCannotVerify {
		t.Fatalf("status = %s", record.Status)
	}
	if record.Reason != ReasonInvalidCIOIDC {
		t.Fatalf("reason = %s", record.Reason)
	}
	if record.TrustScope != TrustScopeLocalObserved {
		t.Fatalf("trust_scope = %s", record.TrustScope)
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
		"ACTIONS_ID_TOKEN_REQUEST_URL":   "https://token.actions.githubusercontent.com/token",
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

func TestFetchGitHubOIDCTokenRejectsSiblingHostAndHTTP(t *testing.T) {
	for _, requestURL := range []string{
		"https://evilactions.githubusercontent.com/token",
		"https://pipelines.actions.githubusercontent.com/token",
		"http://token.actions.githubusercontent.com/token",
	} {
		_, err := FetchGitHubOIDCToken(map[string]string{
			"ACTIONS_ID_TOKEN_REQUEST_URL": requestURL,
		})
		if err == nil {
			t.Fatalf("expected %s to be rejected", requestURL)
		}
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
