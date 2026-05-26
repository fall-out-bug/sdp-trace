# Spec 022: Post-Merge Governance Closure

Status: active clarification complete; implementation tasks opened.

## Objective

Own the remaining governance decisions from Spec 019 without pretending that
missed pre-merge gates can be retroactively satisfied.

## Background

Spec 019 delivered repo realignment, monitoring, advisory gate readiness, PR
review CI enforcement, and locality cleanup work. Some work was merged before
the intended pre-implementation approval gates were recorded. Spec 019 now
splits that residual governance debt into this follow-up spec.

## Clarifications

### Session 2026-05-26

- Q: Has the Spec 019 residual governance decision already been recorded, or
  does Spec 022 still need to choose accept/reject/split? → A: Already
  `split_successor`.
- Q: What closure bar lets Spec 022 move to `complete`? → A: Evidence cited;
  remaining work is none or reviewed successor specs.
- Q: What evidence refresh is required before Spec 022 can claim `complete`? →
  A: Refresh live PR/CI state, then run local consistency checks; unavailable
  live state stays `not_assessed` or `cannot_verify`.
- Q: Where must Spec 022 closure be recorded? → A:
  `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, and
  `docs/roadmap.md`.
- Q: Should Spec 022 prepared backlog become active tasks now? → A: Convert to
  active unchecked tasks now.

## Requirements

- FR-022-001: Preserve the historical fact that Spec 019 pre-implementation
  gates were missed.
- FR-022-002: Treat the existing maintainer decision `split_successor` as the
  current decision for Spec 019 residual governance; Spec 022 MUST NOT reopen
  accept/reject/split unless a new maintainer decision explicitly supersedes
  it.
- FR-022-003: Define what evidence is required before a future post-merge
  governance closure can move to `complete`.
- FR-022-004: Keep `merge_approval`, `maintainer_approval`, `not_assessed`,
  and `cannot_verify` states explicit.
- FR-022-005: Avoid treating CI success, review evidence, or checked task boxes
  as maintainer approval.
- FR-022-006: Spec 022 MAY move to `complete` only when existing
  `split_successor` evidence is cited and remaining work is either explicitly
  recorded as none or filed as reviewed successor specs before implementation.
- FR-022-007: Before Spec 022 claims `complete`, the worker MUST refresh live
  PR/CI state for the cited PRs when available, then run local consistency
  checks. If live PR/CI state is unavailable, the closure record MUST preserve
  that state as `not_assessed` or `cannot_verify` with a concrete reason.
- FR-022-008: Final Spec 022 closure MUST update
  `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, and
  `docs/roadmap.md` in the same change so the decision surface, reality ledger,
  and navigation surface do not drift.
- FR-022-009: Once Spec 022 is explicitly taken into work, prepared backlog rows
  MUST become active unchecked tasks before closure work starts.

## Non-Goals

- No retroactive pre-implementation approval claim.
- No production trust, release approval, or signed external attestation claim.
- No behavior change to existing commands.
- No merge of already-split governance debt back into Spec 019.

## Acceptance Criteria

- The Spec 019 residual governance state is summarized with exact PR and CI
  references.
- Existing `split_successor` evidence is cited from the closure decision
  ledger, roadmap, Spec 019 plan/tasks, and post-merge closure plan.
- Follow-up remediation work, if any, is represented as reviewed tasks before
  implementation starts.
- If no follow-up remediation remains, the spec records that explicitly rather
  than inferring it from checked task boxes or CI success.
- Live PR/CI evidence is refreshed before any `complete` claim, or the missing
  live state is recorded as `not_assessed` or `cannot_verify`.
- Closure state is recorded consistently in the closure decision ledger, spec
  reality ledger, and roadmap.
- Prepared backlog rows are converted to active unchecked tasks before closure
  work is claimed.
