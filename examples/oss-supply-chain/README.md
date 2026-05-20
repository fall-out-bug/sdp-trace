# OSS Supply-Chain Prototype

Status: locally tested, not externally verified
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
# Generate a throwaway key for local testing
(
  set -e
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT
  cd "$TMPDIR"
  openssl genpkey -algorithm RSA -out test-key.pem 2>/dev/null
  in-toto-run \
    --step-name test-wrap \
    --products /dev/null \
    --key test-key.pem \
    -- /bin/true
)
```

## Cosign Local Blob Sign/Verify

```bash
# Run in a subshell to avoid mutating the caller's CWD.
# Transparency-log upload and verification are disabled for local-only testing.
(
  set -e
  export COSIGN_PASSWORD=""
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT
  cd "$TMPDIR"
  cosign generate-key-pair
  printf '{"run":"test"}\n' > run.json
  cosign sign-blob --key cosign.key --yes --tlog-upload=false run.json > run.json.sig
  cosign verify-blob --key cosign.pub --signature run.json.sig --insecure-ignore-tlog run.json
)
```

## SLSA Verifier Negative Path

```bash
# From the repo root:
slsa-verifier verify-artifact \
  --provenance-path examples/oss-supply-chain/local-dsse.json \
  --source-uri local/test \
  /dev/null
```

Expected: failure (truncated signature, no Rekor entry, untrusted key).

## Substitution Boundary

- **What OSS tools replace:** Local command wrapping, blob signing, and
  metadata generation for experiments.
- **What remains sdp-trace-specific:** `witness` contract, `release-proof`
  source-bound verification, `checkpoint` logic, gate verdict semantics.
- **Adapter glue required:** Production equivalence needs OIDC identity
  binding, Rekor transparency log inclusion, and trusted identity policy.
  Local fixtures do not provide any of these.
