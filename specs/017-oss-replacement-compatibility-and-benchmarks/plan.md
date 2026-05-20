# Plan: OSS Replacement Compatibility And Benchmarks

Status: draft

## Workstreams

### WS-017-A: Reproducible Probe Harness

Owned files:

- `tools/osscompat/*`
- `docs/oss-replacement-compatibility.md`
- no product dependency on Node.js/npm

Deliverable:

- A command that runs or explains compatibility probes without mutating tracked
  product artifacts.

### WS-017-B: Schema Compatibility

Owned files:

- `schema/*`
- `examples/flight-recorder/*`
- `docs/agent-entrypoint.md` only if command/schema contract changes

Deliverable:

- Implement the spec decision that live recorder schema compatibility remains
  blocked until `wrap` output conforms to `flight-recorder-run.schema.json` or
  a separate current recorder schema is defined.

### WS-017-C: Policy-As-Code Prototype

Owned files:

- `docs/oss-policy-prototype.md`
- optional `examples/oss-policy/*`

Deliverable:

- OPA/Rego or CUE prototype for one assessment profile, with a clear boundary
  between policy and product verifier behavior.

### WS-017-D: Supply Chain OSS Prototype

Owned files:

- `docs/oss-supply-chain-prototype.md`
- optional `examples/oss-supply-chain/*`

Deliverable:

- in-toto/Cosign/SLSA probe records that show what replaces witness,
  checkpoint, and release-proof pieces, and what does not.

### WS-017-E: Benchmark Harness

Owned files:

- `tools/ossbench/*`
- `docs/oss-benchmark-results.md`

Deliverable:

- Repeatable local benchmarks with environment, command lines, iterations, and
  median/min/max output.

## Verification

```text
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

Additional optional commands depend on installed OSS tools and must report
`blocked` or `not_assessed` if unavailable.

## Pi Handoff Notes

Assign WS-017-B only to a worker allowed to touch schema/examples. Assign
WS-017-C and WS-017-D to separate workers because policy and supply-chain
experiments have different dependencies.
