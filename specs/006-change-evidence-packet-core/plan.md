# Implementation Plan: Change Evidence Packet Core

**Branch**: `006-change-evidence-packet-core` | **Date**: 2026-05-10 |
**Spec**: [spec.md](spec.md)
**Input**: Product Contract v0, 005 product contract direction, and demo 007
dependency.

## Summary

Implement the smallest Go-first `sdp-trace` product slice that can honestly
produce a CTO-readable Change Evidence Packet v0 from structured evidence.

This is the prerequisite for the GitHub OSS demo. Without this slice, the demo
can only show hand-authored Markdown and does not prove product capability.

## Technical Context

**Language/Version**: Go for product code; JSON Schema for contracts; Markdown
for canonical packet projection.
**Primary Dependencies**: Existing schema conventions, Go CLI patterns,
contract/product docs, harness observation structures where reusable.
**Storage**: `schema/`, `internal/packet/`, `internal/packetrender/`,
`examples/change-evidence-packet/`.
**Testing**: `go test ./...`, `jq empty schema/*.json`, fixture validation,
renderer golden tests, negative validator tests, `git diff --check`.
**Project Type**: Product trust artifact.
**Constraints**: No Node/npm/JS/TS. No live GitHub dependency required in P0.

## Constitution Check

| rule | status | evidence |
| --- | --- | --- |
| Spec before implementation | Pass | This slice defines contract before code. |
| Go-first product path | Pass | Product implementation stays in Go. |
| Evidence-backed claims only | Pass | Packet rows require refs or explicit missing state. |
| Preserve missing states | Pass | `cannot_verify` and `not_assessed` are first-class. |
| No approval system | Pass | `PC-DECISION` names owner, not approval. |
| Product independence | Pass | GitHub is an evidence surface, not a required runtime. |

## Product Workstreams

### A1: Schemas

Add JSON schemas for:

- `change-evidence-packet.v0`;
- `evidence-bundle-manifest.v0`;
- GitHub PR evidence input fixture shape.

Schemas must encode required rows, allowed states, evidence refs, resolver
entries, reasons for missing states, and artifact expiry fields.

### A2: Go Models

Add small Go models matching the schemas.

The model layer must keep row state explicit. Do not represent missing evidence
as booleans or derived health scores.

### A3: Validator

Add validator behavior for:

- required row completeness;
- allowed state transitions;
- evidence ref resolver presence;
- expired artifact handling;
- contradictory evidence handling;
- canonical packet artifact vs PR projection.

### A4: Markdown Renderer

Render canonical CTO-readable Markdown:

- metadata;
- row table;
- residual gaps;
- theater section;
- evidence bundle pointer;
- non-approval disclaimer.

Renderer output should be stable enough for golden tests.

### A5: CLI Surface

Add the smallest command surface needed for demo 007.

Preferred shape:

```text
sdp-trace packet validate --bundle <path>
sdp-trace packet render --bundle <path> --out <path>
```

If existing CLI conventions require a different shape, keep the behavior but
record the naming decision in the review ledger.

### A6: Fixtures And Examples

Add fixtures:

- happy path with clean `PC-THEATER: pass`;
- missing verification evidence;
- expired artifact;
- contradictory evidence;
- `agent_claimed_verification` theater finding;
- hand-authored pre-tooling metadata fixture if still needed for migration.

### A7: Trace/Docs

Update product docs only where needed to describe:

- packet contract;
- bundle contract;
- canonical vs projection distinction;
- what the packet does not prove.

Avoid broad enterprise or OSS ecosystem claims.

## Sequencing

1. Land reviewed 006 spec direction.
2. Implement schemas and fixtures.
3. Implement Go models and validator.
4. Implement renderer.
5. Add CLI surface.
6. Run Socratic/implementation review planes.
7. Only then start 007 demo repo packetization.

## Acceptance Evidence

- `go test ./...`
- `jq empty schema/*.json`
- validator positive and negative fixture tests
- renderer golden tests
- `git diff --check`
- review ledger with code/correctness, tracing/evidence, and
  requirements-vs-implementation planes

## Explicit Non-Goals

- GitHub API completeness.
- Signed artifacts.
- Enterprise policy engine.
- Automatic evidence storage.
- Visual dashboard.
- Demo repo mutation.
