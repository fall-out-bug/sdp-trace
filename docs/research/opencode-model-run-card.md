# OpenCode Model Run-Card

Status: recipe; MiniMax has one observed Block 06 run artifact
Spec task: T027

## Scope

This run-card defines the planned OpenCode model slices for MiniMax, Kimi, and GLM. It is not evidence that every model slice works. MiniMax has one observed artifact for the exact Block 06 Kotlin+Bazel fixture slice; Kimi and GLM remain `not_assessed`.

| Slice | Evidence state | Reason code | Artifact reference | External verdict reference | Gap reason | Next required evidence |
|---|---|---|---|---|---|---|
| OpenCode + MiniMax | `observed` | `run_artifact_available` | `examples/pilot-runs/opencode-minimax-kotlin-bazel` | `none` | Evidence is limited to `minimax-coding-plan/MiniMax-M2.5` on the Block 06 Kotlin+Bazel fixture. | Run additional OpenCode+MiniMax rows before broadening beyond this exact slice. |
| OpenCode + Kimi | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed OpenCode run with observed Kimi identity. | Run this card, commit sanitized evidence bundle, provenance, trace snapshot, and export-limitations note. |
| OpenCode + GLM | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed OpenCode run with observed GLM identity. | Run this card, commit sanitized evidence bundle, provenance, trace snapshot, and export-limitations note. |

## Prompt Template

```text
Assess the scoped change for sdp-trace evidence quality.

Repository: <repo-url-or-local-redacted-ref>
Source commit: <40-character-commit-sha>
Scope: <service-or-directory>
Harness: OpenCode
Model family: <MiniMax|Kimi|GLM>
Model version: <observed-version-or-not_assessed>

Produce:
1. evidence-bundle.json with inspected files, commands, and missing evidence as not_assessed.
2. provenance-records.json with operator, harness identity, model identity, source commit, and tool/command summaries when available.
3. trace-snapshot.json linking spec, task, evidence, provenance, and observations.
4. unbacked_claim list for any statement not backed by inspected evidence.

Do not infer support, readiness, compatibility, or policy outcome.
```

## Block 06 Reference Runner

The committed MiniMax proof package was produced from this command shape:

```bash
scripts/run-opencode-minimax-kotlin-bazel-proof.sh \
  --repo examples/pilot-fixtures/kotlin-bazel-service \
  --scope services/example \
  --bazel-target //services/example:compile_hello_jar \
  --bazel-command "bazel build --symlink_prefix=/ //services/example:compile_hello_jar" \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --out .sdp-trace-runs/opencode-minimax-kotlin-bazel/local-run
```

The runner shells out to externally installed OpenCode and Bazel tools. These tools are tested-on environment inputs, not repository dependencies.

## Required Artifacts

| Artifact | Required content | Commit rule |
|---|---|---|
| `evidence-bundle.json` | Inspected files, command summaries, missing evidence, redaction notes, SHA-256 for committed summaries. | Commit sanitized summary only. |
| `provenance-records.json` | Human operator, OpenCode identity/version when available, model family/version when available, source commit, tool/command summaries. | Commit sanitized summary only. |
| `trace-snapshot.json` | Links from spec/task to evidence/provenance/observations. | Commit sanitized summary only. |
| `assessment-input.json` | Optional downstream policy input. | Only when preparing handoff to `sdp-gate` or another external policy consumer. |
| `export-limitations.md` | Missing logs, missing model version, unavailable tool export, or manual capture notes. | Required when any provenance field is unavailable. |

## Provenance Fields

- `operator_identity`
- `harness_family: OpenCode`
- `harness_version`
- `model_family`
- `model_version`
- `source_commit`
- `prompt_sha256`
- `artifact_redaction_status`
- `tool_log_access`
- `command_summary_access`

Unavailable fields must be recorded as `not_assessed` with a reason code such as `missing_export`, `sanitization_pending`, or `discovery_required`.

## Unbacked Claim Capture

Record an `unbacked_claim` item when the model or harness:

- states a positive product outcome or quality conclusion without inspected evidence
- infers build system ownership from dependency metadata alone
- treats missing logs as success
- reports a model version that cannot be observed in the exported artifacts

External verdict words emitted by another tool may be recorded only as external verdict inputs with producer, origin, policy reference when available, and artifact reference.

## Validation

Run after a sanitized run artifact exists:

```bash
node scripts/validate-json-schema.mjs schema/evidence-bundle.schema.json <run>/evidence-bundle.json
node scripts/validate-json-schema.mjs schema/trace.schema.json <run>/trace-snapshot.json
jq -c '.[]' <run>/provenance-records.json | while read -r record; do
  printf '%s\n' "$record" > /tmp/sdp-trace-provenance-record.json
  node scripts/validate-json-schema.mjs schema/provenance-record.schema.json /tmp/sdp-trace-provenance-record.json
done
node scripts/validate-pilot-matrices.mjs
scripts/validate-e2e-pilot-package.sh examples/pilot-runs/opencode-minimax-kotlin-bazel
```

Do not update `docs/model-compatibility.md` to `observed` unless the committed artifact reference points to the sanitized run evidence.
