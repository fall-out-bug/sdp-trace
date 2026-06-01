# Data Model: Post-Merge Governance Closure

This feature models repository governance records, not product runtime data.

## Governance Evidence Summary

Fields:

- `spec_id`: `019`
- `successor_spec_id`: `022`
- `pr_refs`: PR #60 as the missed-gate merge surface and PR #63 as the
  post-merge integration/supersession surface
- `commit_refs`: merge commits and relevant head commits
- `ci_refs`: live CI references; checked-in CI references are context only
  until live-refreshed or externally signed
- `review_refs`: retained review artifacts and dispositions
- `missing_states`: explicit `not_assessed` or `cannot_verify` entries

Validation rules:

- PR #60 merge approval remains `not_assessed`; later CI, review, or closure
  evidence must not be reclassified as pre-merge approval.
- PR #63 evidence may support post-merge closure but must not retroactively
  approve PR #60.

## Maintainer Decision Reference

Fields:

- `decision_id`: `D006`
- `decision_state`: `split_successor`
- `source`: `docs/closure-decision-ledger.md`
- `successor_spec`: `specs/022-post-merge-governance-closure/`
- `superseded_by`: optional future maintainer decision

Validation rules:

- `split_successor` is the current state unless a new maintainer decision
  explicitly supersedes it.
- CI, review, and checked task boxes cannot change this state by inference.

## Remediation Disposition

Fields:

- `residual_work_state`: `none` or `successor_specs_required`
- `successor_specs`: list of reviewed successor spec triplets, if required
- `evidence_refs`: supporting references for no-remediation or split work

Validation rules:

- If residual work remains, successor specs must have reviewed `spec.md`,
  `plan.md`, and `tasks.md` plus a retained review artifact before
  implementation.
- If no residual work remains, the no-remediation state must be recorded
  explicitly.

## Closure Surface Update

Fields:

- `closure_decision_ledger_state`
- `spec_reality_ledger_state`
- `roadmap_state`
- `updated_together`: required same-change update marker

Validation rules:

- The three surfaces must report the same Spec 022 closure state in the same
  change.
- If one surface cannot be updated, closure remains incomplete.

## Verification Record

Fields:

- `live_pr_ci_refresh`: `verified`, `not_assessed`, `cannot_verify`, or
  `failed`
- `doccheck`: command result
- `hygienecheck`: command result
- `diff_check`: command result

Validation rules:

- `complete` requires live PR/CI refresh evidence when GitHub access is
  available. A missing live-state record is allowed only when the concrete
  access failure is recorded.
- Local docs checks must pass before claiming local closure readiness.
