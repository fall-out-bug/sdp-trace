# Block 05 Implementation Plan: Customer Pilot Evidence Package

Status: accepted for implementation; spec review fixes closed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.22`

## Goal

Implement Phase 6 tasks T027-T033 and Phase 7 task T037 as a repository-readable pilot evidence package without native `sdp-trace` support, readiness, or compatibility verdicts.

## Consensus

Consensus is recorded for an evidence-package-first pilot design. Implementation proceeds only after Block 05 spec pi-review findings are recorded in the committed review ledger, mirrored in Beads, and closed. The selected scope is docs, run-cards, matrix updates, and safe fixture placeholders; real customer pilot execution is out of scope.

## File Responsibilities

- `docs/research/opencode-model-run-card.md`: operator run-card for OpenCode + MiniMax/Kimi/GLM slices.
- `docs/research/harness-run-card.md`: operator run-card for Superpowers-style, `gsd`, `gsd2`, and Oh My OpenAgent slices.
- `docs/research/kotlin-bazel-fixture-plan.md`: Kotlin+Bazel evidence requirements and fixture/run boundary.
- `docs/research/customer-pilot-evidence-package.md`: customer-safe pilot package outline.
- `docs/research/run-card-template.md`: generic template updated to avoid native `sdp-trace` pass/fail wording.
- `docs/jvm-bazel-guide.md`: Kotlin+Bazel-specific evidence and anti-heuristics.
- `examples/jvm-bazel/README.md`: fixture status and run instructions.
- `examples/jvm-bazel/evidence-bundle.json`: placeholder or evidence bundle that explicitly keeps real Kotlin+Bazel run behavior `not_assessed`.
- `docs/harness-compatibility-matrix.md`: evidence-state matrix for harness rows.
- `docs/model-compatibility.md`: evidence-state matrix for model rows.
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`: add Block 05 functional/success requirements if the current requirements are not specific enough.
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`: sync T027-T033/T037 and Block 05 review tasks after validation.
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/05-customer-pilot-review-ledger.md`: committed pi-review ledger with findings, severity, disposition, evidence, and optional Beads mirror IDs.
- Release proof artifacts under `examples/contract-foundation/` and `examples/self-trace/`: update only if validation shows manifest/self-attestation hashes must be regenerated after documentation changes.

## Task 1: Spec Review Gate

- Run pi review against:
  - `blocks/05-customer-pilot-evidence-package.md`
  - `blocks/05-customer-pilot-socratic.md`
  - `blocks/05-customer-pilot-implementation-plan.md`
- Record every valid finding in `blocks/05-customer-pilot-review-ledger.md`.
- Mirror every valid finding as a child Beads issue under `sdp-trace-cdn.22`.
- Fix every valid finding before implementation.
- Update the Block 05 spec status from `draft; pi review pending` to `accepted for implementation` only after the finding loop is closed and release proof status is either regenerated or explicitly recorded as stale.

Verification:

```bash
jq empty schema/*.json
git diff --check
npm run validate
rg -n "TBD|external_verdict_recorded|unsupported claim|pass / warn / fail" specs/001-sdp-trace-time-series-evidence-substrate/blocks/05-customer-pilot-evidence-package.md specs/001-sdp-trace-time-series-evidence-substrate/blocks/05-customer-pilot-socratic.md
rg -n "F[0-9]{3}.*open|implementation fix required|proof sync required" specs/001-sdp-trace-time-series-evidence-substrate/blocks/05-customer-pilot-review-ledger.md
```

Expected result: schema syntax and whitespace checks pass; active Block 05 spec files do not use deprecated row states. Full validation either passes or the review ledger explicitly records stale manifest/self-attestation proof as `not_assessed` pending the implementation artifact sync. Ledger entries with `implementation fix required` may remain until the implementation phase.

## Task 2: OpenCode Model Run-Card

- Create `docs/research/opencode-model-run-card.md`.
- Include separate rows for OpenCode + MiniMax, OpenCode + Kimi, and OpenCode + GLM.
- Each row starts as `not_assessed` unless a committed sanitized run artifact exists.
- Include prompt template, required artifacts, provenance fields, `unbacked_claim` capture, validation steps, and matrix update rules.
- State that naming the model is not observed behavior evidence and cannot produce a native support/readiness/compatibility outcome.

Verification:

```bash
rg -n "MiniMax|Kimi|GLM|not_assessed|unbacked_claim|artifact" docs/research/opencode-model-run-card.md
```

Expected result: all required slices and claim boundaries are present.

## Task 3: Harness Run-Card

- Create `docs/research/harness-run-card.md`.
- Include Superpowers-style harnesses, `gsd`, `gsd2`, and Oh My OpenAgent.
- Separate harness evidence from model behavior.
- Include evidence expectations for rules/prompt location, tool log access, hooks, evidence export, manual fallback, and known limitations.
- Ensure rows do not record observed behavior or external compatibility verdicts from planned evaluation.

Verification:

```bash
rg -n "Superpowers-style|gsd|gsd2|Oh My OpenAgent|tool log|hook|evidence export|not_assessed" docs/research/harness-run-card.md
```

Expected result: all requested harness slices and evidence gaps are present.

## Task 4: Kotlin+Bazel Fixture Plan and Guide

- Create `docs/research/kotlin-bazel-fixture-plan.md`.
- Update `docs/jvm-bazel-guide.md` with Kotlin+Bazel-specific evidence requirements:
  - scope-specific Bazel ownership files
  - Kotlin source/rule detection
  - Maven/Gradle metadata caveat
  - scoped service assessment
  - `not_assessed` rules for unsafe or unavailable commands
- Update `examples/jvm-bazel/README.md`.
- Update or add `examples/jvm-bazel/evidence-bundle.json` so the Kotlin+Bazel real run status remains explicitly `not_assessed` until committed run evidence exists.

Verification:

```bash
node -e "JSON.parse(require('fs').readFileSync('examples/jvm-bazel/evidence-bundle.json','utf8')); console.log('ok')"
rg -n "Kotlin|Bazel|not_assessed|BUILD.bazel|MODULE.bazel|Gradle|Maven|design_fixture_only" docs/jvm-bazel-guide.md docs/research/kotlin-bazel-fixture-plan.md examples/jvm-bazel/README.md examples/jvm-bazel/evidence-bundle.json
```

Expected result: JSON parses and docs preserve the Kotlin+Bazel gap boundary.

## Task 5: Customer Pilot Evidence Package Outline

- Create `docs/research/customer-pilot-evidence-package.md`.
- Include pilot objective, scope, customer inputs, `sdp-trace` outputs, package directory shape, redaction rules, validation commands, matrix update rules, review checkpoints, residual `not_assessed` reporting, and downstream policy handoff.
- State that the outline is not a completed pilot result.

Verification:

```bash
rg -n "private customer input|redaction|not_assessed|matrix|sdp-gate|raw customer data|validation|access-neutral" docs/research/customer-pilot-evidence-package.md
```

Expected result: package outline is executable and safe to commit.

## Task 6: Template and Matrix Updates

- Update `docs/research/run-card-template.md` to replace native `pass/warn/fail` wording with `observed`/`not_assessed` evidence states, reason codes, and explicit external verdict handling.
- Update `docs/harness-compatibility-matrix.md` with evidence state, reason code, artifact reference, external verdict reference, gap reason, and next required evidence columns.
- Update `docs/model-compatibility.md` with evidence state, reason code, artifact reference, external verdict reference, gap reason, and next required evidence columns.
- Keep planned rows `not_assessed` unless committed evidence exists.
- Add a deterministic matrix validation step or script that checks:
  - required matrix columns exist
  - `observed` rows have non-empty committed artifact references
  - synthetic fixtures cannot make behavior `observed`
  - native support/readiness/compatibility verdict tokens are either absent or explicitly external-origin

Verification:

```bash
rg -n "pass|fail|warn|blocked|readiness|ready|support|supported|compatible|compatibility|\"status\": \"pass\"|\"status\": \"warn\"|\"status\": \"fail\"" docs/research/run-card-template.md docs/harness-compatibility-matrix.md docs/model-compatibility.md docs/jvm-bazel-guide.md examples/jvm-bazel/evidence-bundle.json
```

Expected result: no native `sdp-trace` support/readiness/compatibility verdicts remain. Any occurrence is explicitly scoped to external origin, legacy file naming, or prohibited wording.

## Task 7: SpecKit Sync

- Update `spec.md` only if the current FR/SC set does not fully encode Block 05 evidence-state rules.
- Update `tasks.md`:
  - T027 through T033 complete only after their artifacts exist and pass validation.
  - T037 complete only after the customer package outline exists.
  - Add Block 05 review/validation tasks if needed for repo-observable process history.
- Do not mark real pilot observed behavior or external compatibility verdicts complete.

Verification:

```bash
rg -n "Block 05|T027|T028|T029|T030|T031|T032|T033|T037|not_assessed|observed" specs/001-sdp-trace-time-series-evidence-substrate
```

Expected result: SpecKit artifacts match implemented docs/examples.

## Task 8: Validation, Release Artifact Sync, and Review

- Run:

```bash
npm run validate
jq empty schema/*.json
git diff --check
```

- If validation reports stale manifest, hash, or self-attestation artifacts, regenerate only the required local proof artifacts and keep external production trust `not_assessed`. If regeneration is intentionally deferred during the spec gate, record the stale proof state as `not_assessed` in the committed review ledger.
- Run implementation pi review after validation passes.
- Register every valid implementation finding in Beads, including minor/P3 items.
- Fix every valid finding and rerun validation.
- Close `sdp-trace-cdn.22` only after docs, fixtures, matrices, validation, pi review, and finding closure are complete.
