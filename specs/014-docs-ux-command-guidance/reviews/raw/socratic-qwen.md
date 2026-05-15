I'll conduct an adversarial review of the 014 spec against the current docs. Let me first examine the full doc set to identify gaps and inconsistencies.
→ Read docs
→ Read docs/reviewer-entrypoint.md
→ Read docs/agent-entrypoint.md
→ Read docs/concepts.md
→ Read docs/agent-onboarding.md
→ Read docs/README.md
✱ Grep "missing_telemetry|not_integrated" in docs · 12 matches
✱ Grep "external_production_trust|repo_baseline_structural|source_bound_local_release" in docs · 22 matches
→ Read docs/adoption-guide.en.md
✱ Grep "observed|pass|fail|not_assessed|cannot_verify" in docs · 370 matches
→ Read docs/claim-authoring.md
Here is the adversarial review. I checked the spec, plan, tasks, and all current docs for the four focus areas plus cold-user UX.
---
## Adversarial Review: 014-docs-ux-command-guidance
### P0 — State Vocabulary Is Not Canonical; It Is Fragmented
**Finding:** The spec's US-002 claims a reviewer can distinguish six states (`not_assessed`, `cannot_verify`, `missing_telemetry`, `observed`, `pass`, `fail`). The actual codebase uses at least **14 distinct state-like terms** across docs, and the "canonical" table in `docs/agent-entrypoint.md:273-284` only lists five.
States found outside the canonical table:
| State | Where defined | Conflict |
|---|---|---|
| `missing_telemetry` | `adoption-guide.en.md:84`, `agent-onboarding.md:73`, `harness-integration.md:37`, `flight-recorder.md:98`, rollout playbooks | Not in canonical table; is it a result state or a telemetry label? |
| `not_integrated` | `agent-onboarding.md:73`, `flight-recorder.md:98-99,142`, `harness-integration.md:37` | Not in canonical table |
| `unsupported` | `flight-recorder.md:98`, `harness-integration.md:37` | Not in canonical table |
| `retention_limited` | `flight-recorder.md:99` | Not in canonical table |
| `warn` | `concepts.md:89` (External Verdict section) | Canonical table has no `warn`; concepts.md is a core doc |
| `stale` | `claim-authoring.md:14` | Allowed in claim tags but absent from canonical table |
| `outside_authority` | `agent-entrypoint.md:211` (authority-envelope row) | Not in canonical table |
| `coverage_satisfied/partial/unresolved` | `reviewer-entrypoint.md:122-123` | pr-review-specific; not explained as sub-states vs result states |
| `complete/partial` | `agent-entrypoint.md:190` (import-transcript) | Source completeness states, not result states — but vocabulary collision |
**Actionable:** The canonical table must either (a) enumerate ALL result-state tokens and classify sub-states/labels separately, or (b) explicitly declare which tokens are result states vs telemetry labels vs completeness markers. A cold user seeing `missing_telemetry` in a report and `not_assessed` in the state table has no way to know they are different vocabularies.
---
### P0 — Profile Selection Has Two Disjoint Lists With No Decision Path
**Finding:** `docs/agent-entrypoint.md:13-30` lists three "profile IDs" (`repo_baseline_structural`, `source_bound_local_release`, `external_production_trust`). `docs/reviewer-entrypoint.md:67-73` lists five "assessment profiles" (`adapter-capture`, `managed-harness`, `forensic-retention`, `ci-artifact-observation`, `authority-envelope`). These are presented as separate lists with no explanation of how they relate.
Additionally:
- `local_dirty_structural_only` is called an "authority scope" (`agent-entrypoint.md:27-30`) but appears alongside profile IDs in the reviewer entrypoint's dirty checkout section (`reviewer-entrypoint.md:99-103`), creating profile/scope confusion.
- `claim-authoring.md:15` uses `repo_baseline` (without `_structural`) — a direct naming mismatch with `repo_baseline_structural`.
- `adoption-guide.en.md:88-97` mixes assessment profiles, commands (`gate`, `witness`, `release-proof`), and conceptual guidance ("Air-gapped guidance") in one table, further blurring the boundary.
**Actionable:** The spec's FR-003 ("which profile do I use?" decision tree) must explicitly separate: (1) trust profiles, (2) assessment profiles, (3) witness kinds, (4) authority scopes. The current docs present them as overlapping lists. The decision tree should be a single flowchart or table, not scattered across three docs.
---
### P1 — Output Location Map Is Implicit, Not Explicit
**Finding:** US-003 requires an output location map. The README's "What It Produces" section (`README.md` lines after "## What It Produces") lists seven artifact types but gives no directory paths, file names, or purpose distinctions. The adoption guide's "What To Inspect" section (`adoption-guide.en.md:100-111`) lists some paths but is incomplete and uses different terminology.
Current output destinations scattered across docs:
- `.sdp-trace-runs/<run-id>/` — run artifacts (mentioned in README, adoption guide, reviewer entrypoint)
- `.sdp-trace-report/` — report package (README, adoption guide)
- `.sdp-trace-report/summary.json`, `evidence-table.json`, `missing-telemetry.json`, `gate-result.json`, `ci-witness.json` (adoption-guide.en.md:101-105)
- `--out <file>` per command (agent-entrypoint command table)
- Query-pack output (`--out <file>`)
- Assessment output (`--out <file>`)
- Witness output (`--out <file>`)
- Release-proof output (`--out <file>`)
**Actionable:** The spec needs a single table mapping: command family → default output path → output file format → purpose → trust boundary. The current state requires reading 3+ docs to reconstruct this.
---
### P0 — Overclaim Rules Are Duplicated, Not Consolidated
**Finding:** US-004 requires one canonical overclaim checklist. Forbidden claims appear in at least **four** documents with overlapping but non-identical wording:
| Location | Content |
|---|---|
| `agent-entrypoint.md:297-305` | "Forbidden Claims" — 5 bullet points |
| `reviewer-entrypoint.md:119-146` | "Gate, Witness, And Release Caveats" + "What You May State From Output" — 10+ bullet points |
| `agent-onboarding.md:75-76` | "Never convert missing evidence into success. Never describe a local verifier fact as production trust." |
| `adoption-guide.en.md:11-13` | Opening paragraph caveats |
The reviewer entrypoint alone has three separate overclaim sections: "Gate, Witness, And Release Caveats" (line 119), "What You May State From Output" (line 135), and the final sentence (line 189). These overlap significantly with agent-entrypoint's "Forbidden Claims."
**Actionable:** FR-005 must produce ONE file (e.g., `docs/overclaim-checklist.md`) with a single enumerated list. All other docs should link to it. The current duplication is a drift risk — if one list is updated and another isn't, they will contradict.
---
### P1 — Reviewer Entrypoint Violates Its Own "Minimal" Claim
**Finding:** `reviewer-entrypoint.md:189` says "This entrypoint is intentionally minimal." It is 190 lines with: a 7-step verification path, exit code contract, 17-command surface list, 5 assessment profiles, 4 witness kinds, air-gapped guidance, harness observation caveats, dirty checkout behavior, not-assessed interpretation, gate/witness/release caveats, allowed-state rules, a 10-row quick reference table, and a 5-step PR review handoff.
A cold user hitting this page for a "five-minute" review will face significant cognitive load. The spec's US-001 ("task-oriented guide rather than a long flat reference") directly addresses this, but the current reviewer entrypoint IS the problem the spec is trying to solve.
**Actionable:** The task-oriented command guide (T010) should be the new reviewer entrypoint. The current reviewer entrypoint should become a reference page linked from the guide, not the first thing a reviewer sees.
---
### P1 — Spec Acceptance Criteria Are Under-Specified
**Finding:** The acceptance criteria say:
- "`docs/reviewer-entrypoint.md` has a short task path before long references" — but the current file already has a verification path section before references. What makes the new one "short"? No metric.
- "State vocabulary is consistent across README, concepts, agent entrypoint, reviewer entrypoint, and adoption guide" — but this is impossible to achieve without resolving the P0 state vocabulary fragmentation identified above. The spec must first define the canonical vocabulary.
- "`go run ./tools/doccheck` passes and covers the command claims it owns" — unclear what "the command claims it owns" means for docs-only changes. Does doccheck validate prose consistency or only command-surface JSON alignment?
**Actionable:** Acceptance criteria need measurable thresholds: e.g., "zero occurrences of `missing_telemetry` outside the canonical state doc," or "doccheck reports no drift between `command-surface` JSON and agent-entrypoint command table."
---
### P2 — Plan Missing Risk Mitigation for Link-Based Consolidation
**Finding:** The plan's Risks section says "Consolidation can hide important caveats if links replace necessary local warnings." This is correct but has no mitigation. If the overclaim checklist moves to a single file and other docs link to it, a reader who only reads one doc (e.g., agent-onboarding) will miss the warnings entirely.
**Actionable:** The plan should specify that critical overclaim warnings remain as one-line inline summaries with a link to the full checklist, not just bare links.
---
### P2 — Tasks Lack Drift Verification for State Terms
**Finding:** T020 says "Remove or link duplicate overclaim/state prose" but there is no task to grep for deprecated or inconsistent state terms across all docs. Given the 14+ state tokens found above, a simple grep audit should be a task.
**Actionable:** Add a task: "Grep all `.md` files for state-like tokens; flag any not in the canonical table for classification."
---
### Summary of Findings
| Severity | Issue | Files |
|---|---|---|
| P0 | State vocabulary is fragmented (14+ tokens, canonical table has 5) | `agent-entrypoint.md:273-284`, `concepts.md:85-93`, `claim-authoring.md:14`, `flight-recorder.md:98-99`, `agent-onboarding.md:73` |
| P0 | Profile selection has two disjoint lists with no decision path | `agent-entrypoint.md:13-30`, `reviewer-entrypoint.md:67-73`, `claim-authoring.md:15` |
| P0 | Overclaim rules duplicated across 4+ docs | `agent-entrypoint.md:297-305`, `reviewer-entrypoint.md:119-146`, `agent-onboarding.md:75-76`, `adoption-guide.en.md:11-13` |
| P1 | No explicit output location map | `README.md`, `adoption-guide.en.md:100-111` |
| P1 | Reviewer entrypoint violates its own "minimal" claim (190 lines) | `reviewer-entrypoint.md` |
| P1 | Acceptance criteria lack measurable thresholds | `spec.md` Acceptance Criteria |
| P2 | Plan has no mitigation for link-based consolidation risk | `plan.md` Risks |
| P2 | Tasks lack state-term drift verification | `tasks.md` Phase 2 |
The spec correctly identifies the problems but under-specifies the canonical vocabulary resolution. Without fixing the state vocabulary fragmentation first, US-002's acceptance criterion ("consistent across all docs") is unachievable.
