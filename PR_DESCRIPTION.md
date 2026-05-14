# Draft PR: 010 Command Package Organization

## Branch
`010-command-package-organization`

## Objective
Reduce CLI navigation debt in `cmd/sdp-trace/` by grouping files by command family while preserving all behavior and quality gates.

## Strategy
Family-prefixed files within the existing flat `package main`, plus a generated index. No subpackages introduced.

Rationale: zero import-cycle risk, zero registration indirection, zero test breakage.

## Scope
- **Completed**: 3 families reorganized into 119 renamed files + 7 merged files.
  - `packet` (66 files, 3 merged) — commit `ce27273`
  - `pr_review` (59 files, 3 merged) — commit `ae260a6`
  - `observe` (5 files, 1 merged) — commit `461b7ff`
- **Pending**: 12 families (~424 `main_*.go` files) recorded as follow-up in `FAMILY_INDEX.md`.

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
- File MI baseline ratchet: PASS (`qualitycheck -fail-only -mi-under 70 -mi-baseline`).
- Function MI baseline ratchet: PASS (`qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline`).
- CRAP < 5 strict: PASS (`crapcheck -threshold 5 -strict-less`).

## Review Evidence
- **Phase 0 spec review**: GLM 5.1 + MiniMax 2.7 — findings synthesized in `reviews/synthesis.md`.
- **Slice 1 review**: GLM 5.1 + MiniMax 2.7 via `codex-subagent` — ACCEPTED.
  - Artifacts: `reviews/slice1-synthesis.md`, `.codex-subagents/runs/run_nwRu9O5MP6/result.md`, `run_bn5qev1MVh/result.md`.

## Files Changed
- `cmd/sdp-trace/packet_*.go` (66 files)
- `cmd/sdp-trace/pr_review_*.go` (59 files)
- `cmd/sdp-trace/observe_*.go` (5 files)
- `cmd/sdp-trace/FAMILY_INDEX.md` (new)
- `specs/010-command-package-organization/tasks.md` (updated)
- `design-note.md` (new)
- `reviews/*.md` (review artifacts)

## Follow-up Work
- Complete prefixing for remaining 12 families (see `FAMILY_INDEX.md`).
- PR-level PI review for slices 2–3 and overall architecture.

## Checklist
- [x] Spec reviewed and approved
- [x] Behavior preserved with snapshot lock
- [x] Quality gates green
- [x] Scoped commits with evidence
- [x] Contributor navigation docs updated (`FAMILY_INDEX.md`)
- [ ] Remaining families (tracked in `FAMILY_INDEX.md`)
- [ ] Final PR-level review
