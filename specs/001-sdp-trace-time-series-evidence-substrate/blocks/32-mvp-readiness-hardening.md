# Block 32: MVP Readiness Hardening

**Status**: Implementation complete locally; PR-level review and final-head CI pending
**Spec package**: `specs/004-mvp-readiness-hardening/`
**Branch/worktree**: `codex/mvp-readiness-spec` in `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-mvp-readiness-spec`

## Problem

Repository review found that `sdp-trace` is close to controlled-pilot MVP, but
the MVP handoff claim is not yet supported across documentation freshness,
example clarity, lint, complexity/CRAP, coverage, and CI enforcement.

The review did not find a need to change the product boundary. It did find that
the current boundary must be easier to verify and harder to accidentally
overclaim.

## Required Scope

- Align authoritative docs with live CLI help.
- Restore English/Russian command-surface parity or clearly route to one
  canonical command contract.
- Classify placeholder examples so they are not confused with MVP evidence.
- Fix current lint failures.
- Define and enforce measurable complexity gates with an honest ratchet toward
  `CRAP < 5`.
- Add focused coverage for MVP-critical zero/low-coverage packages.
- Keep external production trust and GitHub CI state explicit as
  `not_assessed` unless live evidence exists.

## Current Intake Evidence

- `go test ./...`: pass in local checkout.
- `jq empty schema/*.json`: pass in local checkout.
- `git diff --check HEAD`: pass in local checkout.
- `golangci-lint run ./...`: fail.
- `gocyclo -over 4 .`: fail with multiple high-complexity production paths.
- Package coverage includes 0.0% packages and very low coverage in `trace`.
- Documentation and code review subagents both returned `REVISE`.
- External pi review via `zai/glm-5.1` returned `REVISE`.

## Closure Boundary

This block can support controlled-pilot MVP handoff only after reviewed spec
approval, implementation, verification, PR-level review, and fresh CI evidence.
It cannot close external production trust.

## Current Implementation State

- Socratic spec review completed across product/docs, code quality/CRAP, and
  trust-boundary axes; focused re-review returned `APPROVE`.
- User approval to implement was received on 2026-05-10.
- Implementation evidence is recorded in
  `specs/004-mvp-readiness-hardening/implementation-ledger.md`.
- Local lint and tests now pass.
- Strict production `CRAP < 5`, cyclomatic complexity `< 15`, and cognitive
  complexity `< 15` now pass locally for `cmd` and `internal`.
- GitHub CI and PR-level review remain `not_assessed` for the final head until
  PR #37 is pushed and live checks/reviews are observed.
