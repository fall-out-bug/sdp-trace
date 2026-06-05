# PR 73 Scope Decision

Date: 2026-06-05

## Decision

Keep PR 73 as the current active cleanup PR, with Spec 021 and Spec 023 both
represented explicitly in SpecKit artifacts and review evidence, instead of
rewriting history into multiple new PRs after implementation.

## Why This Is Not Hidden Scope

Spec 023 was created after Spec 021 exposed remaining active numbered Go source
files outside the original command-surface slice set. It is a separate SpecKit
artifact because the work was larger than a continuation note inside Spec 021.
The follow-on work was driven by explicit user continuation requests during the
same active repository goal: remove numbered source-file debt from active Go
product paths.

## Rejected Alternatives

- Split completed Spec 023 commits into a new PR after the fact. Rejected now:
  high integration risk, high evidence churn, and no improvement to already
  recorded per-slice reviewability.
- Collapse Spec 023 back into Spec 021. Rejected: that would hide the scope
  increase and make the review trail less honest.
- Treat the aggregate PR as a production trust or release approval. Rejected:
  merge approval, release readiness, and external attestation remain separate
  states.

## Boundary

This decision does not weaken FR-021-001 or FR-023-001 for future work. Future
numbered-file cleanup should open a separate PR before implementation when a new
SpecKit artifact expands the touched package set beyond the current spec.

Spec 021 also expanded internally from the original 11 command-surface slices to
the complete package-by-package slice registry in `plan.md`. This is accepted as
part of the same PR 73 exception because the plan and per-slice evidence files
record the exact touched package or command family for each slice. Future work
should not rely on `plan.md` as a late substitute for an updated spec.

PR 73 is accepted as an exception only because the work already exists as
bounded slices with separate spec, plan, task, evidence, review, and verification
artifacts, and because the current user explicitly requested continuing through
PR review to merge.

Review remediation is also bounded to PR 73: security findings discovered during
OpenCode PR review may add small behavior changes when the alternative is
merging with a known trust/security defect. Those changes must be explicitly
named in the spec and final audit instead of being described as locality-only
file moves.

## Risk If Wrong

- Reviewers may still judge the aggregate PR too large to review completely.
- GitHub diff review is less ergonomic because many moves appear as delete/add
  pairs.
- Future agents may copy the aggregate-PR shape unless this exception is kept
  explicit.

Mitigation: keep the exception local to PR 73, preserve the per-slice evidence,
and require OpenCode PR review plus live checks before merge.

## Verification State

- Scope drift finding: acknowledged and documented as an explicit PR 73
  exception, not silently dismissed.
- Multi-PR strategy for future comparable work: required.
- Merge approval: user_requested_merge in the active Codex thread on
  2026-06-05; release approval remains separate.
- Release readiness: not_assessed.
