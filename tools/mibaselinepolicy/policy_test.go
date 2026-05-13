package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPolicyRejectsExistingBaselineWithProductionGoChange(t *testing.T) {
	err := checkPolicy(policyInput{
		changed: []string{
			"tools/qualitycheck/file-mi-baseline.json",
			"internal/packet/packet.go",
		},
		baselines: []string{"tools/qualitycheck/file-mi-baseline.json"},
		baselineExists: func(path string) (bool, error) {
			return path == "tools/qualitycheck/file-mi-baseline.json", nil
		},
	})
	if err == nil {
		t.Fatalf("checkPolicy returned nil, want mixed baseline/product change rejection")
	}
	if !strings.Contains(err.Error(), "MI baseline changes must be reviewed separately") {
		t.Fatalf("error = %v", err)
	}
}

func TestPolicyAllowsExistingBaselineWithoutProductionGoChange(t *testing.T) {
	err := checkPolicy(policyInput{
		changed: []string{
			"tools/qualitycheck/file-mi-baseline.json",
			"docs/ci-check-policy.md",
		},
		baselines: []string{"tools/qualitycheck/file-mi-baseline.json"},
		baselineExists: func(string) (bool, error) {
			return true, nil
		},
	})
	if err != nil {
		t.Fatalf("checkPolicy returned %v, want pass", err)
	}
}

func TestPolicyAllowsExistingBaselineWithTestOnlyGoChange(t *testing.T) {
	err := checkPolicy(policyInput{
		changed: []string{
			"tools/qualitycheck/file-mi-baseline.json",
			"tools/mibaselinepolicy/policy_test.go",
		},
		baselines: []string{"tools/qualitycheck/file-mi-baseline.json"},
		baselineExists: func(string) (bool, error) {
			return true, nil
		},
	})
	if err != nil {
		t.Fatalf("checkPolicy returned %v, want test-only Go changes allowed", err)
	}
}

func TestProductionGoFileScope(t *testing.T) {
	cases := map[string]bool{
		"cmd/sdp-trace/main.go":                 true,
		"internal/packet/packet.go":             true,
		"tools/mibaselinepolicy/policy.go":      true,
		"tools/mibaselinepolicy/policy_test.go": false,
		"docs/example.go":                       false,
		"tools/readme.md":                       false,
	}
	for path, want := range cases {
		if got := productionGoFile(path); got != want {
			t.Fatalf("productionGoFile(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestPolicyAllowsFirstBaselineIntroduction(t *testing.T) {
	err := checkPolicy(policyInput{
		changed: []string{
			"tools/qualitycheck/function-mi-baseline.json",
			"cmd/sdp-trace/main.go",
		},
		baselines: []string{"tools/qualitycheck/function-mi-baseline.json"},
		baselineExists: func(string) (bool, error) {
			return false, nil
		},
	})
	if err != nil {
		t.Fatalf("checkPolicy returned %v, want first baseline introduction allowed", err)
	}
}

func TestPolicySurfacesBaselineExistenceErrors(t *testing.T) {
	want := errors.New("git unavailable")
	err := checkPolicy(policyInput{
		changed: []string{
			"tools/qualitycheck/function-mi-baseline.json",
			"tools/qualitycheck/main.go",
		},
		baselines: []string{"tools/qualitycheck/function-mi-baseline.json"},
		baselineExists: func(string) (bool, error) {
			return false, want
		},
	})
	if !errors.Is(err, want) {
		t.Fatalf("checkPolicy error = %v, want %v", err, want)
	}
}

func TestRunRejectsMissingBaseRef(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run(nil, strings.NewReader("internal/packet/packet.go\n"), &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit = %d, want usage failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "missing required -base-ref") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestRunRejectsExistingBaselineWithProductionGoChange(t *testing.T) {
	repo := initGitRepo(t)
	writeFile(t, repo, "tools/qualitycheck/file-mi-baseline.json", "{}\n")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-m", "baseline")
	withWorkingDir(t, repo, func() {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		changed := strings.NewReader("tools/qualitycheck/file-mi-baseline.json\ninternal/packet/packet.go\n")
		exit := run([]string{"-base-ref", "HEAD"}, changed, &stdout, &stderr)
		if exit != 1 {
			t.Fatalf("exit = %d, want policy failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "file-mi-baseline.json") {
			t.Fatalf("stderr = %s", stderr.String())
		}
	})
}

func TestRunHonorsCustomBaselineFlag(t *testing.T) {
	repo := initGitRepo(t)
	writeFile(t, repo, "quality/custom-mi-baseline.json", "{}\n")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-m", "custom baseline")
	withWorkingDir(t, repo, func() {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		changed := strings.NewReader("quality/custom-mi-baseline.json\ntools/mibaselinepolicy/main.go\n")
		exit := run([]string{"-base-ref", "HEAD", "-baseline", "quality/custom-mi-baseline.json"}, changed, &stdout, &stderr)
		if exit != 1 {
			t.Fatalf("exit = %d, want policy failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "custom-mi-baseline.json") {
			t.Fatalf("stderr = %s", stderr.String())
		}
	})
}

func TestRepeatedFlagRejectsEmptyValues(t *testing.T) {
	var values repeatedFlag
	if err := values.Set(""); err == nil {
		t.Fatalf("Set empty value returned nil, want error")
	}
}

func TestRunAllowsFirstBaselineIntroduction(t *testing.T) {
	repo := initGitRepo(t)
	writeFile(t, repo, "README.md", "repo\n")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-m", "initial")
	withWorkingDir(t, repo, func() {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		changed := strings.NewReader("tools/qualitycheck/file-mi-baseline.json\ntools/mibaselinepolicy/main.go\n")
		exit := run([]string{"-base-ref", "HEAD"}, changed, &stdout, &stderr)
		if exit != 0 {
			t.Fatalf("exit = %d, want pass; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
		}
	})
}

func TestRunRejectsInvalidBaseRefWhenBaselinePolicyMustInspectSource(t *testing.T) {
	repo := initGitRepo(t)
	writeFile(t, repo, "README.md", "repo\n")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-m", "initial")
	withWorkingDir(t, repo, func() {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		changed := strings.NewReader("tools/qualitycheck/file-mi-baseline.json\ninternal/packet/packet.go\n")
		exit := run([]string{"-base-ref", "definitely-not-a-ref"}, changed, &stdout, &stderr)
		if exit != 1 {
			t.Fatalf("exit = %d, want policy failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "base ref") {
			t.Fatalf("stderr = %s", stderr.String())
		}
	})
}

func initGitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}
	dir := t.TempDir()
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Test User")
	return dir
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}

func writeFile(t *testing.T, dir, path, content string) {
	t.Helper()
	fullPath := filepath.Join(dir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(fullPath), err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", fullPath, err)
	}
}

func withWorkingDir(t *testing.T, dir string, fn func()) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	defer func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working dir: %v", err)
		}
	}()
	fn()
}
