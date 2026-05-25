---
name: sdp-trace-router
description: Route sdp-trace work to project-local trust, review, and audit skills while quarantining generic global skills that can bypass evidence rules.
---

<objective>
Keep every sdp-trace session anchored in portable SpecKit language, live evidence, and project-local trust rules before any generic build, review, or planning skill is used.
</objective>

<when_to_use>
Use this skill at the start of sdp-trace work, when choosing a workflow, when a generic/global skill appears applicable, or when balancing pi, Codex, OpenCode, Kimi, GLM, MiniMax, or other reviewer/model harnesses.
</when_to_use>

<routing>
- Block intake, SpecKit delta, implementation, PR-ready closure, or "берем блок в работу" -> use `sdp-trace-trust-workflow`.
- Block intake, SpecKit delta, implementation, PR-ready closure, or "берем блок в работу" -> use `sdp-trace-trust-workflow`.
- Adversarial, Socratic, PR, implementation, requirements, security, evidence, or provenance review -> use `sdp-trace-quality-audit` or the OmPi `reviewer` agent via the `task` tool.
- Repository polish, quality gates, CRAP, MI, docs/DX/security audit, spec drift, or readiness -> use `sdp-trace-quality-audit`.
- Skill authoring or skill cleanup -> use a skill-authoring workflow, but keep project-local files portable across Codex/OpenCode/Claude-style agents.
</routing>

<global_skill_quarantine>
Do not let generic global skills override sdp-trace trust rules.

- `beads`, `deploy`, `hotfix`, `oneshot`, `init`, `feature`, `prototype`, `vision`, and `prd` are out of the active product path unless the user explicitly asks for that external workflow.
- UI/web skills are out of scope unless the changed surface is a user-facing report, docs view, command output, or accessibility/DX artifact.
- Generic `build`, `bugfix`, `tdd`, `test`, `review`, `security-review`, `lint`, and `debug` may assist only inside a project-local workflow and must not close evidence, CI, review, or release trust by themselves.
- Any checked-in report, review ledger, task checkbox, or proof JSON is local artifact context until live-verified or externally signed.
</global_skill_quarantine>

<model_and_harness_balance>
Use the active harness for scoped integration, local verification, review synthesis, and PR finalization. For large approved implementation blocks, prefer OmPi `task` tool subagents in isolated branches instead of loading the whole implementation loop into the active context. Use `.agents/skills` semantics for portable skill discovery. Use the OmPi `reviewer` agent for adversarial evidence review, preferring non-OpenAI, non-Anthropic, and non-Google models when the repo policy applies and the user authorizes the tool/model.

Treat Kimi, GLM, MiniMax, and other long-horizon or swarm-capable models as advisory reviewers or workers, not trust authorities. Record model, harness, timeout, retry, fallback, and missing-evidence status in the review artifact.
</model_and_harness_balance>

<context_loading_order>
1. `AGENTS.md` and this project-local router.
2. The relevant project-local skill.
3. Relevant SpecKit delta, plan, task, evidence, gate, decision, and provenance docs.
4. Source files, schemas, tests, examples, and fixtures directly involved.
5. Fresh command output from local verification or live external checks.
6. Review prose, checked-in ledgers, and historical reports, marked advisory until verified.
</context_loading_order>

<supporting_files>
- For current model and harness routing notes, read `references/model-routing.md`.
- For context and skill hygiene, read `references/context-loading.md`.
</supporting_files>

<success_criteria>
- A project-local sdp-trace skill is selected before a generic global workflow.
- Any unsupported or external workflow is explicitly marked out of active product scope.
- Every trust claim remains backed by live evidence or is marked `not_assessed` / `cannot_verify`.
</success_criteria>
