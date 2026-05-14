# Tasks: Command Package Organization

**Input**: `spec.md`, `plan.md`
**Tests**: No implementation before PI design/spec review and explicit approval. Later implementation must preserve behavior and rerun full tests, doccheck, strict MI `70`, complexity, CRAP, vet, and `git diff --check`.

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
- [x] T004 Stop for explicit approval before moving code.
  - User approved on 2026-05-14. Explicit consent to proceed with Phase 1.

## Phase 1 - First Slice (packet family)

- [x] T010 Pick one small command family.
  - Chosen: `packet` — smallest well-tested family with discrete boundaries.
- [x] T011 Move or index that family only.
  - Renamed 54 files from `main_0xx_...` → `packet_0xx_...` (removing redundant `packet` infixes).
  - Merged 10 tiny files into 3 consolidated files:
    - `packet_032_requiredflags.go` (5 var blocks)
    - `packet_068_artifacttypes.go` (3 type defs)
    - `packet_096_exits.go` (2 funcs)
  - 5 existing `packet_*.go` files required no rename.
- [x] T012 Prove help, docs, exit codes, and tests remain unchanged.
  - `go build ./...`: PASS
  - `go test -count=1 ./...`: PASS
  - `command-surface` JSON: byte-for-byte identical to snapshot
  - `--help` output: byte-for-byte identical to snapshot
  - `go run ./tools/doccheck`: PASS
  - `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`: PASS
  - `go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline ... cmd internal`: PASS
  - `go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline ... cmd internal tools`: PASS
  - `go run ./tools/crapcheck -threshold 5 -strict-less`: PASS
  - `go vet ./...`: PASS
  - `git diff --check`: PASS
- [x] T013 Run PI adversarial review on slice 1.
  - Models: `zai/glm-5.1`, `minimax/MiniMax-M2.7` via `codex-subagent`.
  - Artifacts: `.codex-subagents/runs/run_nwRu9O5MP6/result.md`, `run_bn5qev1MVh/result.md`, `reviews/slice1-synthesis.md`.
  - Findings: no behavioral drift, merged files correct, boundary clean. Advisory on numbering gaps and baseline coverage for new merged files (addressed: new files MI>70, absolute threshold covers them).
  - Overall disposition: ACCEPTED.

## Phase 2 - Iteration

- [x] T020 Repeat by command family with scoped commits.
  - Slice 2: `pr_review` family — 48 files renamed, 9 tiny files merged into 3.
    Commit `ae260a6`.
  - Slice 3: `observe` family — 2 files renamed, 2 tiny files merged into 1.
    Commit `461b7ff`.
  - 12 families remain pending (see T022).
- [x] T021 Keep command handlers discoverable from one registry.
  - `main_006_commandhandlers.go` untouched; all renamed functions resolve
    correctly within `package main`.
- [x] T022 Record any remaining high-file-count area as advisory debt.
  - `cmd/sdp-trace/FAMILY_INDEX.md` maps completed and pending families.
  - 424 `main_*.go` files remain unprefixed.
  - Pending families: assess, gate (checkpoint/override), witness, export,
    release, query (verify/explain/report), doctor, interaction, wrap,
    envelope, fixture, core.
  - Advisory debt recorded in `FAMILY_INDEX.md`.

## Phase 3 - Closure (Draft PR readiness)

- [x] T030 Run full local verification.
  - All verification commands pass for each slice:
    go build, go test, qualitycheck (cyclo/cognitive/MI), crapcheck,
    go vet, doccheck, git diff --check.
  - Command-surface JSON and --help byte-for-byte identical across all slices.
- [ ] T031 Run PI code/correctness, Clean Architecture, DX, and requirements review.
  - Slice 1 reviewed (GLM 5.1 + MiniMax 2.7) — ACCEPTED.
  - Slices 2–3: no code changes inside files (renames + merges only);
    same reasoning applies. PI review deferred to PR level.
- [x] T032 Update contributor navigation docs.
  - `cmd/sdp-trace/FAMILY_INDEX.md` added as human-navigable index.
  - No changes to `docs/agent-entrypoint.md` (command contract unchanged).
