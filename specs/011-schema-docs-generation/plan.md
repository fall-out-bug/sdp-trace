# Implementation Plan: Schema Documentation Validation

## Summary

Add a Go-first validation or generation path that keeps schema documentation synchronized with current schema files.

## Technical Context

**Language**: Go  
**Dependencies**: Go standard library, existing JSON syntax checks, optional `jq` in CI  
**Verification**: `go test -count=1 ./...`, `jq empty schema/*.json`, schema-doc checker, `git diff --check`  
**Constraints**: no schema semantic changes, no Node.js/npm tooling.

## Phase 0 - PI Spec Review

Run PI review across:

- docs completeness;
- DX for agents/downstream integrators;
- requirements-vs-implementation;
- trust/overclaim risk.

Stop for explicit approval before implementation.

## Phase 1 - Metadata Shape

Choose whether metadata lives in:

- a generated README section;
- a small checked JSON/YAML index;
- comments or annotations in schema files read by a Go tool.

Record why the chosen shape is simplest and least likely to drift.

## Phase 2 - Checker Or Renderer

Implement the smallest Go tool/test that can detect:

- schema file missing from docs/index;
- docs/index entry for missing schema;
- missing status/purpose;
- broken example reference unless marked `not_assessed`.

## Phase 3 - CI And Docs

- Wire the checker into CI or an existing docs check.
- Update docs to explain the difference between syntax validation, docs freshness, and semantic schema coverage.
- Run full verification and PI implementation review.

## Exit Criteria

- Schema docs/index cannot drift silently.
- Current schemas have status and purpose metadata.
- Example coverage gaps are explicit.
- No schema semantic claim is made without semantic tests.
