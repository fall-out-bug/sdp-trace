# Product Classification & Contract Consistency Re-Review

**Reviewer plane**: Product classification and P0 contract consistency
**Reviewed files**: `spec.md`, `plan.md`, `example.md`, `example-local-baseline.md`, `traceability.md`, `tasks.md`
**Prior review reference**: `2026-05-10-full-pi-review.md` findings PCV0-001 through PCV0-018

---

## Verdict

**`APPROVE_FOR_USER_APPROVAL`**

From the product classification plane, the P0 classification gate is now enforceable at the contract level. All ten required-before-user-approval items from the full review are addressed. The direction is ready for explicit user sign-off before implementation.

---

## Findings

| id | severity | file/section | finding | exact fix |
| --- | --- | --- | --- | --- |
| PCR-001 | minor | `spec.md`, P0 Classification Rule | Forward progress definition does not explicitly list `partial -> pass` or `partial -> fail` transitions. The bullet says "from unknown/missing to ..." but does not cover improving an already-assessed row that has partial evidence. | Add "or moves from `partial` to `pass`/`fail` by closing or confirming the material gap" to the forward progress list. Prevents a defender from claiming that only first-entry into a row counts. |
| PCR-002 | minor | `spec.md`, P0 Classification Rule (`buyer_effect`) | The `buyer_effect` field is prose-only and the weakest enforceable field. A team can fill it with vague comfort language ("CTO will have more visibility") without changing any row. | Require that `buyer_effect` cites at least one specific packet section (e.g., "resolves the `PC-THEATER/agent_claimed_verification` finding" or "closes the `PC-VERIFICATION` CI witness gap"). One sentence max. |
| PCR-003 | minor | `traceability.md`, P0 Classification Template | The template uses `packet_rows:` with YAML-like indentation while `spec.md` uses `packet_rows`: in prose. Not a semantic risk, but will cause confusion when implementation generates intake checklists. | Unify the template syntax in both files - either both use structured YAML-like blocks or both use inline prose lists. |
| PCR-004 | info | `example.md` PC-VERIFICATION row | Row state is `partial` which is correct, but the gap/next-evidence column says "Need retained CI artifact, customer witness evidence, or signed verification record." The "or" makes closure easier than the spec's independence rule requires - a signed record could still be the same agent claiming its own tests. | Change "or signed verification record" to "or independent signed verification record witnessed by a non-claiming party" to match the independence definition in the theater section. |

---

## Prior Full-Review Resolution Status (Classification Plane)

| prior finding | verdict in this re-review |
| --- | --- |
| PCV0-001: citation gate without forward progress | **Resolved.** `spec.md` now requires `target_transition` with five explicit forward-progress conditions. FR-005 blocks features that don't progress a row. |
| PCV0-002: permanent `not_assessed` passes via citation | **Resolved.** Forward-progress rule explicitly says "Repeating `not_assessed -> not_assessed` or `cannot_verify -> cannot_verify` without a new evidence surface, narrower claim, or explicit unsupported state does not qualify." |
| PCV0-006: incomplete state vocabulary | **Resolved.** Evidence States table defines all eight states with buyer-readable meanings. |
| PCV0-007: source strength -> implicit score | **Resolved.** Spec says "Projections MUST NOT rank, aggregate, color-score, or present source classes as trust scores, confidence levels, maturity levels, or readiness levels." |
| PCV0-008: undefined profile taxonomy | **Resolved.** Profile Taxonomy table defines required inputs, achievable rows, and commonly unassessed rows per profile. |
| PCV0-011: `PC-RESIDUAL-GAPS` pass is ambiguous | **Resolved.** Rule now says the row can be `pass` only when all known gaps are enumerated with source row, state, reason, and required closure evidence. Examples show `pass` with explicit gap listing. |
| PCV0-017: draft-complete vs. approved status | **Resolved.** `tasks.md` Completion Rule states "Checked draft tasks are draft-complete, not product approval." |

---

## P0 Classification Rule Enforceability Paragraph

The revised P0 classification gate blocks substrate-only work from passing as P0 by requiring two independent conditions: (1) explicit citation of one or more named packet rows via `packet_rows`, and (2) a `target_transition` that demonstrates measurable forward progress on the cited row. Forward progress is defined with five concrete sub-conditions - state transitions, new inspectable evidence surfaces, theater derivation rules, fact/claim separation, or reduced overclaim risk - and the spec explicitly forbids "repeating `not_assessed -> not_assessed` without a new evidence surface" from qualifying. A backlog item that only adds GitHub adapter code but cannot name a changed row state or new evidence surface fails the gate. The classification rule is therefore structurally sound: it forces any P0 claim to be falsifiable by inspection of the named packet row before and after the change. My minor finding PCR-001 (missing `partial -> pass`/`partial -> fail` coverage) does not weaken the gate because such transitions would naturally satisfy the "state moves" or "narrower claim" conditions in practice, but making it explicit would remove a future argument surface.
