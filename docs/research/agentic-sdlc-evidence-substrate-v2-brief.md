# Agentic SDLC Evidence Substrate Brief v2

Status: discussion draft; not committed
Date: 2026-05-05
Inputs:
- `docs/research/agentic-sdlc-evidence-substrate-brief.md`
- `docs/research/harness-telemetry-reviews/round5-socratic-consolidation.md`

This brief is a product framing artifact. It is not source-bound proof,
not product closure evidence, and not a trusted release claim.

## 1. Product Question

The buyer question is:

> Is our AI-assisted delivery control plane gaining or losing control,
> and can we prove which parts of the evidence chain support that answer?

`sdp-trace` does not issue the final degradation verdict. It records
signed observations, explicit gaps, trace chains, retention state, and
verifier outputs that `sdp-gate`, `sdp-report`, or another policy
consumer can use.

Correct product statement:

> `sdp-trace` is a portable evidence substrate for AI-assisted delivery.
> It records what was observed, who observed it, what evidence was bound
> to it, what was missing, and what trust scope the verifier can honestly
> support.

Rejected product statement:

> `sdp-trace` proves that an arbitrary coding agent behaved honestly.

## 2. Product Layers

| Layer | Owns | Must Not Own |
| --- | --- | --- |
| `sdp-trace` | Capture contracts, evidence, provenance, trace, observer identity, missing telemetry, retention, verifier states, query/explain output. | Business pass/fail, readiness, degradation verdicts, waiver policy, opaque health scores. |
| `sdp-gate` | Expected evidence policy, required trust scope, blocking verdicts, waiver rules. | Raw capture and trace storage. |
| `sdp-report` | Time-series aggregation, control posture, CTO/CISO/forensics views. | Cryptographic verification source of truth. |

The first product slice must be usable without `sdp-gate` or
`sdp-report`: `sdp-trace verify`, `sdp-trace query`, and
`sdp-trace explain` must produce a complete single-run view.

## 3. V0 Capture Boundary

V0 uses a command wrapper plus append-only async local event log.

Command shape:

```text
sdp-trace run --task <task-ref> --contract <contract-file> -- <command...>
```

What this captures without harness changes:

- recorder attachment;
- run start/close;
- task/spec/plan digest;
- expected evidence contract digest;
- command argv digest;
- command start/finish time;
- cwd digest or safe path label;
- exit state;
- stdout/stderr retention descriptor;
- source snapshot or git commit/tree digest before and after;
- explicit local gaps for harness/model/gateway/tool-call observers.

What this does not capture in V0:

- prompt bodies;
- response bodies;
- model identity unless provided by adapter/gateway evidence;
- internal tool calls that never cross the wrapper boundary;
- direct file edits made outside the wrapped process unless VCS diff
  evidence is collected after the run.

V0 truth boundary:

> The local wrapper gives `local_observed + partial` evidence. It is
> useful for reconstruction and gap reporting. It is not audit-grade and
> cannot close a gate without CI/OIDC or another external witness.

Bypass behavior:

- If work happens outside the wrapper, there may be no run.
- If a wrapped command reaches CI with missing required observer roles,
  verifier emits `missing_telemetry`.
- If the wrapper starts after work began, verifier emits a
  `late_attach_boundary` and marks pre-attach history `not_assessed`.

V0 local performance contract:

- local writes are append-only and async;
- no network call in the local inner loop;
- p99 local event overhead target: <= 5 ms per event excluding command
  execution time;
- local failure degrades to `cannot_verify` or `partial`, not silent
  success.

## 4. Adapter Interface

Harness, gateway, VCS, and CI integrations emit the same canonical event
shape. Adapters are optional in V0, but their absence must be visible.

Minimum event taxonomy:

| Event | Required Observer Role | Required Fields |
| --- | --- | --- |
| `recorder_attached` | `local_recorder` | run id, recorder version, process id, cwd label/digest, attach mode, local time, nonce |
| `run_started` | `local_recorder` or `harness_adapter` | run id, task ref, contract ref, source snapshot, actor, parent run refs |
| `task_locked` | `human`, `harness_adapter`, or `local_recorder` | task ref, task digest, lock time, author/approver when known |
| `expected_evidence_contract_locked` | `local_recorder`, `ci`, or `human` | contract id, digest, version, policy owner, lock time |
| `harness_identity_observed` | `harness_adapter` | harness id, version, adapter id, attestation state |
| `model_identity_observed` | `gateway` or `harness_adapter` | provider, model, request id/digest, observation source |
| `tool_call_observed` | `harness_adapter` | tool name, call id, args digest, parent span/run id |
| `command_started` | `local_recorder` or `ci` | command id, argv digest, cwd label/digest, start time |
| `command_finished` | `local_recorder` or `ci` | command id, exit state, finish time, stdout/stderr retention |
| `file_mutation_observed` | `local_recorder`, `vcs`, or `ci` | path labels, before/after digests or diff digest, attribution state |
| `test_observed` | `local_recorder` or `ci` | command id, framework, result state, artifact digest/ref |
| `redaction_applied` | `local_recorder`, `ci`, or `gateway` | field class, rule id, redaction manifest digest, pre-write status |
| `policy_override_requested` | `human` | actor, reason, scope, expiry, linked missing evidence, evidence snapshot |
| `requirement_superseded` | `human`, `harness_adapter`, or `local_recorder` | old digest, new digest, reason, author/approver |
| `checkpoint_signed` | `local_recorder`, `ci`, or `external_witness` | sequence range, chain head, signer identity, signing profile |
| `run_closed` | `local_recorder` or `ci` | terminal sequence, chain head, closure state, missing evidence summary |

Each event must carry:

- schema version;
- run id;
- event id;
- sequence;
- event time;
- observer id and role;
- trust scope;
- payload digest;
- previous event hash;
- event hash;
- correlation ids;
- retention descriptor;
- redaction state.

Authority rule:

> A signer can only support trust for event types it is authorized to
> observe. Unauthorized signatures are verification failures, not weaker
> successes.

## 5. Expected Evidence Contract

Missing telemetry is meaningful only when expected evidence is locked
before the run or gate.

Contract shape:

```json
{
  "contract_id": "agent-pr-basic-v1",
  "version": "1.0.0",
  "lock_required_before": "run_started",
  "required_observers": ["local_recorder", "vcs", "ci"],
  "optional_observers": ["harness_adapter", "gateway"],
  "required_events": [
    "recorder_attached",
    "task_locked",
    "expected_evidence_contract_locked",
    "command_started",
    "command_finished",
    "run_closed"
  ],
  "gate_required_events": ["test_observed", "checkpoint_signed"],
  "minimum_trust_scope_for_gate": "ci_witnessed",
  "retention_profile": "digest_safe_default_v1",
  "redaction_profile": "prewrite_digest_default_v1"
}
```

Contract rules:

- Contract digest is recorded in `expected_evidence_contract_locked`.
- Contract changes after lock require `requirement_superseded` or a
  contract supersession event.
- A later weaker contract cannot silently reclassify an older run.
- `sdp-trace` records contract satisfaction facts; `sdp-gate` decides
  whether gaps block.

Verifier must produce a `MissingEvidenceTable`:

| Field | Meaning |
| --- | --- |
| `expected_role` | Observer role required by contract. |
| `expected_event` | Event type required by contract. |
| `observed_state` | `present`, `missing`, `late`, `unauthorized`, `redacted`, `not_assessed`. |
| `sequence_window` | Where the event should have appeared. |
| `reason` | Concrete verifier reason. |
| `policy_ref` | Contract or external policy reference. |

## 6. Privacy, Redaction, And Retention

Default profile: `prewrite_digest_default_v1`.

Defaults:

- no raw prompts;
- no raw responses;
- no raw source snippets;
- no raw stdout/stderr;
- no raw secrets or environment values;
- argv stored as digest plus allowlisted command basename;
- file paths stored as safe labels or digests unless path capture is
  explicitly enabled;
- raw capture is opt-in, path-scoped, and retention-scoped.

Redaction happens before persistence. A secret that reaches the event log
is a product failure, not a retention mode.

Every evidence reference includes:

```json
{
  "retention_mode": "digest_only | sanitized_excerpt | encrypted_raw_ref | external_artifact_ref | not_assessed",
  "replayability": "full | partial | none",
  "expiry_time": "ISO-8601 or null",
  "archival_ref": "URI or null",
  "redaction_manifest_digest": "sha256 or null"
}
```

`redaction_applied` events are part of the trace and must name the rule
id and redaction manifest digest. Aggressive redaction can preserve
privacy while downgrading replayability.

## 7. Signing, Keys, And Witness Scope

The product signs observations and checkpoints, not vague truth claims.

Event chain:

```text
canonical event -> event_hash -> prev_event_hash chain -> checkpoint
```

Checkpoint statement:

```text
in-toto Statement predicate: sdp.trace.checkpoint.v1
DSSE envelope
Sigstore/OIDC, customer PKI, local dev key, or external witness profile
```

Trust scopes:

| Scope | Meaning |
| --- | --- |
| `agent_reported` | Agent self-report only. Never trust upgrade. |
| `local_observed` | Local wrapper/daemon saw an event. Useful but host-forgeable. |
| `harness_observed` | Harness adapter observed lifecycle/tool/model intent. Requires adapter identity. |
| `gateway_observed` | Gateway/provider observed model call metadata. Does not prove local actions. |
| `vcs_observed` | VCS observed commit/diff/PR metadata. |
| `ci_witnessed` | CI verified source, chain, contract, and artifacts, then signed verifier result. |
| `external_witnessed` | External timestamp/transparency/append-only witness anchored chain head. |
| `human_signed` | Human decision, waiver, or approval. |

Key rules:

- Agent-accessible keys cannot sign `local_observed` or stronger scopes.
- Local dev keys can support local structural checks only.
- Gate-grade trust requires CI/OIDC, external witness, or accepted
  customer PKI/private equivalent.
- CI must sign verifier result, source commit, chain head, expected
  contract digest, verifier version, policy profile, timestamp, and CI
  identity.
- If the coding agent runs inside the same CI job as the verifier, the
  verifier must model the topology explicitly and downgrade any event
  whose signer is not independent of the actor being assessed.

Replay/post-hoc protection requires at least one of:

- checkpoint witnessed before gate;
- CI source/time binding;
- external timestamp/transparency receipt;
- customer-approved private equivalent.

Without that, a structurally valid local chain can still be post-hoc.

## 8. Verifier Output

Verifier output is four-axis:

```text
verdict:       pass | fail | cannot_verify | not_assessed
trust_scope:   agent_reported | local_observed | harness_observed |
               gateway_observed | vcs_observed | ci_witnessed |
               external_witnessed | human_signed
completeness:  complete | partial | missing_telemetry | unknown
replayability: full | partial | none | not_assessed
```

Examples:

- `pass + local_observed + partial + partial` is useful but not
  gate-grade.
- `cannot_verify + ci_witnessed + missing_telemetry + partial` is a
  valuable signed gap.
- `fail + local_observed + unknown + none` indicates structural or
  integrity failure.

Verifier must never emit a naked pass without trust scope, completeness,
and replayability.

## 9. Query And Explain Surface

V0 must expose query/explain output, not only JSON.

Required queries:

- `run-summary`: task, contract, source, observer roles, trust scope.
- `timeline`: ordered events with sequence, observer, digest, and gaps.
- `missing-evidence`: MissingEvidenceTable.
- `commands`: command starts/finishes, exit state, retention.
- `files`: VCS/file mutation evidence and attribution state.
- `tests`: test claims vs test evidence.
- `redactions`: redaction rules, manifest digests, replayability impact.
- `witness`: checkpoints, signer identity, witness profile, freshness.
- `overrides`: policy override requests, human signer, scope, expiry.

`sdp-trace explain` must translate failures into developer-readable
causes:

- missing event;
- late attach;
- unauthorized signer;
- chain break;
- payload digest mismatch;
- contract mismatch;
- retention/replayability insufficient;
- redaction unresolved;
- witness missing before gate.

## 10. Developer Workflow

V0 commands:

```text
sdp-trace dry-run --contract <contract-file> -- <command...>
sdp-trace run --task <task-ref> --contract <contract-file> -- <command...>
sdp-trace verify <run-dir>
sdp-trace query <run-dir> --query missing-evidence
sdp-trace explain <run-dir>
```

Emergency path:

- emergency work is not hidden;
- `policy_override_requested` records actor, reason, scope, expiry, and
  linked missing evidence;
- verifier emits `human_signed + partial`, not `complete`;
- `sdp-gate` decides whether the override is acceptable.

Offline path:

- local capture works offline;
- CI/external witness remains `not_assessed` until available;
- delayed CI witness can counter-sign only what it can verify against
  source, chain head, contract, and artifacts.

## 11. Demo Strategy

Demo 0: no-harness-change local observation.

- Run any harness through `sdp-trace run`.
- Output: `pass + local_observed + partial + partial`.
- Show missing harness/model/gateway rows.

Demo 1: dry run and redaction.

- Show what would be captured.
- Include argv with a fake secret.
- Prove secret is redacted before persistence.

Demo 2: CI-witnessed gate.

- CI verifies source, chain head, expected contract, test artifact, and
  verifier version.
- CI signs verifier result.
- Output can become `ci_witnessed + complete` only if contract is met.

Demo 3: missing observer.

- Omit harness/model telemetry.
- Show MissingEvidenceTable.
- Gate usability remains false unless policy accepts the gap.

Demo 4: tamper and chain break.

- Mutate payload, delete event, reorder event.
- Show event-level `explain`.

Demo 5: post-hoc fabrication/replay.

- Generate a valid-looking local trace after work or replay old trace.
- Show downgrade/failure because no pre-gate witness, VCS mismatch,
  timestamp mismatch, nonce mismatch, or external anchor absence.

Demo 6: emergency override.

- Emit `policy_override_requested`.
- Show `human_signed + partial`, not hidden pass.

Demo 7: forensic query.

- Answer what ran, what changed, which tests were evidence, what was
  redacted, and what cannot be assessed.

Demo 8: latency and trace bloat.

- Run high-frequency local events.
- Show p99 overhead and storage limit behavior.

## 12. Known Non-Goals For V0

- No full multi-harness SDK.
- No raw prompt/response capture.
- No opaque health score.
- No native degradation verdict.
- No claim that local-only traces are audit-grade.
- No claim that gateway telemetry proves local file/test behavior.
- No claim that schema validity implies production trust.

## 13. Socratic Review Task

Review this brief adversarially through your assigned persona.

Return:

```text
VERDICT: ACCEPTABLE_WITH_GAPS | CHANGES_REQUIRED | REJECTED

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes/no
- Remaining blockers before implementation:

Critical blockers:
- ...

Major gaps:
- ...

False assumptions:
- ...

Minimum viable correction:
- ...

Questions before implementation:
- ...

Demo changes required:
- ...
```

Tie every finding to at least one area:

- product layering;
- CTO usefulness;
- v0 capture boundary;
- adapter interface;
- evidence model;
- provenance model;
- trace model;
- observer authority;
- signing and verification;
- expected evidence contracts;
- CI/gate anchoring;
- privacy and retention;
- adoption and DX;
- forensic query/replay;
- demo credibility.

Do not reward architecture that weakens UX, DX, or trust clarity. The
target is an 8/10 brief that is good enough to build, not a perfect
general trust architecture.
