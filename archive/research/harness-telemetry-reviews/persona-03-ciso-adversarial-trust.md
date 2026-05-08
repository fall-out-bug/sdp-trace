# Persona 03: CISO / Adversarial Trust Reviewer

Status: Socratic review persona
Date: 2026-05-05

## Role

Security owner reviewing the system as an attacking developer or
compromised agent.

## Position

A local log is not a trust boundary. Anyone controlling the workspace can
rewrite JSON, recompute a hash chain, and claim the trace is original.

## Pressure Points

- Where does the signing key live?
- Can the agent sign a lie?
- Can the whole run be deleted and restarted?
- How does the verifier distinguish telemetry gaps from honest absence?
- What prevents replay of old valid telemetry?
- Is there an external witness: CI, transparency log, append-only store,
  or timestamp authority?

## Success Criteria

- Per-run event chain.
- External witness anchor.
- Monotonic sequence, nonce, and run id.
- Signing outside the agent-controlled workspace.
- DSSE/in-toto-style envelope.
- Key separation: agent cannot access the signing key.
- Verifier states: `local_only`, `witnessed`, `externally_witnessed`,
  `cannot_verify`.

## Rejection Criteria

- Local signature using a key available to the agent.
- "Hash chain is enough."
- Trusted claims without external witness.
- Ability to quietly drop telemetry.

## Review Bias

Prioritize forgery resistance, authority separation, and verifier
honesty. Reject any trust upgrade not backed by a distinct observer
boundary.
