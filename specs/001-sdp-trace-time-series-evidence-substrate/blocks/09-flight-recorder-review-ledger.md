# Block 09 Review Ledger

Status: spec-gate and implementation review findings recorded; valid implementation findings addressed for Block 09 local fixtures
Parent: `09-flight-recorder-trust-kernel.md`

This ledger records the Block 09 pre-implementation review gate. It does not claim recorder behavior exists yet.

## Review Inputs

- Executive Socratic review synthesis in `09-flight-recorder-socratic.md`
- `pi` / `zai/glm-5.1` CIO and Head of Engineering role reviews
- `pi` / `minimax/MiniMax-M2.7` technical executive role review
- `pi` / `kimi-coding/k2p6` CFO role review
- `pi` / `openrouter/deepseek/deepseek-v4-pro` Corporate Architect role review
- `pi` / `openrouter/deepseek/deepseek-v4-flash` COO replacement review after the initial DeepSeek v4-pro COO run hung
- `pi` / `openrouter/qwen/qwen3.6-plus` CPO role review
- `pi` / `openrouter/xiaomi/mimo-v2.5-pro` CISO role review

Notes:

- Some model outputs attempted to inspect or invent local implementation files despite read-only/no-tool prompting. Code-specific claims from those outputs are excluded from this ledger.
- Findings below are product/spec findings, not implementation defects.
- Closure evidence means "the spec and implementation plan now require this behavior." It does not mean the behavior is implemented.

## Findings and Disposition

| ID | Severity | Source/Role | Finding | Disposition | Closure Evidence |
| --- | --- | --- | --- | --- | --- |
| B09-F001 | critical | technical executive / CISO / Corporate Architect | A local hash chain is not a trust boundary because the recorded actor can delete, replace, or recompute a local event log. | Accepted; design now distinguishes local development recording from witnessed recording and forbids accountability claims from local-only chains. | `09-flight-recorder-trust-kernel.md` sections `Executive Review Outcome`, `Trust Model`, `Recorder Modes`; `spec.md` FR-050; implementation plan Slice 3. |
| B09-F002 | critical | technical executive / COO / CISO | Mid-flight attachment is a provenance gap, not proof of the earlier run. | Accepted; design requires explicit `not_assessed` late-attach boundaries and forbids inference before attachment. | `09-flight-recorder-trust-kernel.md` sections `Threats Not Solved by Local Mode`, `Verifier Semantics`; `spec.md` FR-052; implementation plan Slice 4. |
| B09-F003 | critical | CISO / COO / CPO | Voluntary recording is bypassable; the product must not pretend to solve actors choosing not to run the recorder. | Accepted; design separates recorder evidence from downstream enforcement and records missing/tampered traces as evidence for external gates. | `09-flight-recorder-socratic.md` F09-C03; `09-flight-recorder-trust-kernel.md` sections `Threats Not Solved by Local Mode`, `Out of Scope`; implementation plan `Execution Rules`. |
| B09-F004 | critical | CISO / CFO | Redaction conflicts with append-only proof unless retention, authority, and verifier states are explicit. | Accepted; design defines redaction as proof model, not formatting, and requires verifier-visible redaction states. | `09-flight-recorder-trust-kernel.md` section `Redaction Model`; `spec.md` FR-055; implementation plan Slice 5. |
| B09-F005 | major | technical executive / Head of Engineering / CFO | Digest-only evidence is often insufficient for forensic reconstruction. | Accepted; design adds retention modes and forensic profile behavior. | `09-flight-recorder-trust-kernel.md` section `Evidence Capture Requirements`; `spec.md` FR-054; implementation plan Slice 5. |
| B09-F006 | major | CIO / COO / CPO | Raw JSONL is not a usable technical executive, incident, or reviewer surface. | Accepted; design requires query surfaces for run summary, provenance, gaps, commands, file mutations, tests, redactions, and witness state. | `09-flight-recorder-trust-kernel.md` section `Query Surface`; `spec.md` FR-056 and SC-033; implementation plan Slice 6. |
| B09-F007 | major | Corporate Architect / technical executive | Event schema, canonicalization, hash algorithm, and versioning were underspecified. | Accepted; design requires deterministic canonicalization, schema version, hash algorithm, event payload digest, previous hash, and event hash. | `09-flight-recorder-trust-kernel.md` section `Event Chain Requirements`; `spec.md` FR-048 and FR-049; implementation plan Slice 1. |
| B09-F008 | major | technical executive / CISO | Witness verification must distinguish local consistency from external or witnessed accountability evidence. | Accepted; design defines `flight_recorder_local`, `flight_recorder_witnessed`, and `flight_recorder_forensic` profiles. | `09-flight-recorder-trust-kernel.md` sections `Recorder Modes` and `Verifier Semantics`; implementation plan Slices 2-3. |
| B09-F009 | major | CFO / COO | Operating cost, storage, and retention are unresolved and cannot be hidden by the recorder kernel. | Accepted as future-demo/productization obligation; not blocking Block 09 kernel implementation, but demo activation must measure overhead and artifact sizes. | `09-flight-recorder-socratic.md` demo implications include Kotlin+Bazel overhead measurement; Block 09 keeps demo out of scope until kernel exists. |
| B09-F010 | major | CPO / COO | The product cannot claim "no gates" while also requiring consumers to act on missing or tampered traces. | Accepted; design says `sdp-trace` emits verifier states and evidence only; downstream gates own policy decisions. | `09-flight-recorder-trust-kernel.md` sections `Purpose`, `Product Thesis`, `Out of Scope`; `spec.md` FR-002 remains controlling. |
| B09-F011 | major | CISO / CFO | Model identity can be self-reported and must not be treated as verified unless provider/harness evidence exists. | Accepted; design requires requested and observed model identity or explicit gap state. | `09-flight-recorder-trust-kernel.md` states `model_identity_recorded`; Socratic question 6; implementation plan Slice 2 and later adapter/demo work. |
| B09-F012 | major | Head of Engineering / Corporate Architect | Demo work must not start until the recorder kernel proves tamper, witness, late-attach, supersession, and redaction behavior with fixtures. | Accepted; tasks and implementation plan gate the Feature Flag / Entitlements demo behind Block 09. | `tasks.md` Phase 9 activation gate; `09-flight-recorder-implementation-plan.md` Slice 8 exit criteria. |

## Spec-Gate Outcome

Critical and major design findings are accepted and converted into explicit Block 09 requirements.

Spec-gate closure state:

- `B09-F001` through `B09-F008`: closed for design; open as implementation obligations.
- `B09-F009`: accepted as a demo/productization obligation; not required before recorder-kernel implementation starts.
- `B09-F010`: closed for design by preserving the `sdp-trace` / downstream gate boundary.
- `B09-F011`: closed for design; implementation must keep unverifiable model identity explicit.
- `B09-F012`: closed for design by Phase 9 activation gate.

Implementation may start at T097 only if agents treat this ledger as a blocking checklist. No Block 09 implementation closure may be claimed until the implementation review ledger records verifier evidence for the relevant obligations.

## Required Implementation Evidence Before Block 09 Closure

- event mutation fixture fails verification
- event deletion fixture fails verification
- event reordering fixture fails verification
- witnessed profile fails without witness
- witnessed profile fails on witness mismatch
- late attach fixture emits explicit `not_assessed` boundary
- task rewrite fails or is represented only as supersession
- unresolved redaction fails the relevant profile
- query commands expose gaps and witness state without policy verdicts
- dirty/local-only recorder output does not support accountability or external-trust claims

## Docs Sidecar Note

`docs/flight-recorder.md` and `schema/README.md` document the intended Block 09 modes, hash model, boundaries, redaction states, and query commands. This is reviewer-facing guidance only; it does not close any implementation obligation above.

## Implementation Review

Implementation review ran through `pi` after T097-T103 implementation.

Reviewed files included:

- `schema/flight-recorder-event.schema.json`
- `schema/flight-recorder-run.schema.json`
- `schema/flight-recorder-witness.schema.json`
- `scripts/verify-flight-recorder.mjs`
- `scripts/test-flight-recorder.sh`
- `scripts/query-flight-recorder.mjs`
- `docs/flight-recorder.md`
- `schema/README.md`
- `tasks.md`

Review runs:

- `pi` / `minimax/MiniMax-M2.7` strict technical executive/Head of Engineering review: returned `CHANGES_REQUIRED`.
- `pi` / `zai/glm-5.1` strict CISO/Corporate Architect review: returned `CHANGES_REQUIRED`.
- `pi` / `kimi-coding/k2p6` strict COO/CPO/DX review: first high-thinking run hung and was replaced with a narrower low-thinking run; returned `CHANGES_REQUIRED`.

Valid findings and dispositions:

| ID | Severity | Source | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| B09-I001 | major | MiniMax M2.7 | `run.source_summary` and `run.task_summary` could claim pass-like data without being cross-bound to matching events. | Accepted and fixed. | `scripts/verify-flight-recorder.mjs` now emits `source_summary_mismatch` and `task_summary_mismatch`; `scripts/test-flight-recorder.sh` covers both. |
| B09-I002 | major | MiniMax M2.7 / GLM 5.1 | Witness verification needed explicit profile/scope boundary checks beyond schema shape. | Accepted and fixed. | `scripts/verify-flight-recorder.mjs` now emits `witness_scope_mismatch`; tests cover scope mismatch. |
| B09-I003 | major | GLM 5.1 | A witnessed or external run could be downgraded by selecting a weaker CLI profile. | Accepted and fixed. | `scripts/verify-flight-recorder.mjs` now emits `profile_downgrade`; tests cover witnessed-run downgrade to local verification. |
| B09-I004 | minor | GLM 5.1 | `run_closed` was not required to be terminal. | Accepted and fixed. | `scripts/verify-flight-recorder.mjs` now emits `run_closed_not_terminal`; tests cover non-terminal closure. |
| B09-I005 | minor | GLM 5.1 | `forensic_importance` was verifier-significant but unconstrained by schema. | Accepted and fixed. | `schema/flight-recorder-event.schema.json` now constrains `forensic_importance` to `critical`, `standard`, or `not_assessed`. |
| B09-I006 | minor | GLM 5.1 | Witness source-baseline and recorder-version mismatch paths lacked tests. | Accepted and fixed. | `scripts/test-flight-recorder.sh` covers `witness_source_baseline_mismatch` and `witness_recorder_version_mismatch`. |
| B09-I007 | minor | GLM 5.1 | Event timestamps could move backward while preserving hash structure. | Accepted and fixed. | `scripts/verify-flight-recorder.mjs` now emits `event_time_order_mismatch`; tests cover backward event time. |
| B09-I008 | major | Kimi 2.6 | Query/docs mismatch: command timeline and witness state did not expose documented reviewer answers. | Accepted and fixed. | `scripts/query-flight-recorder.mjs` now includes command event hashes, redaction state, finish/exit/retention fields, witness scope, and witness agreement booleans. |
| B09-I009 | major | Kimi 2.6 | Query surface could obscure event type or crash on malformed inputs. | Accepted and fixed. | `scripts/query-flight-recorder.mjs` now preserves authoritative event identity fields after payload spread, includes `event_type` for file mutation/state rows, adds clean JSON read errors, and defensively handles missing redaction state. |

Rejected or non-blocking findings:

- MiniMax `null event_payload` bypass: rejected as a critical finding because schema validation rejects `null` and any `schema_invalid` finding fails standalone and chain verification.
- MiniMax late-attach empty `unavailable_scope`: rejected because schema requires `minItems: 1`.
- Kimi T070 claim/tag comment: acknowledged as pre-existing Block 04 stale-closure tracking outside Block 09 implementation scope; not changed in this block.

Implementation review evidence:

```bash
jq empty schema/*.json examples/flight-recorder/**/*.json
scripts/test-flight-recorder.sh
npm run test:flight-recorder
scripts/verify.sh --profile baseline --allow-dirty --json
git diff --check
```

Current verifier state:

- Flight-recorder tests pass locally with 23 scenarios.
- Baseline structural verifier passes only as `trust_scope: local_dirty_structural_only` because the checkout is dirty.
- `npm run validate:contracts` reaches the new flight-recorder tests but then fails artifact-hash checks for modified tracked files. Per repository trust rules, those source-bound hash artifacts must not be rewritten as part of this dirty implementation cycle.
