package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckForbiddenTrackedPaths(t *testing.T) {
	cases := []struct {
		name     string
		tracked  []string
		want     int
		contains string
	}{
		{
			name:    "clean",
			tracked: []string{"docs/README.md", "cmd/sdp-trace/main.go"},
			want:    0,
		},
		{
			name:     "worktree",
			tracked:  []string{".worktrees/012/main.go"},
			want:     1,
			contains: ".worktrees/012/main.go",
		},
		{
			name:     "codex runs",
			tracked:  []string{".codex-subagents/runs/abc/output.md"},
			want:     1,
			contains: ".codex-subagents/runs/abc/output.md",
		},
		{
			name:     "sdp-trace local",
			tracked:  []string{".sdp-trace-runs/session.json"},
			want:     1,
			contains: ".sdp-trace-runs/session.json",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkForbiddenTrackedPaths(c.tracked)
			if len(got) != c.want {
				t.Fatalf("got %d findings, want %d: %v", len(got), c.want, got)
			}
			if c.want > 0 && !strings.Contains(got[0], c.contains) {
				t.Fatalf("finding %q does not contain %q", got[0], c.contains)
			}
		})
	}
}

func TestCheckRootArtifactClutter(t *testing.T) {
	cases := []struct {
		name     string
		tracked  []string
		want     int
		contains string
	}{
		{
			name:    "clean",
			tracked: []string{"docs/README.md", "specs/010/PR_DESCRIPTION.md"},
			want:    0,
		},
		{
			name:     "root PR_DESCRIPTION",
			tracked:  []string{"PR_DESCRIPTION.md"},
			want:     1,
			contains: "PR_DESCRIPTION.md",
		},
		{
			name:     "root design-note",
			tracked:  []string{"design-note.md"},
			want:     1,
			contains: "design-note.md",
		},
		{
			name:     "root reviews",
			tracked:  []string{"reviews/pi-review.md"},
			want:     1,
			contains: "reviews/",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkRootArtifactClutter(c.tracked)
			if len(got) != c.want {
				t.Fatalf("got %d findings, want %d: %v", len(got), c.want, got)
			}
			if c.want > 0 && !strings.Contains(got[0], c.contains) {
				t.Fatalf("finding %q does not contain %q", got[0], c.contains)
			}
		})
	}
}

func TestCheckRootExecutables(t *testing.T) {
	dir := t.TempDir()
	gitInit(t, dir)

	writeFile(t, dir, "README.md", "# readme\n")
	writeFile(t, dir, "binary.bin", string([]byte{0x7f, 'E', 'L', 'F', 1, 2, 3}))
	gitAdd(t, dir, "README.md")
	gitAdd(t, dir, "binary.bin")
	gitCommit(t, dir, "init")

	tracked := []string{"README.md", "binary.bin"}
	got := checkRootExecutables(dir, tracked)
	if len(got) != 1 {
		t.Fatalf("got %d findings, want 1: %v", len(got), got)
	}
	if !strings.Contains(got[0], "binary.bin") {
		t.Fatalf("finding does not mention binary.bin: %s", got[0])
	}
}

func TestCheckAbsoluteHomePaths(t *testing.T) {
	dir := t.TempDir()
	gitInit(t, dir)

	writeFile(t, dir, "docs/clean.md", "no paths here\n")
	writeFile(t, dir, "docs/bad.md", "run in /home/user/projects\n")
	writeFile(t, dir, "docs/prose.md", "rejects absolute `/home/...` paths\n")
	gitAdd(t, dir, "docs/clean.md")
	gitAdd(t, dir, "docs/bad.md")
	gitAdd(t, dir, "docs/prose.md")
	gitCommit(t, dir, "init")

	tracked := []string{"docs/clean.md", "docs/bad.md", "docs/prose.md"}
	got := checkAbsoluteHomePaths(dir, tracked)
	if len(got) != 1 {
		t.Fatalf("got %d findings, want 1: %v", len(got), got)
	}
	if !strings.Contains(got[0], "docs/bad.md") {
		t.Fatalf("finding does not mention docs/bad.md: %s", got[0])
	}
}

func TestRunPassesCleanRepo(t *testing.T) {
	dir := t.TempDir()
	gitInit(t, dir)
	writeFile(t, dir, "README.md", "# clean\n")
	gitAdd(t, dir, "README.md")
	gitCommit(t, dir, "init")

	if err := runAt(dir); err != nil {
		t.Fatalf("run on clean repo: %v", err)
	}
}

func TestRepoRoot(t *testing.T) {
	root := repoRoot()
	if !strings.Contains(root, "sdp-trace") {
		t.Fatalf("repoRoot() = %q, want path containing sdp-trace", root)
	}
}

func TestExitCode(t *testing.T) {
	var stderr bytes.Buffer
	if got := exitCode(nil, &stderr); got != 0 {
		t.Fatalf("exitCode(nil) = %d, want 0", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("exitCode(nil) wrote %q to stderr, want nothing", stderr.String())
	}
	stderr.Reset()
	if got := exitCode(fmt.Errorf("test error"), &stderr); got != 1 {
		t.Fatalf("exitCode(err) = %d, want 1", got)
	}
	if !strings.Contains(stderr.String(), "test error") {
		t.Fatalf("exitCode(err) wrote %q to stderr, want 'test error'", stderr.String())
	}
}

func gitInit(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "config", "user.email", "test@example.invalid")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git config user.email: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git config user.name: %v\n%s", err, out)
	}
}

func gitAdd(t *testing.T, dir, path string) {
	t.Helper()
	cmd := exec.Command("git", "add", path)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
}

func gitCommit(t *testing.T, dir, msg string) {
	t.Helper()
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
}

func writeFile(t *testing.T, dir, path, content string) {
	t.Helper()
	full := filepath.Join(dir, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}
