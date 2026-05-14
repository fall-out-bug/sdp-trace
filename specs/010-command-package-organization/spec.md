# Feature Specification: Command Package Organization

**Feature Branch**: `010-command-package-organization`
**Created**: 2026-05-13
**Status**: Draft for PI review
**Input**: PR #43 closed numeric quality gates partly through broad same-package file splitting. Review accepted the result but recorded navigation and contributor-DX debt.

## Product Boundary

This slice improves human maintainability and navigation of the CLI implementation without changing user-facing behavior.

Allowed:

- same-behavior Go refactors;
- command-family grouping;
- package or file organization changes;
- small tests that guard behavior while moving code.

Not allowed:

- changing CLI output or exit semantics unless separately specified;
- lowering quality gates;
- using file movement to hide unverified behavior;
- adding non-Go tooling.

## Problem Statement

The codebase now passes strict metrics, but the CLI package contains hundreds of small files. That is acceptable as a metric remediation endpoint, but it is not ideal for new contributors or review agents.

The next improvement should preserve the gates while making command-family ownership visible:

- reviewers should find packet, PR review, assess, witness, query, and recorder command code quickly;
- future edits should not require scanning hundreds of similarly named files;
- package boundaries should not introduce cycles or harness-specific dependencies.

## Core Claim

This slice may claim:

> The CLI implementation is organized by command family or another reviewed structure while preserving command behavior and quality gates.

This slice must not claim:

- a new CLI architecture;
- better metrics unless replayed;
- semantic behavior changes;
- full Clean Architecture completion beyond the reviewed organization scope.

## Required User Stories

### US-001 - Contributor Navigation (P0)

A contributor can identify the implementation area for a command family without scanning hundreds of files.

**Independent Test**: A short repository-local guide or package layout makes command-family ownership explicit, and PI DX review confirms navigation is improved or records remaining gaps.

### US-002 - Behavior Preservation (P0)

Existing CLI commands, help text, exit codes, and docs checks remain unchanged unless a change is explicitly reviewed.

**Independent Test**: Current CLI tests, `go run ./cmd/sdp-trace --help`, and `go run ./tools/doccheck` pass.

### US-003 - Quality Gate Preservation (P0)

CRAP, cognitive complexity, cyclomatic complexity, and strict MI remain green after any regrouping.

**Independent Test**: Replay `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd internal tools`, `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd internal tools`, complexity gates, and coverage-backed CRAP.

### US-004 - Dependency Direction (P0)

If subpackages are introduced, dependency direction stays clear and does not make core trace/evidence packages depend on CLI packages.

**Independent Test**: `go list ./...` and PI architecture review find no dependency cycle or boundary inversion.

## Functional Requirements

- **FR-001**: Choose and document one organization strategy before moving code: command-family subpackages, family-prefixed files, or generated index.
- **FR-002**: Move or index code in small slices so regressions are attributable.
- **FR-003**: Keep command handlers discoverable from a single registry.
- **FR-004**: Preserve existing command names, flags, outputs, and exit code semantics.
- **FR-005**: Record any remaining high-file-count package as an advisory follow-up, not as hidden debt.

## Acceptance Criteria

- No behavior-changing diff is introduced without a specific spec delta.
- Full local verification passes.
- PI review covers code/correctness, Clean Architecture, DX, and requirements-vs-implementation.
- Any advisory follow-ups are recorded separately from blockers.

## PI Review Prompt

Review this spec for whether it can reduce CLI navigation debt without undoing the quality-gate work. Focus on:

- whether the proposed organization is small enough for safe slices;
- whether it risks dependency cycles or import churn;
- whether behavior preservation is testable;
- whether the spec prevents metric-gaming from being mistaken for Clean Code.
