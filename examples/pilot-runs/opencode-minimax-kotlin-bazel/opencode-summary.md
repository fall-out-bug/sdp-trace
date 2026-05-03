# Sanitized OpenCode Output Summary

Raw OpenCode output is not committed. This summary records the inspectable content needed for the Block 06 proof without private logs.

## Model And Run

- Harness: OpenCode `1.14.31`
- Model id: `minimax-coding-plan/MiniMax-M2.5`
- Raw stdout SHA-256: `c8c43b4716c509096074bc94cd8e985c7d3c0ccc70ebcc655533547f68c21711`
- Raw stderr SHA-256: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

## Files Inspected By OpenCode

- `services/example/Hello.kt`
- `services/example/BUILD.bazel`
- `MODULE.bazel`
- `MODULE.bazel.lock`

## Sanitized Assessment Content

- Kotlin evidence was found in `services/example/Hello.kt`.
- Bazel ownership evidence was found in `services/example/BUILD.bazel` and `MODULE.bazel`.
- The target `//services/example:compile_hello_jar` was identified as a `genrule` with `srcs = ["Hello.kt"]`.
- The target command invokes host `kotlinc` through the genrule command. This is observed local fixture behavior, not Bazel-managed Kotlin toolchain proof.

## Unbacked-Claim Handling

- The model stated that no Bazel command output was supplied to inspect. That statement is scoped to the model prompt context only.
- The committed proof package does not rely on the model for Bazel command execution evidence. Bazel command evidence is recorded separately in `evidence/proof-states.json` and `run-report.md`.
- No OpenCode-generated pass/fail, readiness, support, compatibility, or degradation verdict is promoted into `sdp-trace`.
