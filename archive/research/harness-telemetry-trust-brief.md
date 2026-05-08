# Harness-Neutral Telemetry Trust Layer Brief

Status: discussion draft; not committed
Date: 2026-05-05

## 1. Problem

Companies are already adopting AI SDLC through many harnesses: OpenCode, pi, Kilo, Superpowers, Getting Shit Done, internal agents, custom prompts, CI wrappers, and team-specific repo initialization. A CTO will not accept a product that requires replacing the whole AI SDLC process before it can provide value.

The product question is whether `sdp-trace` can act as a harness-neutral telemetry trust layer:

- What did the agent promise to do?
- What actually happened?
- Which model, harness, commands, tools, files, tests, and evidence were involved?
- Which claims are observed, self-reported, missing, or unverifiable?
- Can the timeline be shown to have existed before a gate or merge?
- Can the trace resist mutation, deletion, reordering, replay, and backfill?
- Can the CTO see whether AI-assisted delivery is degrading team control?

## 2. Product Thesis

Agents cannot simply be asked to write trustworthy telemetry. Voluntary envelope compliance is not a trust boundary.

The target model:

- Agents may enrich telemetry with intent, rationale, plan, claims, and self-report.
- Recorder components must capture critical events independently where control points exist.
- Verifier output must distinguish observed evidence from self-reported claims.
- Missing telemetry is a verifier state, not an empty row in a report.
- `sdp-trace` should not be a mandatory harness. It should be a flight recorder, telemetry contract, verifier, witness model, and query surface that can be attached to existing harnesses.

## 3. Telemetry Capture Model

Telemetry can be captured from multiple layers. Each event must declare its source, evidence role, trust scope, and verifier state.

| Layer | Captures Well | Misses / Risks | Trust Strength |
| --- | --- | --- | --- |
| Harness adapter | Task lifecycle, model requested, tool requested, subagent boundaries, agent claims | Harness-specific APIs, version drift, bypass risk | Medium if managed |
| Tool wrapper | Read/write/edit/bash/search invocations, tool arguments, file-level actions | Harness intent and model identity unless correlated | Strong for observed tool actions |
| Shell wrapper | Commands, cwd, argv digest, start/finish time, exit state | Non-shell file edits and model/tool calls | Strong for build/test evidence |
| File watcher / VCS diff | File mutations, path scope, tree/diff digests | Attribution can be weak without tool/shell correlation | Strong for mutation evidence |
| Git / PR metadata | Commit, branch, PR, merge point, review metadata | Local agent behavior before commit | Strong inside repo boundary |
| CI collector | Build/test commands, artifacts, controlled environment, merge-time chain head | Local behavior before CI | Strong for gate evidence |
| Remote witness / append-only store | Chain head existed at time T, anti-backfill anchor | Does not capture events by itself | Strongest anti-forgery boundary |

## 3A. Integration Placement Options

Open question for review: where should `sdp-trace` integrate to capture reliable telemetry with the least adoption friction?

| Integration Point | Captures Well | Misses / Risks | Trust Strength |
| --- | --- | --- | --- |
| OpenCode / pi / Kilo plugin | Agent lifecycle, model choice, prompts/claims, tool calls, task context, subagent boundaries | Harness-specific APIs, version drift, teams may use different harnesses | Medium if plugin cannot be bypassed in managed mode |
| Tool layer wrapper | Read/write/edit/bash/search, file edits, command starts, outputs, exit codes | Misses model reasoning and harness-level intent if not passed down | Strong for observed actions |
| LLM gateway / proxy | Model request/response metadata, model id, token usage, latency, provider, prompt hash, response hash | Misses local shell/file actions; may capture sensitive prompts; cannot prove what agent did with response | Strong for model provenance, weak for code mutation |
| Shell wrapper | Commands, cwd, argv digest, exit code, duration | Misses non-shell file edits and model/tool calls | Strong for build/test evidence |
| File watcher / VCS layer | Changed files, tree/diff digest, scope violations | Attribution to agent/command may be weak without tool/shell correlation | Strong for mutation evidence |
| CI integration | Build/test evidence, merge-time chain head, controlled environment | Only sees post-local state; late for local agent behavior | Strong for gate evidence |
| Remote witness | Chain head existed at time T, anti-backfill anchor | Does not capture events by itself | Strongest anti-forgery boundary |

Review questions:

- Is harness-level integration enough, or must the design assume harness bypass?
- Can OpenCode, pi, and Kilo expose tool calls and subagent events reliably enough?
- Is LLM gateway capture useful as provenance, or too dangerous because of prompt/privacy risk?
- Should model gateway logs be signed independently from harness telemetry?
- Can gateway telemetry prove model identity better than self-reported harness metadata?
- How should gateway request ids be correlated with local tool, shell, file, and CI events?
- What is the minimum viable integration: harness plugin first, gateway first, or local recorder wrapper first?
- What state should verifier emit if gateway telemetry exists but local tool telemetry is missing?
- What state should verifier emit if local tool telemetry exists but gateway telemetry is missing?

Working hypothesis: the architecture should be layered, not either/or.

1. Harness plugin for task, model, tool intent, and subagent boundaries.
2. Tool/shell/file recorder for observed local actions.
3. LLM gateway for model provenance, token/cost telemetry, and prompt/response digests.
4. CI and remote witness for external anchoring.

Gateway integration is valuable, but it cannot replace local recorder telemetry. It can prove a model call happened; it cannot prove the agent changed only these files and ran these tests.

## 4. Agent Cooperation Model

Agents are not trusted writers. They are participants.

Agent-written telemetry is allowed for:

- task interpretation;
- plan;
- expected files;
- claimed tests;
- rationale;
- completion summary.

Agent-written telemetry is not sufficient for:

- command execution;
- file mutation proof;
- test pass/fail;
- final chain integrity;
- witness proof;
- production trust.

If an agent bypasses telemetry:

- managed mode: verifier fails or blocks;
- unmanaged mode: verifier emits `missing_telemetry`;
- late attach: verifier emits `not_assessed` before attach boundary.

## 5. Signing Model

| Level | Signer | Meaning |
| --- | --- | --- |
| Agent-signed | Agent has key | Weak; proves only self-report continuity |
| Local recorder-signed | Local process or daemon | Local integrity; not audit-grade if user controls host |
| CI-signed | CI identity signs chain head and artifacts | Stronger; proves run reached controlled CI |
| External witness-signed | Remote append-only witness, timestamp service, or transparency log | Strongest anti-backfill boundary |

Critical rule: if the agent can access the signing key, signature does not prove truth. It only proves that the agent signed a claim.

Target model:

- local event chain uses hash linking;
- local recorder signs checkpoints;
- CI or remote witness signs chain head;
- external witness is required for accountability-grade claims;
- local-only telemetry must never be described as audit-grade trust.

## 6. Anti-Forgery Model

Threats to handle:

- event mutation;
- event deletion;
- event reordering;
- backfilled telemetry after failure;
- replay of old valid run;
- fake witness record;
- local key theft;
- agent bypassing wrapper;
- CI signing unverified telemetry;
- partial telemetry with successful-looking summary;
- task rewrite after command evidence;
- redaction hiding critical evidence;
- gateway telemetry without local action telemetry;
- local action telemetry without model provenance.

Verifier must not repair these. It must report states:

- `pass`
- `fail`
- `cannot_verify`
- `not_assessed`
- `missing_telemetry`
- `local_only`
- `witnessed`
- `externally_witnessed`

## 7. Adoption Model

The product must support gradual adoption:

1. Observe-only: attach recorder to existing harness; no blocking; show gaps.
2. Managed local: team agrees to wrapper/adapter; missing telemetry becomes local fail.
3. CI witnessed: CI signs run chain head and test/build evidence.
4. External witnessed: chain heads are anchored outside mutable repo/workspace.

This matters because CTOs will reject "replace your harness first".

## 8. CTO Query Surface

CTO does not want raw JSON. Required questions:

- What task was locked?
- Did the task change? Was it superseded or rewritten?
- Which model and harness ran?
- What commands ran?
- Which files changed?
- Were changes inside allowed scope?
- Which tests/builds ran?
- What evidence is missing?
- Was telemetry captured live or attached late?
- Is this local-only or witnessed?
- Can this run support a gate decision?
- Where did the agent claim something that evidence does not support?
- Is the team drifting into unapproved harnesses, agents, prompts, or skills?

No opaque health score.

## 9. Known Current Gaps

Current Block 09 covers:

- event chain;
- witness mismatch;
- late attach boundary;
- supersession;
- redaction states;
- query surface;
- source-bound proof discipline.

Still missing:

- harness adapter contract;
- tool/shell/file capture implementation;
- OpenCode/pi/Kilo integration plan;
- LLM gateway capture and privacy model;
- remote witness protocol;
- CI signing recipe;
- harness registry / approved agent governance;
- CTO metric pack for degradation;
- Kotlin+Bazel real demo evidence.

## 10. Review Task

Review this proposal as an adversarial product/system reviewer.

Return:

```text
VERDICT: CHANGES_REQUIRED | ACCEPTABLE_WITH_GAPS | REJECTED

Critical blockers:
- ...

Major gaps:
- ...

False assumptions:
- ...

Minimum viable changes:
- ...

Questions before demo:
- ...

Attack scenarios not covered:
- ...
```

Tie every finding to at least one area:

- telemetry capture;
- integration placement;
- agent cooperation;
- signing;
- anti-forgery;
- adoption;
- CTO usefulness.
