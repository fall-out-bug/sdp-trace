# Block 24 Review Ledger

Status: Socratic spec review returned `REVISE` across product/demo
credibility, trace/evidence, CI/witness, and privacy/safety planes. Valid
critical and major findings were accepted into the revised SpecKit delta.
Focused re-review returned `APPROVE` across all four planes. The CTO approved
implementation on 2026-05-08 with the clarified case shape: three clean trace
cases and two intentionally dishonest-trace cases.

## Socratic Review Findings

| id | severity | plane | reviewer/source | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- | --- |
| S24-PC-01 | critical | product/demo credibility | MiniMax-M2.7 | Same-owner demo could look like a controlled fixture rather than customer-style proof. | accepted_fixed | Added Owner Independence Gap and report requirements naming what another owner must provide. |
| S24-PC-02 | critical | product/demo credibility | MiniMax-M2.7 | Tiny passing Go test was too low-signal unless the report shows what `sdp-trace` adds beyond raw CI. | accepted_fixed | Replaced tiny Go with the repo-native Feature Flag / Entitlements Kotlin+Bazel demo service and kept the "CI Alone vs sdp-trace" report requirement. |
| S24-PC-09 | critical | CTO correction | user | The reviewed direction still named a tiny Go demo app, which conflicts with Block 24's purpose of replacing toy pilot proof. | accepted_fixed | Demo scope is now Feature Flag / Entitlements Kotlin+Bazel with a deterministic Bazel test; OpenCode/GSD/model-agent execution remains `not_assessed`. |
| S24-PC-03 | critical | product/demo credibility | MiniMax-M2.7 | Negative `cannot_verify` path lacked customer-readable interpretation. | accepted_fixed | AC7 and Socratic Q5 now require customer interpretation and next evidence. |
| S24-PC-04 | major | product/demo credibility | MiniMax-M2.7 | Pilot report quality bar was undefined. | accepted_fixed | Added Pilot Report Contract and answer classes for the nine Block 23 questions. |
| S24-PC-05 | major | product/demo credibility | MiniMax-M2.7 | `gate` output was disclaimed as non-policy but not given a positive meaning. | accepted_fixed | AC5 and report contract now require "Gate Output Meaning" as verifier-derived fact output. |
| S24-PC-06 | major | product/demo credibility | MiniMax-M2.7 | "Repository attachment" overclaimed same-owner/same-provider proof. | accepted_fixed | Goal and plan now use repository integration wording and owner-independence caveat. |
| S24-PC-07 | minor | product/demo credibility | MiniMax-M2.7 | Artifact freshness/retention was unspecified. | accepted_fixed | Artifact Contract now records storage model, retention, expiration, and rerun requirement. |
| S24-PC-08 | minor | product/demo credibility | MiniMax-M2.7 | `find | sort | tail` run-dir selection is operationally fragile. | deferred_not_assessed | Implementation must either make run-dir selection deterministic or record why the pattern is safe; not a spec-approval blocker. |
| S24-TE-01 | major | trace/evidence | ZAI GLM-5.1 | `observed` and `local_observed` conflated capture and attestation axes. | accepted_fixed | Replaced single state table with capture, attestation, and authority-scope axes. |
| S24-TE-02 | major | trace/evidence | ZAI GLM-5.1 | No taxonomy-level guard separated demo evidence from source-bound proof. | accepted_fixed | Added mandatory `authority_scope=demo_pilot_only` for copied Block 24 artifacts and forbade `source_bound_release` without a separate cycle. |
| S24-TE-03 | major | trace/evidence | ZAI GLM-5.1 | "Replayed" implied bit-identical CI evidence reproducibility. | accepted_fixed | AC2 now says re-executed to produce a new structurally comparable run or inspected from recorded refs. |
| S24-TE-04 | major | trace/evidence | ZAI GLM-5.1 | Artifact storage authority model was unspecified. | accepted_fixed | Artifact Contract now distinguishes git-committed workflow, CI artifact store primary evidence, sanitized durable summaries, retention, and expiration. |
| S24-TE-05 | minor | trace/evidence | ZAI GLM-5.1 | Command sketch redirected outputs without explicit exit-code capture. | accepted_fixed | Implementation plan command sketch now records exit files for verify, explain, report, gate, and witness. |
| S24-TE-06 | minor | trace/evidence | ZAI GLM-5.1 | Gate non-policy wording needed a structural report requirement. | accepted_fixed | AC5 and Pilot Report Contract now require "Gate Output Meaning." |
| S24-CW-01 | critical | CI/witness semantics | OpenRouter Qwen 3.6 Plus | No-OIDC negative job lacked exact GitHub Actions permissions boundary. | accepted_fixed | Implementation plan now requires positive `id-token: write` and negative `id-token: none` job permissions, with OIDC env absence check. |
| S24-CW-02 | critical | CI/witness semantics | OpenRouter Qwen 3.6 Plus | Witness verifier behavior on missing OIDC was not stated. | accepted_fixed | Plan records current CLI contract: writes JSON, exits `3` for `cannot_verify`; workflow captures exit without losing artifact. |
| S24-CW-03 | major | CI/witness semantics | OpenRouter Qwen 3.6 Plus | Ambient CI identity could contaminate negative witness path. | accepted_fixed | Negative path is now a separate no-OIDC job and must verify OIDC env vars are absent. |
| S24-CW-04 | major | CI/witness semantics | OpenRouter Qwen 3.6 Plus | Pinned `sdp-trace` source ref materialization was a placeholder. | accepted_fixed | Plan now includes explicit `actions/checkout@v4` shape for `fall-out-bug/sdp-trace` at `<sdp-trace-source-sha>`. |
| S24-CW-05 | major | CI/witness semantics | OpenRouter Qwen 3.6 Plus | Witness output field requirements were undefined. | accepted_fixed | Added Witness Result Fields table with required extracted fields and forbidden raw fields. |
| S24-CW-06 | minor | CI/witness semantics | OpenRouter Qwen 3.6 Plus | Non-GitHub `not_assessed` versus `cannot_verify` rule was not formalized. | accepted_fixed | Evidence Classification now defines when each state applies. |
| S24-CW-07 | minor | CI/witness semantics | OpenRouter Qwen 3.6 Plus | Demo repo checks and `sdp-trace` PR checks were conflated. | accepted_fixed | Acceptance criteria now split demo workflow checks from `sdp-trace` PR checks. |
| S24-PS-01 | critical | privacy/safety | OpenRouter DeepSeek V4 Pro | Redaction scan was prose, not an executable check. | accepted_fixed | Added Safety Scan Contract, command shape, failure contract, and initial denylist pattern file. |
| S24-PS-02 | major | privacy/safety | OpenRouter DeepSeek V4 Pro | Copy-back artifact boundary lacked allow-list and size limits. | accepted_fixed | Safety Scan Contract now limits copied JSON excerpts to allow-listed top-level fields and 40 lines or 2 KB. |
| S24-PS-03 | major | privacy/safety | OpenRouter DeepSeek V4 Pro | Public/private demo repo visibility lacked safety decision criteria. | accepted_fixed | Implementation plan now requires per-data-class public/private decision before publication. |
| S24-PS-04 | major | privacy/safety | OpenRouter DeepSeek V4 Pro | Authenticated provider URL class was undefined. | accepted_fixed | Denylist pattern file covers token/access_token/sig/signature query params and actions.githubusercontent.com URLs with query strings. |
| S24-PS-05 | major | privacy/safety | OpenRouter DeepSeek V4 Pro | No-OIDC negative witness artifact could still expose CI identity fields. | accepted_fixed | Plan requires field allow-list review and scan for `ci-witness-no-oidc.json`. |
| S24-PS-06 | major | privacy/safety | OpenRouter DeepSeek V4 Pro | Safety scan absence had no failure-detection rule. | accepted_fixed | Safety Scan Contract treats skipped/missing/cannot-run scan as `cannot_verify`, not green. |
| S24-PS-07 | minor | privacy/safety | OpenRouter DeepSeek V4 Pro | Build topology paths and Go env paths were not enumerated. | accepted_fixed | Initial denylist includes runner paths, local user path, `GOPATH`, `GOMODCACHE`, and `/pkg/mod/`. |
| S24-PS-08 | minor | privacy/safety | OpenRouter DeepSeek V4 Pro | Copied excerpts lacked tamper-evident link to original artifact. | accepted_fixed | Artifact Contract and report requirements now require workflow run, artifact refs, retention, and digests. |

## Focused Re-Review

| plane | reviewer/source | result | remaining critical/major findings |
| --- | --- | --- | --- |
| product/demo credibility | MiniMax-M2.7 | `APPROVE` | none |
| trace/evidence and claim-boundary | ZAI GLM-5.1 | `APPROVE` | none |
| CI/witness semantics | OpenRouter Qwen 3.6 Plus | `APPROVE` | none |
| privacy/safety | OpenRouter DeepSeek V4 Pro | `APPROVE` | none |
| CTO correction: replace tiny Go with Feature Flag / Entitlements Kotlin+Bazel | MiniMax-M2.7 | `APPROVE` | none |
| CTO correction: Kotlin/Bazel evidence feasibility | OpenRouter Qwen 3.6 Plus | `not_assessed` | unusable response attempted to inspect a demo repo instead of reviewing the supplied SpecKit packet |
| CTO correction: Kotlin/Bazel evidence feasibility replacement | OpenRouter DeepSeek V4 Pro | `APPROVE` | none |

## Implementation Approval

| date | approver/source | decision | implementation constraint |
| --- | --- | --- | --- |
| 2026-05-08 | CTO/user | approved reviewed direction | Implement Feature Flag / Entitlements Kotlin+Bazel demo repo, GitHub Actions first, source-built `sdp-trace`, sanitized copy-back only, three clean cases, and two intentionally dishonest-trace cases. |

## Implementation Evidence

| item | state | evidence |
| --- | --- | --- |
| Demo repo selected | `pass` | `fall-out-bug/sdp-trace-demo-ci-pilot`, private repository, commit `e370d1c00df8a7e7859adc284480563a269e64ca` |
| Demo CI run | `pass` | GitHub Actions run `25548285336` completed successfully on 2026-05-08 |
| Clean cases | `pass` | `clean-feature-flag`, `clean-entitlement-matrix`, `clean-audit-scope` all recorded `observed` with exit `0` |
| CI witness | `pass` | `ci-witness.json` recorded `status=pass`, `established_trust_scope=ci_witnessed`, `reason=ci_identity_present` |
| No-OIDC witness gap | `cannot_verify` | `ci-witness-no-oidc.json` recorded `missing_ci_oidc`, missing identity fields, and exit `3` |
| Dishonest source/run mismatch | `cannot_verify` | `dishonest-source-run-mismatch` recorded `source_run_binding_mismatch` and `ci_witness_not_upgraded` |
| Dishonest stale digest index | `fail` | `dishonest-stale-digest-index` recorded `artifact_digest_mismatch` and `stale_index` |
| Safety scan | `pass` | Denylist pattern sha256 `c5ba21129cbc0c969a2d02b46a15bed1cf8c3d48d51643b4eeb6f899150cbbb7`; CI scans returned `redaction_scan=pass` |
| Raw artifact retention | `pass` | Artifacts `6876152707` and `6876147448`, expiration 2026-05-22 |
| Public/customer owner portability | `not_assessed` | Demo repo is private and same-owner |
| Compiled Kotlin/JVM compatibility | `not_assessed` | Bazel tests inspect target-scoped Kotlin source; they do not compile Kotlin |

Implementation review and PR-level review remain open.

## Implementation Review Disposition

| finding | severity | plane | source | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| I24-TE-01 | major | trace/evidence | OpenRouter DeepSeek V4 Pro | accepted_fixed | Added an Evidence Classification table covering clean runs, report artifacts, witness outputs, dishonest cases, and external production trust. |
| I24-TE-02 | major | trace/evidence | OpenRouter DeepSeek V4 Pro | accepted_fixed | Added exact `rg --pcre2` command, exit semantics, match count `0`, and pass state to report and artifact index. |
| I24-TE-03 | minor | trace/evidence | OpenRouter DeepSeek V4 Pro | accepted_fixed | Added explicit missing-telemetry wording to the no-OIDC customer-question row. |
| I24-RI-01 | minor | requirements-vs-implementation | OpenRouter Xiaomi MiMo V2.5 Pro | accepted_fixed | Added a sanitized `gate-result.json` excerpt showing local/CI/audit state split. |
| I24-RI-02 | minor | requirements-vs-implementation | OpenRouter Xiaomi MiMo V2.5 Pro | accepted_fixed | Added `generated_at` and profile-state extracts for clean and no-OIDC witnesses. |
| I24-RI-03 | minor | requirements-vs-implementation | OpenRouter Xiaomi MiMo V2.5 Pro | accepted_fixed | Exact scan command recorded with exit and match-count semantics. |
| I24-CC-01 | not_assessed | code/correctness | ZAI GLM-5.1 | replaced_unusable | Reviewer returned stale/off-task content for a different Block 24 shape; not counted as review evidence. |
| I24-CC-02 | major | code/correctness | OpenRouter Qwen 3.6 Plus | false_positive | Reviewer judged the denylist pattern file missing because it reviewed only the implementation diff. The file exists at `docs/research/block-24-redaction-denylist.patterns` from commit `5ced37b`; current blob `a6225493d2f935503a684394ad656659fc8715fe`. |
| I24-CC-03 | minor | code/correctness | OpenRouter Qwen 3.6 Plus | accepted_fixed | Added exact `sdp-trace` source commit `f66aa1c4619f0f6a2d56f602da9e0135b00e4a84` to report and artifact index. |
| I24-CC-04 | minor | code/correctness | OpenRouter Qwen 3.6 Plus | accepted_fixed | Added independent SHA-256 digests for `ci-witness.json` and `ci-witness-no-oidc.json`. |
| I24-CC-05 | minor | code/correctness | OpenRouter Qwen 3.6 Plus | accepted_fixed | Added artifact download verification statement using `gh run download 25548285336`. |

Focused re-review is required for the fixed documentation before implementation
review can be counted complete. PR-level review remains open.
