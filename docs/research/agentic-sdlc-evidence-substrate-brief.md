# Agentic SDLC Evidence Substrate Brief

Status: discussion draft; not committed
Date: 2026-05-05
Supersedes: `docs/research/harness-telemetry-trust-brief.md` as the next Socratic review input

This brief is a product framing artifact. It is not source-bound proof,
not product closure evidence, and not a trusted release claim.

## 1. Product Question

The buyer question is not:

> Did this one agent honestly do what it claims?

That question cannot be answered from agent-authored logs or local
telemetry alone.

The buyer question is:

> Is our AI-assisted delivery control plane gaining or losing control,
> and can we prove which parts of the evidence chain support that answer?

A CTO, CISO, or delivery executive needs to see whether agentic
development is drifting toward weaker evidence, weaker provenance,
late-attached traces, unapproved harnesses, manual exceptions, local-only
test claims, or missing telemetry at gates.

`sdp-trace` does not answer the final degradation question itself.
It supplies signed observations, explicit gaps, trace chains, and
verifier states that `sdp-gate`, `sdp-report`, or another external policy
consumer can use to answer it.

## 2. Product Layering

The product family should be split into three layers.

| Layer | Owns | Must Not Own |
| --- | --- | --- |
| `sdp-trace` | Evidence, provenance, trace, observer identity, event chains, signatures, witness scope, verifier states, missing telemetry facts. | Pass/fail policy, readiness, degradation verdicts, business thresholds, overrides. |
| `sdp-gate` | Expected evidence contracts, gate policy, required trust scope, waiver rules, gate verdicts. | Raw capture, low-level recorder implementation. |
| `sdp-report` | Time-series aggregation, control posture reports, trend analysis, CTO/COO/CISO views. | Cryptographic verification source of truth. |

Correct product statement:

> `sdp-trace` is a portable signed evidence substrate for AI-assisted
> delivery. It records what was observed, who observed it, what evidence
> was bound to it, what was missing, and what trust scope the verifier can
> honestly support.

Rejected product statement:

> `sdp-trace` proves that any arbitrary coding agent behaved honestly.

Reason: signatures prove signer identity and payload integrity. They do
not prove semantic truth unless the signer is an authorized observer for
the event type and the observation occurred at the correct trust boundary.

## 3. Industry Defaults To Reuse

Do not invent a bespoke cryptographic or provenance substrate where the
industry already has one.

Adopt:

- OpenTelemetry plus GenAI semantic conventions for runtime spans,
  context propagation, tool calls, model calls, agent spans, metrics, and
  eventual interoperability with Phoenix, LangSmith, Langfuse, and other
  observability systems.
- in-toto Statement as the attestation shape.
- DSSE as the signed envelope for statements and checkpoints.
- Sigstore keyless signing for CI/public trust profiles where available.
- GitHub Actions OIDC or equivalent CI identity for the first gate-grade
  witness.
- Rekor or a Sigstore bundle timestamp/transparency entry when public
  transparency is allowed.
- SLSA vocabulary for source provenance, build provenance, and
  verification summary attestations.
- Customer PKI or private Sigstore as an enterprise equivalent only when
  the identity policy, timestamp/freshness, audit log, protected source
  ref, and approval evidence are explicit.

Build `sdp-trace` above these standards, not beside them.

## 4. Core Data Model

The product has three first-class objects.

### Provenance

Provenance says who or what observed something.

Required dimensions:

- observer id;
- observer role;
- signer identity;
- authority policy;
- capture boundary;
- source commit or source snapshot;
- task/spec/plan identity;
- harness identity when observed;
- model identity when observed;
- CI or external witness identity when present.

### Evidence

Evidence is a concrete artifact or digest-bound fact.

Examples:

- source commit;
- diff or tree digest;
- command argv digest;
- stdout/stderr retention descriptor;
- test result artifact;
- Bazel BEP digest;
- model request id or digest;
- tool call id;
- review approval;
- waiver/override record;
- gate verdict;
- verifier result.

Evidence may be raw, sanitized, encrypted, externally referenced, or
digest-only. Retention state is part of the evidence, not an afterthought.

### Trace

Trace is the ordered chain of events and checkpoints.

Trace requirements:

- deterministic canonical JSON;
- SHA-256 event hash;
- previous event hash;
- monotonic sequence;
- correlation ids across harness, tool, shell, VCS, CI, and gateway;
- explicit late-attach boundaries;
- explicit supersession events when task/spec/plan changes;
- signed checkpoints at trust boundaries;
- verifier-visible missing observer states.

## 5. Observer Model

Observers are not equal. `sdp-trace` must record both the observation and
who was in a position to make it.

| Trust Scope | Observer | Can Support |
| --- | --- | --- |
| `agent_reported` | Agent-authored self-report | Intent, rationale, claimed plan, claimed completion. Never trust upgrade. |
| `local_observed` | Local recorder, wrapper, or daemon | Local reconstruction, command/file/test evidence with local forgery caveat. |
| `harness_observed` | Harness adapter | Task lifecycle, tool intent, model requested, subagent boundaries. |
| `gateway_observed` | LLM gateway or provider-side observer | Model call provenance, request ids, provider/model/cost/token metadata. |
| `vcs_observed` | Git/VCS provider | Commit, branch, PR, diff, review metadata. |
| `ci_witnessed` | CI identity with OIDC or equivalent | Gate-grade source/build/test/verifier evidence when policy matches. |
| `external_witnessed` | Transparency log, timestamp service, or independent witness | Anti-backfill anchor and stronger accountability profile. |
| `human_signed` | Human approver/reviewer | Decision, waiver, review disposition, manual exception. |

Authority rule:

> A signer can only support trust for event types it is authorized to
> observe.

If a local recorder signs a `ci_witnessed` event, verification fails.
If an agent signs a test result, it remains `agent_reported` unless an
authorized observer binds the actual test artifact.

## 6. Signing And Verification Profile

The product should sign observations and checkpoints, not vague truth
claims.

Minimum signing chain:

```text
event canonical JSON
  -> event_hash
  -> hash chain
  -> observer checkpoint statement
  -> DSSE envelope
  -> signer identity / Sigstore bundle / customer PKI equivalent
  -> verifier result
  -> optional gate verdict statement
  -> optional report input digest statement
```

Events do not all need individual signatures in v0. They must be
content-addressed and chained. Trust boundaries sign checkpoints and
verifier outputs.

Required signed checkpoint fields:

- run id;
- sequence range;
- chain head;
- previous checkpoint;
- source commit or snapshot digest;
- task/spec/plan digest;
- observer id;
- observer role;
- signer identity;
- signing profile;
- observed-at time;
- nonce;
- policy profile or expected evidence contract reference when applicable.

Verifier must check:

1. schema validity;
2. canonical hash integrity;
3. chain continuity;
4. payload digest match;
5. signature validity;
6. signer identity;
7. signer authority for observer role and event type;
8. source/task/run correlation;
9. checkpoint freshness before gate decision;
10. expected evidence contract completeness;
11. missing observers are explicitly represented;
12. redaction/retention state is compatible with the selected profile.

Verifier output must be three-axis:

```text
verdict:      pass | fail | cannot_verify | not_assessed
trust_scope:  agent_reported | local_observed | harness_observed |
              gateway_observed | vcs_observed | ci_witnessed |
              external_witnessed | human_signed
completeness: complete | partial | missing_telemetry | unknown
```

`pass + local_observed + partial` is not gate-grade.
`cannot_verify + ci_witnessed + missing_telemetry` is a valuable signed
gap, not a product failure.

## 7. Expected Evidence Contract

Missing telemetry is only meaningful when expectations are explicit.
`sdp-trace` should carry the facts; `sdp-gate` owns whether a gap blocks.

Example expected contract for an agent-attributed PR:

- task/spec locked before work;
- harness identity observed;
- model identity observed or explicitly `not_assessed`;
- tool/shell/file mutation events observed;
- source commit bound;
- test/build evidence observed;
- CI verifier result observed;
- gate checkpoint signed before merge;
- manual exception signed if any required evidence is waived.

Actual evidence may be weaker:

- only local shell events;
- no model gateway event;
- no harness adapter event;
- CI test missing;
- manual exception present.

`sdp-trace` should emit this as signed facts and signed gaps. It should
not silently collapse the run into green/red.

## 8. Demonstration Strategy

The demo must show trust boundaries, not only a happy path.

### Demo 1: Local Trace Is Useful But Not Gate-Grade

Run a command/test sequence through a local recorder.

Expected output:

```text
verdict: pass
trust_scope: local_observed
completeness: partial
gate_usable: false
reason: structurally valid local trace without CI/external witness
```

### Demo 2: CI-Witnessed Gate

Run the same evidence path in CI. CI verifies source commit, chain head,
test/build evidence, expected contract, and signs the verifier result.

Expected output:

```text
verdict: pass
trust_scope: ci_witnessed
completeness: complete
gate_usable: true
```

### Demo 3: Missing Observer

Let the PR reach CI without harness/model telemetry.

Expected output:

```text
verdict: cannot_verify
trust_scope: ci_witnessed
completeness: missing_telemetry
missing:
- harness_observed
- model_identity_observed
```

The point is not to fail dramatically. The point is to make absence
impossible to hide.

### Demo 4: Tamper Attack

Modify an event payload, remove an event, or reorder the chain.

Expected output:

```text
verdict: fail
reason: event_hash mismatch or prev_event_hash mismatch
```

### Demo 5: Post-Hoc Local Fabrication

Generate a nice local trace after the work is already done.

Expected output:

```text
verdict: cannot_verify
trust_scope: local_observed
completeness: partial
reason: no checkpoint witnessed before gate decision
```

## 9. First Build Slice

The first credible implementation slice should not be full agent
governance. It should be a narrow signed evidence substrate.

Build:

- canonical event schema for agentic SDLC observations;
- observer checkpoint schema as in-toto Statement predicate;
- expected evidence contract schema;
- verifier with three-axis output;
- local recorder fixture path;
- CI-witnessed fixture path;
- Sigstore/GitHub OIDC design path, with local simulated CI only clearly
  labeled if real CI is not yet available;
- CTO-readable report for one run showing observed evidence, trust scope,
  completeness, missing telemetry, and gate usability.

Defer:

- full LLM gateway integration;
- raw prompt/response capture;
- multi-harness adapters;
- remote witness service;
- file watcher attribution;
- full degradation dashboard;
- automated blocking outside CI/policy consumers.

Reason: if the event, observer, signing, and verifier model is wrong,
more adapters only create more unauditable telemetry.

## 10. CTO Report Shape

The report should answer concrete questions:

- What task/spec/plan was locked?
- Which source commit is bound?
- Which observer roles were present?
- Which required observers were missing?
- Which commands/tests/builds were observed?
- Which evidence is local-only?
- Which evidence is CI-witnessed?
- Was the chain anchored before the gate?
- Did the verifier fail because of tampering, missing evidence, weak
  trust scope, redaction, or policy mismatch?
- Can an external gate use this trace?

It must not show an opaque health score.

For time-series `sdp-report`, useful degradation indicators include:

- share of PRs with `missing_telemetry`;
- share of test claims that are `local_observed` only;
- share of gate decisions with manual exceptions;
- share of runs missing model or harness provenance;
- share of late-attached traces;
- unapproved harness/model/tool paths;
- correlation between weak evidence states and reverts, CI failures, or
  production incidents.

These are report-layer metrics, not native `sdp-trace` verdicts.

## 11. Known Product Risks

Critical risks:

- Local hash chains can be fabricated after the fact without a witness.
- Local signing is not audit-grade when the developer controls the host.
- CI can accidentally sign unverified telemetry unless verifier result,
  source commit, chain head, and policy context are covered.
- Gateway telemetry proves model calls, not local file/test behavior.
- Local recorder telemetry proves local behavior, not model provenance.
- Missing telemetry is a governance signal, not proof of deletion or
  malicious bypass.
- Raw prompts, responses, argv, stdout, stderr, and file paths may leak
  confidential data.
- Expected evidence contracts can become too strict and create adoption
  friction before teams have adapters installed.

Product discipline:

- never sell local-only as audit-grade;
- never treat agent self-report as a trust upgrade;
- never let schema validity imply production trust;
- never use pass/fail without trust scope and completeness;
- never hide `not_assessed` gaps in prose.

## 12. Socratic Review Task

Review this brief adversarially.

Return:

```text
VERDICT: CHANGES_REQUIRED | ACCEPTABLE_WITH_GAPS | REJECTED

Critical blockers:
- ...

Major gaps:
- ...

False assumptions:
- ...

Standards we should reuse instead of building:
- ...

Minimum viable product correction:
- ...

Questions before implementation:
- ...

Attack scenarios not covered:
- ...

Demo changes required:
- ...
```

Tie every finding to at least one area:

- product layering;
- CTO usefulness;
- evidence model;
- provenance model;
- trace model;
- observer authority;
- signing and verification;
- expected evidence contracts;
- CI/gate anchoring;
- privacy and retention;
- adoption and DX;
- demo credibility.

Do not reward architectural elegance if it weakens UX, DX, or trust
clarity. The acceptable target is an 8/10 brief that is good enough to
build, not an endlessly generalized trust architecture.
