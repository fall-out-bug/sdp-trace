# Block 09 Implementation Plan: Flight Recorder Trust Kernel

Status: ready for agent handoff; implementation not started
Parent: `09-flight-recorder-trust-kernel.md`

## Goal

Implement the minimal recorder trust kernel that can prove local chain structure, witnessed chain anchoring, explicit late-attach gaps, requirement supersession, redaction states, and reviewer-readable query output without adding code-quality verdicts or harness-specific runtime dependencies.

## Execution Rules

- No Feature Flag / Entitlements demo work starts in this block.
- No broad OpenCode, GSD, MiniMax, Kotlin, or Bazel support claim is allowed.
- Local chain consistency must not be described as audit-grade trust.
- Witnessed mode must fail when witness evidence is missing or mismatched.
- Every missing observation is `not_assessed` or `cannot_verify` with a reason.
- Every trust-affecting claim must be verifier-derived.
- Every implementation slice needs tests before or with behavior changes.
- Review findings are recorded in a Block 09 review ledger before closure.

## Proposed File Responsibilities

- `schema/flight-recorder-event.schema.json`: event schema for Block 09 recorder events.
- `schema/flight-recorder-run.schema.json`: run manifest / closure schema.
- `schema/flight-recorder-witness.schema.json`: witness entry schema.
- `examples/flight-recorder/`: positive and negative fixtures.
- `scripts/verify-flight-recorder.mjs` or `.sh`: verifier for local/witnessed/forensic profiles.
- `scripts/test-flight-recorder.sh`: positive and negative verifier tests.
- `scripts/query-flight-recorder.mjs` or `.sh`: reviewer query surface.
- `docs/flight-recorder.md`: product-facing recorder semantics and limits.
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/09-flight-recorder-review-ledger.md`: review and disposition.
- `scripts/validate-contracts.sh`: add only committed fixture validation after the verifier and tests exist.
- `package.json`: add scripts only after command names are stable.

## Slice 0: Spec and Review Ledger Gate

Purpose: ensure implementation starts from an accepted, reviewable spec.

Tasks:

- Create Block 09 review ledger.
- Record executive Socratic findings as spec-gate findings.
- Close or explicitly leave blocking every critical/major design finding.
- Confirm local-only chain is not called trusted.

Verification:

```bash
rg -n "F09-C01|F09-C02|F09-C03|F09-C04|F09-C05|F09-C06" specs/001-sdp-trace-time-series-evidence-substrate/blocks/09-flight-recorder-socratic.md
rg -n "local_development_recorder|witnessed_run_recorder|externally_witnessed_run" specs/001-sdp-trace-time-series-evidence-substrate/blocks/09-flight-recorder-trust-kernel.md
```

Exit criteria:

- Spec-gate findings are visible.
- Implementation starts only after unresolved blockers are accepted as blockers or resolved in spec.

## Slice 1: Event and Witness Schemas

Purpose: define the canonical data contract before verifier code.

Tasks:

- Add event schema with declared canonicalization fields.
- Add run manifest schema.
- Add witness entry schema.
- Add positive local chain fixture.
- Add positive witnessed chain fixture.
- Add negative fixtures:
  - missing `prev_event_hash`
  - mismatched event hash
  - unsupported schema version
  - missing witness for witnessed profile
  - witness chain-head mismatch

Verification:

```bash
jq empty schema/*.json
node scripts/validate-json-schema.mjs schema/flight-recorder-event.schema.json examples/flight-recorder/local-positive/events/000-run-started.json
node scripts/validate-json-schema.mjs schema/flight-recorder-witness.schema.json examples/flight-recorder/witnessed-positive/witness.json
```

Exit criteria:

- Schema fixtures validate.
- Negative fixtures fail for named reasons once the verifier exists.

## Slice 2: Chain Verifier

Purpose: make event mutation, deletion, and reordering detectable.

Tasks:

- Implement verifier profile `flight_recorder_local`.
- Verify event sequence, event hash, previous hash, schema version, and run closure.
- Emit structured states:
  - `event_chain_structurally_valid`
  - `source_baseline_recorded`
  - `task_locked`
  - `run_closed`
- Add tests that mutate an event payload and expect verifier failure.
- Add tests that delete an event and expect verifier failure.
- Add tests that reorder events and expect verifier failure.

Verification:

```bash
scripts/test-flight-recorder.sh --slice chain
scripts/verify-flight-recorder.mjs --profile local examples/flight-recorder/local-positive
```

Exit criteria:

- Local profile passes positive fixture.
- Tamper fixtures fail with stable machine-readable reasons.

## Slice 3: Witness Verification

Purpose: distinguish local consistency from witnessed accountability evidence.

Tasks:

- Implement profile `flight_recorder_witnessed`.
- Verify witness entry binds:
  - run id
  - source baseline hash
  - task hash
  - recorder version
  - final chain head
  - witness timestamp
- Add local file witness fixture as first implementation witness.
- Add negative fixtures for missing witness, wrong run id, wrong chain head, and wrong task hash.

Verification:

```bash
scripts/test-flight-recorder.sh --slice witness
scripts/verify-flight-recorder.mjs --profile witnessed examples/flight-recorder/witnessed-positive
```

Exit criteria:

- Witnessed mode fails without witness evidence.
- Witness mismatch cannot be hidden by a locally valid chain.

## Slice 4: Late Attach and Requirement Supersession

Purpose: make expectation history immutable from attachment forward and honest before attachment.

Tasks:

- Add late-attach fixture with pre-attachment history marked `not_assessed`.
- Add full-run fixture with no late-attach gap.
- Add requirement supersession fixture:
  - original task locked
  - command event occurs
  - requirement superseded
  - later command references new task event
- Add negative fixture that changes task payload after command evidence.

Verification:

```bash
scripts/test-flight-recorder.sh --slice expectations
scripts/query-flight-recorder.mjs --query requirement-timeline examples/flight-recorder/requirement-supersession
```

Exit criteria:

- Late attach is visible in verifier output.
- Task rewrite fails chain verification.
- Supersession is queryable without mutating the original task.

## Slice 5: Evidence Retention and Redaction States

Purpose: avoid false forensic claims and unsafe secret handling.

Tasks:

- Add evidence retention modes:
  - `digest_only`
  - `sanitized_excerpt`
  - `encrypted_raw_ref`
  - `external_artifact_ref`
  - `not_assessed`
- Add redaction states:
  - `not_required`
  - `redacted`
  - `sealed_raw_available`
  - `not_assessed`
  - `cannot_verify`
- Implement verifier rules for unresolved redaction.
- Add fixtures:
  - safe sanitized command output
  - digest-only output accepted in local mode
  - digest-only critical output rejected in forensic mode
  - unresolved redaction rejected

Verification:

```bash
scripts/test-flight-recorder.sh --slice redaction
scripts/verify-flight-recorder.mjs --profile forensic examples/flight-recorder/forensic-positive
```

Exit criteria:

- Redaction cannot silently erase evidence.
- Digest-only evidence is explicitly low-forensic unless profile accepts it.

## Slice 6: Query Surface

Purpose: make the recorder useful to technical executive/reviewer workflows without policy verdicts.

Tasks:

- Implement query commands for:
  - run summary
  - provenance summary
  - late-attach gaps
  - requirement timeline
  - command timeline
  - file mutation summary
  - test evidence summary
  - redaction unresolved
  - witness state
- Ensure queries emit facts and verifier states, not pass/fail readiness decisions.

Verification:

```bash
scripts/query-flight-recorder.mjs --query run-summary examples/flight-recorder/witnessed-positive
scripts/query-flight-recorder.mjs --query gaps examples/flight-recorder/late-attach
scripts/query-flight-recorder.mjs --query witness-state examples/flight-recorder/witnessed-positive
```

Exit criteria:

- Reviewer can answer who/what/when/source/task/model/commands/files/tests/gaps from queries.
- Queries do not emit readiness, support, compatibility, or quality verdicts.

## Slice 7: Validation Wiring and Docs

Purpose: make Block 09 discoverable and reproducible.

Tasks:

- Add stable package scripts after command names settle.
- Wire committed fixture validation into `scripts/validate-contracts.sh`.
- Add `docs/flight-recorder.md`.
- Update `schema/README.md` with flight-recorder schema notes.
- Update `tasks.md` only for actual completed Block 09 work.
- Keep external production trust open unless real external witness evidence exists.

Verification:

```bash
npm run validate
scripts/test-flight-recorder.sh
git diff --check
```

Exit criteria:

- Baseline validation includes committed Block 09 fixtures.
- Docs state local/witnessed/external limits clearly.

## Slice 8: Implementation Review and Handoff to Demo

Purpose: prevent demo work from hiding recorder-kernel gaps.

Tasks:

- Run strict review for schemas, verifier, redaction, witness, and query surface.
- Record findings in Block 09 review ledger.
- Fix valid critical/major findings or leave them blocking.
- Define activation gate for the later Feature Flag / Entitlements demo.

Verification:

```bash
scripts/test-flight-recorder.sh
npm run validate
git diff --check
```

Exit criteria:

- Block 09 can be handed to demo agents only if witnessed profile and tamper fixtures pass.
- Any missing external witness or model identity proof remains explicit and cannot be polished away in demo docs.
