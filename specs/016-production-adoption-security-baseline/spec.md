# Spec 016: Production Adoption And Security Baseline

Status: draft

## Objective

Define the minimum evidence needed before `sdp-trace` can be presented as
production-adoptable beyond a controlled pilot. This spec does not claim
customer production adoption or external production trust.

## Evidence From 2026-05-20 Probe

- GitHub repository: public, 0 stars, 0 forks, 0 issues, 1 pull request,
  no latest release, security policy disabled.
- Current local branch: `main` at
  `464d5e2c6e4a208861fa50b420afb05e4177144c`.
- Command surface: 27 command families; many are `partial`.
- Roadmap: existing specs are draft, blocked, or in progress; no completed spec
  is recorded.
- `go vet ./...`: pass.
- `govulncheck ./...`: pass, no vulnerabilities found.
- `gosec ./...`: 133 findings, including path traversal/file inclusion,
  subprocess launch, integer conversion, and synthetic-secret detections.
  One `G304` false positive was later given a scoped reviewed suppression;
  the current security baseline records 132 remaining findings.
- `gitleaks` on tracked `HEAD` snapshot: 10 findings, apparently fixture/test
  markers or synthetic tokens, but not allowlisted.
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

## Open Questions

- Should `gosec` become a blocking CI gate, advisory CI artifact, or
  slice-specific review input?
- Should security policy be implemented as `SECURITY.md` at repository root or
  under `.github/SECURITY.md`?
- Should ignored local agent run directories be scanned by hygiene tooling, or
  only reported as local clutter?
