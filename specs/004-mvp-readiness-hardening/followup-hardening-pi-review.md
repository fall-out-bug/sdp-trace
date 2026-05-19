# PI Review: Follow-Up Readiness Hardening Delta

**Date**: 2026-05-17
**Harness**: Pi CLI `0.74.1`
**Provider/model**: `kimi-coding/kimi-for-coding`
**Prompt**: `/review-kimi`
**Mode**: non-interactive, read-only tools (`read,grep,find,ls`)
**Artifact**: `specs/004-mvp-readiness-hardening/followup-hardening-spec.md`
**Verdict**: `CONDITIONAL_APPROVE`

## Findings And Dispositions

| Finding | Severity | Disposition | Evidence |
|---|---|---|---|
| Slice 3 lacked an explicit negative test table and did not require pure-Go validation before `exec.Command`. | Critical | `accepted_fixed` | `followup-hardening-spec.md` now requires red tests, pure-Go validation before git commands, schema/Go validator alignment, and a concrete negative-test table. |
| Slice 6 mixed doc-honesty and absolute-MI-improvement paths without default priority. | Major | `accepted_fixed` | `followup-hardening-spec.md` now makes doc honesty the default closure path and treats absolute MI improvement as secondary. |
| Docs and verification for stale MI claims need to remain coupled in the same scoped commit. | Major | `accepted_fixed` | `followup-hardening-spec.md` now states doc closure must be in the same scoped commit as verifier-state evidence when live MI replay fails. |
| The delta made planned closure claims without saying authoritative claims require `sdp-trace-claim` tags later. | Major | `accepted_fixed` | `followup-hardening-spec.md` now states it is a planning delta and not proof of closure. |
| Test-first behavior for releaseproof trust changes was implicit. | Major | `accepted_fixed` | Slice 3 now requires red tests before implementation changes. |
| New Go code should not introduce TODO/FIXME markers. | Minor | `accepted_fixed` | Product boundary now forbids TODO/FIXME markers in new Go code. |

## Remaining Review State

- GLM architecture/trust doubt review: `not_assessed`.
- Qwen or DeepSeek alternate review: `not_assessed`.
- GitHub CI for any future implementation branch: `not_assessed` until queried live for the exact head SHA.
- This review is advisory. Every later implementation finding must still be verified against full files before disposition.

## Re-Review After Amendments

**Date**: 2026-05-17
**Harness**: Pi CLI `0.74.1`
**Provider/model**: `kimi-coding/kimi-for-coding`
**Prompt**: `/review-kimi`
**Mode**: non-interactive, read-only tools (`read,grep,find,ls`)
**Verdict**: `CONDITIONAL_APPROVE` for the spec contract; `BLOCKED` for implementation handoff.

| Finding | Severity | Disposition | Evidence |
|---|---|---|---|
| Previous Slice 3 negative-test and pure-Go-validation gap is resolved at the spec level. | Critical | `accepted_fixed` | Re-review confirmed the concrete rejection table and before-`exec.Command` validation rule are present. |
| Previous Slice 6 MI closure ambiguity is resolved at the spec level. | Major | `accepted_fixed` | Re-review confirmed doc honesty is the default path and absolute MI improvement is secondary. |
| Previous planned-claim ambiguity is resolved at the spec level. | Major | `accepted_fixed` | Re-review confirmed the planning-delta preamble distinguishes acceptance criteria from proof. |
| Releaseproof unsafe source refs remain unimplemented in code. | Critical | `deferred_not_assessed` | This is planned Slice 3 work and must remain open until red tests, pure-Go validation, schema alignment, and focused verification pass. |
| Stale absolute MI pass claims remain unimplemented in parent docs. | Major | `deferred_not_assessed` | This is planned Slice 6 work and must remain open until live MI replay and doc corrections are committed together. |
| Format/import drift remains unimplemented in named files. | Major | `deferred_not_assessed` | This is planned Slice 2 work and must remain open until formatting/import verification passes. |

Re-review conclusion: the planning contract is usable. Implementation remains blocked by the explicit approval gate in `followup-hardening-tasks.md` and by the fact that Slices 2, 3, and 6 are real open work, not spec-complete proof.
