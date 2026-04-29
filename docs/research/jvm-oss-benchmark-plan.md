# JVM OSS Benchmark Plan

Status: draft research protocol
Created: 2026-04-29

This benchmark validates whether `sdp-trace` can capture useful traceability, provenance, evidence, and gate inputs across real JVM repositories and multiple agent harnesses.

## Research Questions

1. Can each harness/model combination build a correct repo map without hallucinating build tools or language conventions?
2. Can it produce an evidence bundle that separates inspected facts from model judgment?
3. Can it handle scoped work inside large JVM repositories?
4. Can it emit `not_assessed` when evidence is missing?
5. Which combinations are strong enough to feed `sdp-gate` readiness decisions?

## Selected Projects

| Project | Source | Primary stack | Build system | Why it is in the set |
|---|---|---|---|---|
| Spring Framework | https://github.com/spring-projects/spring-framework | Java with some Kotlin | Gradle | Enterprise Java baseline, large modular framework, active Spring ecosystem. |
| Apache Kafka | https://github.com/apache/kafka | Java + Scala | Gradle | Large distributed system with Java/Scala split and substantial module structure. |
| JetBrains Kotlin | https://github.com/JetBrains/kotlin | Kotlin + Java | Gradle | Kotlin compiler/tooling stress case; tests Kotlin understanding directly. |
| Apache Flink | https://github.com/apache/flink | Java + Scala | Maven | Large Maven multi-module project with streaming/runtime complexity. |
| Bazel | https://github.com/bazelbuild/bazel | Java + Starlark + Python | Bazel | Bazel-native Java project; validates the exact failure mode raised in discovery. |

## Why Not Elasticsearch In Round 1

Elasticsearch is valuable, but licensing and product-boundary questions make it less clean as the first OSS benchmark. Keep it as a round-2 candidate after the protocol is stable.

## Harness And Model Matrix

| ID | Harness | Model family | Role |
|---|---|---|---|
| C1 | Codex | default GPT-backed runtime | baseline coding-agent behavior |
| CC1 | Claude Code | GLM through configured gateway | Claude harness with non-Claude model stress |
| O1 | OpenCode | GPT | strong hosted baseline |
| O2 | OpenCode | GLM | target sovereign/local-style model |
| O3 | OpenCode | Kimi | target long-context model |
| O4 | OpenCode | MiniMax | target long-context model |
| P1 | Pi | GPT | Pi baseline |
| P2 | Pi | GLM | target sovereign/local-style model |
| P3 | Pi | Kimi | target long-context model |
| P4 | Pi | MiniMax | target long-context model |

MiMo is a target model family, but it is not in the initial runnable matrix until the local routing path is available.

## Experiment Phases

### Phase A: Read-Only Repo Intake

Run every harness/model combination against every project.

Expected artifacts:

- repo map
- detected languages
- detected build system
- scoped-module candidates
- evidence bundle with inspected files
- `not_assessed` list for unavailable checks

Scale: 5 projects x 10 combinations = 50 runs.

### Phase B: Gate Input Synthesis

For each project, select two existing PRs or commits and ask each combination to produce gate inputs.

Expected artifacts:

- trace graph
- evidence bundle
- gate verdict draft
- unsupported-claim list

Scale: 5 projects x 2 changes x 10 combinations = 100 runs.

### Phase C: Bounded Local Patch

For each project, define one tiny local-only patch. Run only the best three combinations from Phase A/B.

Expected artifacts:

- task plan
- changed files
- command log
- evidence bundle
- gate verdict
- reviewer notes

Scale: 5 projects x 3 combinations = 15 runs.

## Scoring Rubric

| Dimension | Pass condition |
|---|---|
| Build detection | Identifies Gradle, Maven, or Bazel correctly. |
| Language detection | Does not apply Java-only rules to Kotlin or Scala without a caveat. |
| Evidence discipline | Claims cite files, commands, diffs, or review artifacts. |
| Scope discipline | Does not scan or edit outside requested module without reason. |
| Structured output | Produces schema-compatible JSON or clearly explains why not. |
| Uncertainty handling | Uses `not_assessed` instead of guessing. |
| Token/runtime practicality | Completes within a usable time and context budget. |

## Output Layout

Recommended local output path:

```text
.sdp-trace-runs/
  <project>/
    phase-a/
      <combo-id>/
        repo-map.md
        evidence-bundle.json
        trace.json
        notes.md
    phase-b/
    phase-c/
```

These run outputs should not be committed until sanitized. Research summaries belong in `docs/research/`.

## Bootstrap

Use the local bootstrap runbook before Phase A:

- [bootstrap-runbook.md](bootstrap-runbook.md)

The project manifest lives at `benchmarks/jvm-oss/projects.json`.

## Round-1 Decision

Proceed with these five projects. They cover Gradle, Maven, Bazel, Java, Kotlin, Scala, compiler/tooling, enterprise framework, and distributed runtime use cases.
