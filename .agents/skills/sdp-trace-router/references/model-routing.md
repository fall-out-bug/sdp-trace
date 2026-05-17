<model_routing>
This file captures routing policy, not benchmark authority. Refresh it when model availability or official harness behavior changes.

<codex>
Use Codex as the default integration agent for this repository:
- small scoped Go changes
- local verification
- PR preparation
- final synthesis from reviewer output

Codex output is not trust authority. It still needs live commands, reviewed diffs, and source-bound evidence.
</codex>

<codex_subagents>
Use `codex-subagent` when the user asks to hand work to Pi/OpenCode/GSD2, when a block needs background workers, or when the implementation would consume substantial Codex context.

- write-capable work: `codex-subagent run <runtime> --isolate worktree --background`
- review work: `codex-subagent panel run pi --profile review`
- inspection: `status`, `events`, `logs`, `result --structured`

Subagent output is work product. Codex must inspect diffs and run verification before integration.
</codex_subagents>

<opencode>
OpenCode discovers `.agents/skills/<name>/SKILL.md`, so project-local skills are the portable format for this repo. Keep frontmatter to `name` and `description` for maximum compatibility; extra frontmatter may be ignored by OpenCode.

Use OpenCode agents with read-only permissions for review roles. Editing agents need explicit file ownership.
</opencode>

<pi_review_models>
For adversarial review, prefer model diversity and cold context:
- GLM-5.1 or current Z.AI coding model for long-horizon engineering review.
- MiniMax M2.7 or current MiniMax coding model for skill-adherence and multi-step workflow review.
- Kimi K2.5 or current official Moonshot coding model for wide-context review; only pin newer versions after primary-source confirmation.

These reviewers are advisory. Record model name, provider/harness, prompt, timeout, retry, replacement, and disposition.
</pi_review_models>

<safety>
Do not grant internet or shell access to a reviewer just because the model is strong. External content can contain prompt injection. Prefer artifact+contract-only prompts and read-only tool permissions.
</safety>
</model_routing>
