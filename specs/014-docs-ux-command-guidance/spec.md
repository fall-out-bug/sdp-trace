# Feature Specification: Docs UX And Command Guidance

**Feature Branch**: `014-docs-ux-command-guidance`
**Created**: 2026-05-14
**Status**: Draft for PI review
**Input**: UX review found command discoverability, state vocabulary, output-location, profile-selection, and overclaim guidance gaps.

## Problem Statement

The CLI and docs are accurate but high-friction. Users see many commands, many evidence states, several scope vocabularies, and multiple output locations. Critical overclaim warnings are repeated across documents rather than presented as one decision aid.

## Core Claim

This slice may claim:

> The user-facing docs provide a guided path for choosing commands, interpreting evidence states, locating outputs, selecting assessment profiles, and avoiding overclaim.

This slice must not claim:

- new verifier behavior unless implemented separately;
- interactive CLI support unless a command is actually added;
- production trust or authority decisions.

## Required User Stories

### US-001 - Command Choice By Task (P0)

A user can choose the next command from a task-oriented guide rather than reading a long flat reference.

**Independent Test**: Reviewer docs answer "I have a run directory, what now?" and "Which assessment profile applies?" without requiring the full command table.

### US-002 - Evidence State Decision Tree (P0)

A reviewer can distinguish result states (`observed`, `pass`, `fail`, `not_assessed`, `cannot_verify`) from telemetry labels, completeness markers, and command-specific sub-states without guessing from exit codes.

**Independent Test**: One canonical state contract defines result states with exit codes and classifies all other state-like tokens. No doc outside that contract redefines or invents result states. `missing_telemetry` is classified as a telemetry label, not a result state. `warn`, `coverage_satisfied`, `coverage_partial`, `coverage_unresolved`, `not_integrated`, and `unsupported` are classified in the contract or removed from non-contract docs.

### US-003 - Output Location Map (P1)

A user can tell where each command writes artifacts and what each artifact is for.

**Independent Test**: Docs include an output map for run dirs, reports, query packs, witness outputs, assessment outputs, and release proof outputs.

### US-004 - Overclaim Checklist (P0)

A reviewer has one canonical checklist for forbidden interpretations and trust-scope escalation.

**Independent Test**: Reviewer entrypoint contains the canonical checklist; README and agent entrypoint link to it.

## Functional Requirements

- **FR-001**: Add a task-oriented command guide or restructure existing docs to include one.
- **FR-002**: Add canonical evidence state vocabulary and exit-code mapping. The canonical contract must classify every state-like token used in the product as either a result state, a telemetry label, a completeness marker, or a command-specific sub-state.
- **FR-003**: Add "which profile do I use?" decision tree. The tree must map trust profile IDs, assessment profiles, witness kinds, and authority scopes in a single decision aid.
- **FR-004**: Add output location reference.
- **FR-005**: Consolidate overclaim rules into one canonical checklist file; all other docs link to it and may retain one-line inline summaries.
- **FR-006**: Mark any future interactive guide command as a separate implementation follow-up, not implied by docs-only work.
- **FR-007**: Add profile taxonomy map that explicitly relates trust profile IDs (`repo_baseline_structural`, `source_bound_local_release`, `external_production_trust`) to assessment profiles, witness kinds, and authority scopes.

## Acceptance Criteria

- `docs/reviewer-entrypoint.md` has a short task path (≤10 lines or a table) before any flat command list longer than 5 entries.
- Exactly one canonical state contract exists (in `docs/agent-entrypoint.md` or a dedicated canonical doc) and defines: (a) result states with exit codes, (b) classification of all other state-like tokens. A grep audit shows zero orphan result-state tokens in non-contract docs.
- One canonical overclaim checklist exists in `docs/overclaim-checklist.md`; README and agent entrypoint link to it.
- One profile decision tree exists and maps all trust profile IDs, assessment profiles, witness kinds, and authority scopes.
- `go run ./tools/doccheck` passes and covers the command claims it owns.
- UX review finds no blocker-level ambiguity in the first-run reviewer path.

## PI Review Prompt

Review whether the proposed docs UX makes command choice and evidence interpretation safer for a cold user. Focus on misleading state language, output confusion, profile selection, and overclaim prevention.
