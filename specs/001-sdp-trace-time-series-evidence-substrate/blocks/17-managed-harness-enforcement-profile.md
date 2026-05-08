# Block 17: Managed Harness Enforcement Profile

Status: spec delta and implementation plan; awaiting explicit approval before
implementation.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13b-capture-boundary-state-dx-baseline.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/16-protected-gate-enforcement-profile.md`
- `archive/research/harness-telemetry-trust-brief.md`
- `archive/research/agentic-sdlc-evidence-substrate-v3-brief.md`

## Goal

Make managed harness enforcement explicit, opt-in, and fail-closed for teams
that agree to run through a registered wrapper or adapter boundary, while
preserving observation-mode value for unmanaged harnesses.

The user-facing outcome is that a platform owner can prove whether a selected
run used an approved managed boundary, whether required wrapper and adapter
telemetry was present, and whether bypass or suppression prevents a managed
profile claim. `sdp-trace` emits verifier facts and deterministic exit behavior;
external CI, `sdp-gate`, or customer policy owns the block/allow decision.

## Problem

Block 16 protects a gate by requiring signed checkpoint, signer authority, CI
witness binding, and required evidence. That still does not prove that local
agent work happened inside a platform-managed harness boundary.

The next product gap is specific:

- unmanaged harnesses can still produce useful observation-mode traces, but
  they cannot satisfy a managed-harness claim;
- wrapper or adapter enrollment is not represented as a first-class contract;
- adapter identity, capability, suppression, and disconnect states are not
  evaluated as one fail-closed profile;
- CI can observe that protected evidence exists, but not that the run used an
  approved harness wrapper or adapter path;
- platform owners need deterministic evidence of bypass, missing adapter
  telemetry, unauthorized adapter identity, or late enrollment.

The weak framing would be "require every team to use the managed harness." That
would break the adoption model. The correct framing is "add a managed profile
that only claims managed enforcement when the selected run was deliberately
enrolled and the required wrapper/adapter evidence verifies."

## Command Naming Risk

The existing `sdp-trace gate` command name is a UX and product-boundary risk
because `sdp-gate` is the complementary policy project. Adding more profiles
under `sdp-trace gate` makes it easier for users to believe `sdp-trace` owns
gate decisions.

Block 17 should not deepen that ambiguity. The preferred new command surface is
profile assessment:

```text
sdp-trace assess --profile managed-harness ...
sdp-trace assess preview --profile managed-harness ...
sdp-trace assess explain <assessment-result.json>
```

Existing Block 14/16 `sdp-trace gate` behavior remains supported for existing
profiles in Block 17. Block 17 must not add `sdp-trace gate --profile
managed-harness`; callers must use `sdp-trace assess --profile
managed-harness`. A later approved migration may add `assess` aliases for
Block 14/16 profiles, but Block 17 does not remove or rename existing commands.
`sdp-gate` remains the policy consumer that turns assessment facts into a gate
decision.

Block 17 must also avoid new schema field names that imply `sdp-trace` owns the
gate decision. New managed fields use assessment wording. Existing Block 14/16
fields with gate wording remain compatibility debt and must not be copied into
new managed fields.

## Non-Goals

- No replacement of OpenCode, pi, Kilo, Superpowers-style harnesses, `gsd`,
  `gsd2`, Oh My OpenAgent, custom agents, or customer wrappers.
- No requirement that observation-mode users enroll in managed mode.
- No harness-specific SDK, plugin runtime, or dependency on a named harness.
- No native merge, release, readiness, degradation, override approval, or
  risk-acceptance decision.
- No automatic approval of managed-mode bypass or emergency override.
- No raw prompt, model response, source snippet, stdout/stderr body, command
  arguments, credential, OIDC request token, adapter secret, gateway token, or
  checkpoint key material in output or persisted artifacts.

## Product Boundary

Block 17 may emit verifier-derived managed profile facts:

- managed profile selected explicitly;
- expected managed boundary registered before run start;
- wrapper enrollment and run launch state;
- adapter registration, identity, capability, activation, disconnect, and
  suppression state;
- observed harness, model, tool-call, file-mutation, and test-provenance event
  coverage where the contract requires those events;
- managed trust scope and completeness state;
- bypass, late enrollment, unauthorized observer, unsupported capability, and
  missing telemetry conditions;
- deterministic exit code suitable for external CI or policy wiring.

Block 17 must not emit a native policy decision. External consumers may block,
permit, escalate, or accept risk based on these facts.

## Threat Model

Block 17 assumes an adversarial or careless local operator may control the
workspace, run directory, adapter registry path, managed policy path, and local
filesystem timestamps. It also assumes an adapter may be buggy, compromised, or
self-claiming.

Managed profile does not defend against a compromised CI identity, compromised
customer policy authority, compromised external witness, malicious kernel, or
hardware-level host compromise. Those remain outside Block 17 and must not be
collapsed into a managed pass.

Managed profile is designed to detect or fail closed for:

- unmanaged runs presented as managed after the fact;
- post-hoc adapter registry or managed policy creation;
- self-claimed or unauthorized adapter identity;
- adapter capability claims not matched by required event evidence;
- adapter disconnect, suppression, or missing telemetry during the run;
- replay of a prior managed run against a different source, run id, policy, or
  chain context;
- agent-reported test evidence presented as executed test proof.

## Managed Enrollment

Managed enrollment is a verifier-visible per-run event, not prose and not a
team-wide assumption.

A run is deliberately enrolled only when the run artifact contains a
`managed_boundary_enrolled` event before child harness execution. The event
must include:

- managed policy digest;
- adapter registry digest;
- selected wrapper id or adapter id;
- enrollment source: `managed_policy`, `ci_config`, `human_signed`, or
  `customer_policy_equivalent`;
- selected managed profile id;
- parent run id;
- run nonce;
- enrollment event digest;
- sequence number in the same recorder chain as the child launch event.

Ordering rule:

- `managed_boundary_enrolled.sequence` must be less than the first child
  `command_started` or equivalent harness launch event sequence.
- If adapter events are written independently from the recorder chain, those
  events cannot prove pre-run enrollment unless they are bound into the same
  run chain by digest before the child launch or by a managed witness that signs
  both the adapter event digest and the launch event digest.
- Post-hoc registry updates, agent-authored enrollment summaries, and events
  without monotonic sequence evidence are `cannot_verify` or `fail`, never
  managed pass.

No repository-global setting or adapter registry entry alone proves enrollment.
Those artifacts define what may be accepted; the run must still show that the
managed boundary was active before work started.

Policy and registry provenance rule:

- managed policy and adapter registry digests must be committed in VCS,
  pinned in CI configuration, human-signed, or covered by an accepted customer
  policy equivalent before the run starts;
- policy or registry artifacts created only inside the run under assessment are
  local reconstruction evidence and cannot authorize managed pass;
- local file modification time is not evidence of pre-run policy or registry
  existence.

## Adapter Authority

Adapter identity is authorized only when the managed policy and adapter
registry agree on the adapter identity and authority boundary.

The managed policy must define:

- `policy_provenance`: `vcs`, `ci_config`, `human_signed`, or
  `customer_policy_equivalent`, plus the binding reference and digest;
- `authorized_adapters`: adapter id, harness id, version constraint, allowed
  signer or authority reference, allowed deployment source, and allowed
  capability ids;
- `authorized_wrappers`: wrapper id, version constraint, deployment source,
  allowed signer or authority reference, and allowed command descriptor labels;
- `required_event_groups`: mapping from condition group to required event
  types and acceptable provenance scopes;
- `suppression_rules`: explicit event groups that may be suppressed, authority
  required for suppression, and whether suppression may satisfy the managed
  profile or only remain visible as `cannot_verify`;
- `witness_binding`: required source, run, policy, registry, wrapper or adapter,
  freshness, chain, and artifact digest fields.

The adapter registry must define:

- adapter id, harness id, version, deployment source, identity state, and
  authority reference;
- capabilities as structured entries, not just free-form strings:
  capability id, version, emitted event types, provenance scope, and required
  payload digest fields;
- allowed event types and schema versions emitted by the adapter.

Identity states:

- `verified`: policy authority reference matches the registry authority
  reference and the verifier can validate the binding.
- `self_claimed`: adapter declared identity but verifier cannot validate the
  authority reference; cannot satisfy managed profile.
- `unauthorized`: adapter identity, signer, version, deployment source, or
  authority reference contradicts policy; managed profile `fail`.
- `not_assessed`: no adapter identity assessment was selected or implemented;
  managed profile `cannot_verify`.

The verifier must not infer adapter authority from adapter id text, harness
name, file path, or agent prose.

Adapter identity is not audit proof. `verified` means the adapter identity
matches the selected managed policy under the configured authority boundary; it
does not prove the adapter implementation is honest, complete, or externally
audited.

## Managed Witness Binding

Managed witness evidence must bind:

- run id and run nonce;
- source repository, ref when available, and source commit or tree digest;
- managed policy digest and policy provenance reference;
- adapter registry digest and registry provenance reference;
- wrapper id or adapter id plus authority reference digest;
- `managed_boundary_enrolled` event digest;
- child harness launch event digest;
- run chain head and event count;
- output artifact digests;
- witness generated time, expiry or freshness window, and witness identity.

Missing freshness fields produce `cannot_verify`. Expired witness evidence,
source mismatch, run id mismatch, run nonce mismatch, policy or registry digest
mismatch, enrollment event mismatch, launch event mismatch, chain-head mismatch,
or artifact digest mismatch produces `fail`.

Managed witness binding still does not produce external audit proof. It does
not establish third-party attestation, append-only transparency, hardware root
of trust, or external witness state. Those remain future external-audit profile
requirements.

## Contract Delta

### Managed Profile Selection

Managed harness evaluation must be explicit. A default `sdp-trace gate` run
remains observation/advisory/protected according to the selected profile and
must not infer managed harness enforcement.

Preferred CLI shape:

```text
sdp-trace assess --profile managed-harness --out <assessment-result.json> --contract <contract.json> --run <run-dir> --adapter-registry <adapter-registry.json> --managed-policy <managed-policy.json> --managed-witness <managed-witness.json>
```

Rules:

- `--profile managed-harness` selects the managed harness enforcement profile.
- Managed profile input must include a managed policy and adapter registry.
- Omitting `--run`, `--adapter-registry`, `--managed-policy`, or
  `--managed-witness` is a usage error: exit `2`, no assessment artifact.
- A named input file that is unreadable or not JSON is a usage error. A
  readable artifact that cannot satisfy the managed verifier emits `fail` or
  `cannot_verify` in the gate result.
- Managed enforcement requires a run artifact whose managed boundary is
  established by `managed_boundary_enrolled` in the same verified run chain
  before the child harness command starts.
- Late attachment, post-hoc adapter registration, post-hoc managed policy
  creation, or agent-authored enrollment evidence cannot satisfy the managed
  profile.
- Unmanaged harnesses remain valid observation-mode inputs and must produce
  `not_integrated`, `unsupported`, `missing_telemetry`, or `not_assessed`
  rather than an inferred managed pass.
- Override records remain visible and non-upgrading. They cannot convert
  bypass, missing telemetry, unauthorized adapter identity, or suppression into
  managed pass.

### Managed Assessment Result Shape

Block 17 extends gate output with managed profile fields:

- `selected_profile`: adds `managed_harness`
- `managed_harness_assessment`: `pass`, `fail`, `cannot_verify`, or
  `not_assessed`
- `managed_boundary`: wrapper and adapter verification summary
- `managed_conditions`: deterministic condition rows

Managed condition rows include:

- `id`
- `state`
- `reason_code`
- `reason`
- `next_action`

Allowed managed condition states are `pass`, `fail`, `cannot_verify`,
`missing_telemetry`, `not_integrated`, `unsupported`, `suppressed`, and
`not_assessed`. Top-level `managed_harness_assessment` never emits
`missing_telemetry`, `not_integrated`, `unsupported`, or `suppressed`; these map
to top-level `cannot_verify` unless a `fail` state is also present.

Top-level precedence:

1. Any `fail` condition produces `managed_harness_assessment: fail`.
2. Else any `cannot_verify`, `missing_telemetry`, `suppressed`,
   `not_integrated`, or `unsupported` condition produces
   `managed_harness_assessment: cannot_verify`.
3. Else any `not_assessed` condition produces
   `managed_harness_assessment: not_assessed`.
4. Else all conditions are `pass` and the top-level assessment is `pass`.

Initial managed condition ids:

- `managed_profile_explicitly_selected`
- `managed_policy_loaded`
- `adapter_registry_loaded`
- `managed_boundary_enrolled_before_run`
- `adapter_identity_authorized`
- `adapter_capabilities_satisfy_contract`
- `adapter_activation_observed`
- `adapter_connection_continuous`
- `required_harness_events_observed`
- `required_tool_events_observed`
- `required_file_mutations_observed`
- `test_provenance_not_agent_reported`
- `suppression_policy_valid`
- `bypass_not_observed`
- `managed_witness_bound`
- `override_does_not_upgrade_managed_profile`

Ordering rule: `managed_conditions` is emitted in the fixed condition-id order
above. Every condition remains visible even when another condition has a
dominant failure. Separate `reasons` and `next_actions` arrays use severity
order `fail`, `cannot_verify`, `missing_telemetry`, `suppressed`,
`not_integrated`, `unsupported`, `not_assessed`, then `pass`, with condition-id
order as the tie breaker.

Stable reason codes must include at least:

- `missing_managed_policy`
- `missing_adapter_registry`
- `run_not_managed`
- `late_enrollment`
- `managed_boundary_not_in_chain`
- `adapter_identity_unauthorized`
- `adapter_capability_missing`
- `adapter_activation_missing`
- `adapter_disconnect_observed`
- `harness_event_missing`
- `tool_event_missing`
- `file_mutation_event_missing`
- `test_provenance_missing`
- `agent_reported_test_not_executed`
- `suppression_unverified`
- `tool_event_suppressed_by_policy`
- `tool_event_suppressed`
- `bypass_observed`
- `managed_witness_missing`
- `managed_witness_mismatch`
- `override_absent_non_upgrading`
- `override_present_non_upgrading`
- `override_upgrade_rejected`

Condition semantics:

- `managed_profile_explicitly_selected` is metadata confirmation. It emits
  `pass` only in managed profile output and is absent from non-managed output.
- `managed_policy_loaded` and `adapter_registry_loaded` refer to readable,
  schema-valid artifacts supplied to assessment. Omitted flags are usage errors
  and produce no artifact; readable but semantically incomplete policy or
  registry files produce condition rows.
- `managed_boundary_enrolled_before_run` means a
  `managed_boundary_enrolled` event appears in the selected run chain before
  child execution and references the supplied policy and registry digests.
- `adapter_connection_continuous` means no adapter disconnect, suppression, or
  unexplained adapter gap overlaps a required managed observation window.
- `adapter_capabilities_satisfy_contract` is derived from the selected
  contract and managed policy. The contract names required event types; the
  managed policy maps those event types into required event groups; the adapter
  registry maps capabilities to emitted event types and provenance scopes.
- `test_provenance_not_agent_reported` passes only when required test evidence
  provenance is `ci_witnessed` or `local_observed` from a registered tool or
  process wrapper that records command descriptor, exit state, and artifact
  digest. `harness_observed` adapter events may correlate the test to harness
  intent, but they cannot by themselves prove executed-test evidence unless
  they are bound to wrapper or CI execution evidence.
- `suppression_policy_valid` passes only when each suppressed required event is
  covered by an explicit suppression rule from a pre-run policy provenance
  source and the selected managed profile says that suppression may satisfy
  that event group. Authorized-but-non-satisfying suppression remains visible
  as `suppressed` and maps to top-level `cannot_verify`. Suppression authorized
  only by a run-local or agent-authored policy is `fail`.
- `override_does_not_upgrade_managed_profile` emits `pass` with
  `override_absent_non_upgrading` when no override is present, `pass` with
  `override_present_non_upgrading` when a valid override is present, and
  `fail` with `override_upgrade_rejected` if an artifact attempts to use an
  override to upgrade managed state.

### Managed Trust Scope

Block 17 introduces the managed profile trust scope:

- `managed_harness_observed`

This scope may be satisfied only when:

- the managed profile was explicitly selected;
- the managed policy and adapter registry are readable and valid;
- managed policy and adapter registry provenance is anchored before the run;
- the run started under a registered wrapper or authorized adapter before child
  execution;
- adapter identity is authorized by policy;
- required adapter capabilities cover the contract's required event types via
  the managed policy event-group mapping;
- required harness/tool/file/test-provenance events are present or explicitly
  allowed as suppressed by policy;
- test evidence required for managed profile is not only `agent_reported`;
- the CI witness or equivalent managed witness binds the run id, source commit,
  run nonce, wrapper/adapter identity digest, managed policy digest, adapter
  registry digest, enrollment event digest, launch event digest, chain head,
  freshness, and output artifact digests.

Disallowed for managed profile pass:

- `agent_reported`
- `local_observed` without registered managed boundary
- `harness_observed` from an unauthorized or self-claimed adapter
- late attachment after child harness start
- post-hoc adapter registry updates
- post-hoc managed policy updates
- missing required adapter capabilities
- missing required events unless suppression is both policy-authorized and
  verifier-visible
- CI witness that cannot bind policy, adapter identity, run id, source, and
  artifacts

## User-Facing Commands

### `sdp-trace assess --profile managed-harness`

Behavior:

1. Check required managed-profile flags and input readability.
2. Load the managed policy and adapter registry.
3. Confirm the selected run began under a registered wrapper or authorized
   adapter boundary.
4. Verify adapter identity and declared capabilities against policy.
5. Evaluate required harness, tool, file-mutation, and test-provenance events.
6. Distinguish missing telemetry, unsupported capability, suppression, bypass,
   and late enrollment.
7. Bind the managed witness to source, run id, run nonce,
   adapter/wrapper identity, managed policy digest, adapter registry digest,
   enrollment and launch event digests, chain head, freshness, and artifacts.
8. Emit deterministic managed profile facts, reasons, and next actions.
9. Write the assessment result only to the explicit output path.
10. Never print or persist raw sensitive fields.

Exit behavior:

- exit `0` only when `managed_harness_assessment` is `pass`;
- exit `1` when `managed_harness_assessment` is `fail`;
- exit `3` when `managed_harness_assessment` is `cannot_verify` or `not_assessed`
  and no failure is present;
- usage errors remain exit `2`.

Exit code table:

| Exit | Meaning | Canonical CI action |
|---|---|---|
| `0` | Managed harness assessment passed for the selected profile. | External CI or `sdp-gate` may consume the artifact as a passing managed-profile fact. |
| `1` | Verifier found contradictory evidence. | Block or escalate under the external policy; fix the run, policy, adapter, or witness. |
| `2` | Usage error: required flag omitted, named input unreadable, or named input not JSON. | Fix command invocation; no assessment artifact is authoritative. |
| `3` | The verifier could not establish managed profile facts, or the selected assessment was not available. | Treat as fail-closed for managed enforcement; either supply missing evidence or lower the claim to observation/advisory. |

### `sdp-trace assess preview --profile managed-harness`

Behavior:

1. Show managed-profile inputs that would be required.
2. Show each managed input path and inspectability status:
   `absent`, `present_readable`, `present_malformed`, or
   `present_unreadable`.
3. Show expected managed boundaries, adapter ids, capability ids, event ids,
   and witness binding fields.
4. Do not evaluate adapter authority, event completeness, or witness binding.
5. Do not emit a managed verdict.
6. Do not write artifacts.

Exit behavior:

- exit `0` when preview input can be inspected and rendered;
- exit `2` for usage errors;
- exit `3` when required input files are named but unreadable or malformed.

`sdp-trace assess explain` explains Block 17 assessment artifacts. For Block 14
or Block 16 gate-result artifacts, `assess explain` exits `3` and reports
`unsupported_pre_assess_artifact` with the detected schema version. Existing
`sdp-trace gate explain` remains the compatibility path for Block 14 and Block
16 artifacts and must state that managed profile fields are absent rather than
inferring managed state.

## Acceptance Criteria

1. A managed profile run that omits required managed input flags exits `2`
   without writing an assessment artifact.
2. A default gate run without `--profile managed-harness` does not evaluate
   managed conditions.
3. A managed profile run with no registered wrapper or adapter emits
   `managed_harness_assessment: fail` with `run_not_managed` or
   `bypass_observed`.
4. A run whose adapter registers after child harness start emits `fail` with
   `late_enrollment`.
5. A self-claimed or unauthorized adapter identity emits `fail`.
6. A registered adapter missing a contract-required capability emits
   `cannot_verify` or `unsupported` at the condition level and exits `3`
   unless another condition fails.
7. Missing required harness, tool, file-mutation, or test-provenance events
   emit deterministic `missing_telemetry` condition rows.
8. Test evidence that is only `agent_reported` cannot satisfy managed profile
   test-provenance requirements.
9. Adapter disconnect or unexplained adapter gap overlapping a required managed
   observation window emits deterministic non-pass state.
10. Suppressed evidence satisfies a managed condition only when suppression is
   policy-authorized by pre-run policy provenance, verifier-visible, and
   allowed by the selected profile.
11. Managed witness mismatch on run id, source commit, run nonce, policy digest,
    registry digest, enrollment event digest, launch event digest, chain head,
    adapter/wrapper identity digest, or artifact digest emits `fail`.
12. Managed witness absence or incomplete binding emits `cannot_verify`.
13. A valid override request is visible and non-upgrading; it does not change
    `managed_harness_assessment` to pass.
14. `assess preview --profile managed-harness` renders missing managed inputs
    without evaluating adapter authority or writing artifacts.
15. `assess explain` and `assess preview` include managed-profile reasons and next
    actions without raw sensitive fields.
16. Existing Block 14 and Block 16 gate-result artifacts remain explainable by
    `gate explain`.
17. All behavior is implemented in Go with test-first coverage.

## Implementation Plan

### Slice A: Contract And Schema Delta

Files:

- `schema/gate-result.schema.json` or successor assessment-result schema
- new `schema/managed-harness-policy.schema.json`
- new `schema/adapter-registry.schema.json`
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`

Tasks:

- add managed profile fields and condition rows to assessment/gate-result
  schema;
- add portable managed policy and adapter registry schemas;
- introduce a new Block 17 assessment-result schema version while preserving
  Block 14/16 read compatibility in `gate explain`;
- add FRs, success criteria, and tasks for managed profile behavior;
- keep native policy decisions out of schema names and docs.

### Slice B: Managed Boundary Evaluation

Files:

- `internal/demo/demo.go`
- `internal/demo/demo_test.go`
- new `internal/managed/managed.go`
- new `internal/managed/managed_test.go`

Tasks:

- evaluate managed profile only when explicitly selected;
- require managed policy, adapter registry, run, and witness inputs;
- detect no managed boundary, late enrollment, out-of-chain enrollment, adapter
  suppression, adapter disconnect, and bypass;
- preserve all managed condition rows even when top-level gate has a dominant
  failure.

### Slice C: Adapter Registry, Capability, And Event Coverage

Files:

- `internal/managed/managed.go`
- `internal/managed/managed_test.go`

Tasks:

- verify adapter identity against the selected managed policy;
- verify required capabilities against the contract;
- verify runtime event coverage separately from declaration-only capability
  matching;
- map missing required harness/tool/file/test events to deterministic condition
  rows;
- prove `agent_reported` test evidence cannot satisfy executed-test evidence.

### Slice D: CLI, Explain, And Preview

Files:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/main_test.go`
- `internal/demo/demo.go`
- `internal/demo/demo_test.go`

Tasks:

- add `sdp-trace assess --profile managed-harness`;
- add managed policy and adapter registry flags;
- add managed profile exit-code tests;
- add assess explain and preview for managed input requirements and next actions;
- preserve `sdp-trace gate explain` read compatibility for Block 14/16 artifacts;
- add secret-like and path-derived sensitive output assertions.

### Slice E: Fixtures, Verification, And Review

Files:

- `examples/block17-managed-harness/`
- `archive/research/block-17-implementation-evidence.md`
- `archive/research/block-17-implementation-review-disposition.md`

Tasks:

- add committed fixtures for unmanaged run, late enrollment, post-hoc policy or
  registry, unauthorized adapter, missing capability, missing harness event,
  missing tool event, missing file mutation event, adapter disconnect,
  agent-reported test evidence, policy-authorized suppression, run-local
  suppression policy, witness missing, stale witness, witness mismatch,
  override present, and valid managed profile;
- run Go-first verification and schema checks;
- run strict review and pi review before PR closure.

## Review Ledger

| Area | State | Disposition |
|---|---|---|
| Spec delta | drafted | Requires user approval and strict review before implementation. |
| Product boundary | drafted | Managed profile is verifier facts plus exit behavior, not native policy ownership. |
| Adoption posture | drafted | Observation mode remains available without managed enrollment. |
| Adapter authority | drafted | Self-claimed or unauthorized adapters cannot satisfy managed profile. |
| Suppression handling | drafted | Suppression is visible and non-upgrading unless policy-authorized for the selected profile. |
| External audit trust | not_integrated | Managed harness evidence is not external audit proof. |

## No-Overclaim Notes

- Managed harness pass is still not a native merge, release, readiness, or risk
  decision.
- Managed mode is opt-in and must not be described as a prerequisite for
  observation-mode value.
- Adapter events are not trustworthy until identity and capability authority
  are verified.
- Adapter identity verification does not prove adapter implementation honesty.
- Missing managed telemetry must not be softened to advisory pass.
- Suppression is not safe redaction unless a verifier can see the suppression
  authority and retention effect.
- `managed_harness_observed` is not `external_witnessed`.
- Managed witness binding is not third-party attestation, transparency-log
  proof, hardware-rooted proof, or external audit proof.
