# Agent Entrypoint

Use these commands and profile IDs to select what `sdp-trace` can prove without introducing harness assumptions.

## Profile Selection

Each assertion is anchored to one of these profile IDs:

- `repo_baseline_structural`
- `source_bound_local_release`
- `external_production_trust`

Do not infer profile from role. Choose the profile directly from the claim you need:

- `repo_baseline_structural` for structural command/gate-set integrity checks.
- `source_bound_local_release` for local DSSE/source-bound checks.
- `external_production_trust` for external production trust checks.

## Command Contract

Only this command set is part of the active entrypoint:

- `go test ./...`
- `go run ./cmd/sdp-trace --help`
- `go run ./cmd/sdp-trace validate-fixtures <fixture-dir>`
- `go run ./cmd/sdp-trace wrap --name <name> -- <command...>`
- `go run ./cmd/sdp-trace verify <run-dir>`
- `go run ./cmd/sdp-trace explain <run-dir>`

Use these as the canonical Block 10-compatible commands.

Do not add aliases, new switches, or workflow-specific wrappers in this block.

## Trust Scope Vocabulary

- Dirty-checkout baseline output is only valid for `local_dirty_structural_only`.
- Treat `local_dirty_structural_only` as `repo_baseline_structural` evidence with dirty checkout constraints.
- Do not use dirty-checkout output to close `source_bound_local_release` or `external_production_trust`.

## Evidence Emission Rules

Text output is sufficient for a first-pass pass/fail look.

Use verifier JSON output when you need full states, reasons, and `result`
values (`observed`, `fail`, `not_assessed`, `cannot_verify`) for
decisioning.

## Exit Code Contract

- `0`: `observed` or `not_assessed`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

- `observed`: verifier evidence met required checks for the selected local profile.
- `fail`: verifier evidence conflicted or was insufficient for one or more required checks in the selected profile.
- `not_assessed`: state was not assessed in this run (scope omission or upstream profile blocker); it does not imply success or evidence.
- `cannot_verify`: verifier could not execute a required check and the profile cannot be concluded from that run.

A checked-in `proof-summary` JSON is an audit artifact, not authority.

Authority is replayed only from live Go verifier output and the canonical command/state contract above.

## Forbidden Claims

Do not emit these in this repo surface:

- `external_production_trust=true` without a live `external_production_trust` profile pass.
- `trusted_contract_release=true` without live external trust closure.
- `production_release_verified=true` outside a concluded `external_production_trust` run.
- Claims that treat `repo_baseline_structural` or `source_bound_local_release` outputs as production trust.
- Dirty-checkout structural output as source-bound or external-trust evidence.
