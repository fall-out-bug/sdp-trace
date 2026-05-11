# Product Contract v0 - Independent Review

**Date**: 2026-05-10
**Scope**: `spec.md`, `plan.md`, `example.md`, `traceability.md`, `tasks.md`, prior 003 findings, `AGENTS.md`

---

## 1. Verdict

**`APPROVE_FOR_USER_REVIEW`**

Product Contract v0 creates a workable hard gate. The mechanism is review-based rather than mechanically enforced, which is appropriate for a contract-only slice. The 003 linkage is claimed as done but not verifiable from this packet - the user should confirm before approval. Remaining issues are real but non-blocking for the go/no-go decision.

---

## 2. What Works

**The hard gate mechanism is concrete and reviewable.** `spec.md:160` ("If `packet_rows` is empty, the item cannot be P0 product progress") is a binary, citable rule. Combined with the intake template in `traceability.md:42-53`, a reviewer has everything needed to reject a substrate-only feature. The five required intake fields (`packet_rows`, `evidence_surface`, `closure_state`, `buyer_effect`, `non_goal`) are specific enough to prevent dilution through vague language.

**All seven open questions from the 003 Socratic review are answered.**
- CTO packet surface -> Markdown + evidence bundle (`spec.md:33-40`)
- Minimum content with gaps -> 11 rows, all valid in `not_assessed` (`example.md:36-48`)
- Russian-market P0 -> Baseline profile requires no GitHub (`spec.md:128-146`)
- P0 theater findings -> 4 codes bound, 4 deferred with visible deferral (`spec.md:109-125`)
- Discovery closure -> Discovery is not P0 unless it maps to rows (`plan.md:140`)
- General-purpose agent boundary -> Software-delivery boundary evidence required (`spec.md:265`)
- Signed attestation -> P2 additive profile, not trust shortcut (`spec.md:236`, `plan.md:141`)

**The example packet demonstrates the two-minute CTO read (SC-001).** The executive summary (`example.md:11-21`) is ten lines covering what changed, what's observed, what's claimed, what's missing, and who decides. The row-state table scans quickly: three `pass`, four `not_assessed`, two `cannot_verify`, one `partial`, one `fail`. The red flags are immediately visible.

**The traceability matrix preserves existing work.** Every substrate area (flight recorder, evidence bundles, harness observation, authority, CI artifacts, PR review, witness/release proof, adapter capture) is mapped as a source for packet rows (`traceability.md:25-36`) rather than dismissed.

**The roadmap reclassification table (`plan.md:130-143`) is the strongest anti-substrate-creep mechanism.** It explicitly names GitHub adapter, GitFlic/local, OpenCode/GSD, `pi`/GSD2, signed attestation, and dashboard as "not P0 by itself" with specific conditions. This directly addresses the repeated failure mode.

---

## 3. Blocking Findings

| id | severity | cited file:line | finding | why it matters | exact fix |
| --- | --- | --- | --- | --- | --- |
| PC-REV-001 | major | `tasks.md:24-28` (T006-T008 checked), `spec.md:241-242` (FR-015) | 003 roadmap linkage is claimed as complete (T006-T008 checked) but the review packet does not include updated 003 files. The hard gate's enforcement power depends on 003 actually referencing this contract as an implementation blocker. | If 003 was not updated, the gate is aspirational in one direction: this contract blocks 003, but 003 does not know it is blocked. Future work on 003 could proceed without packet-row mapping. | Before user approval, either (a) include updated 003 files in the review packet showing the contract reference, or (b) add a T006-verify step in Phase 2 that explicitly confirms the 003 updates exist and cites the specific file/line where the block is recorded. |
| PC-REV-002 | major | `spec.md:83-92` | Eight evidence states are listed but three (`missing_telemetry`, `unsupported`, `not_integrated`) have no definition anywhere in the contract. The example does not use them. An implementer cannot distinguish `not_integrated` from `not_assessed` or `unsupported` from `cannot_verify` without guessing. | Ambiguous states lead to inconsistent row assignments across implementations, which is the exact class of failure the contract exists to prevent. | Add one-line definitions for each state in `spec.md` after line 92. If they are future states not needed for v0, explicitly say so and remove them from the required-state list until they are bound. |

---

## 4. Non-Blocking Concerns

**Theater derivation rules are deferred but PC-THEATER is a required row.** `traceability.md:64` lists "Theater reason-code derivation is not implemented" as a known gap. The 4 P0 reason codes are named with conditions (`spec.md:111-116`), which is enough for manual application, but automated consistency will require explicit derivation rules. Acceptable for contract v0; must be bound before implementation.

**`PC-THEATER` row state semantics are implicit.** The example shows `fail` when theater findings exist (`example.md:45`) and `pass` when the packet avoids scope theater (`example.md:57`). The spec does not state what `pass` means for PC-THEATER (zero findings? all findings closed?) or how the row state aggregates from individual finding states. Not blocking because the example is consistent, but an implementer will need this.

**Backlog gate enforcement is review-based, not mechanical.** The gate works through SpecKit review and human judgment. This is appropriate for v0 and is how most real product gates function. A future schema or linter could enforce `packet_rows` non-empty, but that is implementation work correctly deferred to Phase 3.

**The intake template uses prose YAML in a Markdown code block** (`traceability.md:42-53`). When implementation begins, this needs to become a schema, Go struct, or validated form. The current shape is correct for contract approval.

**SC-005 is tautological.** "A focused Socratic review finds no remaining critical or major ambiguity" is a success criterion for the review itself, not for the product. Not harmful, but not a meaningful success criterion for the contract artifact.

---

## 5. Missing Evidence or `not_assessed` Areas

| area | state | what would close it |
| --- | --- | --- |
| Updated 003 files showing contract reference | `not_assessed` in this review packet | Include updated `003` files or confirm T006-T008 by file:line citation |
| Evidence state definitions for 3 of 8 states | `not_assessed` | One-line definitions or explicit deferral with removal from required list |
| Theater derivation rules for P0 codes | Deferred per `traceability.md:64` | Documented rules and tests before PC-THEATER automation |
| Packet schema / Go model | Deferred per `tasks.md:43` | `change-evidence-packet-v0` schema after contract approval |
| Real (non-`sha256:example`) evidence digests | Deferred (example-only) | Implementation fixtures after approval |
| Actual CTO readability test (SC-001) | `not_assessed` | User or CTO proxy reads example.md and confirms two-minute comprehension |

---

## 6. Scope-Control and Overclaim Risks

**Low overclaim risk.** The contract is unusually disciplined about what it does not claim:
- `example.md:3` - "example-only, not product proof"
- `spec.md:265` - general-purpose agents only with software-delivery boundary
- `spec.md:236` - signed attestation is additive, not a shortcut
- `plan.md:162-169` - six explicit non-goals
- `example.md:77-83` - six explicit "does not prove" statements
- `traceability.md:9-21` - all 11 rows rated "Partial" or "Good substrate", none claimed as "Complete"

**One residual scope risk:** The roadmap reclassification table (`plan.md:130-143`) is plan-level guidance, not spec-level requirement. If future work ignores the plan, the spec alone does not prevent a GitHub adapter from being labeled P0. The actual enforcement is FR-004/FR-005 in the spec plus the intake template in traceability.md, which together are sufficient. But the reclassification table itself has no FR backing it. Consider promoting the reclassification principle into a spec FR in a future revision if drift appears.

---

## 7. One Strongest Reason to Proceed

The contract solves the exact failure it targets. The repeated product failure was "substrate features treated as P0 product progress." Product Contract v0 answers this with a binary rule (`spec.md:160`), concrete intake fields (`traceability.md:42-53`), a worked example, and explicit reclassification of every known integration target as "not P0 by itself" (`plan.md:130-143`). The gate is reviewable, citable, and honest about what it does not enforce mechanically.

---

## 8. One Strongest Reason Not to Proceed Yet

PC-REV-001: the 003 linkage is the gate's enforcement power in one direction. This contract says 003 is blocked, but if 003 does not reference this contract back, the block is one-sided. Before user approval, confirm that T006-T008 produced actual file changes in 003 that a reviewer can cite. Without this, the hard gate has a gap that future 003 work could walk through unintentionally.
