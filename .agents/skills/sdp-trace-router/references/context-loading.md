<context_loading>
Load the smallest context pack that can decide the task.

<default_order>
1. Repository rules: `AGENTS.md`.
2. Project-local skill: router, trust workflow, pi review, or quality audit.
3. Current spec/plan/task/evidence/gate/decision/provenance docs.
4. Exact source, schema, test, fixture, and example files involved.
5. Fresh command output.
6. Historical ledgers and prose, marked advisory.
</default_order>

<avoid>
- Loading global UI, deploy, Beads, or product-planning skills unless the user explicitly asks.
- Treating a checked-in report as current live proof.
- Pasting broad repo output when a focused file or error line is enough.
- Reusing stale model/version assumptions without primary-source verification.
</avoid>

<refresh_triggers>
Refresh context when switching blocks, changing manifest subjects, altering schema/contracts, touching verifier logic, or after context compaction.
</refresh_triggers>
</context_loading>
