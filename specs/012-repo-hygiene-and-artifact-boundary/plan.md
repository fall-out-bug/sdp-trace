# Implementation Plan: Repo Hygiene And Artifact Boundary

## Technical Context

**Language**: Go for checks; Markdown for policy docs
**Dependencies**: Go standard library, existing CI workflow, existing `docs/README.md`
**Verification**: `go test -count=1 ./...`, hygiene check, `git diff --check`, `go run ./tools/doccheck`

## Scope

- Add a small repository hygiene checker or extend an existing Go tool.
- Move or remove PR-local root artifacts.
- Update ignore rules and documentation for local worktrees and subagent runs.

## Non-Goals

- Rewriting historical spec ledgers.
- Removing committed evidence packages that are explicitly part of older specs.
- Introducing Node, Python, or non-Go product tooling.

## Risks

- Over-broad checks could block valid examples.
- Moving review artifacts could break references in PR comments.
- Treating local cleanup as trust evidence would overclaim.

## Review Plan

Run DX and evidence review planes before implementation. Verify every moved file has either a new durable home or is intentionally removed as duplicate PR-only material.
