# Block 05 Design: Customer Pilot Evidence Package and Run-Cards

Status: accepted for implementation; spec pi-review findings closed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.22`
Audience: CTO, CIO, pilot operators, harness evaluators, future `sdp-gate` consumers

## Purpose

Block 05 turns the customer pilot matrix into executable evidence recipes without making support, readiness, or compatibility verdicts.

The block answers the CTO and CIO question:

> Can a pilot operator run the requested harness, model, and Kotlin+Bazel slices and see exactly what evidence exists, what is missing, and which claims lack evidence?

## Product Boundary

`sdp-trace` records pilot evidence, provenance, trace artifacts, and `not_assessed` gaps. It does not decide whether a harness, model, stack, or customer deployment is acceptable.

Run-cards are recipes and evidence contracts, not completed pilot results. Legacy-named compatibility matrices are evidence-state matrices; they must not claim support, readiness, or compatibility as native `sdp-trace` outcomes. Missing pilot runs remain `not_assessed` with a reason code.

Beads is only secondary execution tracking. A repository observer must understand this block from committed SpecKit files, docs, examples, and validation commands.

## In Scope

- OpenCode model run-card covering MiniMax, Kimi, and GLM.
- Harness run-card covering Superpowers-style harnesses, `gsd`, `gsd2`, and Oh My OpenAgent.
- Kotlin+Bazel pilot fixture plan and evidence requirements.
- Customer pilot evidence package outline.
- Harness and model compatibility matrix updates that use observed evidence state, artifact references, `not_assessed`, and reason codes.
- JVM/Bazel guide update for Kotlin+Bazel-specific evidence and anti-heuristics.
- A committed Kotlin+Bazel placeholder or evidence bundle that explicitly remains `not_assessed` until a real run artifact exists.
- Documentation updates that make claims without evidence visible instead of hiding them in prose.
- Spec pi review before implementation, implementation pi review after changes, committed review ledger updates, and Beads registration/closure for every valid finding including minor/P3 items.

## Out of Scope

- Running a real customer pilot in this repository.
- Claiming OpenCode, MiniMax, Kimi, GLM, Superpowers-style patterns, `gsd`, `gsd2`, Oh My OpenAgent, or Kotlin+Bazel support/readiness/compatibility as a native `sdp-trace` outcome.
- Capturing raw customer prompts, customer logs, credentials, secrets, or proprietary source.
- Introducing dependency on OpenCode, Beads, Superpowers, `gsd`, `gsd2`, Oh My OpenAgent, Codex, Claude, GitHub, or any specific runtime.
- Adding `sdp-gate` thresholds, pass/fail policy, readiness decisions, or opaque compatibility scores.

## Recommended Approach

Use an evidence-package-first pilot design.

The run-cards define:

1. The exact prompt or prompt template.
2. Required sanitized artifacts.
3. Required provenance fields.
4. `unbacked_claim` capture.
5. Validation commands where possible.
6. `not_assessed` rules for missing logs, missing identity, incomplete exports, and unrun slices.

The matrices summarize evidence state rather than marketing compatibility:

- `observed` only when a committed sanitized run artifact or evidence summary describes actual observed behavior.
- `not_assessed` when no valid run artifact exists, required export is missing, discovery is still required, or the only artifact is a design/synthetic fixture.

This is stricter than a broad compatibility table, but it preserves trust. A CTO can see whether the pilot path is executable without confusing planned work with proof.

## Pilot Evidence States

Block 05 uses these matrix fields consistently:

| Field | Allowed values | Meaning |
|---|---|---|
| `evidence_state` | `observed`, `not_assessed` | Whether committed evidence describes observed behavior or the row is unassessed. |
| `reason_code` | `run_artifact_available`, `no_run_artifact`, `missing_export`, `discovery_required`, `design_fixture_only`, `sanitization_pending`, `unsafe_to_run`, `external_verdict_available` | Why the row has its current evidence state. |
| `artifact_reference` | committed path or `none` | Evidence artifact supporting the row. `observed` requires a committed sanitized run artifact or evidence summary. |
| `external_verdict_ref` | committed path or `none` | Optional external verdict input with producer, origin, and policy reference when available. |

Synthetic fixtures may support schema or design coverage only. They must not move harness, model, stack, or customer pilot behavior to `observed`.

`pass`, `fail`, `warn`, `ready`, `blocked`, `supported`, `compatible`, and `unsupported` are not native `sdp-trace` pilot verdicts. If a harness or external reviewer emits those words, the run-card must record them as externally produced observations with producer, policy reference when available, artifact reference, and origin.

## Required Run-Card Content

Every Block 05 run-card must include:

- Scope: project, repository reference, source commit, service path, language, build system, harness, model family, and model version when available.
- Prompt: exact prompt text or a deterministic template with named variables.
- Expected artifacts:
  - `evidence-bundle.json`
  - `provenance-records.json`
  - `trace-snapshot.json`
  - `assessment-input.json` when preparing downstream policy handoff
  - sanitized command output or manual evidence note when raw export is unavailable
- Required provenance:
  - human operator
  - harness identity and version when available
  - model family and version when available
  - tool calls or command summaries when available
  - source commit or immutable source reference
  - redaction note and SHA-256 digest for committed artifacts
- `unbacked_claim` capture:
  - any model or harness claim that is not backed by inspected files, commands, or artifacts
  - any inferred build system claim contradicted by repository files
  - any statement that converts missing evidence into success
- Validation:
  - JSON parse/schema validation where a schema exists
  - `not_assessed` checks for missing evidence fields
  - matrix update rule proving no row moves to `observed` without a committed sanitized run artifact or evidence summary

## OpenCode Model Slices

The OpenCode run-card must cover these planned rows:

- OpenCode + MiniMax
- OpenCode + Kimi
- OpenCode + GLM

For each row, the first acceptable evidence state is `not_assessed` unless a committed sanitized run artifact exists. A row may not record model support, readiness, or compatibility from naming the model in a plan.

Required evidence includes model identity, harness identity, prompt text, source reference, produced artifacts, export limitations, and any `unbacked_claim` items.

## Harness Slices

The harness run-card must cover:

- Superpowers-style harnesses
- `gsd`
- `gsd2`
- Oh My OpenAgent

The run-card must separate harness behavior from model behavior. A model failure does not prove a harness failure; a harness export gap does not prove a model quality gap.

Required evidence includes rules/prompt location, tool log access, hook support, evidence export shape, manual fallback path, and known limitations.

## Kotlin+Bazel Slice

The Kotlin+Bazel path must close the current evidence-design gap, not the real behavior proof gap.

Block 05 must define what a real Kotlin+Bazel run must prove:

- Bazel ownership is scope-specific and based on `BUILD` or `BUILD.bazel` target evidence, `MODULE.bazel`, `WORKSPACE`, `WORKSPACE.bazel`, or equivalent Bazel source files tied to the assessed scope.
- `.bazelrc` is supporting configuration evidence only; it must not prove Bazel ownership by itself.
- Kotlin service-language detection is based on `.kt`, `.kts`, `kt_jvm_*`, or Kotlin compiler/toolchain rules tied to the assessed scope.
- Declared Kotlin dependencies may support ecosystem context, but they must not prove the assessed service is Kotlin without source or rule evidence.
- Maven or Gradle metadata inside a scoped Bazel target is dependency metadata only when scoped build evidence proves Bazel ownership.
- Scoped service assessment is preferred over root-level monorepo assessment.
- Build/test commands may remain `not_assessed` when they cannot be run safely or reproducibly.

The fixture placeholder must stay `not_assessed` until a real committed run artifact exists. The Kotlin+Bazel gap is not closed by a synthetic example alone.

## Customer Evidence Package Outline

The customer package must be an outline for a safe pilot handoff. It must list:

- pilot objective and scope
- required private input artifacts from the customer
- outputs produced by `sdp-trace`
- redaction and sanitization rules
- evidence package directory shape
- matrix update rules
- review and approval checkpoints
- residual `not_assessed` reporting
- handoff path to external policy consumers such as `sdp-gate`

Private customer inputs are never committed. Committed artifacts are sanitized summaries, hashes, redaction notes, and access-neutral references only. The outline must not include raw customer data or imply the pilot has already passed.

## UX Requirements

- A pilot operator must be able to run a row without reverse-engineering the repository.
- A CTO/CIO must see the difference between planned, unassessed, and evidence-backed rows in one table.
- Missing evidence must be first-class, not buried in notes.
- The package must make `unbacked_claim` items easy to capture and hard to launder into support/readiness/compatibility wording.
- Matrix rows must point to exact artifact references or explain why the state is `not_assessed`.

## Data Model Impact

No new schema is required for Block 05 unless implementation finds that existing evidence/provenance/trace schemas cannot represent the run-card outputs.

The preferred implementation is documentation and example-first:

- Use existing evidence event, provenance record, trace snapshot, and assessment input schemas where a committed JSON artifact is needed.
- Use Markdown run-cards for operator workflow and matrix evidence rules.
- Keep customer-specific values as placeholders or redacted references, not raw data.

If a new schema becomes necessary, it must be introduced as a separate reviewed task before implementation continues.

## Acceptance

Block 05 design is accepted when:

- The spec explicitly separates run-card recipes from observed behavior evidence and external compatibility verdicts.
- The Socratic dialogue records the main objections and resolutions.
- The implementation plan maps T027-T033 and T037 to concrete docs/examples/matrix changes.
- Pi review has been run on the spec artifacts, every valid finding is recorded in the committed review ledger, every valid finding is mirrored in Beads, and every spec-gate finding is closed before implementation.

Block 05 implementation may start after spec consensus is recorded and pi-review findings on the spec artifacts are closed.

## Delivery State

Current state before implementation:

1. Self-trace and local release proof are complete enough to start pilot evidence design.
2. External production trust remains `not_assessed`.
3. Customer pilot rows are not yet evidence-backed.
4. This block must produce executable run-card and evidence-package artifacts before any external customer compatibility verdict can be recorded.
