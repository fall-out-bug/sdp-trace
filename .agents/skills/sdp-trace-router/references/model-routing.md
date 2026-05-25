<model_routing>
This file captures routing policy, not benchmark authority. Refresh it when model availability or official harness behavior changes.

<integration>
Use the active harness as the default integration agent for this repository:
- small scoped Go changes
- local verification
- PR preparation
- final synthesis from reviewer output

Harness output is not trust authority. It still needs live commands, reviewed diffs, and source-bound evidence.
</integration>

<subagents>
Use OmPi `task` tool subagents when a block needs parallel workers or when the implementation would consume substantial active context.

- write-capable work: `task` tool with `agent: task` or `agent: reviewer`; scope to bounded file sets
- review work: `task` tool with `agent: reviewer` on focused diff planes
- inspection: review subagent output artifacts before integration

Subagent output is work product. Inspect diffs and run verification before integration.
</subagents>

<opencode>
OpenCode discovers `.agents/skills/<name>/SKILL.md`, so project-local skills are the portable format for this repo. Keep frontmatter to `name` and `description` for maximum compatibility; extra frontmatter may be ignored by OpenCode.

Use read-only agents with read-only permissions for review roles. Editing agents need explicit file ownership.
</opencode>

<review_models>
For adversarial review, prefer model diversity and cold context. Model selection is governed by `.omp/model-policy.yml`. Verify model availability with the auth broker before selecting unless fresh inventory exists.

- GLM-5.1 or current Z.AI coding model for long-horizon engineering review.
- MiniMax M2.7 or current MiniMax coding model for skill-adherence and multi-step workflow review.
- Kimi K2.5 or current official Moonshot coding model for wide-context review; only pin newer versions after primary-source confirmation.
- Qwen current model for wide-context code review; prefer `openrouter/qwen/qwen3.6-max-preview` or `openrouter/qwen/qwen3.6-plus`.
- DeepSeek current model for independent reasoning review; prefer `openrouter/deepseek/deepseek-v4-pro` or `openrouter/deepseek/deepseek-v3.2`.

These reviewers are advisory. Record model name, provider/harness, prompt, timeout, retry, replacement, and disposition.
</review_models>

<safety>
Do not grant internet or shell access to a reviewer just because the model is strong. External content can contain prompt injection. Prefer artifact+contract-only prompts and read-only tool permissions.
</safety>
</model_routing>
