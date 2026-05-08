# Block 06 Design: OpenCode + MiniMax + Kotlin+Bazel E2E Product Proof

Status: accepted for implementation; spec pi-review findings closed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-drq`
Audience: technical executive, pilot operator, future external policy consumer consumer, repository observer

## Purpose

Block 06 changes the product proof from contract scaffolding to a real executable value path.

The first proven value is:

> A pilot operator can run an OpenCode + MiniMax assessment against a scoped Kotlin+Bazel target and receive a validated `sdp-trace` evidence package showing what happened, where the evidence came from, what remains `not_assessed`, and what downstream policy consumers can inspect.

This block must not repeat the Block 05 mistake of treating run-cards as proof. A run-card is only a recipe. Product proof requires a real run or an explicit `not_assessed` result explaining why the proof could not be completed.

## Product Boundary

`sdp-trace` may provide a reference runner or adapter that shells out to external tools. It must not make OpenCode, MiniMax, Kotlin, Bazel, Bazelisk, or any provider SDK a repository dependency.

The repository may record:

- the command used
- the external tool identity and version when available
- the model identifier requested and observed
- source reference and scoped target evidence
- sanitized output summaries and hashes
- validation result
- residual `not_assessed` gaps

The repository must not commit credentials, raw customer code, raw prompts that are not approved, raw provider responses with proprietary content, or private logs. Raw run output stays under ignored local run directories until sanitized.

`sdp-trace` still does not decide pass/fail, readiness, compatibility, support, or degradation. If the run produces a downstream verdict, it is recorded as an external verdict input.

## Target Slice

The first target slice is fixed:

| Dimension | Required value |
|---|---|
| Harness | OpenCode |
| Model family | MiniMax |
| Stack | Kotlin |
| Build system | Bazel |
| Scope | One service or fixture Bazel target, not whole-repo inference unless the repository is intentionally tiny |
| Output | Validated `sdp-trace` evidence package and tested-on report |

## In Scope

- A Block 06 spec, Socratic dialogue, implementation plan, pi-review ledger, and Beads mirror.
- A reference runner script that:
  - accepts an existing Kotlin+Bazel repository or fixture path
  - accepts an explicit Bazel target tied to the assessed scope
  - accepts an explicit operator-approved Bazel command for full proof
  - checks OpenCode availability and version
  - checks MiniMax model listing and access as separate proof states
  - checks scoped Kotlin+Bazel target evidence before the model run
  - runs OpenCode with a deterministic prompt against the scoped target
  - captures raw output locally under an ignored directory verified before execution
  - emits a sanitized evidence package
- A validation script for the generated proof package.
- A committed sanitized proof package only after a real OpenCode + MiniMax run exists.
- A tested-on report that states the external environment used for the proof.
- Matrix updates for the exact observed row only after committed sanitized run evidence exists.
- Explicit `not_assessed` proof states when MiniMax credentials, OpenCode export, Bazel execution, or source access are missing.

## Out of Scope

- Adding OpenCode, MiniMax, Kotlin, Bazel, Bazelisk, provider SDKs, or customer harnesses to `package.json`.
- Claiming all MiniMax models, all OpenCode providers, all Kotlin projects, or all Bazel workspaces are supported.
- Running or committing raw customer data.
- Making `sdp-trace` a harness runtime.
- Making `sdp-trace` a policy engine.
- Updating every Block 05 matrix row from one observed slice.

## Required Command UX

The runner should support this shape:

```bash
scripts/run-opencode-minimax-kotlin-bazel-proof.sh \
  --repo /path/to/kotlin-bazel-repo \
  --scope services/example \
  --bazel-target //services/example:unit_test \
  --bazel-command "bazel test //services/example:unit_test" \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --out .sdp-trace-runs/opencode-minimax-kotlin-bazel/manual-run
```

The exact MiniMax model id must remain configurable because OpenCode installations expose different provider ids. On this workstation, `opencode models` shows MiniMax candidates including:

- `opencode/minimax-m2.5-free`
- `minimax-coding-plan/MiniMax-M2`
- `minimax-coding-plan/MiniMax-M2.5`
- `openrouter/minimax/minimax-m2.5`

The runner must print clear status lines for each proof state and must exit non-zero when the requested proof cannot be completed. If it emits a partial package, the package must state that the product proof is not complete.

Raw output path rules:

- `--out` must be under `.sdp-trace-runs/` or another path proven ignored by `git check-ignore`.
- The runner must reject tracked output paths before invoking OpenCode.
- Committed examples must never contain `raw/` output.

Bazel command rules:

- Full Block 06 proof requires `--bazel-command`.
- The runner must not execute model-suggested build/test commands automatically.
- If no operator-approved Bazel command is provided, the package may record candidate commands, but `bazel_commands_executed` remains `not_assessed` and Block 06 remains incomplete.
- The command must target the supplied `--bazel-target` or a narrower target inside the same scope.

OpenCode safety rules:

- The runner must not use `--dangerously-skip-permissions`.
- The runner must set a bounded execution timeout.
- The prompt must instruct OpenCode not to edit files.
- If the assessed repository is a git repository, the runner must record `git status --porcelain` before and after the run and fail the proof if the run mutates files.
- Where OpenCode permission configuration is available, the runner must use proof-specific permissions that avoid edits and bound command execution.

## Proof States

Block 06 uses explicit proof states stored in `evidence/proof-states.json`.

Allowed proof-state values:

- `observed`: the required fact was directly observed and has evidence references.
- `not_observed`: the check ran and the required fact was not observed; this is an observed fact, not a policy failure.
- `not_assessed`: the check did not run or evidence was unavailable; the reason and next required evidence are mandatory.

Each proof state must include `name`, `state`, `evidence_refs`, `reason`, and `next_required_evidence`.

| Proof state | `observed` requires |
|---|---|
| `opencode_available` | `opencode --version` succeeds and version is recorded. |
| `minimax_model_listed` | `opencode models` lists the requested MiniMax model id and the sanitized command output digest is recorded. |
| `minimax_access_verified` | A successful OpenCode run with the requested MiniMax model or authenticated provider evidence is recorded. |
| `kotlin_bazel_target_identified` | A Bazel package/target tied to `--scope` and Kotlin source or Kotlin Bazel rule evidence tied to that target are recorded. |
| `opencode_minimax_run_completed` | `opencode run --model <minimax-id>` exits successfully and produces captured output or export evidence. |
| `bazel_commands_executed` | Build/test command evidence exists from `bazel` or `bazelisk`; absence remains `not_assessed`, not failed. |
| `sdp_trace_package_valid` | The generated evidence, provenance, trace, metric, and assessment artifacts validate. |
| `sanitized_report_committed` | A sanitized report with hashes and redaction notes is committed under an examples or docs path. |

Full Block 06 completion requires every proof state above to be `observed`. If `bazel_commands_executed` remains `not_assessed`, the package may prove OpenCode + MiniMax evidence capture over Kotlin+Bazel target evidence, but it does not complete the OpenCode + MiniMax + Kotlin+Bazel E2E product proof.

Harness/model/build success remains observed evidence or external verdict input, not a native readiness decision.

## Evidence Package Shape

The generated proof package should use this shape:

```text
opencode-minimax-kotlin-bazel/
  README.md
  run-report.md
  redaction-note.md
  evidence/
    proof-states.json
    evidence-events.json
    provenance-records.json
    observations.json
    metric-stream.json
    trace-snapshot.json
  handoff/
    assessment-input.json
  raw/
    ignored-local-output-not-committed
```

Committed packages must omit `raw/`. Raw paths may be referenced by digest or local-only note, but raw contents are not committed.

## Required Prompt

The first prompt must be narrow and evidence-oriented:

```text
Inspect the scoped Kotlin+Bazel target at <scope>.
Do not make readiness, support, pass/fail, or compatibility claims.
Report:
1. Kotlin evidence found or missing.
2. Bazel ownership evidence found or missing.
3. The supplied Bazel target and whether the supplied Bazel command evidence was observed.
4. Files inspected.
5. Claims that are not backed by inspected files or command output.
Return a concise summary suitable for a sanitized evidence record.
```

The prompt can be extended for provider-specific CLI requirements, but it must preserve the no-verdict boundary.

## Acceptance

Block 06 spec is accepted when:

- It states that the first proven value is a real OpenCode + MiniMax + Kotlin+Bazel E2E run, not packaging-only.
- It keeps OpenCode, MiniMax, Kotlin, and Bazel as external tested-on tools, not repository dependencies.
- It defines proof states that separate model run, stack detection, Bazel execution, package validation, and committed report evidence.
- It requires machine-readable proof states in `evidence/proof-states.json`.
- It requires a non-zero outcome or explicit incomplete proof state when the real run cannot complete.
- It requires pi review before implementation and finding closure before code changes.

Block 06 implementation is complete only when:

- The reference runner exists and is documented.
- The validation script exists and is wired into repository validation only for committed examples, not for unavailable local external tools.
- A committed sanitized report proves at least one real OpenCode + MiniMax + Kotlin+Bazel run with assessed Bazel command execution, or the Block remains incomplete.
- Matrices move only the exact observed slice and only to an evidence state supported by the committed report.
- Pi review findings for implementation are registered, fixed, and closed.

## Delivery State

Current implementation state:

1. OpenCode is locally installed and reports version `1.14.31`.
2. MiniMax model id `minimax-coding-plan/MiniMax-M2.5` is visible through `opencode models` and completed a runner invocation.
3. `bazel` is locally installed and reports version `9.1.0 Homebrew`.
4. `kotlin` and `kotlinc` are locally installed and report version `2.3.21`.
5. The committed proof package exists under `examples/pilot-runs/opencode-minimax-kotlin-bazel/`.
6. The committed proof package is complete for the exact OpenCode + MiniMax + Kotlin+Bazel fixture slice only.
