**Top 10 Adoption Blockers**

1. **No local recorder implementation exists.** Section 9 admits tool/shell/file capture is unbuilt. I cannot adopt a flight recorder that has no binary, wrapper, or filesystem watcher. *(telemetry capture, adoption)*

2. **Zero harness adapter contracts.** There is no spec or reference implementation for pi, OpenCode, or Kilo adapters. The integration placement table is a wishlist, not an interface I can ask a team to implement. *(integration placement, adoption)*

3. **No CI signing recipe for any platform.** Section 5 claims CI-signed chain heads are stronger, but Section 9 admits no recipe exists. I have no GitHub Actions, GitLab CI, or Jenkins example to give my platform team. *(signing, adoption)*

4. **Managed/unmanaged mode lacks enforcement mechanics.** Section 7 describes gradual adoption without a feature flag, config schema, or lockdown mechanism. An agent can simply omit the wrapper and the system falls back to unmanaged. *(agent cooperation, adoption)*

5. **LLM gateway integration has no privacy or retention model.** Section 3A notes prompt/privacy risk but offers no PII scrubbing spec, retention policy, or provider contract. Legal will block this before engineering sees it. *(telemetry capture, adoption)*

6. **Remote witness protocol is undefined.** Section 6 lists `externally_witnessed` as a verifier state, but Section 9 admits no remote witness protocol exists. No endpoint, no REST schema, no append-only log format. *(anti-forgery, adoption)*

7. **No Bazel or Kotlin evidence schema.** Section 9 explicitly lists "Kotlin+Bazel real demo evidence" as missing. There is no mapping from Bazel Build Event Protocol or Kotlin test outputs to verifier evidence JSON. *(CTO usefulness, adoption)*

8. **No harness registry for drift detection.** Section 8 asks whether teams are drifting into unapproved harnesses, but Section 9 admits the harness registry is missing. I cannot query what I cannot define. *(anti-forgery, adoption)*

9. **Verifier states are prose, not schema.** Section 6 lists eight states with no JSON schema, state machine, or transition rules. Two different teams will ship incompatible verifiers. *(anti-forgery, adoption)*

10. **No scope enforcement policy engine.** Section 8 asks if file changes were inside allowed scope, but there is no path-allowlist schema, no scope violation detection logic, and no blocking mechanism. *(anti-forgery, CTO usefulness)*

---

**Top 5 Changes Needed Before Kotlin+Bazel Demo**

1. **Ship a Bazel test/build event parser.** Convert `bazel test` XML or Build Event Protocol output into `sdp-trace` evidence JSON. The demo must show real build evidence, not hand-written fixture files. *(telemetry capture, CTO usefulness)*

2. **Commit a harness adapter for one target harness.** Pick pi or OpenCode. Define the hook contract (task start, tool call, subagent boundary) and implement it so the demo captures live agent telemetry rather than agent self-reports. *(integration placement, adoption)*

3. **Provide a working local recorder for shell and file scope.** A shell wrapper or filesystem watcher must emit schema-compliant events during the Kotlin build. The demo needs observed command executions and file mutations. *(telemetry capture, anti-forgery)*

4. **Freeze a verifier JSON schema and demo a verdict run.** The demo must end with the verifier consuming the Kotlin+Bazel trace and emitting a concrete state (`pass`, `missing_telemetry`, etc.) against a committed schema, not a script that prints green text. *(anti-forgery, adoption)*

5. **Cut a v0.1 committed spec.** Change the header from "discussion draft; not committed" to a versioned spec with a schema URI and a date. I will not approve headcount or demo time for a document that can change underneath us the next day. *(adoption)*
