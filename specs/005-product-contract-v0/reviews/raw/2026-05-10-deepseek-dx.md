# Product Contract v0 Review: Implementation Readiness & DX

**Reviewer role**: Implementation readiness and DX reviewer. Scope: can a developer
take this contract and build schema, fixtures, renderer, tests, and intake checks
without architecture drift or implementation ambiguity?

---

## 1. Verdict

`REVISE_BEFORE_USER_REVIEW`

The contract defines the right direction and a strong backlog gate, but the
packet format, derivation rules, evidence bundle, and minimum-data contract are
too ambiguous for a team to start implementing without making decisions that
should already be settled in the product contract.

---

## 2. What works

- The **backlog gate** is real: `packet_rows` as a mandatory field, plus
  `evidence_surface`, `closure_state`, `buyer_effect`, `non_goal` - this is a
  checkable, enforceable template (`spec.md:148-161`, `traceability.md:40-55`).
- The **required rows** and buyer questions (`spec.md:66-80`) form a clean
  semantic foundation. They are stable and provider-neutral.
- The **example packet** (`example.md`) gives a tangible artifact, demonstrates
  `not_assessed`/`cannot_verify` usage, and clearly states it is not product
  proof.
- The **traceability matrix** (`traceability.md`) honestly maps existing
  substrate to rows, preserving current investment while showing gaps.
- Russian enterprise **local-only constraint** is explicit and backed by an
  input allowlist (`spec.md:127-147`).
- Constitution check and **non-goals** are well scoped; no code is authorized
  yet.

---

## 3. Blocking findings

| id | severity | cited file:line | finding | why it matters | exact fix |
|----|----------|-----------------|---------|----------------|-----------|
| IR-001 | **major** | `spec.md:27-37`, `example.md:23-32` | The Markdown report's structure (sections, tables, mandatory fields, optional fields) is not normatively specified. The example suggests a shape, but nothing states that it is the canonical template. | Without a fixed report template, different implementers (or future projections) will produce incompatible artifacts. Diffs between packet versions become meaningless; the CTO surface degrades. | Add a mandatory "Packet Markdown Template" section to `spec.md` (or a separate `packet-format.md`) that lists every required section, table schema, allowed cell values, and rules for missing data. Alternatively, declare `example.md` as the normative template and add a test that every generated packet conforms to it. |
| IR-002 | **major** | `spec.md:35`, `example.md:67-75` | The **evidence bundle** is mentioned as "retained refs, digests, and redaction status" but its format is undefined. Is it a directory of files, a zip, a JSON manifest? How are refs resolved? | The bundle is a core deliverable; without a defined structure, packet generation cannot package evidence consistently, and downstream tools (e.g., signature validation) cannot locate items. | Define the evidence bundle as a directory with a required `manifest.json` (or `bundle-manifest.md`). Specify the manifest schema: `ref`, `digest`, `content_type`, `redacted_flag`, `source_path`. Link bundle entries to packet-row evidence refs. |
| IR-003 | **major** | `spec.md:107-124`, `tasks.md:48` | **Theater reason-code derivation rules** are described only as high-level conditions ("Agent claims tests ... but no independent retained evidence exists"). No mapping from input sources (trace events, harness observations, adapter data) to theater findings is given. | An implementer will invent ad-hoc heuristics, leading to inconsistent findings across releases. The whole "theater" value proposition rests on deterministic, evidence-backed detection. | Add a "Theater Derivation Contract" section that lists for each P0 reason code: input schemas to inspect, decision logic pseudo-code, and what qualifies as "independent retained evidence". Deferred P1 codes can reference a future derivation plan. |
| IR-004 | **major** | `spec.md:79`, `example.md:48` | **PC-RESIDUAL-GAPS** synthesis is undefined. The example row states `pass` because "Packet records missing intent, review, CI witness, authority, and signed evidence", but there is no algorithm for how that list is compiled. | Without a rule, the packet may omit critical gaps or fabricate a "pass" when real gaps are hidden. The CTO expects an honest gap list. | Specify that the row auto-populates from all rows where `state` is not `pass` and from any active theater findings, listing each gap with its source row, current state, and required closure evidence. The row state becomes `pass` only when all such gaps are enumerated; otherwise `cannot_verify`. |
| IR-005 | **major** | `spec.md:79`, `example.md:59-65` | **PC-DECISION** row definition talks about "the next human decision" (singular), but the example and the required output text list multiple decision types (merge, release, risk acceptance, security review). Are all required or only the nearest one? | If the contract forces only a single decision, multi-decision orchestration tools lose context. If all must be present, a pilot with no release process still needs a meaningless row. | Clarify that the row must cover **each** of the four decision types (merge, release, risk acceptance, security review), each with its own owner ref/state. If a decision is not applicable, it may be `not_assessed` with a reason. |
| IR-006 | **major** | `example.md:25-32`, `spec.md` (absent) | **Packet metadata fields** (`packet_id`, `schema`, `generated_from`, `selected_profile`, `redaction_policy`, `packet_state`) appear in the example but are not required by the spec. Are they mandatory, and what are their exact schemas? | Metadata is essential for packet identity, versioning, and downstream trust verification. Missing fields break tooling and audit trails. | Add a "Packet Metadata" requirement block to `spec.md` listing every mandatory metadata field, its type/allowed values, and how it is populated (e.g., `packet_id` from change identity + timestamp, `selected_profile` from the requested profile). |
| IR-007 | **major** | `spec.md:127-147`, `example.md` (entire) | The **Russian baseline profile** says "the packet still has value when it makes missing evidence explicit", but never defines the **minimum data** required to produce a valid packet. Could all rows be `not_assessed`? Does the generator refuse to run? | Without a minimum data contract, the generator's behaviour at the low-input bound is unpredictable. An enterprise pilot with only local Git may get an error instead of a useful gap report. | Define the **Minimum Viable Packet**: must include at least `PC-CHANGE` identity (even if partial) and `PC-DECISION` owner (even if `cannot_verify`). All other rows may be `not_assessed`. The contract should also state that a packet can always be rendered, but the generator must never invent evidence. |

---

## 4. Non-blocking concerns

- The **open decisions table** (`spec.md:258-267`) lists proposed defaults but no
  timeline for when they become final. A developer picking the first artifact
  (Markdown+evidence bundle) might need to know whether the decision is stable or
  subject to change after user review.
- **Attachments section** (`example.md:67-75`) is not anchored in the spec's
  core output definition; it blurs the line between the packet and the evidence
  bundle. Clarifying this would remove duplication.
- **CI default commands** (`AGENTS.md:96-98`) mention `go test ./...`, but the
  contract slice has no Go code; this is not a problem now but may confuse
  newcomers who expect Go in every slice.
- The **source strength classification** (`spec.md:93-104`) is a list of labels
  without a schema for how they appear in packet rows. The example uses a mixed
  notation (evidence refs like `git:abc123..def456`, `harness-run:...`). A
  structured approach would improve DX.

---

## 5. Missing evidence or `not_assessed` areas

- The **focused Socratic review** of this contract has not been run yet
  (`tasks.md:T009-T013`). That review is the intended gate before approval, so
  the contract's own review state is `not_assessed`. This is by design but worth
  noting.
- **No packet schema** (JSON or Go model) exists yet - acknowledged as Phase 3
  task (`tasks.md:T014`). The contract must be clear enough to serve as the spec
  for that schema; the findings above address that.
- **No CI integration specification** for the packet itself: what constitutes a
  "generated packet" artifact, how it is archived, and what integrity checks it
  carries. This may be out of scope for v0, but it should be explicitly deferred.

---

## 6. Scope-control and overclaim risks

- The **traceability matrix** (`traceability.md:9-21`) marks several rows as
  "Good substrate" (Authority, Attestation, Residual Gaps). The risk is that an
  implementer assumes these rows are nearly done and doesn't invest in the
  packet projection layer. The matrix should warn that "Good substrate" does not
  mean the packet row is implemented; it only means the raw data exists.
- The contract repeatedly says **"signed attestation is additive"** (`spec.md:233`,
  `plan.md:141`). This is a strong design decision, but the spec also leaves
  open whether some buyers might require it from day one (`spec.md:276`). The
  contract should explicitly state that v0 ships without signed attestation and
  that such requests will be triaged as separate P2 profiles, not as blockers for
  the baseline contract.
- **"CTO reads one packet in under two minutes"** (`spec.md:166-167`) is a
  success criteria that depends heavily on the template design. If the template
  is not nailed down (IR-001), this criterion is unverifiable.

---

## 7. One strongest reason to proceed

The contract finally forces every future P0 item to **prove its buyer impact** by
mapping to a stable, named packet row. This hard gate, if enforced, will end the
cycle of treating internal plumbing as shippable product - the single most
important structural fix for `sdp-trace` roadmapping.

---

## 8. One strongest reason not to proceed yet

The **packet's concrete shape is not locked down**. Without a normative template
and derivation contracts, the first implementation slice will inevitably invent
its own format and logic. That work will then need to be unwound when the real
product requirements surface, wasting the very delivery credibility this contract
is meant to protect. Nail the artifact first; build second.
