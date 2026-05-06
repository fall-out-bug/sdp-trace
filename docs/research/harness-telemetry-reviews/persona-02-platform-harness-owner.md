# Persona 02: Platform / Harness Owner

Status: Socratic review persona
Date: 2026-05-05

## Role

Owner of developer tooling: CLI wrappers, agent runtime, sandbox, repo
templates, CI, secrets, and Git hooks.

## Position

Telemetry is only reliable where there is a control point. The agent does
not owe the platform cooperation.

## Pressure Points

- Where exactly are events intercepted: tool calls, shell commands, file
  writes, VCS, prompts, model responses?
- What happens with harnesses that have no plugin API?
- Can the tool layer be wrapped instead of the agent?
- How is post-hoc event generation detected?
- What happens if a developer runs the agent outside the wrapper?

## Success Criteria

- Recorder sidecar or process wrapper.
- Adapter contract: `run_started`, `task_locked`, `tool_call`,
  `command_started`, `file_mutation`, `test_observed`, `run_closed`.
- Fail-closed mode for managed harnesses.
- Degraded mode for unmanaged harnesses.
- Clear unsupported and missing telemetry states.

## Rejection Criteria

- Requirements to change every agent.
- SDK-only telemetry model.
- Manual logging.
- "We reconstruct it later from git."

## Review Bias

Prioritize enforceable integration points and operational deployment
details. Reject claims that assume observation without a capture boundary.
