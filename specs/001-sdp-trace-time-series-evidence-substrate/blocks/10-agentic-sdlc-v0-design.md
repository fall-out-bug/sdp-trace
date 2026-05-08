# Block 10: Agentic SDLC Evidence Substrate V0

Status: planning draft; implementation not started
Date: 2026-05-05
Inputs:
- retired research artifact
- retired research artifact
- `09-flight-recorder-trust-kernel.md`

This document is a design input for implementation. It is not
source-bound proof, not product closure evidence, and not a trusted
release claim.

## Goal

Build the first usable `sdp-trace` V0 for agentic SDLC evidence capture:
a portable Go CLI that wraps existing AI development workflows, records
observable events, makes missing telemetry explicit, verifies trace
chain continuity and witness scope, and exports a forensic audit bundle.

The external demo target is:

```text
OpenCode + GSD + Bazel + Kotlin
```

The product promise is narrower:

```text
sdp-trace can show what was observed, by which observer, at what trust
scope, what evidence is replayable, and what expected telemetry is
missing.
```

It must not claim that the agent was honest, that the code is correct, or
that the delivery system is healthy. Those decisions belong to
`sdp-gate`, `sdp-report`, or customer policy consumers.

## Product Boundary

V0 is the evidence substrate, not the governance product.

In scope:

- local process wrapper;
- wrapper composition around existing harnesses;
- Unix socket adapter ingress;
- JSON expected evidence contracts;
- pre-write redaction;
- sequence-continuous local event chain with tamper detection under
  honest-recorder assumptions;
- MissingEvidenceTable;
- verifier output with explicit trust scope;
- CI witness verification path;
- audit bundle export;
- A demo repository may consume the generic product CLI for an
  OpenCode/GSD/Bazel/Kotlin scenario, but those harness/build-system
  artifacts are not part of the `sdp-trace` product repository.

Out of scope:

- retroactive attach to already-running agent processes;
- raw prompt or response capture by default;
- org-wide degradation dashboards;
- opaque health scores;
- local-only audit-grade trust;
- proof of deleted local runs before VCS/CI/preflight expected them;
- Windows adapter transport;
- mandatory SDK adoption by every harness.

## Primary UX

The CTO buyer wants a control layer over an existing AI SDLC. The first
UX must therefore be sidecar-first and harness-neutral:

```bash
sdp-trace wrap --name existing-harness -- <existing command...>
```

When there is no existing wrapper:

```bash
sdp-trace run --task <task-ref> --contract <contract.json> -- <command...>
```

Developer preview:

```bash
sdp-trace dry-run --contract <contract.json> -- <command...>
```

Verification and forensics:

```bash
sdp-trace verify <run-dir>
sdp-trace query <run-dir> --query missing-evidence
sdp-trace explain <run-dir>
sdp-trace export <run-dir> --format audit-bundle
sdp-trace ci-verify --run <run-artifact> --contract <contract.json>
```

The wrapper must preserve stdin, stdout, stderr, TTY behavior, colors,
signals, and child exit code. If `sdp-trace` changes the developer's
normal command semantics, the demo fails even if the trace is correct.

## Trust Model

V0 signs observations, not truth.

Local-only traces assume the recorder process was honest while it ran.
They are sequence-continuous under the recorder's local view; they are
not externally anchored integrity proof. A hostile actor that controls
the recorder from startup can fabricate a locally consistent chain.
CI/external witnessing is what turns a local chain head into gate-usable
evidence.

Required verifier axes:

```text
verdict: observed | fail | cannot_verify | not_assessed
trust_scope: agent_reported | local_observed | harness_observed |
             gateway_observed | vcs_observed | ci_witnessed |
             external_witnessed | human_signed
completeness: complete | partial | missing_telemetry | unknown
replayability: full | partial | none | not_assessed
```

Local recorder evidence can show structural continuity of a local run
under honest-recorder assumptions, but it is not gate-grade. Gate-grade
evidence requires an independent CI or external witness that binds:

- run id;
- source commit or tree digest;
- expected evidence contract digest;
- trace chain head;
- verifier version;
- verifier result;
- CI identity and timestamp;
- witness independence state.

The verifier must reject trust upgrades when a signer is unknown,
unauthorized, self-claimed, or controlled by the assessed agent process.

Human signatures require a configured signing profile. Local named keys
without hardware binding or equivalent protection are recorded as
`human_declared` or customer-accepted local risk, not strong
`human_signed` proof of human presence.

CI contract pinning must be outside the assessed agent's writable source
tree. Repo-committed CI YAML alone is not sufficient if the assessed
agent can modify it in the same change. Acceptable V0 gate-grade pinning
profiles are CI secret/config outside the PR diff, external policy
service, human-signed digest, or customer PKI equivalent.

Witness independence states:

- `independent`: verifier/signing boundary is outside the assessed agent
  process tree and writable workspace.
- `same_job`: verifier runs in the same CI job but a separate process;
  disclose topology and downgrade from gate-grade unless customer policy
  accepts it.
- `same_container`: verifier and agent share container namespace; not
  gate-grade in V0.
- `same_process`: verifier and agent share process memory; treat as
  `agent_reported` or `not_assessed`.
- `not_assessed`: independence could not be determined.

Local timestamps use the system clock plus monotonic ordering where the
platform exposes it. They are not trusted time evidence. Gate-grade time
requires CI or external timestamp/witness evidence.

Run ids are UUIDv7 or 128-bit-random-equivalent identifiers plus a
per-run nonce. Predictable ids are invalid.

Schema versions are strict within V0. Older schema versions are accepted
only when explicitly listed in the verifier compatibility table; unknown
future versions produce `cannot_verify`.

## Event Model

Every event is canonical JSON with:

- schema version;
- run id;
- event id;
- sequence;
- event time;
- observer id;
- observer role;
- claimed role;
- verified role state;
- trust scope;
- payload digest;
- previous event hash;
- event hash;
- correlation ids;
- retention descriptor;
- redaction state.

Canonicalization and hashing:

- JSON canonicalization: RFC 8785 / JSON Canonicalization Scheme.
- `event_hash = sha256(canonical_event_without_event_hash)`.
- `previous_event_hash` is `null` only on `recorder_attached`.
- The first event includes `run_nonce` and `run_id`; these are local
  recorder values, not external anchors.
- Signatures, when present, sign the event hash or checkpoint statement,
  not an implementation-specific serialization.
- Changing the hash preimage is a schema-versioned breaking change.

Required V0 events:

- `recorder_attached`;
- `run_started`;
- `task_locked`;
- `expected_evidence_contract_locked`;
- `command_started`;
- `command_finished`;
- `file_mutation_observed`;
- `test_observed`;
- `redaction_applied`;
- `policy_override_requested`;
- `requirement_superseded`;
- `checkpoint_signed`;
- `verifier_result_observed`;
- `retention_lifecycle_observed`;
- `run_closed`.

Minimum first-milestone payloads:

- `recorder_attached`: recorder version, run id, run nonce, pid, cwd
  digest/label, source snapshot digest, recorder mode, clock disclosure.
- `run_started`: task ref, contract digest or default-contract marker,
  command id, wrapper name, child argv digest, child basename, source
  snapshot digest.
- `expected_evidence_contract_locked`: contract id, contract version,
  contract digest, contract source, lock timing state.
- `command_started`: command id, cwd digest/label, argv digest,
  command basename, start monotonic counter when available.
- `command_finished`: command id, exit code or signal, finish monotonic
  counter when available, stdout retention descriptor, stderr retention
  descriptor.
- `file_mutation_observed`: before source digest, after source digest,
  VCS status/diff digest when available, attribution state.
- `run_closed`: final chain head, closure state, child termination state,
  missing evidence summary digest.

Payloads outside this list are schema work for later slices and must not
block the first milestone.

Adapter lifecycle events:

- `adapter_registered`;
- `adapter_capabilities_declared`;
- `adapter_activated`;
- `adapter_error`;
- `adapter_disconnect`;
- `adapter_suppressed`.

## Adapter Contract

V0 adapter ingress is a Unix domain socket exposed as:

```text
$SDP_TRACE_SOCKET
```

Adapters send JSON events to the recorder. Adapter events can only raise
trust scope after registration and authority validation.

Socket lifecycle:

- the socket is created in a per-run private temp directory before the
  child process starts;
- filesystem socket permissions are user-only, subject to platform
  limits and current user privileges;
- `$SDP_TRACE_SOCKET` is injected into the child environment;
- the recorder drains adapter messages for a bounded grace period after
  child exit before writing `run_closed`;
- messages after `run_closed` are rejected and represented in verifier
  output as late/ignored adapter evidence;
- partial frames, parse failures, disconnects, and suppressed adapters
  emit explicit adapter lifecycle or missing evidence states;
- stale socket paths are removed before bind only inside the recorder's
  own per-run directory.

V0 uses newline-delimited canonical JSON frames over the socket. A future
profile may add length-prefixed frames if needed.

Minimum adapter registration:

- adapter id;
- harness or provider id;
- version;
- capability list;
- allowed event types;
- identity state;
- deployment source;
- registration digest.

Required OpenCode/GSD-oriented optional capabilities:

- `harness_identity_observed`;
- `tool_call_observed`;
- `model_identity_observed`.

If OpenCode or GSD cannot emit adapter events in the demo, `sdp-trace`
must still produce a useful local trace and explicit missing rows.

Adapter authority policy:

- default behavior is no trust upgrade for unlisted adapters;
- policy lists authorized adapter ids, allowed event types, signing
  requirements, and maximum trust scope;
- an adapter's self-declared `identity_state: verified` is ignored unless
  the authority policy and signature verification agree;
- unauthorized adapters may still be recorded as self-claimed local
  observations when policy permits retention, but cannot upgrade trust.

## Expected Evidence Contract

V0 uses JSON contracts only.

The contract defines what evidence is expected, not whether the run is
acceptable. Minimal fields:

```json
{
  "contract_id": "agent-pr-basic-v1",
  "version": "1.0.0",
  "contract_source": "vcs",
  "lock_required_before": "run_started",
  "required_observers": ["local_recorder", "vcs", "ci"],
  "optional_observers": ["harness_adapter", "gateway"],
  "required_events": ["recorder_attached", "task_locked", "command_started", "command_finished", "run_closed"],
  "gate_events": ["test_observed", "checkpoint_signed", "verifier_result_observed"],
  "minimum_gate_trust_scope": "ci_witnessed",
  "retention_profile": "digest_safe_default_v1",
  "redaction_profile": "prewrite_digest_default_v1"
}
```

Agent-authored contracts may support local reconstruction but cannot set
the gate bar for `ci_witnessed` trust unless a CI or human authority pins
or signs the contract digest before the run.

`wrap` accepts `--contract`. If omitted, V0 uses a built-in local default
contract that can only produce local/degraded evidence and must say so in
`dry-run`, `verify`, and `explain`. `run` requires an explicit
`--contract` unless `--use-default-contract` is provided.

Contract locking mechanics:

- load and redact-safe-parse the contract before `run_started`;
- emit `expected_evidence_contract_locked` with contract digest before
  `run_started`;
- refuse `fail_closed` runs if the contract cannot be locked;
- in `degraded` mode, continue with `contract_missing` or
  `contract_locked_late`, never with a silent pass.

## Missing Evidence

The verifier emits a MissingEvidenceTable with:

- expected role;
- expected event;
- observed state;
- sequence or window;
- reason;
- policy reference;
- replayability impact.

Observed states:

```text
present | missing | late | unsupported | unauthorized |
adapter_error | adapter_suppressed | redacted | not_assessed |
expected_run_absent
```

`expected_run_absent` is only emitted when VCS/CI/preflight evidence
indicates that a run artifact was expected for a source change and none
can be found. It is not proof of a deleted local run.

## Redaction And Retention

Default V0 profile:

```text
prewrite_digest_default_v1
```

Defaults:

- no raw prompts;
- no raw responses;
- no raw source snippets;
- no raw stdout/stderr;
- argv stored as digest plus allowlisted command basename;
- path labels or digests by default;
- raw capture only by explicit path-scoped opt-in.

Redaction must happen before persistence. A secret reaching the event log
is a product bug, not a documentation issue.

Retention descriptor:

```json
{
  "retention_mode": "digest_only",
  "replayability": "partial",
  "expiry_time": "2026-06-04T00:00:00Z",
  "archival_ref": null,
  "redaction_manifest_digest": "sha256:..."
}
```

Storage overflow emits `retention_lifecycle_observed` and prevents
`complete` completeness unless the contract explicitly accepts overflow.

Redaction failure mode:

- if a payload cannot be redacted before persistence, the sensitive event
  is not written;
- the recorder writes a sanitized `redaction_failed` event containing
  only rule id, event type, and reason digest;
- verifier downgrades completeness to `partial` or
  `missing_telemetry`;
- `fail_closed` mode aborts the child before launch when redaction policy
  cannot be loaded.

Redaction profiles include allowlisted command basenames, path-label
rules, stdout/stderr retention policy, and opt-in raw capture scopes.

Retention expiry in fixtures must be generated relative to fixture
metadata or marked non-authoritative; absolute example dates are not
policy.

## Implementation Mechanisms

### Transparent Wrapper

When stdio is interactive, `wrap` must use a pseudoterminal on Unix/macOS
so child programs see a TTY and preserve colors/interactive behavior.
When stdio is not interactive, the wrapper may use direct pipe/pass
through. The wrapper records which mode was used.

Termination handling:

- normal exit records exit code;
- signal termination records signal name/number when observable;
- SIGKILL/OOM/unobservable termination records `cannot_verify` closure
  detail;
- parent SIGINT/SIGTERM is forwarded to the child process group.

### File Mutation Observation

V0 does not do syscall interception or real-time filesystem watching.
`file_mutation_observed` means VCS/workspace delta observed by snapshots:

- record source snapshot before child start;
- record source snapshot after child exit;
- for Git workspaces, use tree/diff/status digests;
- for non-Git workspaces, emit `not_assessed` unless a later profile
  defines a snapshotter;
- Bazel sandbox internals are not treated as source mutations.

### CLI Injection Events

The CLI must support structured local observations for demo/preflight
scripts without pretending they are adapter evidence. V0 command shape:

```bash
sdp-trace observe --run <run-dir> --state not_assessed --event <event> --reason <reason>
```

`observe` can only emit `local_observed`, `not_assessed`, or
`missing_telemetry` states unless the caller is authenticated by an
authority policy.

### Run Artifact Layout

The first milestone uses this on-disk layout:

```text
<run-dir>/
  run.json
  events/
    000001-recorder_attached.json
    000002-expected_evidence_contract_locked.json
    000003-run_started.json
    ...
  artifacts/
    stdout.digest
    stderr.digest
  verifier/
    verifier-result.json
    missing-evidence-table.json
    integrity-audit.json
  export/
    audit-bundle.json
```

Rules:

- `run.json` contains run id, schema version, recorder version, created
  time, source snapshot digest, contract digest/default marker, event
  count, final chain head when closed, and closure state.
- `events/*.json` are immutable event files named by zero-padded
  sequence and event type.
- Verifier outputs live under `verifier/` and are not appended to a
  corrupted event chain.
- Export outputs live under `export/` and are reproducible from
  `run.json`, `events/`, `artifacts/`, and verifier outputs.
- Slice B and Slice E must treat this layout as their shared integration
  contract.

First-milestone correlation ids:

- `run_id`: all events.
- `command_id`: command lifecycle and related output/test observations.
- `contract_id`: contract lock and missing evidence rows.
- `source_snapshot_id`: before/after source observations.

Additional correlation domains are deferred until adapter and multi-run
query work.

## Demo Acceptance

The OpenCode + GSD + Bazel + Kotlin demo is acceptable when it shows:

1. `sdp-trace wrap` runs an existing OpenCode/GSD command without
   changing command behavior.
2. A local trace verifies as `local_observed`, not gate-grade.
3. Missing harness/model/gateway evidence is visible, not hidden.
4. Bazel test evidence is recorded when the test command runs.
5. Tampering with events produces `fail` and `integrity_audit`.
6. CI can bind source digest, contract digest, chain head, and verifier
   result into a witness record.
7. `query`, `explain`, and `export` are usable without reading raw JSON.

## Open Questions For Implementation

- Exact Go module layout is not present yet and must be introduced
  without adding another runtime stack.
- Whether GSD has a stable plugin or hook point is not yet assessed.
- Whether OpenCode can emit structured lifecycle events is not yet
  assessed.
- CI witness signing profile must be selected for the demo environment:
  Sigstore keyless, local demo key, or customer PKI equivalent.
- Demo Bazel target must be pinned to a committed Kotlin fixture or a
  customer-provided demo repository before evidence claims are made.
