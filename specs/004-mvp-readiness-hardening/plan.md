# Implementation Plan: MVP Readiness Hardening

**Branch**: `codex/mvp-readiness-spec` | **Date**: 2026-05-10 | **Spec**: [spec.md](spec.md)
**Input**: Repository review findings across documentation, code quality, CRAP risk, and pi-review.

## Summary

Harden `sdp-trace` for a controlled-pilot MVP handoff by turning review findings
into verifiable docs and quality gates. The work must keep the product promise
honest: `sdp-trace` remains an evidence substrate, not a release gate or
production trust authority.

## Technical Context

**Language/Version**: Go CLI, Markdown, JSON Schema Draft 2020-12
**Primary Dependencies**: Go standard library; existing CLI help; `jq`; local
linters available in current environment (`golangci-lint`, `gocyclo`,
`gocognit`)
**Testing**: `go test -count=1 ./...`, `go test -count=1 ./... -coverprofile`,
per-function CRAP baseline/review, `jq empty schema/*.json`, `git diff --check`,
`golangci-lint run ./...`, selected complexity gate, doc command-surface check
**Target Platform**: Portable repository workflows and GitHub Actions
**Project Type**: Documentation, tests, small Go refactors, CI gate hardening
**Constraints**: No Node.js/npm/JavaScript active product tooling; no dependency
on Beads, `sdp_lab`, Operator Mode, agentloop, or a specific harness runtime

## Constitution Check

| Rule | Status | Evidence |
|---|---|---|
| Spec before implementation | Pass | This package is draft spec/plan/tasks only. |
| UX first | Pass | Entrypoint docs and examples are treated as user-facing product surface. |
| DX second | Pass | Quality gates make local and CI checks reproducible. |
| Evidence-backed claims | Pass | Current evidence distinguishes local pass, lint fail, coverage gaps, and complexity hotspots. |
| Preserve `not_assessed` and `cannot_verify` | Pass | CI enforcement gaps use `assessed_gap`; external trust remains explicit if not assessed. |
| No broad production trust claim | Pass | Spec keeps controlled-pilot boundary. |

## Project Structure

```text
specs/004-mvp-readiness-hardening/
├── spec.md
├── plan.md
├── tasks.md
└── socratic-review.md
```

Expected implementation surfaces after approval:

```text
docs/
├── agent-entrypoint.md
├── reviewer-entrypoint.md
├── adoption-guide.en.md
├── adoption-guide.ru.md
├── repository-rollout-playbook.en.md
├── repository-rollout-playbook.ru.md
└── README.md

examples/
├── codex/README.md
├── claude-code/README.md
├── opencode/README.md
└── go-service/README.md

internal/
├── authority/
├── telemetry/
├── posture/
├── harnessobs/
├── trace/
├── contract/
├── policy/
└── export/

.github/workflows/ci.yml
schema/README.md
```

## Design Decisions To Review

| Decision | Rationale | Review risk |
|---|---|---|
| Treat docs freshness as MVP scope | Stale command docs break trust as directly as stale proof JSON. | Could expand into broad docs rewrite; keep only entrypoints/adoption/examples. |
| Keep Russian docs equivalent or explicitly deferred | Partial bilingual parity is worse than English-only authority. | More translation work; avoid inventing new product terms. |
| Compute CRAP instead of using cyclo as a proxy | Existing code cannot reach CRAP<5 immediately without broad refactor. | Must not water down the user's requested bar; record baseline and ratchet plan. |
| Prioritize zero-coverage trust packages | Contract/policy/export/trace are credibility-critical. | Tests may expose stale or dead code; remove or scope honestly. |
| Add CI enforcement only after local commands are stable | Prevent noisy or environment-specific CI failures. | If not added, CI state must remain `not_assessed`. |

## Phase 0 Exit Criteria

Phase 0 is complete only when:

- independent Socratic review produces written findings with
  critical/major/minor classification;
- all critical findings are resolved in the spec or recorded as blockers;
- all major findings are resolved in the spec or converted to explicit tracked
  implementation follow-ups with `not_assessed` state;
- `socratic-review.md` records reviewer, verdict, findings, and dispositions;
- user approval is explicit and not inferred from silence.

## Phases

### Phase 0: Spec Review Gate

Run Socratic spec review on this package before implementation. Resolve or
record all critical/major findings, then stop for explicit user approval.

### Phase 1: Documentation Freshness

Align `docs/agent-entrypoint.md`, `docs/reviewer-entrypoint.md`, README, and
adoption docs with the live CLI surface and controlled-pilot boundary.

### Phase 2: Example Surface Cleanup

Classify first-class examples as real fixtures, pilot evidence packages,
placeholders, retired examples, `not_assessed`, or `cannot_verify`. Remove or
demote anything that reads like proof but is only a future example.

### Phase 3: Gate Baseline, Lint, And Small Hygiene Fixes

Define the initial CRAP/complexity/coverage gate and record baseline before
decomposition. Fix current `golangci-lint` findings in authority and telemetry.
If existing test coverage for a changed package is below the selected floor, the
fix must add at least one focused test for the changed path or record why that
path remains outside the MVP claim surface.

### Phase 4: Complexity And CRAP Hardening

Decompose the highest-risk functions against the selected gate. Record CRAP and
complexity deltas before claiming improvement.

### Phase 5: Coverage Hardening

Add focused tests for zero/low-coverage MVP-critical packages and changed
trust-sensitive paths. Record package-level and function-level coverage deltas.

### Phase 6: CI And Verification

Add or document quality gates. Run local verification and record any GitHub CI
state as `not_assessed` only when outside the selected CI verification scope;
missing CI enforcement for selected gates is `assessed_gap`.

## Acceptance Gate Before Implementation

Implementation must not start until:

- `spec.md`, `plan.md`, and `tasks.md` exist;
- Socratic spec review has usable independent reviewer output;
- critical/major spec findings are fixed or recorded as blockers;
- the user approves the reviewed spec direction.
