# Feature Specification: Repo Hygiene And Artifact Boundary

**Feature Branch**: `012-repo-hygiene-and-artifact-boundary`
**Created**: 2026-05-14
**Status**: Draft for PI review
**Input**: DX/UX/CTO review after PR #46 found root PR artifacts, local worktree leakage, and scattered review evidence paths.

## Problem Statement

The repository currently mixes durable product docs with short-lived PR artifacts and local workflow output. A new contributor can see `PR_DESCRIPTION.md`, `design-note.md`, root `reviews/`, and an untracked `.worktrees/` directory next to the product entrypoints. This weakens the repo boundary and makes it harder to tell which files are source of truth.

## Core Claim

This slice may claim:

> The repository separates durable source/docs/specs from local, PR, and subagent artifacts, and CI prevents common accidental artifact commits.

This slice must not claim:

- product readiness;
- stronger trust evidence;
- cleanup of historical spec records beyond the explicit migration scope.

## Required User Stories

### US-001 - Clean Repository Root (P0)

A new contributor can list the repository root and distinguish product entrypoints from implementation directories without PR-specific clutter.

**Independent Test**: `git ls-files` has no root `PR_DESCRIPTION.md`, root `design-note.md`, root `reviews/*`, root binaries, or checked-in worktree/subagent run output unless explicitly allowlisted.

### US-002 - Artifact Boundary Guardrails (P0)

Accidental binaries, local worktrees, behavior snapshots, and subagent run outputs are blocked or ignored before review.

**Independent Test**: A Go or shell-thin CI check fails on checked-in root executable artifacts, `.worktrees/`, `.codex-subagents/runs/`, `.sdp-trace-*`, and absolute `/home/...` paths in durable docs.

### US-003 - Review Evidence Home (P1)

Review artifacts have one discoverable durable location and do not compete with product docs.

**Independent Test**: `docs/README.md` or a new review-evidence policy points to the canonical location; root review artifacts are migrated or removed.

## Functional Requirements

- **FR-001**: Move PR #46 root artifacts into the owning spec directory or remove them if they duplicate PR body data.
- **FR-002**: Add guardrails for checked-in binaries, worktrees, subagent run directories, local behavior snapshots, and machine-specific absolute paths.
- **FR-003**: Define where durable review synthesis lives versus raw local subagent output.
- **FR-004**: Keep historical records intact unless a file is clearly a local or PR-scoped artifact.
- **FR-005**: Keep checks portable and Go-first; Bash is allowed only as a thin launcher.

## Acceptance Criteria

- Root tracked files are limited to durable repo entrypoints and product metadata.
- `git diff --check`, `go test -count=1 ./...`, and the new hygiene check pass.
- No root binary or local worktree output can be committed without a deliberate allowlist update.
- Docs explain where future PR review artifacts belong.

## PI Review Prompt

Review whether this spec removes local/PR clutter without erasing useful evidence. Focus on artifact classification, CI enforceability, and whether the proposed guardrails could block legitimate portable examples.
