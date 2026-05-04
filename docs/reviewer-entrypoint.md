# Reviewer Entrypoint

Use this path for a first-time reviewer check in under five minutes.

## Verification Path

From a clean checkout, run:

1. `npm run verify:baseline`
2. `npm run verify:source-bound`
3. `npm run verify:external-trust`  
   Known blocker on clean checkout: this command exits with code `1` and reports `external_trust_profile_selected: fail` until external evidence is added. Treat this as an expected in-repo trust blocker, not a setup failure.

If external trust is not in scope, stop at step 2.

Exit code contract:
- `0`: `pass`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

Shared script-form command surface:
- `scripts/verify.sh --profile baseline|source-bound|external-trust [--json] [--allow-dirty] [--version]`
- `npm run verify:baseline` is equivalent to `scripts/verify.sh --profile baseline`
- `npm run verify:source-bound` is equivalent to `scripts/verify.sh --profile source-bound`
- `npm run verify:external-trust` is equivalent to `scripts/verify.sh --profile external-trust`

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

From `npm run verify:*` results, you may only state:

- Which profile command was run.
- Which `result` values were produced for required states.
- Whether a profile ended in live `result: pass` (for `repo_baseline_structural`, `source_bound_local_release`, or `external_production_trust`).

You may not state external production trust guarantees until `external_production_trust` completes with live `result: pass` and `production_release_verified` is supported by its dependent evidence chain.

## Commands and `--json`

Use text output for quick status checks.

Use one of these forms for JSON output:

- `npm run verify:baseline -- --json`
- `npm run verify:source-bound -- --json`
- `npm run verify:external-trust -- --json`
- `scripts/verify.sh --profile baseline --json`
- `scripts/verify.sh --profile source-bound --json`
- `scripts/verify.sh --profile external-trust --json`

Use JSON for:

- evidence packages,
- downstream audit records,
- exact state-level interpretation for `not_assessed` and `cannot_verify`.

## Quick Reference

| Goal | Command | Profile ID | Typical state boundary |
| --- | --- | --- | --- |
| Structural profile | `npm run verify:baseline` | `repo_baseline_structural` | live `result: pass` supports structural assertions only |
| Source-bound profile | `npm run verify:source-bound` | `source_bound_local_release` | live `result: pass` supports local source-bound assertions only |
| External production trust profile | `npm run verify:external-trust` | `external_production_trust` | live `result: pass` supports external-trust assertions only |

This entrypoint is intentionally minimal and is intended to prevent over-claiming from reproducible verifier output.
