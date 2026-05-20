# sdp-trace Spec Roadmap

> Lightweight navigation for active and historical specs.
> Status is editorial lifecycle guidance, not a trust verdict.
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

### Minimal Transitions

These are editorial conventions, not enforced gates:
- `draft` → `pending_review` → `approved` → `in_progress` → `complete`
- `in_progress` may transition to `paused` or `blocked` with a recorded reason.
- `blocked` returns to `in_progress` when the blocker is resolved, not `complete`.
- `complete` is final for a slice; new work opens a new spec.
- `historical` is assigned by maintainers when a spec is archived, not by the original author.

## Active Specs

| Spec | Capability | Status | Blocker / Next Step |
| --- | --- | --- | --- |
| [015](../specs/015-spec-governance-and-roadmap/) | Spec governance, lifecycle taxonomy, roadmap navigation | `in_progress` | Finalize `docs/roadmap.md`, run multi-LLM review, merge |
| [018](../specs/018-core-policy-split-and-pi-delivery/) | Core/policy split and Pi delivery plan | `draft` | Awaiting review |
| [017](../specs/017-oss-replacement-compatibility-and-benchmarks/) | OSS replacement compatibility and benchmarks | `in_progress` | Slice review in progress; all workstreams implemented; final PR review pending |
| [016](../specs/016-production-adoption-security-baseline/) | Production adoption and security baseline | `in_progress` | PR #59 open; repeat review in progress |
| [014](../specs/014-docs-ux-command-guidance/) | Docs UX, command guidance, profile selection | `draft` | Awaiting Socratic review |
| [013](../specs/013-contributor-onboarding-and-verification/) | Contributor onboarding and verification flow | `draft` | Awaiting Socratic review |
| [012](../specs/012-repo-hygiene-and-artifact-boundary/) | Repository hygiene and artifact boundary rules | `draft` | Awaiting Socratic review |
| [011](../specs/011-schema-docs-generation/) | Schema documentation generation and validation | `draft` | Awaiting Socratic review |
| [010](../specs/010-command-package-organization/) | Command package organization and surface structure | `draft` | Awaiting Socratic review |
| [009](../specs/009-machine-readable-command-surface/) | Machine-readable command surface (JSON schema, registry) | `draft` | Awaiting Socratic review |

## Blocked or Paused Specs

| Spec | Capability | Status | Blocker / Next Step |
| --- | --- | --- | --- |
| [001](../specs/001-sdp-trace-time-series-evidence-substrate/) | Time-series evidence substrate, trace format, data model | `draft` | → Blocked on: self-attestation proof incomplete; external production trust blocked until signed release process. See spec tasks T020–T026. `blocks/` directory preserved as evidence. |
| [008](../specs/008-invisible-flight-recorder/) | Invisible flight recorder (wrap command, session capture) | `in_progress` | → Blocked on: post-implementation review recorded; PR/final-head CI evidence pending. See spec blockers. |

## Older Draft Specs (No Active Work)

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| [002](../specs/002-authority-envelope-boundary-observation/) | Authority envelope boundary observation | `draft` | Revised after initial Socratic review; focused re-review pending. No `blocks/` directory. |
| [003](../specs/003-agent-supply-chain-roadmap/) | Agent supply chain roadmap and product positioning | `draft` | Roadmap artifact; Socratic review completed; revisions pending. No `blocks/` directory. |
| [004](../specs/004-mvp-readiness-hardening/) | MVP readiness hardening criteria | `draft` | Revised after initial Socratic review; approval pending. No `blocks/` directory. |
| [005](../specs/005-product-contract-v0/) | Product contract schema and versioning | `draft` | Revised after full review; re-review pending. No `blocks/` directory. |
| [006](../specs/006-change-evidence-packet-core/) | Change evidence packet core format | `draft` | Needs Socratic review before implementation approval. No `blocks/` directory. |
| [007](../specs/007-github-oss-demo-packet/) | GitHub OSS demo packet workflow | `draft` | Needs Socratic review before implementation approval. No `blocks/` directory. |

> **Note**: These specs remain in `draft` per their own `spec.md` files. The roadmap does not override spec source-of-truth status. They are listed here separately because no active work is in progress. When work resumes, move to Active Specs. Upon completion and merge, move to Completed Specs.

## Completed Specs

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| *(none yet)* | | | |

## Historical / Archived Evidence

| Spec | Capability | Status | Notes |
| --- | --- | --- | --- |
| *(none yet — see Blocked or Paused for specs with preserved `blocks/` directories)* | | | |

## Capability Index

Use this to find which spec owns a product surface. A capability may be touched by multiple specs; the listed owner is the primary spec, not an exclusive boundary.

| Capability | Owner Spec(s) | Live? |
| --- | --- | --- |
| Evidence substrate / trace format | 001 | Draft (blocked) |
| Authority envelope / trust boundary | 002 | Draft (old) |
| Product contract schema | 005 | Draft (old) |
| Change evidence packet | 006 | Draft (old) |
| GitHub demo workflow | 007 | Draft (old) |
| Flight recorder / wrap command | 008 | In progress (blocked) |
| Command surface (JSON schema, registry) | 009 | Draft |
| Command package organization | 010 | Draft |
| Schema docs generation | 011 | Draft |
| Repo hygiene / artifact boundary | 012 | Draft |
| Contributor onboarding | 013 | Draft |
| Docs UX / command guidance | 014 | Draft |
| Spec governance / roadmap | 015 | In progress |
| Production adoption / security baseline | 016 | In progress |
| OSS replacement compatibility / benchmarks | 017 | In progress |
| Core/policy split and Pi delivery | 018 | Draft |

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
- Last updated: 2026-05-20.

<!-- sdp-trace-claim: claim=profile_passed; subject=roadmap-001-015-coverage; state=pass; profile=repo_baseline_structural; evidence=command_set:block015-t030 -->
