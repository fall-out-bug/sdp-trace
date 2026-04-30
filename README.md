# sdp-trace

Portable traceability, provenance, evidence, and quality gate contracts for AI-assisted delivery.

## What It Answers

Can we prove that a change was made within a known scope, with known provenance, reviewable evidence, and an explicit quality decision?

## Who It Is For

- CTOs who need confidence that AI-assisted delivery is not degrading quality.
- Team leads who need a shared gate contract without replacing their current coding harness.
- Tool builders integrating agents, SpecKit-style workflows, CI, and PR review.

## What This Repo Contains

- JSON schemas for evidence bundles, traces, gate verdicts, and decision records.
- CTO and team lead docs in English and Russian.
- SpecKit terminology mapping.
- Harness integration examples.
- Stack examples for Go and JVM/Bazel.

## What It Does Not Do

`sdp-trace` does not run a delivery workflow, write code, replace code review, manage tickets, or require SDP Operator Mode.

It is the trust substrate that other tools can use.

## Relationship To sdp-gate

`sdp-gate` uses `sdp-trace` to make readiness decisions for PRs, merge requests, release handoffs, and local gates.

`sdp-trace` does not depend on `sdp-gate`.

```text
sdp-gate -> sdp-trace
sdp-trace -> no SDP runtime
```

## Start Here

- [Current SpecKit feature spec](specs/001-sdp-trace-time-series-evidence-substrate/spec.md)
- [CTO brief, English](docs/cto-brief.en.md)
- [CTO brief, Russian](docs/cto-brief.ru.md)
- [Team lead playbook, English](docs/team-lead-playbook.en.md)
- [Team lead playbook, Russian](docs/team-lead-playbook.ru.md)
- [Adoption ladder](docs/adoption-ladder.md)
- [Core concepts](docs/concepts.md)
- [SpecKit compatibility](docs/speckit-compatibility.md)
- [Schema reference](schema/README.md)

## Minimal Flow

```text
spec -> plan -> task -> change -> evidence -> gate verdict -> decision record
```

The important rule: a verdict may only say what the evidence supports. When evidence is missing, use `not_assessed`.
