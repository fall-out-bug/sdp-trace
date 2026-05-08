# Harness Evidence Matrix

This legacy-named file records harness evidence state. It does not claim harness
support, readiness, or compatibility as a native `sdp-trace` outcome.

For first-time onboarding, start with `docs/README.md` and
`docs/harness-integration.md`. This matrix is an evidence registry, not a user
journey.

Use `observed` only when a committed sanitized run artifact or evidence summary exists. Planned rows remain `not_assessed`.

| Target | Scope | Evidence state | Reason code | Artifact reference | External verdict reference | Gap reason | Next required evidence |
|---|---|---|---|---|---|---|---|
| Superpowers-style harness pattern | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for the harness pattern. | Commit a sanitized evidence/provenance/export-limitations summary before assessment. |
| gsd | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for `gsd`. | Commit a sanitized evidence/provenance/export-limitations summary before assessment. |
| gsd2 | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for `gsd2`. | Commit a sanitized evidence/provenance/export-limitations summary before assessment. |
| Oh My OpenAgent | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for Oh My OpenAgent. | Commit a sanitized evidence/provenance/export-limitations summary before assessment. |
| Hyperpowers | Rules, tool logs, hooks, evidence export | `not_assessed` | `discovery_required` | `none` | `none` | Integration surface has not been inspected enough to define run evidence. | Record discovery notes, then add a run-card row before assessment. |
| Paperclip | Rules, tool logs, hooks, evidence export | `not_assessed` | `discovery_required` | `none` | `none` | Integration surface has not been inspected enough to define run evidence. | Record discovery notes, then add a run-card row before assessment. |
| Claude Code | Harness/model split | `not_assessed` | `discovery_required` | `none` | `none` | Harness behavior and model behavior are not separated in committed evidence. | Define harness evidence dimensions and commit sanitized run summary. |
| Codex | Harness/model split | `not_assessed` | `discovery_required` | `none` | `none` | Harness behavior and model behavior are not separated in committed evidence. | Define harness evidence dimensions and commit sanitized run summary. |
| OpenCode + MiniMax | Harness plus exact MiniMax/Kotlin+Bazel slice | `observed` | `run_artifact_available` | `examples/pilot-runs/opencode-minimax-kotlin-bazel` | `none` | Evidence is limited to `minimax-coding-plan/MiniMax-M2.5` on the committed Kotlin+Bazel fixture; Kimi, GLM, and other OpenCode slices remain `not_assessed`. | Commit additional sanitized model-slice evidence before expanding this row. |
| Kilo | Rules, tool logs, hooks, evidence export | `not_assessed` | `discovery_required` | `none` | `none` | Integration surface has not been inspected enough to define run evidence. | Record discovery notes, then add a run-card row before assessment. |
| Pi | Review-only agent role | `not_assessed` | `discovery_required` | `none` | `none` | Review-agent output is recorded, but harness integration evidence is not defined. | Define whether Pi is a harness target or external reviewer before assessment. |
