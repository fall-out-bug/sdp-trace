# Security Policy

## Supported Versions

`sdp-trace` is pre-1.0 and does not yet have a stable release track.
Only the latest commit on `main` receives security-related updates.

| Version | Supported |
|---------|-----------|
| `main`  | best-effort advisories only |
| older commits | not supported |

## Reporting a Vulnerability

Email: `security+report@fall-out-bug.dev`

If email is not available, use GitHub private vulnerability reporting or draft
a Security Advisory for this repository when that feature is enabled. If no
private channel is visible, open a public issue that asks for a private
disclosure channel but does not include vulnerability details.

**Do not report security issues through public GitHub issues with details.**

Please include:
- Affected file paths and line numbers (or commit hash)
- Steps to reproduce, if applicable
- Whether the finding is in product code, test fixtures, or local tooling
- Whether you believe the issue affects tracked repository files or local ignored clutter

What to expect:
- Acknowledgment within 5 business days
- Initial assessment within 10 business days
- Public disclosure timeline coordinated with the reporter

## Security Scope

`sdp-trace` is a portable evidence substrate. It records evidence and gaps; it
does not make trust decisions, release approvals, or policy choices. The
security scope covers:

- **Input validation**: Command-line arguments, JSON schema validation, and
  file path handling.
- **Evidence integrity**: Digest verification and provenance tracking.
- **Secret handling**: Redaction of sensitive values in trace artifacts.
- **Runtime integration boundaries**: The core binary is local-first. Optional
  witness, release-proof, PR-review, and repo-observer workflows may depend on
  external CI, Git, or customer authority context.

In scope:
- `cmd/sdp-trace/`
- `internal/`
- `tools/`
- `schema/`
- `docs/` (when it affects security claims or trust boundaries)

Out of scope:
- Security decisions made by policy consumers using sdp-trace output.
- Security of systems that wrap or integrate with sdp-trace.
- External infrastructure where sdp-trace artifacts are stored.
- Local agent run directories (`.worktrees/`, `.codex-subagents/`, `.sdp-trace-runs/`)
- Personal developer environments
- Third-party dependencies (report to the upstream project)

## Trust Boundary

`sdp-trace` is an evidence substrate; it does not approve changes, releases,
or production trust decisions. Security reports should focus on:
- Incorrect evidence claims
- Information disclosure in generated artifacts
- Trust boundary violations in witness or authority flows
- Secret leakage in tracked files

## Current Security Posture

See `docs/production-adoption-readiness.md` for the current adoption readiness
matrix and `docs/security-baseline.md` for the scanner findings ledger.

This repository is in controlled-pilot MVP status. It has not undergone external
security audit. Production deployments should follow your organization's security
review process.

## Security Best Practices for Deployments

If you deploy sdp-trace in your environment:

1. Use release artifacts when available, or build from a reviewed source commit.
2. Verify checksums and signatures when available.
3. Run with minimal privileges; sdp-trace does not require root.
4. Store trace artifacts (`.sdp-trace-runs/`) according to your retention policy.
5. Do not commit raw credentials, tokens, or sensitive source to trace artifacts.
6. Review and redact sensitive data before sharing trace packages externally.
7. Treat local filesystem paths in trace metadata as potentially sensitive.

## Scanner Status

Current CI runs the repository's Go, docs, hygiene, schema, and quality gates.
Additional security scanners are tracked in `docs/security-baseline.md`:

- `go vet`: Static analysis; currently pass locally.
- `govulncheck`: Vulnerability database check; currently pass locally.
- `gosec`: Static security analysis; 133 findings classified in
  `docs/security-baseline.md`.
- `gitleaks`: Secret detection; local and tracked-source scans pass with the
  reviewed `.gitleaks.toml` fixture allowlist.

`govulncheck`, `gosec`, and `gitleaks` CI coverage is `not_assessed` until
those jobs are added to `.github/workflows/ci.yml`.

## Local Verification

To reproduce the current security baseline locally:

```bash
go vet ./...
go test -count=1 ./...
govulncheck ./...
gosec ./...
gitleaks detect --source . --config .gitleaks.toml
```
