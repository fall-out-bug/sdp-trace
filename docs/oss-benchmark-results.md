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
| `sdp-trace version` | 4.6 | 4.4 | 5.1 | 10 | `sdp-trace version` | Built-in, measured via `tools/ossbench` |
| `sdp-trace wrap` | 16.1 | 15.8 | 16.4 | 10 | `sdp-trace wrap /bin/true` | Built-in, measured via `tools/ossbench` |
| Shell prototype `wrap` | 6.0 | — | — | 20 | `printf '%s\n' '{"v":"local"}'` | One-shot; minimal JSON, no hash-chain semantics |
| OPA adapter policy eval | 14.0 | — | — | 20 | `opa eval --data adapter.rego --input fixture.json 'data.sdp_trace.adapter.pass'` | One-shot; simplified policy |
| Cosign local verify | 30.5 | — | — | 20 | `cosign verify-blob --key cosign.pub --signature run.json.sig run.json` | One-shot; transparency log ignored |
| `in-toto-run` | 148.0 | — | — | 20 | `in-toto-run --step-name test --products /dev/null --key key.pem -- /bin/true` | One-shot; signed link metadata |
| `check-jsonschema` fixture validation | 271.5 | — | — | 20 | `check-jsonschema --schemafile schema/flight-recorder-run.schema.json examples/...` | One-shot; Python validator startup cost |

**Note:** Min and max are present only for probes measured by the built-in
`tools/ossbench` harness. Rows marked "One-shot" retain `—` because raw
iteration data from the original run was not preserved. This table is a
**local markdown ledger only**. Do not use these numbers for approval decisions.

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
# sdp-trace version (built-in harness probe)
sdp-trace version

# sdp-trace wrap (built-in harness probe)
sdp-trace wrap /bin/true

# Shell prototype wrap
printf '%s\n' '{"v":"local"}'

# OPA adapter policy eval
opa eval --data examples/oss-policy/adapter.rego \
  --input examples/oss-policy/test-fixture.json \
  'data.sdp_trace.adapter.pass'

# Cosign local verify
cosign verify-blob --key /tmp/cosign.pub \
  --signature /tmp/run.json.sig /tmp/run.json

# in-toto-run
in-toto-run --step-name test-wrap --products /dev/null \
  --key /tmp/key.pem -- /bin/true

# check-jsonschema fixture validation
check-jsonschema \
  --schemafile schema/flight-recorder-run.schema.json \
  examples/flight-recorder/local-positive/run.json
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

- **Built-in harness:** `tools/ossbench` now automates the 20-iteration
  protocol for built-in `sdp-trace` probes and emits structured JSON with
  `min_ms`, `max_ms`, `median_ms`, and `iterations`. Run with `-json` for
  machine-readable output.
- **External-tool benchmarks:** Probes requiring `opa`, `cosign`, `in-toto-run`,
  `check-jsonschema`, or `slsa-verifier` are not yet covered by the harness
  because they depend on optional external CLIs.
- **CUE module packaging:** Once `schema/` files are exported as CUE modules,
  add `cue vet` benchmark numbers to this table.
- **SLSA verifier negative benchmark:** Time `slsa-verifier` rejecting local
  DSSE fixture to understand its rejection latency vs. acceptance latency
  for future production scenarios.
