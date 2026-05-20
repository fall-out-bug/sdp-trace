# OSS Policy Prototype

Status: draft
Spec: [017](../specs/017-oss-replacement-compatibility-and-benchmarks/)

This document records the OPA/Rego policy-as-code prototype for an
simplified adapter-capture profile. It is a **local experiment only**;
it does not approve replacing sdp-trace verifier behavior with OPA.

## Probe Result

| Probe | State | Reason |
|---|---|---|
| OPA adapter-capture pass/fail rule | `pass` | Policy evaluates correctly on test fixtures |
| OPA gate verdict replacement | `cannot_verify` | OPA does not understand sdp-trace gate semantics without adapter glue |
| CUE direct validation | `cannot_verify` | Schema refs not packaged as CUE modules |

## Files

- `examples/oss-policy/adapter.rego`
- `examples/oss-policy/test-fixture.json`

## Substitution Boundary

### What OPA Can Replace

- Simplified policy expressions for adapter-capture profiles (e.g.,
  trace_id presence, provenance array bounds).
- External CI policy checks when the input is a flat JSON event.

### What Needs Adapter Glue

- Gate verdict semantics (`pass`, `fail`, `cannot_verify`, `not_assessed`)
  are sdp-trace-specific. OPA returns `true`/`false` or a data structure;
  mapping this to a gate verdict requires an explicit translation layer.
- `sdp-trace-claim` tags are not native OPA concepts.
- Trace provenance (hash chains, recorder profiles, witness references)
  must be flattened into JSON before OPA can reason about them.

### What Remains sdp-trace-Specific

- Evidence collection (`assess`, `verify`, `witness`) is controlled by the
  product verifier, not by OPA evaluation.
- Source-bound proof requires immutable source commits; OPA does not
  verify git hashes or file digests.
- Release-proof and checkpoint semantics are outside OPA's scope.

## Non-Goals

- No production migration to OPA without a spec update.
- No replacement of verifier evidence collection.
- No OPA evaluation of live recorder output until the wrap output/schema
  drift is resolved.
