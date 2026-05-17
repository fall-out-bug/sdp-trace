---
name: sdp-trace-pi-handoff
description: Delegate approved sdp-trace implementation work to Pi or other external coding agents through codex-subagent with worktree isolation, monitoring, review panels, and PR handoff.
---

<objective>
Move large approved implementation blocks out of Codex context and into auditable external subagents while Codex stays the integrator and evidence checker.
</objective>

<when_to_use>
Use this skill when the user asks to send work to Pi, Kimi, OpenCode, GSD2, codex-subagents, external workers, background implementation, or PR-ready delivery through another agent.
</when_to_use>

<contract>
- Commit or explicitly park the reviewed spec before worker launch.
- Build a `context-pack/v1` with `AGENTS.md`, relevant project skills, spec/plan/tasks, and owned files only.
- Launch write-capable workers through `codex-subagent run <runtime> --isolate worktree --background`; pass `--model` when reproducibility, fallback control, or reviewer diversity requires it.
- For Pi runs, prefer explicit provider/model IDs and non-interactive sessionless execution. Local `codex-subagent` must invoke Pi with `--no-session`; if a generated command lacks it, treat the harness as misconfigured before blaming the selected model.
- Record the resolved or requested model/profile. Do not diagnose a stalled Pi run as a model failure unless logs or Pi status prove that root cause.
- Give each worker disjoint file ownership and require `subagent-result/v1`.
- Inspect `status`, `events`, `logs`, `result --structured`, and the worker worktree diff before integration.
- Run independent review panels through `codex-subagent panel run pi` or equivalent read-only runtime.
- Codex integrates or ports worker changes only after checking diffs and fresh verification.
- Do not let workers merge, publish, deploy, close task checkboxes, or claim trust closure.
</contract>

<handoff_workflow>
1. Prepare and commit the spec/review handoff:
   - `git diff --check`
   - `git status --short`
   - scoped commit message, for example `spec: add follow-up readiness hardening plan`
2. Build the worker context:
   - include `AGENTS.md`
   - include `.agents/skills/sdp-trace-router/SKILL.md`
   - include `.agents/skills/sdp-trace-trust-workflow/SKILL.md`
   - include this skill
   - include the approved spec, plan, task, and review files
3. Launch the implementation worker:
   - `codex-subagent context build --subject "<block>" --mode dev --goal "<slice or block goal>" --rule AGENTS.md --rule .agents/skills/sdp-trace-pi-handoff/SKILL.md --file <spec> --file <tasks> --write-allowed --out .codex-subagents/context/<block>.json`
   - `codex-subagent run pi --context-pack .codex-subagents/context/<block>.json --role-template worker --isolate worktree --background --timeout 3600`
   - add `--model <provider/model-id>` when the handoff needs deterministic model selection or a previous run failed with a model/provider-specific error.
   - for Kimi implementation, verify the available Kimi/Moonshot model with `pi --list-models kimi` before launch instead of relying on Pi's default model.
4. Monitor without loading worker context into Codex:
   - `codex-subagent status <run-id>`
   - `codex-subagent events <run-id>`
   - `codex-subagent logs <run-id> --stream stderr`
   - `codex-subagent result <run-id> --structured`
5. Review through separate agents:
   - `codex-subagent panel run pi --context-pack <review-context> --role requirements-reviewer --role code-reviewer --role evidence-reviewer --role security-reviewer --profile review --model <current-provider/model-id> --background`
   - reject hung, empty, generic, or off-task reviewers as `not_assessed`.
6. Integrate:
   - inspect worker diff in its isolated worktree;
   - port or merge only accepted changes;
   - run local verification in the main workspace;
   - commit each verified slice;
   - open/update PR only after review dispositions and verification are recorded.
</handoff_workflow>

<stop_conditions>
Stop and ask the user before continuing when:
- the spec is not reviewed or approval is required by `sdp-trace-trust-workflow`;
- a worker writes outside owned files;
- the worktree cannot be isolated;
- worker output lacks structured result data;
- required review planes cannot run and the user asked for PR-ready closure;
- a worker has no status/result/log/diff progress for several minutes; cancel it, record `cannot_verify`, inspect Pi profile/model/context/tooling evidence, narrow the context, and relaunch with changed run controls only when justified;
- PR creation, merge, publish, or external side effects need credentials or authority not available locally.
</stop_conditions>

<success_criteria>
- The worker ran in an isolated worktree and produced inspectable logs/result/diff.
- Review panels ran independently or missing planes are recorded as `not_assessed` / `cannot_verify`.
- Codex final output cites run IDs, commits, verification commands, PR link or explicit PR blocker, and remaining open evidence states.
</success_criteria>
