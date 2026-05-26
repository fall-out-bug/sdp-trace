# Data Model: Post-Merge Governance Closure

This feature models repository governance records, not product runtime data.

## Governance Evidence Summary

Fields:

- `spec_id`: `019`
- `successor_spec_id`: `022`
- `pr_refs`: PR numbers and URLs for PR #60 and PR #63
- `commit_refs`: merge commits and relevant head commits
- `ci_refs`: live or checked-in CI references
- `review_refs`: retained review artifacts and dispositions
- `missing_states`: explicit `not_assessed` or `cannot_verify` entries

Validation rules:

- PR #60 merge approval remains `not_assessed` unless a new explicit approval
  record is found.
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
  `plan.md`, and `tasks.md` before implementation.
- If no residual work remains, the no-remediation state must be recorded
  explicitly.

## Closure Surface Update

Fields:

- `closure_decision_ledger_state`
- `spec_reality_ledger_state`
- `roadmap_state`
- `updated_together`: boolean

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

- `complete` requires either live PR/CI refresh evidence or an explicit missing
  live-state record.
- Local docs checks must pass before claiming local closure readiness.
