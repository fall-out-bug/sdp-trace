# Block 14: Gate Contract, Explain, And Native Override Event

Status: spec delta and implementation plan.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13b-capture-boundary-state-dx-baseline.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/11-demo-report-gate.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/12-ci-witness-adoption.md`

## Goal

Make advisory gate output operationally useful at the CI boundary without
claiming protected enforcement, audit-grade trust, or native policy ownership.

Block 14 tightens the contract between expected evidence, required runs,
witness binding, emergency override records, and user-facing explanation.
It does not decide merge, release, readiness, degradation, or risk acceptance.
Those decisions remain external policy-consumer concerns.

## Problem

The current gate path can report local contract evidence and CI witness posture,
but the next buyer-visible questions are still weak:

- Was an expected run absent, or merely unmatched?
- Was the CI witness for the same repository, ref, commit, run id, and artifact
  digest set?
- If an emergency override happened, where is the recorded evidence and who
  requested it?
- Why did the gate produce `fail` versus `cannot_verify`, and what should the
  developer or platform owner do next?

Without this block, teams can still use generated gate JSON, but they will need
manual interpretation exactly where the product should reduce ambiguity.

## Non-Goals

- No protected fail-closed enforcement.
- No signed checkpoint, DSSE, in-toto, Sigstore/Rekor, or external transparency
  witness implementation.
- No native `sdp-trace` approval, release, readiness, degradation, or risk
  acceptance decision.
- No GitHub-specific gate core beyond reading the existing generic CI witness
  artifact shape.
- No harness-specific required-run semantics.
- No raw stdout, stderr, prompt, model response, source snippet, credential, or
  token persistence.

## Product Boundary

Block 14 may emit verifier-derived facts:

- local required evidence state;
- required-run presence or absence;
- CI witness binding state;
- override record presence and linkage;
- advisory gate explain output;
- audit-grade state as `cannot_verify` unless later signed/external evidence
  exists.

Block 14 must not emit a native policy decision. External consumers may use the
facts to block, permit, escalate, or accept risk.

## Contract Delta

### Required Runs

The evidence contract may declare required runs in addition to required
evidence. Required runs are observation requirements, not enforcement policy.

Minimal shape:

```json
{
  "required_runs": [
    {
      "id": "verification_run",
      "wrapper_name": "verification-run",
      "required_evidence": ["verification_run_observed"],
      "profile": "observation"
    }
  ]
}
```

Rules:

- `id` is a stable contract-local identifier.
- `wrapper_name` matches the safe wrapper descriptor emitted by the recorder.
- `required_evidence` references top-level contract-declared evidence ids.
- `profile` is `observation`, `advisory_ci`, or `protected_future`.
- `protected_future` must evaluate as `cannot_verify` until Block 15 or later
  supplies signed checkpoint evidence and an external policy consumer enforces
  it.

Missing required runs produce `missing_telemetry`. A run that exists but cannot
be bound to the contract produces `cannot_verify` or `fail` depending on whether
the verifier lacks data or finds contradictory data.

Run-state rules:

| Case | Required-run state |
|---|---|
| No matching run exists | `missing_telemetry` |
| Matching run exists and required evidence is present | `pass` |
| Matching run exists but verifier cannot inspect required fields | `cannot_verify` |
| Matching run exists and contradicts required evidence | `fail` |
| Required run declares `profile: protected_future` | `cannot_verify` until a later signed-checkpoint profile exists |

Ordering rule: required runs are emitted in contract declaration order; observed
runs that are not referenced by a required run are emitted after required runs
by run directory name.

### Native Override Event

Block 14 adds a portable trace event type:

```text
policy_override_requested
```

Required payload fields:

- `override_id`
- `producer`
- `origin`
- `requested_by`
- `reason`
- `source_ref`
- `scope`
- `created_at`

Optional payload fields:

- `external_reference`
- `approver`
- `expires_at`
- `affected_required_runs`
- `affected_evidence`

Rules:

- `origin` is `native_cli` or `external_reference`.
- `producer` names the CLI, adapter, or external system that produced the
  override request record.
- The event records that an override was requested or externally referenced.
- It never converts missing evidence to `pass`.
- It never upgrades `audit_grade_gate`.
- Override presence does not change any required-run or required-evidence
  state; `affected_required_runs` and `affected_evidence` are informational
  links for external policy consumers.
- If `affected_required_runs` or `affected_evidence` names an id that is not
  declared in the selected contract, the override record remains visible but
  its validation state is `cannot_verify`.
- If approval is external, `sdp-trace` records the reference and producer; it
  does not decide whether the approval is acceptable.
- If the event is present but malformed, the override state is `cannot_verify`.
- `policy_override_requested` is a flight-recorder chain event when recorded by
  the native recorder and inherits canonical payload digest, `prev_event_hash`,
  `event_hash`, redaction state, and optional witness reference from the Block
  09 event schema. External-reference imports that cannot prove chain linkage
  remain visible but are capped at `agent_reported` or `cannot_verify` as
  appropriate.

Minimal CLI shape:

```text
sdp-trace override request --out <run-dir> --id <override-id> --by <actor> --reason <text> --source-ref <ref> --scope <scope> [--external-reference <ref>]
```

The command appends a `policy_override_requested` event to the selected run
artifact. It must reject an empty actor, reason, source reference, or scope.

## User-Facing Commands

### `sdp-trace gate`

Block 14 extends gate output, but keeps it advisory.

Required output fields:

- `schema_version`
- `generated_at`
- `local_gate`
- `ci_witness_gate`
- `audit_grade_gate`
- `gate_mode`
- `trust_cap`
- `required_runs`
- `required_evidence`
- `observed_evidence`
- `witness_bindings`
- `override_requests`
- `gate_conditions`
- `reasons`
- `next_actions`
- `runs`

Allowed values:

- gate states: `pass`, `fail`, `cannot_verify`, `not_assessed`
- `gate_mode`: `observation`, `advisory_ci`, `protected_future`
- `trust_cap`: `local_observed`, `ci_witnessed`, `external_witnessed`
- required-run states: `pass`, `fail`, `cannot_verify`,
  `missing_telemetry`, `not_assessed`

`required_runs` shape:

- `id`
- `wrapper_name`
- `profile`
- `state`
- `matched_run_id`
- `reasons`

`gate_conditions` shape:

- `id`
- `state`
- `reason`

Initial condition ids:

- `all_required_runs_present`
- `all_required_evidence_observed`
- `ci_witness_bound_when_required`
- `audit_grade_external_witness_present`

`runs` is the observed run summary inherited from Block 11/12 output. Required
run evaluation lives under `required_runs`; `runs` remains useful for debugging
unmatched or extra observed runs.

Witness binding rules:

| Case | CI witness state |
|---|---|
| Required binding field absent | `cannot_verify` |
| All available binding fields match | `pass` |
| Repository, ref, commit, run id, or artifact digest contradicts current gate input | `fail` |

Partial artifact data is `cannot_verify`, not `fail`, unless a present digest
contradicts the current gate input.

Exit behavior:

- exit `0` only when the local/advisory gate state is `pass`;
- exit `1` for `fail`;
- exit `3` for `cannot_verify` when no `fail` state is present;
- usage errors remain exit `2`.

Exit code is determined from gate dimensions in this order:

1. any applicable emitted gate dimension is `fail` -> exit `1`;
2. otherwise any applicable emitted gate dimension is `cannot_verify` -> exit
   `3`;
3. otherwise the emitted advisory gate dimensions passed -> exit `0`.

Block 14 has no user-selectable gate-dimension filter. A later block may add
one, but this block evaluates every applicable emitted gate dimension.

`audit_grade_gate: cannot_verify` in Block 14 means the audit-grade profile is
not implemented for this evidence set; it is not a native release decision.

### `sdp-trace gate explain`

Usage:

```text
sdp-trace gate explain --gate-result <gate-result.json>
```

Behavior:

1. Load a Block 14 gate-result artifact.
2. Print deterministic human-readable explanation.
3. Include required-run gaps, witness binding gaps, override records, trust cap,
   and next action hints.
4. Never print raw command arguments, stdout, stderr, prompts, source snippets,
   credentials, OIDC request tokens, or model responses.

Next actions are remediation hints only. They may say to supply missing input,
rerun through the required observer, rerun in CI, inspect the policy owner, or
lower the claim. They must not recommend merge, release, approval, readiness,
or risk acceptance.

### `sdp-trace gate preview`

Usage:

```text
sdp-trace gate preview --contract <contract.json> [--witness <witness.json>] <runs-root-or-run-dir>
```

Behavior:

1. Show what gate-relevant fields would be read.
2. Show selected mode and trust cap.
3. Show required runs and evidence ids.
4. Show whether a witness artifact can be inspected locally.
5. Report binding mismatches detectable from available local witness data, but
   do not emit a gate verdict.
6. Do not write report or gate artifacts.
7. Do not claim that the gate will pass.

## Acceptance Criteria

1. A contract with a missing required run emits `missing_telemetry` and the
   advisory gate does not pass.
2. A CI witness for a different repository, ref, commit, run id, or artifact
   digest set emits `cannot_verify` or `fail` with a deterministic reason.
3. A valid override request is visible in gate output and explain output, but
   it does not change missing evidence to pass and does not change
   `audit_grade_gate` away from `cannot_verify`.
4. `gate explain` gives a developer or platform owner the next concrete action
   for `missing_telemetry`, `cannot_verify`, stale witness, source mismatch,
   and override-present states.
5. `gate preview` is read-only and deterministic for identical inputs.
6. Gate, explain, and preview output do not persist or print raw sensitive
   fields.
7. All behavior is implemented in Go with test-first coverage.

## Implementation Plan

### Slice A: Contract And State Model

Files:

- `schema/gate-result.schema.json`
- `internal/trace/contract.go`
- `internal/trace/types.go`
- `internal/demo/demo.go`
- `internal/demo/demo_test.go`

Tasks:

- add or extend Draft 2020-12 gate-result schema for Block 14 output;
- parse optional `required_runs`;
- add gate-mode and required-run state constants;
- evaluate required-run presence before evidence completeness;
- add absent-run, unmatched-run, and protected-future tests.

### Slice B: Witness Binding Checks

Files:

- `internal/demo/demo.go`
- `internal/demo/demo_test.go`
- `internal/witness/witness.go`
- `internal/witness/witness_test.go`

Tasks:

- compare witness source repository, ref, commit, run artifact digests, and
  report artifact digests against the current gate input when available;
- distinguish incomplete witness data from contradictory witness data;
- add negative tests for source mismatch, stale witness, and artifact mismatch.

### Slice C: Override Event

Files:

- `internal/trace/event.go`
- `internal/trace/types.go`
- `internal/recorder/recorder.go`
- `internal/demo/demo.go`
- `cmd/sdp-trace/main.go`
- related tests.

Tasks:

- add `policy_override_requested` event type;
- add one explicit CLI action or external-reference path to record override
  requests;
- include override state in gate output and explain output;
- add malformed override and valid override tests.

### Slice D: Explain And Preview

Files:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/main_test.go`
- `internal/demo/demo.go`
- `internal/demo/demo_test.go`

Tasks:

- add `gate explain`;
- add `gate preview`;
- add deterministic ordering for reasons, next actions, required runs, witness
  bindings, and overrides;
- add negative leak tests with secret-like markers.

Ordering rule: emit contract-declared items in contract order, observed extra
runs by run directory name, witness bindings by binding id, overrides by
`created_at` then `override_id`, and reasons/next actions by severity
(`fail`, `cannot_verify`, `missing_telemetry`, `not_assessed`, `pass`) then id.

### Slice E: Fixtures, Validation, And Review

Files:

- `examples/contract-foundation/` or a Block 14 fixture directory chosen during
  implementation;
- repository validation wiring;
- Block 14 review ledger.

Tasks:

- add committed fixtures for missing required run, unmatched run, stale witness,
  source mismatch, artifact mismatch, valid override request, malformed override
  request, and protected-future requirement;
- wire fixture validation into Go-first repository validation;
- record strict review and pi-review disposition before PR closure.

## Review Ledger

| Area | State | Disposition |
|---|---|---|
| Spec delta | drafted | Requires strict review before implementation closure. |
| Required-run semantics | pass | Covered by Go tests for missing required runs and `protected_future` non-pass behavior. |
| Witness binding | pass | Covered by Go test for repository mismatch and committed CI witness fixture. |
| Override event | pass | Covered by CLI test appending `policy_override_requested` chain event with producer/origin. |
| Explain output | pass | Covered by deterministic CLI test and secret-like negative assertion. |
| Preview output | pass | Covered by read-only CLI test and secret-like negative assertion. |
| Protected enforcement | not_integrated | Explicit non-goal until later block. |
| Audit-grade gate | cannot_verify | Requires later signed/external witness profile. |

## No-Overclaim Notes

- Advisory gate output is useful, but it is not protected enforcement.
- A CI witness JSON file is not external production trust.
- An override record is not approval unless an external policy consumer says so.
- Missing telemetry remains missing even when an override exists.
- Dirty-worktree verification is local structural evidence only.
