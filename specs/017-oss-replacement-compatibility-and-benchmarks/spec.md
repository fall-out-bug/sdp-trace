# Spec 017: OSS Replacement Compatibility And Benchmarks

Status: in_progress

## Objective

Create a reproducible compatibility and performance harness for candidate OSS
replacements: JSON Schema tooling, CUE, OPA/Rego, in-toto, Sigstore/Cosign,
SLSA verifier, and minimal shell/Go prototypes.

This spec tests substitution boundaries and records the replacement decisions
for this phase. It does not approve replacing product code by implementation
accident.

## Evidence From 2026-05-20 Probe

- `check-jsonschema` validates checked flight-recorder fixtures against
  locally rewired schema refs.
- Current live `sdp-trace wrap` output does not validate against
  `schema/flight-recorder-run.schema.json`; required fields and timestamp
  format differ.
- `check-jsonschema` validates `examples/flight-recorder/local-positive/run.json`
  against `schema/flight-recorder-run.schema.json`.
- `OPA` can express a simplified adapter-capture pass rule; the checked-in
  positive fixture evaluates to true.
- `CUE` can import JSON Schema, but direct validation is blocked until schema
  refs are packaged as a CUE module.
- `in-toto-run` can wrap a command, sign link metadata, and record material and
  product hashes (manual-only; no harness probe).
- `cosign` can sign and verify a local `run.json` blob with a local key when
  transparency-log verification is explicitly disabled (manual-only; no harness
  probe).
- `slsa-verifier` rejects the local DSSE fixture because the signature is
  truncated; Rekor rejection is not separately evidenced.

## Benchmark Snapshot

Local 20-iteration benchmark, compiled `sdp-trace`, Linux amd64.
**Note:** Min and max were not preserved from the one-shot run; this table
is provisional scope evidence only.

| Probe | Median ms | Notes |
| --- | ---: | --- |
| `sdp-trace version` | 4.5 | Built-in, measured via `tools/ossbench` |
| `sdp-trace wrap` | 16.1 | Built-in, measured via `tools/ossbench` |
| shell prototype wrap | 6.0 | One-shot; minimal JSON, no hash chain semantics |
| OPA adapter policy eval | 14.0 | One-shot; simplified policy |
| Cosign local verify | 30.5 | One-shot; transparency log ignored |
| `in-toto-run` | 148.0 | One-shot; signed link metadata |
| `check-jsonschema` fixture validation | 271.5 | One-shot; Python validator startup cost |

## Requirements

- FR-017-001: Add a reproducible OSS compatibility command or tool under the
  Go active path or a documented non-product sandbox path.
- FR-017-002: Cover at least these probes: JSON Schema, OPA/Rego, CUE import,
  in-toto command wrapping, Cosign blob signing, SLSA verifier negative path.
- FR-017-003: Record substitution boundaries: what can replace code, what needs
  adapter glue, and what remains sdp-trace-specific.
- FR-017-004: Add benchmark output with median, min, max, iterations, command,
  and environment.
- FR-017-005: Keep benchmarks non-authoritative: they inform scope decisions
  but do not prove production readiness.
- FR-017-006: Do not add Node.js, npm, JavaScript, TypeScript, or `.mjs`
  tooling to the product path.

## Decisions

- No existing product command is replaced by OSS tooling in this phase.
  OSS tools are compatibility probes and candidate substrates only.
- The in-scope OSS candidates are JSON Schema validation, OPA/Rego, CUE,
  in-toto, Cosign, and SLSA verifier. Other tools require a spec update before
  implementation work starts.
- Reproducible compatibility and benchmark tooling belongs under the Go active
  path (`tools/osscompat` and `tools/ossbench`) unless a later spec explicitly
  narrows a probe to docs-only evidence.
- JSON Schema can be used for external CI/example validation now. Live
  recorder output is not schema-compatible until the `wrap` output/schema
  drift is fixed or a current recorder schema is defined.
- OPA/Rego may replace hand-written policy expressions only behind an explicit
  JSON translation layer. It does not replace verifier evidence collection.
- CUE remains import-only until schema refs are packaged as CUE modules.
- in-toto, Cosign, and SLSA tooling may provide signing/provenance substrates.
  They do not replace `witness`, `release-proof`, or gate verdict semantics
  without OIDC, Rekor, and identity-policy evidence.
- Benchmarks are scope evidence only. They cannot create health scores,
  production readiness, or replacement approval.

## Non-Goals

- No automatic migration to OSS replacements.
- No production Sigstore, Rekor, or SLSA trust claim from local fixtures.
- No benchmark health score.
- No long-running external service dependency in default local tests.

## Acceptance Criteria

- Compatibility probes can be rerun from a clean checkout with documented
  prerequisites.
- Probe results distinguish `pass`, `fail`, `cannot_verify`, and
  `not_assessed`.
- The live `wrap` output/schema drift is either fixed or documented as a
  blocker for flight-recorder schema compatibility.
- SLSA/Rekor production trust remains `not_assessed` unless live external
  provenance is provided.
