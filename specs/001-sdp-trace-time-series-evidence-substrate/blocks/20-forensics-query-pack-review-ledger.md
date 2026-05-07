# Block 20 Review Ledger

Status: Socratic spec review completed for approval handoff.
Implementation is not approved.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S20-PB-01 | critical | product-boundary | Query-pack row states included `pass` and `fail`, creating row-level policy verdict risk despite the no native verdict boundary. | Accepted and fixed. Query-pack rows now use query evidence states; upstream `pass`/`fail` may appear only as source condition states, not row verdicts. | MiniMax-M2.7 product-boundary review; `20-forensics-query-pack.md` Result Contract; `spec.md` FR-107 |
| S20-PB-02 | critical | product-boundary | `next_required_evidence` sounded prescriptive and could become policy guidance instead of evidence observation. | Accepted and fixed. Replaced with safe `evidence_gap` descriptions limited to coarse-grained evidence families declared in selected artifacts. | MiniMax-M2.7 product-boundary review; `20-forensics-query-pack.md` Product Boundary and Result Contract; `spec.md` FR-107 |
| S20-PB-03 | major | product-boundary | `safe metadata` and `critical evidence` were undefined, leaving ambiguous scope for timeline and reconstruction rows. | Accepted and fixed. Safe metadata is now enumerated with exclusions; critical evidence is scoped to selected Block 18 policy classifications. | MiniMax-M2.7 product-boundary review; `20-forensics-query-pack.md` Product Boundary and Result Contract |
| S20-PB-04 | major | product-boundary | Explain output could become a natural-language summarizer or encode hidden state through presentation. | Accepted and fixed. Explain output is constrained to stable row rendering from JSON with no hidden severity ordering, color state, indentation state, whitespace state, omitted-section state, inferred conclusions, or new summaries. | MiniMax-M2.7 product-boundary review; `20-forensics-query-pack.md` Result Contract; `spec.md` FR-110 |
| S20-PS-01 | critical | privacy/safety | Result contract required artifact references, conflicting with the raw-reference access-note safety boundary. | Accepted and fixed. Query-pack results now record input artifact digests or path-redacted provider-neutral identifiers already present in upstream artifacts. | Kimi K2P6 privacy/safety review; `20-forensics-query-pack.md` Result Contract; `spec.md` FR-109 |
| S20-PS-02 | critical | privacy/safety | Command and test timeline aggregation could leak command names, executable paths, script paths, and test identifiers. | Accepted and fixed. Timelines now use opaque command/test identifiers only; unsafe command/test labels are forbidden unless public-catalog safe. | Kimi K2P6 privacy/safety review; `20-forensics-query-pack.md` Query Pack Model and Safety Requirements; `spec.md` FR-111 |
| S20-PS-03 | major | privacy/safety | Default stdout output creates ambient persistence in CI logs, terminal scrollback, and shell pipelines. | Accepted and fixed. `--out` is required for forensics query packs; future stdout mode must be explicit and warn to stderr. | Kimi K2P6 privacy/safety review; `20-forensics-query-pack.md` Query Pack Model; `tasks.md` T179 |
| S20-PS-04 | major | privacy/safety | Hazardous fixtures and negative leak assertions could leak realistic tokens through fixtures or failure messages. | Accepted and fixed. Fixture data must be synthetic and negative leak assertions must avoid echoing candidate sensitive strings. | Kimi K2P6 privacy/safety review; `20-forensics-query-pack.md` Acceptance Criteria; `tasks.md` T181; `spec.md` SC-048 |
| S20-TE-01 | critical | tracing/evidence | Missing Block 18/19 artifact handling depended on subjective "misleading row" language with no deterministic per-query required/optional artifact map. | Accepted and fixed. Added required/optional artifact table for each initial query and explicit absent-artifact state mapping. | OpenRouter Qwen tracing/evidence review; ZAI/GLM-5.1 tracing/evidence retry; `20-forensics-query-pack.md` Query Pack Model and Evidence Derivation Rules |
| S20-TE-02 | critical | tracing/evidence | Query-pack row evidence states had no explicit decision tree mapping upstream Block 09/18/19 facts and condition states to output row states. | Accepted and fixed. Added deterministic state propagation table and prohibited unmapped fall-through defaults. | OpenRouter Qwen tracing/evidence review; ZAI/GLM-5.1 tracing/evidence retry; `20-forensics-query-pack.md` Evidence Derivation Rules; `spec.md` FR-107 |
| S20-TE-03 | critical | tracing/evidence | `evidence_gap` used coarse-grained families but no closed vocabulary, leaving free-text evidence gaps. | Accepted and fixed. Added closed `evidence_family` vocabulary and acceptance tests rejecting free-text families. | OpenRouter Qwen tracing/evidence review; `20-forensics-query-pack.md` Result Contract and Acceptance Criteria; `spec.md` FR-107 |
| S20-TE-04 | major | tracing/evidence | Digest-only evidence was prohibited from being called reconstructable, but result rows lacked a machine-readable marker for reconstructability. | Accepted and fixed. Added `reconstructable` boolean and digest-only mapping to `retention_limited` with `reconstructable: false`. | OpenRouter Qwen tracing/evidence review; ZAI/GLM-5.1 tracing/evidence retry; `20-forensics-query-pack.md` Result Contract and Evidence Derivation Rules; `spec.md` FR-108 |
| S20-TE-05 | major | tracing/evidence | The spec did not define which files each query reads, making `input_artifacts` and missing-artifact behavior non-deterministic. | Accepted and fixed. Added required/optional upstream artifact table and required `input_artifacts` union of every upstream file read. | OpenRouter Qwen tracing/evidence review; `20-forensics-query-pack.md` Query Pack Model and Result Contract; `spec.md` FR-109 |
| S20-TE-06 | major | tracing/evidence | Query-scoped row ids were described as opaque but also used for stable sorting, creating ambiguity. | Accepted and fixed. Row ids now use deterministic `<query-short-name>.<NNNN>` format and must not derive from upstream event sequence numbers or shared counters. | OpenRouter Qwen tracing/evidence review; `20-forensics-query-pack.md` Result Contract |
| S20-TE-07 | major | tracing/evidence | Source condition refs had no canonical shape, and Block 09-derived rows do not have Block 18/19 condition ids. | Accepted and fixed. Added `source_ref` formats for Block 09 run/event/witness refs and Block 18/19 condition refs; Block 09 rows omit source condition fields when unavailable. | OpenRouter Qwen tracing/evidence review; ZAI/GLM-5.1 tracing/evidence retry; `20-forensics-query-pack.md` Result Contract |
| S20-TE-08 | major | tracing/evidence | `output_safety` did not distinguish verified absence from design-time omission. | Accepted and fixed. `output_safety` is now a verified absence assertion for serialized JSON and explain output, not a claim about upstream material. | ZAI/GLM-5.1 tracing/evidence retry; `20-forensics-query-pack.md` Result Contract; `spec.md` SC-048 |
| S20-TE-09 | major | tracing/evidence | Acceptance criteria were prose, not machine-checkable enough to enforce fixture coverage and state mappings. | Accepted and fixed. Added machine-checkable fixture matrix requirements and tests rejecting unmapped states, free-text families, bad source refs, missing input artifacts, and invalid explain ordering. | OpenRouter Qwen tracing/evidence review; `20-forensics-query-pack.md` Acceptance Criteria; `tasks.md` T180 |
| S20-TE-10 | major | tracing/evidence | Query groups such as gaps/unverified-claims risked implicit assessment weight through aggregates, ratios, ranked lists, or summary derivation. | Accepted and fixed. Query groups are now flat row arrays only; aggregates require a new pack version and safety review, and `forensics-summary` references other query rows instead of independently deriving state. | OpenRouter Qwen tracing/evidence review; ZAI/GLM-5.1 tracing/evidence retry; `20-forensics-query-pack.md` Query Pack Model and UX Requirements |
| S20-TE-11 | minor | focused tracing/evidence re-review | Focused re-review approved the spec but asked to clarify optional Block 18/19 absence scoping, required-artifact evidence-family override, and Block 09 source refs for task/command/file/test/supersession/claim facts. | Accepted and fixed. Missing optional artifacts are scoped to artifacts referenced by the selected query; missing required Block 18/19 artifacts use query-specific evidence families; Block 09 event source refs now explicitly cover those fact families through `evidence_family`. | OpenRouter Qwen focused tracing/evidence re-review; `20-forensics-query-pack.md` Result Contract and Evidence Derivation Rules |

## Implementation Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |

## Review Evidence State

- Socratic spec review: assessed with MiniMax-M2.7 product-boundary, Kimi K2P6
  privacy/safety, ZAI/GLM-5.1 tracing/evidence retry, and OpenRouter Qwen
  tracing/evidence fallback. The initial ZAI/GLM tracing attempt returned
  unusable non-review output and was not counted as evidence. DeepSeek fallback
  produced no usable output and was not counted as evidence. Critical and major
  findings from usable external reviews were accepted and fixed. OpenRouter Qwen
  focused tracing/evidence re-review returned `APPROVE`; minor clarifications
  were accepted and fixed.
- Code/correctness review: `not_assessed`; implementation not started.
- Tracing/evidence review: `not_assessed`; implementation not started.
- Requirements-vs-implementation review: `not_assessed`; implementation not
  started.
- PR-level review: `not_assessed`; no PR exists for Block 20.
- GitHub CI: `not_assessed`; no PR exists for Block 20.
