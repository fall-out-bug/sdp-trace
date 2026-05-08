# Block 27 Review Ledger: Standard Telemetry Export

Status: PR merge gate review complete; post-merge verification remains pending until PR #22 is merged.

## Spec Review Inputs

- Product-boundary review: `zai/glm-5.1`, local raw output under `.codex-review/`
  scratch, not committed.
- Telemetry/observability review: `kimi-coding/kimi-k2-thinking`, local raw
  output under `.codex-review/` scratch, not committed.
- Trust/evidence review: `minimax/MiniMax-M2.7`, local raw output under
  `.codex-review/` scratch, not committed.

`.codex-review/` is local scratch only. This ledger records accepted/rejected
findings for the repository.

## Findings And Disposition

| ID | Severity | Plane | Finding | Disposition |
| --- | --- | --- | --- | --- |
| S27-PB-01 | minor | product boundary | Gauge metrics used `_total`, which implies counters. | Accepted and fixed. Gauge names use `sdp_trace_posture_refusal` and `sdp_trace_posture_input`. |
| S27-PB-02 | minor | product boundary | Open questions still listed items already resolved by spec text. | Accepted. Unsafe-label behavior was marked resolved; Prometheus-first remains an explicit reviewed choice. |
| S27-PB-03 | minor | product boundary | Socratic findings still read as pending after amendments. | Accepted and fixed in the Socratic record. |
| S27-TR-01 | critical | trust/evidence | Fail-closed unsafe label behavior could be read as contradicting visible `cannot_verify`/`not_assessed` metrics. | Accepted and fixed. The spec separates unsafe value leakage from safe posture facts. |
| S27-TR-02 | major | trust/evidence | Checked-in telemetry examples could become false proof without drift detection. | Accepted and fixed. CLI tests compare live regenerated output with committed `metrics.prom`. |
| S27-TR-03 | major | trust/evidence | No-rows behavior was ambiguous. | Accepted and fixed. Valid empty output emits `# sdp_trace_posture no_rows` and exits 0; malformed input exits `cannot_verify`. |
| S27-TR-04 | major | trust/evidence | Refusal labels could become high-cardinality free text. | Accepted and fixed. Export uses closed `refusal_reason` values from posture output only. |
| S27-TR-05 | minor | trust/evidence | Input schema reference was not explicit. | Accepted and fixed. Spec references `schema/cross-repo-posture-export.schema.json`. |
| S27-TR-06 | minor | trust/evidence | Deterministic output ordering lacked sort details. | Accepted and fixed. Spec and implementation sort by metric family and label set. |
| S27-TR-07 | minor | trust/evidence | `movement_comparable=0` could be mistaken for a degradation verdict. | Accepted and fixed. Spec states it is a factual traceability signal only. |
| S27-TEL-01 | major | telemetry | File output could be scraped while partially written. | Accepted and fixed. CLI writes to a temporary sibling file and atomically renames it. |
| S27-TEL-02 | major | telemetry | Prometheus `# HELP` and `# TYPE` metadata were not required. | Accepted and fixed. Renderer emits both for every family. |
| S27-TEL-03 | major | telemetry | Label value length and series cardinality were unbounded. | Accepted and fixed. Renderer enforces 1,024-byte label values and 10,000 series maximum. |
| S27-TEL-04 | minor | telemetry | Output ordering sort key was not explicit. | Accepted and fixed. Renderer sorts by family name and rendered label set. |
| S27-TEL-05 | minor | telemetry | OpenMetrics EOF terminator was absent. | Accepted and fixed. Renderer emits `# EOF`. |
| S27-TEL-06 | major | telemetry implementation | Duplicate metric name plus label set would produce invalid Prometheus series. | Accepted and fixed. Renderer rejects duplicate series before serialization. |
| S27-TEL-07 | major | telemetry implementation | Reviewer suggested allowing `/` in labels for `org/repo`. | Rejected for current profile. Block 21 safe-label contract forbids `/`, and preserving private-path fail-closed behavior is more important for v1. |
| S27-TEL-08 | minor | telemetry implementation | Carriage returns in labels were not escaped. | Accepted and fixed. Renderer escapes `\r`. |
| S27-DEMO-01 | critical | drift discovery | Current customer demo repo was observed as Java/Bazel despite earlier Kotlin/Bazel framing. | Accepted as separate demo reset blocker, not a Block 27 PR blocker. Block 27 changes only the portable telemetry exporter and does not claim demo readiness. Follow-up: reset `fall-out-bug/sdp-trace-demo-jvm-gsd` immediately after PR #22 merge before making any customer demo claim. |
| PR27-CODE-01 | major | PR-level code/correctness | Reviewer reported possible `dimension_key` contamination through `metric_rows[].dimensions`. | Accepted defensively. The exporter already copied only closed dimensions, and a regression test now rejects malformed `dimension_key` inside the dimensions map. |
| PR27-CODE-02 | major | PR-level code/correctness | Reviewer claimed `os.CreateTemp` creates a world-readable race before chmod. | Rejected as false positive. Go creates temp files with mode `0600`; the CLI writes and closes the sibling temp file before chmod and atomic rename to the destination. |
| PR27-REQ-01 | major | PR-level requirements-vs-implementation | A minimal JSON object with only version/profile fields was treated as a valid empty posture export. | Accepted and fixed. `posture.ValidateExportResult` now enforces required structural fields before telemetry rendering; malformed input exits `cannot_verify` with no partial output. |
| PR27-REQ-02 | minor | PR-level requirements-vs-implementation | Command docs obscured stdout support and the unsafe-label heuristic missed some credential-shaped labels. | Accepted and fixed. English/Russian docs use `--out <file\|->`, and telemetry unsafe-value checks include `bearer`, `api_key`, `access_key`, and `private`. |
| PR27-TRACE-01 | major | PR-level tracing/evidence | Ledger still said implementation review and full verification were pending. | Accepted and fixed. This ledger now records PR-level code, tracing/evidence, and requirements review dispositions; post-merge verification remains explicitly pending until merge. |

## Verification Notes

Focused implementation verification:

- `go test ./internal/telemetry ./cmd/sdp-trace`: passed.

PR-level review planes:

- Code/correctness: revised for PR27-CODE-01; false positive recorded for
  PR27-CODE-02; independent code audit returned `APPROVE`.
- Tracing/evidence: revised for ledger status and demo follow-up rationale;
  initial PR-level tracing/evidence review returned `APPROVE`, then
  Kierkegaard re-review returned `REVISE` for overbroad fixture wording.
  Fixed: telemetry `metrics.prom` drift is checked by live export plus byte
  comparison; the Block 21 posture fixture was regenerated structurally for
  `refusal_rows: []`, not claimed as byte-stable because `generated_at` changes.
- Requirements-vs-implementation: revised for malformed posture input
  fail-closed behavior; replacement Qwen review returned `APPROVE` with minor
  findings fixed.

Fresh verification before final PR update:

- `go test ./internal/posture ./internal/telemetry ./cmd/sdp-trace`: passed.
- Minimal malformed posture JSON exits `cannot_verify` and writes no telemetry.
- Regenerated Block 21 posture fixture now serializes `refusal_rows` as `[]`
  instead of `null`; this is structural fixture repair, not a byte-stable
  posture regeneration claim.
