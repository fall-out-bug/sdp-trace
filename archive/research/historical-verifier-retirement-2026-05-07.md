# Historical Verifier Retirement - 2026-05-07

Status: current retirement record for removed Node/npm/script verifier surfaces.

This document does not rewrite historical evidence. Earlier block ledgers and
self-trace examples may still name `npm`, `package.json`, `scripts/*.sh`, or
`.mjs` commands because those were the verifier surfaces used at the time.
Those references are historical-only and are not current closure evidence.

## Current Verifier Boundary

Current product verification is Go-first:

- `go test ./...`
- `jq empty schema/*.json examples/block14-gate/*.json examples/block15-checkpoint/*.json examples/block16-protected-gate/*.json examples/block17-managed-harness/*.json examples/block18-forensic-retention/*.json examples/block19-adapter-capture/*.json examples/contract-foundation/*.json examples/self-trace/*.json`
- `git diff --check HEAD`
- `go run ./cmd/sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out examples/contract-foundation/contract-release-verification.block18-19.json`

The active source-bound manifest must not list removed `package.json`,
`package-lock.json`, `scripts/*.sh`, or `.mjs` files as current artifacts.

## Retired Surfaces

The following surfaces are retired from current closure evidence:

- `package.json`
- `package-lock.json`
- `scripts/check-artifact-safety.sh`
- `scripts/finalize-source-bound-release.sh`
- `scripts/generate-local-dev-dsse.sh`
- `scripts/query-flight-recorder.mjs`
- `scripts/run-opencode-minimax-kotlin-bazel-proof.sh`
- `scripts/test-e2e-pilot-package.sh`
- `scripts/test-e2e-runner.sh`
- `scripts/test-flight-recorder.sh`
- `scripts/validate-contracts.sh`
- `scripts/validate-e2e-pilot-package.sh`
- `scripts/validate-json-schema.mjs`
- `scripts/validate-pilot-matrices.mjs`
- `scripts/validate-self-trace.sh`
- `scripts/verify-artifact-hashes.sh`
- `scripts/verify-contract-manifest.sh`
- `scripts/verify-flight-recorder.mjs`
- `scripts/verify-release-signature.sh`
- `scripts/verify-self-attestation.sh`

## Closure Rule

Historical ledgers that cite retired verifier commands remain valid only as
records of past work. A current closure claim must use the Go-first verifier
boundary above or explicitly state `cannot_verify` / `not_assessed` with a
reason. Retired command references must not be used to close T070, T169, T174,
or any future trust-sensitive task.
