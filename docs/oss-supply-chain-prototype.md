# OSS Supply-Chain Prototype

Status: in_progress
Spec: [017](../specs/017-oss-replacement-compatibility-and-benchmarks/)

This document records local supply-chain tooling probes: in-toto command
wrapping, Cosign blob signing, and SLSA verifier negative-path testing.
These are **local experiments only**; they do not approve replacing
sdp-trace witness, checkpoint, or release-proof semantics.

## Probe Results

| Probe | State | Reason |
|---|---|---|
| in-toto-run wrap + sign | `cannot_verify` | Manual-only; no automated harness probe run |
| Cosign local blob sign/verify | `cannot_verify` | Manual-only; no automated harness probe run |
| Cosign verify with Rekor | `cannot_verify` | Manual-only; run reproduction command for expected-fail evidence |
| SLSA verifier reject local DSSE | `cannot_verify` | Manual-only; run reproduction command for expected-fail evidence |
| SLSA/Rekor production trust | `not_assessed` | No live external provenance or Rekor inclusion |

## Files

- `examples/oss-supply-chain/local-dsse.json`
- `examples/oss-supply-chain/README.md`

## Substitution Boundary

### What OSS Tools Can Replace

- Local command wrapping and metadata generation for experiments.
- Blob signing and verification with local keys for non-production use.

### What Needs Adapter Glue

- Production trust requires OIDC identity binding, Rekor transparency log
  inclusion, and trusted identity policy verification.
- `witness` contract semantics (what counts as evidence, how checkpoints
  are linked) are sdp-trace-specific.
- `release-proof` source-bound verification requires a clean immutable
  source commit; Cosign local signing does not provide this.

### What Remains sdp-trace-Specific

- `checkpoint` logic and hash-chain semantics.
- Gate verdict interpretation.
- Evidence retention policies and artifact roles.

## Non-Goals

- No production Sigstore, Rekor, or SLSA trust claim from local fixtures.
- No replacement of `witness`, `release-proof`, or `checkpoint` without
  external identity-policy and transparency-log evidence.
- No benchmark health score from supply-chain tool latencies.
