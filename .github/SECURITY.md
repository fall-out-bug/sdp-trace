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

1. **Private disclosure**: Email the maintainers directly with a description of
   the issue, steps to reproduce, and potential impact.
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
security audit. Production deployments should follow your organization's
security review process.

## Security Best Practices for Deployments

If you deploy sdp-trace in your environment:

1. Use release artifacts when available, or build from a reviewed source commit.
2. Verify checksums and signatures when available.
3. Run with minimal privileges; sdp-trace does not require root.
4. Store trace artifacts (`.sdp-trace-runs/`) according to your retention policy.
5. Do not commit raw credentials, tokens, or sensitive source to trace artifacts.
6. Review and redact sensitive data before sharing trace packages externally.

## Scanner Status

Automated security scans are run locally and in CI:

- `go vet`: Static analysis; no findings.
- `govulncheck`: Vulnerability database check; no vulnerabilities found.
- `gosec`: Static security analysis; findings are classified in
  `docs/security-baseline.md`.
- `gitleaks`: Secret detection; findings require allowlist review.

See `docs/security-baseline.md` for the full triage ledger and classification
status.
