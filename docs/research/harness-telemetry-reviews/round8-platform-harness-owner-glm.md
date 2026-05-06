```
VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes
- Are there any critical blockers before implementation? no
- If yes, list only true blockers that prevent starting v0.

Critical blockers:
- None.

Major gaps:
- Section 6 adapter capability example names `harness_identity_observed`,
  `tool_call_observed`, `model_identity_observed` as capabilities and
  allowed event types, but these are not defined in the section 5
  canonical event set. Adapter-emitted event schemas need a home — either
  extend section 5 with an "adapter event" subset or add a forward
  reference to an adapter event appendix. Without it, adapter
  implementers have no contract to code against.
- Transparent proxy spec (section 3) mentions TTY, colors, signals,
  exit code passthrough but does not specify what "transparent" means on
  Windows (no Unix domain socket for adapters, no SIGPIPE semantics).
  V0 can scope to Unix/macOS but should say so explicitly in the product
  boundary.
- `expected_run_absent` detection (section 8) requires "CI expects an
  sdp-trace run artifact" — who registers that expectation and where?
  The predicate is clear but the registration mechanism is unspecified.
  This is likely an `sdp-gate` concern, but v0 needs to know whether
  `sdp-trace ci-verify` emits this or an external system does.

Accepted V0 limitations:
- No attach to already-running processes. Stated and correct.
- No internal tool-call capture without adapter. Correct — the wrapper
  owns the process boundary, not the agent internals.
- Unix domain socket as sole adapter transport. Acceptable for v0 if
  platform scope is explicit.
- No raw prompt/response capture. Acceptable — not the recorder's job
  unless an adapter provides it.
- Local-only traces are not gate-grade. Correct and clearly stated.

Minimum viable correction:
- Add a "V0 adapter event types" table to section 5 or 6 with at least
  the three event types named in the capability example, even if
  payloads are stub `{ "adapter_event_ref": "...", "observed_data": "..." }`.
  This closes the adapter contract gap without blocking implementation.
- Add one sentence to section 3 or 4: "V0 targets Unix and macOS.
  Windows adapter transport is a future profile."

Questions before implementation:
- Is `sdp-trace ci-verify` the component that emits `expected_run_absent`,
  or does an external governance system own that predicate? Affects where
  the CI artifact store lookup logic lives.
- Should the v0 recorder validate adapter event types against the
  adapter's declared `allowed_event_types` at ingestion time (reject
  undeclared types), or just record them with a trust downgrade?
  The former is safer but couples recorder to adapter schema.

Demo changes required:
- Demo 0 should show the wrapper composing around a real harness CLI
  (even a stub), not just `sdp-trace run`. This is the primary adoption
  path and the demo that proves the capture boundary works.
- Demo 3 should show an adapter registering via Unix socket, emitting
  one event, then disconnecting — this proves the adapter contract is
  implementable, not just theoretical.
```
