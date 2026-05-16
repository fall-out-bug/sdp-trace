# Block Intake: 015-spec-governance-and-roadmap

## Objective
- User request / block: Take `specs/015-spec-governance-and-roadmap` into work, work in separate worktree, drive to green PR with multi-LLM review.
- SpecKit delta needed: Define spec lifecycle taxonomy, create lightweight roadmap artifact, set task-file expectations, scope claim-tag enforcement.
- Existing work to land or park first: Working tree is clean (`main` at `171a1e2`). No active block to park.

## Scope
- In scope:
  - Spec lifecycle taxonomy (draft, pending_review, approved, in_progress, paused, blocked, complete).
  - Lightweight roadmap/navigation artifact covering specs 001–015.
  - Task-file expectations for blockers and approval gates.
  - Claim-tag enforcement scope for future authoritative prose.
- Out of scope:
  - Replaying historical evidence.
  - Closing old spec gaps or migrating historical blocks.
  - Rewriting historical evidence packages.
- Manifest subjects:
  - `specs/015-spec-governance-and-roadmap/spec.md`
  - `specs/015-spec-governance-and-roadmap/plan.md`
  - `specs/015-spec-governance-and-roadmap/tasks.md`
  - New roadmap artifact (to be named).

## Trace map
- Spec: `specs/015-spec-governance-and-roadmap/spec.md`
- Plan: `specs/015-spec-governance-and-roadmap/plan.md`
- Task: `specs/015-spec-governance-and-roadmap/tasks.md`
- Evidence: Review dispositions, doccheck, git diff --check, roadmap consistency review.
- Gate: Socratic spec review passed; implementation approval received.
- Decision: Roadmap filename, location, and claim-tag scope to be approved.
- Provenance: Worktree `.worktrees/015-spec-governance-and-roadmap`, branch `015-spec-governance-and-roadmap`.

## Review before implementation
- Socratic/pi-review plane: Adversarial spec review (claim-doubt-cycle) conducted against AGENTS.md trust rules and `docs/claim-authoring.md`.
- Findings: See `specs/015-spec-governance-and-roadmap/review-disposition.md`.
- Disposition: Findings recorded; fixes or explicit `not_assessed` markers required before implementation.
- User approval status: **Pending** — awaiting explicit approval of reviewed spec direction.

## Initial verifier state
- pass: Worktree created; branch checked out; no uncommitted changes on main.
- fail: (none)
- cannot_verify: (none)
- not_assessed: Multi-LLM external review plane (GLM/Qwen/DeepSeek) — harness availability not confirmed yet.
