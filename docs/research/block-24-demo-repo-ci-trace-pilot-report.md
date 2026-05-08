# Block 24 Demo Repository CI Trace Pilot Report

Authority scope: `demo_pilot_only`

## Executive Summary

Block 24 produced a first externalized repository pilot in
`fall-out-bug/sdp-trace-demo-ci-pilot`. GitHub Actions run `25548285336`
checked out the demo repo, checked out `sdp-trace` from
`codex/block-24-demo-ci-trace-pilot`, built the CLI from source, ran three
Bazel-backed clean trace cases, produced report/gate/witness outputs, produced
two intentionally dishonest trace cases, and uploaded CI artifacts.

This is real demo-repository evidence. It is not production trust, not release
proof, not customer deployment readiness, and not proof of non-GitHub
portability.

## Repository And CI

| field | value |
| --- | --- |
| Demo repo | `fall-out-bug/sdp-trace-demo-ci-pilot` |
| Visibility | private |
| Demo commit | `e370d1c00df8a7e7859adc284480563a269e64ca` |
| Workflow | `.github/workflows/sdp-trace-demo.yml` |
| Workflow run | `25548285336` |
| Workflow state | `pass` |
| Primary artifact | `sdp-trace-demo-clean-report` / `6876152707` |
| Negative artifact | `sdp-trace-demo-no-oidc-report` / `6876147448` |
| Artifact expiration | 2026-05-22 |
| `sdp-trace` source ref | `codex/block-24-demo-ci-trace-pilot` |
| `sdp-trace` source commit at run time | `f66aa1c4619f0f6a2d56f602da9e0135b00e4a84` |

Owner-independence remains `not_assessed`: another owner would still need
access to the `sdp-trace` source or a release artifact, permission to edit CI,
artifact-retention policy, OIDC or equivalent witness permission, and privacy
approval for public artifact exposure.

## Demo App And Bazel Scope

The demo app is a small Feature Flag / Entitlements Kotlin service surface.

| evidence | value |
| --- | --- |
| Kotlin source | `app/src/main/kotlin/demo/FeatureFlags.kt` |
| Bazel package | `app/BUILD.bazel` |
| Module marker | `MODULE.bazel` |
| Clean targets | `//app:feature_flag_test`, `//app:entitlement_matrix_test`, `//app:audit_scope_test` |
| Bazelisk version | `v1.28.1` |
| Bazel build label | `9.1.0` |

The Bazel tests are shell-based scope checks over Kotlin source and repo
metadata. This is enough for Block 24 CI trace mechanics, but it is not
evidence of compiled Kotlin/JVM compatibility.

## Clean Cases

| case | command | capture state | attestation state | authority scope |
| --- | --- | --- | --- | --- |
| `clean-feature-flag` | `bazel test //app:feature_flag_test` | `captured` | `ci_witnessed` for witness output; run summary remains `local_observed` | `demo_pilot_only` |
| `clean-entitlement-matrix` | `bazel test //app:entitlement_matrix_test` | `captured` | `ci_witnessed` for witness output; run summary remains `local_observed` | `demo_pilot_only` |
| `clean-audit-scope` | `bazel test //app:audit_scope_test` | `captured` | `ci_witnessed` for witness output; run summary remains `local_observed` | `demo_pilot_only` |

`summary.json` recorded `run_count=3`, `observed_count=3`, and zero failed,
`cannot_verify`, or `not_assessed` run summaries. The GitHub Actions witness
profile separately recorded `status=pass`, `established_trust_scope=ci_witnessed`,
and `reason=ci_identity_present`.

## Evidence Classification

| artifact or case | capture state | attestation state | authority scope |
| --- | --- | --- | --- |
| `clean-feature-flag` run directory | `captured` | `ci_witnessed` through `ci-witness.json`; run summary itself is `local_observed` | `demo_pilot_only` |
| `clean-entitlement-matrix` run directory | `captured` | `ci_witnessed` through `ci-witness.json`; run summary itself is `local_observed` | `demo_pilot_only` |
| `clean-audit-scope` run directory | `captured` | `ci_witnessed` through `ci-witness.json`; run summary itself is `local_observed` | `demo_pilot_only` |
| `verify.txt` | `captured` | `ci_witnessed` as a report artifact digest in `ci-witness.json` | `demo_pilot_only` |
| `explain.txt` | `captured` | `ci_witnessed` as a report artifact digest in `ci-witness.json` | `demo_pilot_only` |
| `summary.json` | `captured` | `local_observed` in the run summary; `ci_witnessed` only as a report artifact digest cross-referenced by `ci-witness.json` | `demo_pilot_only` |
| `gate-result.json` | `captured` | `ci_witnessed` as a report artifact digest in `ci-witness.json`; gate facts still include CI/audit `cannot_verify` | `demo_pilot_only` |
| `ci-witness.json` | `captured` | `ci_witnessed` | `demo_pilot_only` |
| `ci-witness-no-oidc.json` | `captured` | `cannot_verify` | `demo_pilot_only` |
| `dishonest-source-run-mismatch` | `captured` | `cannot_verify` | `demo_pilot_only` |
| `dishonest-stale-digest-index` | `captured` | `fail` | `demo_pilot_only` |
| External production trust | `not_captured` | `not_assessed` | `external_production_trust_not_assessed` |

## Dishonest And Incomplete Cases

| case | state | meaning | evidence needed to raise |
| --- | --- | --- | --- |
| `dishonest-source-run-mismatch` | `cannot_verify` | A copied or edited trace/report artifact with inconsistent source/run binding must not be represented as a clean CI-witnessed run. | Rerun from the intended source commit and produce a witness whose source, run, and report artifact digests all bind to the same workflow run. |
| `dishonest-stale-digest-index` | `fail` | A stale digest/index must be treated as tampered or stale until rerun and reindexed. | Re-download or regenerate the artifact set, recompute digests from current bytes, and publish a fresh index tied to the workflow run. |
| `trace-without-oidc` | `cannot_verify` | The witness command ran in CI without OIDC permission; missing `ACTIONS_ID_TOKEN_REQUEST_TOKEN` and `ACTIONS_ID_TOKEN_REQUEST_URL` prevented `ci_witnessed` trust. | Grant the job `id-token: write`, verify OIDC request env vars are present, and rerun witness collection. |

The no-OIDC witness exited `3`, wrote `ci-witness-no-oidc.json`, and recorded
`unexpected_oidc_env=false`.

## Gate Output Meaning

`gate-result.json` reported:

- `local_gate=pass`
- `ci_witness_gate=cannot_verify`
- `audit_grade_gate=cannot_verify`
- `gate_mode=observation`
- `trust_cap=local_observed`

This is verifier-derived fact output. It is not a native merge decision, release
gate, readiness indicator, risk acceptance, policy decision, or production-trust
claim.

## CI Alone vs sdp-trace

Raw CI showed that jobs passed. `sdp-trace` additionally preserved structured
facts that raw CI logs do not preserve as durable verifier artifacts:

- per-run `stdout_digest` and `stderr_digest`;
- run-level `result`, `trust_scope`, and `closure_state`;
- gate fact split between local observation and missing CI/audit evidence;
- witness profile state, reason codes, OIDC subject, source commit, and artifact
  digests;
- explicit no-OIDC `cannot_verify` state with missing identity fields.

Sanitized excerpt (`authority_scope=demo_pilot_only`):

```json
{
  "local_gate": "pass",
  "ci_witness_gate": "cannot_verify",
  "audit_grade_gate": "cannot_verify",
  "gate_mode": "observation",
  "trust_cap": "local_observed",
  "missing_audit_evidence": [
    "ci_oidc_witness",
    "external_witness_checkpoint"
  ]
}
```

## Redaction And Safety

The CI workflow installed `ripgrep` and ran the denylist scan against both
artifact sets with
`docs/research/block-24-redaction-denylist.patterns`
(`sha256=c5ba21129cbc0c969a2d02b46a15bed1cf8c3d48d51643b4eeb6f899150cbbb7`).
Both scans returned `redaction_scan=pass`.

Exact CI command:

```bash
rg --pcre2 -n -f .sdp-trace-src/docs/research/block-24-redaction-denylist.patterns .sdp-trace-report
```

For both artifact sets, `rg` returned exit `1`, which means no matches. Match
count was `0`.

The committed `sdp-trace` repo receives this sanitized summary and artifact
index only. Raw `.sdp-trace-runs/`, raw `.sdp-trace-report/`, and workflow logs
remain in the demo repo's GitHub Actions artifact store.

Artifact download was verified before this report was committed using
`gh run download 25548285336`.

## Block 23 Evidence Coverage

Block 24 gives direct demo evidence for attaching `sdp-trace` to an existing
repository command: three Bazel commands were wrapped in a separate repository
CI workflow, and the resulting run/report artifacts were uploaded with digests
and 14-day retention.

The run produced both machine-readable and reviewer-readable inspection
surfaces. `verify.txt` was produced and indexed, `explain.txt` was produced and
indexed, and `summary.json` reported three observed runs. `gate-result.json`
stayed fact-only: local pass, CI/audit `cannot_verify`.

The GitHub Actions witness path was exercised for the exact demo topology.
`ci-witness.json` passed with `ci_witnessed`; the no-OIDC path recorded
`cannot_verify`, missing identity fields, missing telemetry for CI OIDC, and
exit `3`.

External production trust, another owner, non-GitHub CI, release binary UX, and
compiled Kotlin/JVM compatibility remain outside Block 24 and are recorded as
`not_assessed`.

## Residual States

| topic | state | reason |
| --- | --- | --- |
| Same-owner private repo portability | `not_assessed` | The demo ran under `fall-out-bug`, not a customer owner. |
| Public artifact inspectability | `not_assessed` | Repo and artifacts are private; access limitation is intentional for the first pilot. |
| Release binary acquisition | `not_assessed` | CI built from source ref. |
| Compiled Kotlin/JVM compatibility | `not_assessed` | Bazel tests inspect Kotlin source scope; they do not compile Kotlin. To assess JVM compatibility, the demo needs a test that compiles Kotlin through `kt_jvm_*` rules or runs JVM bytecode. |
| External production trust | `not_assessed` | No customer PKI, release approval, transparency, or production policy evidence was supplied. |
