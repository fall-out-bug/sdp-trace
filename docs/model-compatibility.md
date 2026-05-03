# Model Evidence Matrix

This legacy-named matrix records model-slice evidence state. It does not claim model support, readiness, or compatibility as a native `sdp-trace` outcome.

## Evaluation Matrix

| Target | Scope | Evidence state | Reason code | Artifact reference | External verdict reference | Gap reason | Next required evidence |
|---|---|---|---|---|---|---|---|
| OpenCode + MiniMax | OpenCode model slice on Block 06 Kotlin+Bazel fixture | `observed` | `run_artifact_available` | `examples/pilot-runs/opencode-minimax-kotlin-bazel` | `none` | Evidence is limited to `minimax-coding-plan/MiniMax-M2.5` on the Block 06 Kotlin+Bazel fixture. | Run additional MiniMax provider/model ids and target stacks before expanding beyond this exact slice. |
| OpenCode + Kimi | OpenCode model slice | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed OpenCode run with observed Kimi identity. | Run `docs/research/opencode-model-run-card.md` and commit sanitized evidence/provenance/export-limitations summary. |
| OpenCode + GLM | OpenCode model slice | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed OpenCode run with observed GLM identity. | Run `docs/research/opencode-model-run-card.md` and commit sanitized evidence/provenance/export-limitations summary. |
| MiMo | Model family discovery | `not_assessed` | `discovery_required` | `none` | `none` | No committed routing path or model identity evidence. | Define harness/model slice before adding run-card assessment. |

## Required Measurements

- context window behavior
- tool-use reliability
- read-before-claim discipline
- structured output validity
- false confidence rate
- `not_assessed` compliance when evidence is missing
