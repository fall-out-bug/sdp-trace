# Flight Recorder

The flight recorder is an evidence sidecar for agentic development runs. It
records ordered events, verifier-visible gaps, and reviewer queries. It does
not decide pass/fail, readiness, accountability, or production trust; downstream
policy consumers own those decisions.

## Recorder Profiles

| Profile | What It Can Support | Limit |
| --- | --- | --- |
| `local_development_recorder` | Local development reconstruction and event-chain consistency. | Local-only chains are mutable by the recorded environment and must not support accountability, audit-grade, or external-trust claims. Witness-related states may be `not_assessed`. |
| `witnessed_run_recorder` | Event-chain consistency plus agreement with a witness record outside the run artifact. | Witness records prove the witness contract, not production trust. Missing or mismatched witness evidence is `fail` or `cannot_verify`, depending on access. |
| `externally_witnessed_run` | Future production-capable accountability profile. | Not implemented unless a real external witness boundary is supplied and verified. |

`forensic_retention` is an assessment profile selected with
`sdp-trace assess --profile forensic-retention`, not a recorder profile and not
a retention mode.

## Hash Model

Recorder hashes are structural evidence, not authority by themselves.

- Each event declares schema version, hash algorithm, canonicalization method, `prev_event_hash`, payload digest, and `event_hash`.
- Canonicalization must be deterministic and must exclude fields whose values are derived from the canonical payload being hashed.
- The verifier must detect event mutation, deletion, and reordering.
- A run closure or witness record binds the run id, source baseline hash, task/expectation hash, recorder version, and chain head.

Changing the event log and recomputing a new local chain is still possible in local mode. That is why witness agreement is the boundary between local reconstruction and stronger evidence.

## Gaps And Supersession

Late attachment is a visible provenance boundary. Events before recorder attachment are `not_assessed` with a reason; profiles must not infer pre-attachment provenance from later events.

Requirement, task, prompt, or expectation changes are supersession events. They preserve the original locked event, link to the replacement, and make the timeline queryable. Rewriting the original task event is tampering, not supersession.

## Redaction And Retention

Recorder artifacts must stay safe to commit. Evidence retention states are:

- `digest_only`
- `sanitized_excerpt`
- `encrypted_raw_ref`
- `external_artifact_ref`
- `not_assessed`

The current forensic-retention contract keeps those retention modes stable.
Recording policy profiles such as `safe_default`, `reviewable_sanitized`,
`encrypted_raw_reference`, and `external_forensic_reference` map event families
to those modes; they are not additional modes.

Redaction states are verifier-visible:

- safe redaction applied before persistence
- encrypted raw evidence exists behind an access-controlled reference
- unresolved redaction
- unverifiable redaction

Forensic retention assessment is explicit:

```bash
sdp-trace assess --profile forensic-retention \
  --run <run-dir> \
  --redaction-policy <policy.json> \
  --out <assessment-result.json>

sdp-trace assess explain \
  --assessment-result <assessment-result.json>
```

The profile evaluates verifier facts only. It must not decide legal,
incident, merge, release, readiness, degradation, or risk-acceptance outcomes.

Unresolved redaction must not be converted into a pass. Profiles that require forensic or accountability evidence must fail unresolved redaction, or emit `cannot_verify` / `not_assessed` with a concrete reason when required evidence is unavailable.

## Adapter Capture

`adapter_capture` is an assessment profile selected with
`sdp-trace assess --profile adapter-capture`. It evaluates whether generic
adapter events are present, bound to the selected run, safe to inspect, and
honest about capture-depth limits.

```bash
sdp-trace assess --profile adapter-capture \
  --run <run-dir> \
  --out <assessment-result.json>

sdp-trace query --query capture-depth <run-dir>
```

Adapter capture supports two binding modes:

- `same_chain`: adapter events are in the recorder chain with hash linkage.
- `adapter_bundle`: adapter events are in a separate digest-bound bundle
  referenced by the run before closure.

Adapter capture does not authorize an adapter, prove managed enrollment, or
upgrade model identity to gateway-observed. It reports `missing_telemetry`,
`unsupported`, `not_integrated`, `not_assessed`, `retention_limited`, and
`cannot_verify` as explicit facts. Query output is a read-only summary and does
not emit a top-level pass/fail policy verdict.

Adapter event artifacts inherit the forensic-retention safety rules. Raw
prompts, raw model
responses, tool input/output bodies, raw command args, stdout/stderr bodies,
adapter configuration, gateway evidence refs, provider tokens, and raw review
bodies must not appear in committed artifacts or preview/query/explain output.

## Query Surface

The current Go-first validation command for committed fixture packages is:

```bash
sdp-trace validate-fixtures examples/agentic-sdlc
```

The query names below describe the product surface expected from flight-recorder
summaries. They are not Node.js command names and must not add Node.js to the
active product path.

```bash
sdp-trace query --query missing-evidence <run-dir>
sdp-trace query --query capture-depth <run-dir>
sdp-trace query-pack --pack forensics-basic-v1 \
  --run <run-dir> \
  --out <query-pack-result.json>
```

Required query names:

| Query | Reviewer Answer |
| --- | --- |
| `run-summary` | Source baseline, task lock, recorder mode, chain state, closure state. |
| `provenance` | Requested and observed model/harness identity, actor, and identity gaps. |
| `gaps` | Late-attach and `not_assessed` boundaries with reasons. |
| `requirement-timeline` | Original requirement events and superseding events. |
| `command-timeline` | Commands, redaction state, exit state, and linked evidence. |
| `file-mutations` | Files changed, digest evidence, and mutation event order. |
| `test-evidence` | Test command evidence and verifier state, without policy verdicts. |
| `redactions` | Redaction and retention issues requiring reviewer attention. |
| `witness-state` | Witness presence, witness agreement, and witness scope. |
| `capture-depth` | Adapter event coverage, missing telemetry, unsupported observers, gateway `not_integrated`, and retention/capture caps. |
| `forensics-basic-v1` query pack | Safe grouped rows for run summary, chain and witness state, redaction and retention issues, capture-depth gaps, timelines, task supersession, unverified claims, and evidence-gap descriptions. |

Queries expose evidence and verifier states only. They must not emit opaque scores or policy verdicts.

Forensics query packs are read-only bundles over existing facts. They must
preserve `not_assessed`, `cannot_verify`, `retention_limited`, `unsupported`,
and `not_integrated` rows rather than hiding them behind a summary. Digest-only
evidence is existence evidence, not reconstructable evidence. Forensics query
packs use opaque command/test identifiers and path-redacted artifact identifiers;
raw command names, executable paths, script paths, test names, and artifact
storage paths are not safe output.
