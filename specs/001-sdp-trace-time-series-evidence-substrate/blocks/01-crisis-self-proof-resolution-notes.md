# Crisis Self-Proof Resolution Notes

Date: 2026-04-30

Trigger: CTO objection: if `sdp-trace` cannot directly prove its own development and contract integrity, it cannot prove anything for a customer.

## Socratic Critic Results

Two clean-context `pi` critics reviewed the current SpecKit package and Block 01 implementation artifacts:

- `01-crisis-self-proof-kimi-critic.json`
- `01-crisis-self-proof-glm-critic.json`

Both critics returned `verdict: "needs_revision"`.

## Accepted Findings

- Self-trace tasks T020-T026 were open, so `sdp-trace` had not yet used its own contracts to describe its own development.
- Block 01 was worded too strongly as a trusted foundation. It is only contract scaffolding until self-trace and self-attestation pass.
- The local DSSE/OpenSSL proof demonstrates envelope binding and digest verification mechanics, not production Sigstore/Rekor/OIDC identity.
- `source_commit: "working-tree-block-01"` is acceptable only in local development fixtures and cannot support a trusted product claim.
- Metric catalog T011 remains open, so pilot metric comparability is not yet proven.

## Spec Changes Made

- Added the Non-Negotiable Self-Proof Rule to `spec.md`.
- Added FR-029 through FR-031 and SC-016 through SC-018 to require self-trace, separate proof states, and immutable source references before customer pilot claims.
- Changed `plan.md` to treat self-proof as the first product proof, with separate self-trace and self-attestation milestones.
- Promoted Phase 5 to Priority P0 and added T052-T056 for self-trace validation, crisis review evidence, accountability, metrics, and negative fixtures.
- Added Phase 5A with T057-T061 for self-attested contract release proof.
- Reworded Block 01 status and purpose so it no longer claims product proof.
- Updated README and CTO briefs to state that customer trust is blocked until this repository traces itself.

## Current Product State

`sdp-trace` has validated contract scaffolding.

`sdp-trace` has not yet proven product viability.

The next implementation block must be Self-Trace v0. Pilot matrix work is blocked until self-trace passes, and production trust claims are blocked until self-attestation records the actual proof state without hiding missing external trust anchors.
