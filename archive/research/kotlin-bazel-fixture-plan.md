# Kotlin+Bazel Fixture Plan

Status: fixture plan; Block 06 fixture behavior observed for one exact slice
Spec task: T029

## Purpose

This plan defines what a real Kotlin+Bazel pilot run must prove. The design fixture alone does not prove Kotlin+Bazel behavior. The Block 06 proof package records one observed OpenCode + MiniMax run against the committed fixture target.

## Evidence Boundary

The committed `examples/jvm-bazel/evidence-bundle.json` is still a design fixture. It may demonstrate the shape of a safe evidence bundle, but it does not prove real Kotlin+Bazel run behavior.

The committed package under `examples/pilot-runs/opencode-minimax-kotlin-bazel/` is the first observed fixture run for this exact target:

- Fixture: `examples/pilot-fixtures/kotlin-bazel-service`
- Scope: `services/example`
- Target: `//services/example:compile_hello_jar`
- Command: `bazel build --symlink_prefix=/ //services/example:compile_hello_jar`
- Tested-on report: `archive/research/opencode-minimax-kotlin-bazel-proof-report.md`

## Required Scope Evidence

| Evidence type | Acceptable evidence | Not enough by itself |
|---|---|---|
| Bazel ownership | `BUILD` or `BUILD.bazel` target for the assessed scope; `MODULE.bazel`, `WORKSPACE`, or `WORKSPACE.bazel` tied to that target. | `.bazelrc` alone or Bazel files in unrelated directories. |
| Kotlin service language | `.kt` or `.kts` source in the assessed scope; `kt_jvm_*` rule; Kotlin compiler/toolchain rule tied to the target. | Kotlin libraries, Kotlin stdlib dependency, or generated metadata alone. |
| Scoped command | `bazel test <target>` or `bazel build <target>` with sanitized output, or explicit reason why the command was unsafe/unavailable. | Root-level monorepo command with no service scope. |
| Dependency metadata | Maven/Gradle files may be recorded as dependency metadata when scoped Bazel evidence proves Bazel ownership. | Inferring Maven/Gradle ownership from metadata inside a Bazel-owned target. |
| Missing command output | `not_assessed` with `unsafe_to_run`, `missing_export`, or `sanitization_pending`. | Treating unavailable output as success. |

## Observed Block 06 Run Artifact

The first observed run commits a sanitized package with:

- `evidence/proof-states.json`
- `evidence/evidence-events.json`
- `provenance-records.json`
- `trace-snapshot.json`
- command summary and command digests for the scoped Bazel target
- redaction note for withheld raw output
- explicit no-verdict boundary

## Placeholder Rule

For any Kotlin+Bazel target other than the exact Block 06 fixture slice:

- matrix state remains `not_assessed`
- reason code is `design_fixture_only` or `no_run_artifact`
- artifact reference for observed behavior is `none`
- no support, readiness, compatibility, or policy result is recorded as a native `sdp-trace` outcome
