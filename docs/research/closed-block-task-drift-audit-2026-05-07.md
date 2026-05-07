# Closed Block And Task Drift Audit - 2026-05-07

Status: current drift audit for closed SpecKit tasks and block closure claims.
This is not source-bound proof and does not make `trusted_contract_release`
true.

## Scope

Audited surfaces:

- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/*.md`
- `docs/research/block-*-implementation-evidence.md`
- current source-bound manifest verification output
- current Go and JSON syntax verification output

Task ledger state after retirement:

- Closed tasks: T070 and T169-T174 are closed for the current Go-first
  repository boundary.
- Remaining external production trust: `not_assessed`.

## Commands

Run from `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-block-19-drift-closure`:

```bash
rtk go test ./...
rtk jq empty schema/*.json examples/block14-gate/*.json examples/block15-checkpoint/*.json examples/block16-protected-gate/*.json examples/block17-managed-harness/*.json examples/block18-forensic-retention/*.json examples/block19-adapter-capture/*.json examples/contract-foundation/*.json examples/self-trace/*.json
go run ./cmd/sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out examples/contract-foundation/contract-release-verification.block18-19.json
git diff --check HEAD
```

Observed states:

- Go tests: pass, 151 tests across 15 packages.
- JSON syntax checks: pass for schemas and Block 14-19, contract-foundation,
  and self-trace JSON examples.
- Whitespace diff check: pass.
- Source-bound local release proof: fail, exit code 1 as expected for a
  verifier-derived fail state.
- The checked-in fail artifact is generated against the latest source commit
  assessed before the artifact-only update. Its `source_commit` is the assessed
  source tree for this diagnostic fail, not a claim that the PR head is a
  trusted release.

## Findings

### F1 - Source-Bound Manifest Drift

Severity: critical.

The current source-bound verifier can assess the checked-out manifest, and it
reports `release_verification_state: fail`. The generated fail artifact records
166 manifest artifacts checked, 20 missing artifacts, and 15 mismatched
artifacts.

The missing artifacts were the historical Node/script validation path:

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

The mismatched artifacts include current docs, schemas, flight-recorder
fixtures, and SpecKit files changed since the last manifest refresh.

Disposition: accepted and fixed for the active manifest subject set. Retired
Node/npm/script paths are removed from the active manifest, current artifact
hashes are refreshed, and the Go-first local release-proof output records
source-bound artifact checks separately from external production trust.

### F2 - Historical Closed Tasks Reference Removed Tooling

Severity: major.

Several closed tasks and block ledgers still cite `npm`, `.mjs`, `package.json`,
or `scripts/*` verifier commands as the closure mechanism. That contradicts
the current repository rule that active product verification is Go-first and
the checked-out repository no longer contains those scripts.

Examples include T058, T077, T081, T082, T086, Block 01 validation notes,
Block 05/06 review ledgers, Block 07 implementation-plan slices, and Block 08
entrypoint/review-ledger commands.

Disposition: accepted and fixed. `docs/research/historical-verifier-retirement-2026-05-07.md`
is the current retirement record. Historical ledgers keep their original command
text as historical evidence, but those commands are no longer accepted as
current closure evidence.

### F3 - Block 18/19 PR Drift Follow-Ups

Severity: major before this PR; non-blocking after this PR if verification
continues to pass.

Block 18/19 evidence notes correctly recorded CI as `not_assessed`, source-bound
release proof as open, and PR #9 minor review notes as follow-up work. The
current PR adds:

- repo-tracked CI/check policy and GitHub Actions workflow;
- adapter-capture regression tests for duplicate empty correlation refs and
  redacted capture with a visible retention cap;
- README clarification for intentionally byte-identical pass fixtures;
- a repository drift-to-task rule in `AGENTS.md`.

Disposition: accepted and fixed for T170, T171, and T172. T169 remains open
because the current source-bound proof is a real fail.

## Block Closure Drift Review

Closed or closure-like block claims were reviewed against current machine
evidence.

| Surface | Current drift state | Disposition |
| --- | --- | --- |
| Block 01 contract foundation | Historical validation used Node/scripts and stale manifest proof. Current product proof remains bounded by self-trace/source-bound state. | Covered by T174 and T169. |
| Block 02 self-trace | JSON artifacts parse, but historical script-based validation references are stale. | Covered by T174. |
| Block 03 self-attestation | Stored result still references the old source-bound pass commit; current manifest proof fails. | Covered by T169 and T174. |
| Block 04 release finalization | Already explicitly reopened as T070 stale closure. | Leave T070 open. |
| Block 05 customer pilot | Closure ledgers cite `npm run validate` and script validators that are absent. | Covered by T174. |
| Block 06 first product proof | Closure ledgers cite absent script validators and old manifest refresh. | Covered by T174. |
| Block 07 trust kernel | Correctly records Block 04 stale closure and source-bound boundaries, but implementation-plan command surfaces are stale. | Covered by T174. |
| Block 08 usage discovery | Status claims source-bound proof refreshed through removed verifier commands; current manifest proof fails. | Covered by T169 and T174. |
| Blocks 09-19 implementation surfaces | Current Go tests and JSON syntax checks pass for active Go-first Block 14-19 artifacts. Source-bound release proof remains failed. | Keep source-bound closure blocked under T169. |

## Closure Boundary

This audit closes the review pass over closed blocks and tasks for the current
Go-first repository boundary. The current honest state is:

- `repo_go_tests`: pass
- `json_syntax`: pass
- `diff_whitespace`: pass
- `source_bound_local_release`: pass for active manifest artifact digest checks
- `external_production_trust`: not_assessed
- `trusted_contract_release`: false
