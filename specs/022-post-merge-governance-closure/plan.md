# Implementation Plan: Post-Merge Governance Closure

**Branch**: `codex/022-post-merge-governance-closure` |
**Date**: 2026-05-26 |
**Spec**: [spec.md](spec.md)

**Input**: Feature specification from
`specs/022-post-merge-governance-closure/spec.md`

**Setup Note**: Initial planning happened on `codex/install-github-speckit`,
where `.specify/scripts/bash/setup-plan.sh --json` could not resolve an active
feature branch. Implementation now runs in the isolated worktree branch
`codex/022-post-merge-governance-closure`; `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks`
resolves this feature to `specs/022-post-merge-governance-closure/`.

## Summary

Spec 022 closes the residual governance decision split from Spec 019 without
retroactively approving the missed PR #60 pre-implementation gates. The work is
docs/governance-only: refresh live PR/CI evidence where available, cite the
existing `split_successor` decision, decide whether residual remediation is
none or filed as reviewed successor specs, then synchronize the closure decision
ledger, spec reality ledger, and roadmap.

## Technical Context

**Language/Version**: Docs-only change. Go 1.22 repository checks remain the
verification baseline; no Go product code is planned.

**Primary Dependencies**: Git, GitHub CLI for live PR/CI refresh when
available, Go repository tools already present in `tools/`.

**Storage**: Repository Markdown files under `specs/` and `docs/`; no new
database, generated JSON, or persistent runtime store.

**Testing**: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, and
`git diff --check` are required for the planned docs changes. `go test ./...`
and `go vet ./...` remain repository default checks; for this docs-only slice
they may be recorded as `not_assessed` only if no product code, schema, command,
or executable examples change.

**Target Platform**: Repository-local Spec Kit governance docs consumed by
agents, maintainers, and reviewers.

**Project Type**: Go CLI/tooling plus JSON schemas, Markdown docs, and portable
examples; this feature touches only Markdown governance artifacts.

**Performance Goals**: Not applicable; no runtime path changes.

**Constraints**: Harness-independent evidence recorder, not an approval
authority; evidence-backed or explicitly `not_assessed` / `cannot_verify`; no
opaque health scores; no Node/npm/JS/TS in active product path.

**Scale/Scope**: One governance closure spec plus three navigation/decision
surfaces: `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, and
`docs/roadmap.md`.

## Constitution Check

*GATE: Passed before Phase 0 research. Re-check after Phase 1 design.*

- Portable evidence recorder, not authority: pass. Spec 022 records governance
  closure state and explicitly does not approve merge, release, production
  trust, or risk acceptance.
- Evidence state honesty: pass. Missing or unavailable live PR/CI state remains
  `not_assessed` or `cannot_verify`.
- SpecKit-compatible trace flow: pass. Spec, plan, tasks, evidence, decision,
  trace, and provenance surfaces are named and updated in repository-local
  artifacts.
- Go-first, small product path: pass. No product code, dependency, schema, or
  mixed-toolchain change is planned.
- Separate review, verification, and approval: pass. The plan cites the existing
  `split_successor` decision without converting PR #60 CI/review evidence into
  retroactive approval.
- Authority and scope constraints: pass. Closure state is scoped to repository
  governance and does not claim production trust.

## Project Structure

### Documentation (this feature)

```text
specs/022-post-merge-governance-closure/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── spec.md
└── tasks.md
```

### Source Code (repository root)

```text
docs/
├── closure-decision-ledger.md
├── spec-reality-ledger.md
└── roadmap.md

specs/
├── 019-repo-realignment-monitoring-gate-readiness/
└── 022-post-merge-governance-closure/
```

**Structure Decision**: Keep closure planning in Spec 022 and update the three
existing governance surfaces. Do not introduce a new contract directory because
no public command, schema, or API contract changes.

## Phase 0: Research

See [research.md](research.md).

## Phase 1: Design

See [data-model.md](data-model.md) and [quickstart.md](quickstart.md).

No `/contracts` artifacts are generated for this feature because the planned
change is repository governance documentation only.

## Workstreams

### WS-022-A: Governance Evidence Summary

Owned files:

- `specs/019-repo-realignment-monitoring-gate-readiness/`
- `docs/spec-reality-ledger.md`
- `docs/closure-decision-ledger.md`
- `docs/roadmap.md`

Deliverable:

- Summarize PR #60, PR #63, and Spec 019 review/CI evidence without converting
  it into approval.
- Refresh live PR/CI state for the cited PRs when available; if unavailable,
  preserve the missing live evidence as `not_assessed` or `cannot_verify`.
- In this implementation worktree GitHub CLI authentication is available, so
  PR #60 and PR #63 live state must be refreshed before any `complete` claim.
- Keep closure decision, reality, and roadmap surfaces synchronized.

### WS-022-B: Existing Maintainer Decision

Owned files:

- `specs/022-post-merge-governance-closure/`
- `docs/closure-decision-ledger.md`

Deliverable:

- Cite the existing `split_successor` maintainer decision and preserve it as
  the current Spec 019 residual-governance outcome unless a new maintainer
  decision explicitly supersedes it.

### WS-022-C: Remediation Disposition

Owned files:

- successor specs only if needed
- `docs/spec-reality-ledger.md`
- `docs/roadmap.md`

Deliverable:

- Create reviewed implementation tasks only for residual work that remains
  after applying the existing `split_successor` decision.
- If no residual work remains, record that state explicitly with the evidence
  cited by WS-022-A and WS-022-B.
- If successor specs are required, their `spec.md`, `plan.md`, and `tasks.md`
  must pass the same pre-implementation review loop before any successor
  implementation can be treated as reviewed.

## Post-Design Constitution Check

- Portable evidence recorder, not authority: pass. Designed artifacts only
  record closure evidence and missing states.
- Evidence state honesty: pass. The design requires live refresh or explicit
  `not_assessed` / `cannot_verify`.
- SpecKit-compatible trace flow: pass. Tasks map directly to closure evidence,
  decision, and synchronized docs.
- Go-first, small product path: pass. No new dependencies or product code.
- Separate review, verification, and approval: pass. Existing missed gates stay
  historical evidence; no retroactive approval claim is introduced.
- Authority and scope constraints: pass. Scope remains repository governance.

## Complexity Tracking

No constitution violations.

## Verification

```text
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

Run broader checks only if implementation expands beyond docs:

```text
go test ./...
go vet ./...
go run ./tools/schemadoc
```

## Review Cadence

Before documentation implementation starts, review this spec/plan/task set with
independent model lanes and fix all retained findings. After each user story,
run a focused review of the changed surfaces before proceeding to the next
story. Before PR-ready claims, run a full diff review plus spec, constitution,
product, quality, Clean Architecture, Clean Code, SOLID, DRY, and YAGNI drift
review. Missing, failed, or off-task review lanes remain `cannot_verify` or
`not_assessed` and cannot be counted as LGTM.

Review artifacts must record model, provider, harness, date, prompt class,
timeout, retries, fallback, and whether the lane was exact, unavailable, or
substituted. External review lanes must run only after the reviewed changes are
committed or in an isolated copy, because review harnesses may reset local
uncommitted changes.
