# Feature Specification: Spec Governance And Roadmap Navigation

**Feature Branch**: `015-spec-governance-and-roadmap`
**Created**: 2026-05-14
**Status**: Draft for PI review
**Input**: DX/CTO review found parallel specs, inconsistent spec lifecycle markers, and no concise ownership map for current product work.

## Problem Statement

The repo has many historical block records and active SpecKit-shaped directories. This is valuable evidence, but it is difficult to tell which spec owns which product surface, which specs are active, blocked, or complete, and where new work should attach. Task files also vary in lifecycle detail.

## Core Claim

This slice may claim:

> The repository has a lightweight roadmap and spec lifecycle policy that tells contributors which specs are active, blocked, complete, or historical.

This slice must not claim:

- that historical evidence was replayed;
- that blocked specs are resolved;
- that roadmap status is a trust verdict.

## Required User Stories

### US-001 - Spec Ownership Map (P0)

A reviewer can see which spec owns command surface, docs generation, repository hygiene, onboarding, docs UX, schema docs, and product contract work.

**Independent Test**: A one-page roadmap maps current specs to owned capabilities, status, blockers, and next step.

### US-002 - Lifecycle Labels (P0)

A contributor can tell whether a spec is `draft`, `pending_review`, `approved`, `in_progress`, `paused`, `blocked`, or `complete`.

**Independent Test**: New specs use a consistent status taxonomy in `spec.md` and task blockers are explicit.

### US-003 - Historical Evidence Boundary (P1)

Historical block ledgers remain discoverable but are not mistaken for current live proof.

**Independent Test**: Roadmap distinguishes active specs from historical block records and links to trust rules for checked-in evidence.

### US-004 - Claim Tag Adoption Plan (P1)

The repo has a path to enforce claim tags for future authoritative prose without breaking historical artifacts.

**Independent Test**: A plan names the Markdown scopes where claim tags are required and where historical files are exempt.

## Functional Requirements

- **FR-001**: Add a short roadmap/navigation artifact for current specs.
- **FR-002**: Define lifecycle status labels and blocker notation.
- **FR-003**: Update current draft specs to use the taxonomy, or record follow-ups.
- **FR-004**: Define where claim tags are required for new authoritative claims.
- **FR-005**: Avoid rewriting historical evidence packages unless required by a scoped follow-up.

## Acceptance Criteria

- Roadmap covers specs 001 through 015 at the capability/status level.
- Each new spec created after this slice has consistent status and tasks.
- Reviewers can identify active next work without reading every historical block.
- Any claim-tag enforcement plan is scoped to new or touched files unless separately approved.

## PI Review Prompt

Review whether this governance spec improves roadmap clarity without turning checked-in prose into false authority. Focus on lifecycle labels, blocker visibility, and claim-tag enforcement scope.
