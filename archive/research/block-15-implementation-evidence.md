# Block 15 Implementation Evidence

Date: 2026-05-06

Scope:

- signed checkpoint and trusted checkpoint policy schemas;
- local detached Ed25519 checkpoint creation and verification;
- payload digest, signature, run binding, nonce binding, source/task/contract
  binding, and chain-head verification;
- monotonic checkpoint-set verification;
- signer policy public-key binding;
- `checkpoint create` and `checkpoint verify` CLI commands;
- no-overclaim gate regression for local signed checkpoints.

## Verification

Live local structural checks:

- `rtk go test ./...` -> pass, 71 tests across 11 packages.
- `rtk jq empty schema/*.json examples/block15-checkpoint/*.json` -> pass.

These checks are local structural evidence in a dirty worktree. They are not
source-bound release proof and do not establish protected enforcement or
external production trust.

## Implementation Notes

- Checkpoint signatures cover canonical checkpoint payload bytes, not the
  outer JSON envelope.
- `payload_digest` is recomputed from canonical payload bytes before a
  checkpoint binding is trusted.
- Missing run nonce returns `cannot_verify`; nonce mismatch returns a binding
  failure.
- `sequence > 0` requires `previous_checkpoint_digest`; sequence `0` must not
  declare one.
- `VerifySet` verifies each checkpoint signature/digest/binding before
  accepting sequence linkage.
- `Verify` rejects mismatched checkpoint schema version, profile, hash
  algorithm, and canonicalization constants before using checkpoint fields.
- Trusted checkpoint policy requires a public-key binding for the signer.
- Local development signatures remain capped at `local_signed` and do not
  upgrade `protected_future` or `audit_grade_gate`.

## Residual State

- External witness trust remains `not_integrated`.
- CI-isolated signer authority remains `cannot_verify` until CI binding context
  is implemented.
- Source-bound release proof regeneration is not performed in this
  implementation chunk.
