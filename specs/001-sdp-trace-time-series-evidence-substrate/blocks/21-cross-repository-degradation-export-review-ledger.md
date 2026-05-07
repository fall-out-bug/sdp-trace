# Block 21 Review Ledger

Status: Socratic spec review assessed. Critical and major findings from
product-boundary, tracing/evidence, and privacy/safety planes were accepted and
fixed. Focused re-reviews returned `APPROVE` with no remaining critical or
major findings. Implementation was approved by the user after disclosure of the
Block 20 dependency drift; implementation and PR-level review are in progress.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S21-PB-01 | major | product-boundary | Explain mode lacked stable traversal, truncation, and large-output behavior, making output non-reproducible. | Accepted and fixed. Explain now has stable traversal order, no truncation/pagination/hidden filters for v1, and must re-run output-safety checks before printing. | MiniMax-M2.7 product-boundary review; `21-cross-repository-degradation-export.md` CLI Boundary |
| S21-PB-02 | major | product-boundary | CLI error taxonomy was underspecified for usage, partial exports, unwritable output, and malformed inputs. | Accepted and fixed. Added developer-facing closed error taxonomy and stderr boundaries. | MiniMax-M2.7 product-boundary review; CLI Boundary |
| S21-PB-03 | major | product-boundary | Digest manifest safety contract did not cover digest/path substitution and safe logging. | Accepted and fixed. Added path-redacted artifact id grammar, digest-set handling, raw digest-manifest path safety class, and unsafe path refusal. | MiniMax-M2.7 product-boundary review; Input Model, Aggregation Rules, Safety Requirements |
| S21-PB-04 | major | product-boundary | `comparable: false` could be missed in large exports and non-comparable reasons were free text. | Accepted and fixed. Added `movement_summary`, metric row refs, dimension key, closed comparison basis, and closed non-comparable reasons. | MiniMax-M2.7 product-boundary review; Result Contract, Aggregation Rules |
| S21-PB-05 | major | product-boundary | Block 20 dependency status contradiction must be resolved before implementation claims full confidence. | Accepted as historical drift for Block 21 dependency purposes after explicit user approval. This is not a new Block 20 closure claim. | MiniMax-M2.7 product-boundary review; Activation Gate; user approval before implementation |
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
| I21-CC-01 | minor | code/correctness | Duration-style freshness boundaries such as `P7D` silently disabled freshness checks in v1. | Accepted and fixed. v1 now rejects duration freshness boundaries with a non-nil error and has regression coverage. | ZAI/GLM-5.1 code/correctness review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| I21-RS-01 | major | requirements/safety | Selection manifest artifact paths were passed to file reads without rejecting absolute paths, traversal, or URL-like references. | Accepted and fixed. Query-pack, digest-manifest, and posture-signal paths are checked before read; unsafe paths create explicit refusal rows. | Kimi K2P6 requirements/safety review and focused re-review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| I21-RS-02 | major | requirements/safety | `input_id` and `time_window` bypassed safe-label validation despite being emitted into schema-constrained and explain-rendered fields. | Accepted and fixed. Both fields now pass through safe-label validation before aggregation. | Kimi K2P6 requirements/safety review and focused re-review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| I21-RS-03 | major | requirements/safety | Explain output safety did not reject `token` or `credential` substrings even though labels rejected them. | Accepted and fixed. Explain safety now blocks those substrings while allowing the closed safety-class enum value `credential_or_token`; regression safety tests cover unsafe labels and explain output. | Kimi K2P6 requirements/safety review and focused re-review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| I21-RS-04 | minor | requirements/safety | Digest manifest path was checked for unsafe characters but not bound to the selected query-pack path. | Accepted and fixed. Digest manifest artifact path must match the selected query-pack basename before the digest can be trusted. | Kimi K2P6 requirements/safety review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| I21-RS-05 | minor | requirements/safety | Selection, posture signal, and digest manifest schema versions were accepted without explicit version checks. | Accepted and fixed. v1 now checks all three schema versions and treats unsupported versions as non-trusted inputs or unsupported selection errors. | Kimi K2P6 requirements/safety review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| I21-RS-06 | major | focused requirements/safety re-review | Unsupported digest manifest schema returned a nil error because of Go error shadowing, allowing unsupported digest manifests to be trusted. | Accepted and fixed. `verifyDigestManifest` now returns an explicit non-nil error, with regression test `TestBuildRejectsUnsupportedDigestManifestSchema`. | Kimi K2P6 focused re-review returned `APPROVE`; `internal/posture/posture.go`; `internal/posture/posture_test.go` |

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| PR21-RS-01 | major | requirements/safety | ID counters could exceed four digits and violate schema patterns for metric, movement, and refusal rows. | Accepted and fixed. Schema patterns now allow four or more digits for those deterministic ids and refs. | Kimi K2P6 PR-level requirements/safety review; `schema/cross-repo-posture-export.schema.json` |
| PR21-RS-02 | major | requirements/safety | Omitted `handoff` serialized as JSON `null`, violating the required object schema. | Accepted and fixed. Build normalizes nil handoff to `{}` and validates handoff keys/values for safe output. | Kimi K2P6 PR-level requirements/safety review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| PR21-RS-03 | major | requirements/safety | Core repo labels could bypass safe-label validation when excluded from `dimension_exposure_policy`. | Accepted and fixed. Core labels are now validated unconditionally, and active grouping keys must be allowed by the exposure policy. | Kimi K2P6 PR-level requirements/safety review; focused re-reviews returned `APPROVE`; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| PR21-RS-04 | major | requirements/safety | Block 20 dependency disposition remained unresolved in the review ledger. | Accepted and fixed. Ledger now records explicit user approval after disclosure as historical drift for Block 21 dependency purposes, without claiming Block 20 closure. | Kimi K2P6 PR-level requirements/safety review; review ledger |
| PR21-RS-05 | minor | requirements/safety | Windows absolute paths were not rejected by selection path validation. | Accepted and fixed. Selection path validation rejects drive-letter absolute paths and has regression coverage. | Kimi K2P6 PR-level requirements/safety review; `internal/posture/posture.go`; `internal/posture/posture_test.go` |
| PR21-RS-06 | minor | requirements/safety | Explain unreadable-result error used `selection_unreadable`, which named the wrong artifact. | Accepted and fixed. CLI now emits `result_unreadable`. | Kimi K2P6 PR-level requirements/safety review; `cmd/sdp-trace/main.go` |

## Review Evidence State

- Socratic spec review: assessed with MiniMax-M2.7 product-boundary,
  OpenRouter Qwen tracing/evidence, and Kimi K2P6 privacy/safety planes.
  Initial reviews returned `REVISE`; accepted critical and major findings were
  fixed. Focused product-boundary, privacy/safety, and tracing/evidence
  re-reviews returned `APPROVE` with no remaining critical or major findings.
- Block 20 dependency disposition: accepted as historical drift for Block 21
  dependency purposes by explicit user approval after disclosure. This is not a
  new Block 20 closure claim.
- Code/correctness review: assessed with ZAI/GLM-5.1. Initial review returned
  `APPROVE` with minor findings; accepted freshness-boundary issue was fixed.
- Tracing/evidence review: assessed with MiniMax-M2.7. Review returned
  `APPROVE` with no critical or major findings.
- Requirements-vs-implementation review: assessed with Kimi K2P6. Initial
  review returned `REVISE` with major safety findings. Findings were accepted
  and fixed; focused re-review found one digest-manifest schema-version bug.
  That bug was fixed, and second focused re-review returned `APPROVE`.
- PR-level review: assessed on PR #14 with ZAI/GLM-5.1 code/correctness,
  MiniMax-M2.7 tracing/evidence, and Kimi K2P6 requirements/safety planes.
  Code/correctness and tracing/evidence returned `APPROVE` with minor
  observations. Requirements/safety returned `REVISE`; accepted critical/major
  findings were fixed. Focused re-reviews returned `APPROVE`.
- GitHub CI: not_assessed until final PR head checks report.
