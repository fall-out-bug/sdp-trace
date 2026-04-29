# Model Compatibility

`sdp-trace` does not require a specific model, but gate quality depends on model behavior when a model contributes findings.

## Evaluation Matrix

| Model family | Scout | Gate findings | PR review | JSON/schema | JVM/Bazel | Notes |
|---|---|---|---|---|---|---|
| GLM | TBD | TBD | TBD | TBD | TBD | Target family |
| MiniMax | TBD | TBD | TBD | TBD | TBD | Target family |
| Kimi | TBD | TBD | TBD | TBD | TBD | Target family |
| MiMo | TBD | TBD | TBD | TBD | TBD | Target family |

## Required Measurements

- context window behavior
- tool-use reliability
- read-before-claim discipline
- structured output validity
- false confidence rate
- `not_assessed` compliance when evidence is missing
