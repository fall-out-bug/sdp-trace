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
    evidence-bundle.json
    provenance-records.json
    trace-snapshot.json
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
node scripts/validate-json-schema.mjs schema/evidence-bundle.schema.json customer-pilot/evidence/evidence-bundle.json
node scripts/validate-json-schema.mjs schema/trace.schema.json customer-pilot/evidence/trace-snapshot.json
jq -c '.[]' customer-pilot/evidence/provenance-records.json | while read -r record; do
  printf '%s\n' "$record" > /tmp/sdp-trace-provenance-record.json
  node scripts/validate-json-schema.mjs schema/provenance-record.schema.json /tmp/sdp-trace-provenance-record.json
done
if [ -f customer-pilot/handoff/assessment-input.json ]; then
  node scripts/validate-json-schema.mjs schema/assessment-input.schema.json customer-pilot/handoff/assessment-input.json
fi
if [ -f customer-pilot/handoff/external-verdict-input.json ]; then
  node scripts/validate-json-schema.mjs schema/external-verdict-input.schema.json customer-pilot/handoff/external-verdict-input.json
fi
node scripts/validate-pilot-matrices.mjs
```

Use repository-local validation for committed examples:

```bash
npm run validate
git diff --check
```

## Residual `not_assessed` Reporting

The package must list every missing or withheld evidence item with:

- evidence field
- reason code
- operator or accountable owner
- next required evidence
- whether the gap blocks observed behavior evidence
- whether an external policy consumer needs the field
