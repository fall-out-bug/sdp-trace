# Customer Pilot Evidence Package Outline

Status: outline; not a completed pilot result
Spec task: T037

## Objective

Create a safe package shape for a customer pilot that records evidence, provenance, trace links, residual `not_assessed` gaps, and optional external policy inputs without committing raw customer data.

## Scope

Included pilot slices:

- OpenCode + MiniMax
- OpenCode + Kimi
- OpenCode + GLM
- Superpowers-style harness pattern
- `gsd`
- `gsd2`
- Oh My OpenAgent
- Kotlin+Bazel scoped service target

## Private Customer Inputs

Private customer input artifacts are never committed to this repository.

| Input | Handling rule | Committed substitute |
|---|---|---|
| Source code | Stays in customer environment. | Redacted path, source commit or customer-approved immutable reference, and file hash when approved. |
| Prompts | Stays in customer environment unless approved for release. | Prompt SHA-256 and redaction note. |
| Tool logs | Sanitized before export. | Summary, omitted-field list, and hash of approved sanitized artifact. |
| Build/test output | Sanitized before export. | Command summary, status origin, and redaction note. |
| Human approvals | Customer approval system remains authoritative. | Access-neutral approval reference. |

## Package Shape

```text
customer-pilot/
  README.md
  scope.md
  redaction-note.md
  run-cards/
    opencode-model-run-card.md
    harness-run-card.md
    kotlin-bazel-fixture-plan.md
  evidence/
    evidence-events.json
    provenance-records.json
    trace-snapshot.json
    observations.json
    metric-stream.json
    proof-states.json
    export-limitations.md
  matrices/
    harness-compatibility-matrix.md
    model-compatibility.md
  handoff/
    assessment-input.json
    external-verdict-input.json
```

The directory shape is illustrative. Committed examples in this repository must use sanitized summaries or placeholders only.

## Outputs Produced By `sdp-trace`

- evidence events or evidence bundles
- provenance records
- observations
- trace snapshot
- assessment input for downstream policy consumers
- matrix evidence state updates
- residual `not_assessed` report

`sdp-trace` does not approve the pilot, decide readiness, or decide support. If `sdp-gate` or another external policy system emits a verdict, record it as an external verdict input with producer, origin, policy reference when available, and artifact reference.

## Review Checkpoints

| Checkpoint | Required evidence |
|---|---|
| Scope lock | Customer-approved scope, source reference, and redaction constraints. |
| Run preflight | Selected run-card, operator identity, expected artifacts, and safety constraints. |
| Evidence packaging | Sanitized artifacts, hashes, omitted-field list, and `not_assessed` reasons. |
| Matrix update | Evidence state, reason code, artifact reference, gap reason, and next required evidence. |
| Policy handoff | Optional assessment input and external verdict inputs. |

## Validation Commands

```bash
jq empty customer-pilot/evidence/*.json
if [ -f customer-pilot/handoff/assessment-input.json ]; then
  jq empty customer-pilot/handoff/assessment-input.json
fi
if [ -f customer-pilot/handoff/external-verdict-input.json ]; then
  jq empty customer-pilot/handoff/external-verdict-input.json
fi
go test ./...
jq empty schema/*.json
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
git diff --check
```

The current Go command surface does not provide a general schema validator for
arbitrary customer pilot package directories. If a pilot package is not
represented by an implemented verifier profile or fixture validator, record
package-level schema validation as `not_assessed` with reason
`validator_not_implemented`; do not replace that gap with retired Node/script
validation.

## Residual `not_assessed` Reporting

The package must list every missing or withheld evidence item with:

- evidence field
- reason code
- operator or accountable owner
- next required evidence
- whether the gap blocks observed behavior evidence
- whether an external policy consumer needs the field
