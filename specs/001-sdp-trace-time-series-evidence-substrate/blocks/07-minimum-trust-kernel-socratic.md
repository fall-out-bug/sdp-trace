# Block 07 Socratic Review Synthesis

Status: draft; pi review incorporated
Parent: `07-minimum-trust-kernel.md`

## Review Inputs

This synthesis consolidates independent review threads focused on:

- Minimum Trust Kernel design
- security and forgery resistance
- verifier UX and testing
- current repository claim consistency

No reviewer was asked to edit files.

## Core Challenge

If real trust is deferred to a future block, `sdp-trace` looks like a structured logger. A technical executive or CISO can reasonably reject it because the product is asking customers to trust a trust substrate that cannot yet prove its own claims.

## Resolved Design Direction

Use proof profiles and live verifier authority instead of a maturity ladder.

A maturity ladder such as `L0-L4` is rejected because it compresses orthogonal claims into marketing shorthand. Schema validity, self-trace, source-bound release, observed slice evidence, and external production trust are different claims. They must not be collapsed into one level label.

## Socratic Questions and Resolutions

### Q1: What is the authoritative artifact when docs and JSON disagree?

Resolution: live verifier output wins. Persisted proof summaries, release verification JSON, reports, matrices, task checkboxes, Beads mirrors, and review ledgers are indexes or audit artifacts. They cannot override a current verifier run.

### Q2: Can a block be complete if validation currently fails?

Resolution: no. A task that claims a command was verified is false if the current tree fails that command. Historical success may be recorded as prior evidence, but it cannot close the current checkout.

### Q3: Can local DSSE or private key possession establish production trust?

Resolution: no. Local attestation can support source-bound local release candidate proof. Production trust requires a named external profile and matching external evidence.

### Q4: Why is `production_release_verified` dangerous if `trusted_contract_release` remains false?

Resolution: it is a hidden trust claim. Any production verification boolean must be constrained by the same required state set as `trusted_contract_release`, or consumers can launder trust through a secondary field.

### Q5: Should incomplete Block 06 evidence fail baseline validation?

Resolution: no. Incomplete optional slices can be honest artifacts. They should validate as incomplete packages while failing any completed-proof profile.

### Q6: What prevents forged `matched` statuses?

Resolution: schema shape is insufficient. The canonical verifier must compute matched states from canonical subjects. Checked-in proof summaries and verification JSON are untrusted input unless regenerated, re-verified, or externally signed.

### Q7: Should pilot artifacts be part of the Minimum Trust Kernel manifest?

Resolution: not initially. A too-wide manifest increases stale proof risk. The kernel should first prove the contract set needed by downstream consumers, then let pilot evidence reference a trusted contract version.

### Q8: What must an external consumer reject without reading prose?

Resolution: any stale, forged, mismatched, missing, not assessed, expired, rollback-suspected, wrong signer, wrong identity, wrong profile, or contradicted trust claim.

## Findings to Carry Into Review Ledger

1. critical: Current self-attestation baseline is not reproducible.
2. critical: Current Block 04 delivery state overclaims closure while source commit artifact verification mismatches.
3. critical: `production_release_verified.value: true` is not independently constrained enough to prevent hidden trust claims.
4. major: Baseline validation is coupled to optional incomplete pilot proof.
5. major: E2E package validation conflates package validity with completed proof.
6. major: JSON Schema cannot prove that checked-in `matched` states were honestly computed.
7. major: Evidence/provenance cross-reference integrity is not fully verifier-enforced.
8. major: A verifier cannot prove its own checkout is trustworthy without an external trust anchor; the spec must name the bootstrap assumption.
9. major: Prose claim scanning is not robust; authoritative claims need machine-readable claim tags.
10. minor: The first kernel manifest surface should be smaller than the current broad manifest.

## Design Consequence

The next work must not be "finish Block 06" or "add more examples". The next work must establish live verifier authority, then repair stale claims under that authority.

Implementation may start only after the spec and implementation plan have been reviewed, findings are recorded, and critical/major spec-gate findings are closed.
