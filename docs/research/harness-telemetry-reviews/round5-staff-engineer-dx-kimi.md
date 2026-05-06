VERDICT: CHANGES_REQUIRED

Critical blockers:
- Default capture policy is undefined. Acknowledging that "raw prompts, responses, argv, stdout, stderr, and file paths may leak confidential data" is not a policy. Without a default-redact or digest-only rule, developers will treat this as a secret exfiltration vector and refuse adoption. (privacy and retention, adoption and DX)
- No first-class emergency path in the first slice. Emergency fixes are real; if the trace substrate blocks or gates reject without a documented `policy_override_requested` event type that is explicit and auditable, developers will bypass or disable the tool. (adoption and DX, expected evidence contracts)
- Latency and offline workflow are unspecified. The brief assumes CI/OIDC/Sigstore for gate-grade trust, but local agent loops cannot tolerate network round-trips per event. If local chaining adds friction to every file save or tool call, it will be routed around. (adoption and DX, trace model)

Major gaps:
- No local preview or dry-run mode. Developers need to see exactly what evidence types, digests, and correlations will be captured before any data is hashed or written. (adoption and DX, evidence model)
- Verifier explainability is gate/CTO-centric, not developer-centric. The three-axis output is useful for policy consumers, but there is no human-readable `why did my trace fail?` report with event names and specific mismatch details. A developer debugging a false positive needs prose, not just a `fail` verdict. (signing and verification, adoption and DX)
- Redaction is modeled as retention-state metadata rather than a capture boundary filter. Retention state "is part of the evidence, not an afterthought" implies the secret is already in the payload. Secrets must be filtered before hashing and writing. (privacy and retention, evidence model)
- No friction budget or performance requirements for the hash chain. High-frequency agent loops (file watchers, rapid tool calls) will generate thousands of events; unbounded SHA-256 chaining and canonical JSON serialization is a real tax on every operation. (trace model, adoption and DX)
- argv, stdout, stderr, and file paths are listed as evidence types without default exclusion or sanitization rules. These are secret-bearing by default and should not be raw-captured unless explicitly opted into a scoped profile. (evidence model, privacy and retention)

False assumptions:
- That developers will tolerate local hash-chain overhead even though the brief correctly labels local traces as non-gate-grade. If local is weak and adds latency, it is dead weight unless it is zero-cost and non-blocking. (trace model, adoption and DX)
- That developers understand DSSE/in-toto/SLSA vocabulary well enough to debug local trace failures. These standards are correct for CI boundaries but are opaque for a developer running a local recorder. (adoption and DX, signing and verification)
- That "canonical JSON" and per-event SHA-256 hashing is cheap in high-frequency agent workflows. In practice, this adds serialization and I/O cost on every event. (trace model, adoption and DX)

Standards we should reuse instead of building:
- Reuse OpenTelemetry context propagation and baggage for correlation IDs across harness/tool/shell/CI boundaries rather than a custom hash-chain at the agent-loop level. (trace model, product layering)
- Reuse existing OpenTelemetry/LangSmith/Langfuse redaction and sanitization patterns for secrets in spans rather than inventing a new retention-state schema that still captures the secret. (privacy and retention, evidence model)

Minimum viable product correction:
- Default evidence profile must be `digest-only` for argv/stdout/stderr/file paths; raw capture must be opt-in and path-scoped. Redaction must happen before write, not as post-hoc metadata. (privacy and retention, evidence model)
- Local recorder must be append-only, async, and non-blocking, with offline-capable ephemeral signing that allows CI counter-signing later. It must not require network access or block the agent loop. (trace model, adoption and DX)
- Add a developer-facing `explain` command that translates verifier failures into human-readable event names, authority mismatches, and chain breaks without requiring DSSE expertise. (signing and verification, adoption and DX)
- Make `policy_override_requested` a first-class event type in the expected evidence contract, requiring `human_signed` justification, rather than burying it as a manual exception. (expected evidence contracts, adoption and DX)
- Ship `sdp-trace --dry-run` in the first slice so developers can preview capture scope before running a task. (adoption and DX, evidence model)

Questions before implementation:
- What is the p99 latency budget per tool call or file mutation? If it is >5 ms, this will not survive an agentic inner loop. (adoption and DX, trace model)
- Does the local recorder run as a daemon, library, or wrapper? Who owns the process lifecycle and crash recovery? (adoption and DX, product layering)
- What is the offline fallback when Sigstore/OIDC is unreachable? Is the trace unobserved or locally signed with delayed witness? (signing and verification, CI/gate anchoring)
- How does a developer redact a specific argv flag or env var without disabling the recorder entirely? (privacy and retention, adoption and DX)
- Can a developer pause local tracing for a sensitive task without breaking the chain or creating a suspicious gap? (trace model, adoption and DX)

Attack scenarios not covered:
- Secret exfiltration via the evidence substrate: a compromised recorder or misconfigured adapter streams raw stdout/stderr to a shared observability backend, using the trace as the exfil channel. (privacy and retention, evidence model)
- Denial-of-service via trace bloat: a runaway agent or malicious tool generates millions of events, exhausting local storage or CI verifier compute budgets. (trace model, adoption and DX)
- Emergency hotfix abandonment: a developer disables the recorder because the override workflow is too slow or unclear, resulting in a total audit gap rather than a signed override. (adoption and DX, expected evidence contracts)

Demo changes required:
- Demo 1 must include a `--dry-run` preview showing exactly what evidence is captured and what is redacted. (adoption and DX, demo credibility)
- Demo 3 (Missing Observer) must demonstrate the emergency override path, emitting `policy_override_requested` and resulting in `human_signed + partial`, proving the gap is explicit, not hidden. (expected evidence contracts, demo credibility)
- Add a latency/friction demo: a high-frequency local agent loop with the recorder enabled, proving negligible overhead. (adoption and DX, demo credibility)
- Add a redaction demo: an argv containing a secret is captured as digest-only because the recorder rejects raw evidence, falling back safely. (privacy and retention, demo credibility)
