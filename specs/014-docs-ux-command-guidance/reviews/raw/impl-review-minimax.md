✱ Glob "docs/*.md" 31 matches
→ Read docs/overclaim-checklist.md
→ Read docs/concepts.md
→ Read docs/reviewer-entrypoint.md
✱ Grep "pass.*fail|not_assessed|cannot_verify|observed" 2958 matches
→ Read docs/claim-authoring.md
✱ Grep "stale" 210 matches
→ Read docs/agent-entrypoint.md [limit=50, offset=270]
→ Read docs/agent-entrypoint.md [limit=50, offset=319]
Issues found:
**1. `stale` is an unclassified state token — `docs/claim-authoring.md:14`**
`stale` appears as an allowed `state` value in claim grammar but is not classified in the canonical state contract at `docs/agent-entrypoint.md:275-341`. The contract requires all non-result-state tokens to be classified (telemetry labels, completeness markers, PR-review sub-states, external verdict sub-states, integration labels, or authority scope labels). `stale` belongs to none of those categories.
**2. Overclaim checklist is duplicated with drift risk — `docs/reviewer-entrypoint.md:141-170`**
The "Overclaim Checklist" section (lines 141-170) copies the canonical overclaim checklist verbatim with a note that it is "authoritative only when it matches the canonical file." This pattern creates drift risk: a future edit to either file without synchronizing the other will make the statement false. The canonical file at `docs/overclaim-checklist.md` has no such safeguard.
**3. `stale` appears in `specs/001-sdp-trace-time-series-evidence-substrate/data-model.md:110`**
```
signature_status`: valid, invalid, missing, stale, not_assessed
```
Here `stale` appears in a schema data model as a signature status value, also unclassified relative to the canonical contract.
