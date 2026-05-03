# Harness Evidence Matrix

This legacy-named matrix records harness evidence state. It does not claim harness support, readiness, or compatibility as a native `sdp-trace` outcome.

Use `observed` only when a committed sanitized run artifact or evidence summary exists. Planned rows remain `not_assessed`.

| Target | Scope | Evidence state | Reason code | Artifact reference | External verdict reference | Gap reason | Next required evidence |
|---|---|---|---|---|---|---|---|
| Superpowers-style harness pattern | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for the harness pattern. | Run `docs/research/harness-run-card.md` and commit sanitized evidence/provenance/export-limitations summary. |
| gsd | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for `gsd`. | Run `docs/research/harness-run-card.md` and commit sanitized evidence/provenance/export-limitations summary. |
| gsd2 | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for `gsd2`. | Run `docs/research/harness-run-card.md` and commit sanitized evidence/provenance/export-limitations summary. |
| Oh My OpenAgent | Rules, tool logs, hooks, evidence export | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run artifact for Oh My OpenAgent. | Run `docs/research/harness-run-card.md` and commit sanitized evidence/provenance/export-limitations summary. |
| Hyperpowers | Rules, tool logs, hooks, evidence export | `not_assessed` | `discovery_required` | `none` | `none` | Integration surface has not been inspected enough to define run evidence. | Record discovery notes, then add a run-card row before assessment. |
| Paperclip | Rules, tool logs, hooks, evidence export | `not_assessed` | `discovery_required` | `none` | `none` | Integration surface has not been inspected enough to define run evidence. | Record discovery notes, then add a run-card row before assessment. |
| Claude Code | Harness/model split | `not_assessed` | `discovery_required` | `none` | `none` | Harness behavior and model behavior are not separated in committed evidence. | Define harness evidence dimensions and commit sanitized run summary. |
| Codex | Harness/model split | `not_assessed` | `discovery_required` | `none` | `none` | Harness behavior and model behavior are not separated in committed evidence. | Define harness evidence dimensions and commit sanitized run summary. |
| OpenCode + MiniMax | Harness plus exact MiniMax/Kotlin+Bazel slice | `observed` | `run_artifact_available` | `examples/pilot-runs/opencode-minimax-kotlin-bazel` | `none` | Evidence is limited to `minimax-coding-plan/MiniMax-M2.5` on the Block 06 Kotlin+Bazel fixture; Kimi, GLM, and other OpenCode slices remain `not_assessed`. | Run `docs/research/opencode-model-run-card.md` for additional OpenCode model slices before expanding this row. |
| Kilo | Rules, tool logs, hooks, evidence export | `not_assessed` | `discovery_required` | `none` | `none` | Integration surface has not been inspected enough to define run evidence. | Record discovery notes, then add a run-card row before assessment. |
| Pi | Review-only agent role | `not_assessed` | `discovery_required` | `none` | `none` | Review-agent output is recorded, but harness integration evidence is not defined. | Define whether Pi is a harness target or external reviewer before assessment. |
