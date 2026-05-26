# Block 33 Review Ledger: Repository CI Trust System

Status: imported draft. Historical review notes below came from stale branch
`origin/codex/block-33-repo-ci-trust-system`; raw review outputs were not
committed and are `cannot_verify` in the current checkout. Run a fresh Socratic
review before treating Block 33 as approved for implementation.

## Intake Notes

Block 33 starts from the decision to improve CI for the `sdp-trace` repository
itself before designing customer-repository CI templates.

Current constraints:

- Block 32 PR-review CI is already integrated on current `main` and remains
  evidence-only.
- PR-level review evidence, green CI, and ready state are not merge approval.
- Missing GitHub checks stay `not_assessed`, not green.
- Source-bound proof is separate from ordinary PR verification.
- `sdp-trace` owns evidence facts, not downstream branch-protection or release
  policy.

## Historical Socratic Review Findings

Raw review outputs were local scratch under `.codex-review/block33-socratic/`
on the stale branch and are not committed. Treat these rows as advisory draft
context, not current review evidence.

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S33-PB-01 | critical | product-boundary / maintainer UX | Open questions controlled implementation scope and branch protection, so the draft was not approval-ready. | accepted_fixed | Replaced open questions with reviewed decisions: `contract-validate` starts evidence-only/documented-only, source-bound release and trace-evidence workflows are deferred, and latency budgets are explicit. |
| S33-PB-02 | critical | product-boundary / maintainer UX | The draft lacked a current validator inventory, making `contract-validate` aspirational. | accepted_fixed | Added Current Validator Inventory with ready, partial, `not_assessed`, and manual/evidence-only surfaces. |
| S33-PB-03 | major | product-boundary / maintainer UX | Evidence-only `cannot_verify` behavior was ambiguous and could hide gaps behind green checks. | accepted_fixed | State mapping now says required and trust-sensitive evidence workflows fail on `cannot_verify`; observation-mode jobs may pass only with visible Step Summary `cannot_verify`. |
| S33-PB-04 | major | product-boundary / maintainer UX | Branch-protection promotion criteria for `contract-validate` were undefined. | accepted_fixed | Added concrete promotion and backout criteria: concrete local validators, no `not_assessed`, five clean runs, no overrides, separate review, and false-positive demotion rule. |
| S33-PB-05 | major | product-boundary / maintainer UX | The two-PR implementation recommendation did not map to tasks. | accepted_fixed | Implementation plan now maps PR 1 to T247 plus docs and PR 2 to T248-T250, with T251 spanning both. |
| S33-TE-01 | critical | trace/evidence | Required jobs could include unavailable surfaces as `not_assessed`, weakening branch protection. | accepted_fixed | Required jobs now exclude `not_assessed` surfaces; unavailable validators are outside required scope and tracked as future work. |
| S33-TE-02 | critical | trace/evidence | Tooling commands for formatting, YAML, and JSON were not pinned or exact. | accepted_fixed | Spec now names `gofmt -l`, `find ... jq empty`, and pinned Go-based `actionlint` command. |
| S33-TE-03 | critical | trace/evidence / source-bound | Source-bound release workflow over `contract-manifest.example.json` risked overclaiming an example manifest as release proof. | accepted_narrower | Renamed Layer 4 to source-bound release boundary, deferred the workflow, and kept ordinary PR CI away from source-bound release proof. The repo still references the example manifest for the existing source-bound cycle, but Block 33 will not create a release workflow over it. |
| S33-TE-04 | major | trace/evidence | `trace-evidence-artifacts` lacked concrete artifact allowlist and safety gate. | accepted_narrower | Deferred `trace-evidence-artifacts` to a later block; Block 33 may document only the boundary. |
| S33-DX-01 | major | DX / maintainability | Diff-scoped "changed examples" validation could create local/CI drift. | accepted_fixed | Plan now requires static paths or full-directory commands; diff filtering can only be an optimization after the full local command is defined. |
| S33-DX-02 | major | DX / maintainability | Block 33 could scope-creep into writing new validators. | accepted_fixed | Plan now limits Block 33 `contract-validate` to existing committed Go code or reviewed standard tooling. |

## Current Review State

- Product-boundary / UX review: `cannot_verify` for the current checkout.
- Trace/evidence review: `cannot_verify` for the current checkout.
- DX/maintainability review: `cannot_verify` for the current checkout.
- Focused re-review: `cannot_verify` for the current checkout.
- Critical findings remaining: `not_assessed` until fresh review.
- Major findings remaining: `not_assessed` until fresh review.
