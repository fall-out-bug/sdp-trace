# Implementation Plan: Docs UX And Command Guidance

## Technical Context

**Language**: Markdown; optional small Go doccheck extension
**Dependencies**: Existing docs and command-surface output
**Verification**: doccheck, docs grep for deprecated state terms, cold-reader review, `git diff --check`

## Scope

- Reviewer and operator docs only unless a small doccheck extension is required.
- No CLI behavior changes in the first slice.

## Non-Goals

- Adding an interactive guide command.
- Rewriting all docs.
- Changing exit-code semantics.

## Risks

- Consolidation can hide important caveats if links replace necessary local warnings.
- State vocabulary must reflect actual CLI behavior, not desired behavior.

## Review Plan

Run UX, evidence, and requirements review planes. Verify every state and profile name against current command-surface/docs.
