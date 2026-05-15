<model_routing>
This file captures routing policy, not benchmark authority. Refresh it when model availability or official harness behavior changes.

<codex>
Use Codex as the default implementation and integration agent for this repository:
- scoped Go changes
- local verification
- PR preparation
- final synthesis from reviewer output

Codex output is not trust authority. It still needs live commands, reviewed diffs, and source-bound evidence.
</codex>

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
