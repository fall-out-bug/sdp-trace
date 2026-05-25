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
| `opa` | Rego policy evaluation (requires v1.0+ for `import rego.v1`) | `not_assessed` |
| `cue` | CUE import / validation | `not_assessed` |
| `in-toto-run` | Signed supply-chain link generation | `not_assessed` |
| `cosign` | Local DSSE/blob signing & verification | `not_assessed` |
| `openssl` | Throwaway key generation for in-toto/Cosign local tests | `not_assessed` |
| `slsa-verifier` | SLSA provenance verifier negative test | `not_assessed` |

## Compatibility Probes

The table below reflects the **current environment**, where no optional
external CLI tools are installed. Run `go run ./tools/osscompat -json` to
produce live harness output for your environment. When the required optional
tools are available, the expected results are shown in the
**Expected (tool available)** column.

| Proxy / Tool | Capability Tested | Current Result | Expected (tool available) | Status |
|---|---|---|---|---|
| `check-jsonschema` | Validate `examples/flight-recorder/local-positive/run.json` against `flight-recorder-run.schema.json` | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `check-jsonschema` present [^1] |
| `check-jsonschema` | Validate live `run.json` artifact from `sdp-trace wrap` against `run-manifest.schema.json` | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `check-jsonschema` present [^1] |
| OPA/Rego | Express simplified adapter-capture pass rule | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `opa` present [^1] |
| OPA/Rego | Combined negative fixture (both rules broken) | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `opa` present [^1] |
| OPA/Rego | Negative trace_id rule only | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `opa` present [^1] |
| OPA/Rego | Negative provenance rule only | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `opa` present [^1] |
| CUE | JSON Schema import to stdout | `not_assessed` | `pass` | Not assessed in current environment; expected `pass` when `cue` present [^1] |
| CUE | Validate flight-recorder fixture via imported CUE | `cannot_verify` | `cannot_verify` | Docs-only/manual; no automated harness probe. Direct validation is blocked until schema refs are packaged as a CUE module |
| in-toto | Wrap command, sign link metadata, record material/product hashes | `not_assessed` | `cannot_verify` | Automated probe returns `cannot_verify` with reproducible reason when `in-toto-run` is present; no automated conformance verdict is issued [^1] [^2] |
| Cosign | Sign/verify local `run.json` blob | `not_assessed` | `cannot_verify` | Automated probe returns `cannot_verify` with reproducible reason when `cosign` is present; no automated conformance verdict is issued [^1] [^2] |
| Cosign | Verify with transparency log / Rekor | `not_assessed` | `cannot_verify` | Manual-only expected-fail probe; run reproduction command for evidence |
| SLSA verifier | Accept local DSSE fixture as production SLSA evidence | `not_assessed` | `cannot_verify` | Automated probe returns `cannot_verify` with reproducible reason when `slsa-verifier` is present; intentionally invalid local fixture (see reproduction command) [^1] [^2] |

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

### JSON Schema live wrap manifest validation

```bash
# Build the binary from the repo root, run it from an isolated temp dir, then
# validate the generated run.json artifact against the current live manifest
# schema.
(
  set -euo pipefail
  command -v check-jsonschema >/dev/null || { echo "check-jsonschema not found"; exit 1; }
  REPO_ROOT=$(git rev-parse --show-toplevel)
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT
  cd "$REPO_ROOT"
  go build -o "$TMPDIR/sdp-trace" ./cmd/sdp-trace
  cd "$TMPDIR"
  # Preflight: confirm the tool and schema are functional against a known-good
  # live-manifest fixture.
  check-jsonschema \
    --schemafile "$REPO_ROOT/schema/run-manifest.schema.json" \
    "$REPO_ROOT/examples/agentic-sdlc/local-wrap-positive/run.json"
  WRAP_OUT=$("$TMPDIR/sdp-trace" wrap true)
  if ! printf '%s' "$WRAP_OUT" | grep -qE '^run_dir: '; then
    echo "ERROR: unexpected wrap stdout: $WRAP_OUT"
    exit 1
  fi
  RUN_DIR=$(printf '%s' "$WRAP_OUT" | awk '{print $2}')
  RUN_JSON="$TMPDIR/$RUN_DIR/run.json"
  if [ ! -f "$RUN_JSON" ]; then
    echo "ERROR: run.json not found at $RUN_JSON"
    exit 1
  fi
  check-jsonschema \
    --schemafile "$REPO_ROOT/schema/run-manifest.schema.json" \
    "$RUN_JSON"
)
```

### OPA/Rego policy evaluation

```bash
# Evaluate a simplified adapter-capture rule against a test fixture.
# Assert the result is exactly true.
(
  set -e
  cd examples/oss-policy || exit 1
  RESULT=$(opa eval --data adapter.rego \
    --input test-fixture.json \
    --format raw \
    'data.sdp_trace.adapter.pass')
  if [ "$RESULT" != "true" ]; then
    echo "ERROR: expected true, got: $RESULT"
    exit 1
  fi
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
  set +e
  OUT=$(cosign verify-blob --key cosign.pub --signature run.json.sig run.json 2>&1)
  STATUS=$?
  set -e
  if [ "$STATUS" -eq 0 ]; then
    echo "ERROR: expected Rekor verification to fail"
    exit 1
  fi
  if [ "$STATUS" -ne 1 ]; then
    echo "ERROR: cosign exited abnormally (status=$STATUS); not a Rekor-missing failure"
    echo "$OUT"
    exit 1
  fi
  # Sanity-check that the failure relates to the transparency log.
  if ! printf '%s' "$OUT" | grep -qiE 'rekor|tlog|transparency'; then
    echo "WARNING: cosign failed but output does not mention Rekor/transparency log"
    echo "$OUT"
  fi
)
```

### SLSA verifier negative path (expected fail)

```bash
# Attempt to verify an intentionally invalid local DSSE fixture as production
# SLSA evidence. Expected to fail; the exact failure mode depends on the
# verifier version and is not separately evidenced.
(
  set -e
  command -v slsa-verifier >/dev/null || { echo "slsa-verifier not found"; exit 1; }
  REPO_ROOT=$(git rev-parse --show-toplevel)
  FIXTURE="$REPO_ROOT/examples/oss-supply-chain/local-dsse.json"
  if [ ! -f "$FIXTURE" ]; then
    echo "ERROR: fixture not found at $FIXTURE"
    exit 1
  fi
  set +e
  OUT=$(slsa-verifier verify-artifact \
    --provenance-path "$FIXTURE" \
    --source-uri local/test \
    /dev/null 2>&1)
  STATUS=$?
  set -e
  if [ "$STATUS" -eq 0 ]; then
    echo "ERROR: expected SLSA verification to fail"
    exit 1
  fi
  if [ "$STATUS" -ne 1 ]; then
    echo "ERROR: slsa-verifier exited abnormally (status=$STATUS); not an expected-validation failure"
    echo "$OUT"
    exit 1
  fi
)
```

## Substitution Boundaries

### JSON Schema (`check-jsonschema`)
- **Can replace:** External CI/schema validation steps against published `schema/` refs or example fixtures.
- **Adapter glue required:** Live `sdp-trace wrap` output is validated against
  `run-manifest.schema.json`. The richer `flight-recorder-run.schema.json`
  remains a separate run-level profile artifact and must not be inferred from
  the live manifest alone.
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
| JSON Schema fixture alignment | `not_assessed` | `check-jsonschema` not available in current environment; expected `pass` when tool present |
| OPA policy prototype | `not_assessed` | `opa` not available in current environment; expected `pass` when tool present |
| in-toto wrap probe | `cannot_verify` | Automated probe reports `cannot_verify` when `in-toto-run` present; see `tools/osscompat/probe.go` and `probe_test.go` |
| Cosign local sign probe | `cannot_verify` | Automated probe reports `cannot_verify` when `cosign` present; see `tools/osscompat/probe.go` and `probe_test.go` |
| SLSA verifier negative probe | `cannot_verify` | Automated probe reports `cannot_verify` when `slsa-verifier` present; see `tools/osscompat/probe.go` and `probe_test.go` |


All `not_assessed` states remain open until external, reproducible evidence is provided. Local fixture success does not imply production readiness or external trust. The in-toto, Cosign, and SLSA verifier probes are explicitly preserved as `cannot_verify` with reproducible state reporting: when the required tool is present, the probe runs and records an explicit, actionable reason; when the tool is absent, the harness reports `not_assessed`. Test coverage for these states lives in `tools/osscompat/probe_test.go` (`TestRunInTotoWrap_Direct`, `TestRunCosignLocalSign_Direct`, `TestRunSLSANegative_Direct`).

[^1]: Automated probes in `tools/osscompat` run the full validation when
safe to do so (JSON Schema fixture check, live wrap drift confirmation,
OPA eval against checked-in fixtures, CUE import). For probes that mutate
state or require external services, the tool reports `cannot_verify` and
the doc records the manual result. Run the reproduction commands above for
the actual validation.

[^2]: Direct probe function tests assert the `cannot_verify` state and verify
that the reason string contains explicit, actionable guidance for both the
tool-present and tool-missing cases. See `tools/osscompat/probe_test.go`.
