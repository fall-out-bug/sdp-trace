# Core Concepts

`sdp-trace` is intentionally small. It defines the trust contract around delivery, not the delivery workflow itself.

## Spec

The intended behavior and acceptance criteria for a change.

The spec may come from GitHub SpecKit, an issue, a product brief, a ticket, or a human-written note.

## Plan

The intended execution approach.

Plans are useful evidence because they establish what work was expected before code changed.

## Task

An executable unit of work.

In full SDP this may map to a Beads issue or workstream. In `sdp-trace`, it is just a task.

## Change

The concrete code, docs, config, schema, or infrastructure modification.

A change should link back to at least one task or spec.

## Provenance

The origin chain behind the change.

Examples:

- human author
- coding agent
- model family and model version when available
- tool calls and commands
- CI runner
- review actor

Provenance does not imply quality. It only records origin.

## Evidence

Proof that can be inspected.

Examples:

- tests
- CI runs
- static analysis
- review comments
- command output
- changed files
- screenshots
- manual sign-off

Evidence should be referenced, not paraphrased into vague claims.

## Gate

A check applied to evidence.

Examples:

- required tests pass
- changed files match scope
- security-sensitive files have review
- generated code is identified
- docs were updated when public behavior changed

## Verdict

The gate outcome:

- `pass`: evidence satisfies the gate
- `warn`: evidence exists but risk remains
- `fail`: evidence proves the gate is not satisfied
- `not_assessed`: evidence is missing or insufficient

## Decision Record

The final decision and rationale.

Decision records are for humans. They explain why a change was accepted, blocked, returned for work, or overridden.
