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

--- END REVIEWER-ENTRYPOINT ---

--- BEGIN AGENT-ENTRYPOINT (state contract) ---
## State And Exit Code Contract

### Result States

These are the verifier result states. Every command that reports verifier
outcome uses one of these. They map to exit codes.

| Result state | Exit code | Meaning |
| --- | --- | --- |
| `observed` | `0` | Verifier evidence met required checks for the selected local profile. |
| `pass` | `0` | Selected profile concluded successfully where the command contract uses pass/fail states. |
| `fail` | `1` | Verifier evidence conflicted or was insufficient for required checks. |
| `not_assessed` | `0` | State was outside the run scope; it does not imply success or evidence. May return `0` only when the command completed and the unassessed state is explicitly outside the selected profile or run scope. |
| `cannot_verify` | `3` | Verifier could not execute a required check or lacked required evidence. |

Exit code `2` is reserved for usage error / invalid command invocation and is
not a verifier result state.

### Telemetry Labels

These labels describe evidence availability, not verifier outcomes. They do
not have exit-code mappings.

- `missing_telemetry`: a telemetry stream or metric was expected but not found.
  Used by query-pack and managed-harness capture-depth reporting. It is not a
  verifier result state; the corresponding verifier result may be `not_assessed`
  or `cannot_verify` depending on whether the telemetry was required.

### Completeness Markers

These describe source or input completeness, not verifier outcomes.

- `complete`, `partial`: source completeness for `interaction import-transcript`.
  They describe the imported data set, not a verifier pass/fail.

### PR-Review Sub-States

These are command-specific coverage states reported by `pr-review`. They are
not verifier result states and do not have exit-code mappings.

- `coverage_satisfied`: the review packet reached the coverage threshold for the profile.
- `coverage_partial`: some planes were reviewed but coverage did not reach the threshold.
- `coverage_unresolved`: coverage could not be determined.

### External Verdict Sub-States

These appear in policy-consumer output and concept docs. They are not verifier
result states.

- `warn`: evidence exists but risk remains. Defined in `docs/concepts.md` as an
  External Verdict value. The corresponding verifier result is typically
  `observed` or `pass` with an advisory note.

### Integration And Adapter Labels

These describe integration or adapter status. They are not verifier result states.

- `not_integrated`: an expected adapter or integration is absent.
- `unsupported`: a format, schema version, or configuration is not supported.

### Authority Scope Labels

These describe the reporting boundary, not a verifier result.

- `outside_authority`: an observed action is outside the caller-selected
  authority envelope. Emitted by `assess --profile authority-envelope`.
- `local_dirty_structural_only`: dirty-checkout structural output. Valid only
  for local shape/debug inspection; cannot support source-bound or external
  trust closure.

A checked-in proof JSON is an audit artifact, not authority. Authority is replayed
only from live Go verifier output and the canonical command/state contract above.

## Air-Gapped Fixture Guidance

Air-gapped evidence is a fixture and customer-policy pattern, not a native
`witness --kind air-gapped` command. Use `customer-pki` or an accepted private
equivalent with explicit authority policy, payload digest, freshness evidence,
and retained audit references. If those are absent, record `not_assessed` or
`cannot_verify`; do not claim external production trust.


--- END AGENT-ENTRYPOINT ---

--- BEGIN OUTPUT-LOCATION-MAP ---
# Output Location Map

This table maps each command family to its default output location, format,
and trust boundary. For the full command contract, see
`docs/agent-entrypoint.md`.

## Run Artifacts

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `wrap` | `--output-dir .sdp-trace-runs/<name>/` | JSON + metadata | Record one command as a trace run | Local observation only |
| `run` | `--output-dir .sdp-trace-runs/<task-ref>/` | JSON + metadata | Task-linked trace run | Local observation; missing contract evidence visible |
| `observe setup` | `--out <run-dir>/` | JSON | Setup metadata before harness run | Session-profile bounded |
| `observe collect` | normalizes into `<run-dir>/` | JSON | Harness output after run | `cannot_verify` if declared output missing |
| `observe session` | `--out <run-dir>/` | JSON | Convenience wrapper for setup + collect | Same as setup + collect |
| `harness observe` | `--out <run-dir>/` | JSON | Import local harness export | Reads explicit files only; unsafe content fails before write |

## Reports And Summaries

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `report` | `--out .sdp-trace-report/` | JSON + markdown | Package observed data and gaps | Report presence is not proof of completeness |
| `gate` | `--out .sdp-trace-report/gate-result.json` | JSON | Advisory gate facts | Not a native merge/release/risk decision |
| `explain` | stdout | Markdown | Human-readable run explanation | Does not upgrade trust scope |
| `harness summarize` | stdout | Markdown | Human summary of harness validation | Non-authoritative |
| `assess explain` | stdout | Markdown | Explain assessment result | Unsupported schema may give `cannot_verify` |
| `query-pack explain` | stdout | Markdown | Explain forensic query-pack | No new evidence created |
| `envelope summarize` | `--out summary.json` | JSON | Summarize delivery trace envelope | Read-only over refs |
| `interaction summarize` | `--out summary.json` | JSON | Summarize interaction events | Friction counts are facts, not scores |

## Query And Assessment

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `query` | stdout | JSON | Missing evidence or capture depth | Missing rows are not passes |
| `query-pack` | `--out <file>` | JSON | Forensic query package | Limited by retained/redacted evidence |
| `assess --profile <profile>` | `--out <file>` | JSON | Profile-specific assessment | Facts only; policy owns block/allow |
| `assess preview` | stdout | JSON | Preview required inputs | Read-only; does not evaluate authority |

## Witness And Release

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `witness --kind <kind>` | `--out <file>` | JSON | CI or customer witness artifact | CI witness is not production trust by itself |
| `release-proof` | `--out <file>` | JSON | Source-bound local release proof | Narrower than external production trust |

## Cross-Repo And Telemetry

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `export cross-repo-posture` | `--out <file>` | JSON | Cross-repo evidence posture | Degradation decisions remain outside |
| `export telemetry` | `--out <file\|->` | Prometheus text | Telemetry export | Dashboards/alerts remain downstream |

## PR Review

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `pr-review packet` | `--out <dir>/` | JSON + files | Build frozen PR packet | Packet digest binds to inputs |
| `pr-review run` | `--out <dir>/` | JSON | Run review planes | Raw output digest recorded as `raw_output_ref` |
| `pr-review synthesize` | `--out <file>` | JSON | Synthesize runs | Aggregation only |
| `pr-review validate` | `--out <file>` | JSON | Validate against ledger | Completeness check |
| `pr-review summarize` | `--out <file>` | JSON | Summarize validation | Not merge approval |
| `pr-review check` | `--out <dir>/` | JSON + files | End-to-end PR review | Review-record completeness only |

## Packet Commands

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `packet build-pr` | `--out <dir>/` | JSON + files | Build live PR packet | `PC-VERIFICATION` must bind to workflow evidence |
| `packet build-github` | `--out <file>` | JSON | Build from fixture | Backfill/fixture authority only |
| `packet validate` | stdout | JSON | Validate bundle | Structural validation |
| `packet check-demo` | stdout | JSON | Demo-check bundle | Limited to first-packet minimum bar |
| `packet render` | `--out <file>` | Markdown | Render bundle | Row states and residual gaps are not approval |

--- END OUTPUT-LOCATION-MAP ---

--- BEGIN PROFILE-SELECTION-GUIDE ---
# Profile Selection Guide

This guide maps the three profile taxonomies used in `sdp-trace` and helps you
choose the right one. For the canonical state and exit-code contract, see
`docs/agent-entrypoint.md`.

## Taxonomy Overview

| Taxonomy | What it describes | Examples |
| --- | --- | --- |
| **Trust profile ID** | What level of trust the evidence can support | `repo_baseline_structural`, `source_bound_local_release`, `external_production_trust` |
| **Assessment profile** | Which kind of evidence quality check to run | `adapter-capture`, `managed-harness`, `forensic-retention`, `ci-artifact-observation`, `authority-envelope` |
| **Witness kind** | Which CI or customer system provided identity evidence | `github-actions`, `gitlab-ci`, `buildkite`, `customer-pki` |
| **Authority scope** | The reporting boundary for a committed package | `demo_pilot_only`, `local_dirty_structural_only` |

Trust profile IDs and authority scopes are **not** commands. Assessment profiles
are selected with `sdp-trace assess --profile <profile>`. Witness kinds are
selected with `sdp-trace witness --kind <kind>`.

## Trust Profile IDs

Choose the trust profile from the claim you need to make, not from your role.

| Trust profile ID | Use when | What it proves |
| --- | --- | --- |
| `repo_baseline_structural` | You need structural command, fixture, and local trace integrity. | Local shape and debug inspection. |
| `source_bound_local_release` | You need local manifest, source commit, artifact digest, and DSSE/source-bound checks. | The built artifact matches the source commit and manifest. |
| `external_production_trust` | You need external identity, protected source, transparency or customer audit evidence, approval, and production release verification. | The full external trust chain is closed. |

**Rule**: Do not use a lower trust profile to claim a higher one. Dirty-checkout
output is valid only under `local_dirty_structural_only` (an authority scope,
not a profile ID) and cannot close `source_bound_local_release` or
`external_production_trust`.

## Assessment Profiles

Choose the assessment profile from the evidence question you need answered.

| Question | Assessment profile | Typical command |
| --- | --- | --- |
| Did the adapter capture enough evidence? Is there overclaim risk? | `adapter-capture` | `sdp-trace assess --profile adapter-capture --out assessment.json --run <run-dir>` |
| Does the managed harness evidence satisfy policy, registry, and witness inputs? | `managed-harness` | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run <run-dir> --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` |
| Can retained evidence support forensic reconstruction? | `forensic-retention` | `sdp-trace assess --profile forensic-retention --out assessment.json --run <run-dir> --redaction-policy redaction.json` |
| Are selected artifact families CI-uploaded facts or lower-authority facts? | `ci-artifact-observation` | `sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` |
| Do observed actions stay inside a caller-selected authority envelope? | `authority-envelope` | `sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` |

Assessment profiles produce **facts only**. Block/allow, readiness, and policy
decisions remain with the downstream consumer.

## Witness Kinds

Choose the witness kind from the CI or customer system that produced the run.

| System | Witness kind | Typical command |
| --- | --- | --- |
| GitHub Actions with OIDC | `github-actions` | `sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` |
| GitLab CI | `gitlab-ci` | `sdp-trace witness --kind gitlab-ci --out gitlab-witness.json --witness-envelope envelope.json .sdp-trace-runs` |
| Buildkite | `buildkite` | `sdp-trace witness --kind buildkite --out buildkite-witness.json --witness-envelope envelope.json .sdp-trace-runs` |
| Customer PKI or private-equivalent | `customer-pki` | `sdp-trace witness --kind customer-pki --out customer-pki-witness.json --customer-pki-authority-policy policy.json --customer-pki-public-cert cert.pem --customer-pki-payload-digest <sha256> --customer-pki-freshness-evidence freshness.json .sdp-trace-runs` |

A CI witness file is **not** external production trust by itself. It binds
available evidence to CI identity when the required OIDC or envelope evidence
exists.

## Authority Scopes

Authority scopes describe the boundary of a report or package, not a verifier
result.

| Authority scope | Meaning |
| --- | --- |
| `demo_pilot_only` | Sanitized demo-repository evidence. Supports pilot mechanics only. |
| `local_dirty_structural_only` | Dirty-checkout structural output. Local shape/debug inspection only. |

## Decision Flow

1. **What do you need to prove?** → Choose a **trust profile ID**.
2. **What evidence do you have?** → Choose an **assessment profile** to check it.
3. **Where did the run happen?** → Choose a **witness kind** if CI identity is available.
4. **Is the checkout clean?** → If dirty, the authority scope is `local_dirty_structural_only`.

Do not mix scopes: a `repo_baseline_structural` result plus a `github-actions`
witness does not become `external_production_trust` unless every required
external-trust check passes live.

--- END PROFILE-SELECTION-GUIDE ---

--- BEGIN OVERCLAIM-CHECKLIST ---
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

--- END OVERCLAIM-CHECKLIST ---

--- BEGIN CONCEPTS (external verdict) ---
## External Verdict

The externally produced gate or policy outcome:

- `pass`: evidence satisfies the gate
- `warn`: evidence exists but risk remains
- `fail`: evidence proves the gate is not satisfied
- `not_assessed`: the gate was outside the selected scope or intentionally not
  evaluated in this run
- `cannot_verify`: the selected gate could not be verified because required
  evidence, environment, freshness, or consistency was missing

Missing required evidence for a selected gate is `cannot_verify` or `fail`, not
`pass`. Missing optional or out-of-scope evidence is `not_assessed`.

`warn` is an External Verdict sub-state, not a verifier result state. The
verifier result states (`observed`, `pass`, `fail`, `not_assessed`,
`cannot_verify`) and their exit-code mappings are defined in
`docs/agent-entrypoint.md`. When an External Verdict is `warn`, the underlying
verifier result is typically `observed` or `pass` with an advisory note.


--- END CONCEPTS ---


## Contract (rules the artifact must satisfy)

- A cold user must be able to choose the next command from a task-oriented guide rather than reading a long flat reference.
- State vocabulary must be consistent: result states (observed, pass, fail, not_assessed, cannot_verify) with exit codes; all other tokens classified as telemetry labels, completeness markers, PR-review sub-states, external verdict sub-states, integration labels, or authority scope labels.
- No doc outside the canonical contract redefines or invents result states.
- One canonical overclaim checklist exists; other docs link to it.
- An output location map exists mapping command family → default output path → format → purpose → trust boundary.
- A profile decision tree exists mapping trust profile IDs ↔ assessment profiles ↔ witness kinds ↔ authority scopes.

## Review Prompt (apply this lens)

You are a cold user who has never used sdp-trace before. You just checked out the repo and opened the reviewer entrypoint.

1. Can you find the next command for your task in under 30 seconds?
2. Can you tell what exit code means what without reading the entire agent entrypoint?
3. Can you find where outputs are written without scanning every command table?
4. Can you tell which profile to use without guessing from three disjoint lists?
5. Can you find the overclaim checklist without reading 4 different docs?
6. Are there any remaining ambiguous state terms that look like result states but aren't defined in the canonical contract?
7. Are there any docs that still duplicate the overclaim checklist in a way that could drift?

Return only actionable issues with file/line or artifact references, or state that you cannot find any after checking the contract. Do not validate. Do not summarize.
