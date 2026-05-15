# Reviewer Entrypoint

Use this path for a first-time reviewer check in under five minutes. For the
full bilingual command/profile surface, see `docs/agent-entrypoint.md` and
`sdp-trace --help`.

For the demo-repository pilot evidence package, read
`examples/pilot-runs/opencode-minimax-kotlin-bazel/README.md` before inspecting
the retained package. Treat that package as an exact observed slice, not broad
OpenCode, MiniMax, Kotlin, or Bazel support.

## Quick Reference — I Have A Run Directory, What Now?

| Goal | Command | Typical state boundary |
| --- | --- | --- |
| Verify the run | `sdp-trace verify <run-dir>` | `observed` supports local structural assertions only |
| Find missing evidence | `sdp-trace query --query missing-evidence <run-dir>` | Missing evidence remains visible, not passed |
| Build a forensic package | `sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out query-pack.json` | Limited by retained/redacted evidence |
| Explain the run | `sdp-trace explain <run-dir>` | Explanation does not upgrade trust scope |
| Assess adapter capture | `sdp-trace assess --profile adapter-capture --out assessment.json --run <run-dir>` | Can fail if adapter events are absent |
| Assess managed harness | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run <run-dir> --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` | Policy owns block/allow |
| Assess forensic retention | `sdp-trace assess --profile forensic-retention --out assessment.json --run <run-dir> --redaction-policy redaction.json` | Digest-only or missing retention may fail |
| Assess CI artifacts | `sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` | Facts only; checked-in claims cannot satisfy `ci_uploaded` |
| Assess authority envelope | `sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` | Authority facts only; policy owns consequences |
| Build a report | `sdp-trace report --out .sdp-trace-report .sdp-trace-runs` | Packages observed data and gaps |
| Witness CI run | `sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` | CI witness is not production trust by itself |
| Check release proof | `sdp-trace release-proof --manifest <file> --out release-proof.json` | Local source-bound proof only |
| Run automated PR review | `sdp-trace pr-review check --out review --repo-id <safe-id> --change-ref pr-123 --base <sha> --head <sha> --diff change.diff --profile examples/pr-review/trust-sensitive-default.profile.json` | Review-record completeness only; not merge approval |

For output locations, see [`docs/output-location-map.md`](output-location-map.md).
For profile selection, see [`docs/profile-selection-guide.md`](profile-selection-guide.md).

## Verification Path

From a clean checkout, run:

1. `go test -count=1 ./...`
2. For a source checkout, define `sdp_trace() { go run ./cmd/sdp-trace "$@"; }`.
   After installing a release binary, use `sdp-trace` directly.
3. `sdp_trace --help` for a source checkout, or `sdp-trace --help` for a release binary.
4. `sdp_trace validate-fixtures examples/agentic-sdlc` for a source checkout, or `sdp-trace validate-fixtures examples/agentic-sdlc` for a release binary.
5. Create or inspect a run with
   `sdp_trace wrap --name smoke --output-dir .sdp-trace-runs/smoke -- /bin/echo ok`
   for a source checkout, or the same command with `sdp-trace` after installing
   a release binary.
6. Verify that run with `sdp_trace verify .sdp-trace-runs/smoke` or
   `sdp-trace verify .sdp-trace-runs/smoke`.
7. If documentation changed, compare command examples against `sdp_trace --help`
   or `sdp-trace --help`.

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

## Reviewer Command Surface

This is the reviewer subset for fast orientation. The full command surface is
authoritative in [Agent Entrypoint](agent-entrypoint.md) and `sdp-trace --help`.
When reviewing command docs, compare against both.

- `version`, `wrap`, `run`, `dry-run`, `preview`, `doctor`
- `command-surface`
- `install repo-observer`
- `interaction relay`, `interaction import-transcript`, `interaction summarize`
- `observe setup`, `observe collect`, `observe session`
- `harness observe`, `harness validate`, `harness summarize`
- `envelope summarize`
- `verify`, `explain`, `query`
- `query-pack`, `query-pack explain`
- `export cross-repo-posture`, `export cross-repo-posture explain`, `export telemetry`
- `assess`, `assess preview`, `assess explain`
- `report`, `gate`, `witness`, `release-proof`, `pr-review`
- `packet build-pr`, `packet build-github`, `packet validate`, `packet check-demo`, `packet render`
- `validate-fixtures`

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

First-run observation commands use a session profile to bound setup and
collection. `observe setup` writes setup metadata before delivery,
`observe collect` normalizes declared harness output after the normal harness
command, and `observe session` is a convenience wrapper for one controlled
command. They do not inject stdin, relay prompts, retain stdout/stderr bodies by
default, or turn missing harness output into a pass.

## Dirty Checkout Behavior

- Clean checkout: verifier trust scope is the stated profile (`repo_baseline_structural`, `source_bound_local_release`, or `external_production_trust`).
- Dirty checkout without a command-supported dirty allowance: required clean-source checks may return `cannot_verify`.
- Dirty structural output may support only the `local_dirty_structural_only`
  authority scope.
- Do not use dirty output to conclude `source_bound_local_release` or
  `external_production_trust`.

## Not-Assessed Interpretation

`not_assessed` means the selected run did not assess that state.

What it allows:

- Continue using the command output with that state held back.
- Ask for the missing evidence or rerun against a scope that can assess it.

What it does not allow:

- Treating the state as passed.
- Using it as external trust closure.

## Overclaim Checklist

See [`docs/overclaim-checklist.md`](overclaim-checklist.md) for the canonical
overclaim checklist. The summary below is authoritative only when it matches
the canonical file.

- `pr-review` emits review-record evidence over a frozen PR packet. It can
  report `coverage_satisfied`, `coverage_partial`, `coverage_unresolved`,
  `not_assessed`, or `cannot_verify`, but it does not approve, merge, mark
  ready, release, accept risk, or replace human approval.
- `gate` emits verifier-derived facts and deterministic states. It does not own
  merge, release, readiness, degradation, override approval, or risk acceptance.
- `witness` binds available CI or customer-PKI evidence. A CI witness file is
  not external production trust, a transparency log, or a release approval by
  itself.
- `release-proof` can establish `source_bound_local_release` only when the
  source commit and manifest subjects match. It does not prove
  `external_production_trust`, `trusted_contract_release`, or
  `production_release_verified`.

From verifier results, you may only state:

- Which command/profile was run.
- Which `result` or state values were produced.
- Whether the selected profile concluded with live `pass` or `observed`.
- Which states remain `not_assessed` or `cannot_verify`, with the emitted reason.

You may not state external production trust guarantees until
`external_production_trust` completes with live `pass` and
`production_release_verified` is supported by its dependent evidence chain.

## Manual External PR Review Handoff

For `manual_external` PR review planes, a usable `findings_reported` or
`no_findings` status requires retained reviewer output. A bare PR comment or
hand-edited status is not enough.

Reviewer output must be JSON matching `schema/pr-review-result.schema.json` and
must echo the packet digest, plane, and role. The review runner records the raw
output digest as `raw_output_ref`; validation counts the plane only after that
digest-bound output exists.

Minimum handoff steps:

1. Build or reuse a frozen packet directory with `packet/packet.json`.
2. Give the reviewer the packet digest, plane, role id, diff ref, context refs,
   and validation criteria.
3. Store the reviewer JSON output in a file outside the packet directory.
4. Use a profile role whose `command` prints that JSON file, then run
   `sdp-trace pr-review run --packet <packet-dir> --profile <profile> --out <runs-dir>`.
5. Run `sdp-trace pr-review synthesize`, `validate`, and `summarize` against the
   resulting run set and ledger.

If the reviewer output is absent, empty, off-task, malformed, lacks retained raw
output, or targets a different packet digest, record the plane as
`not_assessed` or `cannot_verify`. Do not treat it as sign-off.

This entrypoint is intentionally minimal and is intended to prevent over-claiming
from reproducible verifier output.
