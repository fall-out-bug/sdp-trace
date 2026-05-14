# Design Note: Command Package Organization

**Feature**: 010-command-package-organization
**Date**: 2026-05-14
**Status**: Reviewed, awaiting implementation approval

## Chosen Strategy

**Family-prefixed files within the existing `package main` flat package, plus a generated index.**

No subpackages are introduced in `cmd/sdp-trace/`.

## Rationale

1. `cmd/sdp-trace/` is `package main` with zero internal consumers. Introducing subpackages adds exported handler types, registration indirection, and potential init-order fragility for no coupling benefit.
2. Same-package renames carry **zero import-cycle risk**, **zero behavior-change risk**, and **zero test-scope breakage**.
3. The existing `commandHandlers` registry (`main_006_commandhandlers.go`) continues to work unchanged; no registration API is needed.
4. Quality gates (MI, cyclomatic, cognitive, CRAP) are measured per-function or per-file; same-package renames do not affect metrics.
5. Reviewers and agents can navigate by prefix (`packet_`, `pr_review_`, `assess_`, etc.) without scanning 500+ similarly named files.

## Families

Enumerated from the `commandHandlers` registry and grouped by ownership:

| Family | Commands / Scope | Example prefix |
|--------|------------------|----------------|
| `core` | Registry, dispatch, flags, shared helpers, version, command-surface | `core_` |
| `wrap` | wrap, run, dry-run, preview | `wrap_` |
| `doctor` | doctor, install | `doctor_` |
| `interaction` | interaction relay/import/summarize | `interaction_` |
| `observe` | observe, harness | `observe_` |
| `query` | verify, explain, query, query-pack, report | `query_` |
| `assess` | assess (all profiles) | `assess_` |
| `gate` | gate, override, checkpoint | `gate_` |
| `witness` | witness | `witness_` |
| `export` | export (telemetry, cross-repo-posture) | `export_` |
| `release` | release-proof | `release_` |
| `pr_review` | pr-review (packet, run, synthesize, validate, summarize, check) | `pr_review_` |
| `packet` | packet (build-pr, build-github, validate, check-demo, render) | `packet_` |
| `envelope` | envelope summarize | `envelope_` |
| `fixture` | validate-fixtures | `fixture_` |

## Pre-requisite: Merge Tiny Files

~72 files in `cmd/sdp-trace/` contain ≤8 lines. Many are metric-splitting artifacts from PR #43. Before prefixing, merge tiny files that belong to the same family into slightly larger files (still under complexity thresholds) so the reorganization reduces file count, not just reshuffles it.

Rules for merging:
- Only merge files from the **same family**.
- Merged file must still pass cyclomatic ≤10, cognitive ≤10, CRAP <5, and MI baseline.
- Do not merge test files with source files.

## Behavior Lock

Before any file moves:
1. Run `go run ./cmd/sdp-trace command-surface > .sdp-trace-cmd-surface-before.json`.
2. Run `go run ./cmd/sdp-trace --help > .sdp-trace-help-before.txt`.
3. Record `git rev-parse HEAD`.

After each family slice, diff the snapshots. Any delta is a regression until explicitly reviewed.

## Index Generation

After all families are prefixed, generate `cmd/sdp-trace/FAMILY_INDEX.md` mapping each prefix to its commands and primary files. This is the human-navigable contract.

## Slice Order

1. `packet` — small, well-tested (`packet_cli_test.go`), discrete boundaries.
2. `pr_review` — medium, well-tested (`pr_review_cli_test.go`), but many files.
3. `observe` — medium, includes harness subcommands.
4. `assess` — large, many profiles.
5. `gate` — medium.
6. `witness` — small.
7. `export` — small.
8. `release` — small.
9. `query` — medium.
10. `doctor` — small.
11. `interaction` — small.
12. `wrap` — medium.
13. `envelope` — small.
14. `fixture` — small.
15. `core` — last, because it contains shared infrastructure used by all other families.

## Stop Condition

Phase 2 stops when all 15 families above are prefixed and indexed. Any remaining advisory debt is recorded in `specs/010-command-package-organization/advisory-debt.md`.

## Risk Summary

| Risk | Mitigation |
|------|------------|
| Behavior change | Snapshot lock; same-package renames only. |
| Quality gate regression | Merge only same-family tiny files; re-run gates after each slice. |
| Import cycle | Not possible — no new packages. |
| Test breakage | Tests remain in `package main`; no import changes. |
| Metric gaming | Aggregate complexity unchanged; prefixes are renames, not splits. |
