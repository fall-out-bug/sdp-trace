# OSS Benchmark Results

Status: in_progress
Spec: [017](../specs/017-oss-replacement-compatibility-and-benchmarks/)

This document records local benchmark numbers for OSS tools and `sdp-trace`
commands. Numbers are **non-authoritative**: they inform scope and complexity
decisions but do not prove production readiness.

## Environment

- Platform: Linux amd64
- CPU: AMD Ryzen 7 5800X 8-Core Processor exposed as 6 vCPUs
- Memory: 47 GiB available host memory
- Iterations: 20 per probe
- Reported metric: **median** wall-clock time per invocation
- Background load: not recorded

## Benchmark Table

| Probe | Median (ms) | Min (ms) | Max (ms) | Iterations | Exact Command | Notes |
|---|---:|---:|---:|---|---|---|
| Shell prototype `wrap` | 6.0 | — | — | 20 | `bash -c 'echo {"v":"local"}'` | Minimal JSON, no hash-chain semantics |
| `sdp-trace verify` | 8.0 | — | — | 20 | `sdp-trace verify` | Existing local run |
| OPA adapter policy eval | 14.0 | — | — | 20 | `opa eval --data adapter.rego --input fixture.json 'data.sdp_trace.adapter.pass'` | Simplified policy |
| `sdp-trace wrap` | 26.0 | — | — | 20 | `sdp-trace wrap /bin/true` | Local `/bin/true` |
| Cosign local verify | 30.5 | — | — | 20 | `cosign verify-blob --key cosign.pub --signature run.json.sig run.json` | Transparency log ignored |
| `in-toto-run` | 148.0 | — | — | 20 | `in-toto-run --step-name test --products /dev/null --key key.pem -- /bin/true` | Signed link metadata |
| `check-jsonschema` fixture validation | 271.5 | — | — | 20 | `check-jsonschema --schemafile schema/flight-recorder-run.schema.json examples/...` | Python validator startup cost |

**Note:** Min and max values are marked `—` because the raw iteration data
from the 2026-05-20 one-shot run was not preserved. This table is a **local
markdown ledger only** and does not satisfy FR-017-004 until `tools/ossbench`
produces reproducible structured output with full statistics. Do not use these
numbers for approval decisions.

## Observations

**Compiled Go is fastest for core operations.** The shell prototype (6.0 ms)
and `sdp-trace verify` (8.0 ms) represent floor latency for minimal JSON
handling and local verification. External tools incur higher base overhead.

**OPA policy evaluation is lightweight.** At 14 ms median, Rego evaluation
on a simplified adapter-capture rule adds negligible overhead compared to
`sdp-trace wrap` (26 ms). This supports using OPA as a policy layer if
the translation adapter remains small.

**CUE performance not measured.** Direct CUE validation is blocked until
`schema/` files are packaged as CUE modules. CUE compatibility is recorded
in the [compatibility doc](oss-replacement-compatibility.md) under Status
`cannot_verify` for this dimension.

**Supply-chain tools are heavier.** `in-toto-run` (148 ms) and Cosign
local verify (30.5 ms) are significantly slower than core `sdp-trace`
commands. This is expected: both tools perform cryptographic operations
(key generation, signing, metadata serialization) that `sdp-trace` wrap
does not replicate. The overhead is acceptable for release-proof or
witness steps, but too high for per-event recording.

`check-jsonschema` (271.5 ms) is dominated by Python interpreter startup.
For CI schema checks this is acceptable; it is not suitable as a per-event
recording validator in hot paths.

## Commands Measured

The following commands were invoked from the repository root in a subshell
unless noted otherwise. They are provided so a reader can reproduce the
measurement protocol, not to claim the exact same medians will hold on
another machine.

```bash
# Shell prototype wrap
bash -c 'echo {"v":"local"}'

# sdp-trace verify
sdp-trace verify

# OPA adapter policy eval
opa eval --data examples/oss-policy/adapter.rego \
  --input examples/oss-policy/test-fixture.json \
  'data.sdp_trace.adapter.pass'

# sdp-trace wrap
sdp-trace wrap /bin/true

# Cosign local verify
cosign verify-blob --key /tmp/cosign.pub \
  --signature /tmp/run.json.sig /tmp/run.json

# in-toto-run
in-toto-run --step-name test-wrap --products /dev/null \
  --key /tmp/key.pem -- /bin/true

# check-jsonschema fixture validation
check-jsonschema \
  --schemafile schema/flight-recorder-run.schema.json \
  examples/flight-recorder/local-wrap-positive/run.json
```

## Non-Authoritative Disclaimer

These numbers come from a local Linux amd64 environment with 20 iterations
per probe. They do not account for:

- Concurrent load or multi-process contention
- Network latency for external services (Rekor, OCI registry, OIDC provider)
- Production data sizes or artifact volumes
- Cold starts vs. warmed caches beyond the measured median
- System-level jitter from other processes
- Different CPU architectures or constrained environments (containers, CI)
- Hardware and background-load differences; the original run did not capture a
  full performance profile

Benchmark results must not be used to assign health scores, readiness gates,
or production trust decisions. They are scope-informative only.

## Follow-Ups

- **Reproducible harness:** T017-050 tracks adding a checked-in benchmark
  harness before this table is used for approval. The current table is a local
  markdown ledger only.
- **CUE module packaging:** Once `schema/` files are exported as CUE modules,
  add `cue vet` benchmark numbers to this table.
- **SLSA verifier negative benchmark:** Time `slsa-verifier` rejecting local
  DSSE fixture to understand its rejection latency vs. acceptance latency
  for future production scenarios.
- **Reproducible benchmark harness:** A Go tool under `tools/ossbench/` could
  automate the 20-iteration protocol and emit structured JSON output.
