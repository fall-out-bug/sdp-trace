```text
VERDICT: CHANGES_REQUIRED

Critical blockers:

- No adapter interface contract is defined anywhere. The brief lists
  what a harness adapter *owns* (Section 2) but not what it must emit,
  how it plugs in, or what the minimum integration surface is. The
  persona requires `run_started`, `task_locked`, `tool_call`,
  `command_started`, `file_mutation`, `test_observed`, `run_closed` as
  an adapter contract. The brief has no equivalent. Without this, the
  first build slice cannot be built and the demo cannot demonstrate a
  real capture boundary. (observer authority, evidence model, adoption
  and DX)

- The "local recorder" in Demo 1 and the first build slice is undefined.
  Is it a shell wrapper, a process tracer, a file watcher, a git hook,
  a PTY interceptor, an eBPF probe? The brief says "local recorder
  fixture path" which reads like test data, not a running interceptor.
  The single most important integration point for v0 is unspecified.
  (evidence model, demo credibility, CI/gate anchoring)

- No operational model for "developer runs agent outside the wrapper."
  Section 11 lists this as a risk but the brief provides no mechanism to
  detect or represent it. The correct answer is likely: the run emits
  `harness_observed: missing` and `local_observed` only. But the brief
  does not say this, and no event in the trace model represents "no
  harness was present." The `missing_telemetry` concept applies to
  individual observers within a run, not to the absence of the capture
  layer itself. (observer authority, trace model, provenance model)

Major gaps:

- Event type taxonomy is absent. Section 4 lists evidence *examples* but
  the trace model does not enumerate event types or bind each type to
  its valid observer roles. Without this mapping, the authority rule in
  Section 5 ("a signer can only support trust for event types it is
  authorized to observe") cannot be enforced by the verifier.
  (observer authority, signing and verification)

- No fail-closed or degraded mode specification. Managed harnesses
  should fail closed (block the run if the recorder cannot attach).
  Unmanaged harnesses should degrade explicitly. The brief is silent on
  operational modes. If the recorder cannot intercept tool calls, what
  does the adapter emit? What does the trace look like? (adoption and DX,
  evidence model)

- Demo 5 (post-hoc fabrication) shows the correct verdict but does not
  explain the detection mechanism. "No checkpoint witnessed before gate
  decision" requires a witness that existed at the time. In the local
  case there is no such witness. The demo should show what the verifier
  actually checks to reach this verdict: absence of harness_observed,
  absence of gateway_observed, timestamp analysis, or explicit
  "recorder was not attached at run start" fact. (demo credibility,
  signing and verification)

- The brief says "Events do not all need individual signatures in v0.
  They must be content-addressed and chained." This is acceptable for
  ci_witnessed checkpoints but under-specifies the tamper detection
  story for local_observed runs where no checkpoint signer exists until
  CI. An attacker who controls the local host can rewrite the entire
  chain before CI sees it. The brief should state explicitly: local
  chains are tamper-detectable only after first external witness, not
  before. (signing and verification, provenance model)

- Expected evidence contract schema is listed in the first build slice
  but the contract contents (Section 7) are described in prose with no
  schema shape. Who defines the contract? Where does it live? Is it a
  file in the repo, a CI config, a policy document? (adoption and DX,
  expected evidence contracts)

False assumptions:

- The brief assumes events will be emitted by named observers without
  specifying the interception mechanism. Observation requires a capture
  boundary. "harness_observed" assumes a harness adapter is attached and
  running before agent work begins. The brief does not specify the
  lifecycle: when the adapter starts, how it attaches, how it survives
  agent crashes, or how it signals "I was here from the beginning."
  (evidence model, provenance model)

- The brief assumes "no checkpoint witnessed before gate decision" is
  detectable in Demo 5 but does not explain what timestamps or anchors
  the verifier uses. If all timestamps are local and self-reported,
  the verifier cannot distinguish real-time from post-hoc without an
  external clock source. (signing and verification, CI/gate anchoring)

Standards we should reuse instead of building:

- The brief correctly adopts OpenTelemetry, in-toto, DSSE, Sigstore,
  SLSA. One gap: the adapter contract should map directly to
  OpenTelemetry spans with GenAI semantic conventions rather than
  defining a bespoke event vocabulary. The brief does not say whether
  events are OTel spans or a custom shape. If custom, justify why OTel
  spans are insufficient. (adoption and DX, evidence model)

- Consider reusing SLSA's `byDefinition` vs `byDigest` pattern for
  evidence binding instead of inventing a new binding model.
  (evidence model)

Minimum viable product correction:

- Add one section: "Adapter Interface" with the minimum event set,
  required fields per event type, valid observer roles per event type,
  and the lifecycle (start, attach, heartbeat or liveness, detach,
  close). This is not deferred work. It is the build slice.
  (adoption and DX, observer authority)

- Replace "local recorder fixture path" with a concrete recorder design:
  what it wraps, how it is invoked, what it emits, how it signals
  "recorder attached at run start." Even if the v0 implementation is
  a shell wrapper stub, the design must be stated. (demo credibility,
  evidence model)

- Add one event type: `recorder_attached` with observer role
  `local_observed`, emitted before any other event, to anchor the
  recorder's presence and start time. This makes Demo 5's detection
  mechanism concrete. (trace model, signing and verification)

Questions before implementation:

- Is the local recorder a shell wrapper, a process monitor, a git hook,
  or something else? This determines the entire capture boundary and
  what "local_observed" actually means in practice. (evidence model)

- Who creates and stores the expected evidence contract? Is it committed
  to the repo, configured in CI, or generated by the harness?
  (expected evidence contracts, adoption and DX)

- What happens when two observers disagree? If the harness says
  `tool_call` happened but the local recorder shows no corresponding
  `command_started`, what does the verifier emit?
  (signing and verification, observer authority)

Attack scenarios not covered:

- Agent spawns a subprocess that bypasses the recorder. The brief does
  not address whether the recorder intercepts child processes or only
  the top-level shell. (evidence model, provenance model)

- Developer modifies the local recorder binary or config before the run.
  The brief assumes the recorder is trusted software but does not
  address recorder integrity. (signing and verification, provenance
  model)

- Agent writes files directly via API or library call rather than
  through a wrapped shell command. File mutation events require either
  a filesystem watcher or an agent SDK hook. The brief defers "file
  watcher attribution" but this is core to local_observed credibility.
  (evidence model, demo credibility)

Demo changes required:

- Demo 1 must specify the recorder mechanism and show the
  `recorder_attached` event as the first event in the chain.
- Demo 3 must show the trace emitted when the recorder is not attached
  at all, not just when model/harness telemetry is missing.
- Demo 5 must state the specific verifier check that detects post-hoc
  fabrication, not just the verdict.
- All demos should include the concrete event types being emitted, not
  just the verdict summary.
```
