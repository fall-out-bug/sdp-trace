# Trust and Evidence Semantics Review: Product Contract v0

**Reviewer**: Trust and evidence semantics plane
**Targets**: `spec.md`, `plan.md`, `example.md`, `traceability.md`, `tasks.md`
**Standard**: Evidence state preservation, theater code integrity, signed-attestation discipline, missing state handling, packet self-reference risk

---

## 1. Verdict

**REVISE_BEFORE_USER_REVIEW**

Product Contract v0 has a coherent core structure. The row taxonomy, evidence state vocabulary, and backlog gate are sound. However, three critical issues and two major issues prevent it from being ready:

1. `PC-THEATER` row state semantics are inverted-treating expected theater findings as a `fail` condition
2. `PC-VERIFICATION` conflates harness-observed evidence with unverifiable state
3. Evidence state definitions for `missing_telemetry` and `not_integrated` are absent from `spec.md`
4. The traceability matrix overclaims `PC-THEATER` and `PC-AUTHORITY` readiness using implementation artifacts rather than theater reason codes
5. Deferred P1 theater codes create an unresolved theater detection completeness gap

---

## 2. What Works

| area | what's solid | cited lines |
| --- | --- | --- |
| Row taxonomy | 11 row ids cover the key buyer questions without overreaching | `spec.md:67-79` |
| Evidence state vocabulary | Core states (`pass`, `partial`, `fail`, `not_assessed`, `cannot_verify`) are well-defined and correctly distinguished | `spec.md:81-93` |
| Source strength vs trust | Clear separation: source strength describes origin, not credibility | `spec.md:94-105` |
| Signed attestation discipline | FR-012 correctly bans signing as a shortcut to upgrade missing evidence | `spec.md:247` |
| Authority non-verdict | Packet names decision owners but does not approve merge/release/compliance | `spec.md:75`, `plan.md:45` |
| Backlog gate fields | `packet_rows`, `evidence_surface`, `closure_state`, `buyer_effect`, `non_goal` are enforceable | `spec.md:148-158` |
| Example theater table | Correctly shows `scope_theater: pass` and `ci_theater: not_assessed` for unavailable CI | `example.md:52-57` |
| Deferred P1 findings visibility | FR-009 correctly requires deferred theater codes to be visible in `PC-RESIDUAL-GAPS` | `spec.md:228-229` |

---

## 3. Blocking Findings

| id | severity | location | finding | why it matters | fix |
| --- | --- | --- | --- | --- | --- |
| TE-001 | **critical** | `spec.md:81-93` | Evidence states `missing_telemetry` and `not_integrated` are listed but never defined. `missing_telemetry` appears nowhere in examples; `not_integrated` appears only in the state list. | A buyer or engineer cannot correctly assign these states without a definition. `missing_telemetry` likely means "no telemetry was produced"; `not_integrated` likely means "capability exists but was not wired for this run"-but these are guesses. | Add explicit definitions to `spec.md:81-93`. Example: `missing_telemetry` = "required telemetry was not produced or retained"; `not_integrated` = "the capability exists in the substrate but was not connected or invoked for this evidence run." |
| TE-002 | **critical** | `example.md:42` | `PC-VERIFICATION` is `cannot_verify` when harness-observed evidence exists. The row says "harness observed a verification command" but assigns `cannot_verify`. If harness evidence exists, the state should be `partial` or `fail` based on whether the evidence is sufficient-not `cannot_verify`. | `cannot_verify` means the packet has no information, not that it has partial information. Using `cannot_verify` when partial evidence exists conflates the two and understates what the packet does contain. This affects how buyers interpret the row. | Separate harness-observed evidence from retained CI witness. If harness observed test execution but no CI artifact is retained: `partial` with theater finding `ci_theater`. If no harness observation exists: `cannot_verify`. |
| TE-003 | **critical** | `example.md:45` and `spec.md:108-117` | `PC-THEATER` row state is `fail`, but theater findings are expected output when evidence is incomplete. The example lists `agent_claimed_verification` and `unbound_intent` as theater findings-the exact evidence the packet is designed to surface-then marks the theater row `fail`. This inverts the semantics. | A `fail` row state implies the packet is broken or the assessment failed. But theater findings ARE the expected output of `PC-THEATER`. The theater findings table already correctly shows `scope_theater: pass` for no-overclaim and `ci_theater: not_assessed` for unavailable CI. The row state contradicts the findings table. | Change `PC-THEATER` row state semantics. If theater assessment was run and findings exist: `partial`. If theater assessment was run and no findings: `pass`. If theater assessment cannot be run due to missing inputs: `not_assessed`. Remove `fail` from valid theater row states, or redefine it as "theater assessment was run and found critical evidence gaps requiring deferral." |
| TE-004 | **major** | `traceability.md:18` | `PC-THEATER` claims "partial" coverage with substrate inputs: "existing negative fixtures, managed harness failures, adapter capture overclaim failures." These are implementation artifacts, not theater reason codes. | Theater reason codes are the first-class P0 output: `agent_claimed_verification`, `unbound_intent`, `ci_theater`, `scope_theater`. Current substrate provides test cases and failure patterns-not the theater derivation logic itself. This overclaims readiness and masks that theater reason-code derivation is not yet implemented. | Revise `PC-THEATER` coverage to: "No theater reason-code derivation. Need documented rules mapping evidence conditions to P0 theater codes and tests for each code." |
| TE-005 | **major** | `traceability.md:17` | `PC-AUTHORITY` claims "Good substrate" for `PC-AUTHORITY`, but the same row says "Need packet projection that says authority state without making policy decisions." | "Good substrate" is a readiness claim. The identified gap-packet projection without policy verdicts-is not a trivial translation. If that gap exists, coverage is not "good"-it's partial with a known shape. | Change `PC-AUTHORITY` coverage from "Good substrate" to "Partial: authority envelope schema exists, but buyer-facing projection without policy verdicts is not defined." |

---

## 4. Non-Blocking Concerns

| id | severity | location | concern | reasoning |
| --- | --- | --- | --- | --- |
| TE-006 | minor | `spec.md:108-117` | P1 theater codes are deferred without explaining how P0 theater detection remains meaningful without `actor_laundering` and `review_theater`. | Actor laundering and review theater are among the highest-risk trust failures in agent supply chains. Deferring them to P1 leaves a significant detection gap that is not explicitly named as residual risk. |
| TE-007 | minor | `example.md:52-57` | Theater finding `agent_claimed_verification` severity is `major` even in a context where harness-observed evidence exists and CI is unavailable. | In a no-CI local baseline, the agent's own test claims may be the strongest available evidence signal. Severity classification should consider the evidence profile, not just the finding type. |
| TE-008 | minor | `spec.md:94-102` | Source strength descriptors are well-explained but the evidence bundle attachment format is not defined. | Rows cite `evidence refs` but the attachment schema (how digests, bundles, and retained forms are stored) is not specified. The backlog gate and CTO readability depend on this. |
| TE-009 | minor | `traceability.md:57-65` | Known gaps section lists "no generated packet command," "no packet schema," "no theater derivation" but does not map these to closure states with concrete exit criteria. | Gaps are named but the transition from gap to evidence is abstract. For `PC-THEATER` theater derivation, the closure criterion should name the exact test input-output pairs. |

---

## 5. Missing Evidence / `not_assessed` Areas

| area | current state | what would close it |
| --- | --- | --- |
| `missing_telemetry` and `not_integrated` state definitions | `not_assessed` | Explicit definitions in spec.md with example usage |
| Theater reason-code derivation rules | `not_assessed` | Documented mapping: evidence condition -> reason code, with test fixtures for each P0 code |
| Packet self-reference integrity | `not_assessed` | Explicit acknowledgment in spec that the packet itself requires binding to source evidence, or a separate packet integrity row |
| Evidence bundle attachment schema | `not_assessed` | Schema or Go model defining how refs, digests, and retained forms are stored alongside the packet |
| Deferred P1 theater completeness impact | `not_assessed` | Explicit statement of what P0 theater detection cannot catch without P1 codes, visible in `PC-RESIDUAL-GAPS` |

---

## 6. Scope-Control and Overclaim Risks

### Overclaim: Theater Detection Completeness

The spec presents P0 theater codes (`agent_claimed_verification`, `unbound_intent`, `ci_theater`, `scope_theater`) as if theater detection is close to complete. But `actor_laundering`, `review_theater`, `artifact_theater`, and `human_approval_theater` are P1. In a real agent supply chain, actor laundering and review theater are high-frequency, high-impact failures. Deferring them to P1 means the P0 theater row will miss the most damaging patterns.

**Risk**: Buyers read `PC-THEATER` and trust that theater findings are comprehensive. They are not.

**Mitigation needed**: `PC-RESIDUAL-GAPS` must explicitly name the deferred P1 theater codes and their current non-coverage as a visible residual gap, not just mention deferral in prose.

### Overclaim: Packet as Self-Evidence

The packet describes evidence gaps, theater findings, and missing states. But nothing in the contract binds the packet itself to the evidence it references. A malicious or buggy packet generator could:
- Claim theater findings that are not derived from the cited evidence refs
- Mark rows `pass` when evidence does not support it
- Reference stale or non-existent evidence bundles

**Risk**: The packet itself becomes evidence theater-signed, authoritative, and wrong.

**Mitigation needed**: Either add a packet integrity row that binds the packet to source evidence digests, or explicitly state in `PC-RESIDUAL-GAPS` that packet self-reference integrity is not yet in scope.

### Underclaim: `PC-VERIFICATION` Conflation

The `PC-VERIFICATION` row conflates three distinct evidence situations:

1. Agent claims tests passed, no harness observation, no CI artifact -> `cannot_verify`
2. Agent claims tests passed, harness observed command, no CI artifact -> `partial` + `ci_theater`
3. Agent claims tests passed, harness observed command, retained CI artifact matches -> `pass`

The current example collapses (1) and (2) into `cannot_verify`, understating what the packet does observe and overstating the evidence gap.

---

## 7. Strongest Reason to Proceed

The core product contract structure is sound. The row taxonomy, evidence state vocabulary, backlog gate fields, and Russian enterprise baseline are well-designed and internally consistent. The identified issues are fixable without changing the contract's fundamental direction. The example packet is concrete enough to guide implementation and the traceability matrix correctly preserves current work. If TE-001 through TE-005 are resolved, the contract creates a real gate against substrate-only P0 work.

---

## 8. Strongest Reason Not to Proceed Yet

`PC-THEATER` row state semantics are inverted. Theater findings are expected output, not a failure condition. The current `fail` state for a packet that correctly identifies two major theater findings contradicts the packet's own findings table and signals the wrong thing to a buyer. This is not a wording issue-it is a semantic error that would cause implementers to design the theater row incorrectly.
