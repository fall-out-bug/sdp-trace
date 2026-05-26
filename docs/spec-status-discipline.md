# Spec Status Discipline

This repository tracks spec work across separate status axes. Do not collapse
them into a single "ready" or "complete" claim.

## Status Axes

| Axis | Question | Allowed Values |
| --- | --- | --- |
| `spec_state` | Is the written spec reviewed and approved for implementation? | `draft`, `pending_review`, `approved`, `in_review`, `blocked`, `historical` |
| `task_state` | Do `tasks.md` checkboxes match repository reality? | `not_started`, `partial`, `checked_complete`, `blocked`, `not_assessed` |
| `implementation_state` | Does implementation for the claimed scope exist in the repository? | `not_started`, `partial`, `implemented`, `not_assessed` |
| `review_state` | Has the required review completed with remaining findings resolved? | `not_assessed`, `pending`, `changes_requested`, `lgtm`, `blocked` |
| `merge_state` | Is the scoped implementation merged to `main` and reconciled locally? | `not_started`, `pr_open`, `merged`, `not_assessed` |
| `trust_state` | Is the claimed gate/verdict backed by current evidence? | `pass`, `fail`, `cannot_verify`, `not_assessed` |

## Closure Rule

A spec can be called `complete` only when all relevant axes are closed for the
claimed scope:

- `task_state=checked_complete`;
- `implementation_state=implemented`;
- `review_state=lgtm` or an explicitly accepted review disposition exists;
- `merge_state=merged`;
- `trust_state` is evidence-backed, or explicitly out of scope for that spec.

If any axis is missing, keep the spec out of the formal `complete` bucket and
name the missing axis. For example: `implemented; review not_assessed` is valid,
but `complete` is not.

## Roadmap Rule

`docs/roadmap.md` is the navigation surface, not the trust authority. It must
show both:

- the formal lifecycle state; and
- the current repository reality when tasks or implementation are ahead of the
  formal lifecycle state.

The roadmap must not use `draft` as shorthand for "nothing implemented" and
must not use checked task boxes as merge, review, or trust approval.

## Evidence Rule

Use live commands, current PR state, or direct file inspection for claims about
current repository state. Checked-in ledgers, review files, and JSON artifacts
are context until replayed or externally signed.

<!-- sdp-trace-claim: claim=profile_passed; subject=spec-status-discipline-doc; state=pass; profile=repo_baseline_structural; evidence=state:claim_tags_consistent -->
