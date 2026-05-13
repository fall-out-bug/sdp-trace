# Implementation Plan: Command Package Organization

## Summary

Improve CLI maintainability by organizing command implementation around command families while preserving behavior and quality gates.

## Technical Context

**Language**: Go
**Dependencies**: Go standard library and current internal packages
**Verification**: full tests, doccheck, strict MI `70.1`, complexity gates, coverage-backed CRAP, `go vet ./...`, `git diff --check`
**Constraints**: no CLI behavior changes, no non-Go tooling, no dependency inversion.

## Phase 0 - Organization Design Review

Before moving code, write a short design note choosing one strategy:

- command-family subpackages;
- command-family file prefixes with generated index;
- hybrid grouping for only the largest families.

Run PI review on the design before implementation.

## Phase 1 - One Vertical Slice

Move or index one command family with strong tests, preferably a smaller family first. Verify:

- command still dispatches;
- help/docs unchanged;
- quality gates still pass.

## Phase 2 - Repeat By Family

Proceed family by family. Do not batch unrelated command families in one commit.

Candidate families:

- packet;
- PR review;
- assess;
- witness;
- observe/harness;
- release proof;
- query/report/gate.

## Phase 3 - Final Audit

- Run full verification.
- Run PI review planes.
- Update contributor-facing docs only if the organization changed how future contributors should navigate the code.

## Exit Criteria

- Reviewers can identify command-family ownership quickly.
- Quality gates remain green.
- No command behavior changes are introduced accidentally.
- Remaining navigation debt is recorded explicitly.
