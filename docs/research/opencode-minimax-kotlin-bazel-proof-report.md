# OpenCode + MiniMax + Kotlin+Bazel Proof Report

Status: incomplete package with observed subset for the exact Block 06 slice
Spec task: T083
Proof package: `examples/pilot-runs/opencode-minimax-kotlin-bazel/`

## Tested Slice

| Dimension | Observed value |
|---|---|
| Harness | OpenCode `1.14.31` |
| Model | `minimax-coding-plan/MiniMax-M2.5` |
| Stack | Kotlin `2.3.21` |
| Build system | Bazel `9.1.0 Homebrew` |
| Fixture | `examples/pilot-fixtures/kotlin-bazel-service` |
| Scope | `services/example` |
| Bazel target | `//services/example:compile_hello_jar` |
| Bazel command | `bazel build --symlink_prefix=/ //services/example:compile_hello_jar` |

## Artifact References

- Proof states: `examples/pilot-runs/opencode-minimax-kotlin-bazel/evidence/proof-states.json`
- Evidence events: `examples/pilot-runs/opencode-minimax-kotlin-bazel/evidence/evidence-events.json`
- Provenance: `examples/pilot-runs/opencode-minimax-kotlin-bazel/evidence/provenance-records.json`
- Trace snapshot: `examples/pilot-runs/opencode-minimax-kotlin-bazel/evidence/trace-snapshot.json`
- Assessment input: `examples/pilot-runs/opencode-minimax-kotlin-bazel/handoff/assessment-input.json`
- Human report: `examples/pilot-runs/opencode-minimax-kotlin-bazel/run-report.md`

## Proof State Summary

The committed package records `completion_state: incomplete`. These proof states are currently `observed`:

- `opencode_available`
- `minimax_model_listed`
- `minimax_access_verified`
- `kotlin_bazel_target_identified`
- `opencode_minimax_run_completed`
- `bazel_commands_executed`

These required proof states remain `not_assessed` in `proof-states.json`:

- `sdp_trace_package_valid`
- `sanitized_report_committed`

## Boundary

This is evidence for one tested slice only. It does not claim all OpenCode providers, all MiniMax ids, all Kotlin projects, or all Bazel workspaces are covered.

`sdp-trace` records the observed evidence and handoff shape. It does not produce a pass/fail, readiness, support, compatibility, or degradation verdict.
