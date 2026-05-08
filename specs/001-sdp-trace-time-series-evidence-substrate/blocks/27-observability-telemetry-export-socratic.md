# Block 27 Socratic Review: Standard Telemetry Export

Review date: 2026-05-08

Spec under review:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/27-observability-telemetry-export.md`

## Buyer Necessity Review

### Q1. Does this block preserve the core SDP identity?

Yes, if the block stays focused on provenance, evidence, and trace facts. The
core SDP identity is not dashboards. It is security, observability, and
verifiability through inspectable delivery facts.

The spec preserves the boundary: `sdp-trace` emits metric values, movements,
refusals, and input trust states. Downstream security, observability,
governance, reporting, or policy systems decide how to use them.

Disposition: acceptable.

### Q2. Is this materially different from "just keep CI artifacts"?

Yes. CI artifacts are storage. This block turns selected `sdp-trace` posture
facts into a standard metrics surface:

- stable metric names;
- safe labels;
- deterministic output;
- movement deltas;
- refusal/input-trust visibility;
- direct compatibility with standard telemetry consumers.

The value is reusable fact export, not a new dashboard.

Disposition: acceptable.

### Q3. Is this too close to marketing/sales material?

The spec must avoid marketing docs, buyer report language, and sales framing in
repo artifacts. It should name the technical interoperability question, not
encode a sales narrative.

The examples should be technical downstream fixtures, not pitch decks:

- Prometheus text output;
- Prometheus config;

Disposition: acceptable with wording discipline.

## Technical Review

### Q4. Is Prometheus text the right MVP target?

Prometheus text is the right first format because it is dependency-light,
Go-friendly, file-compatible, and enough for Grafana through Prometheus or node
exporter textfile collector.

OpenTelemetry should not be in MVP unless implemented as a documented future
profile. Adding OTel now risks pulling dependency and semantic-convention work
into a block that should be small.

Disposition: Prometheus text first; OpenTelemetry explicitly deferred.

### Q5. Should metric rows become generic metrics or first-class metric names?

The draft proposes generic row metrics like
`sdp_trace_posture_metric_numerator{metric_id="cannot_verify_rows"}`.

This is safer for schema stability and avoids adding a new Prometheus metric
family every time the posture catalog changes. The trade-off is less ergonomic
Grafana queries.

For MVP, generic row metrics are acceptable. If Grafana examples become too
awkward, a later profile can add first-class aliases without breaking the generic
contract.

Disposition: keep generic metrics for MVP.

### Q6. Should exporter support stdout?

The draft requires `--out`. For observability pipelines, stdout is useful:

```text
sdp-trace export telemetry ... --out -
```

This is still Go-first and avoids file-only friction. It should be added to the
spec before implementation.

Disposition: major finding, spec should add `--out -` support.

### Q7. How should unsafe labels behave?

Failing the whole export is safer for trust because partial omission can hide the
fact that telemetry was incomplete. However, observability consumers also need a
visible refusal signal.

Preferred MVP behavior:

- unsafe posture input labels cause `cannot_verify` and no metrics output when
  unsafe data could leak;
- safe refusal metrics can only be emitted when the unsafe value is not echoed
  and the source can be represented by a generated safe id.

The initial implementation may fail closed. Do not silently drop bad series.

Disposition: major finding, spec should require fail-closed for unsafe labels.

### Q8. Do Grafana and Alertmanager examples violate the boundary?

They are not required for the product block. Keeping them out of Block 27 is
cleaner: Prometheus text output plus a minimal Prometheus consumption fixture is
enough to prove standard telemetry export.

Grafana dashboards, Alertmanager rules, Telegram routing, and CTO-facing reports
belong to separate demonstration streams that consume Block 27 output.

Disposition: scope narrowed.

### Q11. Do checked-in telemetry examples risk becoming false proof?

Yes, unless they are generated from live fixture output and verified for drift.
`sdp-trace` already treats checked-in proof JSON as non-authoritative unless
live-verified or externally signed. The same rule applies here: a committed
`metrics.prom` snapshot is documentation and regression fixture, not authority.

Disposition: accepted. The spec now requires regenerated output snapshots and a
verification path that fails on example drift.

## Requirements-vs-Implementation Risk Review

### Q9. Can current Block 21 data support the telemetry export without new source collection?

Yes. `internal/posture` already emits metric rows, movement rows, refusal rows,
and input selections. Block 27 can be a converter over existing output.

Risk: Block 21 metrics may not cover all future CTO questions, but that is not a
Block 27 blocker. This block exposes existing facts; it does not invent new
measurements.

Disposition: acceptable.

### Q10. Does this answer "Prometheus for agents" exactly?

No, and it should not claim that. It answers "standard telemetry export for
provenance/evidence/trace posture." It does not observe every agent action
directly, and it should not claim that.

The product cannot see unwrapped work; it can expose missing telemetry as a
metric.

Disposition: wording corrected.

## Review Findings

### Confirmed: stdout support is specified

`--out -` writes Prometheus text to stdout. This makes the exporter easier to
pipe into file collectors and avoids unnecessary temp files.

Disposition: fixed in spec.

### Confirmed: fail-closed unsafe label behavior is specified

Unsafe label values must not produce partial silently omitted telemetry. The
first implementation fails closed unless a later profile can emit a safe refusal
metric without echoing unsafe data.

Disposition: fixed in spec.

### Critical: Separate unsafe value leaks from `cannot_verify` as a fact

The first external trust/evidence review found that fail-closed unsafe-label
behavior could be read as contradicting the requirement that `cannot_verify`,
`not_assessed`, untrusted input, and refusal facts remain visible as metrics.

Disposition: accepted and fixed. The spec now distinguishes safe posture states
that must emit telemetry from unsafe label values that would leak sensitive data
and therefore fail closed with no partial metrics output.

### Major: Checked examples need drift verification

The first external trust/evidence review found that checked-in telemetry examples
could become false proof if not regenerated from live fixture output.

Disposition: accepted and fixed. The spec now requires live fixture regeneration
and drift detection for committed telemetry snapshots.

### Minor: Prometheus `_total` suffix was wrong for gauges

The product-boundary review found that `sdp_trace_posture_refusal_total` and
`sdp_trace_posture_input_total` used a counter suffix while the values are
snapshot gauges.

Disposition: accepted and fixed. The metric names are now
`sdp_trace_posture_refusal` and `sdp_trace_posture_input`.

### Major: Atomic file writes are required for textfile collectors

The telemetry review found that file output could be scraped while partially
written.

Disposition: accepted and fixed. The spec now requires writing to a temporary
sibling file and atomically renaming it to the target path.

### Major: Prometheus family metadata and bounded cardinality are underspecified

The telemetry review found missing `# HELP`/`# TYPE` lines and unbounded label
cardinality/value length.

Disposition: accepted and fixed. The spec now requires family metadata, `# EOF`,
deterministic ordering, a 10,000-series cap, 1,024-byte label value cap, and
fail-closed behavior on limit breach.

### Minor: Avoid broad "Prometheus for agents" claim in repo docs

Use "standard telemetry export for provenance/evidence/trace posture" or
"posture telemetry export." This prevents overclaim that `sdp-trace` observes
all agent activity or owns dashboards.

Resolution required before implementation docs.

## Approval Recommendation

Approve implementation only after the spec is amended for:

1. `--out -` stdout support.
2. fail-closed unsafe label behavior.
3. scoped wording around "AI delivery evidence" rather than broad "agents."
