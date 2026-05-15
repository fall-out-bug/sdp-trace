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

The adapter-capture path uses portable adapter events. Harness adapters should
emit generic event families rather than product-specific concepts:

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

Missing adapter support is not a failure by itself. The verifier result is
`not_assessed` or `cannot_verify` with a concrete reason. Telemetry labels
such as `missing_telemetry` and integration labels such as `unsupported` or
`not_integrated` describe the nature of the gap but are not verifier result
states. See `docs/agent-entrypoint.md` for the canonical state contract. Agent-reported and harness-observed test
claims never become executed test evidence without CI or registered wrapper/tool
execution proof.

Adapter artifacts must stay safe to commit. Do not persist raw prompts, model
responses, command args, stdout/stderr bodies, tool input/output bodies, adapter
configuration, gateway secrets, provider tokens, authenticated URLs, or raw
review bodies.

## Evidence State

Do not maintain static compatibility matrices for harnesses, models, languages,
or build tools. They collapse different questions into one table and invite
false support claims.

Record evidence at the run or package level instead:

- [Agent Entrypoint](agent-entrypoint.md) defines the current command and state
  contract.
- [Reviewer Entrypoint](reviewer-entrypoint.md) defines how to inspect verifier
  output without overclaiming it.
- [OpenCode + MiniMax + Kotlin/Bazel proof package](../examples/pilot-runs/opencode-minimax-kotlin-bazel/README.md)
  is one exact observed slice, not general support for OpenCode, MiniMax,
  Kotlin, or Bazel.
- [JVM And Bazel Guide](jvm-bazel-guide.md) documents the current JVM/Bazel
  fixture boundary.

Validation must measure tool-use reliability, structured output discipline,
context handling, retained evidence, redaction behavior, and evidence-grounded
claims for the exact observed run. Broader support, readiness, or compatibility
language requires a separate downstream verdict; it is not a native
`sdp-trace` claim.
