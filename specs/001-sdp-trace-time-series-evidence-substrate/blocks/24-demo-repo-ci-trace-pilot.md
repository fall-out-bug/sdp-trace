# Block 24: Demo Repository CI And Trace Pilot

Status: Socratic-reviewed SpecKit delta pending CTO approval. Implementation is
blocked until the CTO explicitly approves the reviewed direction.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/23-mvp-closure-drift-and-readiness.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot-implementation-plan.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot-socratic.md`
- `docs/research/block-23-mvp-closure-package.md`
- `docs/customer-questions.en.md`
- `docs/customer-questions.ru.md`
- `docs/agent-entrypoint.md`
- `docs/reviewer-entrypoint.md`

## Goal

Prove `sdp-trace` as a CI-integrated sidecar on a small separate demo
repository through a real CI-backed workflow with inspectable trace outputs,
without relying on the retired Block 06 toy runner/validator scripts.

The block must answer the practical customer question Block 23 left open:

> Can a team integrate `sdp-trace` with a repository, run CI, and inspect trace,
> evidence, report, gate, and witness artifacts that honestly distinguish
> observed facts from `not_assessed` or `cannot_verify` gaps?

## Scope

The demo repository may be separate from this repository. `sdp-trace` remains
the portable product substrate; demo-specific application code, CI workflow
files, raw build logs, and raw model output belong in the demo repository unless
they are sanitized examples intentionally copied here.

Minimum pilot surface:

- one separate demo repository under the same GitHub owner for the first pilot,
  provisionally named `fall-out-bug/sdp-trace-demo-ci-pilot`;
- a small but real Feature Flag / Entitlements Kotlin+Bazel service with one
  deterministic Bazel test command;
- scoped Kotlin+Bazel evidence in the pilot report: assessed Bazel target,
  `BUILD` or `BUILD.bazel`, `MODULE.bazel` or workspace marker when present,
  Kotlin source or `kt_jvm_*` rule tied to the selected target, and captured
  Bazel or Bazelisk version;
- GitHub Actions as the first CI target, because the current product already
  ships a `github-actions` witness profile and the repository remote is GitHub;
- `sdp-trace` built from an explicit source ref inside the demo CI job, or a
  pinned release binary if one exists when implementation starts;
- no runtime or module dependency from this repository to the demo repository;
- `sdp-trace wrap` around the selected build/test command, with a deterministic
  run root such as `.sdp-trace-runs/`;
- committed or externally referenced `.sdp-trace-runs/` artifacts in the demo
  repository, plus a sanitized artifact index copied or linked from this repo;
- `verify`, `explain`, `report`, `gate`, and `witness` outputs over the captured
  run root;
- at least one trace path showing an observed successful command;
- at least one negative or incomplete path showing `not_assessed`,
  `cannot_verify`, missing telemetry, or witness gaps;
- a redaction boundary for logs, command output, paths, tokens, and personal
  data;
- a customer-readable pilot report mapping artifacts to the Block 23 customer
  questions.

## Proposed Decisions For Review

| decision | proposed direction | reason | residual state |
| --- | --- | --- | --- |
| Demo location | Separate repository under `fall-out-bug` | Proves first repository integration without making `sdp-trace` depend on a fixture directory inside itself. Same owner keeps access, CI, and cleanup cheap for the first pilot. | Owner-independent portability remains `not_assessed`. |
| CI provider | GitHub Actions first | The current product has a shipped `github-actions` witness path and the Block 23 remote is GitHub. This tests live CI without reopening Block 22 enterprise-provider semantics. | GitLab, Buildkite, customer PKI, and air-gapped pilot execution remain `not_assessed`. |
| Demo app | Feature Flag / Entitlements Kotlin+Bazel service with a deterministic Bazel test | Matches the repo's existing external demo target and customer-requested JVM/Kotlin/Bazel path. It is still narrow enough to keep Block 24 focused on CI trace, not application complexity. | OpenCode/GSD/model-agent execution remains `not_assessed` unless explicitly added in a later block. |
| Tool acquisition | Build `sdp-trace` from a pinned source ref in CI | Makes the tested product version explicit and replayable before release binaries exist. | Binary distribution UX remains `not_assessed` unless a release artifact is used. |
| Artifact split | Raw run/report/witness artifacts live in the demo repo or CI artifact store; this repo records sanitized report, digest index, links, and review disposition only | Prevents this repo from becoming a raw-log archive and keeps source-bound release proof separate from demo evidence. | Demo evidence is pilot evidence, not source-bound product proof. |
| Negative path | A no-OIDC or intentionally incomplete witness job records `cannot_verify`, and a local-only run records missing witness as `not_assessed` or `cannot_verify` | Shows honest gap behavior without fabricating test failures or unsafe secrets. | The exact negative job shape is finalized after review. |

## Owner Independence Gap

Same-owner demo evidence is useful only as a first integration proof. It shows
that `sdp-trace` can be built from an explicit source ref, run inside a separate
repository's CI workflow, and produce inspectable artifacts. It does not prove
that a different company, GitHub owner, access model, or CI provider can adopt
the same workflow without changes. The pilot report must include an explicit
owner-independence gap section naming what a different owner must supply:

- access to the `sdp-trace` source ref or release artifact;
- permission to add CI workflow steps;
- artifact retention policy;
- OIDC or equivalent witness permission if CI witness is in scope;
- privacy approval for any public artifact or repo visibility.

## Non-Goals

- Do not resurrect Block 06 `scripts/*` or `npm` validation as current product
  proof.
- Do not claim external production trust.
- Do not claim production readiness, policy enforcement, customer deployment
  readiness, or enterprise CI support.
- Do not put raw customer data, raw model responses, secrets, CI tokens, private
  paths, or long raw logs into this repository.
- Do not make `sdp-trace` depend on the demo repository, Beads, Operator Mode,
  a specific coding agent, or a specific CI provider.
- Do not turn demo gate output into a native `sdp-trace` policy decision.
- Do not claim GitHub Actions support beyond the exact observed demo topology.
- Do not commit raw CI logs, raw OIDC tokens, raw JWT bodies, raw command
  output streams, or private filesystem paths into this repository.
- Do not use checked-in demo artifacts to update the source-bound release
  manifest unless a separate source-bound release cycle is explicitly approved.

## Acceptance Criteria

1. The demo repository, source ref, CI workflow ref, and `sdp-trace` source ref
   are recorded in the pilot report.
2. The demo repo run can be re-executed from documented commands to produce a
   new structurally comparable run, or inspected from recorded artifact refs.
3. CI produces current evidence for at least one `sdp-trace`-wrapped command.
4. The pilot report records the selected Bazel target, Bazel/Bazelisk version,
   Kotlin target evidence, and source files/rules tied to the assessed scope.
5. `sdp-trace verify` and `explain` are run against the captured run directory
   and their outputs are linked from the pilot report.
6. `sdp-trace report` and `gate` are run against the demo run set and the pilot
   report contains a required "Gate Output Meaning" section stating that gate
   output is verifier-derived fact output, not a native merge decision, release
   gate, readiness indicator, risk acceptance, or production-trust claim.
7. `sdp-trace witness --kind github-actions` is run in CI with its exact status,
   trust scope, reason codes, and missing identity fields recorded.
8. A negative or incomplete witness path records `not_assessed` or
   `cannot_verify` with a concrete reason such as missing CI OIDC, missing
   source binding, missing run binding, or local-only evidence. The report must
   explain what a customer should conclude from that state and what evidence
   would be needed to raise the claim.
9. The pilot report maps the artifacts to the nine Block 23 customer questions
   and classifies each answer as `direct_demo_evidence`,
   `partial_demo_evidence`, or `not_assessed`.
10. The report includes a "CI Alone vs sdp-trace" section showing at least one
   fact visible in `sdp-trace` output that raw CI logs alone do not preserve as
   a structured artifact, such as missing witness identity, authority scope,
   gate fact state, artifact digest, or retained missing telemetry.
11. The report distinguishes:
   - observed command execution;
   - missing telemetry;
   - local-only evidence;
   - CI-witnessed evidence;
   - external production trust that remains `not_assessed`.
12. Redaction checks prove committed artifacts do not contain CI tokens, OIDC
   tokens, provider credentials, private key material, authenticated provider
   URLs, raw model payloads, raw logs, private filesystem paths, or unsafe
   personal identifiers.
13. Any artifact copied back into `sdp-trace` is sanitized, source-attributed,
   tagged with `authority_scope: demo_pilot_only`, and marked as demo evidence
   rather than source-bound product proof.
14. Demo repository GitHub workflow checks are recorded as `pass`, `fail`,
   `not_assessed`, or `cannot_verify`; absent checks are not called green.
15. The `sdp-trace` PR checks for Block 24 summary artifacts are recorded
   separately as `pass`, `fail`, `not_assessed`, or `cannot_verify`; absent
   checks are not called green.
16. Review runs across code/correctness, trace/evidence, and
    requirements-vs-implementation planes before implementation closure and
    again at PR level.

## Artifact Contract

The implementation plan may refine names, but the reviewed direction needs this
minimum artifact split:

| location | artifact | storage model | authority |
| --- | --- | --- | --- |
| demo repo | `.github/workflows/sdp-trace-demo.yml` | git-committed workflow | CI execution recipe for the pilot only. |
| demo repo CI | `.sdp-trace-runs/` | GitHub Actions artifact store, primary evidence during review | Observed demo run evidence; report records workflow run id, artifact name/id when available, retention days, and expiration or rerun requirement. |
| demo repo CI | `.sdp-trace-report/` with report, gate result, witness result, verify/explain outputs, exit-code files, and safety scan output | GitHub Actions artifact store, primary evidence during review | Demo pilot evidence and gap state; if artifact expires before review, the state becomes `cannot_verify` until rerun. |
| demo repo | sanitized report/index copies, if public-safe | git-committed sanitized evidence only | Durable demo summary, not raw authority for original CI artifact bytes. |
| `sdp-trace` repo | `docs/research/block-24-demo-repo-ci-trace-pilot-report.md` | git-committed sanitized summary | Customer-readable summary, links, states, owner-independence gap, and Block 23 question mapping. |
| `sdp-trace` repo | `docs/research/block-24-demo-repo-ci-artifact-index.md` | git-committed sanitized index | Sanitized artifact references and digests, not raw logs. |
| `sdp-trace` repo | `specs/.../blocks/24-demo-repo-ci-trace-pilot-review-ledger.md` | git-committed ledger | Review findings, dispositions, re-review state, CI state, and residual gaps. |

## Evidence Classification

Block 24 artifacts must carry three separate classifications. Do not collapse
them into one status string.

| axis | allowed values | meaning |
| --- | --- | --- |
| `capture_state` | `captured`, `not_captured`, `failed` | Whether `sdp-trace` captured or generated the selected artifact. |
| `attestation_state` | `local_only`, `ci_witnessed`, `ci_witness_attempted_incomplete`, `not_assessed`, `cannot_verify`, `fail` | Who or what can vouch for the artifact. |
| `authority_scope` | `demo_pilot_only`, `source_bound_release`, `external_production_trust_not_assessed` | The trust boundary the artifact may support. |

Rules:

- `demo_pilot_only` is mandatory for every Block 24 demo artifact copied into
  this repository.
- `source_bound_release` is forbidden for Block 24 demo artifacts unless a
  separate source-bound release cycle is approved and completed.
- `external_production_trust_not_assessed` is mandatory for the external
  production trust row in the pilot report.
- `not_assessed` applies when a provider/profile/trust state is explicitly out
  of scope for the selected pilot.
- `cannot_verify` applies when a selected in-scope check ran or was requested
  but lacked required identity, OIDC, binding, freshness, artifact, environment,
  or structural evidence.
- `ci_witness_attempted_incomplete` applies when CI witness collection ran and
  produced a structured non-pass result without a hard mismatch.

Block 24 may show `attestation_state: ci_witnessed` for the exact demo run if
the shipped verifier emits passing CI-bound witness facts. It must keep external
production trust as `external_production_trust_not_assessed`.

## Witness Result Fields

The pilot report and artifact index must extract only these fields from witness
JSON unless the review ledger approves a narrower exception:

| field | requirement |
| --- | --- |
| `kind` | Must be `github-actions` for the positive Block 24 CI path. |
| `status` | Must be one of `pass`, `fail`, `cannot_verify`, or `not_assessed`. |
| `trust_scope` / `established_trust_scope` | Must be recorded without upgrading beyond the verifier output. |
| `reason` / `reason_codes` | Required for non-pass states and recommended for all states. |
| `generated_at` | Required for freshness context. |
| `missing_identity_fields` | Required when identity or OIDC evidence is missing. |
| `run_artifacts` / `report_artifacts` | Record paths and SHA-256 digests only. |
| `profile_states` | Required when present in verifier output. |
| `output_safety` | Required when present in verifier output. |

Raw OIDC JWTs, CI request tokens, authenticated URLs, actor emails, and private
paths must not be copied into this repository.

## Pilot Report Contract

The report must include these sections:

- demo repo, owner, visibility, source refs, `sdp-trace` source ref, workflow
  run id, artifact retention, and owner-independence gap;
- selected Bazel target, Bazel/Bazelisk version, Kotlin source/rule evidence,
  and target-scoped build ownership evidence;
- "CI Alone vs sdp-trace" contrast;
- "Gate Output Meaning" disclaimer;
- evidence classification table with `capture_state`, `attestation_state`, and
  `authority_scope`;
- negative-path customer interpretation;
- artifact index summary;
- redaction scan command, pattern digest, result, and residual safety state;
- nine Block 23 customer-question rows.

The nine question rows must use this quality bar:

| answer class | allowed meaning |
| --- | --- |
| `direct_demo_evidence` | The exact Block 24 demo artifacts answer the question for the selected demo run. |
| `partial_demo_evidence` | The demo artifacts answer part of the question and name the missing evidence needed to raise the claim. |
| `not_assessed` | The selected demo does not attempt the question's provider/profile/trust state. |

## Safety Scan Contract

Block 24 implementation must use an executable denylist scan before any demo
artifact is committed or linked as customer-inspectable evidence. The planned
pattern file is `docs/research/block-24-redaction-denylist.patterns`.

Minimum command shape:

```bash
rg --pcre2 -n -f docs/research/block-24-redaction-denylist.patterns <artifact-root> <sanitized-copy-root>
```

Expected result: no matches and exit code `1` from `rg`. Any match blocks
closure until the artifact is redacted or the review ledger records a justified
false positive. If the scan is skipped, absent from CI output, or cannot run,
artifact safety is `cannot_verify`, not green.

The implementation must record:

- exact scan command;
- SHA-256 digest of the pattern file;
- artifact roots scanned;
- exit state;
- match count;
- disposition for every match or false positive.

Copied JSON excerpts must use an allow-list of top-level fields from the
Witness Result Fields section, stay below 40 lines or 2 KB per artifact, and be
tagged as `sanitized demo evidence; authority_scope=demo_pilot_only`.

## Open Design Questions For Socratic Review

1. Is same-owner separate repo enough for the first portability proof, or does
   it under-prove customer-style attachment?
2. Is GitHub Actions the right first target, or does it overfit the product demo
   to the current repository host?
3. Is the Feature Flag / Entitlements Kotlin+Bazel service still small enough
   for a CI trace pilot, or does it pull Block 24 back into broad stack/demo
   work?
4. Which artifacts must stay only in the demo repo or CI artifact store, and
   which sanitized summaries should be copied into `sdp-trace`?
5. What is the smallest negative scenario that proves missing telemetry or
   witness gaps without fabricating a product failure?
6. Does building `sdp-trace` from source in demo CI prove attachment well enough,
   or does the pilot need a release binary/install path?
7. How should the pilot report avoid turning `gate` and `witness` facts into
   policy, support, readiness, or production-trust claims?
8. What exact state keeps non-GitHub CI providers honest: `not_assessed` because
   out of scope, or `cannot_verify` because the current pilot cannot run them?

## Verification Plan

Minimum local checks in `sdp-trace` during spec intake:

```bash
go test ./...
go vet ./...
staticcheck ./...
golangci-lint run ./...
jq empty schema/*.json
git diff --check HEAD
go run ./cmd/sdp-trace --help
```

Minimum demo checks, to be finalized after Socratic review and run from the demo
repository with `SDP_TRACE_BIN` pointing at the pinned `sdp-trace` binary:

```bash
mkdir -p .sdp-trace-report
"$SDP_TRACE_BIN" wrap --name demo-test --output-dir .sdp-trace-runs -- <demo-command>
"$SDP_TRACE_BIN" verify <demo-run-dir>
"$SDP_TRACE_BIN" explain <demo-run-dir>
"$SDP_TRACE_BIN" report --out .sdp-trace-report .sdp-trace-runs
"$SDP_TRACE_BIN" gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs
"$SDP_TRACE_BIN" witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
```

Every command must record its actual environment, source refs, artifact refs,
exit state, and residual `not_assessed`/`cannot_verify` states.

## Approval Boundary

Approval of this reviewed SpecKit direction authorizes implementation planning
and demo-repository work for the exact scoped pilot above. It does not authorize
external production trust claims, enterprise CI support claims, policy
enforcement claims, or broad customer deployment readiness claims.
