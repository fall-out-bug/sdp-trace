# Product Contract v0 - Adversarial Re-Review

## 1. Verdict

**REVISE_BEFORE_USER_APPROVAL**

The revision is strong. All eighteen prior findings are substantively addressed. The remaining gaps are narrow but two of them sit on exact boundaries the contract is trying to protect: the anti-approval-system boundary and the anti-score boundary. Fixing both is a small, scoped change; shipping them unfixed to the CTO risks the same kind of semantic drift the contract exists to prevent.

---

## 2. Top Findings

| id | severity | file/section | finding | exact fix |
| --- | --- | --- | --- | --- |
| RRV-001 | major | `spec.md` -> "Row-Specific Rules" / PC-DECISION / `example-local-baseline.md` Decision Ownership table | `PC-DECISION` owner states are not enumerated. The local-baseline example shows merge-readiness owner state as `pass` with reason "Bound to local merge-owner policy ref." A CTO scanning the example will read "merge readiness: pass" as "merge is approved." This is the exact approval-system drift the contract forbids. `spec.md` defines explicit state vocabularies for `PC-AUTHORITY` (4 states) and `PC-RESIDUAL-GAPS` (computed rules), but `PC-DECISION` has no row-specific rule section at all. | Add a "PC-DECISION" subsection under "Row-Specific Rules" enumerating owner states - e.g., `owner_identified`, `cannot_verify`, `not_assessed` - and explicitly state that `pass` (or `owner_identified`) means "an owner ref is bound to a named source," not "the decision has been made or approved." Update both example Decision Ownership tables to use the enumerated states. |
| RRV-002 | major | `spec.md` -> "Packet Metadata" table, `example.md:31`, `example-local-baseline.md:28` | `packet_state` is a required metadata field with no defined value vocabulary. Both examples set it to `not_assessed for production trust`, which is free-form prose. Because this field sits at the top of the packet and is the first scalar a reader sees, it is the most likely vector for a future global trust score, health score, or readiness verdict - exactly what the contract prohibits. | Either (a) define `packet_state` as a metadata-status field with allowed values such as `draft`, `under_review`, `reviewed` and add an explicit rule that `packet_state` MUST NOT be a trust/health/readiness score, or (b) remove the field from the required template and let packet-level status live in `PC-RESIDUAL-GAPS` where it is already decomposed by row. |
| RRV-003 | minor | `spec.md` -> "Profile Taxonomy" last paragraph | "A combined profile still preserves the weakest rows" implies an ordinal relationship between states (weakest approx worst). This partially contradicts the categorical-source-strength rule and the anti-ranking rule, even though the sentence refers to row-state combination, not source classes. | Replace "weakest rows" with an explicit combination rule: per-row state is the maximum achievable state across combined profiles, constrained by each profile's own row limits. Or simply state that a combined profile takes the per-row state from the profile that provides the most restrictive assessment for that row, where `fail` > `partial` > `pass` in restrictiveness, and `not_assessed`/`cannot_verify`/`unsupported`/`not_integrated` do not upgrade. |
| RRV-004 | minor | `spec.md` -> "Open Decisions With Proposed Defaults" table | "First change-host rich adapter: GitHub" is clearly a proposed default with rationale, but when quoted out of context in a slide or executive summary it will become "sdp-trace targets GitHub first." | Add a parenthetical: "(default for fastest rich-evidence path; not a product-ontology commitment)." Alternatively, move this table out of the normative spec into `plan.md` where defaults belong. |
| RRV-005 | minor | `spec.md` -> "P0 Classification Rule" fifth forward-progress bullet | "A projection becomes less likely to overclaim because a required non-goal or residual gap is added" is the most prose-like of the five forward-progress types. A team could claim P0 by adding non-goal text without changing any evidence surface or row state. | Narrow to: "a projection adds a required non-goal or residual gap that a previous packet omitted, and the gap names a specific row, evidence surface, and closure condition." This forces the non-goal to be concrete and row-bound. |

---

## 3. Prior Findings Resolution

All eighteen findings from the first full review (`2026-05-10-full-pi-review.md`) are substantively resolved:

| prior id | resolved? | how |
| --- | --- | --- |
| PCV0-001 | yes | Forward-progress rule with five concrete types; "Repeating `not_assessed` -> `not_assessed`" explicitly excluded. |
| PCV0-002 | yes | Same rule covers the permanent-`not_assessed` bypass. |
| PCV0-003 | yes | Seven required sections in order; projections must not omit or rename them. |
| PCV0-004 | yes | Evidence bundle manifest table with eight fields, plus the ref-absence rule. |
| PCV0-005 | yes | Theater trigger conditions, trigger evidence, independence/retained definitions, and `PC-THEATER` row-state rules table. |
| PCV0-006 | yes | Evidence states table with meaning for each state. |
| PCV0-007 | yes | Explicit categorical rule plus MUST NOT prohibition on ranking/scoring. |
| PCV0-008 | yes | Profile taxonomy table with required inputs, achievable rows, commonly-not-assessed rows, and notes. |
| PCV0-009 | yes | `PC-VERIFICATION` now `partial`; harness evidence present but no independent retained CI/customer witness. |
| PCV0-010 | yes | `PC-THEATER` state rules table; example uses `partial` for triggered findings, not `fail`. |
| PCV0-011 | yes | `PC-RESIDUAL-GAPS` computation rules; `pass` means enumeration is complete. |
| PCV0-012 | yes | Owner ref sources enumerated; local-baseline example shows policy-bound owner. |
| PCV0-013 | yes | `example-local-baseline.md` added with TeamCity, local Git, no OpenCode/GSD. |
| PCV0-014 | yes | `PC-AUTHORITY` limited to four explicit fact states. |
| PCV0-015 | yes | Traceability matrix consistently says "Partial" or "Partial substrate." |
| PCV0-016 | yes | "It never upgrades `not_assessed` or `cannot_verify` row states." |
| PCV0-017 | yes | Completion Rule in `tasks.md` distinguishes draft-complete from approval. |
| PCV0-018 | yes | Re-review packet includes `003` linkage files. |

---

## 4. P0 Classification Rule Effectiveness

The rule is **effective as written**, with one remaining weak point.

The five forward-progress types are concrete enough to block most substrate-only P0 claims:

- A new evidence surface must be named and inspectable (bullets 1-2).
- A theater derivation rule must have trigger evidence, not just prose (bullet 3).
- Row precision requires separating facts from claims (bullet 4).
- Bullet 5 is the weakest: adding non-goal text is close to prose-only progress (see RRV-005).

The explicit exclusion of "`not_assessed` -> `not_assessed`" and "`cannot_verify` -> `cannot_verify`" without a new evidence surface closes the main loophole from the first review.

A substrate feature *can* still qualify as P0 - but only when it exposes a new inspectable surface that feeds a named row and moves that row's state. An internal refactor, a schema that is never rendered into a packet, or a discovery report that produces no retained evidence surface will not pass. This is the correct boundary.

The one scenario the rule does not fully prevent is a team that names a row, names a trivial evidence surface (e.g., "the existence of this config file"), and claims `not_assessed` -> `partial` on the basis that the config file can now be inspected. The `buyer_effect` and `non_goal` fields provide some friction here, but the gate's strength ultimately depends on the reviewer enforcing that the evidence surface produces *retained, resolvable* evidence, not just that a file exists. The spec's bundle manifest rules (ref absent from manifest = `cannot_verify`) provide the enforcement hook. This is adequate for a contract-level spec; implementation-level validation should harden it further.

---

**Summary**: Fix RRV-001 (PC-DECISION owner states) and RRV-002 (packet_state vocabulary), then this is ready for CTO approval. The other three findings are polish.
