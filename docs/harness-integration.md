# Harness Integration

`sdp-trace` should work with existing harnesses instead of replacing them.

## Integration Contract

A harness integration should provide:

- scope input
- actor identity
- model identity when available
- tool calls or command log
- changed files or diff reference
- evidence references
- assessment input output location
- optional external verdict input location when a policy consumer runs

Block 19 adapter capture adds a portable adapter-event path. Harness adapters
should emit generic event families rather than product-specific concepts:

- `run_started`
- `task_locked`
- `task_superseded`
- `tool_call`
- `command_started`
- `file_mutation`
- `model_call_observed`
- `test_observed`
- `run_closed`

Stable product behavior depends on generic fields such as tool family,
provenance scope, capture state, run binding, source baseline, and
redaction/retention metadata. Adapter-local labels may be retained as sanitized
metadata, but they are not stable product contract members.

Missing adapter support is not a failure by itself. It must be reported as
`missing_telemetry`, `unsupported`, `not_integrated`, `not_assessed`, or
`cannot_verify` with a concrete reason. Agent-reported and harness-observed test
claims never become executed test evidence without CI or registered wrapper/tool
execution proof.

Adapter artifacts must stay safe to commit. Do not persist raw prompts, model
responses, command args, stdout/stderr bodies, tool input/output bodies, adapter
configuration, gateway secrets, provider tokens, authenticated URLs, or raw
review bodies.

## Harness Families To Validate

- Superpowers
- Hyperpowers
- gsd / gsd2
- Oh My OpenAgent
- Paperclip
- Codex
- Claude Code
- OpenCode
- Kilo
- Pi

## Model Families To Validate

- GLM
- MiniMax
- Kimi
- MiMo

Validation must measure tool-use reliability, structured output discipline, context handling, and evidence-grounded claims.
