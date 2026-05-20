# OSS Replacement Compatibility

Status: in_progress
Spec: [017](../specs/017-oss-replacement-compatibility-and-benchmarks/)

This document records local compatibility probes for candidate OSS replacements
and maps substitution boundaries. It is a **docs-only first slice**; it proves
that probes can run locally, not that sdp-trace should replace any product code.

## Scope

- Covers JSON Schema validation, OPA/Rego policy, CUE import, in-toto
  command wrapping, Cosign blob signing, and SLSA verifier negative path.
- Non-goals: no automatic migration, no production Sigstore/Rekor/SLSA trust
  claim from local fixtures, no benchmark health score, no Node.js/npm tooling.

## Prerequisites

To reproduce probes, the following CLI tools must be available in `$PATH`.
Missing tools render that probe `not_assessed`.

| Tool | Purpose | Verdict if missing |
|---|---|---|
| `check-jsonschema` (Python) | JSON Schema fixture validation | `not_assessed` |
| `opa` | Rego policy evaluation | `not_assessed` |
| `cue` | CUE import / validation | `not_assessed` |
| `in-toto-run` | Signed supply-chain link generation | `not_assessed` |
| `cosign` | Local DSSE/blob signing & verification | `not_assessed` |
| `openssl` | Throwaway key generation for in-toto/Cosign local tests | `not_assessed` |
| `slsa-verifier` | SLSA provenance verifier negative test | `not_assessed` |

## Compatibility Probes

| Proxy / Tool | Capability Tested | Probe Result | Status |
|---|---|---|---|
| `check-jsonschema` | Validate `examples/flight-recorder/local-positive/run.json` against local schema refs | `pass` | Checked fixture conforms to schema [^1] |
| `check-jsonschema` | Validate live `sdp-trace wrap` output vs `flight-recorder-run.schema.json` | `cannot_verify` | Live wrap output differs from schema; safe auto-run is not implemented so the automated probe reports `cannot_verify` [^1] |
| `check-jsonschema` | Validate `examples/flight-recorder/local-positive/run.json` against `flight-recorder-run.schema.json` | `pass` | Checked fixture conforms to schema [^1] |
| OPA/Rego | Express simplified adapter-capture pass rule | `pass` | Policy evaluates the checked-in pass fixture as expected [^1] |
| CUE | JSON Schema import to stdout | `pass` | `cue import` succeeds against `schema/flight-recorder-run.schema.json` without mutating the working tree [^1] |
| CUE | Validate flight-recorder fixture via imported CUE | `cannot_verify` | Direct validation is blocked until schema refs are packaged as a CUE module |
| in-toto | Wrap command, sign link metadata, record material/product hashes | `cannot_verify` | Manual-only; no automated harness probe run [^1] |
| Cosign | Sign/verify local `run.json` blob | `cannot_verify` | Manual-only; no automated harness probe run [^1] |
| Cosign | Verify with transparency log / Rekor | `fail` | Expected for local-only fixtures; no Rekor entry exists |
| SLSA verifier | Accept local DSSE fixture as production SLSA evidence | `fail` | Expected: truncated signature prevents verification before any Rekor check |

## Reproduction Commands

The following commands can be run from a clean checkout at the repository root.
Each command uses subshell isolation where the working directory matters.
Missing tools produce a non-zero exit code; probe scripts should catch this
and report `not_assessed`.

### JSON Schema fixture validation

```bash
# Validate a representative flight-recorder fixture
(
  check-jsonschema \
    --schemafile schema/flight-recorder-run.schema.json \
    examples/flight-recorder/local-positive/run.json
)
```

### JSON Schema live wrap drift (expected fail)

```bash
# This command is expected to fail until T017-020 resolves the drift.
# Use an isolated temp directory so wrap does not leave state behind.
(
  set -e
  REPO_ROOT=$(git rev-parse --show-toplevel)
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT
  cd "$TMPDIR"
  WRAP_OUT=$(mktemp -p "$TMPDIR")
  trap 'rm -rf "$TMPDIR" "$WRAP_OUT"' EXIT
  go run "$REPO_ROOT/cmd/sdp-trace" wrap /bin/true > "$WRAP_OUT"
  check-jsonschema \
    --schemafile "$REPO_ROOT/schema/flight-recorder-run.schema.json" \
    "$WRAP_OUT"
)
```

### OPA/Rego policy evaluation

```bash
# Evaluate a simplified adapter-capture rule against a test fixture
(
  cd examples/oss-policy || exit 1
  opa eval --data adapter.rego --input test-fixture.json \
    'data.sdp_trace.adapter.pass'
)
```

### CUE JSON Schema import

```bash
# Import JSON Schema types into CUE (does not validate sdp-trace artifacts)
(
  cue import --package sdptrace schema/flight-recorder-run.schema.json -o -
)
```

### in-toto command wrapping

```bash
# Generate a throwaway key and wrap a command
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

### Cosign local blob sign/verify

```bash
# Sign and verify a local JSON blob with a generated ephemeral key.
# Transparency-log upload and verification are disabled for local-only testing.
(
  set -e
  export COSIGN_PASSWORD=""
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT
  cd "$TMPDIR"
  printf '{"run":"test"}\n' > run.json
  cosign generate-key-pair
  cosign sign-blob --key cosign.key --yes --tlog-upload=false run.json > run.json.sig
  cosign verify-blob --key cosign.pub --signature run.json.sig --insecure-ignore-tlog run.json
)
```

### Cosign Rekor verification (expected fail)

```bash
# Attempt to verify a locally-signed blob against the public Rekor log.
# Expected to fail because the blob was never uploaded to Rekor.
(
  set -e
  export COSIGN_PASSWORD=""
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT
  cd "$TMPDIR"
  printf '{"run":"test"}\n' > run.json
  cosign generate-key-pair
  cosign sign-blob --key cosign.key --yes --tlog-upload=false run.json > run.json.sig
  # This command omits --insecure-ignore-tlog so Rekor verification is attempted.
  if cosign verify-blob --key cosign.pub --signature run.json.sig run.json; then
    echo "ERROR: expected Rekor verification to fail"
    exit 1
  fi
)
```

### SLSA verifier negative path (expected fail)

```bash
# Attempt to verify a local DSSE fixture as production SLSA evidence.
# Expected to fail because the signature is truncated and no Rekor entry exists.
(
  slsa-verifier verify-artifact \
    --provenance-path examples/oss-supply-chain/local-dsse.json \
    --source-uri local/test \
    /dev/null
)
```

## Substitution Boundaries

### JSON Schema (`check-jsonschema`)
- **Can replace:** External CI/schema validation steps against published `schema/` refs or example fixtures.
- **Adapter glue required:** Live `sdp-trace wrap` output currently mismatches `flight-recorder-run.schema.json`. Requires either fixing the recorder schema alignment or a separate current recorder schema.
- **Remains sdp-trace-specific:** Embedded Go validation, hash-chain semantics, and internal recorder profiles.

### OPA/Rego
- **Can replace:** Policy-as-code expression for simplified assessment profiles (e.g., adapter-capture pass/fail).
- **Adapter glue required:** OPA does not natively understand sdp-trace gate verdicts, trace provenance, or `sdp-trace-claim` tags. A translation layer or JSON-mapping layer is required.
- **Remains sdp-trace-specific:** Product verifier behavior (`assess`, `verify`, `witness`) is not a direct OPA evaluation target; the verifier controls evidence collection.

### CUE
- **Can replace:** Schema authoring and validation in a future workflow if schemas are exported as CUE modules.
- **Adapter glue required:** `schema/` JSON files must be packaged or generated as CUE modules before `cue vet` can validate sdp-trace artifacts.
- **Remains sdp-trace-specific:** Current schema contracts are JSON Schema-first; CUE compatibility is import-only until module packaging is implemented.

### Supply Chain OSS (in-toto / Cosign / SLSA)
- **Can replace:** Local command wrapping and blob signing experiments.
- **Adapter glue required:** Production equivalence requires OIDC identity binding, Rekor transparency log inclusion, and trusted identity policy verification. Local fixtures do not provide this.
- **Remains sdp-trace-specific:** `witness` contract, `release-proof` source-bound verification, and `checkpoint` logic are sdp-trace domains. OSS tools provide the signing/timestamping substrate but not the evidence interpretation or gate verdicts.

## Trust And Verification Status

| Domain | Status | Reason |
|---|---|---|
| SLSA production trust | `not_assessed` | No live external SLSA provenance or Rekor inclusion provided |
| Rekor integration | `not_assessed` | Transparency log requires external service and OIDC identity |
| CUE module packaging | `cannot_verify` | `schema/` refs not yet packaged as CUE modules |
| Cosign production trust | `not_assessed` | Local key signing excludes keyless/OIDC/transparency verification |
| JSON Schema fixture alignment | `pass` | Checked examples validate; live wrap output mismatch documented |
| OPA policy prototype | `pass` | Local evaluation of the checked-in positive and negative fixtures returns expected true/false |

All `not_assessed` states remain open until external, reproducible evidence is provided. Local fixture success does not imply production readiness or external trust.

[^1]: Automated probes in `tools/osscompat` run the full validation when
safe to do so (JSON Schema fixture check, OPA eval against checked-in
fixtures, CUE import). For probes that mutate state or require external
services, the tool reports `cannot_verify` and the doc records the manual
result. Run the reproduction commands above for the actual validation.
