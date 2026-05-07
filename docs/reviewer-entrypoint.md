# Reviewer Entrypoint

Use this path for a first-time reviewer check in under five minutes.
For the full current command surface, see `docs/agent-entrypoint.md` and
`go run ./cmd/sdp-trace --help`.

## Verification Path

From a clean checkout, run:

1. `go test ./...`
2. `go run ./cmd/sdp-trace --help`
3. `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`
4. Create or inspect a run with `go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok`.
5. Verify that run with `go run ./cmd/sdp-trace verify <run-dir>`.

External trust is not part of this quick path. Use the documented
`external_production_trust` profile path before making production trust claims.

Exit code contract:
- `0`: `observed` or `not_assessed`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

Shared command surface:
- `go test ./...`
- `go run ./cmd/sdp-trace validate-fixtures <fixture-dir>`
- `go run ./cmd/sdp-trace wrap --name <name> -- <command...>`
- `go run ./cmd/sdp-trace verify <run-dir>`
- `go run ./cmd/sdp-trace explain <run-dir>`

If any command returns exit code `3`, clean the checkout first for source-bound/external conclusions.

`--json` is required when you need exact state output for any `not_assessed` or `cannot_verify` path.

## Dirty Checkout Behavior

- Clean checkout: verifier trust scope is the stated profile (`repo_baseline_structural`, `source_bound_local_release`, or `external_production_trust`).
- Dirty checkout without `--allow-dirty`: returns `cannot_verify` for required checks when a required clean source is needed.
- Dirty checkout with `--allow-dirty`: may still support `repo_baseline_structural` in `local_dirty_structural_only` scope only.
- Do **not** use `--allow-dirty` to conclude `source_bound_local_release` or `external_production_trust`.
- Dirty source-bound or external-trust conclusion from `--allow-dirty` is not allowed.

## Not-Assessed Interpretation

`not_assessed` means the selected profile/state run did not assess that state.

What it allows:

- Continue using the command output, but with that state held back.
- Ask for the missing evidence or run against a scope that can assess it.

What it does not allow:

- Treating the state as passed.
- Using it as external trust closure.

## External Trust Gap

In this repository, `external_attestation_present`, `external_identity_policy_matched`, and `transparency_or_audit_verified` are downstream of `external_trust_profile_selected: fail`.

When that state is `fail`, those states are blocked before independent assessment, and `production_release_verified` remains failed by profile-closure logic.

## What You May State from Output

From verifier results, you may only state:

- Which profile command was run.
- Which `result` values were produced for required states.
- Whether a profile ended in live `result: pass` (for `repo_baseline_structural`, `source_bound_local_release`, or `external_production_trust`).

You may not state external production trust guarantees until `external_production_trust` completes with live `result: pass` and `production_release_verified` is supported by its dependent evidence chain.

## Commands and `--json`

Use text output for quick status checks.

Use `go run ./cmd/sdp-trace verify <run-dir>` for JSON verifier output.

Use JSON for:

- evidence packages,
- downstream audit records,
- exact state-level interpretation for `not_assessed` and `cannot_verify`.

## Quick Reference

| Goal | Command | Profile ID | Typical state boundary |
| --- | --- | --- | --- |
| Local trace verification | `go run ./cmd/sdp-trace verify <run-dir>` | `local_observed` | live `result: observed` supports local structural assertions only |

This entrypoint is intentionally minimal and is intended to prevent over-claiming from reproducible verifier output.
