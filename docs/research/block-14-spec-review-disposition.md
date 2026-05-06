# Block 14 Spec Review Disposition

Date: 2026-05-06

Scope:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/14-gate-contract-explain-override.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`

Review models:

- MiniMax-M2.7, no tools, no context files.
- ZAI GLM-5.1, no tools, no context files.

## MiniMax-M2.7 Findings

| Severity | Finding | Disposition |
|---|---|---|
| Critical | Required evidence schema and run-local references were underdefined. | Accepted. Block 14 now states required-run `required_evidence` references top-level contract-declared evidence ids and adds required-run state rules. |
| Critical | Exit code wording for `cannot_verify` was ambiguous. | Accepted. Exit code rules now define fail before cannot-verify before pass. |
| Critical | Explain next actions risked policy overclaim. | Accepted. Next actions are constrained to remediation hints and must not recommend merge, release, approval, readiness, or risk acceptance. |
| Critical | Override ingestion command was missing. | Accepted. Added minimal `sdp-trace override request` command shape and required field rejection rules. |
| Major | Trust cap and required-run states needed explicit mapping. | Accepted. Added required-run state table and witness binding state table. |
| Major | Deterministic ordering was not specified. | Accepted for required runs and observed runs; implementation must carry this through for reasons and next actions. |
| Major | Witness inspection and audit-grade meaning were unclear. | Accepted. Preview now reports locally detectable binding mismatches without verdicts; audit-grade `cannot_verify` is explicitly not a release decision. |

## ZAI GLM-5.1 Findings

| Severity | Finding | Disposition |
|---|---|---|
| Major | `policy_override_requested` was not connected to flight-recorder chain semantics. | Accepted. Native recorder events now inherit Block 09 chain fields; external imports remain capped. |
| Major | Override effect on required-run state was undefined. | Accepted. Override presence is informational and does not change required-run or required-evidence state. |
| Major | Gate-result schema task was missing. | Accepted. Added T107 for Draft 2020-12 gate-result schema. |
| Major | `gate_conditions` and `runs` shapes were undefined. | Accepted. Added field shapes and initial condition ids. |
| Major | Exit code semantics across local, CI, and audit gate dimensions were ambiguous. | Accepted. Added selected-dimension precedence rules. |
| Major | Committed fixture task was missing for SC-034 through SC-036. | Accepted. Added T113 fixture task and Slice E. |
| Major | Override event lacked `producer` and `origin`. | Accepted. Added both required fields and allowed origin values. |
| Minor | Partial witness binding behavior needed a table. | Accepted. Added witness binding table. |
| Minor | Preview mismatch surfacing was ambiguous. | Accepted. Preview can report locally detectable mismatches but cannot emit a verdict. |

## Residual State

Post-fix GLM-5.1 review reported no remaining critical or major findings.

Accepted minor clarifications after post-fix review:

- deterministic ordering for reasons and next actions;
- invalid override references produce `cannot_verify` on the override record;
- Block 14 evaluates every applicable emitted gate dimension and has no
  user-selectable dimension filter.

Required implementation proof remains `not_assessed` until Go tests, fixtures,
schema validation, and fresh review run against code.
