# sdp-trace

Portable traceability, provenance, evidence, accountability, and contract-integrity substrate for AI-assisted delivery.

## What It Answers

Can we prove what happened, where the evidence came from, what changed over time, who is accountable, and whether the contract used to assess it was trusted?

Current proof status: Block 01 has contract scaffolding and validation fixtures. `sdp-trace` is not ready for customer pilot claims until the repository validates its own self-trace and self-attestation artifacts.

## Who It Is For

- CTOs who need inspectable movement data rather than opaque quality claims.
- Team leads who need a shared evidence contract without replacing their current coding harness.
- Tool builders integrating agents, SpecKit-style workflows, CI, and PR review.

## What This Repo Contains

- JSON schemas for evidence, provenance, observations, metric movement, accountability, contract manifests, assessment inputs, external verdict inputs, and legacy compatibility artifacts.
- CTO and team lead docs in English and Russian.
- SpecKit terminology mapping.
- Harness integration examples.
- Stack examples for Go and JVM/Bazel.

## What It Does Not Do

`sdp-trace` does not run a delivery workflow, write code, replace code review, manage tickets, decide pass/fail, decide readiness, decide degradation, or require SDP Operator Mode.

It is the trust substrate that other tools can use.

It must prove itself first. If this repository cannot trace its own development under its own contracts, it must not ask a customer to trust the substrate.

## Relationship To sdp-gate

`sdp-gate` uses `sdp-trace` contracts to make readiness decisions for PRs, merge requests, release handoffs, and local gates.

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
- [Agent entrypoint (current verifier contract)](docs/agent-entrypoint.md)
- [Reviewer entrypoint (current proof scope)](docs/reviewer-entrypoint.md)
- [Adoption ladder](docs/adoption-ladder.md)
- [Core concepts](docs/concepts.md)
- [Process metric catalog](docs/process-metric-catalog.md)
- [SpecKit compatibility](docs/speckit-compatibility.md)
- [Schema reference](schema/README.md)

## Minimal Flow

```text
idea -> spec -> plan -> task -> change -> evidence -> provenance -> accountability -> metric movement -> assessment input
```

External policy consumers can turn an assessment input into a gate verdict. `sdp-trace` records those verdicts only as external inputs. When evidence is missing, use `not_assessed`.
