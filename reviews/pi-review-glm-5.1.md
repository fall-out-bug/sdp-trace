Warning: No models match pattern "kimi-coding/k2p6"
Now I have a thorough understanding of the codebase. Let me produce the adversarial review.

---

# Adversarial Review: Command Package Organization Spec

## Summary (ordered by severity)

### F-1. **HIGH — Spec proposes refactoring a single-package flat namespace that has no import coupling problem**

The entire `cmd/sdp-trace/` is `package main` — 548 non-test files, ~17.7K lines, 751 functions. There are **zero import cycles possible** today because all files share one package. The spec's FR-004 ("If subpackages are introduced, dependency direction stays clear") guards against a problem that does not yet exist and introduces risk for no proven benefit.

**Evidence**: Every file begins with `package main`. The `commandHandlers` registry in `main_006_commandhandlers.go` is a flat `map[string]commandHandler` with 27 entries. Subpackages would break this into separate packages that must export handler types, requiring an interface or registration indirection that currently does not exist.

**Risk**: Introducing subpackages into `cmd/sdp-trace/` would require either:
- (a) exported handler types + a registration API, or
- (b) a separate `cmd/sdp-trace/internal/` layer with reverse imports.

Both add indirection for a CLI entry point that currently has **no internal consumers**. The `internal/` packages already own the domain logic; the CLI is pure wiring.

**Disposition**: **advisory** — The spec should explicitly constrain Phase 0 to evaluate whether *any* subpackage strategy produces measurable navigation improvement over prefixed files in the existing flat package, given that the real problem is 71 files under 8 lines each that exist only to satisfy MI/complexity tooling.

### F-2. **HIGH — 71 files ≤8 lines are metric-splitting artifacts, not organizational debt**

Of 548 non-test source files, **71 (13%) contain 8 lines or fewer**, and 36 of those contain 5 lines or fewer. These files exist because PR #43 split single functions across files to satisfy MI/CRAP thresholds. Examples:

- `main_285_checkpointverifystringflags.go` (3 lines): just `var checkpointVerifyStringFlags = []string{"run", "checkpoint", "policy"}`
- `main_037_prreviewpacketrequiredflags.go` (5 lines): just a `[]requiredCLIFlag` slice literal
- `main_089_successfulhttpstatus.go` (5 lines): one boolean function
- `main_535_rest.go` (5 lines): one struct literal

**Evidence**: `main_037_prreviewpacketrequiredflags.go`, `main_285_checkpointverifystringflags.go`, `main_089_successfulhttpstatus.go`, `main_535_rest.go` — confirmed by `wc -l` output above.

The spec's "Problem Statement" frames this as navigation debt, but the root cause is metric-driven splitting, not organizational choice. Reorganizing these into command-family subpackages would relocate the atomized files without solving the real problem: a single logical function's flag declarations or type definitions are spread across multiple files.

**Disposition**: **accepted** — Finding is real. The spec should acknowledge that file-count reduction (merging tiny same-family files) is a prerequisite to any reorganization; moving 71 tiny files into subpackages creates subpackages full of tiny files.

### F-3. **MEDIUM — The spec's three strategy options are not equal in risk**

FR-001 lists three strategies: subpackages, prefixed files with generated index, or hybrid. The spec does not state a preference, deferring to Phase 0. But the risk profiles are dramatically different:

| Strategy | Import cycle risk | Behavior change risk | Navigation improvement |
|---|---|---|---|
| Subpackages | High (new exports, registration API) | Medium (init ordering, test scope) | Medium |
| Prefixed files + generated index | None (same package) | Low | High (grep-friendly) |
| Hybrid | Medium | Medium | Medium |

The existing code already has meaningful named files (`harness_cli.go`, `harness_observe_cli.go`) alongside the numbered files. The natural low-risk path is *renaming numbered files with family prefixes within the same package* plus an index file.

**Evidence**: `harness_*.go` files already demonstrate that named, non-numbered files work in this package. The `main_006_commandhandlers.go` registry is importable from any file in `package main` without change.

**Disposition**: **accepted** — Spec should rank strategies by risk in Phase 0, not present them as equivalent. The highest-risk strategy (subpackages) should require explicit justification over the zero-risk alternative.

### F-4. **MEDIUM — Behavior preservation is testable but the spec has no concrete test plan**

US-002 says "Existing CLI commands, help text, exit codes, and docs checks remain unchanged." The existing test suite confirms this is testable — `command_surface_test.go` covers handler coverage and output determinism, and the command surface JSON machinery in `main_536`-`main_546` provides a machine-readable contract.

However, the spec's tasks (T012) only say "Prove help, docs, exit codes, and tests remain unchanged" without specifying:
- What pre/post snapshot comparison mechanism?
- Whether the `commandSurfaceJSON` test is the canonical behavior lock?
- How to detect help-text regressions (the existing `commandSurfaceDrift` test checks schema, not help strings).

**Evidence**: `command_surface_test.go` contains `TestCommandSurfaceJSONIsDeterministic`, `TestCommandSurfaceCoversAllHandlers`, and `TestCommandSurfaceUsageDrift`. But `TestCommandSurfaceUsageDrift` compares against a checked-in fixture — if you move files and update the fixture in the same commit, the drift test passes by definition.

**Disposition**: **accepted_fixed** — The spec should require a pre-move snapshot of command-surface JSON and help output, and task T012 should specify that the snapshot is taken *before* any file movement and compared *after* via a separate verification step, not via the existing drift test that can be updated alongside the code.

### F-5. **MEDIUM — Phase 2 "Repeat by family" has no stop condition**

The spec says "Proceed family by family" but does not define when to stop. Based on the codebase, there are at least 12 distinct command families (packet, pr-review, assess, witness, query, gate, checkpoint, doctor, harness, recorder, interaction, export) plus shared infrastructure (flag parsing, JSON output, command surface). That's potentially 12+ sequential slice-and-verify cycles, each requiring a scoped commit and quality gate run.

**Evidence**: The `commandHandlers` map in `main_006_commandhandlers.go` lists 27 commands. The `internal/` directory contains 23 subpackages.

**Disposition**: **accepted** — The spec should bound Phase 2: define the families upfront in Phase 0, count them, and state that the work is done when all families are organized. Otherwise "repeat by family" can drift indefinitely.

### F-6. **LOW — 16 test files for 548 source files (2.9% test-file ratio) is a fragility signal**

The `cmd/sdp-trace/` directory has only 16 `*_test.go` files covering 563 source files. While some tests are integration-weight (e.g., `pr_review_cli_test.go` at 403 lines), the test coverage ratio suggests that behavior preservation during refactoring relies heavily on compilation and a few broad integration tests, not on targeted per-handler unit tests.

**Evidence**: `ls cmd/sdp-trace/*_test.go | wc -l` → 16. These include `command_surface_test.go`, `crap_hotspot_test.go`, `harness_cli_test.go`, `pr_review_cli_test.go`, etc.

**Disposition**: **advisory** — Not a spec defect, but the spec should note that the safety of moving code is proportional to test density per family. Families without dedicated CLI tests (checkpoint, gate, doctor) should be treated as higher-risk moves.

### F-7. **LOW — Task numbering has gaps**

Tasks go T001–T004, then T010–T012, T020–T022, T030–T032. The gaps (T005–T009, T013–T019, etc.) suggest either incomplete decomposition or reserved slots. This is cosmetic but makes it harder to track completeness.

**Disposition**: **rejected_false_positive** — Gap numbering is intentional for phase grouping. Not a trust issue.

### F-8. **LOW — The spec claims "small tests that guard behavior while moving code" as allowed, but this could re-introduce the metric-driven splitting problem**

The Product Boundary allows "small tests that guard behavior while moving code." Given that the current file explosion was caused by metric-driven splitting, adding more small test files without constraints could worsen the file-count problem that triggered this spec.

**Disposition**: **advisory** — Any new test files should be justified per family, not per function.

---

## `not_assessed` Areas

- **Import cycle detection at scale**: I did not run `go vet` or a cycle detector against hypothetical subpackage layouts because no code has been written yet. The risk assessment in F-1 is based on structural reasoning, not live verification.
- **MI/CRAP metrics under merged files**: I did not test whether merging the 71 tiny files into family-grouped files would cause MI or CRAP regressions. The spec assumes quality gates stay green, but the metric tooling's behavior with larger files is unverified.
- **Test coverage percentage per command family**: I counted test files but did not measure line coverage per family. Some families may have near-zero coverage.

---

## Overall Disposition

The spec identifies a real DX problem but underestimates its root cause (metric-driven atomization) and overestimates the benefit of subpackages. The safest path — family-prefixed file renaming within the same `package main` — is listed as one of three equivalent options when it should be the default. The spec is **implementable** after addressing F-1 (constrain subpackage justification), F-2 (acknowledge merging as prerequisite), F-3 (rank strategies by risk), and F-4 (concrete snapshot-based behavior lock).
