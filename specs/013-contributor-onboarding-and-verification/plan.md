# Implementation Plan: Contributor Onboarding And Verification

## Technical Context

**Language**: Markdown, Go for doc drift checks if needed
**Dependencies**: Existing docs and `tools/doccheck`
**Verification**: `go run ./tools/doccheck`, `go test -count=1 ./...`, smoke command replay, `git diff --check`

## Scope

- Create or update one contributor quick-start path.
- Reduce duplicated smoke commands.
- Extend doccheck only where duplication remains intentional.

## Non-Goals

- Redesigning all docs.
- Changing CLI behavior.
- Claiming CI or external trust from local smoke results.

## Risks

- Over-shortening may hide trust caveats.
- Duplicated command snippets can drift again unless checked.

## Review Plan

Run DX and UX review planes. Require a cold-reader pass before closure.
