# Schema Reference

These schemas define the portable `sdp-trace` contract.

## Schemas

| Schema | Purpose |
|---|---|
| `trace.schema.json` | Links specs, tasks, changes, evidence, gates, decisions, and actors. |
| `evidence-bundle.schema.json` | Captures reviewable proof for a scoped change. |
| `gate-verdict.schema.json` | Records one gate result. |
| `decision-record.schema.json` | Records the final human or automated decision. |

## Validation

Basic JSON syntax check:

```bash
jq empty schema/*.json
```

Full JSON Schema validation is intentionally not pinned yet. The first implementation track should choose a validator and freeze schema IDs.
