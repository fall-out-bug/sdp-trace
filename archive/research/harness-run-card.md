# Harness Run-Card

Status: recipe; no run artifacts committed
Spec task: T028

## Scope

This run-card defines planned harness evidence capture for the Superpowers-style harness pattern, `gsd`, `gsd2`, and Oh My OpenAgent. It does not introduce a dependency on any harness runtime and does not record harness support, readiness, or compatibility.

| Harness target | Evidence state | Reason code | Artifact reference | External verdict reference | Gap reason | Next required evidence |
|---|---|---|---|---|---|---|
| Superpowers-style harness pattern | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run showing rules, tool logs, hooks, or export behavior. | Run a scoped assessment and commit sanitized evidence/provenance/export-limitations summary. |
| `gsd` | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run showing rules, tool logs, hooks, or export behavior. | Run a scoped assessment and commit sanitized evidence/provenance/export-limitations summary. |
| `gsd2` | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run showing rules, tool logs, hooks, or export behavior. | Run a scoped assessment and commit sanitized evidence/provenance/export-limitations summary. |
| Oh My OpenAgent | `not_assessed` | `no_run_artifact` | `none` | `none` | No committed run showing rules, tool logs, hooks, or export behavior. | Run a scoped assessment and commit sanitized evidence/provenance/export-limitations summary. |

## Prompt Template

```text
Assess the scoped change for sdp-trace evidence capture.

Repository: <repo-url-or-local-redacted-ref>
Source commit: <40-character-commit-sha>
Scope: <service-or-directory>
Harness target: <Superpowers-style|gsd|gsd2|Oh My OpenAgent>
Model family/version: <observed-or-not_assessed>

Capture:
1. rules or prompt location
2. tool log access
3. hook availability or absence
4. evidence export shape
5. manual fallback path when export is missing
6. unbacked_claim list

Do not infer model behavior from harness export behavior, and do not infer harness behavior from model output quality.
```

## Evidence Dimensions

| Dimension | Required evidence | `not_assessed` reason code |
|---|---|---|
| Rules or prompt location | File path, command output, or manual evidence note. | `discovery_required` |
| Tool log access | Exported log summary or explicit unavailable note. | `missing_export` |
| Hook availability | Hook configuration, command output, or explicit unavailable note. | `discovery_required` |
| Evidence export | File list and schema/parse status where available. | `missing_export` |
| Manual fallback | Operator note with redaction status and source commit. | `sanitization_pending` |
| Harness identity | Observed harness name/version or explicit unavailable note. | `missing_export` |

## Required Artifacts

- `evidence-bundle.json`
- `provenance-records.json`
- `trace-snapshot.json`
- `export-limitations.md`
- optional external verdict input when a separate policy tool emits one

Every artifact committed to the repository must be sanitized. Raw prompts, tool logs, customer source, customer output, and credentials stay outside the repository.

## Validation

```bash
jq empty <run>/evidence/*.json
if [ -f <run>/handoff/assessment-input.json ]; then
  jq empty <run>/handoff/assessment-input.json
fi
go test ./...
jq empty schema/*.json
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
git diff --check
```

The current Go command surface does not provide a general schema validator for
arbitrary harness package directories. If the package is not represented by an
implemented verifier profile or fixture validator, record package-level schema
validation as `not_assessed` with reason `validator_not_implemented`; do not
replace that gap with retired Node/script validation.

Update `docs/harness-compatibility-matrix.md` only with observed evidence state and artifact references. Planned rows remain `not_assessed`.
