# Round 5: Agentic SDLC Evidence Substrate Socratic Consolidation

Status: discussion draft; not committed
Date: 2026-05-05

Inputs:

- `archive/research/agentic-sdlc-evidence-substrate-brief.md`
- `archive/research/harness-telemetry-reviews/persona-01-cto-buyer.md`
- `archive/research/harness-telemetry-reviews/persona-02-platform-harness-owner.md`
- `archive/research/harness-telemetry-reviews/persona-03-ciso-adversarial-trust.md`
- `archive/research/harness-telemetry-reviews/persona-04-staff-engineer-dx-skeptic.md`
- `archive/research/harness-telemetry-reviews/persona-05-compliance-forensics-lead.md`
- `archive/research/harness-telemetry-reviews/round5-cto-buyer-minimax.md`
- `archive/research/harness-telemetry-reviews/round5-platform-harness-owner-glm.md`
- `archive/research/harness-telemetry-reviews/round5-ciso-adversarial-trust-grok.md`
- `archive/research/harness-telemetry-reviews/round5-staff-engineer-dx-kimi.md`
- `archive/research/harness-telemetry-reviews/round5-compliance-forensics-qwen.md`

Run note: the first CISO attempt used DeepSeek through `pi` and produced
no output before hanging. It was stopped and replaced with xAI/Grok. This
consolidation treats the Grok artifact as the CISO review input.

This file is a human consolidation of Socratic persona outputs. It is
not source-bound proof, not product closure evidence, and not a trusted
release claim.

## Overall Verdict

CHANGES_REQUIRED.

The new brief is materially stronger than Round 4 because it separates
`sdp-trace`, `sdp-gate`, and `sdp-report`, adopts in-toto/DSSE/Sigstore
instead of inventing signatures, and makes `missing_telemetry` first
class.

It is still not good enough to build because the first integration
boundary is not defined. The brief says "local recorder", "adapter",
"expected evidence contract", and "query surface", but does not specify
the minimum executable contract that makes those terms real.

The next revision must shift from trust architecture vocabulary to a
concrete v0 capture and investigation contract.

## Socratic Progression

1. CTO buyer: the brief still lacks day-one value for a company that will
   not change harnesses.
2. Platform owner: the brief has no concrete adapter/recorder interface,
   so there is no enforceable capture point.
3. CISO reviewer: the brief still lets local chains look stronger than
   they are unless key separation, witness, replay, and signer authority
   are explicit.
4. Staff engineer: privacy, latency, emergency workflow, and false
   positive explainability are not optional DX details; they are adoption
   gates.
5. Forensics lead: signed events are insufficient without retention,
   redaction audit, replay/query, and missing evidence tables.

## Accepted Product Corrections

| Finding | Raised By | Severity | Disposition | Required Change |
| --- | --- | --- | --- | --- |
| No zero-instrumentation or near-zero-instrumentation day-one path. | CTO | critical | accepted | Add Demo 0 and a v0 local observation path that works without harness rewrites, even if it only yields `local_observed + partial`. |
| Adapter contract is missing. | Platform | critical | accepted | Add a minimum adapter interface with event types, required fields, lifecycle, observer roles, and failure modes. |
| Local recorder mechanism is undefined. | Platform, Staff Engineer | critical | accepted | Replace "local recorder fixture path" with a concrete v0 mechanism: command wrapper, daemon, or sidecar. |
| Privacy/redaction is not a policy. | CTO, Staff Engineer, Forensics | critical | accepted | Define default capture profile before implementation: digest-only/redacted by default, raw capture opt-in and scope-bound. |
| Redaction must happen before persistence. | Staff Engineer, Forensics | critical | accepted | Treat redaction as a capture-boundary filter plus signed redaction manifest, not post-hoc metadata. |
| Key separation is underspecified. | CISO | critical | accepted | State which keys are outside agent/workspace control and which trust scopes cannot be signed locally. |
| External witness is mandatory for gate-grade claims. | CISO, CTO, Platform | critical | accepted | Keep local-only as reconstruction evidence; require CI/OIDC or external witness for gate-grade trust. |
| Post-hoc fabrication demo lacks detection mechanism. | CTO, Platform, Forensics | critical | accepted | Demo 5 must show the exact check: missing pre-gate checkpoint, VCS/timestamp mismatch, nonce mismatch, or external witness absence. |
| Expected evidence contract is prose, not a contract. | CTO, Platform, CISO | critical | accepted | Add schema shape, storage location, version lock, authority, and immutability/supersession rules. |
| CI can sign unverified or co-located agent telemetry. | CTO, CISO | critical | accepted | CI witness must sign verifier result, source commit, chain head, policy profile, and verifier identity; co-located agent-in-CI topology must be downgraded or explicitly modeled. |
| Missing telemetry needs a table, not a list. | Forensics | major | accepted | Verifier/report output must include a MissingEvidenceTable with expected role, absent evidence, sequence/window, reason, and policy reference. |
| Query/replay surface is missing. | Forensics | critical | accepted | Add a minimal forensic query contract before claiming investigation usefulness. |
| Emergency path is absent. | Staff Engineer | critical | accepted | Add `policy_override_requested` / waiver event with `human_signed` reason, scope, time, and evidence snapshot. |
| Developer explainability is missing. | Staff Engineer, Platform | major | accepted | Add `sdp-trace explain` style output: event-level cause, missing observer, authority mismatch, chain break, redaction issue. |
| Latency/offline budget is unspecified. | Staff Engineer | major | accepted | Define local recorder as async/offline-capable with a p99 overhead target and no network dependency in the inner loop. |
| Retention is underspecified. | Forensics, Staff Engineer | critical | accepted | Every evidence reference needs retention descriptor: mode, expiry, archival location/ref, redaction manifest digest, replayability. |
| Harness identity can be spoofed. | CTO, CISO | major | accepted | Add harness identity attestation or classify unverified harness identity as `agent_reported` / `local_observed`, not `harness_observed`. |
| Retroactive contract weakening is possible. | CTO | critical | accepted | Expected evidence contract must be locked before run/gate; changes require supersession events and cannot silently reclassify old runs. |

## Rejected Or Modified Suggestions

| Suggestion | Source | Disposition | Reason |
| --- | --- | --- | --- |
| Defer signing chain until observer capture and retention are validated. | CTO | modified | Do not build production Sigstore first, but keep checkpoint/signature shape in v0 fixtures. Signing semantics drive trust boundaries; implementation can start local/simulated as long as outputs are labeled non-gate-grade. |
| Make Sigstore/Rekor mandatory for all checkpoints. | CISO | modified | Correct for production/public gate-grade trust, too strict for private/air-gapped or local reconstruction. The brief should require an accepted external trust profile, not public Rekor in every environment. |
| Use W3C Verifiable Credentials for observer authority. | Forensics | defer | Useful candidate, but not needed for first build slice. Authority policy can start as a simple verifier config as long as it is signed/versioned and not hidden. |
| Use CloudEvents or OTel Logs instead of bespoke canonical JSON. | Forensics | investigate | Strong DX/interoperability argument. Need a design decision: either canonical events are OTel-compatible, or the brief must justify why OTel spans/logs cannot be the canonical event shape. |

## Required Brief Revision

The next brief should add these sections before any implementation plan:

1. **V0 Capture Boundary**
   - Choose command wrapper, daemon, sidecar, or another concrete local
     recorder.
   - Define how it starts, attaches, heartbeats, flushes, and closes.
   - Define what is visible when the agent bypasses it.

2. **Adapter Interface**
   - Minimum event set:
     `recorder_attached`, `run_started`, `task_locked`,
     `expected_evidence_contract_locked`, `tool_call_observed`,
     `command_started`, `command_finished`, `file_mutation_observed`,
     `test_observed`, `policy_override_requested`,
     `requirement_superseded`, `redaction_applied`, `run_closed`.
   - For each event: required fields, valid observer roles, allowed trust
     scopes, correlation ids, retention requirements.

3. **Expected Evidence Contract**
   - Schema shape.
   - Storage and lock point.
   - Version/supersession rules.
   - Who can author/approve it.
   - How missing observers become MissingEvidenceTable rows.

4. **Privacy And Retention Profile**
   - Default digest-only or redacted capture.
   - No raw prompt/source/stdout/stderr by default.
   - Pre-write redaction.
   - Redaction manifest.
   - Evidence retention descriptor.
   - Replayability state.

5. **Key And Authority Policy**
   - Which keys live where.
   - Which signers can sign which observer roles.
   - What happens when signer authority is missing.
   - Why local signing never upgrades to gate-grade.

6. **CI/Witness Profile**
   - CI signs verifier result, source commit, chain head, expected
     evidence contract, verifier version, policy profile, and timestamp.
   - Co-located agent-in-CI topology must be explicitly downgraded or
     split by execution identity.

7. **Forensic Query Contract**
   - Query names and expected answers.
   - MissingEvidenceTable.
   - Timeline reconstruction.
   - Evidence replayability.
   - Links to commit, PR, CI, and signed checkpoints.

8. **Developer Explainability**
   - `explain` output.
   - dry-run/preview output.
   - false-positive debugging.
   - emergency override flow.
   - latency/offline budget.

## Demo Set Required By Round 5

The next demo plan should include:

1. **Demo 0: No-Harness-Change Local Observation**
   - Run any agent/harness through the v0 local recorder.
   - Output `local_observed + partial`.
   - Show what was captured and what remains missing.

2. **Demo 1: Local Trace With Dry Run And Redaction**
   - Show `--dry-run`.
   - Show default redaction/digest-only behavior.
   - Show local trace is not gate-grade.

3. **Demo 2: CI-Witnessed Gate**
   - CI signs verifier result, source commit, chain head, expected
     contract, verifier version, and policy profile.
   - Output `ci_witnessed + complete` only if expected evidence is
     satisfied.

4. **Demo 3: Missing Observer With MissingEvidenceTable**
   - Show absent harness/model/gateway/local events as explicit rows.
   - Include unmanaged/degraded mode semantics.

5. **Demo 4: Tamper And Chain Break**
   - Mutate, delete, and reorder events.
   - Show event-level verifier explanation.

6. **Demo 5: Post-Hoc Fabrication / Replay**
   - Recreate a valid-looking local trace after the fact or replay an old
     trace against a new PR.
   - Show failure/downgrade by missing pre-gate witness, VCS mismatch,
     timestamp mismatch, nonce mismatch, or explicit absence of external
     anchor.

7. **Demo 6: Emergency Override**
   - Emit `policy_override_requested`.
   - Show `human_signed + partial`, not hidden bypass.

8. **Demo 7: Forensic Query**
   - Given a trace, answer what ran, which files changed, which tests
     were evidence, what was redacted, and what cannot be assessed.

9. **Demo 8: Latency / Trace Bloat Guard**
   - Run a high-frequency local loop.
   - Show bounded overhead, async writing, and storage limits.

## Product Decision

The product should keep the Round 5 thesis:

> `sdp-trace` records signed observations and signed gaps so external
> gate/report consumers can reason about AI SDLC control posture.

But the next revision must stop treating "recorder", "adapter",
"contract", and "query" as obvious nouns. They are the product. Until
they are specified, the brief remains a strong strategy memo and a weak
implementation target.

## Recommended Next Step

Revise `agentic-sdlc-evidence-substrate-brief.md` into a buildable v0
brief with:

- concrete local recorder mechanism;
- adapter/event taxonomy;
- expected evidence contract;
- pre-write privacy/redaction policy;
- CI witness authority profile;
- forensic query output;
- developer dry-run/explain output.

Then run one more Socratic round focused only on whether that v0 is
small enough to build without losing the CTO value.
