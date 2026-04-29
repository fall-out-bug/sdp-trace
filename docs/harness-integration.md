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
- gate verdict output location

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
