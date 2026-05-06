# Agentic SDLC Evidence Substrate Brief v4

Status: discussion draft; not committed
Date: 2026-05-05
Inputs:
- `docs/research/agentic-sdlc-evidence-substrate-v3-brief.md`
- `docs/research/harness-telemetry-reviews/round7-socratic-consolidation.md`

This brief is a product framing artifact. It is not source-bound proof,
not product closure evidence, and not a trusted release claim.

## 1. Product Boundary

`sdp-trace` is a portable evidence substrate for AI-assisted delivery.
It records observed events, observer provenance, evidence retention,
explicit gaps, trace integrity, witness scope, and verifier output.

It does not decide business pass/fail, readiness, degradation, waiver
acceptance, or opaque health scores. Those belong to `sdp-gate`,
`sdp-report`, or another policy consumer.

Compared with existing artifacts:

- CI logs show what CI ran.
- Git diff shows what changed.
- Review comments show what humans discussed.
- `sdp-trace` adds signed provenance, explicit missing telemetry,
  trust-scope classification, tamper evidence, retention/replayability
  state, and a query surface that connects the pieces.

## 2. V0 Product Promise

V0 is not a full agent governance platform. It is a buildable recorder
and verifier contract.

V0 supports:

- local command-wrapper observation;
- wrapper composition for existing harness commands;
- CI-attached observation;
- optional adapter events;
- expected evidence contracts;
- pre-write redaction;
- four-axis verifier output;
- MissingEvidenceTable;
- developer `dry-run`, `verify`, `query`, and `explain`;
- signed verifier result in CI.

V0 does not support:

- reliable attach to already-running processes;
- full multi-harness SDK;
- raw prompt/response capture;
- org-wide degradation dashboard;
- local-only audit-grade claims.

## 3. Invocation And Capture Boundary

Primary local command:

```text
sdp-trace run --task <task-ref> --contract <contract-file> -- <command...>
```

Wrapper composition:

```text
sdp-trace wrap --name <wrapper-name> -- <existing-wrapper> <args...>
```

`wrap` still launches the process from the beginning. It is not a
retroactive `attach <pid>`. If work started before the wrapper, the run
must emit `late_attach_boundary` or `expected_run_absent` where
detectable.

If a team already has a harness CLI or launch wrapper, prefer
`sdp-trace wrap -- <existing-harness-cli> <args...>`. Use
`sdp-trace run` directly only when no intermediate wrapper exists.

`wrap` must be a transparent process proxy: stdin, stdout, stderr, TTY
mode, colors, signals, and exit code pass through as if the command ran
without `sdp-trace`. On SIGINT/SIGTERM, the recorder forwards the signal
to the child and emits `run_closed` with the observed signal or
`cannot_verify` if closure cannot be written.

CI-attached command:

```text
sdp-trace ci-verify --run <run-artifact> --contract <contract-file>
```

V0 local recorder captures:

- recorder attachment;
- task and contract digests;
- source snapshot before and after;
- command start/finish;
- argv digest and allowlisted command basename;
- exit state;
- stdout/stderr retention descriptors;
- VCS diff/tree digest after command;
- explicit missing harness/model/gateway/tool-call rows.

V0 runs fully offline. Local capture, local verification, local query,
and local explain do not require network access. CI signing and external
witnessing are additive.

Local recorder performance and output:

- local writes are append-only and async;
- no network call in the local inner loop;
- p99 wrapper overhead target: <= 5 ms per event excluding child command
  execution time;
- normal recorder stderr output is at most one start line and one summary
  line unless `--verbose` is enabled.

V0 local recorder does not capture:

- internal tool calls that never cross the wrapped process boundary;
- prompts or responses;
- model identity unless adapter/gateway evidence exists;
- file writes after run close.

## 4. Operational Modes

V0 has two modes:

```text
mode: degraded | fail_closed
```

`degraded`:

- continues execution when recorder/adapter evidence is missing;
- emits `partial`, `missing_telemetry`, or `cannot_verify`;
- never silently upgrades trust.

`fail_closed`:

- aborts before launching the child process if required local recorder
  setup fails;
- exits non-zero if the expected contract cannot be locked before run;
- is intended for managed CI or managed harness environments.

Mode is recorded in `recorder_attached` and `run_started`.

## 5. Canonical Event Set

Every event includes:

- schema version;
- run id;
- event id;
- sequence;
- event time;
- observer id;
- observer role;
- claimed role;
- verified role state: `verified | self_claimed | unauthorized | not_assessed`;
- trust scope;
- payload digest;
- previous event hash;
- event hash;
- correlation ids;
- retention descriptor;
- redaction state.

Required V0 events:

| Event | Purpose |
| --- | --- |
| `recorder_attached` | First event; binds run id, nonce, recorder version, process id, workspace/source digest, mode, local clock state. |
| `run_started` | Binds task ref, contract ref, source snapshot, actor, wrapper command, parent run refs. |
| `task_locked` | Records task/spec/plan digest and lock provenance. |
| `expected_evidence_contract_locked` | Records contract digest, version, source, policy owner, and lock timing state. |
| `command_started` | Records command id, argv digest, cwd label/digest, start time. |
| `command_finished` | Records command id, exit state, finish time, stdout/stderr retention. |
| `file_mutation_observed` | Records VCS/tree/diff evidence and attribution state. |
| `test_observed` | Records test command, result, artifact digest, and execution environment digest. |
| `redaction_applied` | Records pre-write redaction rule and manifest digest. |
| `policy_override_requested` | Records signed human override request and linked missing evidence. |
| `requirement_superseded` | Records task/contract/spec change with old/new digests. |
| `checkpoint_signed` | Records sequence range, chain head, signer identity, signing profile. |
| `verifier_result_observed` | Records verifier version, input digests, four-axis output, MissingEvidenceTable digest. |
| `retention_lifecycle_observed` | Records archive, expiry, deletion, or key-retention lifecycle events. |
| `run_closed` | Terminal local event; records chain head, missing evidence summary, closure state. |

Optional adapter lifecycle events:

- `adapter_registered`;
- `adapter_capabilities_declared`;
- `adapter_activated`;
- `adapter_error`;
- `adapter_disconnect`;
- `adapter_suppressed`.

These events let the verifier distinguish:

- unsupported telemetry;
- adapter absent;
- adapter crashed;
- adapter suppressed;
- observer unauthorized.

V0 adapter transport:

- adapters emit JSON events to a Unix domain socket named by
  `$SDP_TRACE_SOCKET`;
- the recorder creates the socket with user-only permissions;
- unsupported platforms may use loopback HTTP in a later profile, but
  Unix socket is the V0 transport.

## 6. Adapter Identity And Capability Contract

Adapters must register before their events can upgrade trust scope.

Adapter registration includes:

- adapter id;
- provider/harness id;
- version;
- capability list;
- signing identity or unsigned/self-claimed state;
- allowed event types;
- deployment source;
- registration digest.

Capability example:

```json
{
  "adapter_id": "opencode-local-v1",
  "harness": "opencode",
  "capabilities": [
    "harness_identity_observed",
    "tool_call_observed",
    "model_identity_observed"
  ],
  "identity_state": "self_claimed | verified",
  "allowed_event_types": ["tool_call_observed", "model_identity_observed"]
}
```

If identity is not verified, events remain `self_claimed` and do not
upgrade beyond the verifier's honest floor.

## 7. Expected Evidence Contract

Contract is locked before `run_started` for normal runs.

If contract is absent or locked late:

- verifier emits `contract_missing` or `contract_locked_late`;
- missing evidence rows are still generated from any later contract;
- gate-grade trust is unavailable unless CI policy explicitly accepts
  the late-lock state.

Gate-grade contract provenance requires one of:

- contract committed in VCS before the source commit being verified;
- contract digest pinned in CI config;
- human-signed contract approval;
- accepted customer policy equivalent.

Agent-authored contracts can support local reconstruction but cannot set
the evidence bar for `ci_witnessed` gate-grade trust.

Minimal contract fields:

```json
{
  "contract_id": "agent-pr-basic-v1",
  "version": "1.0.0",
  "contract_source": "vcs | ci_config | human_signed | agent_reported",
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

`sdp-trace` records satisfaction facts. `sdp-gate` decides whether a gap
blocks.

V0 contract format is JSON. YAML or other formats are future profiles.

For `ci_witnessed`, a VCS-sourced contract must be pinned by digest in CI
configuration or human-signed before the run. If the contract digest is
not pinned or signed, the maximum trust scope is downgraded even if all
events are present.

## 8. Missing Evidence And No-Run Detection

Verifier emits a MissingEvidenceTable with:

- expected role;
- expected event;
- observed state;
- sequence/window;
- reason;
- policy reference;
- replayability impact.

Observed states:

```text
present | missing | late | unsupported | unauthorized |
adapter_error | adapter_suppressed | redacted | not_assessed |
expected_run_absent
```

`expected_run_absent` is emitted when source/VCS/CI evidence indicates a
change that should have had a trace, but no matching run artifact exists.
This is a governance gap, not proof of malicious deletion.

Concrete V0 predicate:

- a VCS commit or PR touches files in the contract/task scope;
- CI expects an `sdp-trace` run artifact for that commit or PR;
- no run artifact carries a matching task ref and source commit/tree
  digest in the CI artifact store.

If a developer deletes a local run before any VCS/CI/preflight observer
expects it, `sdp-trace` cannot prove the deleted run existed. That is an
accepted V0 limitation unless a pre-commit hook, CI preflight, or
mandatory upload policy is deployed.

## 9. Privacy, Redaction, Retention

Default profile: `prewrite_digest_default_v1`.

Defaults:

- no raw prompts;
- no raw responses;
- no raw source snippets;
- no raw stdout/stderr;
- argv stored as digest plus allowlisted command basename;
- file paths stored as safe labels or digests unless explicitly enabled;
- raw capture is opt-in, path-scoped, and retention-scoped.

Redaction happens before persistence. A secret that reaches the event log
is a bug.

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

Storage behavior:

- local default floor: at least 100 MB or 10,000 events per run,
  whichever is hit first; implementations may raise this;
- overflow must emit `retention_lifecycle_observed` with
  `storage_limit_reached`;
- overflow cannot silently drop events and still report `complete`.
- V0 retention descriptors must include non-null `expiry_time` derived
  from the selected retention profile.

## 10. Signing, Keys, Witness

Local V0 key profile:

- per-run ephemeral key generated by the recorder process;
- stored only in recorder memory while the run is active;
- public key written into the run metadata;
- supports local structural continuity only;
- host compromise or agent access to process memory defeats it;
- never upgrades to gate-grade trust.

Role authentication:

- observer role strings are claims until signer authority policy
  validates them;
- unauthorized or unknown signers produce `self_claimed` or
  `unauthorized`, never trust upgrade;
- verified adapter/CI/human roles require an authority policy entry.

Human override profile:

- V0 accepts `human_signed` only through a configured signing profile:
  local named key, SSO/OIDC assertion, WebAuthn, or customer equivalent;
- unsigned human names are recorded as `human_declared`, not
  `human_signed`.

CI witness signs:

- verifier result;
- source commit/tree digest;
- chain head;
- expected contract digest;
- verifier version;
- policy profile;
- CI job identity;
- timestamp;
- input artifact digests.

If the agent and verifier run in the same CI job, topology is recorded.
Events whose signer is not independent of the actor being assessed are
downgraded.

CI witness key boundary:

- the signing key or OIDC signing operation must occur in a process,
  service, or CI step that is not a fork/spawn child of the agent process;
- verifier records `witness_independence: independent | same_job |
  same_process | not_assessed`;
- only `independent` can support `ci_witnessed` gate-grade trust.

CI witness verification must bind run id, source commit/tree digest,
workspace digest, contract digest, chain head, and verifier result. A
local trace replayed against another source commit is downgraded unless
that exact tuple was CI- or externally witnessed.

## 11. Verifier Output

Machine output is four-axis:

```text
verdict: observed | fail | cannot_verify | not_assessed
trust_scope: agent_reported | local_observed | harness_observed |
             gateway_observed | vcs_observed | ci_witnessed |
             external_witnessed | human_signed
completeness: complete | partial | missing_telemetry | unknown
replayability: full | partial | none | not_assessed
```

V0 deliberately avoids a generic local `pass` label. Human-readable
output must lead with trust scope:

```text
Trust scope: local_observed
Verdict: observed
Completeness: partial
Replayability: partial
Gate usable: false
Reason: no CI or external witness before gate
```

`verifier_result_observed` is added to the chain and signed by the
verifier identity. CI-witnessed results are signed by CI.

## 12. Query And Explain

Required queries:

- `run-summary`;
- `timeline`;
- `missing-evidence`;
- `commands`;
- `files`;
- `tests`;
- `redactions`;
- `witness`;
- `overrides`;
- `retention`.

Audit export:

```text
sdp-trace export <run-dir> --format audit-bundle
```

The export contains event chain, chain head, verifier result,
MissingEvidenceTable, witness assertions, retention manifest, and
integrity audit record when present.

If the chain is corrupted, the verifier writes a separate
`integrity_audit` record rather than appending to the corrupted chain.
The record includes failure type, failing event ref when known, verifier
version, and input digest.

`sdp-trace explain <run-dir>` handles:

- no run directory;
- empty run directory;
- no events;
- corrupt event chain;
- local-only trace with no witness;
- missing contract;
- contract locked late;
- storage overflow;
- redaction downgrade;
- unauthorized signer;
- expected run absent.

Empty/no-run output is explicit:

```text
No sdp-trace run found.
Trust scope: not_assessed
Completeness: missing_telemetry
Suggested next step: run through sdp-trace run or configure CI preflight.
```

## 13. Developer Workflow

Commands:

```text
sdp-trace dry-run --contract <contract-file> -- <command...>
sdp-trace run --task <task-ref> --contract <contract-file> -- <command...>
sdp-trace wrap --name <wrapper-name> -- <existing-wrapper> <args...>
sdp-trace verify <run-dir>
sdp-trace query <run-dir> --query missing-evidence
sdp-trace explain <run-dir>
sdp-trace export <run-dir> --format audit-bundle
sdp-trace ci-verify --run <run-artifact> --contract <contract-file>
```

`dry-run` is a redaction and contract simulation, not proof that the
future run will emit the same events. It must say this.

Dry-run example:

```text
Simulation only: no event chain was written.
Command basename: deploy
Argv retention: digest_only
Raw argv: not stored
Redactions:
- --token=<redacted> by rule secret.flag.token
Contract: agent-pr-basic-v1 sha256:...
Expected trust without CI: local_observed / partial / partial
```

Telemetry opt-out:

- explicit opt-out emits no local run;
- CI/VCS may later emit `expected_run_absent`;
- emergency work should use `policy_override_requested`, not silent
  bypass.

## 14. Demo Set

Demo 0: wrapper-composed local observation.

- Run existing harness command through `sdp-trace wrap`.
- Output: `observed + local_observed + partial + partial`.
- Show missing harness/model/gateway rows.
- Show transparent TTY/stdout/stderr/signal/exit-code passthrough.

Demo 0.25: fully offline local run.

- Run with local minimal contract and no CI/network.
- Show local verification and explain output.
- Output remains not gate-grade.

Demo 0.5: CI-attached run.

- Run `sdp-trace run` inside CI.
- CI signs verifier result.
- Show fastest path to gate-grade evidence.

Demo 1: dry run and redaction.

- Include fake secret in argv.
- Show dry-run warning and pre-write redaction behavior.

Demo 2: CI-witnessed gate.

- CI verifies source, chain, contract, test artifact, verifier version.
- Output can become `ci_witnessed + complete` only if contract is met.
- Include a relaxed agent-authored contract attempt; verifier rejects
  trust upgrade because contract digest is not pinned/signed.

Demo 3: missing and suppressed observers.

- Omit adapter once.
- Suppress adapter once.
- Show distinct MissingEvidenceTable states.

Demo 4: tamper and chain break.

- Mutate, delete, reorder events.
- Show `explain`.
- Show `sdp-trace export` containing an `integrity_audit` record.
- Include accidental post-run file move as a non-malicious mismatch.

Demo 5: post-hoc fabrication/replay.

- Replay old trace or regenerate local trace after work.
- Show failure/downgrade through missing witness, VCS mismatch,
  timestamp mismatch, process nonce mismatch, or expected-run gap.
- Include a trace whose workspace digest matches the current commit but
  lacks a witnessed tuple binding run id, commit digest, contract digest,
  and chain head.

Demo 6: human override.

- Unsigned human declaration stays `human_declared`.
- Signed override becomes `human_signed + partial`, not complete.

Demo 7: forensic cold query.

- Query after retention lifecycle/key metadata event.
- Show what remains replayable and what is not assessed.

Demo 8: latency and overflow.

- Run high-frequency local events.
- Show p99 overhead, storage cap, and overflow behavior.

## 15. Round 7 Review Task

Review this v3 brief through your assigned persona.

Return:

```text
VERDICT: ACCEPTABLE_WITH_GAPS | CHANGES_REQUIRED | REJECTED

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes/no
- Are there any critical blockers before implementation? yes/no
- If yes, list only true blockers that prevent starting v0.

Critical blockers:
- ...

Major gaps:
- ...

Accepted V0 limitations:
- ...

Minimum viable correction:
- ...

Questions before implementation:
- ...

Demo changes required:
- ...
```

Tie findings to:

- CTO usefulness;
- v0 capture boundary;
- adapter interface;
- expected evidence contracts;
- observer authority;
- signing and verification;
- privacy and retention;
- adoption and DX;
- forensic query/replay;
- demo credibility.

Convergence target: `ACCEPTABLE_WITH_GAPS`, implementation start `yes`,
and no critical blockers. Remaining major gaps may become implementation
tasks or explicit V0 limitations.
