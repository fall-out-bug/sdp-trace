# Block 16: Protected Gate Enforcement Profile

Status: implemented and merged in PR #6; closure reflected in SpecKit tasks
after post-merge status reconciliation.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/14-gate-contract-explain-override.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/15-signed-checkpoint-replay-resistance.md`
- `schema/gate-result.schema.json`
- `schema/signed-checkpoint.schema.json`
- `schema/trusted-checkpoint-policy.schema.json`

## Goal

Make a protected gate profile explicit, deterministic, and fail-closed for
missing or invalid protected evidence while keeping the actual merge, release,
readiness, degradation, override approval, or risk-acceptance decision outside
`sdp-trace`.

The user-facing outcome is that a platform owner can wire `sdp-trace gate` into
CI and know that the protected profile will never pass from local-only evidence,
missing checkpoint evidence, unverifiable signer authority, stale witness data,
or override presence alone.

## Problem

Block 14 made gate output explainable but advisory. Block 15 made signed
checkpoints replay-resistant, but intentionally kept local signatures below
protected trust.

The next product gap is specific:

- protected gate requirements are still represented as `protected_future`;
- signed checkpoint verification facts are not evaluated as one gate profile;
- local signed evidence can be inspected, but there is no fail-closed profile
  that rejects it for protected use;
- CI owners need deterministic exit behavior and next actions without
  `sdp-trace` becoming the policy owner.

The weak framing would be "implement protected enforcement in `sdp-trace`."
That would violate the repository boundary. The correct framing is "implement a
protected verification profile whose fail-closed result can be enforced by an
external pipeline or policy consumer."

## Non-Goals

- No native merge, release, readiness, degradation, override approval, or risk
  acceptance decision.
- No GitHub-specific protected branch or pull-request API integration.
- No dependency on `sdp-gate`, Beads, agentloop, Operator Mode, Claude, Codex,
  OpenCode, GitHub, or any harness runtime.
- No Sigstore/Rekor or external witness service implementation unless already
  represented through a portable evidence artifact.
- No automatic approval of overrides.
- No raw stdout, stderr, prompt, model response, source snippet, credential,
  OIDC request token, or secret-like value in output or persisted artifacts.

## Product Boundary

Block 16 may emit verifier-derived protected profile facts:

- protected profile selected explicitly;
- required run and evidence coverage;
- CI witness binding state;
- signed checkpoint verification state;
- signer authority state;
- trust cap and profile state;
- override record presence;
- deterministic exit code suitable for CI wiring.

Block 16 must not emit a native policy decision. External consumers may use the
profile facts and exit code to block, permit, escalate, or accept risk.

## Contract Delta

### Protected Profile Selection

Protected evaluation must be explicit. A default `sdp-trace gate` run remains
observation or advisory unless the user selects the protected profile.

Minimal CLI shape:

```text
sdp-trace gate --profile protected --out <gate-result.json> --contract <contract.json> --checkpoint <checkpoint.json> --checkpoint-policy <policy.json> --witness <ci-witness.json> <runs-root-or-run-dir>
```

Rules:

- `--profile protected` selects the protected gate profile.
- Protected profile input must include a signed checkpoint and a
  trusted-checkpoint policy.
- For `sdp-trace gate --profile protected`, omitting required input flags
  (`--checkpoint`, `--checkpoint-policy`, or `--witness`) is a usage error:
  exit `2`, no gate artifact. A named file that is unreadable or not JSON is
  also a usage error. A readable artifact that cannot satisfy the protected
  verifier emits `fail` or `cannot_verify` in the gate result.
- Missing checkpoint evidence inside readable protected-profile input produces
  `cannot_verify`, not `pass`.
- Missing trusted-checkpoint policy evidence inside readable protected-profile
  input produces `cannot_verify`, not `not_assessed`, because protected profile
  requires signer authority evidence.
- Local-development checkpoint authority fails protected profile use even when
  the checkpoint signature itself verifies, because local-signed trust is
  explicitly outside protected pass scope.
- CI isolated signer authority may pass only when the verifier can bind the CI
  witness to the selected source, ref, commit, run id, and artifact digests.
- External witness authority remains `not_integrated` unless a portable
  external witness verification artifact exists and a later approved block
  implements that verifier. Block 16 recognizes the state but cannot use it for
  protected pass.
- Override records remain visible and produce an override-present condition;
  they never convert protected profile failures or missing evidence to pass.

### Protected Gate Result Shape

Block 16 extends gate output with protected profile fields:

- `selected_profile`: `observation`, `advisory_ci`, or `protected`
- `protected_gate`: `pass`, `fail`, `cannot_verify`, or `not_assessed`
- `checkpoint_verification`: signed-checkpoint verification summary
- `protected_conditions`: deterministic condition rows

Protected condition rows include:

- `id`
- `state`
- `reason_code`
- `reason`
- `next_action`

Allowed protected condition states are `pass`, `fail`, `cannot_verify`,
`missing_telemetry`, `not_integrated`, and `not_assessed`. Top-level
`protected_gate` never emits `missing_telemetry` or `not_integrated`; those
lower-level states map to top-level `cannot_verify` unless a `fail` state is
also present.

Initial protected condition ids:

- `protected_profile_explicitly_selected`
- `all_required_runs_present`
- `all_required_evidence_observed`
- `ci_witness_bound`
- `witness_freshness_valid`
- `checkpoint_signature_valid`
- `checkpoint_run_binding_valid`
- `checkpoint_signer_authorized`
- `protected_trust_scope_satisfied`
- `override_does_not_upgrade_profile`

Ordering rule: `protected_conditions` is emitted in the fixed condition-id order
above, and every condition remains visible even when another condition has a
dominant failure. Separate `reasons` and `next_actions` arrays use severity
order `fail`, `cannot_verify`, `missing_telemetry`, `not_integrated`,
`not_assessed`, then `pass`, with condition-id order as the tie breaker inside
each severity. Reason codes must be stable identifiers; users should not need
to parse prose to distinguish `missing_checkpoint`, `missing_policy`,
`stale_witness`, `source_mismatch`, or `local_signed_not_protected`.

Condition semantics:

- `ci_witness_bound` means repository, ref, commit, run id, and artifact digests
  match the selected gate input where those fields are required by the profile.
  Any required witness field that is absent is `cannot_verify`; any required
  witness field that is present and contradicts the selected gate input is
  `fail`. An empty artifact digest list is `cannot_verify`; a non-empty digest
  list with no matching digest is `fail`.
- `witness_freshness_valid` means supplied witness evidence is not stale under
  the selected witness freshness fields. If freshness fields are absent,
  protected profile emits `cannot_verify`; if they contradict the selected run
  context or declared freshness window, protected profile emits `fail`.
  Explicitly unbounded, null, or unknown freshness windows are treated as
  `cannot_verify` because protected profile cannot confirm the evidence is not
  stale. Expired witness evidence, witness timestamps after their declared
  expiry, witness timestamps before the selected run context when that context
  requires post-run witnessing, and incompatible freshness-window declarations
  are contradictions and emit `fail`. `cannot_verify` is still fail-closed for
  CI use: it exits non-zero and cannot satisfy protected profile.
- `checkpoint_run_binding_valid` means Block 15 checkpoint run binding,
  chain binding, source binding, and nonce binding all pass against the selected
  run context. CI witness binding remains a separate condition.
- `checkpoint_signer_authorized` means a supplied trusted-checkpoint policy
  authorizes the checkpoint signer for the protected profile authority boundary.
  It evaluates to `pass` only when the policy explicitly authorizes the signer
  within protected trust scope. If the policy exists but signer authority is
  `not_assessed` or `not_integrated`, the condition emits `cannot_verify` with
  an authority-gap reason. If the policy is missing entirely, protected profile
  emits `cannot_verify` with a missing-policy reason before treating signer
  authorization as satisfied.
- `protected_trust_scope_satisfied` is derived from checkpoint trust scope plus
  witness binding. It is stricter than the output `trust_cap`; `trust_cap`
  describes the strongest observed evidence, while this condition says whether
  that evidence satisfies protected use. In Block 16 it passes only when
  checkpoint verification passes with CI signer authority, CI witness binding
  passes, witness freshness passes, and required run/evidence conditions pass.
- `override_does_not_upgrade_profile` means override records are visible and
  non-upgrading. A well-formed Block 14 override request passes this condition
  only as "non-upgrading"; malformed, expired, or contract-unlinked override
  records are visible and `cannot_verify`, but still cannot upgrade the profile
  or reduce the severity of any prior `fail` or `cannot_verify` condition.

Override state table:

| Override state | Condition state | Reason |
|---|---|---|
| No override records present | `pass` | No override request is available to upgrade the profile. |
| Well-formed, contract-linked override present | `pass` | Override request is visible and non-upgrading. |
| Malformed, expired, or contract-unlinked override present | `cannot_verify` | Override record cannot be validated; it remains visible and non-upgrading. |

### Protected Trust Scope

Block 16 implemented protected pass scope:

- `ci_signed`

Disallowed for protected pass:

- `local_observed`
- `local_signed`
- `untrusted_shape_only`
- `external_witnessed` until an external witness verifier is implemented in a
  later approved block

Rules:

- `pass + local_signed` checkpoint verification with local-development
  authority is protected profile `fail`, because the authority boundary is
  explicitly outside protected pass scope.
- A signer policy mismatch is `fail`.
- A missing signer policy is `cannot_verify`.
- A supplied signer policy that leaves signer authority `not_assessed` or
  `not_integrated` maps to protected profile `cannot_verify`.
- A CI signer without source/run/artifact witness binding is `cannot_verify`.
- A CI witness source, ref, commit, run id, or artifact digest contradiction is
  `fail`.
- External witness values remain `cannot_verify` or `not_integrated` until a
  concrete external witness artifact and verifier are in scope.
- Trust scope values are the Block 15 checkpoint trust scopes plus Block 14
  gate `trust_cap` values. Block 16 must keep the protected schema enum explicit
  when extending `schema/gate-result.schema.json`.

### Schema Compatibility

Block 16 must not silently break Block 14 consumers. Implementation must
introduce a new Block 16 gate-result schema version for protected profile output
and keep `gate explain` able to read Block 14 artifacts without requiring
`selected_profile`, `protected_gate`, `checkpoint_verification`, or
`protected_conditions` to be present.

The reader must detect gate-result version from `schema_version`. For pre-Block
16 artifacts, protected fields are treated as absent and protected profile is
not inferred. `gate explain` must not default old artifacts to protected or
advisory policy conclusions. Schema version naming follows the existing
Block 14 gate-result convention; the normative compatibility requirement is
version separation plus Block 14 read compatibility.

### Exit Behavior

For `--profile protected`:

- exit `0` only when `protected_gate` is `pass`;
- exit `1` when `protected_gate` is `fail`;
- exit `3` when `protected_gate` is `cannot_verify` or `not_assessed` and no
  failure is present;
- usage errors remain exit `2`.

This makes the command CI-friendly without claiming `sdp-trace` owns the
organization's merge or release policy.

If mixed states appear, `fail` dominates `cannot_verify`, `not_integrated`, and
`not_assessed`; `cannot_verify` dominates `not_integrated` and `not_assessed`
for the emitted `protected_gate` state. The original lower-level state remains
visible in condition rows and checkpoint verification summary.

## User-Facing Commands

### `sdp-trace gate --profile protected`

Behavior:

1. Check required protected-profile flags and input file readability.
2. Evaluate required runs and required evidence from the selected contract.
3. Verify the supplied signed checkpoint against the selected run.
4. Check signer authority against the supplied trusted-checkpoint policy.
5. Bind the CI witness to repository, ref, commit, run id, and artifact digests
   when CI authority is required.
6. Emit deterministic protected profile facts, reasons, and next actions.
7. Write the gate result only to the explicit output path.
8. Never print or persist raw sensitive fields.

Next actions may say to supply a checkpoint, supply a signer policy, rerun under
the required CI observer, fix witness binding, regenerate checkpoint evidence,
or lower the claim. They must not recommend merge, release, approval, readiness,
or risk acceptance.

### `sdp-trace gate preview --profile protected`

Behavior:

1. Show which protected-profile inputs would be required.
2. Show each protected input path and inspectability status:
   `absent`, `present_readable`, `present_malformed`, or `present_unreadable`.
3. Show next actions for absent, malformed, or unreadable inputs.
4. Show selected profile and trust cap.
5. Do not verify signatures, signer authority, or witness binding.
6. Do not emit a protected verdict.
7. Do not write artifacts.

Exit behavior:

- exit `0` when preview input can be inspected and rendered;
- exit `2` for usage errors;
- exit `3` when required input files are named but unreadable or malformed.

`gate explain` remains an explanation command: it exits `0` when the selected
gate-result artifact can be parsed and explained, `2` for usage errors, and `3`
for unreadable or unsupported gate-result artifacts. It does not reuse the gate
artifact's protected pass/fail exit code.

Protected explain output must include:

- selected profile and top-level protected gate state when present;
- checkpoint verification summary before condition rows;
- protected condition rows in deterministic order;
- Block 14 gate conditions when the artifact also includes them;
- next actions using stable reason codes and no raw sensitive values.

For Block 14 artifacts, explain must render the existing Block 14 output and
state that protected profile fields are absent rather than inferring a protected
state.

## Acceptance Criteria

1. A protected profile run that omits required protected input flags exits `2`
   without writing a gate artifact.
2. A protected profile run with readable input that lacks checkpoint evidence
   emits `protected_gate: cannot_verify` and exits `3`.
3. A default gate run without `--profile protected` does not evaluate protected
   conditions.
4. A protected profile run with a valid local-development checkpoint does not
   pass protected gate.
5. A protected profile run with signer policy mismatch emits `protected_gate:
   fail` and exits `1`.
6. A protected profile run with CI signer authority but missing or incomplete CI
   witness binding emits `protected_gate: cannot_verify`.
7. A protected profile run with CI witness source, ref, commit, run id, or
   artifact digest contradiction emits `protected_gate: fail`.
8. A stale CI witness emits deterministic `fail`; absent or unbounded freshness
   emits `cannot_verify`; neither can satisfy protected profile.
9. A malformed override request remains visible and cannot reduce a failure to
   `cannot_verify` or change `protected_gate` to pass.
10. A valid override request is visible, but it does not change
   `protected_gate` to pass.
11. `gate preview --profile protected` renders missing protected inputs without
   verifying signatures or writing artifacts.
12. `gate explain` and `gate preview` include protected-profile reasons and next
   actions without raw sensitive fields.
13. Existing Block 14 gate-result artifacts remain explainable after Block 16.
14. All behavior is implemented in Go with test-first coverage.

## Implementation Plan

### Slice A: Contract And Schema Delta

Files:

- `schema/gate-result.schema.json`
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`

Tasks:

- add protected profile fields and condition rows to gate-result schema;
- introduce a new Block 16 gate-result schema version with explicit read
  compatibility in `gate explain` for Block 14 artifacts;
- add FRs, success criteria, and tasks for protected profile behavior;
- keep native policy decisions out of schema names and docs.

### Slice B: Protected Profile Evaluation

Files:

- `internal/demo/demo.go`
- `internal/demo/demo_test.go`
- `internal/checkpoint/checkpoint.go`
- `internal/checkpoint/checkpoint_test.go`

Tasks:

- evaluate protected profile only when explicitly selected;
- require checkpoint and signer policy for protected profile;
- map checkpoint verification states to protected profile states;
- prove local-development signed checkpoint evidence cannot pass protected
  profile;
- preserve all protected condition rows even when top-level protected gate has a
  dominant failure.

### Slice C: CI Witness And Authority Binding

Files:

- `internal/demo/demo.go`
- `internal/demo/demo_test.go`
- `internal/witness/witness.go`
- `internal/witness/witness_test.go`

Tasks:

- bind CI authority to existing witness source/run/artifact checks;
- distinguish incomplete CI authority context from contradictory witness data;
- add negative tests for missing witness, incomplete witness, and source or
  artifact mismatch;
- add stale witness tests for absent freshness data and contradictory freshness
  data.

### Slice D: CLI, Explain, And Preview

Files:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/main_test.go`
- `internal/demo/demo.go`
- `internal/demo/demo_test.go`

Tasks:

- add `--profile protected`;
- add checkpoint and policy flags to protected gate evaluation;
- add protected profile exit-code tests;
- extend explain and preview for protected input requirements and next actions;
- add secret-like and path-derived sensitive output assertions.

### Slice E: Fixtures, Verification, And Review

Files:

- `examples/block16-protected-gate/`
- `archive/research/block-16-implementation-evidence.md`
- `archive/research/block-16-implementation-review-disposition.md`

Tasks:

- add committed fixtures for missing checkpoint, local-development checkpoint,
  local-development checkpoint with invalid run binding, signer mismatch,
  missing CI witness, absent freshness, stale CI witness, CI source mismatch,
  CI artifact mismatch, malformed override with a trust-scope failure, valid
  CI-authority protected profile, and override-present protected profile;
- run Go-first verification and schema checks;
- run strict review and pi review before PR closure.

## Review Ledger

| Area | State | Disposition |
|---|---|---|
| Spec delta | drafted | Requires user approval and strict review before implementation. |
| Product boundary | drafted | Protected profile is verifier facts plus exit behavior, not native policy ownership. |
| Local signed checkpoint | drafted | Must not pass protected profile. |
| CI authority binding | drafted | Requires witness binding before protected pass. |
| External witness trust | not_integrated | Reserved until a portable external witness verifier exists. |
| Override approval | not_integrated | Override records remain visible but non-upgrading. |

## No-Overclaim Notes

- Protected profile pass is still not a native merge or release decision.
- CI may enforce the command exit code, but CI owns that enforcement policy.
- A local signed checkpoint is useful replay evidence, not protected authority.
- Missing protected evidence must not be softened to advisory pass.
- Override presence is evidence of an override request, not approval.
