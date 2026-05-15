→ Read docs
✱ Glob "docs/**/*.md" in . · 32 matches
→ Read docs/reviewer-entrypoint.md
→ Read docs/agent-entrypoint.md
→ Read docs/output-location-map.md
→ Read docs/profile-selection-guide.md
→ Read docs/overclaim-checklist.md
Now let me check for duplicate overclaim content and ambiguous state terms across other docs.
✱ Grep "(external_production_trust=true|trusted_contract_release=true|production_release_verified=true)" in docs · 6 matches
✱ Grep "(Forbidden Claims|overclaim|does not decide|merge approval|release readiness|production trust authority)" in docs · 42 matches
✱ Grep "(ready|approved|passed|success|verified)" in docs · 46 matches
→ Read docs/concepts.md
→ Read docs/claim-authoring.md
✱ Grep "\b(passed|failed)\b" in docs · 16 matches
✱ Grep "(offline_dev|stale|ready|blocked|acceptable)" in docs · 24 matches
Here are the actionable findings against the contract:
---
**1. concepts.md:83-102 — External Verdict reuses canonical result-state names with identical descriptions**
The External Verdict section defines `pass`, `fail`, `not_assessed`, `cannot_verify` with the same wording as the canonical verifier result states in `agent-entrypoint.md:280-286`, plus `warn`. A cold reader cannot distinguish these from result states without reading the disclaimer at line 98. This violates: *"No doc outside the canonical contract redefines or invents result states."*
**Fix**: Rename External Verdict values to a distinct vocabulary (e.g., `verdict_pass`, `verdict_warn`, `verdict_fail`, `verdict_not_assessed`, `verdict_cannot_verify`), or restructure the section to reference the canonical states rather than re-listing them.
---
**2. process-metric-catalog.md:33-34, 44 — Metric enum values `passed`/`failed` look like result states**
`schema_validation_state`, `artifact_safety_scan_state`, and `build_test_evidence_state` use values `passed`, `failed`, `not_assessed`. These are metric catalog enums, not verifier result states, but the near-identical naming (`passed` vs canonical `pass`, `failed` vs canonical `fail`) combined with `not_assessed` (which IS a canonical result state) creates ambiguity. Violates review prompt question 6.
**Fix**: Use distinct names (e.g., `evidence_present`/`evidence_absent`) or add an explicit disclaimer that these are metric enum values, not verifier result states.
---
**3. output-location-map.md — Missing command families that write files**
The map omits:
- `doctor --profile <profile>` → writes `--out <file>` (agent-entrypoint.md:187)
- `install repo-observer` → writes `--out <file>` (agent-entrypoint.md:188)
- `harness validate` → writes `--out <validation.json>` (agent-entrypoint.md:196)
Violates: *"An output location map exists mapping command family → default output path → format → purpose → trust boundary."* (incomplete coverage)
**Fix**: Add rows for these three command families.
---
**4. agent-entrypoint.md:354-366 + reviewer-entrypoint.md:141-170 — Embedded overclaim summaries create drift risk**
Both files contain substantial forbidden-claims summaries alongside the link to `docs/overclaim-checklist.md`. While both include "authoritative only when it matches" caveats, the contract states: *"One canonical overclaim checklist exists; other docs link to it."* The embedded summaries are large enough (5+ bullet points each) that they could drift.
**Fix**: Replace embedded summaries in both files with a single-line link to `docs/overclaim-checklist.md`, or reduce to one non-driftable sentence.
