# Tasks: OSS Replacement Compatibility And Benchmarks

Status: draft

## Phase 0 - Review

- [ ] T017-001 Review substitution candidates and decide which are in scope for
  first implementation.
- [ ] T017-002 Decide whether compatibility tooling is product Go code under
  `tools/` or docs-only reproducible shell snippets.

## Phase 1 - Pi-Ready Workstreams

- [ ] T017-010 WS-017-A: Create reproducible OSS compatibility harness.
  Pi ownership: `tools/osscompat/*` or one docs page.
- [ ] T017-020 WS-017-B: Resolve live `wrap` output vs flight-recorder schema
  compatibility. Pi ownership: schema/examples plus focused recorder tests.
- [ ] T017-030 WS-017-C: Build OPA/Rego or CUE assessment-profile prototype.
  Pi ownership: docs/examples for policy prototype only.
- [ ] T017-040 WS-017-D: Build in-toto/Cosign/SLSA supply-chain prototype.
  Pi ownership: docs/examples for supply-chain prototype only.
- [ ] T017-050 WS-017-E: Add benchmark harness and benchmark report.
  Pi ownership: `tools/ossbench/*` or benchmark docs only.

## Phase 2 - Integration

- [ ] T017-060 Run all compatibility probes available in the local environment.
- [ ] T017-070 Mark unavailable external services as `not_assessed` or
  `cannot_verify`; do not infer pass from local fixture success.
- [ ] T017-080 Update roadmap and docs index after accepted implementation.
