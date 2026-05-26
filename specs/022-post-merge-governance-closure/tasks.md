# Tasks: Post-Merge Governance Closure

Status: active clarification complete; implementation tasks opened.

## Active Tasks

- [ ] T022-001 Summarize PR #60, PR #63, and Spec 019 governance evidence,
  refreshing live PR/CI state when available. Verification: exact PR, commit,
  CI, and review references recorded; unavailable live state remains
  `not_assessed` or `cannot_verify`.
- [ ] T022-002 Cite the existing `split_successor` maintainer decision for Spec
  019 residual governance. Verification: decision row, roadmap, Spec 019
  plan/tasks, and post-merge closure plan references match.
- [ ] T022-003 Prepare remediation specs only for work still required after the
  existing split decision, or explicitly record that no residual remediation
  remains. Verification: successor specs have reviewed triplets before
  implementation, or the no-remediation state has cited evidence.
- [ ] T022-004 Update closure decision ledger, spec reality ledger, and roadmap
  together. Verification: the three surfaces report the same Spec 022 closure
  state.
- [ ] T022-005 Verify closure docs after decision recording. Verification:
  `go run ./tools/doccheck`; `go run ./tools/hygienecheck`; `git diff --check`.
