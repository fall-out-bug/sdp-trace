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
| `slsa-verifier` | SLSA provenance verifier negative test | `not_assessed` |

## Compatibility Probes

| Proxy / Tool | Capability Tested | Probe Result | Status |
|---|---|---|---|
| `check-jsonschema` | Validate flight-recorder fixtures against local schema refs | `pass` | Local refs resolve; validation succeeds on checked examples |
| `check-jsonschema` | Validate live `sdp-trace wrap` output vs `flight-recorder-run.schema.json` | `fail` | Required fields and timestamp format differ |
| `check-jsonschema` | Validate representative assessment, gate, and release examples | `pass` | Example fixtures conform to schema |
| OPA/Rego | Express simplified adapter-capture pass/fail rule | `pass` | Policy evaluates and correctly detects `test_provenance_not_overclaimed` failure fixture |
| CUE | JSON Schema import | `cannot_verify` | CUE can import JSON Schema types, but direct validation is blocked until schema refs are packaged as a CUE module |
| in-toto | Wrap command, sign link metadata, record material/product hashes | `pass` | Link metadata generated and signed locally |
| Cosign | Sign/verify local `run.json` blob | `pass` | Works with local key when transparency log verification is explicitly disabled |
| Cosign | Verify with transparency log / Rekor | `fail` | Expected for local-only fixtures; no Rekor entry exists |
| SLSA verifier | Accept local DSSE fixture as production SLSA evidence | `fail` | Expected: no matching Rekor entries found |

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
(
  sdp-trace wrap /bin/true > /tmp/wrap-run.json
  check-jsonschema \
    --schemafile schema/flight-recorder-run.schema.json \
    /tmp/wrap-run.json || true
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
  cue import --package sdptrace schema/flight-recorder-run.schema.json \
    -o /tmp/flight-recorder-run.cue || true
)
```

### in-toto command wrapping

```bash
# Generate a throwaway key and wrap a command
(
  cd /tmp
  openssl genpkey -algorithm RSA -out test-key.pem 2>/dev/null
  in-toto-run \
    --step-name test-wrap \
    --products /dev/null \
    --key test-key.pem \
    -- /bin/true || true
)
```

### Cosign local blob sign/verify

```bash
# Sign and verify a local JSON blob with a generated ephemeral key
(
  cd /tmp
  cosign generate-key-pair || true
  cosign sign-blob --key cosign.key --yes run.json > run.json.sig || true
  cosign verify-blob --key cosign.pub --signature run.json.sig run.json \
    || true
)
```

### SLSA verifier negative path (expected fail)

```bash
# Attempt to verify a local DSSE fixture as production SLSA evidence.
# Expected to fail because no Rekor entry exists.
(
  slsa-verifier verify-artifact \
    --provenance-path examples/oss-supply-chain/local-dsse.json \
    --source-uri local/test \
    /dev/null || true
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
| OPA policy prototype | `pass` | Local evaluation passes/fails as expected on test fixtures |

All `not_assessed` states remain open until external, reproducible evidence is provided. Local fixture success does not imply production readiness or external trust.
