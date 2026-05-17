<model_policy>
Use model diversity to reduce shared blind spots, not to outsource judgment.

<roles>
- Codex/GPT: implementation review, Go-specific integration, local verification synthesis.
- OpenCode: portable skill execution and read-only reviewer harnesses using `.agents/skills`.
- GLM-5.1 or current Z.AI coding model: long-horizon engineering, optimization, and architecture doubt. Prefer `openrouter/z-ai/glm-5.1` or `zai/glm-5.1` after confirming credentials; avoid `glm-4.6` as a normal reviewer default.
- MiniMax M2.7 or current MiniMax coding model: skill adherence, multi-step workflow review, and tool-loop critique.
- Kimi K2.5 or current official Moonshot coding model: wide-context code review and alternate reasoning. Do not pin K2.6 in durable policy without primary-source confirmation.
- Qwen current coding/reasoning model: wide-context code review. Prefer `openrouter/qwen/qwen3.6-max-preview` or `openrouter/qwen/qwen3.6-plus`; avoid `qwen3-coder` as a normal reviewer default.
- DeepSeek current reasoning/coding model: independent reasoning review. Prefer `openrouter/deepseek/deepseek-v4-pro` or `openrouter/deepseek/deepseek-v3.2`; avoid `deepseek-chat-v3.1` as a normal reviewer default.
</roles>

<requirements>
- Before selecting a reviewer family, run `pi --list-models qwen3.6`, `pi --list-models deepseek`, `pi --list-models glm`, or the relevant family query unless a fresh model inventory was already captured in this session.
- Use provider-qualified IDs such as `openrouter/z-ai/glm-5.1`; bare or ambiguous IDs can route through the wrong provider and fail on missing credentials.
- Record exact model/provider/harness, date, prompt class, timeout, retries, and fallback.
- Prefer read-only permissions for reviewers.
- Do not count hung, empty, generic, or off-task output.
- Replace failed reviewers when the review plane is required; otherwise mark the plane `cannot_verify` or `not_assessed`.
- Never treat benchmark claims as evidence that a review finding is correct.
</requirements>
</model_policy>
