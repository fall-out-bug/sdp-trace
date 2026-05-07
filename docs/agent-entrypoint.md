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

The active entrypoint is the Go CLI reported by `go run ./cmd/sdp-trace --help`.
Use this document to interpret proof scope; use `--help` to confirm exact flags
before adding or reviewing command examples.

- `go test ./...`
- `go run ./cmd/sdp-trace --help`
- `go run ./cmd/sdp-trace wrap --name <name> [--contract <file>] -- <command...>`
- `go run ./cmd/sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>`
- `go run ./cmd/sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>`
- `go run ./cmd/sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>`
- `go run ./cmd/sdp-trace doctor [--contract <file>]`
- `go run ./cmd/sdp-trace verify <run-dir>`
- `go run ./cmd/sdp-trace explain <run-dir>`
- `go run ./cmd/sdp-trace query --query <missing-evidence|capture-depth> <run-dir>`
- `go run ./cmd/sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>`
- `go run ./cmd/sdp-trace query-pack explain --result <file>`
- `go run ./cmd/sdp-trace assess --profile <adapter-capture|managed-harness|forensic-retention> [profile inputs]`
- `go run ./cmd/sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention> [profile inputs]`
- `go run ./cmd/sdp-trace assess explain --assessment-result <file>`
- `go run ./cmd/sdp-trace report --out <dir> <runs-root-or-run-dir>`
- `go run ./cmd/sdp-trace gate --out <file> <runs-root-or-run-dir>`
- `go run ./cmd/sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>`
- `go run ./cmd/sdp-trace release-proof --manifest <file> --out <file>`
- `go run ./cmd/sdp-trace validate-fixtures [root-dir]`

Do not add aliases, hidden switches, or workflow-specific wrappers as product
entrypoints unless this document and `--help` are updated in the same change.

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
