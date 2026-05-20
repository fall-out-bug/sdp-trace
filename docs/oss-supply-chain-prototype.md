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
| in-toto-run wrap + sign | `pass` | Link metadata generated and signed locally |
| Cosign local blob sign/verify | `pass` | Works with local key when transparency log is disabled |
| Cosign verify with Rekor | `fail` | Expected: no Rekor entry for local fixture |
| SLSA verifier accept local DSSE | `fail` | Expected: no matching Rekor entries |
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
