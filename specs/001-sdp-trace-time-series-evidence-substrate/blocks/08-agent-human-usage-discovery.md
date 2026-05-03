# Block 08: Agent and Human Usage Discovery

Status: unblocked for design review after Block 07 `closed_for_block_08_activation`
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Audience: implementation agents, CTO, CIO, CISO, repository observers

## Purpose

Block 08 will make `sdp-trace` understandable from two first-use paths:

1. An agent needs to know which verifier profile to run, which evidence to emit, and which claims it must not make.
2. A human reviewer needs to know what the repository currently proves, what remains `not_assessed`, and how to verify that in minutes.

Both paths must converge on the same verifier output. If the agent path and the human path disagree, the verifier wins and the documentation is wrong.

## Activation Gate

Do not start Block 08 design review, implementation planning, or code changes until Block 07 has either:

- passed the required trust profiles with live verifier evidence, or
- remained open with a verifier-derived blocking state that Block 08 can explain honestly.

Block 08 must not become a workaround for missing Block 07 proof.

Current activation basis:

- `repo_baseline_structural`: pass
- `source_bound_local_release`: pass
- `external_production_trust`: fail with explicit missing external evidence states

This is sufficient for discovery because the external production trust blocker is explicit. It is not sufficient for any customer-facing production trust claim.

## Future Scope

- Agent entrypoint: portable usage contract independent of Codex, Claude, OpenCode, Pi, GitHub, Beads, or any harness runtime.
- Human entrypoint: short CTO/CISO/reviewer path that explains proof scope, dirty checkout behavior, incomplete slices, and external trust state.
- Shared command surface: one documented route to current verifier profiles and proof-summary states.
- Claim discipline: no onboarding doc may claim support, readiness, compatibility, production trust, or completion without verifier-backed state.
- Decomposition check: do not add more than one Block 08 skill unless the module is split.

## Non-Goals Before Block 07

- No new agent workflow.
- No new onboarding claims.
- No CLI UX changes.
- No README rewrite that suggests the product is easier to trust than Block 07 proves.
