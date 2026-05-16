# Implementation Plan: Spec Governance And Roadmap Navigation

## Technical Context

**Language**: Markdown; optional Go check for status headings later
**Dependencies**: Existing specs, docs map, claim authoring docs
**Verification**: doccheck, manual roadmap review, `git diff --check`

## Scope

- Add a lightweight roadmap/navigation artifact at `docs/roadmap.md`.
- Define spec lifecycle taxonomy and status caveat (status ≠ trust verdict).
- Add task-file expectations for blockers and approval gates.
- Define claim-tag enforcement scope for new/touched files only.

## Non-Goals

- Replaying historical evidence.
- Closing old spec gaps.
- Migrating every historical block to the new taxonomy.

## Risks

- Roadmap can become stale unless ownership is clear.
  - **Mitigation**: FR-06 freshness rule — update on spec open/status change; owner = spec author.
- Lifecycle labels can look like trust verdicts unless caveated.
  - **Mitigation**: Explicit caveat in spec.md US-002 linking to `docs/claim-authoring.md`.

## Review Plan

- Socratic spec review completed; findings recorded in `review-disposition.md`.
- Implementation authorized after user approval.
- External multi-LLM review (GLM/Qwen/DeepSeek) deferred to PR stage.
