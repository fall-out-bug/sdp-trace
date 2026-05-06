# Implementation Plan: sdp-trace Time-Series Evidence Substrate

**Branch**: `001-sdp-trace-time-series-evidence-substrate` | **Date**: 2026-04-30 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-sdp-trace-time-series-evidence-substrate/spec.md`

## Summary

Reframe `sdp-trace` as the portable evidence, provenance, observation, accountability, contract-integrity, and time-series trace substrate for AI-assisted delivery. The implementation creates SpecKit-visible artifacts first, then evolves schemas, examples, docs, and pilot run-cards so `sdp-gate` can later apply policy without `sdp-trace` owning gate decisions.

The roadmap now treats self-proof as the first product proof, not as a later nicety. Block 01 may establish contract scaffolding, but it cannot establish product viability until Block 02 self-trace and Block 03 self-attestation are validated.

## Technical Context

**Language/Version**: JSON Schema Draft 2020-12, Markdown, shell validation commands
**Primary Dependencies**: Go; `jq` for JSON parse checks; SHA-256 digest verification; target signing profile using in-toto Statement, DSSE, and Sigstore/Cosign keyless verification where available
**Storage**: Files committed to the repository; local ignored pilot outputs under `.sdp-trace-runs/`
**Testing**: `jq empty schema/*.json`; Draft 2020-12 validation over committed examples with pinned local `ajv@8.20.0`
**Target Platform**: Portable repository artifacts
**Project Type**: Schema and documentation substrate
**Performance Goals**: Validation commands should run locally in seconds over committed artifacts
**Constraints**: No dependency on `sdp_lab`, Beads, Operator Mode, agentloop, OpenCode, GitHub, Claude, Codex, or any specific harness runtime
**Scale/Scope**: Pilot covers OpenCode model slices, selected harness families, and JVM/Kotlin/Bazel evidence paths

## Constitution Check

The repository rules act as the constitution for this feature.

| Rule | Status | Evidence |
|---|---|---|
| Use SpecKit terms first | Pass | Spec, plan, task, evidence, gate, decision, trace, provenance are used throughout this package. |
| Keep `sdp-trace` independent | Pass | Plan excludes runtime dependencies on `sdp_lab`, Beads, Operator Mode, agentloop, and harness runtimes. |
| Evidence-backed claims only | Conditional | Contract scaffolding is validated, but product viability remains blocked until committed self-trace and self-attestation artifacts validate. |
| No opaque health scores | Pass | Metric catalog explicitly excludes aggregate health scores. |
| Do not implement without plan/spec | Pass | This SpecKit package is the implementation gate. |

## Project Structure

### Documentation (this feature)

```text
specs/001-sdp-trace-time-series-evidence-substrate/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── socratic-resolution-notes.md
├── socratic-judge-result.json
├── contracts/
│   └── sdp-trace-sdp-gate-boundary.md
└── tasks.md
```

### Repository Artifacts Expected Later

```text
schema/
├── accountability.schema.json
├── risk-classification.schema.json
├── contract-manifest.schema.json
├── contract-release-verification.schema.json
├── evidence-event.schema.json
├── observation.schema.json
├── metric-stream.schema.json
└── assessment-input.schema.json

examples/
├── opencode/
├── superpowers/
├── github-speckit/
├── self-trace/
└── jvm-bazel/

docs/
├── accountability-model.md
├── contract-release-signing.md
├── concepts.md
├── evidence-policy.md
├── harness-compatibility-matrix.md
├── model-compatibility.md
└── jvm-bazel-guide.md
```

**Structure Decision**: Keep feature planning under `specs/`; keep product docs under `docs/`; keep portable schemas under `schema/`; keep concrete examples under `examples/`; keep ignored raw run outputs under `.sdp-trace-runs/`.

## Contract Decisions From Socratic Review

This SpecKit package was challenged by clean-context `pi` critics using GLM, MiniMax, and Kimi. The converged blockers are resolved by these plan decisions:

| Blocker | Plan decision |
|---|---|
| CTO asks "are we degrading?" while `sdp-trace` must not issue verdicts | `sdp-trace` records prior/current values, deltas, dimensions, evidence coverage, and `not_assessed` gaps. A yes/no degradation verdict is external. |
| External verdicts and evidence strength can blur into native trace judgments | Native trace entities record observations and samples only. External verdicts, scores, and strength assertions use an explicit external verdict input shape with producer and origin. |
| JSON Schema draft and validator were selected too late | New schemas target Draft 2020-12. Validation strategy is pinned before schema authoring: local `ajv@8.20.0` plus `jq` syntax checks. |
| Evidence references can leak sensitive or unverifiable artifacts | Committed examples require sanitized summaries, SHA-256 hashes when artifacts are committed, redaction notes, and `integrity_status`. Raw local outputs stay ignored under `.sdp-trace-runs/`. |
| `sdp-gate` inherits contracts without versioning policy | Every new schema gets a semver `schema_version`; additive optional changes are minor, required/semantic changes are major. `schema/trace.schema.json` remains a compatibility path until a replacement and migration note exist. |
| CEO asks who is accountable when AI-assisted work fails | Accountable artifacts carry human-held DRI, approver, risk owner, escalation path, approval reference, and line of defense. AI actors can produce or review but cannot be sole accountable owners. |
| Humans or models can simplify the contract and still pass JSON validation | Trusted contract release requires a manifest with SHA-256 digests, a trusted identity policy, and release verification under `sdp-trace-signature/sigstore-dsse-keyless-v1` or an explicit private equivalent. JSON-valid but unsigned, unauthorized-signer, stale, or digest-mismatched artifacts are not trusted releases. |

## Phase 0: Research

Research is captured in [research.md](research.md). It maps useful `sdp_lab` ideas into portable `sdp-trace` terms and rejects runtime or policy coupling.

## Phase 1: Design

Design outputs:

- [data-model.md](data-model.md) defines the entities and relationships.
- [contracts/sdp-trace-sdp-gate-boundary.md](contracts/sdp-trace-sdp-gate-boundary.md) defines product ownership boundaries.
- [quickstart.md](quickstart.md) defines how a repository observer validates the artifacts.

## Phase 2: Task Breakdown

Task breakdown lives in [tasks.md](tasks.md). Beads issues mirror these tasks for execution discipline only:

| SpecKit task area | Beads mirror |
|---|---|
| Extract portable ideas from `sdp_lab` | `sdp-trace-cdn.1` |
| Define `sdp-trace` / `sdp-gate` boundary | `sdp-trace-cdn.2` |
| Define human accountability and oversight classification | `sdp-trace-cdn.2` |
| Design time-series observation contract | `sdp-trace-cdn.3` |
| Design evidence event and provenance contract | `sdp-trace-cdn.4` |
| Define process metric catalog | `sdp-trace-cdn.5` |
| Create harness and model run-cards | `sdp-trace-cdn.6` |
| Add JVM/Kotlin/Bazel fixture plan | `sdp-trace-cdn.7` |
| Add schema validation and fixture strategy | `sdp-trace-cdn.8` |
| Add contract manifest and release verification strategy | `sdp-trace-cdn.8` |
| Self-trace this feature with the new contracts | `sdp-trace-cdn.12` |
| Self-attest the contract release proof | `sdp-trace-cdn.13` |
| Build customer pilot evidence package outline | `sdp-trace-cdn.9` |
| Update compatibility matrices from evidence | `sdp-trace-cdn.10` |
| Reframe CTO and team docs | `sdp-trace-cdn.11` |

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| Multiple schema artifacts | Separate evidence events, observations, metric streams, and assessment inputs have different consumers and validation needs. | One large schema would make `sdp-gate` inheritance less clear and would force unrelated fields into every artifact. |
| Pilot matrix across harness/model/stack | Customer pilot explicitly asks for OpenCode, MiniMax/Kimi/GLM, Superpowers/GSD-style harnesses, and JVM/Kotlin/Bazel. | A single Codex or JVM baseline would not prove portability. |

## Self-Trace Milestone

`sdp-trace` starts tracing itself immediately after the minimal contract set exists and before any customer pilot or compatibility claim.

This milestone is a hard product gate for the repository narrative. If it fails, the product is not ready to ask a customer to trust the substrate.

Minimum prerequisite tasks:

- boundary and public language aligned: T005-T007
- accountability and risk classification schemas drafted: T044-T045
- contract manifest and release verification schemas drafted: T046-T047
- trusted identity policy and signing proof tasks complete: T049-T050
- observation and metric stream schemas drafted: T008-T009
- evidence, provenance, and assessment input schemas drafted: T015-T017
- at least one valid example and one `not_assessed` example exist
- schema draft, validator command, and schema versioning policy documented: T034
- artifact safety and integrity rules documented: T040
- negative fixtures exist for AI-as-sole-accountable-owner, modified manifest digest mismatch, and unauthorized signer identity: T048-T049

Self-trace output:

```text
examples/self-trace/
├── evidence-events.json
├── provenance-records.json
├── observations.json
├── metric-stream.json
├── trace-snapshot.json
└── assessment-input.json
```

The first self-trace is not a gate decision. It is a consumer test proving the contracts can describe the development of this SpecKit feature itself.

## Self-Attestation Milestone

After self-trace validates, the next mandatory milestone is self-attestation of the contract release.

Minimum proof states:

- `schema_valid`: committed JSON artifacts validate against committed schemas
- `digest_verified`: the contract manifest matches the checked-out artifacts
- `locally_attested`: DSSE/private-equivalent verification proves envelope binding, manifest digest verification, trusted identity policy match, freshness, and source-content verification against the selected source reference
- `externally_attested`: an external trust anchor, such as GitHub OIDC plus Sigstore/Rekor or accepted customer PKI, verifies signer identity and audit evidence
- `production_release_verified`: immutable source commit, protected ref or release tag, workflow identity, trusted identity policy, freshness, and rollback checks all match

`locally_attested` is useful engineering evidence only when all local source, digest, signature, identity-policy, and freshness checks are assessed. It is not a substitute for `externally_attested` or `production_release_verified` when making a product trust claim.

## Verification Commands

Syntax verification:

```bash
jq empty schema/*.json
```

Draft 2020-12 validation strategy (Block 10 compatible):

```bash
go test ./...
go run ./cmd/sdp-trace validate-fixtures examples/github-speckit
```

T036 will generalize this into a repository command that validates committed examples while excluding `.git/`, `.beads/`, `.sdp-trace-runs/`, and `benchmarks/repos/`.
