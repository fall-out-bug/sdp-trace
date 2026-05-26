# Block 25 Review Ledger

Status: initialized for Socratic spec review.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S25-PB-01 | critical | product/customer credibility | "Customer-ready" wording conflicted with explicit non-goals and could imply production readiness. | Accepted and fixed. Goal now says demonstrable technical pilot scoped to a selected compiled target and states one target does not assess the monorepo. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Goal |
| S25-PB-02 | critical | product/customer credibility | Multi-target monorepo scope limitation was not explicit enough. | Accepted and fixed. Required demo behavior and README requirements now state one passing selected target does not imply other targets, BUILD packages, or full monorepo surface. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Required Demo Behavior |
| S25-PB-03 | major | product/customer credibility | Artifact retention policy was vague. | Accepted and fixed. CI artifact contract and README requirements now require configured retention duration and preservation responsibility. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` CI Artifact Contract |
| S25-PB-04 | major | product/customer credibility | `sdp-trace` source fetch vs released binary path was underspecified. | Accepted and fixed. Spec now requires source fetch or released binary mechanism and exact source commit/ref or artifact version. Released binary UX remains `not_assessed` unless separately evaluated. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Documentation Requirements |
| S25-PB-05 | major | product/customer credibility | Hung role review behavior was not actionable. | Accepted and fixed narrower than suggested. Spec keeps trust workflow semantics: hung/off-task reviews are `not_assessed` and must be replaced; closure requires usable role outputs. | MiniMax-M2.7 Socratic review; `25-compiled-ci-demo-pilot.md` Review Plan |
| S25-PB-06 | minor | product/customer credibility | Synthetic Feature Flag / Entitlements domain was not named. | Accepted and fixed. Required Demo Behavior now calls it a synthetic entitlements service surface. | MiniMax-M2.7 Socratic review |
| S25-PB-07 | minor | product/customer credibility | Residual states omitted CI credential and artifact auth scope. | Accepted and fixed. Residual states now include CI credential shape, source-fetch token policy, OIDC token access scope, and artifact upload/download authentication. | MiniMax-M2.7 Socratic review |
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
| S25-TE-06 | minor | tracing/evidence | Sanitized report format was unconstrained. | Accepted and fixed. Spec now requires Markdown under retired research artifacts, not a schema-validated fixture under `examples/`. | ZAI/GLM-5.1 Socratic review |
| S25-TE-07 | minor | tracing/evidence | Demo README needed exact `sdp-trace` ref/commit evidence. | Accepted and fixed. Documentation requirements now require source fetch/released binary mechanism and exact source commit/ref or artifact version. | ZAI/GLM-5.1 Socratic review |
| S25-SP-01 | major | security/privacy | PR and merge checklist omitted security/privacy review. | Accepted and fixed. PR-level review now includes security/privacy. | DeepSeek Socratic review; implementation plan PR And Merge |
| S25-SP-02 | major | security/privacy | Implementation plan did not configure or verify redaction scan despite acceptance criterion. | Accepted and fixed. Slice 2 now owns `scripts/redaction-scan.sh`, records scan details, and expects pass over downloaded artifact sets. | DeepSeek Socratic review |
| S25-SP-03 | minor | security/privacy | Sensitive-output prohibition mentioned committed artifacts but not uploaded CI artifacts. | Accepted and fixed. Non-goals and integration audit now include uploaded CI artifacts. | DeepSeek Socratic review |

## Implementation Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| I25-technical executive-01 | minor | technical executive customer | Stale digest failure could read like an unresolved defect rather than expected dishonest-case evidence. | Accepted and fixed. The report now says stale digest is expected fail evidence and explains the payload mutation. Focused re-review returned `APPROVE` with no findings. | MiniMax-M2.7 review; Mimo focused re-review; retired research artifact Negative Evidence Cases |
| I25-technical executive-02 | minor | technical executive customer | Clean artifact and no-OIDC artifact entry-count asymmetry was not self-documenting. | Accepted and fixed. Demo artifacts now include `artifact-manifest.json`, and the report explains clean versus no-OIDC artifact shape. Focused re-review returned `APPROVE` with no findings. | MiniMax-M2.7 review; Mimo focused re-review; demo run `25555299371` |
| I25-ENG-01 | minor | Head of Engineering | Fixed `/tmp/sdp-trace-demo` workspace was conventionally safe on hosted runners but not exclusive. | Accepted and fixed. Workflow now uses `mktemp -d /tmp/sdp-trace-demo.XXXXXX` and passes the resulting `DEMO_WORKDIR` to later steps. Focused re-review returned `APPROVE`. | ZAI/GLM-5.1 review and focused re-review; demo commit `8d99c13491121a99c5c4cd984ec708dcc1f5025c` |
| I25-ENG-02 | minor | Head of Engineering | Artifact-index temp file could cross filesystems if staged from default temp dir. | Accepted and fixed. `write-artifact-index.sh` now stages the temp file in the artifact root's parent directory, outside the indexed root and on the same filesystem. Focused re-review returned `APPROVE`. | ZAI/GLM-5.1 review and focused re-review |
| I25-ENG-03 | minor | Head of Engineering | Artifact entry-count asymmetry needed an artifact-local explanation for future audits. | Accepted and fixed with `artifact-manifest.json` in both artifact roots. Focused re-review returned `APPROVE`. | ZAI/GLM-5.1 review and focused re-review; downloaded artifacts from run `25555299371` |
| I25-SEC-01 | major | Head of InfoSec | Source-fetch token policy and artifact download authentication were `not_assessed`; reviewer requested they be fully assessed before approval. | Rejected as a Block 25 closure blocker, accepted as residual scope. The Block 25 spec/report explicitly keep source-fetch token policy, artifact download auth, owner independence, production trust, and released-binary UX as `not_assessed`; Block 25 does not claim those surfaces. Focused InfoSec re-review accepted the scope and returned `APPROVE` with no critical/major findings. | DeepSeek review and focused re-review; retired research artifact Trust State |

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| PR25-CODE-01 | none | code/docs correctness | No stale references, broken wording, or accidental green claims detected from the PR evidence packet. | `pass` | Mimo PR-level review on PR #19 |
| PR25-TRACE-01 | none | tracing/evidence | No evidence overclaim found; `cannot_verify`, `fail`, `pass`, and `not_assessed` states remain distinct. | `pass` | MiniMax-M2.7 PR-level tracing/evidence review on PR #19 |
| PR25-REQ-01 | none | requirements-vs-implementation | No unmet Block 25 acceptance criteria or scope drift found from the PR evidence packet. | `pass` | Qwen PR-level requirements review on PR #19 |
| PR25-SEC-01 | none | security/privacy | No leakage or misleading security claim found within Block 25 scope; explicit residual non-goals remain `not_assessed`. | `pass` | MiniMax-M2.7 PR-level security focused review on PR #19 |

## Current Active Demo Role Review Findings

These reviews target the active `fall-out-bug/sdp-trace-demo-jvm-gsd` evidence
packet after T211-T215 closure, not the retired `sdp-trace-demo-ci-pilot`
artifact track.

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| T216-TE-01 | none | technical executive customer | No customer-credibility overclaim, missing residual state, or pilot-claim blocker found in the active evidence packet. | `pass` | `pi` / `openrouter/minimax/minimax-m2.7`, 2026-05-26, role review, no tools, no session; verdict `APPROVE`, no findings. |
| T216-ENG-01 | none | Head of Engineering | No replayability, path consistency, CI/run/artifact consistency, or task-closure blocker found in the active evidence packet. | `pass` | `pi` / `openrouter/z-ai/glm-4.7`, 2026-05-26, role review, no tools, no session; verdict `APPROVE`, no findings. |
| T216-SEC-01 | major | Head of InfoSec | The sanitized report did not include redaction scan command, pattern-file digest, scanned roots, exit code, or state, making the redaction evidence unreplayable. | Accepted and fixed. `docs/reviews/block25-redaction-patterns.txt` now records the pattern set with SHA-256 `494d868e528f8a017b0c320aead26ca227d70d2c31d955b1ff0d0b5e77ca52b3`; `docs/reviews/block25-jvm-gsd-demo-sanitized-report.md` records the `rg` command template, scanned roots, exit code `1`, and `pass` state for no matches without embedding the local artifact path. | `pi` / `openrouter/deepseek/deepseek-v4-pro`, 2026-05-26, role review, no tools, no session; local scan over the downloaded artifact work directory returned exit code `1`. |
| T216-SEC-02 | minor | Head of InfoSec | Artifact download authentication remains `not_assessed`. | Accepted as residual scope. The active report keeps artifact/download authentication-adjacent trust outside Block 25 proof and does not upgrade it to `pass`. | `pi` / `openrouter/deepseek/deepseek-v4-pro`, 2026-05-26, role review, no tools, no session. |
| T216-SEC-RR-01 | none | Head of InfoSec focused re-review | No critical or major security/privacy finding remained after the redaction-scan evidence fix. | `pass` | `pi` / `openrouter/deepseek/deepseek-v4-pro`, 2026-05-26, focused re-review, no tools, no session; verdict `APPROVE`, no critical or major findings. |

## Current Review Evidence State

- Socratic spec review: initial review assessed with MiniMax-M2.7,
  OpenRouter Qwen, ZAI/GLM-5.1, and OpenRouter DeepSeek planes. Valid findings
  have been accepted and fixed. Focused re-review returned `APPROVE` on all four
  planes with no remaining critical or major findings. Minor focused
  observations about JVM pinning, unique artifact-index paths, recursive
  enumeration, CI log secrecy, and redaction-scan output were accepted and
  folded into the spec/plan.
- Implementation review: `pass`; technical executive customer, Head of Engineering, and Head of
  InfoSec role reviews have no remaining critical or major findings after
  fixes and focused re-review.
- PR-level review: `pass`; code/docs correctness, tracing/evidence,
  requirements-vs-implementation, and security/privacy planes have no
  remaining critical or major findings.
- PR CI: GitHub Actions `verify` passed on PR #19 head
  `d9e83620fd1bf221c0cca2ddb50be19ae65e5208` before this ledger update.
- Demo repo CI: `pass` for run `25555299371` on
  `fall-out-bug/sdp-trace-demo-ci-pilot@8d99c13491121a99c5c4cd984ec708dcc1f5025c`.
- Artifact index digest verification: `pass` for downloaded clean and no-OIDC
  artifact roots from run `25555299371`.
- Redaction scan: `pass` for downloaded clean and no-OIDC artifact roots from
  run `25555299371`.
- Negative evidence states: no-OIDC witness gap `cannot_verify` with
  `missing_ci_oidc`; stale digest fixture `fail` with
  `artifact_digest_mismatch`; source/run mismatch fixture `fail`.
- Gate output: `fail`; `gate.exit` was `3`, so gate output is not green closure
  evidence.
- External production trust: `not_assessed`.
