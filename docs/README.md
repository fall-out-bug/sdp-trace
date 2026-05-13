# Documentation Map

Use this page as the documentation entrypoint. It separates product docs from
working specs and historical implementation records.

## First-Time Path

1. [Install](install.md): binary-first setup plus source-checkout commands.
2. [Core Concepts](concepts.md): vocabulary and product boundary.
3. [Agent Onboarding](agent-onboarding.md): the single link to give a coding
   agent before it works in this repository.
4. [Agent Entrypoint](agent-entrypoint.md): current command, state, trust-scope,
   authority-scope, and exit-code contract.
5. [Reviewer Entrypoint](reviewer-entrypoint.md): quick verification path and
   overclaim rules.
6. [Harness Integration](harness-integration.md): how existing workflows feed
   trace evidence without being replaced.
7. [Schema Reference](../schema/README.md): the portable JSON contracts.

## Reader Shortcuts

- New contributor: read onboarding, concepts, and the agent entrypoint before
  editing files.
- Reviewer: start with the reviewer entrypoint, then inspect the evidence
  policy and relevant examples.
- Harness maintainer: read harness integration, flight recorder, and schema
  reference.
- Engineering leader: read the repository README, run the reviewer five-minute
  verification path, then inspect the adoption guide, CI check policy, and
  accountability model.

The same product boundary applies to every path: `sdp-trace` structures
evidence and gaps; it does not approve changes.

## Governance And Rollout Docs

- [Adoption Guide, English](adoption-guide.en.md)
- [Adoption Guide, Russian](adoption-guide.ru.md)
- [Repository Rollout Playbook, English](repository-rollout-playbook.en.md)
- [Repository Rollout Playbook, Russian](repository-rollout-playbook.ru.md)
- [Accountability Model](accountability-model.md)
- [Evidence Policy](evidence-policy.md)
- [CI Check Policy](ci-check-policy.md)

## Engineering Docs

- [Harness Integration](harness-integration.md)
- [Claim Authoring](claim-authoring.md)
- [Spec Drift Register](spec-drift-register.md)
- [Contract Release Signing](contract-release-signing.md)
- [Process Metric Catalog](process-metric-catalog.md)
- [Flight Recorder](flight-recorder.md)
- [SpecKit Compatibility](speckit-compatibility.md)
- [JVM And Bazel Guide](jvm-bazel-guide.md)

## Working Artifacts

- `specs/`: active and historical working specs. Current repository records are
  SpecKit-shaped, but the product contract can map evidence from other planning
  flows.
- `examples/`: sanitized fixtures and pilot evidence packages. Treat each
  example according to its README and recorded evidence state.

## Current Product Boundary

`sdp-trace` records evidence and missing evidence. It does not own policy
decisions. Readiness, degradation, merge blocking, release approval, and risk
acceptance belong to CI, customer governance, release management, or another
external policy consumer.
