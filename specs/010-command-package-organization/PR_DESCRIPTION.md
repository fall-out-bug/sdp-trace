# Draft PR: 010 Command Package Organization

## Branch
`010-command-package-organization`

## Objective
Reduce CLI navigation debt in `cmd/sdp-trace/` by grouping all files by command family while preserving all behavior and quality gates.

## Strategy
Family-prefixed files within the existing flat `package main`, plus a generated index. No subpackages introduced.

Rationale: zero import-cycle risk, zero registration indirection, zero test breakage.

## Scope
- **Completed**: All 15 families reorganized into 537 renamed files + 12 merged files across 17 scoped commits.
  - `packet` (66 files, 3 merged) — commit `ce27273`
  - `pr_review` (59 files, 3 merged) — commit `ae260a6`
  - `observe` (5 files, 1 merged) — commit `461b7ff`
  - `witness` (24 files, 1 merged) — commit `405cb15`
  - `release` (5 files, 1 merged) — commit `4aafd4e`
  - `export` + `interaction` + `envelope` + `wrap` + `fixture` (49 files) — commit `eb2cb04`
  - `query` (22 files, 1 merged) — commit `0f035d1`
  - `assess` (55 files) — commit `5a1a529`
  - `gate` (103 files) — commit `bf21a3a`
  - `doctor` (52 files) — commit `87597b8`
  - `core` (100 files) — commit `cd1ab63`
  - gofmt fix (review finding F-1) — commit `22cbd8a`
  - remove binary + trailing whitespace — commit `1ab762a`
  - revert baseline JSON to origin/main — commit `d117e12`
  - remove codex-subagent config — commit `c1ec86b`
  - revert command-surface registry rename — commit `dd917fa`
- **Intentionally unchanged**: `main_536...main_546` command-surface registry files and `main_test.go` remain as `main_` because they are below the MI threshold and tracked in the baseline. Future refactor PR can raise their MI above 70 and rename them.

## Behavior Preservation Evidence
- `command-surface` JSON: byte-for-byte identical to baseline (snapshot `5302ed0`).
- `--help` output: byte-for-byte identical to baseline.
- `go test -count=1 ./...`: PASS across all slices.
- `go build ./...`: PASS across all slices.
- `go vet ./...`: PASS.
- `go run ./tools/doccheck`: PASS.
- `git diff --check`: PASS (no whitespace errors).

## Quality Gate Evidence
- Cyclomatic complexity <= 10: PASS (`qualitycheck -fail-only -cyclo-over 10`).
- Cognitive complexity <= 10: PASS (`qualitycheck -fail-only -cognitive-over 10`).
- CRAP < 5 strict: PASS (`crapcheck -threshold 5 -strict-less`, fresh pipeline replayed).
- MI baseline change policy: PASS (`mibaselinepolicy -base-ref origin/main` — no baseline changes in this diff).
- File/function MI baseline ratchet: PASS (uses origin/main baseline unchanged).

## Review Evidence
- **Phase 0 spec review**: GLM 5.1 + MiniMax 2.7 — findings synthesized in `reviews/synthesis.md`.
- **Slice 1 review**: GLM 5.1 + MiniMax 2.7 via `codex-subagent` — ACCEPTED.
  - Artifacts: `reviews/slice1-synthesis.md`.
- **Slices 2–14 review**: Subagent panel runs for parallel batches — ACCEPTED with advisory conditions.
- **Final PR-level review** (`panel_FRY3ES0Hyc`, 2026-05-14):
  - **Reviewer**: ADVISORY APPROVED with conditions.
    - F-1 (gofmt): **FIXED** in commit `22cbd8a`.
    - F-2 (File MI baseline incomplete): pre-existing debt from PR #44/009; no regression.
    - F-3 (Function MI baseline schema): pre-existing tooling question; no regression.
    - F-4 (CRAP not verified): addressed — fresh coverage pipeline replayed.
  - **Evidence-reviewer**: **ACCEPTED**. All gates verified live. No veto conditions.
- **Post-review fixes** (addressing blocking comments in PR #46):
  - Removed checked-in ELF binary `sdp-trace` + added to `.gitignore` (`1ab762a`).
  - Fixed trailing whitespace in `design-note.md`, `reviews/pi-review-glm-5.1.md`, `reviews/pi-review-minimax-2.7.md` (`1ab762a`).
  - Reverted baseline JSON changes; no baseline changes in this PR (`d117e12`).
  - Removed `.agents/codex-subagents/` harness-specific config (`c1ec86b`).
  - Reverted rename for `main_536...main_546` command-surface registry files to avoid MI baseline sequencing deadlock (`dd917fa`).

## Files Changed
- `cmd/sdp-trace/FAMILY_*.go` (537 renamed + 12 merged files across all families)
- `cmd/sdp-trace/FAMILY_INDEX.md` (new, all 15 families)
- `specs/010-command-package-organization/tasks.md` (updated)
- `design-note.md` (new)
- `reviews/*.md` (review artifacts)
- `.gitignore` (add `sdp-trace`)

## Dependencies
- None. PR is self-contained and mergeable independently.

## Follow-up Work
- Future refactor PR: raise MI of `main_536...main_546` command-surface registry files above 70, then rename to `core_` prefix.
- CI run on this PR.

## Checklist
- [x] Spec reviewed and approved
- [x] Behavior preserved with snapshot lock
- [x] Quality gates green (local)
- [x] Scoped commits with evidence
- [x] All 15 families completed
- [x] Contributor navigation docs updated (`FAMILY_INDEX.md`)
- [x] Final PR-level review (panel_FRY3ES0Hyc)
- [x] Review findings addressed (F-1 gofmt fixed)
- [x] Binary removed from PR
- [x] Trailing whitespace fixed
- [x] No baseline changes in this PR
- [x] Harness-specific config removed
- [x] Command-surface registry files kept as main_ to avoid baseline deadlock
- [ ] CI verification on PR
- [ ] **Explicit user approval before merge**
