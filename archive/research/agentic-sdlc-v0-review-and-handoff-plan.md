# Agentic SDLC V0 Review And Handoff Plan

Status: planning draft
Date: 2026-05-05

Inputs:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-design.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-demo-readiness.md`
- `archive/research/agentic-sdlc-evidence-substrate-v4-brief.md`
- `archive/research/harness-telemetry-reviews/round8-final-convergence.md`

This plan prepares review and parallel handoff. It is not proof that
implementation or demo execution has completed.

## GPT Spark Handoff Package

Primary task for GPT Spark workers:

```text
Implement Block 10 Agentic SDLC Evidence Substrate V0 in parallel slices.
Do not broaden product claims. Do not close trust claims without live
verification. Product code and active validation must be Go. Do not add
Node/npm/JavaScript/TypeScript/.mjs tooling. Migrate or remove Node from
the active path in this batch. Bash is allowed only as a justified thin
launcher.
```

Recommended Spark work allocation:

1. Schema worker: Slice A.
2. CLI/recorder worker: Slice B.
3. Privacy/DX worker: Slice C.
4. Adapter worker: Slice D.
5. Verifier worker: Slice E.
6. Forensics worker: Slice F.
7. CI/trust worker: Slice G.
8. Demo integration worker: Slice H.
9. Review coordinator: Slice I.

Workers are not alone in the codebase. They must avoid reverting others'
changes, keep ownership boundaries clear, and adapt to changes already
made by other workers.

Quality bar for every worker:

- Clean Architecture boundaries;
- Clean Code;
- TDD for behavior changes;
- CRAP below 5 for changed Go code;
- no TODO/FIXME markers;
- no new Node dependency;
- no unjustified Bash logic.

## pi Review Targets

Review the design before implementation and again after the first local
wrap milestone.

Preferred models from the allowed pi set:

- GLM: `zai/glm-5.1`
- MiniMax: `minimax/MiniMax-M2.7`
- Kimi: `kimi-coding/k2p6`

OpenAI models should not be used for adversarial pi review in this repo.

## GLM Review Prompt

Role: platform / harness owner and implementation critic.

Focus:

- wrapper/process boundary;
- Unix socket adapter feasibility;
- Go module/file ownership;
- testability;
- OpenCode/GSD demo realism;
- mismatch with existing Block 06/09 artifacts.

Command shape:

```bash
rtk pi --model zai/glm-5.1 --no-tools --no-context-files --no-session -p @specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-design.md @specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md @specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-demo-readiness.md "Review Block 10 as a platform/harness owner. Return critical, major, minor findings. Prioritize implementation blockers, unclear control points, broken DX, and demo infeasibility. Do not praise. Tie each finding to a file and section."
```

## MiniMax Review Prompt

Role: CISO / adversarial trust reviewer.

Focus:

- forgery and replay;
- local key weakness;
- witness independence;
- signer authority;
- late/post-hoc trace fabrication;
- overclaiming gate-grade trust;
- missing evidence semantics.

Command shape:

```bash
rtk pi --model minimax/MiniMax-M2.7 --no-tools --no-context-files --no-session -p @specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-design.md @specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md @archive/research/agentic-sdlc-evidence-substrate-v4-brief.md "Review Block 10 as an adversarial CISO. Return critical, major, minor findings. Attack signing, provenance, witness, replay, and local trust claims. Identify any language that overclaims proof. Do not propose broad new scope unless needed to remove a false trust claim."
```

## Kimi Review Prompt

Role: one-file micro-reviewer.

Focus:

- implementation plan clarity;
- missing first milestone tasks;
- parallel ownership conflicts;
- hidden dependencies.

Command shape:

```bash
rtk pi --model kimi-coding/k2p6 --no-tools --no-context-files --no-session -p @specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md "Micro-review this implementation plan only. Return at most 10 findings. Focus on blockers for parallel workers: unclear file ownership, missing tests, hidden dependencies, and claims that cannot be implemented from the plan."
```

## Review Output Files

Save outputs as:

- `archive/research/harness-telemetry-reviews/round9-block10-glm.md`
- `archive/research/harness-telemetry-reviews/round9-block10-minimax.md`
- `archive/research/harness-telemetry-reviews/round9-block10-kimi.md`
- `archive/research/harness-telemetry-reviews/round9-block10-synthesis.md`

The synthesis must classify every valid finding:

- critical;
- major;
- minor;
- accepted limitation;
- invalid / out of scope.

## Implementation Start Gate

Implementation may start when:

- Block 10 design and implementation plan exist;
- GLM, MiniMax, and Kimi review outputs are recorded or explicitly
  marked unavailable;
- valid critical findings are resolved or converted into explicit
  blockers;
- first milestone remains scoped to:

```text
local wrap -> event chain -> verify -> missing evidence -> explain -> tamper fail
```

Do not wait for org-wide reporting, full OpenCode/GSD adapters, or
production CI/OIDC signing before starting the first milestone.
