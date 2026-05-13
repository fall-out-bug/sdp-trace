package repoobserver

func hookScript(name string) string {
	// Hooks collect metadata and diagnostics only; they do not enforce
	// sdp-trace policy or block user operations.
	return `#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$repo_root"
event="` + name + `"
ts="$(date -u +%Y%m%dT%H%M%SZ)"
out_dir=".sdp-trace/hooks/${event}/${ts}"
mkdir -p "$out_dir"

{
  printf 'event=%s\n' "$event"
  printf 'git_head=%s\n' "$(git rev-parse --verify HEAD 2>/dev/null || printf unknown)"
  printf 'git_branch=%s\n' "$(git branch --show-current 2>/dev/null || printf unknown)"
} > "$out_dir/metadata.env"

git status --short > "$out_dir/status.txt" || true
git diff --cached --name-status > "$out_dir/staged-files.txt" || true
git diff --check > "$out_dir/diff-check.txt" 2>&1 || true
`
}
