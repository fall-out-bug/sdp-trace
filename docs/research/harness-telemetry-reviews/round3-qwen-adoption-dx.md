# Round 3 — Adoption/DX Review

**Source text**: `docs/research/harness-telemetry-trust-brief.md`
**Reviewer role**: CPO + CTO buyer evaluating purchase/adoption decision; refuses harness replacement mandate; assumes heterogeneous team using pi, OpenCode, Kilo, and internal CI wrappers concurrently
**Prior rounds**: 4 Round 1 reviews (all CHANGES_REQUIRED), 3 Round 2 forgery assessments (0 of 25 attacks prevented, 2 detected against naive adversary only)

---

## VERDICT: REJECTED (unpurchasable in current form)

Not because the problem is wrong. The problem is correct and urgent. The product described does not exist yet. What exists is a threat inventory with architecture aspirations. I cannot ask my team to adopt a "trust layer" that defines no integration contract, emits no queryable output, and whose verifier correctly detects only naive hash tampering.

I would fund this effort. I would not deploy it.

---

## Adoption Blockers

| # | Blocker | Area | Why Fatal |
|---|---------|------|-----------|
| B1 | "Harness-neutral" requires harness cooperation but defines none | Integration | The brief states sdp-trace must not be a mandatory harness, but every layer of value (task boundary, model identity, tool intent, subagent correlation) requires harness-specific API integration with no adapter contract defined. Teams cannot adopt "attach recorder" when the recorder has no harness API to attach to. |
| B2 | No installable artifact exists | Adoption | No SDK, no CLI, no daemon, no schema, no package. The brief describes 4 capture layers and 4 signing levels but provides zero executable components. A team cannot "try" this in any form. |
| B3 | Observe-only mode provides zero trust signal | Adoption | Section 7 promises "observe-only: no blocking; show gaps." Round 2 red team confirms the verifier emits `local_only` + `pass` for fabricated post-hoc chains. An observe-only deploy will show clean runs for both honest and dishonest agents. The CTO dashboard will look green regardless. This destroys product credibility on first exposure. |
| B4 | Gateway privacy risk blocks enterprise procurement | Integration | Section 3A proposes prompt/response digests at LLM gateway. Enterprises will not route proprietary codebase prompts through a new telemetry proxy without a documented privacy boundary, PII scrubbing model, and data residency guarantee. None exist. Legal review will kill this before engineering evaluates it. |
| B5 | Local signing key management is unsolved and unsolvable on unmanaged laptops | Signing | Round 1 reviews unanimously flagged this. Agent and recorder share user context on macOS/Linux dev machines. No hardware-bound attestation, no TEE, no remote signing service is specified. Local signatures prove only "something with user access signed this." This is not a trust boundary and must not be sold as one. |
| B6 | The CTO query surface has no implementation | CTO usefulness | Section 8 lists 13 questions the CTO wants answered. Zero have a defined data source, index model, or query mechanism. The brief offers no CLI query, no API, no dashboard, and no report format. A CTO cannot evaluate a product whose output they cannot see. |
| B7 | Completeness contract absence makes all verdicts meaningless | Anti-forgery | Round 2 established: verifier cannot detect missing telemetry without an expected-event manifest. Any run — honest, bypassed, or fabricated — can emit `pass` if the chain it produces is internally consistent. The core value proposition ("can the CTO see whether AI-assisted delivery is degrading team control?") requires detecting absence. The system cannot. |

---

## First-Week Value (what a team could actually get)

Given the current design, here is the honest first-week experience:

**If a team tries to adopt today:**

1. Install nothing (no package exists).
2. Read the brief (good threat model, no executable component).
3. Manually wrap shell commands with `scripts/verify.sh` (current repo baseline).
4. Get local structural evidence for committed code (already exists in this repo).
5. Cannot correlate model calls, tool calls, or subagent events.
6. Cannot distinguish live from post-hoc telemetry.
7. Cannot answer any CTO query from Section 8.

**Honest first-week value: structural evidence for committed file changes, using existing repo tools.** This is useful but is not a "trust layer." It is a git diff with a hash.

**What the brief implies for first-week value:** "Observe-only: attach recorder to existing harness; no blocking; show gaps." This requires a recorder that does not exist. The implied promise is false.

---

## Required Minimum Integration (what must exist before any team can adopt)

The brief lists 7 integration points. The minimum viable integration is not "all of them" or "layered." It is **one layer that works independently** and proves the model.

**Recommendation: shell wrapper first.**

Rationale:
- Works across every harness (pi, OpenCode, Kilo, CI, custom scripts) with zero harness-specific code.
- Captures the evidence a CTO needs most: what commands ran, did tests pass, did builds succeed.
- Strongest local evidence for build/test integrity (Section 3 and 3A both rate this "Strong").
- Lowest privacy risk (no prompts, no model responses, no source content).
- Implementable as a single portable binary or POSIX script.

**What the shell wrapper must emit (minimum schema):**

```json
{
  "event_type": "shell_command",
  "event_id": "uuid",
  "session_id": "uuid",
  "timestamp_start": "ISO8601",
  "timestamp_end": "ISO8601",
  "cwd": "string",
  "argv": ["string"],
  "argv_digest": "sha256",
  "exit_code": "integer",
  "stdout_digest": "sha256",
  "stderr_digest": "sha256",
  "source": "shell_wrapper",
  "evidence_role": "observed",
  "prev_hash": "sha256"
}
```

This single layer, working reliably, would:
- Answer "what commands ran" and "which tests/builds ran" from the CTO query surface.
- Provide build/test evidence strength sufficient for local gate decisions.
- Work today, with no harness cooperation, no gateway, no CI changes.

**Everything else — harness plugin, gateway, file watcher, CI witness, remote witness — ships in later versions.**

The brief's "layered, not either/or" thesis is architecturally correct but commercially wrong. A product cannot ship "layers." It ships one thing that works, then adds layers. The shell wrapper is the one thing.

---

## Unsafe Telemetry (what must never be captured or must be transformed)

| Telemetry | Risk | Required Action |
|-----------|------|-----------------|
| Full LLM prompts | Proprietary code, secrets, PII, customer data in prompt context | **Never capture full prompts.** Hash-only, with salted SHA-256 minimum, with documented minimum entropy floor. |
| Full LLM responses | Same as prompts; responses may contain generated secrets or keys | **Never capture full responses.** Response hash only. |
| Environment variables from shell events | API keys, tokens, credentials in `env` | **Redact known secret patterns** (KEY, TOKEN, SECRET, PASSWORD) before capture. Document the redaction policy. |
| File content in file watcher events | Source code, configuration, credentials | **Capture path + diff digest only.** Never capture file content in telemetry events. |
| Stdout/stderr from shell events | May contain test output (safe) or secrets (unsafe) | **Digest-only by default.** Opt-in full capture for specific known-safe commands (e.g., `npm test`). |

**Product rule: telemetry must never be a data exfiltration vector.** If the CTO's security team discovers that an agent telemetry system stored proprietary prompts in a queryable database, the product is dead.

---

## CTO Query Surface (honest assessment)

Section 8 lists 13 questions. Here is which can be answered by the shell wrapper alone:

| CTO Question | Shell Wrapper | Full Design |
|---|:---:|:---:|
| What task was locked? | ❌ | Partial (harness adapter + task hash checkpoint required) |
| Did the task change? | ❌ | Partial (task immutability mechanism required, not specified) |
| Which model and harness ran? | ❌ (harness adapter / gateway required) | Partial |
| What commands ran? | ✅ | ✅ |
| Which files changed? | ❌ (file watcher required) | ✅ |
| Were changes inside allowed scope? | ❌ | Partial (scope checkpoint required) |
| Which tests/builds ran? | ✅ (by parsing argv) | ✅ |
| What evidence is missing? | ❌ (completeness contract required) | ❌ (completeness contract not specified anywhere) |
| Was telemetry captured live or attached late? | ❌ | ❌ (liveness proof not specified anywhere) |
| Is this local-only or witnessed? | Partially (if witness protocol exists) | Partially (witness protocol not specified) |
| Can this run support a gate decision? | ❌ | Partially (CI signing recipe required) |
| Where did agent claim something unsupported? | ❌ | ❌ (agent-vs-recorder correlation not specified) |
| Is the team drifting into unapproved harnesses? | ❌ | ❌ (harness registry not specified) |

**Scorecard: shell wrapper answers 2 of 13. Full design, if implemented as described, answers maybe 6 of 13. Seven questions require mechanisms that are listed as missing or undefined.**

**Product implication:** The CTO value proposition overclaims by roughly 2×. The brief presents 13 questions as "required" when fewer than half are supportable even in the aspirational design. The honest query surface for v1 is 4-5 questions, all evidence-focused:
1. What commands ran?
2. Which tests/builds ran and what was the exit state?
3. Is the command chain internally consistent?
4. Is this chain local-only or witnessed?
5. For witnessed runs: was the chain head anchored before the gate decision?

**Cut the other 8 questions from the v1 brief.** They are roadmap items, not current product claims.

---

## Demo Scope Cuts (what to cut to ship a credible demo)

The current brief scope:
- 7 capture layers
- 4 signing levels
- 8 verifier states (conflated across 3 axes)
- 13 CTO queries
- 13 threat categories with zero mitigations
- 11 attack scenarios from red team, 2 detected

**Demo scope — cut to the following:**

| Keep | Cut | Rationale |
|------|-----|-----------|
| Shell wrapper capture | Harness plugin | No harness adapter contract exists; shell wrapper is harness-neutral by nature |
| Local hash chain | Remote witness protocol | Witness protocol is undefined; demo the local chain first |
| Verifier: pass/fail/cannot_verify with hash integrity check | Gateway capture, CI signing, external witness | All require protocols that don't exist |
| 3 CTO queries: commands, tests, chain integrity | Remaining 10 CTO queries | Cannot demo what doesn't exist |
| One attack demonstrated: event mutation detected by hash chain | Remaining 24 attacks | Demonstrate that the one mechanism that works actually works |
| Schema: shell_command event (JSON) | All other event types | One working event beats 7 undefined events |

**Demo claim language:**
> "sdp-trace v0.1 captures shell commands during agent sessions, builds a tamper-evident hash chain, and verifies chain integrity post-session. It detects event mutation, deletion, and reordering within the captured chain. It does not yet capture model provenance, file mutations, or external witness anchors. Local-only chains cannot resist host-level forgery."

This is honest. It is also dramatically less impressive than the current brief. **It is the only claim the design can currently support.**

---

## Non-Negotiable Product Changes

1. **Drop "harness-neutral trust layer" from the tagline until an adapter contract exists.**
   Current reality: sdp-trace is a structural evidence verifier for committed code and (aspirational) a shell command recorder. Call it that. "Trust layer" implies a boundary that cannot be enforced without at least one working integration contract.

2. **Split the product into telemetry layer (v0.x) and trust layer (v1.x).**
   The adoption ladder (Section 7) acknowledges this implicitly but the brief markets the whole stack as one product. Stage 1-2 are telemetry: capture and show. Stage 3-4 are trust: verify and enforce. Selling "trust" at stage 1 is overclaiming and will destroy credibility with enterprise buyers who check.

3. **Define the verifier state model as 3 axes before any further feature work.**
   Round 1 (Mimo, GLM) and Round 2 (all three) unanimously required: verdict × trust-scope × completeness. The current 8-state single enum is ambiguous and produces false positives for fabricated chains. This is not an implementation detail; it is the core product interface. Every CTO query output depends on it.

4. **Write and publish one JSON schema for one event type.**
   The shell wrapper event schema above is the recommendation. Publish it in `schema/`. Validate it with `jq empty`. Make it the baseline for all subsequent event types. Without a schema, the "harness-neutral contract" is prose, not a specification, and cannot be implemented by any harness vendor.

5. **Write a completeness contract stub.**
   Even if minimal: "a complete run must contain at least one observed event from a recorder source." Without this, the verifier emits `pass` for empty or fabricated runs, and the product is demonstrably unreliable.

6. **Remove LLM gateway prompt/response capture from Slice 2.**
   The privacy risk blocks enterprise procurement. There is no mitigation defined. Defend the decision to cut it: "Gateway provenance is valuable but prompt/response capture creates unacceptable data exfiltration risk without a documented sanitization pipeline. Gateway integration is deferred to a follow-up brief with explicit privacy model."

7. **State explicitly: local-only telemetry is not audit-grade.**
   The brief says this in Section 5 but buries it. It must be a top-level product disclaimer, displayed in every verifier output for local-only runs: "This chain was verified for local structural integrity only. It cannot resist host-level forgery and is not suitable for audit or gate decisions."

8. **Define the managed mode enforcement mechanism or remove managed mode.**
   "Managed mode: verifier fails or blocks" is the brief's most aspirational claim with the least implementation detail. If the verifier cannot block the harness (it cannot, post-hoc), then managed mode is a policy commitment, not a technical enforcement. State that clearly: "Managed mode is an organizational agreement to reject runs with missing_telemetry. The verifier does not block execution."

---

## Summary

The problem is real, the architecture direction is correct, and the threat model is comprehensive. But the product does not yet exist as something a team can adopt. The brief describes a destination — a 4-layer, 4-signing-level, externally-witnessed telemetry system — without defining the first step as a shippable artifact.

**What a CTO buyer needs:**
- One working recorder (shell wrapper).
- One event schema (shell_command).
- One verifier output format (3-axis state).
- One honest claim: "tamper-evident local command telemetry."
- One demonstrated attack that it prevents.

**What the brief currently offers:**
- 7 integration points, none implemented.
- 8 verifier states, conflated across 3 axes.
- 13 CTO queries, less than half supportable.
- 13 threats named, zero mitigated.
- "Trust layer" branding for a system whose red-team record is 0 of 25 attacks prevented.

**Close the gap between claim and capability before asking teams to adopt.** The shell wrapper is the bridge.
