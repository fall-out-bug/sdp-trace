# Block 25 Review Ledger

Status: initialized for Socratic spec review.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S25-PB-01 | critical | product/buyer credibility | "Buyer-ready" wording conflicted with explicit non-goals and could imply production readiness. | Accepted and fixed. Goal now says demonstrable technical pilot scoped to a selected compiled target and states one target does not assess the monorepo. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Goal |
| S25-PB-02 | critical | product/buyer credibility | Multi-target monorepo scope limitation was not explicit enough. | Accepted and fixed. Required demo behavior and README requirements now state one passing selected target does not imply other targets, BUILD packages, or full monorepo surface. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Required Demo Behavior |
| S25-PB-03 | major | product/buyer credibility | Artifact retention policy was vague. | Accepted and fixed. CI artifact contract and README requirements now require configured retention duration and preservation responsibility. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` CI Artifact Contract |
| S25-PB-04 | major | product/buyer credibility | `sdp-trace` source fetch vs released binary path was underspecified. | Accepted and fixed. Spec now requires source fetch or released binary mechanism and exact source commit/ref or artifact version. Released binary UX remains `not_assessed` unless separately evaluated. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Documentation Requirements |
| S25-PB-05 | major | product/buyer credibility | Hung role review behavior was not actionable. | Accepted and fixed narrower than suggested. Spec keeps trust workflow semantics: hung/off-task reviews are `not_assessed` and must be replaced; closure requires usable role outputs. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Review Plan |
| S25-PB-06 | minor | product/buyer credibility | Synthetic Feature Flag / Entitlements domain was not named. | Accepted and fixed. Required Demo Behavior now calls it a synthetic entitlements service surface. | MiniMax-M2.7 Socratic review |
| S25-PB-07 | minor | product/buyer credibility | Residual states omitted CI credential and artifact auth scope. | Accepted and fixed. Residual states now include CI credential shape, source-fetch token policy, OIDC token access scope, and artifact upload/download authentication. | MiniMax-M2.7 Socratic review |
| S25-ER-01 | critical | engineering/replayability | Compiled JVM target shape was underspecified and allowed a Java-only escape hatch. | Accepted and fixed. Spec now requires pinned `kt_jvm_library` plus `kt_jvm_test`, JVM version, and test framework documentation. | Qwen Socratic review; `25-compiled-ci-demo-pilot.md` Required Demo Behavior |
| S25-ER-02 | critical | engineering/replayability | Artifact index verification lacked a concrete reproducible contract. | Accepted and fixed. Spec now defines deterministic JSON entries with relative paths, SHA-256, size bytes, sorted order, no self-entry, and exact verifier behavior. | Qwen Socratic review; `25-compiled-ci-demo-pilot.md` CI Artifact Contract |
| S25-ER-03 | major | engineering/replayability | Local vs CI artifact generation boundary was ambiguous. | Accepted and fixed. Spec now says local artifacts are debugging only and cannot be cited; only numbered GitHub Actions artifacts on selected source commit count as Block 25 proof. | Qwen Socratic review |
| S25-ER-04 | major | engineering/replayability | Negative cases lacked concrete triggers and observable output distinctions. | Accepted and fixed. Spec now defines no-OIDC, stale digest, absent source/run binding, and contradictory source/run binding triggers and expected states. | Qwen Socratic review |
| S25-ER-05 | major | engineering/replayability | Artifact index verifier ownership/language was unspecified. | Accepted and fixed. Implementation plan now requires `scripts/verify-artifact-index.sh` as a portable shell script using `sha256sum` or `shasum -a 256`. | Qwen Socratic review; implementation plan Slice 2 |
| S25-ER-06 | major | engineering/replayability | Redaction scan tool and scope were underspecified. | Accepted and fixed. Spec and plan now require exact command, pattern digest, artifact roots, exit code/state, and sensitive pattern classes across downloaded artifact sets. | Qwen Socratic review; DeepSeek Socratic review |
| S25-ER-07 | minor | engineering/replayability | Role review criteria and timing were underspecified. | Accepted and fixed by closure rule and implementation plan review steps. | Qwen Socratic review |
| S25-ER-08 | minor | engineering/replayability | Dishonest cases could be confused with same-job post-processing. | Accepted for implementation planning. Negative cases now require distinct trigger/output contracts; implementation may choose separate jobs where useful. | Qwen Socratic review |
| S25-TE-01 | major | tracing/evidence | Source/run mismatch collapsed `cannot_verify` and `fail`. | Accepted and fixed. Spec now separates absent/unresolvable binding (`cannot_verify`) from provably contradictory binding (`fail`). | ZAI/GLM-5.1 Socratic review |
| S25-TE-02 | major | tracing/evidence | Stale-digest fixture independence was not testable. | Accepted and fixed. Spec now requires mutating a non-index artifact after clean index generation while leaving the clean index entry unchanged. | ZAI/GLM-5.1 Socratic review |
| S25-TE-03 | minor | tracing/evidence | Index file's own integrity after move was unaddressed. | Accepted and fixed. Spec states index integrity is not self-indexed and must be handled by CI artifact service metadata/downloaded archive digest where available plus verification of listed files. | ZAI/GLM-5.1 Socratic review |
| S25-TE-04 | minor | tracing/evidence | T216 did not explicitly require recording review findings. | Accepted and fixed. T216 now requires recording all findings and dispositions in the Block 25 review ledger. | ZAI/GLM-5.1 Socratic review |
| S25-TE-05 | minor | tracing/evidence | AC4 needed explicit non-self-index assertion. | Accepted and fixed. AC4 and CI artifact contract now require assertion that the index path is absent from listed entries. | ZAI/GLM-5.1 Socratic review |
| S25-TE-06 | minor | tracing/evidence | Sanitized report format was unconstrained. | Accepted and fixed. Spec now requires Markdown under `docs/research/`, not a schema-validated fixture under `examples/`. | ZAI/GLM-5.1 Socratic review |
| S25-TE-07 | minor | tracing/evidence | Demo README needed exact `sdp-trace` ref/commit evidence. | Accepted and fixed. Documentation requirements now require source fetch/released binary mechanism and exact source commit/ref or artifact version. | ZAI/GLM-5.1 Socratic review |
| S25-SP-01 | major | security/privacy | PR and merge checklist omitted security/privacy review. | Accepted and fixed. PR-level review now includes security/privacy. | DeepSeek Socratic review; implementation plan PR And Merge |
| S25-SP-02 | major | security/privacy | Implementation plan did not configure or verify redaction scan despite acceptance criterion. | Accepted and fixed. Slice 2 now owns `scripts/redaction-scan.sh`, records scan details, and expects pass over downloaded artifact sets. | DeepSeek Socratic review |
| S25-SP-03 | minor | security/privacy | Sensitive-output prohibition mentioned committed artifacts but not uploaded CI artifacts. | Accepted and fixed. Non-goals and integration audit now include uploaded CI artifacts. | DeepSeek Socratic review |

## Implementation Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| _pending_ | _pending_ | _pending_ | Implementation not started. | `not_assessed` | Block 25 activation gate |

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| _pending_ | _pending_ | _pending_ | PR not opened. | `not_assessed` | Block 25 activation gate |

## Current Review Evidence State

- Socratic spec review: initial review assessed with MiniMax-M2.7,
  OpenRouter Qwen, ZAI/GLM-5.1, and OpenRouter DeepSeek planes. Valid findings
  have been accepted and fixed. Focused re-review returned `APPROVE` on all four
  planes with no remaining critical or major findings. Minor focused
  observations about JVM pinning, unique artifact-index paths, recursive
  enumeration, CI log secrecy, and redaction-scan output were accepted and
  folded into the spec/plan.
- Implementation review: `not_assessed`.
- PR-level review: `not_assessed`.
- Demo repo CI: `not_assessed` for Block 25.
- Artifact index digest verification: `not_assessed` for Block 25.
- External production trust: `not_assessed`.
