# Documentation Map

Use this page as the documentation entrypoint. It separates product docs from
working specs and historical research artifacts.

## First-Time Path

1. [Core Concepts](concepts.md): vocabulary and product boundary.
2. [Agent Entrypoint](agent-entrypoint.md): current command, state, trust-scope,
   authority-scope, and exit-code contract.
3. [Reviewer Entrypoint](reviewer-entrypoint.md): quick verification path and
   overclaim rules.
4. [Harness Integration](harness-integration.md): how existing workflows feed
   trace evidence without being replaced.
5. [Schema Reference](../schema/README.md): the portable JSON contracts.

## Governance And Rollout Docs

- [One-Minute Brief, English](cto-brief.en.md)
- [One-Minute Brief, Russian](cto-brief.ru.md)
- [Adoption Guide, English](cto-adoption-guide.en.md)
- [Adoption Guide, Russian](cto-adoption-guide.ru.md)
- [Repository Rollout Playbook, English](team-lead-playbook.en.md)
- [Repository Rollout Playbook, Russian](team-lead-playbook.ru.md)
- [Customer Questions, English](customer-questions.en.md)
- [Customer Questions, Russian](customer-questions.ru.md)
- [Accountability Model](accountability-model.md)
- [Evidence Policy](evidence-policy.md)
- [CI Check Policy](ci-check-policy.md)

## Engineering Docs

- [Harness Integration](harness-integration.md)
- [Claim Authoring](claim-authoring.md)
- [Contract Release Signing](contract-release-signing.md)
- [Process Metric Catalog](process-metric-catalog.md)
- [Flight Recorder](flight-recorder.md)
- [SpecKit Compatibility](speckit-compatibility.md)
- [JVM And Bazel Guide](jvm-bazel-guide.md)

## Evidence Matrices

These files are evidence-state registries. They are not broad compatibility
claims.

- [Harness Evidence Matrix](harness-compatibility-matrix.md)
- [Model Evidence Matrix](model-compatibility.md)

## Working And Historical Artifacts

- `specs/`: active and historical SpecKit artifacts. Use these for development
  context, not as a first-time product guide.
- `archive/research/`: research notes, review outputs, pilot reports, and
  development evidence. These files can explain why decisions were made, but
  they are not authoritative proof unless current verifier output or accepted
  external signatures bind them.
- `examples/`: sanitized fixtures and pilot evidence packages. Treat each
  example according to its README and recorded evidence state.

## Current Product Boundary

`sdp-trace` records evidence and missing evidence. It does not own policy
decisions. Readiness, degradation, merge blocking, release approval, and risk
acceptance belong to downstream consumers such as `sdp-gate`, CI, customer
governance, or release management.
