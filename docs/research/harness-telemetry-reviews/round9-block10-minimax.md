# Round 9: Block 10 MiniMax Review

Status: pi review output summary
Date: 2026-05-05
Model: `minimax/MiniMax-M2.7`
Role: CISO / Adversarial Trust Reviewer

Reviewed:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-design.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md`
- `docs/research/agentic-sdlc-evidence-substrate-v4-brief.md`

This is a review artifact, not source-bound proof or closure evidence.

## Verdict

`ACCEPTABLE_WITH_GAPS`

The design can start v0 implementation, but local trust wording and
authority boundaries required correction before handoff.

## Critical Findings

### C1: Ephemeral local keys provide no integrity assurance

Local keys generated and held by the recorder process cannot provide
strong integrity if the recorder is compromised or controlled from
startup.

Disposition: accepted. Block 10 now states that local-only traces are
sequence-continuous under honest-recorder assumptions, not externally
anchored integrity proof.

### C2: First event hash has no external anchor

A locally generated first event, run id, and nonce can anchor a fabricated
local chain if the recorder is hostile from startup.

Disposition: accepted. Block 10 now discloses that the first event uses
local recorder values and gate-grade trust requires CI/external witness
binding.

### C3: `human_signed` accepted insecure local key defaults

Local named keys readable by the agent workspace/process could allow
impersonation.

Disposition: accepted. Block 10 now requires a configured signing
profile and treats unprotected local named keys as declaration or
customer-accepted local risk, not strong human presence proof.

### C4: Contract pinning in CI config was ambiguous

Repo-committed CI YAML is insufficient if the assessed agent can modify
it in the same change.

Disposition: accepted. Block 10 now requires CI-secret/config outside
the PR diff, external policy service, human-signed digest, or customer
PKI equivalent for gate-grade pinning.

## Major Findings

- Witness independence needs finer states than `same_job` and
  `same_process`.
- Redaction was asserted without a failure mode.
- Local clock state was underspecified.
- `expected_run_absent` only works where CI/preflight expects a run.
- Adapter registration lacked authority policy and could allow fake
  adapters to self-upgrade trust.

Disposition: accepted. Block 10 now defines witness independence states,
redaction failure behavior, local clock limitations, adapter authority
policy, and expected-run limits.

## Minor Findings

- Run id generation method missing.
- Schema version compatibility policy missing.
- Socket permission model incomplete.
- Signal/exit semantics incomplete.

Disposition: accepted. Block 10 now specifies UUIDv7/random-equivalent
run ids, strict schema compatibility, socket permission limits, and
termination handling requirements.
