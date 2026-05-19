# Spec 017: OSS Replacement Compatibility And Benchmarks

Status: draft

## Objective

Create a reproducible compatibility and performance harness for candidate OSS
replacements: JSON Schema tooling, CUE, OPA/Rego, in-toto, Sigstore/Cosign,
SLSA verifier, and minimal shell/Go prototypes.

This spec tests substitution boundaries. It does not decide to replace product
code by itself.

## Evidence From 2026-05-20 Probe

- `check-jsonschema` validates checked flight-recorder fixtures against
  locally rewired schema refs.
- Current live `sdp-trace wrap` output does not validate against
  `schema/flight-recorder-run.schema.json`; required fields and timestamp
  format differ.
- `check-jsonschema` validates representative assessment, gate, and release
  example files.
- `OPA` can express a simplified adapter-capture pass/fail rule and detects the
  `test_provenance_not_overclaimed` failure fixture.
- `CUE` can import JSON Schema, but direct validation is blocked until schema
  refs are packaged as a CUE module.
- `in-toto-run` can wrap a command, sign link metadata, and record material and
  product hashes.
- `cosign` can sign and verify a local `run.json` blob with a local key when
  transparency-log verification is explicitly disabled.
- `slsa-verifier` does not accept the local DSSE fixture as production SLSA
  evidence; it fails because no matching Rekor entries are found.

## Benchmark Snapshot

Local 20-iteration benchmark, compiled `sdp-trace`, Linux amd64:

| Probe | Median ms | Notes |
| --- | ---: | --- |
| shell prototype wrap | 6.0 | Minimal JSON, no hash chain semantics |
| `sdp-trace verify` | 8.0 | Existing local run |
| OPA adapter policy eval | 14.0 | Simplified policy |
| `sdp-trace wrap` | 26.0 | Local `/bin/true` |
| Cosign local verify | 30.5 | Transparency log ignored |
| `in-toto-run` | 148.0 | Signed link metadata |
| `check-jsonschema` fixture validation | 271.5 | Python validator startup cost |

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
