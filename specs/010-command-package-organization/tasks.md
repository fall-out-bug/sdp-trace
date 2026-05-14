# Tasks: Command Package Organization

**Input**: `spec.md`, `plan.md`
**Tests**: No implementation before PI design/spec review and explicit approval. Later implementation must preserve behavior and rerun full tests, doccheck, strict MI `70.1`, complexity, CRAP, vet, and `git diff --check`.

## Phase 0 - PI Review

- [x] T001 Run PI review on whether this organization work is worth doing now.
  - Models: `zai/glm-5.1`, `minimax/MiniMax-M2.7`.
  - Artifacts: `reviews/pi-review-glm-5.1.md`, `reviews/pi-review-minimax-2.7.md`, `reviews/synthesis.md`.
  - Key findings: flat-package family prefixes are safest; ~72 tiny metric-split files need merging first; snapshot behavior lock required.
- [x] T002 Choose one organization strategy and record the rationale.
  - Strategy: family-prefixed files within existing `package main`, plus generated index.
  - Rationale: zero import-cycle risk, zero registration indirection, zero test breakage, preserves all quality gates.
  - Recorded in: `design-note.md`.
- [x] T003 Review the chosen strategy for dependency-cycle and behavior-preservation risk.
  - Dependency cycle: not possible — no new packages introduced.
  - Behavior preservation: same-package renames only; snapshot lock via `command-surface` JSON and `--help` output before/after each slice.
  - Verified MiniMax false-positive claims about missing spec and undefined MI against `docs/agent-entrypoint.md`, `.github/workflows/ci.yml`, `tools/qualitycheck/`.
- [ ] T004 Stop for explicit approval before moving code.

## Phase 1 - First Slice

- [ ] T010 Pick one small command family.
- [ ] T011 Move or index that family only.
- [ ] T012 Prove help, docs, exit codes, and tests remain unchanged.

## Phase 2 - Iteration

- [ ] T020 Repeat by command family with scoped commits.
- [ ] T021 Keep command handlers discoverable from one registry.
- [ ] T022 Record any remaining high-file-count area as advisory debt.

## Phase 3 - Closure

- [ ] T030 Run full local verification.
- [ ] T031 Run PI code/correctness, Clean Architecture, DX, and requirements review.
- [ ] T032 Update contributor navigation docs only if the navigation contract changed.
