# Tasks: OSS Replacement Compatibility And Benchmarks

Status: in_progress

## Phase 0 - Review

- [x] T017-001 Verify the spec-approved substitution candidates remain the only
  in-scope OSS tools for first implementation.
- [x] T017-002 Implement compatibility tooling under the Go active path unless
  a later spec explicitly narrows a probe to docs-only evidence.

## Phase 1 - Pi-Ready Workstreams

- [x] T017-010 WS-017-A: Create reproducible OSS compatibility harness.
  Pi ownership: `tools/osscompat/*`.
- [x] T017-020 WS-017-B: Verify and document live `wrap` output vs
  flight-recorder schema compatibility drift. The `osscompat` harness confirms
  the drift when `check-jsonschema` is available (`fail`); when the tool is
  absent the probe reports `not_assessed`. Pi ownership: schema/examples plus
  focused recorder tests. Fixing the drift remains out of scope.
- [x] T017-030 WS-017-C: Build OPA/Rego or CUE assessment-profile prototype.
  Pi ownership: docs/examples for policy prototype only.
- [ ] T017-040 WS-017-D: Build in-toto/Cosign/SLSA supply-chain prototype.
  Pi ownership: docs/examples for supply-chain prototype only.
  Status: docs/examples complete; automated harness probes for in-toto wrapping,
  Cosign blob signing, and SLSA negative path remain open. Current probes are
  manual-only (`cannot_verify`). Full automated coverage is deferred to a
  future spec update.
- [x] T017-050 WS-017-E: Add benchmark harness and benchmark report.
  Pi ownership: `tools/ossbench/*` and benchmark docs.
  Note: the harness satisfies FR-017-004 structurally for sdp-trace built-in
  benchmarks. The markdown table in `docs/oss-benchmark-results.md` carries
  reproducible numbers from the harness (see spec benchmark snapshot).
  External-tool benchmark rows are explicitly deferred; they are not automated
  in this slice and are not claimed as complete.

## Phase 2 - Integration

- [x] T017-060 Run all compatibility probes available in the local environment.
  Verification: `go test -count=1 ./tools/osscompat` and `go test -count=1
  ./tools/ossbench` pass locally. Example harness output (current environment,
  all optional tools absent): `go run ./tools/osscompat -json` reports ten
  `not_assessed` probes because no optional external tools are installed; this
  is the expected behavior for an environment without the optional CLI
  prerequisites. `go run ./tools/ossbench -json` reports reproducible
  `min_ms`/`max_ms`/`median_ms` for both built-in benchmarks.
- [x] T017-070 Mark unavailable external services as `not_assessed` or
  `cannot_verify`; do not infer pass from local fixture success.
- [ ] T017-080 Update roadmap and docs index after accepted implementation.
- [ ] T017-090 Keep the `wrap` output vs
  `schema/flight-recorder-run.schema.json` drift open as a blocker until a
  source-bound fix (schema update, wrap output change, or new recorder schema)
  lands in a subsequent spec. Fixture schema compatibility must not be claimed
  as live recorder compatibility.
