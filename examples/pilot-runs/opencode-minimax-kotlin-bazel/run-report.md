# OpenCode + MiniMax + Kotlin+Bazel Run Report

Completion state: complete

## Tested-On Environment

| Field | Value |
|---|---|
| Repository | examples/pilot-fixtures/kotlin-bazel-service |
| Source ref | working-tree-scope:e7afea8291ec62f74543045fb39549d6f786d227eff9a732907ee3245618c9ae |
| Source content sha256 | e7afea8291ec62f74543045fb39549d6f786d227eff9a732907ee3245618c9ae |
| Scope | services/example |
| Bazel target | //services/example:compile_hello_jar |
| Bazel command | bazel build --symlink_prefix=/ //services/example:compile_hello_jar |
| Model | minimax-coding-plan/MiniMax-M2.5 |
| OpenCode version | 1.14.31 |
| Bazel version | bazel 9.1.0 Homebrew |
| Kotlin version | Kotlin version 2.3.21-release-298 (JRE 17.0.19+0) |
| kotlinc version | info: kotlinc-jvm 2.3.21 (JRE 17.0.19+0) |

## Source Artifacts

| Path | sha256 |
|---|---|
| MODULE.bazel | b0d7c78060dc50d9aeba514e968bf84ff2ceb2362aff0eb6890ec8d48321c016 |
| services/example/BUILD.bazel | a1a3ea84601a5897ab6b382eac3e18bda6a0a797690130cbf30000d5af8190a8 |
| services/example/Hello.kt | 9aff584dc92a33627445c05cbd7c92cd822292170090a2a81092213a6fb70670 |

## Command Results

| Command | Started | Ended | Exit code | stdout sha256 | stderr sha256 |
|---|---|---|---|---|---|
| OpenCode run | 2026-05-01T14:29:15Z | 2026-05-01T14:29:43Z | 0 | c8c43b4716c509096074bc94cd8e985c7d3c0ccc70ebcc655533547f68c21711 | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| Bazel command | 2026-05-01T14:29:43Z | 2026-05-01T14:29:44Z | 0 | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 | 2ea0839a6e384a346efa4aa6b37b6e658fa83766b63517c845c18a464e553e76 |

## Proof States

| Proof state | State | Reason |
|---|---|---|
| opencode_available | observed | opencode --version succeeded |
| minimax_model_listed | observed | requested MiniMax model id appears in opencode models output |
| minimax_access_verified | observed | successful opencode run verifies access to requested model id |
| kotlin_bazel_target_identified | observed | bazel query succeeded and target rule output ties Kotlin/Bazel files to the supplied target |
| opencode_minimax_run_completed | observed | opencode run completed with requested MiniMax model id and produced captured output |
| bazel_commands_executed | observed | operator-approved bazel command completed |
| sdp_trace_package_valid | observed | Package validation passed for committed proof package |
| sanitized_report_committed | observed | Sanitized report committed to repository |

## Boundary

sdp-trace records observed evidence for this exact slice only. It does not produce a pass/fail, readiness, support, compatibility, or degradation verdict.
