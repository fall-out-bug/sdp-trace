# Core Concepts

`sdp-trace` is intentionally small. It defines the trust contract around
delivery, not the delivery workflow itself.

The practical question is always the same: what can be replayed, who produced
it, what is missing, and who owns the next decision?

## Spec

The intended behavior and acceptance criteria for a change.

The spec may come from SpecKit, an issue, a product brief, a ticket, or a
human-written note. Tool-specific mappings belong in compatibility and
integration docs, not in the core vocabulary.

## Plan

The intended execution approach.

Plans are useful evidence because they establish what work was expected before code changed.

## Task

An executable unit of work.

It may map to a ticket, checklist item, or workflow step. In `sdp-trace`, it is
just a task.

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

A policy check applied to evidence by CI, release governance, customer
governance, or another external policy consumer.

Examples:

- required tests pass
- changed files match scope
- security-sensitive files have review
- generated code is identified
- docs were updated when public behavior changed

`sdp-trace` may record gate facts, but the gate owner remains external.

## External Verdict

The externally produced gate or policy outcome. These are policy-consumer
concepts, not verifier result states. The canonical verifier result states
(`observed`, `pass`, `fail`, `not_assessed`, `cannot_verify`) and their
exit-code mappings are defined in `docs/agent-entrypoint.md`.

- `pass`: evidence satisfies the gate
- `warn`: evidence exists but risk remains
- `fail`: evidence proves the gate is not satisfied
- `not_assessed`: the gate was outside the selected scope or intentionally not
  evaluated in this run
- `cannot_verify`: the selected gate could not be verified because required
  evidence, environment, freshness, or consistency was missing

Missing required evidence for a selected gate is `cannot_verify` or `fail`, not
`pass`. Missing optional or out-of-scope evidence is `not_assessed`.

`warn` is an External Verdict sub-state, not a verifier result state. When an
External Verdict is `warn`, the underlying verifier result is typically
`observed` or `pass` with an advisory note.

## Assessment Input

A package of evidence, observations, metric streams, accountability,
`not_assessed` gaps, and external verdict inputs prepared for a policy consumer.

`sdp-trace` owns assessment input structure. It does not decide whether the
package passes a policy.

## Accountability

Human-held ownership for the next step.

AI actors may produce, review, critique, or judge artifacts, but they cannot be the sole accountable owner, approver, risk owner, or escalation owner.

## Contract Manifest

A versioned list of schema, docs, validation script, fixture, source commit, approval, and compatibility-note digests.

Schema-valid files are not automatically trusted contract releases. Trusted release status requires manifest digest verification and signature or approved private-equivalent verification.

## Decision Record

The final external decision and rationale.

Decision records are for humans or policy consumers. They explain why a change was accepted, blocked, returned for work, or overridden.
