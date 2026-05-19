# Production Adoption Readiness

This document tracks the evidence needed before `sdp-trace` can be presented
as production-adoptable beyond a controlled pilot. It does not claim
production trust, release approval, or customer adoption.

## Trust Scope

`sdp-trace` records evidence and gaps. It does not approve merge, release,
risk acceptance, or production trust. Every row below separates what is
known from what remains `not_assessed` or `cannot_verify`.

## Adoption Readiness Matrix

| Dimension | Current State | Evidence | Open Gap |
| --- | --- | --- | --- |
| **Public repository adoption** | Available as open-source project | GitHub repository, commit history | External customer adoption evidence: `not_assessed` |
| **Pilot capability** | Controlled-pilot MVP | Local trace packages, repo-observable evidence, assessment profiles, witness artifacts, source-bound release checks | Command surface has `partial` families; see `docs/agent-entrypoint.md` |
| **Release readiness** | Local source-bound release proof available | `sdp-trace release-proof` verifies manifest subjects against source commit | Live CI at final head: `not_assessed`; no published releases yet |
| **External production trust** | Not claimed | N/A | Requires external trust profile pass; currently `not_assessed` |
| **Security policy** | Local policy file drafted | `.github/SECURITY.md` | GitHub security policy publication is `not_assessed` until merged and visible on GitHub |
| **Security findings triage** | In progress | `gosec` 132 findings, tracked-source `gitleaks` passes with reviewed fixture allowlist; see `docs/security-baseline.md` | `gosec` findings must be fixed, narrowed, or reviewed before adoption claim improves |

## Scanner Findings Summary

| Tool | Finding Count | Status |
| --- | --- | --- |
| `go vet` | 0 | `pass` |
| `govulncheck` | 0 vulnerabilities | `pass` |
| `gosec` | 132 findings | `needs_triage`; classification ledger in `docs/security-baseline.md` |
| `gitleaks` | 0 findings with reviewed `.gitleaks.toml` | `verified`; default-config fixture hits remain documented in `docs/security-baseline.md` |

See `docs/security-baseline.md` for the full triage ledger.

## Open `not_assessed` Areas

The following areas cannot be claimed from local checks alone:

- **External customer adoption**: No evidence of production use by external customers.
- **Live CI at final head**: CI workflow exists but live run at HEAD is not verified in this document.
- **External security audit**: No external security review has been completed.
- **gosec findings classification**: 132 findings require reviewed disposition; must be labeled blocker, false positive, accepted fix, or deferred advisory before readiness improves.
- **gitleaks CI coverage**: Local and tracked-source scans pass with `.gitleaks.toml`; live CI coverage is `not_assessed` until the job exists.
- **GitHub security policy**: Local `.github/SECURITY.md` exists; GitHub publication is `not_assessed`.

## What This Document Does Not Claim

- No production trust decision has been made.
- No merge or release approval is implied.
- No claim that `gitleaks` or `gosec` exhaustively prove absence of secrets or vulnerabilities.
- No claim that local checks alone support production readiness.
- No external customer adoption evidence.

## Related Documents

- [Adoption Guide](adoption-guide.en.md): implementation model and evidence boundaries.
- [Repository Rollout Playbook](repository-rollout-playbook.en.md): team setup and evidence contracts.
- [Security Baseline](security-baseline.md): scanner triage ledger.
- [Overclaim Checklist](overclaim-checklist.md): forbidden claims and trust-scope rules.

## Verification Commands

```text
go vet ./...
go test -count=1 ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

Additional scanner commands (results recorded separately):

```text
govulncheck ./...
gitleaks detect --source . --redact
gosec ./...
```

## Maintenance

Update this matrix when:

- a new completed spec is recorded in `docs/roadmap.md`;
- scanner findings are classified or resolved;
- a live CI run completes and is recorded;
- external customer adoption evidence becomes available;
- a GitHub security policy is enabled and configured.
