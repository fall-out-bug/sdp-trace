# Block 25 Compiled CI Demo Pilot Report

Status: implementation evidence captured; role review and PR-level review remain `not_assessed`.

## Authority Scope

This report records a narrow demo-pilot signal for one selected compiled
Kotlin/JVM Bazel target in `fall-out-bug/sdp-trace-demo-ci-pilot`.

It does not establish production trust, release trust, owner independence,
non-GitHub portability, released-binary installation UX, broad JVM/Bazel
compatibility, or monorepo coverage. Those states remain `not_assessed`.

## Source And CI Run

| Field | Value |
| --- | --- |
| Demo repository | `fall-out-bug/sdp-trace-demo-ci-pilot` |
| Demo branch | `codex/block-25-compiled-ci-demo` |
| Demo source commit | `e7af52aca71af76635118392919d7a01fddf6c3e` |
| `sdp-trace` source ref | `codex/block-25-compiled-ci-demo` |
| `sdp-trace` source commit recorded by CI | `c28e683f8438580a89911038add11a4c976620f2` |
| Workflow | `sdp-trace demo pilot` |
| GitHub Actions run | `25554876944` |
| Run URL | `https://github.com/fall-out-bug/sdp-trace-demo-ci-pilot/actions/runs/25554876944` |
| Run event | `workflow_dispatch` |
| Run state | `success` |
| Created at | `2026-05-08T12:12:15Z` |
| Updated at | `2026-05-08T12:13:55Z` |

The workflow ran wrapped commands from `/tmp/sdp-trace-demo` so uploaded run
events did not persist the GitHub-hosted runner checkout path.

## Compiled Target Evidence

Primary behavior target:

- `//app:feature_flags_behavior_test`
- Bazel rule shape: `kt_jvm_library` plus `kt_jvm_test`
- Behavior covered: plan-specific feature enablement, audit-log entitlement
  behavior, seat-overage warning behavior, and unknown-flag denial
- JVM: Bazel remote JDK 17 through `.bazelrc`
- Kotlin rules: `rules_kotlin` 2.3.20

Local verification in the demo repository:

```text
bazel test //...
//app:entitlement_matrix_test                                   (cached) PASSED
//app:feature_flag_test                                         (cached) PASSED
//app:feature_flags_behavior_test                               (cached) PASSED
//app:audit_scope_test                                                   PASSED
```

The CI clean job wrapped the compiled target as
`clean-compiled-feature-flags`. Source/scope shell checks remained secondary
metadata checks.

## Downloaded Artifact Verification

Artifacts were downloaded locally from run `25554876944` into
`/tmp/sdp-trace-demo-run-25554876944` and verified from the downloaded bytes.

| Artifact | Artifact id | Expires at | Downloaded index entries | Local recompute state |
| --- | ---: | --- | ---: | --- |
| `sdp-trace-demo-clean-report` | `6878833487` | `2026-05-22T12:13:51Z` | 65 | `pass` |
| `sdp-trace-demo-no-oidc-report` | `6878830704` | `2026-05-22T12:13:41Z` | 13 | `pass` |

Verification commands:

```text
scripts/verify-artifact-index.sh /tmp/sdp-trace-demo-run-25554876944/sdp-trace-demo-clean-report /tmp/sdp-trace-demo-run-25554876944/sdp-trace-demo-clean-report/artifact-index.json
scripts/verify-artifact-index.sh /tmp/sdp-trace-demo-run-25554876944/sdp-trace-demo-no-oidc-report /tmp/sdp-trace-demo-run-25554876944/sdp-trace-demo-no-oidc-report/artifact-index.json
```

Both commands returned:

```text
artifact_index_verify=pass
```

The artifact index is root-level JSON with `schema_version:
demo-artifact-index-v1`, `authority_scope: demo_pilot_only`, sorted relative
paths, SHA-256 digests, file sizes, and no self-index entry.

## Redaction Scan

The CI redaction scan covered report and run directories before upload. The
downloaded artifacts record:

```text
redaction_scan=pass
match_count=0
```

This scan used the Block 24 denylist pattern file. It checks for token-like
strings, provider credentials, private key material, raw OIDC material,
authenticated artifact URLs, and known private local/runner path patterns.

## Negative Evidence Cases

| Case | Artifact evidence | Expected state | Observed state |
| --- | --- | --- | --- |
| no-OIDC witness gap | `sdp-trace-demo-no-oidc-report/report/ci-witness-no-oidc.json` | `cannot_verify` | `cannot_verify`, reason `missing_ci_oidc`, exit `3` |
| stale digest | `sdp-trace-demo-clean-report/report/dishonest/stale-digest-index.json` | `fail` | `fail`, reason codes `artifact_digest_mismatch`, `stale_index`, verifier exit `1` |
| source/run mismatch | `sdp-trace-demo-clean-report/report/dishonest/source-run-mismatch.json` | `fail` | fixture state `fail`, reason codes `source_run_binding_mismatch`, `source_commit_contradiction` |

The stale digest fixture was independent of the clean artifact-index
implementation: it generated an index over `payload.txt`, then changed
`payload.txt` without changing the index. The verifier failed only on the
payload digest mismatch.

## Trust State

| Surface | State | Evidence |
| --- | --- | --- |
| Demo CI run | `pass` | GitHub Actions run `25554876944` |
| Compiled selected target | `pass` | `//app:feature_flags_behavior_test` in local and CI wrapped runs |
| Clean CI witness | `pass` | `ci-witness.json`, established scope `ci_witnessed` |
| Artifact index verification | `pass` | Downloaded clean and no-OIDC artifacts recomputed locally |
| Redaction scan | `pass` | `redaction_scan=pass`, `match_count=0` in both artifact roots |
| No-OIDC witness case | `cannot_verify` | `missing_ci_oidc`, exit `3` |
| Stale digest case | `fail` | digest mismatch fixture, verifier exit `1` |
| Gate output | `fail` | `gate.exit` was `3`; this is not green closure evidence |
| Production trust | `not_assessed` | outside Block 25 scope |
| Broad JVM/Bazel compatibility | `not_assessed` | one selected target only |
| Non-GitHub portability | `not_assessed` | GitHub Actions-only demo |
| Released binary acquisition UX | `not_assessed` | CI built `sdp-trace` from source |
| Role review | `not_assessed` | pending CTO buyer, Head of Engineering, and Head of InfoSec review |
| PR-level review | `not_assessed` | PR not opened |

## Current Interpretation

Block 25 now has implementation evidence for the selected compiled JVM/Bazel
demo path and deterministic CI artifact indexing. It is not closed: role review,
PR-level review, and final `sdp-trace` PR verification remain open.
