# OSS Supply-Chain Prototype

Status: `local_pass`  
Spec: [017](../../specs/017-oss-replacement-compatibility-and-benchmarks/)

This directory contains minimal supply-chain tooling probes:
in-toto link metadata, Cosign blob signing, and SLSA verifier
negative-path testing. These are local experiments only.

## Files

| File | Purpose |
|---|---|
| `local-dsse.json` | Local DSSE-like fixture for SLSA verifier negative test |

## Prerequisites

- `in-toto-run`
- `cosign`
- `slsa-verifier`

## in-toto Command Wrapping

```bash
in-toto-run \
  --step-name test-wrap \
  --products /dev/null \
  --key /tmp/test-key.pem \
  -- /bin/true
```

## Cosign Local Blob Sign/Verify

```bash
cd /tmp
cosign generate-key-pair
echo '{"run":"test"}' > run.json
cosign sign-blob --key cosign.key --yes run.json > run.json.sig
cosign verify-blob --key cosign.pub --signature run.json.sig run.json
```

## SLSA Verifier Negative Path

```bash
slsa-verifier verify-artifact \
  --provenance-path examples/oss-supply-chain/local-dsse.json \
  --source-uri local/test \
  /dev/null
```

Expected: failure because no Rekor entry exists.

## Substitution Boundary

- **What OSS tools replace:** Local command wrapping, blob signing, and
  metadata generation for experiments.
- **What remains sdp-trace-specific:** `witness` contract, `release-proof`
  source-bound verification, `checkpoint` logic, gate verdict semantics.
- **Adapter glue required:** Production equivalence needs OIDC identity
  binding, Rekor transparency log inclusion, and trusted identity policy.
  Local fixtures do not provide any of these.
