# Implementation Plan: Spec Governance And Roadmap Navigation

## Technical Context

**Language**: Markdown; optional Go check for status headings later
**Dependencies**: Existing specs, docs map, claim authoring docs
**Verification**: doccheck, manual roadmap review, `git diff --check`

## Scope

- Add a lightweight roadmap/navigation artifact.
- Define spec lifecycle taxonomy.
- Add task-file expectations for blockers and approval gates.

## Non-Goals

- Replaying historical evidence.
- Closing old spec gaps.
- Migrating every historical block to the new taxonomy.

## Risks

- Roadmap can become stale unless ownership is clear.
- Lifecycle labels can look like trust verdicts unless caveated.

## Review Plan

Run CTO/product, DX, and evidence review planes. Verify labels against current files before marking any status.
