# Canonical Overclaim Checklist

This is the canonical overclaim and forbidden-interpretation checklist for
`sdp-trace`. Other docs may summarize or link here; they must not contradict it.

## What sdp-trace Does Not Decide

`sdp-trace` records evidence and gaps. It does **not** decide:

- merge approval
- release readiness
- risk acceptance or override
- degradation decisions
- production trust authority
- whether a team may ship

Policy decisions belong to CI, release governance, customer governance, or
another external policy consumer that already owns the decision.

## Forbidden Claims

Do not emit these without the required live evidence:

1. `external_production_trust=true` without a live
   `external_production_trust` profile pass.
2. `trusted_contract_release=true` without live external trust closure.
3. `production_release_verified=true` outside a concluded
   `external_production_trust` run.
4. Claims that treat `repo_baseline_structural` or
   `source_bound_local_release` outputs as production trust.
5. Dirty-checkout structural output as source-bound or external-trust evidence.

## What You May State From Verifier Output

From verifier results, you may only state:

- Which command and profile were run.
- Which `result` or state values were produced.
- Whether the selected profile concluded with live `pass` or `observed`.
- Which states remain `not_assessed` or `cannot_verify`, with the emitted reason.

You may **not** state external production trust guarantees until
`external_production_trust` completes with live `pass` and
`production_release_verified` is supported by its dependent evidence chain.

## Command-Specific Caveats

- `pr-review` emits review-record evidence. It reports coverage and finding
  states, but it does not approve, merge, mark ready, release, accept risk, or
  replace human approval.
- `gate` emits verifier-derived facts and deterministic states. It does not own
  merge, release, readiness, degradation, override approval, or risk acceptance.
- `witness` binds available CI or customer-PKI evidence. A CI witness file is
  not external production trust, a transparency log, or a release approval by
  itself.
- `release-proof` can establish `source_bound_local_release` only when the
  source commit and manifest subjects match. It does not prove
  `external_production_trust`.
- `assess` emits assessment facts. Block/allow, authority, and readiness
  decisions remain downstream.

## Trust Scope And Authority Scope

Keep these vocabularies separate:

- **Result state**: the verifier outcome (`observed`, `pass`, `fail`, `not_assessed`, `cannot_verify`).
- **Trust scope**: the evidence boundary (`local_observed`, `ci_witnessed`, `external_witnessed`).
- **Authority scope**: the reporting boundary for a committed package (`demo_pilot_only`, `local_dirty_structural_only`).

A checked-in proof JSON is an audit artifact, not authority. Authority is replayed
only from live Go verifier output and the canonical state contract.
