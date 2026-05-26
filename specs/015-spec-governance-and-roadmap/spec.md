# Feature Specification: Spec Governance And Roadmap Navigation

**Feature Branch**: `015-spec-governance-and-roadmap`
**Created**: 2026-05-14
**Status**: `in_progress` — implementation active on branch `015-spec-governance-and-roadmap`
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

**Verification**: This claim is satisfied when `docs/roadmap.md` exists, lists every spec 001–015 with capability and status, and `docs/claim-authoring.md` governs authoritative claims. Roadmap status is prose guidance, not machine proof; authoritative claims require `sdp-trace-claim` tags per `docs/claim-authoring.md`.

## Required User Stories

### US-001 - Spec Ownership Map (P0)

A reviewer can see which spec owns command surface, docs generation, repository hygiene, onboarding, docs UX, schema docs, and product contract work.

**Independent Test**: A one-page roadmap maps current specs to owned capabilities, status, blockers, and next step.

### US-002 - Lifecycle Labels (P0)

A contributor can tell whether a spec is `draft`, `pending_review`, `approved`, `in_progress`, `paused`, `blocked`, `complete`, or `historical`.

**Independent Test**: New specs use a consistent status taxonomy in `spec.md` and task blockers are explicit.

**Caveat**: A status of `complete` or `approved` describes editorial lifecycle state, not a trust verdict. Trust state (pass, fail, `not_assessed`) is expressed only through `sdp-trace-claim` tags or live verifier output. See `docs/claim-authoring.md`.

### US-003 - Historical Evidence Boundary (P1)

Historical block ledgers remain discoverable but are not mistaken for current live proof.

**Independent Test**: Roadmap distinguishes active specs from historical block records and links to trust rules for checked-in evidence.

**Concrete Rule**: Historical block records live in a spec's `blocks/` directory and are marked `historical` in the roadmap with no live status transition. They are exempt from claim-tag requirements. Live work uses the spec root (`spec.md`, `plan.md`, `tasks.md`).

### US-004 - Claim Tag Adoption Plan (P1)

The repo has a path to enforce claim tags for future authoritative prose without breaking historical artifacts.

**Independent Test**: A plan names the Markdown scopes where claim tags are required and where historical files are exempt.

**Scope Rule**: Claim tags are required for new authoritative claims in any file created or materially modified after this slice. Historical specs (001–014), their `blocks/` directories, and existing checked-in review JSON are exempt unless separately approved for migration.

### US-005 - Multi-Axis Status Discipline (P0)

A contributor can distinguish spec review, task completion, implementation,
review closure, merge state, and trust state without treating one field as a
complete readiness verdict.

**Independent Test**: Roadmap guidance names separate status axes and lists
implemented-but-not-formally-closed specs separately from formal completion.

## Functional Requirements

- **FR-001**: Add a short roadmap/navigation artifact for current specs at `docs/roadmap.md`.
- **FR-002**: Define lifecycle status labels and blocker notation.
- **FR-003**: Update current draft specs to use the taxonomy, or record follow-ups.
- **FR-004**: Define where claim tags are required for new authoritative claims.
- **FR-005**: Avoid rewriting historical evidence packages unless required by a scoped follow-up.
- **FR-006**: Roadmap freshness: update `docs/roadmap.md` when a new spec is opened or an active spec's status changes. Owner = the spec author or current block worker.
- **FR-007**: Roadmap status must not collapse implementation, review, merge,
  and trust evidence into one lifecycle label.

## Acceptance Criteria

- Roadmap covers specs 001 through 015 at the capability/status level.
- Each new spec created after this slice has consistent status and tasks.
- Reviewers can identify active next work without reading every historical block.
- Reviewers can identify specs that are implemented but not formally closed.
- Any claim-tag enforcement plan is scoped to new or touched files unless separately approved.

## PI Review Prompt

Review whether this governance spec improves roadmap clarity without turning checked-in prose into false authority. Focus on lifecycle labels, blocker visibility, and claim-tag enforcement scope.
