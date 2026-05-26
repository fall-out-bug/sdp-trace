# sdp-trace Spec Roadmap

> Lightweight navigation for active and historical specs.
> Status is editorial lifecycle guidance, not a trust verdict.
> For multi-axis status rules, see `docs/spec-status-discipline.md`.
> For authoritative claims, see `docs/claim-authoring.md`.

## Legend

| Status | Meaning |
| --- | --- |
| `draft` | Spec written, not yet reviewed or approved for implementation. |
| `pending_review` | Spec in Socratic or adversarial review; approval blocked on findings. |
| `approved` | Review complete, implementation authorized. |
| `in_progress` | Implementation active on a feature branch or worktree. |
| `paused` | Intentionally paused; resume criteria recorded in tasks. |
| `blocked` | Cannot proceed until external dependency or blocker resolved. |
| `complete` | Implementation merged to `main`; maintenance mode only. |
| `historical` | Block records preserved for evidence; not live work. |
| `retired_superseded` | Planning artifact retired because later specs or blocks replaced it; not a product completion claim. |

These labels are navigation shortcuts. They do not replace the separate
`spec_state`, `task_state`, `implementation_state`, `review_state`,
`merge_state`, and `trust_state` axes defined in
[`spec-status-discipline.md`](spec-status-discipline.md).

### Minimal Transitions

These are editorial conventions, not enforced gates:
- `draft` → `pending_review` → `approved` → `in_progress` → `complete`
- `in_progress` may transition to `paused` or `blocked` with a recorded reason.
- `blocked` returns to `in_progress` when the blocker is resolved, not `complete`.
- `complete` is final for a slice; new work opens a new spec.
- `historical` is assigned by maintainers when a spec is archived, not by the original author.
- `retired_superseded` is assigned when a stale planning artifact is closed as
  non-authoritative and mapped to successor specs or blocks.

## Current Reality Snapshot

Current source tree snapshot, based on direct inspection of `specs/*/tasks.md`:

| Measure | Current Value |
| --- | --- |
| Spec directories | 19 |
| SpecKit triplets (`spec.md`, `plan.md`, `tasks.md`) | 19 / 19 |
| Checked task boxes | 585 / 605 |
| Specs with all task boxes checked | 9 |
| Formal `complete` roadmap rows | 0 |

Interpretation: the repository has substantial implemented work. Formal
completion remains open because review, merge, or trust axes are not uniformly
closed or represented. Do not read `draft` as "not implemented"; read the
reality notes and the spec-specific task state.

## Active Specs

| Spec | Capability | Status | Blocker / Next Step |
| --- | --- | --- | --- |
| [015](../specs/015-spec-governance-and-roadmap/) | Spec governance, lifecycle taxonomy, roadmap navigation | `in_progress` | Finalize `docs/roadmap.md`, run multi-LLM review, merge |
| [019](../specs/019-repo-realignment-monitoring-gate-readiness/) | Repo realignment, monitoring, and gate readiness | `in_progress` | PR #60 merged partial implementation to `main` without recorded approval; integration PR #63 superseded PR #62 and PR #31 by combining post-merge closure with PR-review CI enforcement and merged after final-head CI passed. Maintainer approval remains `not_assessed`. See `post-merge-closure-plan.md`. |
| [018](../specs/018-core-policy-split-and-pi-delivery/) | Core/policy split and Pi delivery plan | `in_progress` | Workstreams A-E and integration verification are mapped in `tasks.md` (9 / 11 checked); maintainer review remains `not_assessed`; follow-up implementation specs remain open. |
| [017](../specs/017-oss-replacement-compatibility-and-benchmarks/) | OSS replacement compatibility and benchmarks | `in_progress` | Task ledger is 11 / 11 checked after controlled supply-chain prototype and docs index closure. Optional external tools remain `not_assessed` when absent and supply-chain conformance remains `cannot_verify`; final maintainer closure remains pending. See `docs/spec-reality-ledger.md`. |
| [016](../specs/016-production-adoption-security-baseline/) | Production adoption and security baseline | `in_progress` | PR #59 merged; external audit and production adoption evidence remain `not_assessed`. |
| [014](../specs/014-docs-ux-command-guidance/) | Docs UX, command guidance, profile selection | `draft` | Awaiting Socratic review |
| [013](../specs/013-contributor-onboarding-and-verification/) | Contributor onboarding and verification flow | `in_progress` | Task ledger is 12 / 12 checked after contributor quick-start reconciliation, smoke replay, doccheck coverage, and closure review. Formal maintainer closure remains pending. |
| [012](../specs/012-repo-hygiene-and-artifact-boundary/) | Repository hygiene and artifact boundary rules | `draft` | Awaiting Socratic review |
| [011](../specs/011-schema-docs-generation/) | Schema documentation generation and validation | `draft` | Awaiting Socratic review |
| [010](../specs/010-command-package-organization/) | Command package organization and surface structure | `draft` | Awaiting Socratic review |
| [009](../specs/009-machine-readable-command-surface/) | Machine-readable command surface (JSON schema, registry) | `draft` | Awaiting Socratic review |

## Blocked or Paused Specs

| Spec | Capability | Status | Blocker / Next Step |
| --- | --- | --- | --- |
| [001](../specs/001-sdp-trace-time-series-evidence-substrate/) | Time-series evidence substrate, trace format, data model | `draft` | → Blocked on: self-attestation proof incomplete; external production trust blocked until signed release process. Historical CRAP/MI failures were remediated by later quality work; current baseline quality gates pass, while absolute MI remains an assessed gap. See `docs/spec-reality-ledger.md`. `blocks/` directory preserved as evidence. |
| [008](../specs/008-invisible-flight-recorder/) | Invisible flight recorder (wrap command, session capture) | `in_progress` | → Blocked on: post-implementation review recorded; PR/final-head CI evidence pending. See spec blockers. |

## Implemented But Not Formally Closed

These specs have all current task boxes checked in `tasks.md`, but they are not
formal `complete` rows until review, merge, and trust closure are represented
for the claimed scope.

| Spec | Capability | Task State | Missing Closure Axis |
| --- | --- | --- | --- |
| [008](../specs/008-invisible-flight-recorder/) | Invisible flight recorder | 26 / 26 checked | PR/final-head CI and merge closure not represented |
| [009](../specs/009-machine-readable-command-surface/) | Machine-readable command surface | 14 / 14 checked | PI review / approval not represented |
| [010](../specs/010-command-package-organization/) | Command package organization | 14 / 14 checked | PI review / approval not represented |
| [011](../specs/011-schema-docs-generation/) | Schema documentation validation | 14 / 14 checked | PI review / approval not represented |
| [012](../specs/012-repo-hygiene-and-artifact-boundary/) | Repository hygiene and artifact boundary | 12 / 12 checked | PI review / approval not represented |
| [014](../specs/014-docs-ux-command-guidance/) | Docs UX and command guidance | 15 / 15 checked | PI review / approval not represented |
| [013](../specs/013-contributor-onboarding-and-verification/) | Contributor onboarding and verification flow | 12 / 12 checked | Formal maintainer closure not represented |
| [015](../specs/015-spec-governance-and-roadmap/) | Spec governance and roadmap navigation | 17 / 17 checked | Formal post-merge closure not represented |
| [016](../specs/016-production-adoption-security-baseline/) | Production adoption and security baseline | 10 / 10 checked | External audit, customer adoption, signed release, and production trust remain `not_assessed` |

## Older Draft Specs (No Active Work)

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| [002](../specs/002-authority-envelope-boundary-observation/) | Authority envelope boundary observation | `draft` | Task ledger is 35 / 35 checked after PR-level closure review, explicit merge approval, PR #64 merge as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`, and post-merge `main` CI run `26453881873`. No `blocks/` directory. |
| [003](../specs/003-agent-supply-chain-roadmap/) | Agent supply chain roadmap and product positioning | `retired_superseded` | Roadmap artifact retired as stale planning; task ledger is 42 / 42 checked as superseded, not implemented. No `blocks/` directory. |
| [004](../specs/004-mvp-readiness-hardening/) | MVP readiness hardening criteria | `draft` | Task ledger is 43 / 43 checked after PR-level review, named reviewer sign-off, explicit merge approval, PR #64 merge as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`, and post-merge `main` CI run `26453881873`; absolute MI remains an assessed gap. See `docs/spec-reality-ledger.md`. No `blocks/` directory. |
| [005](../specs/005-product-contract-v0/) | Product contract schema and versioning | `in_progress` | Task ledger is 19 / 20 checked after mapping implementation placeholders to Spec 006 packet artifacts. Historical explicit approval remains `not_assessed`. No `blocks/` directory. |
| [006](../specs/006-change-evidence-packet-core/) | Change evidence packet core format | `in_progress` | Implementation, local verification, Socratic closure review, implementation review, and PR-level evidence are mapped in `tasks.md` (26 / 27 checked). Historical pre-implementation approval remains `not_assessed`. No `blocks/` directory. |
| [007](../specs/007-github-oss-demo-packet/) | GitHub OSS demo packet workflow | `draft` | Task ledger is 21 / 22 checked after external demo evidence review, v1 baseline tag, and buyer rehearsal. Explicit demo-track approval remains open. No `blocks/` directory. |

> **Note**: These specs remain in `draft` per their own `spec.md` files. The roadmap does not override spec source-of-truth status. They are listed here separately because no active work is in progress. When work resumes, move to Active Specs. Upon implementation, review, merge, and trust closure, move to Formally Closed Specs.

## Formally Closed Specs

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| *(none yet)* | | | |

## Historical / Archived Evidence

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| *(none yet — see Blocked or Paused for specs with preserved `blocks/` directories)* | | | |

## Capability Index

Use this to find which spec owns a product surface. A capability may be touched by multiple specs; the listed owner is the primary spec, not an exclusive boundary.

| Capability | Owner Spec(s) | Current Repository Reality |
| --- | --- | --- |
| Evidence substrate / trace format | 001 | Largely implemented; Block 21 and Block 32 closure refreshes recorded, while demo/first-run work remains open |
| Authority envelope / trust boundary | 002 | PR-level review complete; merge/post-merge closure open |
| Product contract schema | 005 | Contract and implementation placeholders are artifact-complete via Spec 006; approval remains `not_assessed` |
| Change evidence packet | 006 | Implemented and reviewed locally; historical pre-implementation approval remains `not_assessed` |
| GitHub demo workflow | 007 | Spec reviewed; approval and demo-repository implementation remain open |
| Flight recorder / wrap command | 008 | Implemented; formal closure open |
| Command surface (JSON schema, registry) | 009 | Implemented; review closure open |
| Command package organization | 010 | Implemented; review closure open |
| Schema docs generation | 011 | Implemented; review closure open |
| Repo hygiene / artifact boundary | 012 | Implemented; review closure open |
| Contributor onboarding | 013 | Implemented local / reviewed; formal maintainer closure remains pending |
| Docs UX / command guidance | 014 | Implemented; review closure open |
| Spec governance / roadmap | 015 | Implemented; formal closure open |
| Production adoption / security baseline | 016 | Implemented locally; production trust remains `not_assessed` |
| OSS replacement compatibility / benchmarks | 017 | Implemented local / trust-bounded; optional external tools and external supply-chain trust remain `not_assessed` / `cannot_verify` |
| Core/policy split and Pi delivery | 018 | Workstreams A-E implemented locally; maintainer review and follow-up-spec decision remain open |
| Repo realignment / monitoring / gate readiness | 019 | Partial; HITL and closure blockers remain |

## Claim-Tag Enforcement Scope

- **Required**: New authoritative claims in any file created or materially modified after spec 015.
- **Exempt**: Historical specs (001–014), their `blocks/` directories, and existing checked-in review JSON.
- **Rule**: See `docs/claim-authoring.md` for tag grammar and allowed values.

## Task-File Expectations

Every spec directory must contain `spec.md`, `plan.md`, and `tasks.md`.

### Required in `tasks.md`

- Phase grouping (e.g., Phase 0 – Review, Phase 1 – Implementation).
- Explicit blocker notation for any task that cannot proceed.
- Approval gate markers when a phase requires review before continuation.
- Checkbox state that matches the current repository reality.

### Blocker Notation

Use this format for blocked tasks:

```markdown
- [ ] T### Description → Blocked on: <concrete reason>.
```

Example:

```markdown
- [ ] T020 Run external integration tests → Blocked on: OIDC signing environment not yet provisioned.
```

### Approval Gates

When a phase requires explicit approval before the next phase begins:

```markdown
- [ ] T### Approval gate: <what must be approved> → Awaiting: <who>.
```

### Historical Exemption

Historical block records in `blocks/` are exempt from these expectations. They are preserved as evidence, not live planning artifacts.

## Roadmap Freshness

- Update this file when a new spec is opened or an active spec's status changes.
- Owner: the spec author or current block worker.
- Last updated: 2026-05-26.

<!-- sdp-trace-claim: claim=profile_passed; subject=roadmap-001-015-coverage; state=pass; profile=repo_baseline_structural; evidence=command_set:block015-t030 -->
