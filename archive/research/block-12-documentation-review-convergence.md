# Block 12 Documentation Review Convergence

Date: 2026-05-05

Reviewed documents:

- `docs/cto-adoption-guide.en.md`
- `docs/cto-adoption-guide.ru.md`
- `docs/team-lead-playbook.en.md`
- `docs/team-lead-playbook.ru.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/12-ci-witness-adoption.md`

## Review Personas

| Persona | Model | Initial Result | Final Result |
|---|---|---|---|
| CTO buyer | MiniMax-M2.7 | Critical gaps: non-cooperating agents, policy-ready overclaim, undefined control mechanism. | No critical/major findings. |
| Platform / Harness Owner | GLM-5.1 | Critical gaps: capture boundary, fail-closed enforcement, tamper evidence. | No critical/major findings. |
| CISO / Adversarial Trust | Kimi K2P6 | Critical gap: plain JSON witness could be mistaken for external trust. | No critical/major findings. |
| Staff Engineer / DX Skeptic | Qwen3.6 Plus | Critical gaps: emergency path, redaction/privacy, OIDC token lifecycle. | No critical/major findings; one minor future enhancement for native override ergonomics. |
| Compliance / Forensics Lead | DeepSeek V4 Pro | Critical/major gaps: signed timeline, query surface, retention, redaction audit trail. | No critical/major findings. |

## Accepted Documentation Corrections

- Defined the Block 12 control layer as evidence capture, missing telemetry, and
  CI-witnessed reporting, not automatic merge blocking or external audit proof.
- Explained what is visible when agents do not cooperate: missing expected
  evidence at CI/report/gate boundaries, not full local activity.
- Added trust-state interpretation for `local_observed`, `ci_witnessed`, and
  future `external_witnessed`.
- Stated that Block 12 CI witness is a CI-generated JSON artifact, not a DSSE
  envelope, transparency log entry, or signed forensic timeline.
- Warned that developer-committed witness files are not authority.
- Documented capture boundary: process lifecycle and command-level metadata, not
  arbitrary harness internals.
- Documented product gaps: fail-closed managed harness enforcement, deletion and
  replay detection, signed timeline, query/dashboard analytics, redaction audit
  trail, retention implementation, non-GitHub CI profiles, and native
  `policy_override_requested` events.
- Added privacy guidance: default digest-only stdout/stderr; no raw prompt,
  source, or model-response capture in Block 12.
- Added emergency-change guidance: use external policy/change-management records
  and keep `cannot_verify`/`missing_telemetry` visible.
- Added offline/failure-mode guidance and a gate debugging checklist.
- Added retention guidance for `.sdp-trace-report/` and `.sdp-trace-runs/`.

## Residual Minor

Staff Engineer / DX review recommends a future native override affordance such
as `gate --override-reason` that writes a `policy_override_requested` trace
event linked to the external policy system.

Disposition: accepted as future product work. Block 12 documentation now states
that native `policy_override_requested` is not implemented and that overrides
must remain external records referencing report artifacts.

## Convergence State

No persona has remaining critical or major documentation findings after the
second review pass.
