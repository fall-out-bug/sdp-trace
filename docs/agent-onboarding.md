# Agent Onboarding

Use this page as the link to give a coding agent before it works in this
repository.

## Product Boundary

`sdp-trace` is a portable evidence recorder for AI-assisted delivery. It records
what happened, which evidence exists, what is missing, where the evidence came
from, and which human owns the next decision.

It does not replace the team's planning system, harness, agent, CI, review
process, or release governance. It also does not decide whether a change may be
merged, released, accepted, or overridden.

Origin note: `sdp-trace` was extracted from delivery evidence work in
`sdp_lab`. That history is not a dependency and should not be treated as
important runtime context.

## Workflow Mapping

`sdp-trace` uses a small shared vocabulary:

```text
spec -> plan -> task -> change -> evidence -> provenance -> accountability -> assessment input
```

The source workflow may be SpecKit, gsd, Superpowers, Oh My OpenAgent, a ticket
tracker, a repository template, or a custom delivery process. Map the workflow
into the shared vocabulary, then record evidence. Do not claim that the source
workflow itself is supported unless the current run produced evidence for that
exact integration.

Useful mappings:

| Source workflow concept | sdp-trace concept |
| --- | --- |
| Feature brief, issue, GSD goal, SpecKit spec | `spec` |
| Delivery plan, implementation plan, runbook | `plan` |
| Ticket, task, bead, checklist item | `task` |
| Diff, docs change, config change, generated artifact | `change` |
| Tests, CI run, review, command output, retained artifact | `evidence` |
| Human, agent, model, tool, command, source chain | `provenance` |
| DRI, approver, risk owner, escalation owner | `accountability` |
| Evidence package for CI, governance, or another policy consumer | `assessment input` |

## First Commands

Follow the [Contributor Quick Start](contributor-quickstart.md) for the current
canonical Go-first smoke path. It includes environment checks, failure routing,
and expected verifier states.

Use these docs while working:

- [Core Concepts](concepts.md): vocabulary and product boundary.
- [Agent Entrypoint](agent-entrypoint.md): command, state, trust-scope,
  authority-scope, and exit-code contract.
- [Reviewer Entrypoint](reviewer-entrypoint.md): quick verification path and
  overclaim checklist.
- [Harness Integration](harness-integration.md): adapter and workflow evidence
  expectations.
- [Schema Reference](../schema/README.md): portable JSON contracts.

## Evidence Rules

- A local run is local evidence only.
- A CI witness is CI evidence only.
- A customer-PKI witness is customer-authority evidence only for the stated
  payload and policy.
- A checked-in JSON file is an audit artifact until the current verifier replays
  it or an accepted external signature binds it.
- Missing evidence remains `not_assessed`, `cannot_verify`,
  `missing_telemetry`, `not_integrated`, or a concrete failure reason.

Never convert missing evidence into success. Never describe a local verifier
fact as production trust.

## What To Produce

For a repository change, preserve enough information for another reviewer to
replay the evidence boundary:

- the source spec, plan, and task identifiers;
- the command or adapter path used to observe the work;
- the retained run directory or report package;
- test, CI, review, and artifact references;
- any `not_assessed` or `cannot_verify` gaps;
- the human-held DRI, approver, risk owner, and escalation path when the change
  affects trust claims.

If the workflow is not observable by `wrap`, `run`, or an adapter, say that
directly and record the missing boundary. `sdp-trace` is more useful when it is
honest about gaps than when it pretends to see work it did not observe.
