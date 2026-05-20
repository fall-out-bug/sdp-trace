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
- [x] T017-020 WS-017-B: Document live `wrap` output vs flight-recorder schema
  compatibility drift and add structural evidence. Pi ownership: schema/examples
  plus focused recorder tests. Fixing the drift remains out of scope.
- [x] T017-030 WS-017-C: Build OPA/Rego or CUE assessment-profile prototype.
  Pi ownership: docs/examples for policy prototype only.
- [x] T017-040 WS-017-D: Build in-toto/Cosign/SLSA supply-chain prototype.
  Pi ownership: docs/examples for supply-chain prototype only.
- [x] T017-050 WS-017-E: Add benchmark harness and benchmark report.
  Pi ownership: `tools/ossbench/*` and benchmark docs.

## Phase 2 - Integration

- [x] T017-060 Run all compatibility probes available in the local environment.
- [x] T017-070 Mark unavailable external services as `not_assessed` or
  `cannot_verify`; do not infer pass from local fixture success.
- [ ] T017-080 Update roadmap and docs index after accepted implementation.
- [x] T017-090 Keep the live `wrap` output vs
  `schema/flight-recorder-run.schema.json` drift open as a blocker until
  T017-020 lands; fixture schema compatibility must not be claimed as live
  recorder compatibility.
