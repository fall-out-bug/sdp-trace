# Block 13 Product Gap Review Convergence

Date: 2026-05-06

Reviewed artifact:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`

## Review Personas

| Persona | Model | Initial Result | Disposition |
|---|---|---|---|
| CTO buyer | MiniMax-M2.7 | Revise: major concern that enforcement and required-run manifest could violate sidecar-first adoption. | Accepted; roadmap now separates observation, advisory gate, protected gate, managed harness, and external audit modes. |
| Platform / Harness Owner | GLM-5.1 | Revise: capture boundary was assumed, and managed wrapper architecture appeared too late. | Accepted; roadmap now adds Block 13B with explicit interception architecture, state taxonomy, doctor, preview, determinism, offline state, and overhead budget. |
| CISO / Adversarial Trust | Kimi K2P6 | Revise: gate before trust anchor creates forgery window; signer isolation and monotonic sequence were underspecified. | Accepted; Block 14 is advisory until signed checkpoints exist; Block 15 now requires signer isolation, monotonic sequence, and caps at `ci_witnessed`. |
| Staff Engineer / DX Skeptic | Qwen3.6 Plus | Revise: missing offline mode, latency budget, emergency ergonomics, deterministic output, local preview, and early doctor. | Accepted; Block 13B and Block 14 now own these DX requirements before protected enforcement. |
| Compliance / Forensics Lead | DeepSeek V4 Pro | Revise: test evidence provenance, witness-before-merge, PR linkage, approval identity, retention enforcement, and replay semantics were insufficient. | Accepted; roadmap now adds test provenance, PR/MR and approval references, checkpoint-to-merge binding, retention enforcement, visible `not_assessed`, and a non-goal for deterministic re-execution. |

## Accepted Corrections

- Added operating modes:
  - observation mode;
  - advisory gate mode;
  - protected gate mode;
  - managed harness mode;
  - external audit mode.
- Added interception architecture table:
  - process wrapper;
  - adapter socket/API;
  - tool-level wrapper;
  - VCS/PR observer;
  - CI observer;
  - external witness.
- Added Block 13B before gate work, because gate contracts must not require
  evidence the active capture boundary cannot observe.
- Changed Block 14 from protected enforcement to advisory gate, explain, native
  override, deterministic output, preview, and emergency recording.
- Moved protected enforcement to Block 16 after signed checkpoints.
- Added signer isolation, monotonic sequence verification, replay resistance,
  and `ci_witnessed` cap for Block 15.
- Added source, PR, CI, merge-event, artifact, and checkpoint correlation for
  protected gates.
- Made managed harness enforcement explicitly opt-in; unmanaged harnesses retain
  observation-mode value without enrollment.
- Moved redaction and retention before broad adapter capture expansion.
- Added a cross-cutting pre-write safety floor from Block 13B onward: no raw
  prompts, model responses, source snippets, stdout/stderr bodies, tokens,
  secrets, or OIDC request tokens are persisted by default.
- Added test evidence provenance so `agent_reported` test claims cannot become
  executed test evidence.
- Added PR/MR, approval, and review references to the query and provenance
  surface.
- Added visible `not_assessed` rows and clarified that the forensics product is
  an evidence timeline, not guaranteed re-execution of arbitrary agent side
  effects.

## Review Output Files

- `archive/research/block-13-cto-buyer-minimax-review.md`
- `archive/research/block-13-platform-glm-review.md`
- `archive/research/block-13-ciso-kimi-review.md`
- `archive/research/block-13-staff-qwen-review.md`
- `archive/research/block-13-forensics-deepseek-review.md`
- `archive/research/block-13-cto-buyer-minimax-review-2.md`
- `archive/research/block-13-platform-glm-review-2.md`
- `archive/research/block-13-ciso-minimax-review-2.md`
- `archive/research/block-13-staff-qwen-review-2.md`
- `archive/research/block-13-staff-qwen-review-3.md`
- `archive/research/block-13-forensics-deepseek-review-2.md`

## Second-Pass Results

| Persona | Model | Result |
|---|---|---|
| CTO buyer | MiniMax-M2.7 | No critical or major findings. |
| Platform / Harness Owner | GLM-5.1 | No critical or major findings. |
| CISO / Adversarial Trust | MiniMax-M2.7 | No critical or major findings. Kimi K2P6 second pass returned an empty artifact and was replaced. |
| Staff Engineer / DX Skeptic | Qwen3.6 Plus | One remaining major: pre-write redaction must be a cross-cutting floor, not delayed to Block 18. Accepted and fixed. |
| Compliance / Forensics Lead | DeepSeek V4 Pro | No critical or major findings. |

Narrow Staff Engineer follow-up after the redaction correction:

- Model: Qwen3.6 Plus.
- Result: no critical or major findings.

## Convergence State

The second review pass had no remaining critical or major findings after the
Staff Engineer redaction correction was incorporated and re-reviewed. The
roadmap is now product-converged enough to drive SpecKit breakdown.
