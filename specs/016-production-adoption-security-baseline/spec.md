# Spec 016: Production Adoption And Security Baseline

Status: in_progress

## Objective

Define the minimum evidence needed before `sdp-trace` can be presented as
production-adoptable beyond a controlled pilot. This spec does not claim
customer production adoption or external production trust.

## Evidence From 2026-05-20 Probe

- GitHub repository: public, 0 stars, 0 forks, 0 issues, 1 pull request,
  no latest release, security policy disabled.
- Current local branch: `main` at `464d5e2` (base), rebased onto `18cd4a0`.
- Command surface: 18 command families; 6 are `complete`, 11 are `partial`,
  1 is `not_assessed`.
- Roadmap: existing specs are draft, blocked, or in progress; no completed spec
  is recorded.
- `go vet ./...`: pass.
- `govulncheck ./...`: pass, no vulnerabilities found.
- `gosec ./...`: 132 findings (133 originally; one `G304` in `internal/capturedepth`
  has a scoped reviewed `#nosec` suppression), including path traversal/file inclusion,
  subprocess launch, integer conversion, and synthetic-secret detections.
- `gitleaks` on tracked `HEAD` snapshot: 10 findings with default config;
  0 findings with reviewed `.gitleaks.toml` allowlist.
- External customer adoption evidence: not_assessed.
- Live CI at final head: not_assessed.

## User Stories

1. As a security reviewer, I can see which findings are accepted, false
   positive, advisory, or deferred, with file/line evidence and verification.
2. As a production adopter, I can tell whether a command family is stable,
   partial, fixture-only, or not production-ready.
3. As a maintainer, I can run a reproducible local security baseline without
   scanning ignored local agent clutter as product evidence.
4. As a release owner, I can see why current artifacts do or do not support a
   release/adoption claim.

## Requirements

- FR-016-001: Add a production adoption readiness matrix that separates
  public repository adoption facts, pilot capability, release readiness, and
  external customer adoption.
- FR-016-002: Add a security triage ledger for `gosec`, `gitleaks`,
  `govulncheck`, and `go vet` results.
- FR-016-003: Treat secrets found only in local ignored clutter as local hygiene
  findings, not repository proof; tracked findings must be triaged.
- FR-016-004: Add or update a GitHub security policy path or explicitly record
  why it remains `not_assessed`.
- FR-016-005: Classify each `gosec` finding family as blocker, accepted false
  positive, accepted fix, or deferred advisory.
- FR-016-006: Add tests or fixture allowlists for intentional synthetic secret
  markers without hiding real secret leaks.
- FR-016-007: Do not claim production adoption, production trust, customer
  trust, or release readiness from local checks alone.

## Decisions

- `gosec` is advisory for this phase, not a blocking CI gate. The current
  132 findings remain classified until a later security slice fixes,
  narrows, or suppresses individual call sites with evidence.
- `govulncheck` and `go vet` are required local verification commands.
  `gitleaks` with the reviewed repository config is the selected secret scan
  for local and tracked-source checks.
- The security policy lives under `.github/SECURITY.md`. A root
  `SECURITY.md` is not required for this phase.
- Tracked synthetic secret fixtures are handled by narrow, reviewed
  path-and-regex allowlists. Broad scanner suppression is not allowed.
- Ignored local agent run directories are local clutter, not repository proof.
  They may be scanned during working-tree hygiene, but tracked-source scans
  are the evidence for repository claims.
- No readiness state may improve beyond controlled pilot while external
  customer adoption, live CI at final head, and external security audit remain
  `not_assessed`.

## Non-Goals

- No production trust decision.
- No merge/release approval.
- No claim that gitleaks/gosec exhaustively prove absence of secrets or
  vulnerabilities.
- No broad rewrite of the command surface.

## Acceptance Criteria

- A checked-in readiness/security matrix exists and names every
  `not_assessed` production-adoption area.
- Security findings have dispositions with file/line evidence.
- Tracked synthetic secret fixtures are either made scanner-safe or covered by
  a reviewed allowlist.
- `go vet ./...`, `govulncheck ./...`, and the selected secret scan complete in
  a documented local command path.
- Any remaining `gosec` findings are explicitly classified before readiness is
  improved beyond controlled pilot.

## Resolved Questions

- `gosec`: advisory for this phase; blocker if new G101/G204 appear after baseline.
- Security policy: `.github/SECURITY.md` per GitHub community standard.
- Ignored local agent run directories: scanned by hygiene tooling locally,
  reported as local clutter, not repository proof.

## Change Log

- 2026-05-20: Spec intake completed. All workstreams implemented.
  Readiness matrix, security baseline, gitleaks allowlist, and security
  policy created. All tasks complete; remaining `not_assessed` areas
  explicitly recorded.
