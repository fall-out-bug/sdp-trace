# Block 15: Signed Checkpoint And Replay Resistance

Status: spec delta and implementation plan.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13b-capture-boundary-state-dx-baseline.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/14-gate-contract-explain-override.md`
- `schema/flight-recorder-witness.schema.json`
- `schema/ci-witness.schema.json`

## Goal

Make a run checkpoint independently replayable enough for a verifier to detect
post-hoc mutation, source/run replay, and stale checkpoint reuse.

Block 15 adds a signed checkpoint artifact and verifier states for checkpoint
shape, signature, chain binding, source binding, nonce binding, monotonic
sequence, signer authority, and replay freshness. It does not implement
protected merge enforcement or audit-grade external witness trust.

## Problem

Block 14 can say that a protected future gate cannot verify before signed
checkpoint evidence exists. The missing product surface is now specific:

- what exactly is signed;
- how the signed payload binds to the run chain head, source snapshot, task,
  run id, nonce, and checkpoint sequence;
- how a verifier reports replay to another run/source/task;
- how a verifier reports an old checkpoint being reused as current evidence;
- how signer authority is represented without pretending local key possession
  is an external trust boundary.

Without this block, teams can still produce local chain hashes and CI witness
JSON, but a reviewer cannot tell whether a checkpoint belongs to this run or
was copied from a previous run.

## Non-Goals

- No native merge, release, readiness, degradation, override approval, or risk
  acceptance decision.
- No protected fail-closed gate behavior.
- No Sigstore/Rekor integration in the active product path.
- No dependency on GitHub, Claude, Codex, OpenCode, Beads, or a harness runtime.
- No claim that a local private key is isolated from the recorded agent.
- No external timestamp authority or append-only transparency log.
- No raw stdout, stderr, prompt, model response, source snippet, credential, or
  token persistence.

## Product Boundary

Block 15 may emit verifier-derived facts:

- checkpoint payload shape and canonical digest;
- detached signature validity for the implemented profile;
- run id, chain head, source snapshot, task hash, nonce, and sequence binding;
- signer authority state from an explicit trusted-checkpoint policy;
- replay/freshness state from verifier-supplied context;
- local trust cap when signer isolation or external witness evidence is absent.

Block 15 must not emit a native policy decision. External consumers may use the
facts to block, permit, escalate, or accept risk.

## Checkpoint Artifact

Schema file:

```text
schema/signed-checkpoint.schema.json
```

Minimal artifact shape:

```json
{
  "schema_version": "block15-signed-checkpoint-v1",
  "checkpoint_id": "checkpoint-001",
  "run_id": "run-abc",
  "sequence": 0,
  "checkpoint_time": "2026-05-06T00:00:00Z",
  "profile": "sdp-trace-checkpoint/ed25519-detached-v1",
  "canonicalization": {
    "algorithm": "json-canonicalization-scheme",
    "version": "1.0.0"
  },
  "hash_algorithm": "sha256",
  "payload": {
    "run_id": "run-abc",
    "run_nonce": "run-nonce-run-abc",
    "event_chain_head": "<sha256>",
    "event_count": 5,
    "source_snapshot_digest": "<sha256>",
    "source_snapshot_state": "git_tree_clean",
    "task_hash": "<sha256>",
    "contract_digest": "<sha256>",
    "previous_checkpoint_digest": "",
    "replay_context": {
      "repository": "not_assessed",
      "ref": "not_assessed",
      "commit_sha": "not_assessed"
    }
  },
  "payload_digest": "<sha256>",
  "signature": {
    "algorithm": "ed25519",
    "signature": "<base64>",
    "public_key": "<base64>"
  },
  "signer_identity": {
    "signer_id": "local-dev",
    "authority": "local_development_key",
    "key_isolation": "not_assessed"
  }
}
```

Rules:

- `payload_digest` is computed from canonical `payload` only.
- The detached signature signs the canonical `payload` bytes, not the outer
  JSON envelope.
- `run_id`, `event_chain_head`, `event_count`, `source_snapshot_digest`,
  `source_snapshot_state`, `task_hash`, `contract_digest`, and `run_nonce`
  must match the selected run artifact when verifier context is available.
- `sequence` is monotonic within the supplied checkpoint set. Duplicate,
  missing, or descending sequences fail the checkpoint-set verifier.
- `previous_checkpoint_digest` is empty only for sequence `0`; later
  checkpoints must reference the previous checkpoint payload digest.
- A checkpoint copied to another run, source snapshot, task, contract, or nonce
  fails binding verification.
- Local `ed25519` is a development verifier profile. It can prove content was
  signed by the matching private key; it cannot prove signer isolation.
- `key_isolation: not_assessed` keeps protected and audit-grade trust capped.

## Trusted Checkpoint Policy

Schema file:

```text
schema/trusted-checkpoint-policy.schema.json
```

The policy names allowed signer identities and the expected authority boundary.
Initial authority values:

- `local_development_key`
- `ci_isolated_job`
- `external_witness_service`

Rules:

- Missing policy leaves signer authority `not_assessed`.
- Policy mismatch is `fail`.
- `local_development_key` can verify signature integrity but cannot raise trust
  above local signed evidence.
- `ci_isolated_job` requires explicit CI identity binding; if that context is
  unavailable, signer authority is `cannot_verify`.
- `external_witness_service` is reserved for later external witness work and
  remains `not_integrated` in Block 15.

## Verifier Result

Block 15 writes checkpoint verification as machine-readable facts:

```json
{
  "schema_version": "block15-checkpoint-verification-v1",
  "checkpoint_id": "checkpoint-001",
  "run_id": "run-abc",
  "result": "pass",
  "trust_scope": "local_signed",
  "signature_state": "pass",
  "chain_binding_state": "pass",
  "source_binding_state": "pass",
  "nonce_binding_state": "pass",
  "sequence_state": "pass",
  "signer_authority_state": "not_assessed",
  "replay_freshness_state": "not_assessed",
  "reasons": []
}
```

Allowed states:

- `pass`
- `fail`
- `cannot_verify`
- `not_assessed`
- `not_integrated`

Trust scopes:

- `local_signed`
- `ci_signed`
- `external_witnessed`
- `untrusted_shape_only`

Rules:

- Overall result is `fail` if any binding, signature, sequence, or authority
  check fails.
- Overall result is `cannot_verify` if required verifier context is missing and
  no check failed.
- Overall result is `not_assessed` only when no verification profile is selected
  or implemented.
- Overall result `pass` is scoped. `pass + local_signed` is not protected
  enforcement and not audit-grade external trust.

## User-Facing Commands

### `sdp-trace checkpoint create`

Usage:

```text
sdp-trace checkpoint create --run <run-dir> --out <checkpoint.json> --private-key <ed25519-private-key.json> --signer-id <id> [--policy <policy.json>]
```

Behavior:

1. Load and verify the local run chain before signing.
2. Build the canonical checkpoint payload from run metadata and recorder
   events.
3. Sign only the canonical payload bytes.
4. Write one checkpoint artifact.
5. Do not print raw command args, stdout, stderr, prompts, source snippets,
   credentials, OIDC request tokens, or model responses.

### `sdp-trace checkpoint verify`

Usage:

```text
sdp-trace checkpoint verify --run <run-dir> --checkpoint <checkpoint.json> [--policy <policy.json>]
```

Behavior:

1. Validate checkpoint shape.
2. Verify payload digest and detached signature.
3. Recompute run chain binding from the selected run.
4. Check run id, nonce, source snapshot, task hash, contract digest, event
   count, and event chain head.
5. Check signer identity against the selected policy when supplied.
6. Emit deterministic JSON result and write no files unless an explicit output
   flag is added in a later block.

## Acceptance Criteria

1. A valid local Ed25519 checkpoint over an untampered run verifies as
   `pass + local_signed` with signer authority `not_assessed` unless a policy is
   supplied.
2. Tampering with a signed checkpoint payload fails payload digest or signature
   verification.
3. Replaying a checkpoint against a different run id, source snapshot, task,
   contract digest, nonce, event count, or chain head fails binding
   verification.
4. A checkpoint set with duplicate, missing, or descending sequence numbers
   fails monotonic sequence verification.
5. A signer not allowed by the selected trusted-checkpoint policy fails signer
   authority verification.
6. A local-development signer cannot upgrade `audit_grade_gate` or protected
   gate output to pass.
7. Checkpoint create and verify output does not persist or print raw sensitive
   fields.
8. All behavior is implemented in Go with test-first coverage.

## Implementation Plan

### Slice A: Schemas And Spec Wiring

Files:

- `schema/signed-checkpoint.schema.json`
- `schema/trusted-checkpoint-policy.schema.json`
- `schema/README.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`

Tasks:

- define signed checkpoint and trusted policy schemas;
- add FRs and tasks for checkpoint binding, signature verification, replay
  checks, and no-overclaim gate behavior;
- keep external witness and protected enforcement out of scope.

### Slice B: Checkpoint Domain And Verifier

Files:

- `internal/checkpoint/checkpoint.go`
- `internal/checkpoint/checkpoint_test.go`
- `internal/trace/types.go`
- `internal/trace/store.go`

Tasks:

- derive canonical checkpoint payload from a run artifact;
- sign and verify detached Ed25519 signatures;
- compute payload digest deterministically;
- validate run id, nonce, source snapshot, task hash, contract digest, event
  count, and chain head bindings;
- evaluate monotonic checkpoint sequences.

### Slice C: CLI

Files:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/main_test.go`

Tasks:

- add `checkpoint create`;
- add `checkpoint verify`;
- add usage text;
- add secret-like negative output assertions.

### Slice D: Gate Integration Without Enforcement

Files:

- `internal/demo/demo.go`
- `internal/demo/demo_test.go`
- `schema/gate-result.schema.json`

Tasks:

- surface signed-checkpoint verification facts only when supplied by a later
  explicit gate flag or artifact path;
- keep `protected_future` and `audit_grade_gate` `cannot_verify` without an
  external policy consumer;
- prove local signed checkpoint evidence does not convert protected gate state
  to pass.

### Slice E: Fixtures, Verification, And Review

Files:

- `examples/block15-checkpoint/`
- retired research artifact
- retired research artifact

Tasks:

- add committed positive and negative fixtures;
- run Go-first verification and schema checks;
- record strict review and pi-review disposition before PR closure.

## Review Ledger

| Area | State | Disposition |
|---|---|---|
| Spec delta | pass | Captured in this file plus FR-065 through FR-070 and T116-T125. |
| Signed checkpoint schema | pass | Implemented in `schema/signed-checkpoint.schema.json`. |
| Trusted checkpoint policy schema | pass | Implemented in `schema/trusted-checkpoint-policy.schema.json`; policy public-key binding is required. |
| Signature verification | pass | Implemented in Go with detached Ed25519 signature and canonical payload digest checks. |
| Replay binding checks | pass | Run id, nonce, source snapshot, task hash, contract digest, event count, and chain head mismatches fail or cannot verify. |
| Monotonic sequence checks | pass | Duplicate, missing, descending, and forged checkpoint set cases are covered by Go tests. |
| Gate enforcement | not_integrated | Explicit non-goal; external policy consumer required. |
| External witness trust | not_integrated | Reserved for later external witness profile. |

## No-Overclaim Notes

- A valid local signature proves only that the matching key signed the payload.
- Local key possession is not signer isolation.
- `local_signed` is stronger than unsigned local JSON but weaker than
  `ci_signed` and `external_witnessed`.
- A checkpoint copied from another run must fail replay binding even if its
  signature is valid.
- `audit_grade_gate` remains `cannot_verify` until external witness and policy
  consumer work exists.
