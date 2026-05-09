# Reviewer Entrypoint

Use this path for a first-time reviewer check in under five minutes. For the
full bilingual command/profile surface, see `docs/agent-entrypoint.md` and
`go run ./cmd/sdp-trace --help`.

For the demo-repository pilot evidence package, read
`examples/pilot-runs/opencode-minimax-kotlin-bazel/README.md` before inspecting
the retained package. Treat that package as an exact observed slice, not broad
OpenCode, MiniMax, Kotlin, or Bazel support.

## Verification Path

From a clean checkout, run:

1. `go test ./...`
2. `go run ./cmd/sdp-trace --help`
3. `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`
4. Create or inspect a run with `go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok`.
5. Verify that run with `go run ./cmd/sdp-trace verify <run-dir>`.
6. If documentation changed, compare command examples against `go run ./cmd/sdp-trace --help`.

External production trust is not part of this quick path. Use a live
`external_production_trust` profile path before making production trust claims.

## Exit Code Contract

Use `docs/agent-entrypoint.md` as the canonical state, trust-scope, authority
scope, and exit-code contract. The short exit summary is:

- `0`: `observed`, `pass`, or explicitly scoped `not_assessed`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

If any command returns exit code `3`, inspect the emitted reason and do not
upgrade the state in prose.

## Shared Command Surface

The current top-level command set is:

- `wrap`, `run`, `dry-run`, `preview`, `doctor`
- `harness observe`, `harness validate`, `harness summarize`
- `verify`, `explain`, `query`
- `query-pack`, `query-pack explain`
- `export cross-repo-posture`, `export cross-repo-posture explain`
- `assess`, `assess preview`, `assess explain`
- `report`, `gate`, `witness`, `release-proof`, `validate-fixtures`

Current assessment profiles:

- `adapter-capture`
- `managed-harness`
- `forensic-retention`
- `ci-artifact-observation`
- `authority-envelope`

Current witness kinds:

- `github-actions`
- `gitlab-ci`
- `buildkite`
- `customer-pki`

Air-gapped evidence is represented through customer policy/private-equivalent
guidance and fixtures, not as a separate `witness --kind` value.

Harness observation commands import and validate explicit local harness event
exports. They do not run OpenCode, GSD, MiniMax, GitHub, provider APIs, or any
other harness. Treat missing harness event families as `not_assessed` or
`cannot_verify`, not as feature delivery evidence.

## Dirty Checkout Behavior

- Clean checkout: verifier trust scope is the stated profile (`repo_baseline_structural`, `source_bound_local_release`, or `external_production_trust`).
- Dirty checkout without a command-supported dirty allowance: required clean-source checks may return `cannot_verify`.
- Dirty structural output may support only the `local_dirty_structural_only`
  authority scope.
- Do not use dirty output to conclude `source_bound_local_release` or `external_production_trust`.

## Not-Assessed Interpretation

`not_assessed` means the selected run did not assess that state.

What it allows:

- Continue using the command output with that state held back.
- Ask for the missing evidence or rerun against a scope that can assess it.

What it does not allow:

- Treating the state as passed.
- Using it as external trust closure.

## Gate, Witness, And Release Caveats

- `gate` emits verifier-derived facts and deterministic states. It does not own
  merge, release, readiness, degradation, override approval, or risk acceptance.
- `witness` binds available CI or customer-PKI evidence. A CI witness file is
  not external production trust, a transparency log, or a release approval by
  itself.
- `release-proof` can establish `source_bound_local_release` only when the
  source commit and manifest subjects match. It does not prove
  `external_production_trust`, `trusted_contract_release`, or
  `production_release_verified`.

## What You May State From Output

From verifier results, you may only state:

- Which command/profile was run.
- Which `result` or state values were produced.
- Whether the selected profile concluded with live `pass` or `observed`.
- Which states remain `not_assessed` or `cannot_verify`, with the emitted reason.

You may not state external production trust guarantees until
`external_production_trust` completes with live `pass` and
`production_release_verified` is supported by its dependent evidence chain.

## Quick Reference

| Goal | Command | Typical state boundary |
| --- | --- | --- |
| Local trace verification | `go run ./cmd/sdp-trace verify <run-dir>` | `observed` supports local structural assertions only |
| Missing evidence review | `go run ./cmd/sdp-trace query --query missing-evidence <run-dir>` | Missing evidence remains visible, not passed |
| Forensic package review | `go run ./cmd/sdp-trace query-pack explain --result <file>` | Explanation of retained evidence only |
| Managed harness review | `go run ./cmd/sdp-trace assess explain --assessment-result <file>` | Assessment facts; external policy owns block/allow |
| Authority envelope review | `go run ./cmd/sdp-trace assess --profile authority-envelope --authority-package <file> --out <file>` | Authority facts only; external policy owns consequences |
| CI/customer witness review | `go run ./cmd/sdp-trace witness --kind <kind> --out <file> <runs-root-or-run-dir>` | CI/customer-bound evidence, not production trust by itself |
| Source-bound release review | `go run ./cmd/sdp-trace release-proof --manifest <file> --out <file>` | Local source-bound proof only |

This entrypoint is intentionally minimal and is intended to prevent over-claiming
from reproducible verifier output.
