# OSS Replacement Compatibility

Status: draft
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
Missing tools render that probe `blocked`.

| Tool | Purpose | Verdict if missing |
|---|---|---|
| `check-jsonschema` (Python) | JSON Schema fixture validation | `blocked` |
| `opa` | Rego policy evaluation | `blocked` |
| `cue` | CUE import / validation | `blocked` |
| `in-toto-run` | Signed supply-chain link generation | `blocked` |
| `cosign` | Local DSSE/blob signing & verification | `blocked` |
| `slsa-verifier` | SLSA provenance verifier negative test | `blocked` |

## Compatibility Probes

| Proxy / Tool | Capability Tested | Probe Result | Status |
|---|---|---|---|
| `check-jsonschema` | Validate flight-recorder fixtures against local schema refs | `pass` | Local refs resolve; validation succeeds on checked examples |
| `check-jsonschema` | Validate live `sdp-trace wrap` output vs `flight-recorder-run.schema.json` | `fail` | Required fields and timestamp format differ |
| `check-jsonschema` | Validate representative assessment, gate, and release examples | `pass` | Example fixtures conform to schema |
| OPA/Rego | Express simplified adapter-capture pass/fail rule | `pass` | Policy evaluates and correctly detects `test_provenance_not_overclaimed` failure fixture |
| CUE | JSON Schema import | `partial` | CUE can import JSON Schema types |
| CUE | Direct validation against sdp-trace `schema/` refs | `cannot_verify` | Blocked: schema refs must be packaged as a CUE module first |
| in-toto | Wrap command, sign link metadata, record material/product hashes | `pass` | Link metadata generated and signed locally |
| Cosign | Sign/verify local `run.json` blob | `pass` | Works with local key when transparency log verification is explicitly disabled |
| Cosign | Verify with transparency log / Rekor | `fail` | Expected for local-only fixtures; no Rekor entry exists |
| SLSA verifier | Accept local DSSE fixture as production SLSA evidence | `fail` | Expected: no matching Rekor entries found |

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
