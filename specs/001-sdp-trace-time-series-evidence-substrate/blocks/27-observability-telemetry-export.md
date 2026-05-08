# Block 27: Standard Telemetry Export

Status: implementation in progress after reviewed spec direction.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/21-cross-repository-degradation-export.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/26-ci-artifact-observation-contract.md`
- `docs/process-metric-catalog.md`
- `schema/cross-repo-posture-export.schema.json`

## Goal

Expose `sdp-trace` provenance, evidence, and trace posture facts in standard
telemetry formats so downstream security, observability, governance, reporting,
or policy systems can consume them without making `sdp-trace` a gate, dashboard,
alerting system, or policy engine.

The product behavior is telemetry export:

```text
sdp-trace posture JSON -> standard telemetry output -> external consumer
```

Dashboards, alerts, and reports are demonstration or downstream consumer layers.
They are not part of the `sdp-trace` product boundary.

## Technical Question

"Can provenance/evidence/trace facts produced by `sdp-trace` be consumed by
standard operational tooling without bespoke JSON integration?"

Block 27 answers only the telemetry export layer:

- current and previous posture values are observable as metrics;
- movement values are observable as metrics;
- missing, stale, untrusted, or non-comparable inputs remain observable;
- external systems can graph, alert, report, or decide policy using their own
  thresholds and context.

`sdp-trace` must not answer "good", "bad", "degraded", "acceptable", "block",
or "page the CTO". It emits measured facts.

## Problem

Block 21 already exports deterministic cross-repository posture JSON with
`metric_rows`, `movement_rows`, and refusal facts. That is useful for a custom
consumer, but it is not immediately attachable to common observability tooling.

The missing MVP surface is small:

- a stable metric-name mapping from posture rows to standard telemetry names;
- a Prometheus-compatible text export path;
- optional OpenTelemetry-compatible JSON or line protocol guidance, if it can be
  implemented without adding a runtime dependency;
- minimal downstream fixtures showing that a standard consumer can ingest the
  metrics.

Without this, `sdp-trace` remains a bespoke artifact producer. With it, it can
act as a standard telemetry source for provenance/evidence/trace posture.

## Non-Goals

- No native degradation verdict, health score, risk score, grade, color state, or
  severity ranking.
- No native alert threshold, SLO, paging, Telegram delivery, Slack delivery,
  dashboard, report, or incident creation.
- No bundled Grafana, Prometheus, Alertmanager, Datadog, or OpenTelemetry server.
- No network push in the first implementation.
- No dependency on Node.js, npm, JavaScript, TypeScript, or `.mjs`.
- No raw command output, prompt text, model response body, source snippets,
  private paths, authenticated URLs, token-like values, or personal identifiers
  in exported labels.
- No reinterpretation of `not_assessed`, `cannot_verify`, or `fail` as policy
  decisions.

## Input Model

The first implementation should consume an existing Block 21 export:

```text
sdp-trace export telemetry \
  --profile prometheus-text-v1 \
  --cross-repo-posture <cross-repo-posture-export.json> \
  --out <metrics.prom>
```

`--out -` writes Prometheus text to stdout for shell pipelines and textfile
collector handoff.

The input must be a `cross-repo-evidence-posture-v1` result with schema version
`block21-cross-repo-posture-export-v1`; see
`schema/cross-repo-posture-export.schema.json`.

Unsupported schema versions return `cannot_verify`. Malformed input returns
`cannot_verify`.

Normal posture facts such as `not_assessed_input`, `cannot_verify_input`,
`untrusted_input`, and `refusal_rows` must remain visible as metrics when their
labels are safe. Unsafe label values are a different case: if a value that would
be serialized into the telemetry payload contains a credential, token, private
path, authenticated URL, raw personal identifier, or another unsafe value, the
export fails closed and emits no partial metrics. A later profile may emit safe
refusal metrics for unsafe inputs only when the unsafe value is not echoed and
the source can be represented by a generated safe id.

## Output Model

The MVP output is Prometheus text exposition format.

Metric names must use a stable `sdp_trace_` prefix and encode posture facts, not
policy:

| Source | Metric name | Type | Value |
| --- | --- | --- | --- |
| `metric_rows[].numerator` | `sdp_trace_posture_metric_numerator` | gauge | numerator |
| `metric_rows[].denominator` | `sdp_trace_posture_metric_denominator` | gauge | denominator |
| `metric_rows[].not_assessed_count` | `sdp_trace_posture_metric_not_assessed` | gauge | not assessed count |
| `movement_rows[].current_value` | `sdp_trace_posture_movement_current` | gauge | current value |
| `movement_rows[].previous_value` | `sdp_trace_posture_movement_previous` | gauge | previous value |
| `movement_rows[].delta` | `sdp_trace_posture_movement_delta` | gauge | signed delta |
| `movement_rows[].comparable` | `sdp_trace_posture_movement_comparable` | gauge | 1 when comparable, 0 when not comparable |
| `refusal_rows` | `sdp_trace_posture_refusal` | gauge | count by reason/state/window |
| `input_selection` | `sdp_trace_posture_input` | gauge | count by input trust state |

Labels must be closed and safe:

- `metric_id`
- `metric_version`
- `dimension_key`
- `time_window`
- `repo`
- `team`
- `service`
- `harness`
- `change_type`
- `input_trust_state`
- `refusal_reason`
- `comparable`
- `non_comparable_reason`

Labels absent in the source row must be omitted, not synthesized.

The exporter must escape Prometheus labels correctly and must reject unsafe label
values rather than writing them. Refusal labels must use the closed
`refusal_reason` values already present in the posture export. Free-text errors,
timestamps, raw exception messages, provider URLs, or other high-cardinality
values must not become label values.

Metrics must be emitted in stable order: sort by metric name lexically, then by
label set rendered with label names in lexical order and label values escaped.
This makes checked examples and downstream textfile collector output
deterministic.

Every metric family must include `# HELP` and `# TYPE` lines before its data
points. Output must end with `# EOF` followed by a newline.

The first implementation must enforce bounded telemetry output:

- no more than 10,000 series per export;
- no label value longer than 1,024 UTF-8 bytes;
- `time_window` remains a bounded categorical label from the posture export, not
  an unbounded raw timestamp.

If these limits are exceeded, the exporter exits `cannot_verify` with safe reason
text and emits no partial metrics.

## Downstream Examples

Block 27 should include examples under `examples/block27-observability-telemetry/`
that are explicitly downstream fixtures, not product dependencies:

- `prometheus/metrics.prom`: expected text output for the existing Block 21
  fixture.
- `prometheus/prometheus.yml`: minimal scrape or file-based example suitable for
  local Prometheus.
- `README.md`: how to regenerate the metrics and where responsibility moves
  from `sdp-trace` to downstream consumers.

Grafana dashboards, Alertmanager rules, Telegram routing, and CTO-facing reports
belong to separate demonstration streams. They may consume Block 27 metrics, but
they should not be required acceptance artifacts for this block.

## CLI Contract

Preferred command:

```text
sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture <file> --out <file>
```

`--out -` writes to stdout.

When writing to a file path, the exporter must write to a temporary sibling file
and atomically rename it to the target path so textfile collectors cannot scrape
torn output.

Optional explain command:

```text
sdp-trace export telemetry explain --result <metrics.prom>
```

The explain command is optional for MVP. If implemented, it must not parse policy
or thresholds; it should only summarize metric families and label sets.

## State And Exit Behavior

- Successful export exits `0`.
- Usage errors exit with the existing usage exit code.
- Unsupported or unreadable input exits `cannot_verify`.
- Unsafe label values or unsafe output exit `cannot_verify` with safe reason text
  only and no partial metrics output in the first implementation.
- Valid input with zero metric, movement, refusal, and input rows exits `0` and
  emits only the safe Prometheus comment
  `# sdp_trace_posture no_rows`. Invalid or malformed input with missing required
  fields exits `cannot_verify`.
- `sdp_trace_posture_movement_comparable=0` is a factual traceability signal
  that movement cannot be compared under the selected profile. It is not a
  degradation, readiness, health, or risk verdict.

## Acceptance Criteria

1. A Block 21 posture export can be converted into Prometheus text format by a
   Go-first CLI path.
2. Metric names and label names are deterministic, documented, and stable for
   `prometheus-text-v1`.
3. Unsafe label values are rejected or redacted before output; raw sensitive data
   is not emitted.
4. `not_assessed`, `cannot_verify`, stale input, untrusted input, and refusal
   facts remain visible as metrics.
5. Checked examples are regenerated from live fixture output and verification
   fails if regenerated output differs from the committed snapshot.
6. The examples show standard telemetry consumption without making `sdp-trace`
   own dashboards, thresholds, reports, or alerts.
7. Tests cover at least:
   - valid Prometheus output from the Block 21 fixture;
   - unsupported input schema;
   - unsafe label refusal;
   - deterministic output ordering;
   - movement delta and comparable flag rendering;
   - refusal/input trust state rendering.

## Implementation Plan Sketch

1. Add a small Go package, likely `internal/telemetry`, that converts
   `posture.ExportResult` to Prometheus text.
2. Add CLI routing under `export telemetry`.
3. Add tests using existing Block 21 fixture output.
4. Add examples under `examples/block27-observability-telemetry/`; generate
   committed output snapshots from the live fixture path and add a verification
   path that detects drift.
5. Update command docs and schema docs only where the active command surface
   changes.
6. Run the standard verification:
   - `go test ./...`
   - `jq empty schema/*.json examples/block27-observability-telemetry/**/*.json`
     where shell support permits, otherwise enumerate files explicitly;
   - `git diff --check`.

## Open Questions For Socratic Review

1. Is Prometheus text enough for MVP, or should OpenTelemetry be specified now?
2. **Resolved:** unsafe values that would leak into telemetry fail closed with no
   partial metrics output. Safe `not_assessed`, `cannot_verify`, `untrusted`, and
   refusal facts still emit metrics because those states are facts, not unsafe
   values.
3. Should the exporter produce a single snapshot file for node exporter textfile
   collector, or should it also support stdout for pipelines?
4. Are row-oriented generic metric names acceptable, or should every posture
   metric become a first-class Prometheus metric such as
   `sdp_trace_cannot_verify_rows`?
5. What is the smallest downstream fixture that proves standard consumption
   without smuggling dashboard/report/gate scope into the product?
