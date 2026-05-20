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

- `check-jsonschema` fixture validation: expected `pass` when tool available;
  in the current environment (no optional tools installed) the probe reports
  `not_assessed`. When available, it validates checked flight-recorder fixtures
  against locally rewired schema refs.
- Live `sdp-trace wrap` CLI stdout does not validate against
  `schema/flight-recorder-run.schema.json`. When `check-jsonschema` is
  available, the `osscompat` harness confirms this drift with a temp-dir probe
  that builds from source, runs `wrap`, captures stdout, and reports `fail`
  because the stdout is not schema JSON. When `check-jsonschema` is absent,
  the probe reports `not_assessed`.
- `OPA` adapter-capture rule evaluation: expected `pass` when tool available;
  in the current environment the probe reports `not_assessed`. When available,
  the checked-in positive fixture evaluates to true and negative fixtures
  evaluate to false.
- `CUE` JSON Schema import: expected `pass` when tool available; in the
  current environment the probe reports `not_assessed`. Direct validation is
  blocked until schema refs are packaged as a CUE module.
- `in-toto-run` command wrapping: automated harness probe reports
  `cannot_verify` (manual-only). No automated conformance verdict is issued.
- `cosign` local blob sign/verify: automated harness probe reports
  `cannot_verify` (manual-only). No automated conformance verdict is issued.
- `slsa-verifier` local DSSE fixture: automated harness probe reports
  `cannot_verify` (manual-only expected-fail). No automated conformance
  verdict is issued.

## Benchmark Snapshot

Local 20-iteration benchmark, compiled `sdp-trace` from source, Linux amd64.
All rows below are reproducible via `tools/ossbench`.

| Probe | Median ms | Min ms | Max ms | Command |
| --- | ---: | ---: | ---: | --- |
| `sdp-trace version` | 4.60 | 4.39 | 5.03 | `sdp-trace version` |
| `sdp-trace wrap` | 16.08 | 15.64 | 17.53 | `sdp-trace wrap /bin/true` |

The harness builds `sdp-trace` from source into a temp directory on every run;
the displayed command is a display name (`filepath.Base`), while the actual
binary path is recorded in JSON `binary_path` and `binary_source`. Reproduce
with: `go run ./tools/ossbench -json -n 20`.

One-shot historical numbers for external tools (OPA, Cosign, in-toto,
check-jsonschema) are not included because min/max were not preserved and
the harness does not yet automate those probes. See `docs/oss-benchmark-results.md`
for scope context.

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
