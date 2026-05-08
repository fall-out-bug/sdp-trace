# Block 08: Agent and Human Usage Discovery

Status: implemented discovery + validation state; source-bound proof refreshed, external trust intentionally open
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Audience: implementation agents, technical executive, CIO, CISO, repository observers

## Purpose

Block 08 defines how an external operator (agent or human) discovers what `sdp-trace` can prove right now, without adding runtime behavior.

The block is scoped to:

- agent and human first-use entrypoints,
- one shared command surface,
- explicit trust scope and blocked claim boundaries.

It does not add CLI UX, agent workflow, or README rewrites.

Both entrypoints must converge on the same live verifier output. If they diverge, the contract must be corrected.

## Activation Gate and Trust Snapshot

Block 08 was activated from a clean-checkout verifier run with these outcomes:

- `repo_baseline_structural`: `pass`
- `source_bound_local_release`: `pass`
- `external_production_trust`: `fail`

The external-trust blocker is explicit in live `scripts/verify.sh --profile external-trust` state output:

- `external_trust_profile_selected`: `fail`
- `external_attestation_present`: `not_assessed`
- `external_identity_policy_matched`: `not_assessed`
- `transparency_or_audit_verified`: `not_assessed`
- `production_release_verified`: `fail`

This block does not resolve external production trust; it only makes the current blocker discoverable and reproducible through verifier runs.
That is sufficient for discovery. It is not sufficient for any customer-facing production-trust claim.
In a dirty checkout, `repo_baseline_structural` and `source_bound_local_release` may instead return `cannot_verify`, and that does not change the activation basis that was confirmed from a clean checkout.

## Agent Entrypoint

### 1) Profile Selection

Agents must select one verifier profile explicitly from existing commands:

- `repo_baseline_structural` for structural proof only.
- `source_bound_local_release` for local release binding evidence.
- `external_production_trust` for external signer/transparency/audit checks.

Select the profile from the claim under review, not from the harness or tool that happens to invoke the verifier.
Do not infer a profile from assumptions. The profile is an explicit input.

### 2) Command Contract

The command surface is fixed and pre-existing:

- `npm run verify:baseline`
- `npm run verify:source-bound`
- `npm run verify:external-trust`
- `scripts/verify.sh --profile baseline|source-bound|external-trust [--json] [--allow-dirty] [--version]`

No new command options are introduced in this block.

### 3) Evidence Emission Rules

Verifier output is the only claim source.

- `--json` is allowed for machine workflows; text output is valid for manual first-pass checks.
- A profile command emits verifier states with results (`pass`, `fail`, `not_assessed`, `cannot_verify`).
- `pass` is authoritative only for states required by the selected profile.
- `not_assessed` never means success.
- `fail` and `cannot_verify` must block any trust claim attached to that selected profile.

Human operators should treat checked-in `proof-summary` JSON as examples until re-verified with a live command.

## Human Entrypoint

### 5-minute Verification Path

1. Run `npm run verify:baseline`.
2. If required by role, run `npm run verify:source-bound`.
3. If external trust is required, run `npm run verify:external-trust` (from a clean checkout this is currently expected to exit `1`; see `External Trust Gap` below). Do not use `--allow-dirty` for this command.
4. Use the live states and exit code to decide what can be claimed. See Claim Discipline / Forbidden Claims below before making any trust or completion statement.
5. If any state shows `not_assessed`, read `Not-Assessed Interpretation` below before drawing conclusions.

If `npm run verify:external-trust` returns exit code `1` from a clean checkout, that is the expected current repository state. Do not attempt to debug it as a repository problem during first-pass discovery, and do not proceed with any external-trust claim. Only `repo_baseline_structural` and `source_bound_local_release` currently support live claims.
If any command returns exit code `3`, the required check could not verify. Clean the checkout first. Use `--allow-dirty` only for local structural checks, not for source-bound or external-trust conclusions.
Do not use `--allow-dirty` for first-pass verification when you need source-bound or external-trust conclusions.

### Proof Scope in Plain Terms

- `repo_baseline_structural` means repo structure and core contract checks passed live verification in the clean-checkout activation run.
- `source_bound_local_release` means local DSSE/source-bound checks passed live verification in the clean-checkout activation run.
- `external_production_trust` is currently blocked: a clean-checkout live run exits `1` because external production trust evidence is required and currently unresolved.

### Dirty Checkout Behavior

- Without `--allow-dirty`, any dirty checkout returns exit code `3` and cannot support profile-level claims.
- With `--allow-dirty`, `repo_baseline_structural` can still return `pass`, but only in `trust_scope=local_dirty_structural_only`.
- `source_bound_local_release` and `external_production_trust` still execute from a dirty checkout, but they do not yield an authoritative `pass`; they fail closed with exit code `3` / `cannot_verify` from `source_checkout_clean`.
- Treat `--allow-dirty` results as local development evidence only.
- `--allow-dirty` never supports source-bound or external-trust claims; see Claim Discipline / Forbidden Claims before making any statement from dirty-checkout output.

### Not-Assessed Interpretation

`not_assessed` means the selected profile did not assess the state in this run and does not imply proof exists. This state appears in live `--json` output when a state is outside the selected profile's scope or remains explicitly unassessed. It can mean either the profile excluded the state or upstream scope/blockers prevented assessment. Treat it as: no claim is available for this state in this run.

## Shared Command Surface

All onboarding uses the same execution path:

| User Goal | Command | Selected Profile | Exit Result |
| --- | --- | --- | --- |
| Structural check | `npm run verify:baseline` | `repo_baseline_structural` | `0` pass / `1` fail / `2` usage / `3` cannot verify |
| Local release binding | `npm run verify:source-bound` | `source_bound_local_release` | `0` pass / `1` fail / `2` usage / `3` cannot verify |
| External trust verification | `npm run verify:external-trust` | `external_production_trust` | `0` pass / `1` fail / `2` usage / `3` cannot verify |

Equivalent script form:

| Form | Example |
| --- | --- |
| baseline | `scripts/verify.sh --profile baseline --json` |
| source-bound | `scripts/verify.sh --profile source-bound --json` |
| external-trust | `scripts/verify.sh --profile external-trust --json` |

Use `--json` when you need the complete state set, including `not_assessed` fields and machine-readable reasons.

## External Trust Gap

Block 08 surfaces, but does not close, the external production gap.

- The current gap is real-time visible and blocks `external_production_trust`.
- The gap is closed only when all external states required by `external_production_trust` pass in one live run.
- Convergence is convergence of verifier scope, not promised capabilities outside current profile.

Current required external states and status:

| State | Result |
| --- | --- |
| `external_trust_profile_selected` | `fail` |
| `external_attestation_present` | `not_assessed` |
| `external_identity_policy_matched` | `not_assessed` |
| `transparency_or_audit_verified` | `not_assessed` |
| `production_release_verified` | `fail` |

The `not_assessed` states here follow the rule above: they are unassessed in this run and do not imply proof exists or will exist. The three `not_assessed` external states are downstream of `external_trust_profile_selected: fail`. `production_release_verified` fails because the required external-trust states do not pass in the same live run.

## Claim Discipline / Forbidden Claims

Block 08 inherits Block 07 claim authority and `docs/claim-authoring.md`.

Design-level requirements:

- Do not state support/compatibility/readiness/completion claims without live verifier-backed states.
- Do not state production trust under `source_bound_local_release` or `repo_baseline_structural`.
- Do not treat local structural output as external trust evidence.
- Use explicit `sdp-trace-claim` tags for machine-intended authoritative claims.

Forbidden claims:

- `trusted_contract_release=true` and `production_release_verified=true` outside external production trust pass.
- Capability-level verdict language that is not tied to a profile + state result.
- Any wording that turns a current blocker into an implied future trust claim, such as treating missing external evidence as already on the way to trusted closure.

Portability boundary:

- Agent and human entrypoints depend only on local verification output and this repository’s portable command set.
- No dependency is introduced on Codex, OpenCode, Beads, or any external harness runtime for discovery.

## Acceptance Conditions for T087-T090

- **T087:** A design review can proceed only if both agent and human sections are present, use the same profile vocabulary, and point to one shared command surface.
- **T088:** Entry points explicitly define profile selection, command contract, evidence emission rules, and forbidden claims.
- **T089:** Human entry point includes five-minute verification, dirty-checkout warning, `not_assessed` interpretation, and explicit external trust-gap explanation.
- **T090:** A fresh-agent pre-implementation review is complete before implementation planning and is recorded in the review ledger.
