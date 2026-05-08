# Block 17 Spec Review Disposition

Date: 2026-05-06

Reviewed artifact:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/17-managed-harness-enforcement-profile.md`

## Review Planes

| Plane | Model | State |
|---|---|---|
| CTO buyer / product boundary | MiniMax-M2.7 | changes_required |
| Platform harness owner | ZAI/GLM-5.1 | changes_required |
| CISO / forensics | MiniMax-M2.7 | changes_required |

## Valid Findings And Disposition

| Severity | Finding | Disposition |
|---|---|---|
| critical | New Block 17 fields reused `gate` wording and weakened the `sdp-trace` / `sdp-gate` boundary. | Fixed by making `sdp-trace assess --profile managed-harness` the primary Block 17 command and replacing `managed_harness_gate` with `managed_harness_assessment`. Existing `gate` wording remains Block 14/16 compatibility debt only. |
| critical | Managed enrollment was undefined and could be interpreted as a team-wide or prose state. | Fixed by defining `managed_boundary_enrolled` as a per-run event with policy digest, registry digest, wrapper/adapter id, enrollment source, profile id, run id, run nonce, event digest, and sequence before child launch. |
| critical | Adapter identity authority was undefined and self-attestation could satisfy managed profile. | Fixed by adding managed policy authority fields, adapter registry authority fields, explicit identity states, and a rule that identity cannot be inferred from adapter id, harness name, file path, or agent prose. |
| critical | Post-hoc adapter registry or managed policy creation was not detectable. | Fixed by requiring pre-run policy and registry provenance through VCS, CI config, human signature, or customer policy equivalent; local mtime is explicitly rejected as evidence. |
| critical | Suppression authorization was self-referential and attacker-definable. | Fixed by requiring suppression authorization from pre-run policy provenance and making run-local or agent-authored suppression policy a failure. |
| critical | Managed witness binding could be replayed or mistaken for audit proof. | Fixed by requiring run nonce, policy and registry provenance, enrollment event digest, launch event digest, chain head, event count, freshness, and witness identity; added explicit no-overclaim exclusions for external audit proof. |
| major | `assess` versus `gate` compatibility path was ambiguous. | Fixed by specifying Block 17 does not add `gate --profile managed-harness`; `assess` is primary for Block 17, while existing `gate explain` remains compatibility path for Block 14/16. |
| major | Required managed event sources and capability mapping were undefined. | Fixed by defining contract required event types, managed policy event groups, adapter registry capability-to-event mapping, and runtime event coverage separate from declaration-only capability matching. |
| major | Test provenance could be fabricated by an adapter while avoiding `agent_reported`. | Fixed by requiring test proof from CI or registered tool/process wrapper execution evidence; adapter events can correlate intent but do not by themselves prove execution. |
| major | Adapter disconnect during run was named but not evaluated. | Fixed by adding `adapter_connection_continuous`, `adapter_disconnect_observed`, and acceptance criteria for disconnect or unexplained gaps. |
| major | Exit code `3` was not actionable for CI. | Fixed by adding a consolidated exit-code table and requiring external policy to treat exit `3` as fail-closed for managed enforcement unless the claim is lowered. |
| minor | Override non-upgrading states were underspecified. | Fixed by defining absent, present, and upgrade-rejected override reason codes. |

## Remaining Scope Boundaries

- Block 17 does not implement external audit proof.
- Block 17 does not make managed enrollment mandatory for observation-mode users.
- Block 17 does not replace existing harnesses.
- Block 17 does not remove existing Block 14/16 `sdp-trace gate` commands.
- External CI, `sdp-gate`, or customer policy remains responsible for block,
  allow, approval, release, readiness, and risk decisions.

## Re-Review Requirement

Narrow re-review over the revised Block 17 spec and this disposition:

| Plane | Model | Result |
|---|---|---|
| Platform harness owner | ZAI/GLM-5.1 | `NO_REMAINING_CRITICAL_OR_MAJOR` |
| Product/trust adversarial | MiniMax-M2.7 | `NO_REMAINING_CRITICAL_OR_MAJOR` |

Implementation may start. Remaining risk is implementation correctness, not
spec-direction ambiguity.
