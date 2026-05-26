# Spec 022: Post-Merge Governance Closure

Status: draft follow-up split from Spec 019.

## Objective

Own the remaining governance decisions from Spec 019 without pretending that
missed pre-merge gates can be retroactively satisfied.

## Background

Spec 019 delivered repo realignment, monitoring, advisory gate readiness, PR
review CI enforcement, and locality cleanup work. Some work was merged before
the intended pre-implementation approval gates were recorded. Spec 019 now
splits that residual governance debt into this follow-up spec.

## Requirements

- FR-022-001: Preserve the historical fact that Spec 019 pre-implementation
  gates were missed.
- FR-022-002: Decide whether already-merged Spec 019 work is accepted as-is,
  rejected, or split into further remediation specs.
- FR-022-003: Define what evidence is required before a future post-merge
  governance closure can move to `complete`.
- FR-022-004: Keep `merge_approval`, `maintainer_approval`, `not_assessed`,
  and `cannot_verify` states explicit.
- FR-022-005: Avoid treating CI success, review evidence, or checked task boxes
  as maintainer approval.

## Non-Goals

- No retroactive pre-implementation approval claim.
- No production trust, release approval, or signed external attestation claim.
- No behavior change to existing commands.
- No merge of already-split governance debt back into Spec 019.

## Acceptance Criteria

- The Spec 019 residual governance state is summarized with exact PR and CI
  references.
- A maintainer decision is recorded for accepting, rejecting, or further
  splitting the residual governance state.
- Follow-up remediation work, if any, is represented as reviewed tasks before
  implementation starts.

