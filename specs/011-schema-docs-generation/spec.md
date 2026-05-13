# Feature Specification: Schema Documentation Validation

**Feature Branch**: `011-schema-docs-generation`  
**Created**: 2026-05-13  
**Status**: Draft for PI review  
**Input**: PR #43 DX review found that schema documentation is manually maintained and can drift from `schema/*.json`.

## Product Boundary

This slice prevents schema documentation drift. It does not redesign schemas or change their semantics.

Allowed:

- Go validation or rendering tools;
- checked generated Markdown if reviewed;
- CI checks that compare schema docs with schema files.

Not allowed:

- Node.js/npm schema tooling;
- changing schema behavior without a separate spec;
- claiming semantic correctness from documentation freshness alone.

## Problem Statement

The repository has many JSON schemas and a human-readable schema README. JSON syntax is checked, but docs can drift when schemas are added, renamed, retired, or changed.

Agents and downstream integrators need a reliable index that answers:

- which schemas are current;
- what each schema is for;
- which examples or fixtures exercise it;
- which schemas are historical or not assessed.

## Core Claim

This slice may claim:

> Schema documentation is checked or generated from current schema files so missing, stale, or extra schema entries are detected.

This slice must not claim:

- every schema has complete semantic tests;
- every example validates against every schema;
- external trust or compatibility certification.

## Required User Stories

### US-001 - Schema Index Freshness (P0)

A contributor adding, renaming, or removing a schema must update docs or a generated index, and CI catches drift.

**Independent Test**: A focused test or Go tool fails when a schema file is missing from the index or the index names a missing schema.

### US-002 - Schema Purpose Metadata (P0)

Each current schema has enough metadata for an agent to understand its purpose and whether it is active, historical, or not assessed.

**Independent Test**: The checker rejects entries without purpose/status metadata.

### US-003 - Example Coverage Mapping (P1)

Where examples exist, the index names representative examples or marks example coverage `not_assessed`.

**Independent Test**: The checker can distinguish present example refs from `not_assessed` refs and fails broken file refs.

### US-004 - No Product Tooling Regression (P0)

The validation/generation path remains Go-first and does not add Node tooling.

**Independent Test**: CI uses Go and existing shell/JQ checks only.

## Functional Requirements

- **FR-001**: Add a Go checker or renderer for schema documentation.
- **FR-002**: The checker must detect missing schema entries, extra schema entries, and broken example refs.
- **FR-003**: Every schema entry must include status and purpose.
- **FR-004**: If generated docs are committed, generation must be deterministic.
- **FR-005**: Docs must distinguish syntax validation from semantic validation.

## Acceptance Criteria

- Existing `jq empty schema/*.json` continues to pass.
- New schema-doc validation passes locally and in CI.
- `go test -count=1 ./...` passes.
- PI review covers docs completeness, DX, and requirements-vs-implementation.

## PI Review Prompt

Review this spec for whether it prevents schema documentation drift without expanding into schema redesign. Focus on:

- whether status/purpose/example metadata is enough for agents;
- whether the checker scope is precise and Go-only;
- whether the spec avoids treating docs freshness as semantic schema proof.
