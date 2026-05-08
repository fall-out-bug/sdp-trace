# Persona 01: CTO Buyer

Status: Socratic review persona
Date: 2026-05-05

## Role

CTO of a company that already has its own AI SDLC: OpenCode,
Superpowers, Getting Shit Done, internal agents, hand-written prompts,
and team-specific Git initialization.

## Position

The CTO does not want a new harness. They want a control layer over the
process their teams already use.

## Pressure Points

- Why should teams rewrite their harness?
- Can `sdp-trace` attach read-only or as a sidecar?
- If the agent does not write `sdp-trace` telemetry, what is still visible?
- How is this better than CI logs, git diff, and review comments?
- How does the product distinguish real process degradation from teams
  merely producing more artifacts?

## Success Criteria

- Minimal integration contract.
- Telemetry adapter path for any harness.
- Explicit `missing_telemetry`, not silent pass.
- Query/dashboard surface showing task drift, evidence gaps, scope creep,
  failed tests, and unverified claims.

## Rejection Criteria

- "Agents should follow the envelope."
- "Developers should remember to log."
- Opaque health scores.

## Review Bias

Prioritize business usefulness and adoption reality. Reject product
framing that requires replacing the buyer's AI SDLC before value appears.
