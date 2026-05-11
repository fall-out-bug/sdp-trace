# Product Contract v0 Traceability Matrix

This matrix maps current `sdp-trace` substrate to Change Evidence Packet v0
rows. It is intentionally conservative: a current capability may feed a packet
row without proving the full row end to end.

## Packet Row Coverage

| packet row | current substrate inputs | current coverage | remaining gap |
| --- | --- | --- | --- |
| `PC-CHANGE` | `schema/delivery-trace-envelope.schema.json`, `internal/repoobserver`, `examples/block28-repo-observer/` | Partial | Need packet-level change identity contract for local Git, GitHub, GitFlic, GitLab, and artifact-only flows. |
| `PC-INITIATOR` | `schema/decision-record.schema.json`, trace/provenance examples, harness session metadata | Partial | Need source task/issue/prompt-boundary binding row and safe redaction rule. |
| `PC-AGENT-ROUTE` | `internal/harnessobs`, `examples/harness-observation/`, `schema/adapter-event.schema.json`, `schema/adapter-capture-run.schema.json` | Partial | Need route projection from raw substrate into buyer-readable chain. |
| `PC-MUTATION` | `internal/harnessobs`, `internal/repoobserver`, authority observed actions, Git evidence examples | Partial | Need packet row that separates mutation existence from actor/tool/model attribution. |
| `PC-VERIFICATION` | `schema/evidence-event.schema.json`, `schema/evidence-bundle.schema.json`, `internal/harnessobs`, `internal/ciartifact`, block17 and block26 examples | Partial | Need packet row that separates agent-claimed, harness-observed, CI-witnessed, and missing verification. |
| `PC-REVIEW` | `internal/prreview`, `examples/pr-review/`, reviewer entrypoint docs | Partial | Need packet row for reviewer plane, independence state, retained result, and absent review state. |
| `PC-AUTHORITY` | `schema/authority-envelope.schema.json`, `schema/authority-evaluation.schema.json`, `internal/authority`, `docs/authority-envelope.md` | Partial substrate | Need packet projection that says authority state without making policy decisions. |
| `PC-THEATER` | Existing negative fixtures, managed harness failures, adapter capture overclaim failures, review ledger discipline | Partial | Need first-class theater reason-code rows and derivation rules. |
| `PC-ATTESTATION` | `schema/witness-profile-result.schema.json`, `internal/witness`, `internal/releaseproof`, `docs/contract-release-signing.md`, block15/16 examples | Partial substrate | Need additive packet profile language and private/customer witness baseline. |
| `PC-DECISION` | `schema/accountability.schema.json`, `docs/accountability-model.md`, decision records | Partial | Need buyer-facing owner row that does not become approval or blame. |
| `PC-RESIDUAL-GAPS` | `not_assessed`/`cannot_verify` patterns across schemas and examples | Partial substrate | Need packet-level gap summary grouped by decision relevance. |

## Existing Work Reclassification

| current area | product role under Contract v0 | not product progress until |
| --- | --- | --- |
| Flight recorder / trace substrate | Evidence source for multiple rows | It renders into packet rows. |
| Evidence bundle schemas | Attachment and evidence-row source | Packet references the bundle and exposes missing states. |
| Harness observation | Agent route, mutation, verification evidence source | It fills `PC-AGENT-ROUTE`, `PC-MUTATION`, or `PC-VERIFICATION`. |
| Authority envelope | Authority row source | It is projected into `PC-AUTHORITY` without policy verdicts. |
| CI artifact observation | Verification/witness source | It fills `PC-VERIFICATION` with retained CI witness refs. |
| PR review profiles | Review row source | It fills `PC-REVIEW` with plane, result, and independence state. |
| Witness and release proof | Attestation row source | It fills `PC-ATTESTATION` as additive evidence only. |
| Adapter capture | Route/source capture | It fills route or mutation rows without broad support claims. |
| Query packs | Investigation aid | A packet uses or links query results. |
| Dashboard/report UI | Projection | Packet semantics are already stable. |

## P0 Classification Template

Every P0 candidate must include this block:

```text
packet_rows:
  - PC-...
evidence_surface:
  - <artifact, API, schema, fixture, or command output>
start_state:
  <pass|partial|fail|not_assessed|cannot_verify|missing_telemetry|unsupported|not_integrated>
target_transition:
  <what changes in the row after the work>
buyer_effect:
  <packet row, section, or theater reason code that becomes clearer>
non_goal:
  <what this still does not prove>
```

If the candidate cannot fill `packet_rows` or cannot show forward progress on a
row, it is not P0 product progress. Repeating `not_assessed -> not_assessed` or
`cannot_verify -> cannot_verify` without a new evidence surface, narrower
claim, or explicit unsupported state does not qualify.

## Known Gaps

| gap | affected area | closure evidence |
| --- | --- | --- |
| No generated packet command | Implementation after approval | Go implementation that renders packet from retained inputs. |
| No packet schema | Implementation after contract approval | `change-evidence-packet.schema.json` or equivalent Go model. |
| No generated local enterprise fixture | Russian baseline confidence | Implementation fixture using local Git plus self-hosted/change-host refs. |
| Theater reason-code derivation is not implemented | `PC-THEATER` automation | Documented rules and tests for P0 theater codes. |
| Decision owner row is not bound | `PC-DECISION` confidence | Policy, task, or change-host owner ref with missing-state handling. |
| Static HTML projection absent | Demo polish | Projection generated from canonical packet without changing semantics. |
