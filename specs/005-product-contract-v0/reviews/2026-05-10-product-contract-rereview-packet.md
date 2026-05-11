# Re-Review Packet: Product Contract v0

Date: 2026-05-10
Scope: `specs/005-product-contract-v0/` plus linkage text in
`specs/003-agent-supply-chain-roadmap/`.

## Review Target

Review the revised Product Contract v0 after the first full pi-review.

Primary files:

- `specs/005-product-contract-v0/spec.md`
- `specs/005-product-contract-v0/plan.md`
- `specs/005-product-contract-v0/example.md`
- `specs/005-product-contract-v0/example-local-baseline.md`
- `specs/005-product-contract-v0/traceability.md`
- `specs/005-product-contract-v0/tasks.md`

Linkage files:

- `specs/003-agent-supply-chain-roadmap/plan.md`
- `specs/003-agent-supply-chain-roadmap/tasks.md`

Prior review ledger:

- `specs/005-product-contract-v0/reviews/2026-05-10-full-pi-review.md`

## Revision Focus

The revision tries to make Product Contract v0 less slogan-like and less
approval-system-like:

- Change Evidence Packet v0 is the first buyer-facing artifact.
- P0 product work must cite packet rows and show a target transition.
- Substrate, adapter, and discovery work are useful, but not P0 product
  progress unless they improve packet rows.
- The canonical packet is Markdown plus evidence bundle manifest.
- Local/self-hosted enterprise operation is a baseline, not a later luxury.
- Evidence states preserve `not_assessed`, `cannot_verify`, `partial`,
  `unsupported`, and `not_integrated`.
- Source-strength classes are categorical and must not become a score.
- Theater findings are derived findings with trigger evidence; they are not a
  global failure score.
- Signed attestation is additive and does not upgrade underlying row states.

## Review Question

Is the revised Product Contract v0 good enough to take back to the CTO/user for
explicit approval as the product contract direction, before any implementation
work?

Do not review implementation code. There is no implementation in this slice.
