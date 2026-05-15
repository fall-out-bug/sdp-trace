→ Read docs
✱ Glob "docs/*.md" 31 matches
✱ Glob "docs/**/*.md" 32 matches
→ Read docs/claim-authoring.md
→ Read docs/accountability-model.md
→ Read docs/canonical-vs-projection.md
→ Read docs/concepts.md
→ Read docs/evidence-policy.md
→ Read docs/adoption-guide.en.md
→ Read docs/ci-check-policy.md
→ Read docs/change-evidence-packet.md
→ Read docs/authority-envelope.md
→ Read docs/agent-entrypoint.md
→ Read docs/ci-artifact-observation-downstream.md
→ Read docs/harness-integration.md
✱ Grep "assessed_gap" 75 matches
✱ Grep "`stale`" 5 matches
→ Read docs/overclaim-checklist.md
→ Read docs/adoption-ladder.md
Here are the actionable issues found against the contract:
**Issue 1 — `stale` is an unclassified state token in claim-authoring.md:14**
`stale` appears as an allowed `state` value alongside `pass`, `fail`, `not_assessed`, `cannot_verify` but is not classified in the canonical state contract (`docs/agent-entrypoint.md:275-341`). The contract requires all non-result-state tokens to be classified (telemetry labels, completeness markers, PR-review sub-states, external verdict sub-states, integration labels, or authority scope labels). `stale` belongs to none of those categories. A cold user sees it listed among verifier-like state values with no classification anchor. Already flagged in `specs/014-docs-ux-command-guidance/reviews/raw/socratic-qwen.md:26`.
**Issue 2 — `assessed_gap` is used as a status token in product docs but never classified**
Used in `docs/agent-entrypoint.md:95,103` and `docs/ci-check-policy.md:103` as a standing gate status (e.g., "Maintainability Index remains `assessed_gap`"), 75 matches across the repo. Not classified in the canonical state contract. A cold user reading the agent entrypoint sees `assessed_gap` used in backtick-formatted status prose — could be misread as a verifier result state alongside the five canonical ones. Its definition exists only in `specs/004-mvp-readiness-hardening/spec.md:85`, which is a spec doc, not the canonical docs contract.
**No issues found for questions 1, 3, 4, 5, and 7**:
- Q1: The quick-reference table at the top of reviewer-entrypoint.md is task-oriented and findable in under 30s. ✓
- Q3: `docs/output-location-map.md` is linked from reviewer-entrypoint.md and provides a clean table mapping command → default output → trust boundary. ✓
- Q4: `docs/profile-selection-guide.md` is linked and has a decision flow mapping trust profiles, assessment profiles, witness kinds, and authority scopes. ✓
- Q5: The canonical `docs/overclaim-checklist.md` is linked from reviewer-entrypoint.md, agent-entrypoint.md, and adoption-guide.en.md. No doc duplicates it without linking to it. ✓
- Q7: No doc duplicates the overclaim checklist content without an explicit "authoritative only when matches" guard. All summaries link to the canonical file and condition their authority. ✓
