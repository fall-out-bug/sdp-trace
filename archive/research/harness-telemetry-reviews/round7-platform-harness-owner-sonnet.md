```text
VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes
- Are there any critical blockers before implementation? no

Critical blockers:
- None. local recorder + run/wrap commands + fail-closed/degraded modes +
  adapter lifecycle events + MissingEvidenceTable are sufficient to start.

Major gaps:

1. Adapter transport protocol unspecified.
   Section 6 defines the adapter capability contract (what adapters
   declare) but not the IPC mechanism (how adapters push events into the
   recorder: Unix socket, HTTP, stdin pipe, file tail, shared memory?).
   Harness owners cannot implement an adapter without this. Needs even a
   single concrete V0 transport, marked extensible.

2. expected_run_absent detection is hand-wavy.
   Section 8 says "source/VCS/CI evidence indicates a change that should
   have had a trace." The predicate is undefined. What VCS evidence? Commit
   exists with no matching run-id tag? CI job ran but no artifact uploaded?
   Without a concrete detection rule, CI pipelines cannot implement the
   check, and Demo 5 has no enforceable baseline.

3. wrap vs run distinction is underdeveloped for integrators.
   Section 3 distinguishes them but gives no guidance on which to use when
   the harness already has its own wrapper (e.g., opencode's CLI). A
   harness owner deploying this needs a decision tree: wrap the harness CLI,
   replace the launch script, or inject via shell alias? One concrete
   integration example is missing.

4. Storage cap is "implementation-defined."
   Section 9 says "local default cap: implementation-defined, but must be
   explicit in recorder_attached." That defers a real operational question.
   Platform owners need a floor (e.g., 100 MB or 10k events) to size
   sandboxes and set overflow policy before deployment.

Accepted V0 limitations:
- No reliable attach to already-running processes. Honest and correct.
- No full multi-harness SDK. Wrapper-first is the right V0 scope.
- No raw prompt/response. Defensible default; opt-in path exists.
- Local-only trust never becomes gate-grade. Correctly stated.
- Agent-authored contracts cannot set the gate-grade evidence bar. Correct.

Minimum viable correction:
- Add one sentence to Section 6: nominate a single V0 adapter transport
  (e.g., "V0 adapters emit events over a Unix domain socket at
  $SDP_TRACE_SOCKET; full transport spec is a separate document").
- Add one concrete predicate to Section 8 for expected_run_absent (e.g.,
  "a VCS commit touches files in the task scope but no run artifact carries
  a matching task-ref and commit digest in the CI artifact store").
- Add a single integration decision note to Section 3: "If the harness
  has its own launch wrapper, prefer sdp-trace wrap -- <harness-cli> <args>;
  use sdp-trace run directly only when no intermediate wrapper exists."
- Pick a concrete storage floor in Section 9 or mark it a V0 impl task
  with a filed issue number.

Questions before implementation:
- Is the V0 adapter transport Unix socket, HTTP on loopback, or something
  else? This gates any harness integration work.
- Does expected_run_absent require CI artifact store access, or can it be
  approximated from a local git log scan? The answer changes the verifier
  deployment footprint.

Demo changes required:
- Demo 0 (wrapper-composed local observation): add the concrete wrap
  invocation for at least one real harness CLI (e.g., opencode or cursor
  CLI) so it is not purely hypothetical.
- Demo 5 (post-hoc fabrication): must show the concrete expected_run_absent
  trigger, not just "VCS mismatch." Otherwise the demo cannot be
  reproduced by a platform team verifying the claim.
- Demo 3 (missing vs suppressed): add adapter_absent as a third distinct
  row in the MissingEvidenceTable output, separate from adapter_error and
  adapter_suppressed. The brief defines all three states; the demo should
  show all three.
```
