---

# Product Contract v0 Adversarial Review

Reviewer role: Adversarial overclaim and scope-control reviewer.

---

## 1. Verdict

**`REVISE_BEFORE_USER_REVIEW`**

The contract is a significant step forward and the direction is correct. However, I found **3 critical and 5 major loopholes** that an adversarial agent or pressured product owner could exploit to route substrate work back through the P0 gate, leak support claims, or drift toward surveillance. These must be closed before asking for user approval.

---

## 2. What Works

- The core idea - packet rows as a hard backlog gate - is sound and directly addresses the repeated failure mode described in `AGENTS.md`.
- Evidence states are well-defined and explicitly resist collapsing into scores (`spec.md:82-105`).
- The Russian enterprise baseline profile is declared early and has teeth (`spec.md:128-147`).
- The traceability matrix honestly preserves current work without overclaiming (`traceability.md:9-21`).
- The example packet is a strong teaching artifact: it shows `not_assessed`, `cannot_verify`, and `fail` states without pretending they're resolved (`example.md:36-48`).
- Theater reason codes are concrete and scoped to P0 vs P1 (`spec.md:109-125`).
- The backlog intake template in `traceability.md:42-53` is enforceable in principle.

---

## 3. Blocking Findings

| id | severity | cited file:line | finding | why it matters | exact fix |
|---|---|---|---|---|---|
| **OC-001** | critical | `spec.md:150-160` (Backlog Gate), `traceability.md:42-55` (Backlog Intake Template) | The backlog gate requires `packet_rows`, `evidence_surface`, `closure_state`, `buyer_effect`, and `non_goal` - but there is **no required format, schema, or validation rule** for these fields. An agent can fill them with prose that *names* a row id without actually improving the row. | This is the primary loophole. "Cites a row" != "fills or improves the row." Without a validation rule, an agent can write `packet_rows: [PC-VERIFICATION]` and `closure_state: not_assessed -> not_assessed` and claim P0 status because the row was *cited*. | Add a gate rule: `closure_state` must show a **forward transition** (from `not_assessed`/`cannot_verify` to `observed`, `partial`, `pass`, or documented `unsupported`/`not_integrated`) or the item is not P0. Also add: "repeating `not_assessed` or `cannot_verify` as closure does not satisfy the gate." |
| **OC-002** | critical | `spec.md:70-79` (Required Packet Rows table), `spec.md:59-60` (Input Boundary) | The row for `PC-INITIATOR` allows `not_assessed` as a valid state, and the spec says missing inputs "must remain `not_assessed`" (`spec.md:59`). Combined with no forward-transition requirement, **every row can remain `not_assessed` forever** while the backlog gate is technically "satisfied" because rows were cited. | This is the core gate-softening loophole. The contract creates the appearance of a hard gate while allowing permanent `not_assessed` on every row. An adversarial agent can ship 20 features that all cite `PC-INITIATOR` and never actually bind an initiator. | Add an explicit rule: "A P0 backlog item must show evidence of forward progress on at least one cited row. Repeated `not_assessed` with no new evidence surface is not P0 progress. If no forward progress is possible, the item is substrate or discovery, not P0." |
| **OC-003** | critical | `spec.md:109-117` (P0 Theater codes), `example.md:45` (PC-THEATER row) | The `agent_claimed_verification` theater code detects when "no independent retained evidence exists" - but the contract **does not define what counts as "independent" or "retained."** An agent could claim harness-observed verification is "independent" because it came from a different process than the coding agent. | Without independence criteria, the theater detection itself becomes theater. The most important P0 theater finding (`agent_claimed_verification`) is undefined enough to be gamed. | Add a one-line independence rule: "Independent means the verification artifact was produced by a system or person that is not the same agent runtime that performed the change, and the artifact is retained (not ephemeral, not claimed-from-memory)." Cite this in the theater code definition and in `PC-VERIFICATION`. |
| **OC-004** | major | `spec.md:23-25` (Product Language), `spec.md:103-105` (Source strength) | Source strength is defined as "not a trust score" but there is **no enforcement preventing it from being used as one downstream**. A projection (CLI summary, HTML dashboard) could easily rank `signed` > `ci_witnessed` > `harness_observed` > `agent_claimed` and present it as a de facto score. | This is a slow-burn overclaim vector. The spec says "not a trust score" in prose, but the ordered taxonomy *looks* like a score. Without an explicit anti-ranking rule, projections will rank them. | Add: "Source strength classes MUST NOT be ordered, ranked, aggregated, or presented as a trust score or confidence level in any projection. They are categorical, not ordinal." |
| **OC-005** | major | `spec.md:76-77` (PC-AUTHORITY row), `spec.md:258-266` (Open Decisions) | `PC-AUTHORITY` says "no merge, blame, employment, or policy verdict" - but the row *does* accept an "authority evaluation state." The boundary between "authority evaluation state" and "policy verdict" is not defined. An adversarial agent could present an authority evaluation that functionally *is* a policy verdict (e.g., "authority: compliant") while claiming it's just a "state." | This is the surveillance-drift vector. If authority evaluation states can encode compliance/blame semantics, the packet becomes an employee monitoring artifact under a different name. | Add: "`PC-AUTHORITY` states are limited to: `within_declared`, `exceeded_declared`, `not_assessed`, `cannot_verify`. These states describe the relationship between observed actions and a declared authority envelope. They do not describe compliance, blame, fitness, or policy adherence." |
| **OC-006** | major | `example.md:31` (selected_profile), `spec.md:51-52` (Evidence profiles), `spec.md:128-147` (Russian Enterprise Baseline) | The spec names four evidence profiles: `local`, `change-host-rich`, `harness-observed`, `signed` (`spec.md:51-52`). The example uses `local-enterprise-baseline-v0` (`example.md:31`). **These profile names are not the same set.** The profile taxonomy is undefined - what profiles exist, what inputs each requires, and which states each can reach. | Without a defined profile taxonomy, an agent can invent a profile that claims P0 sufficiency while omitting required inputs. The Russian enterprise baseline works only if the profile contract is explicit about what it can and cannot verify. | Add a profile definitions table: profile id, required inputs, rows it can reach `pass`/`partial` on, rows that must be `not_assessed`, and maximum achievable evidence state per row. Minimum: `local-enterprise-baseline-v0`, `change-host-rich-v0`, `harness-observed-v0`, `signed-v0`. |
| **OC-007** | major | `spec.md:240-241` (FR-014, FR-015), `tasks.md:11-20` (Phase 0 tasks checked off) | FR-014 says "one example packet marked example-only, not product proof" and FR-015 says roadmap items "MUST be blocked from implementation approval until they map to Product Contract v0 rows." But the example packet (`example.md`) is already checked off as complete, and tasks T001-T008 are checked - **yet the contract itself is not reviewed or approved.** The contract is doing exactly what it warns against: presenting draft work as progress. | The contract's own tasks violate its trust rules. Checking off T001-T008 before review approval is the same pattern the contract is designed to prevent. | Uncheck T001-T008 or add a note: "These tasks are draft-complete but not reviewed or approved. They become P0 progress only after Socratic review and user approval." Alternatively, add a task status column with `draft`, `reviewed`, `approved`. |
| **OC-008** | major | `AGENTS.md:29` (Quality Bar), `spec.md:109-117` (Theater codes) | The quality bar says "Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`." The theater codes are presented as gate-verdict claims (the packet *finds* theater), but **theater derivation rules are not defined** - `traceability.md:64` explicitly lists "Theater reason-code derivation is not implemented" as a known gap. | Theater findings in the example packet are asserted without derivation rules. This violates the repo's own quality bar. An agent could produce theater findings that are themselves theater. | Add to the contract: "Theater findings are claims. Each finding must cite the specific missing evidence or inconsistency that triggers the reason code. A theater finding without a cited trigger is itself an overclaim." Mark the example theater findings as `example-only_derivation` until rules are implemented. |

---

## 4. Non-Blocking Concerns

| id | concern | file:line |
|---|---|---|
| NB-001 | The `PC-DECISION` row says "Named role or owner ref" but does not specify whether the owner must *accept* the decision or merely be *named*. Naming != owning. | `spec.md:78` |
| NB-002 | The example packet's `PC-RESIDUAL-GAPS` row has state `pass` - but a "pass" on residual gaps implies gaps are acceptable. Should this be a different semantic, like `recorded`? | `example.md:48` |
| NB-003 | The spec says "general-purpose agents: in scope only with software-delivery boundary evidence" (`spec.md:266`) but does not define what that boundary evidence is. This is the SR-007 finding from the prior review, and it is still open. | `spec.md:266` |
| NB-004 | `FR-012` says "Signed attestation MUST be additive evidence over a packet, not a shortcut that upgrades missing evidence into trust." This is correct but not enforceable without a rule: "A signed attestation over a packet with `not_assessed` or `cannot_verify` rows does not change those row states." | `spec.md:235-236` |
| NB-005 | The backlog gate fields (`packet_rows`, `evidence_surface`, etc.) have no schema or Go type. Implementation will need to decide: are these strings, arrays, structured objects? This is fine for a contract-only slice but will block implementation planning. | `spec.md:150-160` |
| NB-006 | The review prompt asks "Is the example packet useful to a CTO?" but the example packet uses internal terms like `harness-run`, `trace:session-command-digest`, and `external_assertion:task-id:PAY-1842`. A CTO would need a glossary or the packet needs a lay-readable evidence column. | `example.md:36-48` |

---

## 5. Missing Evidence or `not_assessed` Areas

| area | status | what would close it |
|---|---|---|
| Profile taxonomy definition | `not_assessed` | A table mapping profile id -> required inputs -> achievable row states. |
| Theater derivation rules | `not_assessed` | Documented trigger conditions per reason code. |
| Gate validation mechanism | `not_assessed` | Schema, Go validation, or CI check that enforces forward-transition rule. |
| Software-delivery boundary evidence | `not_assessed` | Minimum evidence conditions for general-purpose agent tracking (SR-007 carry-forward). |
| Signed attestation additive rule | `not_assessed` | Enforceable rule that signing does not upgrade `not_assessed` rows. |
| Projection anti-ranking rule | `not_assessed` | Explicit rule preventing source-strength ranking in projections. |

---

## 6. Scope-Control and Overclaim Risks

1. **Gate softening via citation loophole** (OC-001, OC-002): The most critical risk. Without a forward-transition requirement, the backlog gate is a citation gate, not a progress gate. This is the exact failure mode the contract was designed to prevent.

2. **Theater-as-theater** (OC-003, OC-008): Theater findings are the contract's most differentiating feature, but without derivation rules or independence criteria, they risk becoming the same kind of prose-only claim the repo already failed on.

3. **Source-strength score creep** (OC-004): The taxonomy is ordinal in appearance. Without an explicit anti-ranking rule, projections will rank and score. This is a slow-burn surveillance vector.

4. **Authority evaluation drift** (OC-005): The line between "authority evaluation state" and "policy verdict" is the line between traceability and employee monitoring. The contract names the boundary but does not define it.

5. **Profile invention** (OC-006): Without a defined profile taxonomy, any agent can invent a profile that claims P0 sufficiency while omitting required inputs. The Russian enterprise baseline works only if the profile contract is explicit.

6. **Contract self-exemption** (OC-007): The contract's own tasks are checked off before review. This is the pattern it warns against. If the contract cannot hold itself to its own rules, it will not hold the backlog.

---

## 7. One Strongest Reason to Proceed

The contract correctly identifies the core failure: substrate work was being presented as product progress. The packet-row gate, if hardened, is the right mechanism. The direction is sound, the traceability matrix is honest, and the example packet teaches the right lesson. Fixing the gate loopholes (especially OC-001 and OC-002) is a tractable problem, not a reframe.

---

## 8. One Strongest Reason Not to Proceed Yet

**The gate is not actually a gate yet.** Without a forward-transition requirement (OC-001, OC-002), the backlog gate is a citation requirement, not a progress requirement. An adversarial agent can cite every row, claim P0, and ship 20 features that never move any row from `not_assessed`. This is the exact failure the contract was built to prevent, and the current language does not close the loophole. Approving the contract before fixing this would ratify a gate that does not gate.

---

## Disposition Summary

| finding | disposition |
|---|---|
| OC-001: No forward-transition rule | **Blocker** - must fix before approval |
| OC-002: Permanent `not_assessed` loophole | **Blocker** - must fix before approval |
| OC-003: No independence criteria for theater | **Blocker** - must fix before approval |
| OC-004: Source-strength score creep | Major - must fix before approval |
| OC-005: Authority evaluation drift | Major - must fix before approval |
| OC-006: Undefined profile taxonomy | Major - must fix before approval |
| OC-007: Contract self-exemption | Major - must fix before approval |
| OC-008: Theater findings without derivation | Major - must fix before approval |
| NB-001 through NB-006 | Non-blocking - address during implementation planning or as follow-ups |

---

## Next Steps

1. Add forward-transition rule to backlog gate (OC-001, OC-002).
2. Add independence criteria for theater codes (OC-003).
3. Add anti-ranking rule for source strength (OC-004).
4. Define `PC-AUTHORITY` state vocabulary (OC-005).
5. Define profile taxonomy table (OC-006).
6. Fix contract's own task status to match its trust rules (OC-007).
7. Add derivation-citation requirement for theater findings (OC-008).
8. Re-run adversarial review after fixes.
9. Stop for user approval.
