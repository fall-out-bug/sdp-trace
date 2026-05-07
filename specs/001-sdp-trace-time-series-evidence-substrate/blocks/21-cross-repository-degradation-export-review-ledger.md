# Block 21 Review Ledger

Status: Socratic spec review assessed. Critical and major findings from
product-boundary, tracing/evidence, and privacy/safety planes were accepted and
fixed. Focused re-reviews returned `APPROVE` with no remaining critical or
major findings. Implementation remains blocked pending explicit user approval
of the reviewed spec direction and Block 20 dependency disposition.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S21-PB-01 | major | product-boundary | Explain mode lacked stable traversal, truncation, and large-output behavior, making output non-reproducible. | Accepted and fixed. Explain now has stable traversal order, no truncation/pagination/hidden filters for v1, and must re-run output-safety checks before printing. | MiniMax-M2.7 product-boundary review; `21-cross-repository-degradation-export.md` CLI Boundary |
| S21-PB-02 | major | product-boundary | CLI error taxonomy was underspecified for usage, partial exports, unwritable output, and malformed inputs. | Accepted and fixed. Added developer-facing closed error taxonomy and stderr boundaries. | MiniMax-M2.7 product-boundary review; CLI Boundary |
| S21-PB-03 | major | product-boundary | Digest manifest safety contract did not cover digest/path substitution and safe logging. | Accepted and fixed. Added path-redacted artifact id grammar, digest-set handling, raw digest-manifest path safety class, and unsafe path refusal. | MiniMax-M2.7 product-boundary review; Input Model, Aggregation Rules, Safety Requirements |
| S21-PB-04 | major | product-boundary | `comparable: false` could be missed in large exports and non-comparable reasons were free text. | Accepted and fixed. Added `movement_summary`, metric row refs, dimension key, closed comparison basis, and closed non-comparable reasons. | MiniMax-M2.7 product-boundary review; Result Contract, Aggregation Rules |
| S21-PB-05 | major | product-boundary | Block 20 dependency status contradiction must be resolved before implementation claims full confidence. | Accepted as unresolved approval dependency. The Activation Gate keeps this as a blocker for implementation approval unless explicitly accepted as historical drift. | MiniMax-M2.7 product-boundary review; Activation Gate |
| S21-TE-01 | critical | tracing/evidence | Witness-scope metrics referenced fields absent from the Block 20 result contract. | Accepted and fixed. Added posture signal manifest as the source for witness-scope metrics and required absent source fields to emit not-assessed/cannot-verify metric rows instead of invented facts. | OpenRouter Qwen tracing/evidence review; Input Model, Metric Catalog |
| S21-TE-02 | critical | tracing/evidence | Movement rows lacked refs back to source metric rows. | Accepted and fixed. Movement rows now include metric id/version, dimension key, current/previous metric row refs, and closed comparison basis. | OpenRouter Qwen tracing/evidence review; Result Contract |
| S21-TE-03 | critical | tracing/evidence | Denominator definition was circular because active grouping keys were undefined. | Accepted and fixed. Added selected grouping set id, active grouping keys, first-profile grouping sets, and denominator definition against active grouping keys. | OpenRouter Qwen tracing/evidence review; Result Contract, Aggregation Rules |
| S21-TE-04 | critical | tracing/evidence | Override and contract-change metrics depended on upstream facts not in Block 20. | Accepted and fixed. Added posture signal manifest and validated external verdict input as explicit sources; absent fields remain not_assessed/cannot_verify. | OpenRouter Qwen tracing/evidence review; Input Model, Metric Catalog |
| S21-TE-05 | major | tracing/evidence | Freshness boundary was referenced but not parameterized. | Accepted and fixed. Repository selection manifest now carries `freshness_boundary`. | OpenRouter Qwen tracing/evidence review; Input Model |
| S21-TE-06 | major | tracing/evidence | Metric catalog/version relationship was unspecified. | Accepted and fixed. Metric version changes now require export profile version bump. | OpenRouter Qwen tracing/evidence review; Result Contract |
| S21-TE-07 | major | tracing/evidence | Source artifact digest refs cardinality/format was unconstrained. | Accepted and fixed. Metric rows use sorted source input refs plus digest-set hash; full per-input digests live in input selection. | OpenRouter Qwen tracing/evidence review; Aggregation Rules |
| S21-TE-08 | major | focused tracing/evidence re-review | `refusal_rows` required closed reason codes but did not enumerate them. | Accepted and fixed. Added closed `refusal_reason` enum and required schema/tests to reject other strings. | OpenRouter Qwen focused tracing/evidence re-review; Result Contract, Aggregation Rules |
| S21-TE-09 | major | focused tracing/evidence re-review | `output_safety` class identifiers were not a closed enum. | Accepted and fixed. Added closed `output_safety.sensitive_class` enum matching the aggregate export sensitive-class list. | OpenRouter Qwen focused tracing/evidence re-review; Result Contract |
| S21-TE-10 | major | focused tracing/evidence re-review | `not_assessed_count` on metric rows was required but undefined across Block 20 rows, posture signal fields, and optional inputs. | Accepted and fixed. Defined metric-row `not_assessed_count` for Block 20-derived metrics and posture/external-input-derived metrics, excluding refused stale/malformed/untrusted inputs. | OpenRouter Qwen focused tracing/evidence re-review; Metric Catalog |
| S21-PS-01 | critical | privacy/safety | Safe labels were delegated to upstream artifacts without Block 21 verification. | Accepted and fixed. Block 21 now validates safe labels itself with slug grammar, unsafe class rejection, and dimension exposure policy. | Kimi K2P6 privacy/safety review; Input Model, Acceptance Criteria |
| S21-PS-02 | critical | privacy/safety | Free-text refusal and non-comparable reasons could leak parser errors, paths, URLs, or secrets. | Accepted and fixed. Reasons are closed-enum only, and free-text refusal/exception strings are a forbidden output class. | Kimi K2P6 privacy/safety review; Result Contract, Safety Requirements |
| S21-PS-03 | major | privacy/safety | Safety test requirements were not concrete enough to prove leak detection. | Accepted and fixed. Added reserved synthetic prefixes, regex-class checks, URL/path checks, and hashed negative assertions. | Kimi K2P6 privacy/safety review; Safety Requirements |
| S21-PS-04 | major | privacy/safety | Dimension-level aggregation can reveal sensitive organizational topology beyond raw PII. | Accepted and fixed. Added `dimension_exposure_policy` as a safety boundary and refusal for dimensions outside policy. | Kimi K2P6 privacy/safety review; Input Model, Safety Requirements |
| S21-PS-05 | major | privacy/safety | External verdict inputs could carry unsafe payloads into override metrics. | Accepted and fixed. Block 21 may count only closed validated fields; malformed, unsafe, or payload-bearing external verdict inputs become `cannot_verify_input`. | Kimi K2P6 privacy/safety review; Safety Requirements |
| S21-PS-06 | major | privacy/safety | Explain output had no secondary safety barrier. | Accepted and fixed. Explain must re-run output-safety check on rendered bytes before printing. | Kimi K2P6 privacy/safety review; CLI Boundary |

## Implementation Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |

## Review Evidence State

- Socratic spec review: assessed with MiniMax-M2.7 product-boundary,
  OpenRouter Qwen tracing/evidence, and Kimi K2P6 privacy/safety planes.
  Initial reviews returned `REVISE`; accepted critical and major findings were
  fixed. Focused product-boundary, privacy/safety, and tracing/evidence
  re-reviews returned `APPROVE` with no remaining critical or major findings.
- Block 20 dependency disposition: unresolved. Block 21 implementation remains
  blocked unless the user resolves the Block 20 approval/status contradiction
  or explicitly accepts it as historical drift for this dependency.
- Code/correctness review: not_assessed; no implementation exists.
- Tracing/evidence review: not_assessed; no implementation exists.
- Requirements-vs-implementation review: not_assessed; no implementation exists.
- PR-level review: not_assessed; no PR exists.
- GitHub CI: not_assessed; no PR exists.
