# Spec Review Pack: 017 OSS Replacement Compatibility And Benchmarks

## Artifact
The following files from `specs/017-oss-replacement-compatibility-and-benchmarks/`:

### spec.md
```
# Spec 017: OSS Replacement Compatibility And Benchmarks

Status: draft

## Objective
Create a reproducible compatibility and performance harness for candidate OSS
replacements: JSON Schema tooling, CUE, OPA/Rego, in-toto, Sigstore/Cosign,
SLSA verifier, and minimal shell/Go prototypes.

This spec tests substitution boundaries and records the replacement decisions
for this phase. It does not approve replacing product code by implementation
accident.

## Evidence From 2026-05-20 Probe
- check-jsonschema validates checked flight-recorder fixtures against locally rewired schema refs.
- Current live sdp-trace wrap output does not validate against schema/flight-recorder-run.schema.json; required fields and timestamp format differ.
- check-jsonschema validates representative assessment, gate, and release example files.
- OPA can express a simplified adapter-capture pass/fail rule and detects the test_provenance_not_overclaimed failure fixture.
- CUE can import JSON Schema, but direct validation is blocked until schema refs are packaged as a CUE module.
- in-toto-run can wrap a command, sign link metadata, and record material and product hashes.
- cosign can sign and verify a local run.json blob with a local key when transparency-log verification is explicitly disabled.
- slsa-verifier does not accept the local DSSE fixture as production SLSA evidence; it fails because no matching Rekor entries are found.

## Benchmark Snapshot
Local 20-iteration benchmark, compiled sdp-trace, Linux amd64:
| Probe | Median ms | Notes |
| shell prototype wrap | 6.0 | Minimal JSON, no hash chain semantics |
| sdp-trace verify | 8.0 | Existing local run |
| OPA adapter policy eval | 14.0 | Simplified policy |
| sdp-trace wrap | 26.0 | Local /bin/true |
| Cosign local verify | 30.5 | Transparency log ignored |
| in-toto-run | 148.0 | Signed link metadata |
| check-jsonschema fixture validation | 271.5 | Python validator startup cost |

## Requirements
- FR-017-001: Add a reproducible OSS compatibility command or tool under the Go active path or a documented non-product sandbox path.
- FR-017-002: Cover at least these probes: JSON Schema, OPA/Rego, CUE import, in-toto command wrapping, Cosign blob signing, SLSA verifier negative path.
- FR-017-003: Record substitution boundaries: what can replace code, what needs adapter glue, and what remains sdp-trace-specific.
- FR-017-004: Add benchmark output with median, min, max, iterations, command, and environment.
- FR-017-005: Keep benchmarks non-authoritative: they inform scope decisions but do not prove production readiness.
- FR-017-006: Do not add Node.js, npm, JavaScript, TypeScript, or .mjs tooling to the product path.

## Decisions
- No existing product command is replaced by OSS tooling in this phase.
- The in-scope OSS candidates are JSON Schema validation, OPA/Rego, CUE, in-toto, Cosign, and SLSA verifier.
- Reproducible compatibility and benchmark tooling belongs under the Go active path (tools/osscompat and tools/ossbench) unless a later spec explicitly narrows a probe to docs-only evidence.
- JSON Schema can be used for external CI/example validation now. Live recorder output is not schema-compatible until the wrap output/schema drift is fixed.
- OPA/Rego may replace hand-written policy expressions only behind an explicit JSON translation layer. It does not replace verifier evidence collection.
- CUE remains import-only until schema refs are packaged as CUE modules.
- in-toto, Cosign, and SLSA tooling may provide signing/provenance substrates. They do not replace witness, release-proof, or gate verdict semantics without OIDC, Rekor, and identity-policy evidence.
- Benchmarks are scope evidence only. They cannot create health scores, production readiness, or replacement approval.

## Non-Goals
- No automatic migration to OSS replacements.
- No production Sigstore, Rekor, or SLSA trust claim from local fixtures.
- No benchmark health score.
- No long-running external service dependency in default local tests.

## Acceptance Criteria
- Compatibility probes can be rerun from a clean checkout with documented prerequisites.
- Probe results distinguish pass, fail, cannot_verify, and not_assessed.
- The live wrap output/schema drift is either fixed or documented as a blocker.
- SLSA/Rekor production trust remains not_assessed unless live external provenance is provided.
```

### plan.md
```
# Plan: OSS Replacement Compatibility And Benchmarks

## Workstreams
### WS-017-A: Reproducible Probe Harness
Owned files: tools/osscompat/*, docs/oss-replacement-compatibility.md
Deliverable: A command that runs or explains compatibility probes without mutating tracked product artifacts.

### WS-017-B: Schema Compatibility
Owned files: schema/*, examples/flight-recorder/*, docs/agent-entrypoint.md only if contract changes
Deliverable: Implement the spec decision that live recorder schema compatibility remains blocked until wrap output conforms or a separate current recorder schema is defined.

### WS-017-C: Policy-As-Code Prototype
Owned files: docs/oss-policy-prototype.md, optional examples/oss-policy/*
Deliverable: OPA/Rego or CUE prototype for one assessment profile, with clear boundary between policy and product verifier behavior.

### WS-017-D: Supply Chain OSS Prototype
Owned files: docs/oss-supply-chain-prototype.md, optional examples/oss-supply-chain/*
Deliverable: in-toto/Cosign/SLSA probe records that show what replaces witness, checkpoint, and release-proof pieces, and what does not.

### WS-017-E: Benchmark Harness
Owned files: tools/ossbench/*, docs/oss-benchmark-results.md
Deliverable: Repeatable local benchmarks with environment, command lines, iterations, and median/min/max output.

## Verification
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check

## Pi Handoff Notes
Assign WS-017-B only to a worker allowed to touch schema/examples. Assign WS-017-C and WS-017-D to separate workers because policy and supply-chain experiments have different dependencies.
```

### tasks.md
```
# Tasks: OSS Replacement Compatibility And Benchmarks

## Phase 0 - Review
- [ ] T017-001 Verify the spec-approved substitution candidates remain the only in-scope OSS tools for first implementation.
- [ ] T017-002 Implement compatibility tooling under the Go active path unless a later spec explicitly narrows a probe to docs-only evidence.

## Phase 1 - Pi-Ready Workstreams
- [ ] T017-010 WS-017-A: Create reproducible OSS compatibility harness. Pi ownership: tools/osscompat/*.
- [ ] T017-020 WS-017-B: Resolve live wrap output vs flight-recorder schema compatibility. Pi ownership: schema/examples plus focused recorder tests.
- [ ] T017-030 WS-017-C: Build OPA/Rego or CUE assessment-profile prototype. Pi ownership: docs/examples for policy prototype only.
- [ ] T017-040 WS-017-D: Build in-toto/Cosign/SLSA supply-chain prototype. Pi ownership: docs/examples for supply-chain prototype only.
- [ ] T017-050 WS-017-E: Add benchmark harness and benchmark report. Pi ownership: tools/ossbench/* and benchmark docs.

## Phase 2 - Integration
- [ ] T017-060 Run all compatibility probes available in the local environment.
- [ ] T017-070 Mark unavailable external services as not_assessed or cannot_verify; do not infer pass from local fixture success.
- [ ] T017-080 Update roadmap and docs index after accepted implementation.
- [ ] T017-090 Keep the live wrap output vs schema/flight-recorder-run.schema.json drift open as a blocker until T017-020 lands.
```

## Rules from AGENTS.md
- Target product code is Go. No Node.js, npm, JavaScript, TypeScript, or .mjs tooling in active product path.
- Bash allowed only as thin command launcher.
- New Go code must be small, readable, testable, covered by focused tests. No TODO/FIXME markers.
- Root router under 100 lines; module over 10 skills is too large.
- Machine proof wins over prose. No deferred trust closure.
- Source-bound proof requires clean immutable source commit.

## Review Instructions
1. Read the artifact (spec, plan, tasks) carefully.
2. Look for inconsistencies, overclaim, hidden coupling, trust boundary violations, and rule violations (especially Node.js/TS tooling sneaking in, module size limits, TODO markers).
3. Check whether requirements map cleanly to deliverables and whether acceptance criteria are verifiable.
4. Report findings with severity (Critical / Important / Advisory) and concrete evidence.
5. If zero findings, output exactly `LGTM`.
