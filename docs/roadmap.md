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

## Current Reality Snapshot

Current source tree snapshot, based on direct inspection of `specs/*/tasks.md`:

| Measure | Current Value |
| --- | --- |
| Spec directories | 19 |
| SpecKit triplets (`spec.md`, `plan.md`, `tasks.md`) | 19 / 19 |
| Checked task boxes | 500 / 605 |
| Specs with all task boxes checked | 8 |
| Formal `complete` roadmap rows | 0 |

Interpretation: the repository has substantial implemented work. Formal
completion remains open because review, merge, or trust axes are not uniformly
closed or represented. Do not read `draft` as "not implemented"; read the
reality notes and the spec-specific task state.

## Active Specs

| Spec | Capability | Status | Blocker / Next Step |
| --- | --- | --- | --- |
| [015](../specs/015-spec-governance-and-roadmap/) | Spec governance, lifecycle taxonomy, roadmap navigation | `in_progress` | Finalize `docs/roadmap.md`, run multi-LLM review, merge |
| [019](../specs/019-repo-realignment-monitoring-gate-readiness/) | Repo realignment, monitoring, and gate readiness | `in_progress` | PR #60 merged partial implementation to `main` without recorded approval; integration PR #63 now supersedes PR #62 and PR #31 by combining post-merge closure with PR-review CI enforcement. PR #63 final-head CI is `not_assessed` in checked-in docs; maintainer approval remains `not_assessed`. See `post-merge-closure-plan.md`. |
| [018](../specs/018-core-policy-split-and-pi-delivery/) | Core/policy split and Pi delivery plan | `in_progress` | Workstreams A-E and integration verification are mapped in `tasks.md` (9 / 11 checked); maintainer review remains `not_assessed`; follow-up implementation specs remain open. |
| [017](../specs/017-oss-replacement-compatibility-and-benchmarks/) | OSS replacement compatibility and benchmarks | `in_progress` | Slice review in progress; workstreams A–C and E implemented; WS-017-D automated probes remain open (manual-only). **CRAP and MI gates PASS** for `tools/ossbench` and `tools/osscompat` after WS-019-B refactor. Live `wrap` manifest now has `schema/run-manifest.schema.json`; richer flight-recorder profile schema remains separate. See `docs/spec-reality-ledger.md`. Final PR review pending. |
| [016](../specs/016-production-adoption-security-baseline/) | Production adoption and security baseline | `in_progress` | PR #59 merged; external audit and production adoption evidence remain `not_assessed`. |
| [014](../specs/014-docs-ux-command-guidance/) | Docs UX, command guidance, profile selection | `draft` | Awaiting Socratic review |
| [013](../specs/013-contributor-onboarding-and-verification/) | Contributor onboarding and verification flow | `draft` | Awaiting Socratic review |
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
| [015](../specs/015-spec-governance-and-roadmap/) | Spec governance and roadmap navigation | 17 / 17 checked | Formal post-merge closure not represented |
| [016](../specs/016-production-adoption-security-baseline/) | Production adoption and security baseline | 10 / 10 checked | External audit, customer adoption, signed release, and production trust remain `not_assessed` |

## Older Draft Specs (No Active Work)

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| [002](../specs/002-authority-envelope-boundary-observation/) | Authority envelope boundary observation | `draft` | Revised after initial Socratic review; focused re-review pending. No `blocks/` directory. |
| [003](../specs/003-agent-supply-chain-roadmap/) | Agent supply chain roadmap and product positioning | `draft` | Roadmap artifact; Socratic review completed; revisions pending. No `blocks/` directory. |
| [004](../specs/004-mvp-readiness-hardening/) | MVP readiness hardening criteria | `draft` | Revised after initial Socratic review; approval pending. Current baseline CRAP/MI gates pass, but absolute MI remains an assessed gap and the spec is not approved. See `docs/spec-reality-ledger.md`. No `blocks/` directory. |
| [005](../specs/005-product-contract-v0/) | Product contract schema and versioning | `draft` | Revised after full review; re-review pending. No `blocks/` directory. |
| [006](../specs/006-change-evidence-packet-core/) | Change evidence packet core format | `in_progress` | Implementation and local verification artifacts are mapped in `tasks.md` (21 / 27 checked). Spec review, approval, implementation-review, and PR-level review evidence remain open / `not_assessed`. No `blocks/` directory. |
| [007](../specs/007-github-oss-demo-packet/) | GitHub OSS demo packet workflow | `draft` | Needs Socratic review before implementation approval. No `blocks/` directory. |

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
| Evidence substrate / trace format | 001 | Largely implemented but blocked on explicit trust closure and open tasks |
| Authority envelope / trust boundary | 002 | Mostly checked; re-review / closure open |
| Product contract schema | 005 | Partial |
| Change evidence packet | 006 | Implemented locally; task ledger now maps product artifacts, but review/approval/fresh verification remain open |
| GitHub demo workflow | 007 | Partial planning / demo work |
| Flight recorder / wrap command | 008 | Implemented; formal closure open |
| Command surface (JSON schema, registry) | 009 | Implemented; review closure open |
| Command package organization | 010 | Implemented; review closure open |
| Schema docs generation | 011 | Implemented; review closure open |
| Repo hygiene / artifact boundary | 012 | Implemented; review closure open |
| Contributor onboarding | 013 | Spec-only / not implemented |
| Docs UX / command guidance | 014 | Implemented; review closure open |
| Spec governance / roadmap | 015 | Implemented; formal closure open |
| Production adoption / security baseline | 016 | Implemented locally; production trust remains `not_assessed` |
| OSS replacement compatibility / benchmarks | 017 | Partial; automated external probes remain open / `not_assessed` |
| Core/policy split and Pi delivery | 018 | Workstreams A-E implemented locally; maintainer review and integration verification remain open |
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
