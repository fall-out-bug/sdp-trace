# Block 24 Implementation Plan: Demo Repository CI And Trace Pilot

Status: Socratic-reviewed plan approved for implementation on 2026-05-08.

## Product Thesis

Block 24 is not another internal fixture. It is the first small externalized
repository integration proof: a team can point `sdp-trace` at a repository, run
CI, and inspect trace, evidence, report, gate, and witness artifacts while gaps
stay explicit.

The demo must be boring enough to replay and real enough to answer customer
pressure questions. It should test the current product surfaces, not invent new
ones in the demo.

## Proposed Workstreams

| workstream | write scope | expected output | closure test |
| --- | --- | --- | --- |
| WS0 spec review | Block 24 spec, Socratic file, review ledger | reviewed SpecKit direction with dispositions | no critical/major Socratic findings remain unresolved |
| WS1 demo repo selection | external demo repo docs/workflow, sanitized pointer in this repo | repository URL, ownership, visibility decision, CI provider, source refs, privacy boundary | repo can be cloned or inspected from documented refs or access limits are recorded |
| WS2 demo app and CI | external demo repo only | Feature Flag / Entitlements Kotlin+Bazel service, deterministic Bazel test, GitHub Actions workflow | CI runs the selected Bazel test command |
| WS3 trace/report/gate/witness capture | external demo repo and CI artifacts | `.sdp-trace-runs/`, `.sdp-trace-report/`, gate result, witness result, verify/explain outputs | current CI run emits artifacts with exact states |
| WS4 negative and safety path | external demo repo workflow plus sanitized safety summary | two intentionally dishonest-trace cases, no-OIDC or incomplete witness output, redaction scan, public/private decision record | missing witness/OIDC, stale digest, source/run mismatch, or overclaim states are `cannot_verify`/`fail`/`not_assessed`, not hidden |
| WS5 sdp-trace summary artifacts | `docs/research/block-24-*`, Block 24 review ledger | customer-readable report, artifact index, review disposition | report maps all nine Block 23 questions |
| WS6 implementation and PR review | changed files only | code/correctness, trace/evidence, requirements-vs-implementation review evidence | review findings fixed or dispositioned; PR-level review repeated |

WS1 through WS4 may happen in the external demo repository after WS0 approval.
This repository should receive only sanitized summaries, digest indexes, review
ledgers, and docs needed to make the pilot inspectable.

## Demo Repository Shape

Proposed first repo:

- owner: `fall-out-bug`
- name: `sdp-trace-demo-ci-pilot`
- visibility: public if safe; private is acceptable only if the pilot report
  records access limitations
- app: Feature Flag / Entitlements Kotlin+Bazel service with one narrow domain
  path, for example "user has plan X, flag Y resolves to enabled/disabled"
- test: deterministic `bazel test //...` or one explicitly named Bazel target
- stack evidence: captured `bazel version` or `bazelisk version`, selected target
  label, target-scoped `BUILD` or `BUILD.bazel`, Kotlin source or `kt_jvm_*`
  rule tied to the target, and module/workspace marker when present
- CI: GitHub Actions workflow with one normal witness job and one intentionally
  incomplete witness path if needed for negative evidence

The demo repo must not import this repository as an application dependency. The
CI job may clone or check out `fall-out-bug/sdp-trace` at an explicit ref and
build the CLI as a tool.

## Command Sketch

In the positive GitHub Actions witness job:

```bash
mkdir -p .sdp-trace-tools .sdp-trace-report
go build -o .sdp-trace-tools/sdp-trace ./.sdp-trace-src/cmd/sdp-trace
export SDP_TRACE_BIN="$PWD/.sdp-trace-tools/sdp-trace"

"$SDP_TRACE_BIN" wrap \
  --name demo-test \
  --output-dir .sdp-trace-runs \
  -- bazel test //...

RUN_DIR="$(find .sdp-trace-runs -mindepth 1 -maxdepth 1 -type d | sort | tail -1)"

set +e
"$SDP_TRACE_BIN" verify "$RUN_DIR" > .sdp-trace-report/verify.txt
printf "%s\n" "$?" > .sdp-trace-report/verify.exit
"$SDP_TRACE_BIN" explain "$RUN_DIR" > .sdp-trace-report/explain.txt
printf "%s\n" "$?" > .sdp-trace-report/explain.exit
"$SDP_TRACE_BIN" report --out .sdp-trace-report .sdp-trace-runs
printf "%s\n" "$?" > .sdp-trace-report/report.exit
"$SDP_TRACE_BIN" gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs
printf "%s\n" "$?" > .sdp-trace-report/gate.exit
"$SDP_TRACE_BIN" witness \
  --kind github-actions \
  --out .sdp-trace-report/ci-witness.json \
  --report-dir .sdp-trace-report \
  .sdp-trace-runs
printf "%s\n" "$?" > .sdp-trace-report/ci-witness.exit
set -e
```

The actual workflow must materialize the pinned `sdp-trace` source ref before
the build step, for example:

```yaml
permissions:
  contents: read

jobs:
  trace-with-oidc:
    permissions:
      contents: read
      id-token: write
    steps:
      - uses: actions/checkout@v4
        with:
          path: demo
      - uses: actions/checkout@v4
        with:
          repository: fall-out-bug/sdp-trace
          ref: <sdp-trace-source-sha>
          path: demo/.sdp-trace-src
```

The pilot report must record `<sdp-trace-source-sha>`, demo repo commit SHA,
workflow run id, job id/name, run attempt, artifact retention period, selected
Bazel target, Bazel/Bazelisk version, Kotlin source/rule evidence, and
target-scoped build ownership evidence.

If the final workflow uses `sdp-trace run --task`, it must use
`--use-default-contract` or an explicit contract file and record why `run` is
better than `wrap` for the customer story.

## Negative Path

The implementation must include two intentionally dishonest-trace cases plus
one explicit witness gap path:

- clean baseline: three honest cases capture successful wrapped commands with
  distinct names and scopes;
- dishonest case 1: a copied or edited trace/report artifact intentionally
  mismatches the observed run/source binding and must not be represented as a
  clean CI-witnessed run;
- dishonest case 2: a stale or tampered digest/index entry intentionally fails
  the artifact-integrity or source/run-binding story;
- witness gap: a separate GitHub Actions job without OIDC permission writes
  `ci-witness-no-oidc.json` with `status: "cannot_verify"` and reason
  `missing_ci_oidc`;
- fallback: a local-only run or incomplete witness envelope records
  `not_assessed` or `cannot_verify` with exact missing fields.

The negative paths must not fabricate app failure, leak a token, or make the
overall pilot look broken. They exist to prove dishonest or incomplete trace
rendering is explicit and does not get upgraded into trust.

Negative job YAML must keep OIDC unavailable at job scope:

```yaml
jobs:
  trace-without-oidc:
    permissions:
      contents: read
      id-token: none
```

The implementation must verify that `ACTIONS_ID_TOKEN_REQUEST_URL` and
`ACTIONS_ID_TOKEN_REQUEST_TOKEN` are absent in the negative job before running
`witness`. The current CLI contract writes a witness JSON record and exits `3`
for `cannot_verify`; the workflow must capture that exit code without losing
the artifact.

## Sanitized Artifacts In This Repository

Expected files after implementation:

- `docs/research/block-24-demo-repo-ci-trace-pilot-report.md`
- `docs/research/block-24-demo-repo-ci-artifact-index.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot-review-ledger.md`

These files may link to demo repo commits, workflow runs, and CI artifacts. If
copied JSON snippets are needed, they must be shortened, sanitized, and clearly
labeled as demo evidence. Raw logs, raw JWTs, provider tokens, private paths,
and raw model payloads stay out of this repository.

## Safety Checks

Before any Block 24 implementation closure, run a redaction scan over every
artifact copied into this repository and over every demo artifact intended for
public inspection. The scan must fail or block closure if it finds:

- CI request tokens, provider credentials, or bearer tokens;
- raw JWT bodies;
- private key material;
- authenticated provider URLs;
- raw command logs beyond short sanitized excerpts;
- raw model request/response payloads;
- private filesystem paths;
- customer data or unsafe personal identifiers.

Negative assertions must avoid echoing real secrets in failure output.

Required command shape:

```bash
rg --pcre2 -n -f docs/research/block-24-redaction-denylist.patterns <artifact-root> <sanitized-copy-root>
```

Expected result is no matches and `rg` exit code `1`. If the command cannot run,
the safety state is `cannot_verify`. If it finds matches, closure is blocked
until every match is redacted or recorded as a false positive in the review
ledger.

Every copied JSON excerpt must be field-allowlisted, less than 40 lines or 2 KB
per artifact, and tagged as sanitized demo evidence with
`authority_scope=demo_pilot_only`.

The public/private demo repo decision must be recorded before repo creation or
publication. The decision record must list workflow filenames, runner OS labels,
repo/org names in OIDC subject or source fields, artifact URLs, commit
timestamps, and whether each class is safe to expose publicly. If any class is
unsafe, the demo repo stays private and the pilot report records the resulting
inspectability limit.

## Review And Verification

Spec gate:

```bash
go test ./...
go vet ./...
staticcheck ./...
golangci-lint run ./...
jq empty schema/*.json
git diff --check HEAD
go run ./cmd/sdp-trace --help
```

Implementation gate:

- rerun the spec gate in `sdp-trace`;
- verify the demo repo CI run state from GitHub, with absent checks recorded as
  `not_assessed`;
- inspect the demo repo artifacts from the exact workflow run;
- run the redaction denylist scan and record pattern-file digest and exit state;
- run code/correctness, trace/evidence, and requirements-vs-implementation
  review planes on the implementation;
- repeat those review planes at PR level for `sdp-trace`;
- record every finding with `accepted_fixed`, `accepted_narrower`,
  `false_positive`, `deferred_not_assessed`, or `unresolved_blocker`.

## Drift And Regression Checks

Block 24 implementation must check drift against:

- Block 06 retired-runner language: no current closure may cite retired
  `scripts/*` or `npm` validation as proof;
- Kotlin+Bazel evidence: the demo must cite target-scoped Bazel and Kotlin
  evidence, not `.bazelrc`, Maven/Gradle metadata, or synthetic fixtures alone;
- Block 22 witness semantics: GitHub Actions witness facts must not overclaim
  enterprise/provider-neutral support;
- Block 23 customer questions: all nine rows must be answered or explicitly
  marked `not_assessed`/`cannot_verify`;
- source-bound proof boundary: demo evidence must not be called
  `source_bound_local_release` proof.

## Approval Record

The CTO approved implementation on 2026-05-08 after Socratic review and fixes,
with the clarified case shape:

1. separate same-owner demo repo;
2. GitHub Actions first;
3. Feature Flag / Entitlements Kotlin+Bazel demo service;
4. source-built `sdp-trace` tool in CI;
5. sanitized report and artifact-index copy-back only;
6. three clean trace cases;
7. two intentionally dishonest-trace cases;
8. no-OIDC or incomplete witness negative path where needed for witness gap
   semantics.
