# Security Policy

## Supported Versions

| Version | Supported          | Notes |
| ------- | ------------------ | ----- |
| current | security reports accepted | Current development branch; published support window is not yet established |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please
report it responsibly.

**Do not report security issues through public GitHub issues.**

To report a security concern:

1. **Private disclosure**: Use GitHub private vulnerability reporting or draft
   Security Advisory for this repository when available. If that channel is not
   visible, open a public issue that asks for a private disclosure channel but
   does not include vulnerability details.
2. **Expected response time**: We aim to acknowledge within 7 days.
3. **What to expect**: We will assess the report, determine severity, and work
   with you on a fix or mitigation.
4. **Disclosure**: We request that you give us reasonable time to address the
   issue before any public disclosure.

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

## Out of Scope

- Security decisions made by policy consumers using sdp-trace output.
- Security of systems that wrap or integrate with sdp-trace.
- External infrastructure where sdp-trace artifacts are stored.

## Current Security Posture

See `docs/production-adoption-readiness.md` for the current adoption readiness
matrix and `docs/security-baseline.md` for the scanner findings ledger.

This repository is in controlled-pilot MVP status. It has not undergone external
security audit. A dedicated security email/contact is `not_assessed` until the
maintainers publish one. Production deployments should follow your
organization's security review process.

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
- `gosec`: Static security analysis; 132 findings require triage before it can
  become a CI gate.
- `gitleaks`: Secret detection; local and tracked-source scans pass with the
  reviewed `.gitleaks.toml` fixture allowlist.

`govulncheck`, `gosec`, and `gitleaks` CI coverage is `not_assessed` until
those jobs are added to `.github/workflows/ci.yml`.

See `docs/security-baseline.md` for the full triage ledger and classification
status.
