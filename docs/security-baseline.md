# Security Baseline

Local security evidence for `sdp-trace` as of 2026-05-20.
This document is advisory; it does not prove absence of vulnerabilities
or guarantee production readiness.

## Probe Metadata

| Field | Value |
| --- | --- |
| Date | 2026-05-20 |
| Scanner target | Local branch snapshot after final review fixes |
| Go version | `1.22` from `go.mod` |
| `gosec` version | `dev` build from local OSS tool cache |

The scanner evidence was produced locally on this branch. Re-run the commands
below before using this ledger in a later PR; raw scanner JSON is intentionally
not checked in.

## Tool Results Summary

| Tool | Version | Result | Disposition |
|------|---------|--------|-------------|
| `go vet ./...` | go1.22 | pass | green |
| `govulncheck ./...` | latest | 0 known vulnerabilities | green |
| `gosec ./...` | dev | 132 findings | advisory — classified below |
| `gitleaks detect` (default config, tracked source) | v8.30.1 | 10 findings | all triaged false positive below |
| `gitleaks detect` (default config, working tree) | v8.30.1 | 12 findings | 10 tracked + 2 local ignored clutter |
| `gitleaks detect` (with `.gitleaks.toml`) | v8.30.1 | 0 findings | verified |

## `gosec` Finding Family Classification

| Rule | Count | Severity | Classification | Rationale |
|------|-------|----------|----------------|-----------|
| G304 | 60 | MEDIUM | **deferred advisory** | Tools read user-specified files (schemas, baselines, evidence). Path is the tool's explicit input. One `G304` in `internal/capturedepth` has a scoped reviewed `#nosec` suppression; the remaining 59 call sites need per-call-site audit before this family can be marked accepted. |
| G301 | 28 | MEDIUM | **accepted** | `os.MkdirAll` with `0o755` creates group-readable evidence directories. Intentional for shared local inspection. |
| G306 | 24 | MEDIUM | **accepted** | `os.WriteFile` with `0o644` creates group-readable evidence files. Intentional for shared local inspection. |
| G204 | 11 | MEDIUM | **deferred advisory** | Git subprocess launch with variable args. Most callers use hardcoded safe args; `gitOutput`/`runGit` variadic helpers and `mibaselinepolicy` git helpers pass external values. Review needed: validate all callers sanitize inputs before passing to git helpers. |
| G703 | 5 | HIGH | **deferred advisory** | Path traversal in repoobserver install helpers. Paths derive from install targets. Review needed: verify caller validates paths are within repository boundary. |
| G101 | 2 | HIGH | **accepted false positive** | String constants `ReasonLocalHooksBypassable` and policy map `authorityPreviewSafety` are not credentials. |
| G115 | 1 | HIGH | **accepted false positive** | `byte(r)` is guarded by `r <= 127`; overflow impossible. |
| G302 | 1 | MEDIUM | **accepted** | `os.Chmod(tmpName, 0o644)` in atomic text write matches other generated evidence file modes. |

### Deferred Advisory Detail

**G204 — Subprocess launch with variable**
- Files: `internal/trace/source_snapshot.go`, `internal/repoobserver/git.go`, `tools/mibaselinepolicy/git.go`
- Risk: If external input reaches `gitOutput`/`runGit` args or `mibaselinepolicy` ref/path parameters without validation, malicious flags could be injected.
- Required follow-up:
  - `internal/trace/source_snapshot.go:14-15` calls `gitOutput(cleanBase, "rev-parse", "HEAD^{tree}")` and `gitOutput(cleanBase, "status", "--porcelain")` — hardcoded safe args.
  - `internal/repoobserver/status_git.go:6-7` calls `gitOutput(repoRoot, "rev-parse", "--verify", "HEAD")` and `gitOutput(repoRoot, "branch", "--show-current")` — hardcoded safe args.
  - `internal/repoobserver/surface_hooks.go:18` calls `gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath")` — hardcoded safe args.
  - `internal/repoobserver/install_hooks_path.go:13,23` calls `gitOutput` and `runGit` with hardcoded config args.
  - `tools/mibaselinepolicy/git.go` callers: verify `ref` is a valid git ref (SHA or tag) and `path` is within repository before passing to `gitFileExistsAtRef` and `gitCommitExists`.

**G304 — Variable file inclusion via path**
- Files: 60 call sites across `tools/`, `internal/`, and `cmd/sdp-trace/`.
- Risk: If a path derives from external input without validation, file reads could escape the intended boundary.
- Current state: One call site in `internal/capturedepth/capture_depth.go:37` has a reviewed `#nosec G304` because `runDir` is a caller-selected local evidence root and query output does not echo provider refs.
- Required follow-up: Audit the remaining 59 call sites. For each, either confirm the path is a hardcoded or caller-validated tool input, or add a scoped `#nosec` with a recorded rationale.

**G703 — Path traversal via taint analysis**
- Files: `internal/repoobserver/install_target_overwrite.go`, `internal/repoobserver/install_gitignore_update.go`
- Risk: If `target.path` or `path` parameters contain `../` sequences, writes could escape the repository.
- Required follow-up: Verify callers validate paths are within the repository root before calling install helpers.

## `gitleaks` Findings Triage

Default-config scans:
- Tracked-source snapshot (`git archive HEAD`): **10 findings**, all in tracked files.
- Working-tree scan (includes local ignored clutter): **12 findings** (the same 10 tracked + 2 additional hits in `.codex-subagents/` local run clutter).
- None are live secrets.

The 2 working-tree-only hits are local ignored clutter, not repository evidence.
They are excluded from the triage table below.

Configured scan (with reviewed `.gitleaks.toml`): **0 findings**.

A reviewed `.gitleaks.toml` allowlist covers intentional synthetic markers.

| Category | Files | Rule | Disposition |
|----------|-------|------|-------------|
| Synthetic JWT sentinel (redaction tests) | `internal/witness/profiles_test.go` | `jwt` / `generic-api-key` | **accepted false positive** — test fixture verifying redaction behavior |
| Review dedupe keys (fixture labels) | `examples/self-trace/evidence-events.json`, `examples/self-trace/negative-native-policy-field.json`, `examples/self-trace/assessment-input.json` | `generic-api-key` | **accepted false positive** — review event labels, not credentials |
| Private-key marker (historical diff) | `specs/004-mvp-readiness-hardening/pr-review/ec8db52/packet/inputs/diff.patch` | `private-key` | **accepted false positive** — sanitized diff showing redaction sentinel changes |

## Local Ignored Clutter Policy

The following paths are ignored by `.gitignore` and must NOT be treated as
repository security evidence:

- `.worktrees/`
- `.codex-subagents/`
- `.sdp-trace-runs/`
- `.sdp-trace-report/`
- `.sdp-trace-*`

Scan them with local hygiene tooling, but report findings as **local hygiene**
only. If a finding in an ignored path is a real secret, rotate it; do not
allowlist it in repository scans.

## Reporting

For vulnerability reporting guidance, see [`.github/SECURITY.md`](../.github/SECURITY.md).

## Verification Commands

Reproduce the current baseline:

```bash
# Green gates
go vet ./...
go test -count=1 ./...
govulncheck ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check

# Advisory scans
gosec ./...

# Default-config scans (for count verification)
gitleaks detect --source . --no-git
mkdir -p /tmp/tracked-scan && rm -rf /tmp/tracked-scan/*
git archive --format=tar HEAD | tar -xf - -C /tmp/tracked-scan
gitleaks detect --source /tmp/tracked-scan --no-git

# Configured tracked-source scan
mkdir -p /tmp/tracked-scan && rm -rf /tmp/tracked-scan/*
git archive --format=tar HEAD | tar -xf - -C /tmp/tracked-scan
gitleaks detect --source /tmp/tracked-scan --no-git --config .gitleaks.toml

# Configured scan (reviewed allowlist)
gitleaks detect --source . --config .gitleaks.toml
```

## Required Follow-Up

- Review all `G703` and `G204` call sites first because they touch
  path, network, external-input, and command-execution boundaries.
- Review remaining `G304` call sites (59 of 60) and add scoped `#nosec` with
  rationale where the path is a validated tool input.
- Re-run `gitleaks` with `.gitleaks.toml` after each fixture or scanner-rule
  change; do not broaden allowlist patterns without review.
- Add CI jobs for `govulncheck`, `gosec`, and `gitleaks` only after local
  baselines are stable enough to avoid turning known fixture hits into noise.
- Re-run this baseline after any security, trust, path, witness, release-proof,
  repo-observer, or PR-review change.

## Change Log

- 2026-05-20: Initial security baseline for spec 016. All `gosec` families
  classified. Tracked-source `gitleaks` (10 findings) and working-tree
  `gitleaks` (10 tracked + 2 local clutter) triaged as false positive.
  Two `gosec` families (G204, G703) deferred as advisory pending caller audit;
  `G304` deferred advisory (1 reviewed `#nosec`, 59 remaining call sites).
