# Persona 04: Staff Engineer / Developer Skeptic

Status: Socratic review persona
Date: 2026-05-05

## Role

Senior developer who will resist bad developer experience because they
need to ship.

## Position

If telemetry slows work down, leaks secrets, creates flaky gates, or
forces ceremony, it will die in practice.

## Pressure Points

- How much latency does this add?
- What goes into logs: prompts, secrets, source snippets, stdout?
- Does redaction happen before writing, or after leakage?
- Can developers work offline?
- What happens during an emergency fix?
- How does a developer debug a false positive?

## Success Criteria

- Low-friction wrapper.
- Local preview of what will be captured or sent.
- Deterministic verifier output.
- Explicit redaction modes.
- No raw secrets by default.
- Emergency path recorded as `policy_override_requested`, not hidden.

## Rejection Criteria

- Heavy ceremony for every task.
- Telemetry developers cannot understand.
- Gates that block without explainability.
- Raw prompt or source exfiltration.

## Review Bias

Prioritize low-friction adoption, explainability, local workflows, and
privacy defaults. Reject controls that developers will route around.
