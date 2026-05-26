# Documentation Map

Use this page as the documentation entrypoint. It separates product docs from
working specs and historical implementation records.

## First-Time Path

1. [Install](install.md): binary-first setup plus source-checkout commands.
2. [Core Concepts](concepts.md): vocabulary and product boundary.
3. [Agent Onboarding](agent-onboarding.md): the single link to give a coding
   agent before it works in this repository.
4. [Contributor Quick Start](contributor-quickstart.md): run the canonical
   local smoke path and verify your environment.
5. [Agent Entrypoint](agent-entrypoint.md): current command, state, trust-scope,
   authority-scope, and exit-code contract.
6. [Reviewer Entrypoint](reviewer-entrypoint.md): quick verification path.
7. [Output Location Map](output-location-map.md): where each command writes
   artifacts.
8. [Profile Selection Guide](profile-selection-guide.md): which profile to use.
9. [Overclaim Checklist](overclaim-checklist.md): canonical forbidden-claims
   and trust-scope rules.
10. [Harness Integration](harness-integration.md): how existing workflows feed
    trace evidence without being replaced.
11. [Schema Reference](../schema/README.md): the portable JSON contracts.

## Reader Shortcuts

- New contributor: read the quick start, onboarding, concepts, and the agent
  entrypoint before editing files.
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

- [Spec Roadmap](roadmap.md): current spec statuses, capability ownership, and lifecycle labels.
- [Spec Status Discipline](spec-status-discipline.md): separate spec, task,
  implementation, review, merge, and trust states.
- [Spec Closure Route](spec-closure-route.md): audit of specs 001-019 and the
  route for closing or deferring each one.
- [Adoption Guide, English](adoption-guide.en.md)
- [Adoption Guide, Russian](adoption-guide.ru.md)
- [Production Adoption Readiness](production-adoption-readiness.md): what is known, pilot-capable, and `not_assessed`.
- [Repository Rollout Playbook, English](repository-rollout-playbook.en.md)
- [Repository Rollout Playbook, Russian](repository-rollout-playbook.ru.md)
- [Accountability Model](accountability-model.md)
- [Evidence Policy](evidence-policy.md)
- [Security Baseline](security-baseline.md): local security scan results and triage.
- [Security Policy](../.github/SECURITY.md): vulnerability reporting and scope.
- [CI Check Policy](ci-check-policy.md)

## Engineering Docs

- [Harness Integration](harness-integration.md)
- [Claim Authoring](claim-authoring.md)
- [Command Stability Matrix](command-stability-matrix.md)
- [Package Ownership Map](package-ownership-map.md)
- [Extension Boundary Plan](extension-boundary-plan.md)
- [Spec Drift Register](spec-drift-register.md)
- [Contract Release Signing](contract-release-signing.md)
- [Process Metric Catalog](process-metric-catalog.md)
- [Flight Recorder](flight-recorder.md)
- [SpecKit Compatibility](speckit-compatibility.md)
- [JVM And Bazel Guide](jvm-bazel-guide.md)
- [OSS Replacement Compatibility](oss-replacement-compatibility.md)
- [OSS Benchmark Results](oss-benchmark-results.md)

## Working Artifacts

- `specs/`: active and historical working specs. Current repository records are
  SpecKit-shaped, but the product contract can map evidence from other planning
  flows. Each spec directory may contain its own `reviews/` subdirectory for
  durable review synthesis tied to that block.
- `examples/`: sanitized fixtures and pilot evidence packages. Treat each
  example according to its README and recorded evidence state.

## Review Evidence

Durable review synthesis lives with the spec it assesses:
`specs/<block-name>/reviews/` for block-scoped PI review rounds and synthesis.

Raw local subagent output, temporary worktrees, and PR-scoped working notes
belong outside the tracked tree. The repository hygiene check enforces this by
rejecting tracked `.worktrees/`, `.codex-subagents/runs/`, `.sdp-trace-*`, root
`PR_DESCRIPTION.md`, root `design-note.md`, root `reviews/`, root
executables or binaries, and absolute local paths in durable docs.

## Current Product Boundary

`sdp-trace` records evidence and missing evidence. It does not own policy
decisions. Readiness, degradation, merge blocking, release approval, and risk
acceptance belong to CI, customer governance, release management, or another
external policy consumer.
