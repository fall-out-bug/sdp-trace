# Block 05 Socratic Dialogue: Customer Pilot Evidence Package

Date: 2026-05-01
Block: `05-customer-pilot-evidence-package`
Beads mirror: `sdp-trace-cdn.22`

## Consensus Candidate

Block 05 should implement run-card and evidence-package artifacts for the requested pilot slices, but it must not claim harness, model, stack, or customer readiness until committed evidence exists.

## Consensus Result

Consensus is accepted for spec review. The selected scope is an evidence-package-first pilot design:

1. Create executable run-cards for OpenCode model slices and harness slices.
2. Define Kotlin+Bazel evidence requirements and keep the fixture `not_assessed` until a real run exists.
3. Update legacy-named compatibility matrices with evidence state, reason codes, artifact references, and gaps only.
4. Build a customer pilot evidence package outline that is safe to commit and does not contain raw customer data.

Implementation is blocked until pi-review findings on the spec artifacts are recorded in the committed review ledger, mirrored in Beads, and closed.

## Q1: Are run-cards enough to record observed behavior or an external compatibility verdict?

**Critic**: A run-card proves that we know what to run, not that the harness or model works.

**Answer**: No. A run-card is a recipe and evidence contract. Observed behavior requires committed sanitized run artifacts or evidence summaries. Support, readiness, and compatibility are external verdicts, not native `sdp-trace` outcomes.

**Resolution**: Matrix rows created from run-cards alone must be `not_assessed`, not `observed`.

## Q2: Should matrices use vague placeholder states or `not_assessed`?

**Critic**: Vague placeholder states can hide missing evidence. `not_assessed` is stronger but assumes the expected evidence is already defined.

**Answer**: Use `not_assessed` for both missing runs and discovery gaps, with reason codes such as `no_run_artifact`, `missing_export`, or `discovery_required`.

**Resolution**: Block 05 docs must remove vague placeholder states and each matrix row must state the evidence reason.

## Q3: Can OpenCode model rows record support if the run-card names MiniMax, Kimi, and GLM?

**Critic**: Naming a model in a plan can look like coverage to a non-technical reviewer.

**Answer**: No. The row can only say the slice is planned and unassessed until a committed artifact records a run with observed model identity. Support would still be an external verdict input.

**Resolution**: OpenCode + MiniMax/Kimi/GLM rows start as `not_assessed` with required evidence listed.

## Q4: Should Superpowers, `gsd`, `gsd2`, and Oh My OpenAgent share one harness verdict?

**Critic**: A shared verdict would blur very different export surfaces and make unbacked claims harder to spot.

**Answer**: No. The run-card may share a template, but each harness row needs its own evidence state for rules location, tool logs, hooks, and evidence export.

**Resolution**: The harness matrix must keep per-harness evidence states and artifact references.

## Q5: Does a synthetic Kotlin+Bazel fixture close the Kotlin+Bazel gap?

**Critic**: A synthetic example can prove schema shape, but not real-world detection or build behavior.

**Answer**: No. It can close the documentation and placeholder gap only. Real Kotlin+Bazel behavior remains `not_assessed` until a real committed run artifact exists.

**Resolution**: The fixture placeholder must explicitly say `not_assessed` with `design_fixture_only`, and docs must prohibit calling T031 a real stack proof.

## Q6: Can a run-card use `pass`, `warn`, and `fail`?

**Critic**: The existing template uses those words, but `sdp-trace` must not own native pass/fail decisions.

**Answer**: Not as native `sdp-trace` verdicts. If an external tool emits those statuses, record them as external verdict inputs with producer and policy reference when available.

**Resolution**: Block 05 run-cards should prefer `observed` or `not_assessed` with reason codes. Any `pass/warn/fail` text must be scoped to external origin and linked as evidence input, not as a row state.

## Q7: Should Block 05 introduce a new run-card schema?

**Critic**: A schema could make package validation stronger, but it also expands scope and creates another contract to maintain.

**Answer**: Not yet. Existing evidence/provenance/trace schemas can represent the JSON artifacts. Markdown run-cards are sufficient for operator workflow in this block.

**Resolution**: Do not add a new schema unless implementation proves existing contracts cannot express required artifacts. If that happens, stop and review the new schema as a separate task.

## Q8: Can customer pilot package docs include real customer examples?

**Critic**: Real examples are useful, but raw customer prompts, source, logs, or secrets cannot be committed.

**Answer**: No raw customer data. Use placeholders, redacted references, hashes, and sanitization notes.

**Resolution**: The package outline must define safe input/output boundaries and a redaction checklist. Private customer inputs are never committed; committed artifacts are sanitized summaries, hashes, redaction notes, and access-neutral references only.

## Q9: What is the main UX risk?

**Critic**: A table full of `not_assessed` can look like failure or noise instead of an honest evidence map.

**Answer**: The table must separate planned rows from observed evidence and explain the next artifact needed for each gap.

**Resolution**: Matrices need columns for evidence state, reason code, artifact reference, gap reason, and next required evidence.

## Q10: What consensus is required before implementation?

**Critic**: The customer pilot scope is broad enough to sprawl into runtime integrations and real pilot execution.

**Answer**: Implementation is limited to docs, run-cards, matrix updates, and safe fixture placeholders. Real pilot execution is out of scope.

**Resolution**: Proceed only after spec pi review is complete, the committed review ledger is updated, and all valid spec-gate findings are closed in Beads.
