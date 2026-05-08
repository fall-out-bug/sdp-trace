# Round 4: Consolidation

Status: discussion draft; not committed
Date: 2026-05-05
Inputs:
- `archive/research/harness-telemetry-trust-brief.md`
- `archive/research/harness-telemetry-reviews/round1-minimax-cto-trust.md`
- `archive/research/harness-telemetry-reviews/round1-glm-platform-architect.md`
- `archive/research/harness-telemetry-reviews/round1-deepseek-operational.md`
- `archive/research/harness-telemetry-reviews/round1-mimo-constrained-agent.md`
- `archive/research/harness-telemetry-reviews/round2-minimax-forgery-red-team.md`
- `archive/research/harness-telemetry-reviews/round2-glm-forgery-red-team.md`
- `archive/research/harness-telemetry-reviews/round2-deepseek-forgery-red-team.md`
- `archive/research/harness-telemetry-reviews/round3-qwen-adoption-dx.md`
- `archive/research/harness-telemetry-reviews/round3-mimo-adoption-dx.md`
- `archive/research/harness-telemetry-reviews/round3-kimi-head-eng-micro.md`
- `archive/research/harness-telemetry-reviews/round3-minimax-commercial-ops.md`

This file is a human consolidation of adversarial discussion outputs.
It is not source-bound proof, not product closure evidence, and not a
trusted release claim.

## Product Correction

The initial concept is directionally right but overclaims the current
product boundary.

Correct product statement:

> `sdp-trace` is a portable evidence recorder and verifier for
> AI-assisted delivery. It records observed events, event provenance,
> evidence gaps, and witness scope. It does not make unverifiable work
> true. It makes unsupported claims, missing evidence, and weak trust
> boundaries visible enough for governance gates.

Rejected product statement:

> `sdp-trace` is a harness-neutral trust layer that can prove agentic
> delivery truth after being attached to any workflow.

Reason: without live capture, expected-event contracts, and a witness
anchor, local telemetry can be internally consistent while still being
post-hoc fabricated. The product must show this honestly instead of
turning structural consistency into trust.

## Core Decision

Split the product into two explicit trust bands.

1. Local observation band.
   - Captures local session evidence.
   - Useful for debugging, incident reconstruction, and gap reporting.
   - Emits `local_only` trust scope.
   - Must be displayed as non-audit-grade.
   - Must not support gate closure by itself.

2. Witnessed gate band.
   - Requires CI/OIDC identity or an external witness to anchor chain
     heads before the gate decision.
   - Can support governance gates when verifier verdict, trust scope,
     and completeness contract all meet policy.
   - Must reject or downgrade missing, late, unverifiable, or
     uncorrelated telemetry.

This split resolves the strongest critique from Round 2 and Round 3:
local telemetry is valuable, but it is not a trust boundary.

## Consolidated Findings

| Finding | Source | Severity | Disposition | Reason | Demo Impact |
| --- | --- | --- | --- | --- | --- |
| Single verifier enum conflates verdict, trust scope, and completeness. | Round 1 MiMo, Round 1 GLM, Round 2 GLM, Round 3 Qwen | critical | accepted | Consumers must distinguish `pass` from `local_only` and from `partial`. Otherwise the product creates false assurance. | blocker |
| Local recorder signatures are not a trust boundary when agent and recorder share user context. | Round 1 all, Round 2 all, Round 3 MiMo | critical | accepted | A signature proves only key access unless key identity, process isolation, and lifecycle are specified. | blocker |
| Local hash chain cannot prove live capture. | Round 1 MiMo, Round 2 all | critical | accepted | A post-hoc chain can be internally consistent. Live proof needs CI/OIDC, remote witness, or checkpoint anchoring before gate time. | blocker |
| Missing telemetry must be first-class output, but it cannot prove whether a run was deleted or never existed. | Round 2 all | critical | accepted | Absence is useful governance signal, not forensic certainty. Do not overstate deletion detection. | demo requirement |
| Expected-event/completeness contract is required before any meaningful verifier verdict. | Round 1 GLM, Round 2 all, Round 3 Qwen/MiMo | critical | accepted | Without required events, an empty or selectively captured run can look clean. | blocker |
| No harness adapter contract exists. | Round 1 all, Round 3 Kimi/Qwen/MiniMax | critical | accepted | "Harness-neutral" cannot mean "no adapter contract." It must mean portable schemas plus small adapters. | demo requirement |
| No installable artifact exists. | Round 3 Kimi/Qwen/MiniMax | critical | accepted | Teams cannot adopt a brief. The next build slice needs at least one executable recorder/verifier path. | blocker |
| LLM gateway capture is valuable but privacy-blocked without a data model. | Round 1, Round 2, Round 3 all | critical | accepted | Prompt/response capture or naive hashes can leak proprietary data. Gateway work needs a separate privacy spec. | defer from demo |
| Gateway-only telemetry cannot prove local action. | Round 2 all | critical | accepted | A model call proves model activity, not file changes, commands, tests, or scope compliance. | demo caveat |
| Local action telemetry without model provenance is incomplete for agent provenance. | Round 2 all | major | accepted | It can support build/test evidence, but cannot prove which model produced the plan or edit. | demo caveat |
| CI must not sign unverified telemetry. | Round 2 all, Round 3 MiniMax/MiMo | critical | accepted | CI signature must cover verifier result, chain head, source commit, and policy context, not raw untrusted uploads. | blocker |
| Shell wrapper is the fastest visible artifact but weak as a developer-laptop enforcement mechanism. | Round 3 Qwen, Round 3 MiMo | major | accepted with split | Use it for local observation and demo visibility, not as the primary audit boundary. | demo requirement with caveat |
| CI/OIDC is the first credible gate anchor. | Round 3 MiMo, Round 3 MiniMax | critical | accepted | CI identity is easier to govern than local laptop keys and aligns with gate decisions. | blocker |
| File watcher is noisy and expensive. | Round 1 MiMo, Round 3 MiMo | major | defer | Useful later for attribution, but too expensive and flaky for first credible demo. | defer from demo |
| Raw stdout/stderr, argv, file paths, prompts, and responses are sensitive telemetry. | Round 3 all | critical | accepted | The product must default to digests and redaction, with explicit opt-in for raw capture. | blocker |
| CTO query surface overclaims current capability. | Round 1, Round 3 Qwen/MiniMax | critical | accepted | v0 should answer a small set of evidence questions instead of pretending to answer full delivery governance. | demo requirement |
| Managed mode is undefined. | Round 1 MiMo, Round 3 MiniMax/Kimi | critical | accepted | If blocking is policy/CI behavior, say so. The recorder cannot magically block a harness after the fact. | blocker |
| Remote witness protocol is undefined. | Round 1, Round 3 Kimi/MiniMax | major | defer | Needed for audit-grade local-session proof, but CI/OIDC can be first witnessed gate anchor. | future |
| Harness registry is governance, not only schema. | Round 1 MiMo, Round 3 Kimi/MiniMax | major | defer | Needed to detect drift into unapproved harnesses, but not required for first Kotlin+Bazel evidence demo. | future |
| Kotlin+Bazel demo needs real Bazel evidence parsing. | Round 3 Kimi | major | accepted | Demo must consume real `bazel test`/BEP or test XML output, not hand-written fixtures. | demo requirement |
| Verifier failure explanations are required for DX. | Round 3 MiMo | major | accepted | Developers need event-level cause: missing event, broken hash, trust-scope downgrade, or policy mismatch. | demo requirement |
| Event volume, retention, and access control are unscoped. | Round 1 MiMo, Round 3 MiMo/MiniMax | major | accepted | Even local telemetry can leak sensitive operational data. | demo caveat |

## v0 Demo Boundary

The Kotlin+Bazel demo must be honest and narrow.

Keep:
- One local recorder path for observed command execution.
- One event schema for command/build/test evidence.
- One Bazel evidence parser using real Bazel output.
- One verifier state model with three axes:
  - verdict: `pass`, `fail`, `cannot_verify`, `not_assessed`
  - trust scope: `local_only`, `ci_witnessed`, `externally_witnessed`
  - completeness: `complete`, `partial`, `missing_telemetry`, `unknown`
- One CI/OIDC or CI identity anchor for gate-grade evidence.
- One CTO report with plain-language caveats.

Cut:
- LLM gateway capture.
- Prompt/response telemetry.
- Remote witness service.
- Multi-harness support.
- File watcher attribution.
- Full CTO dashboard.
- Unapproved harness governance.
- Agent self-reported telemetry as trusted evidence.

Allowed demo claim:

> `sdp-trace` can record observed local build/test events, hash-chain
> them, verify structural integrity, and show whether the evidence is
> local-only or CI-witnessed. For CI-witnessed runs, it can support a
> gate decision only when the expected-event contract is satisfied.

Forbidden demo claim:

> `sdp-trace` proves that an arbitrary agent honestly performed a task
> across any harness.

## Minimum Spec Work Before Implementation

The next development block should produce these artifacts before
building the demo:

1. Event schema.
   - Session identity.
   - Event identity.
   - Source identity.
   - Event type.
   - Timestamp.
   - Payload digest fields.
   - Previous event hash.
   - Optional correlation ids.

2. Verifier state model.
   - Separate verdict, trust scope, and completeness.
   - Define composition rules when layers disagree.
   - Define display language for local-only and partial evidence.

3. Completeness contract.
   - Required events for local build/test run.
   - Required events for CI-witnessed gate run.
   - Downgrade rules for missing, late, or uncorrelated events.

4. Signing and anchoring profile.
   - Local signatures are structural only.
   - CI/OIDC anchor signs verified chain head, source commit, verifier
     version, and policy profile.
   - Agent-owned keys are never trusted as recorder identity.

5. Privacy and retention profile.
   - Digest-only by default.
   - No raw prompts or raw responses in v0.
   - Argv/stdout/stderr redaction rules.
   - Local store permissions and retention.

6. CTO report contract.
   - Do not show opaque scores.
   - Show verdict, trust scope, completeness, missing evidence,
     source commit, commands/tests observed, and gate usability.

## Integration Placement Decision

The discussion produced a useful split:

| Placement | v0 decision | Reason |
| --- | --- | --- |
| OpenCode/pi/Kilo adapter | defer except one optional spike | Valuable for task/model/tool intent, but no adapter contract is defined yet. |
| Tool wrapper | defer | Harness internals vary; wrapper points may not exist. |
| Shell wrapper | use for local observation only | Fastest visible artifact, but easy to bypass and weak as enforcement. |
| Git hook / commit metadata | use as supporting context | Ties evidence to source commit without pretending tool-level causality. |
| Bazel output parser | use | Demo needs real build/test evidence. |
| CI/OIDC | use for gate anchor | Strongest practical first witness boundary. |
| LLM gateway | defer | Useful for model provenance, but privacy and correlation are not ready. |
| Remote witness | defer | Needed later for stronger local-session liveness proof. |

## Open Questions For Next Block

1. Which CI identity is the first target: GitHub Actions OIDC, GitLab CI,
   Jenkins, or a local simulated CI witness for demo only?
2. What exact Bazel output format will the demo parse first: Build Event
   Protocol JSON, XML test logs, or command exit evidence only?
3. Is the first local recorder a shell wrapper, a command launcher CLI, or
   an explicit `sdp-trace run -- <command>` wrapper?
4. What is the minimum CTO report: one Markdown report, JSON API, or CLI
   table?
5. Do we require task/spec identity in v0, or do we restrict v0 to
   build/test evidence and source commit identity?

## Recommended Next Step

Lock the next block as "telemetry event contract and witnessed gate
profile" rather than "harness integration."

Reason: if the state model, completeness contract, and witness profile
are wrong, every adapter and gateway integration will only produce more
untrustworthy data. The demo should prove one small chain honestly,
not simulate full governance.
