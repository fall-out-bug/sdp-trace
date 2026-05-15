<model_policy>
Use model diversity to reduce shared blind spots, not to outsource judgment.

<roles>
- Codex/GPT: implementation review, Go-specific integration, local verification synthesis.
- OpenCode: portable skill execution and read-only reviewer harnesses using `.agents/skills`.
- GLM-5.1 or current Z.AI coding model: long-horizon engineering, optimization, and architecture doubt.
- MiniMax M2.7 or current MiniMax coding model: skill adherence, multi-step workflow review, and tool-loop critique.
- Kimi K2.5 or current official Moonshot coding model: wide-context code review and alternate reasoning. Do not pin K2.6 in durable policy without primary-source confirmation.
</roles>

<requirements>
- Record exact model/provider/harness, date, prompt class, timeout, retries, and fallback.
- Prefer read-only permissions for reviewers.
- Do not count hung, empty, generic, or off-task output.
- Replace failed reviewers when the review plane is required; otherwise mark the plane `cannot_verify` or `not_assessed`.
- Never treat benchmark claims as evidence that a review finding is correct.
</requirements>
</model_policy>
