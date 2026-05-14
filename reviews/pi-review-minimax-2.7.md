Warning: No models match pattern "kimi-coding/k2p6"


Based on my analysis of the codebase, here is the severity-ordered risk review:

---

# Risk Review: CLI Package Organization Refactoring

## Q1: Which Strategy is Safest?

**Hybrid (family-prefixed files in flat `package main`, with optional subpackage extraction for truly isolated helpers) is safest**, with caveats.

### Evidence from codebase:
- All 579 files share **one package** (`package main`)
- Critical shared types: `flagSet`, `commandHandler`, `subcommandHandler`, `commandSurfaceSchema`
- Handler dispatch via `commandHandlers` map in `main_006_commandhandlers.go`
- `commandSurfaceSchema` in `main_536_commandsurfaceschema.go` is the behavioral contract

### Why subpackages fail the trust bar:

| Risk | Subpackages | Family-prefixed flat | Hybrid |
|------|-------------|---------------------|--------|
| Breaking shared types | **CRITICAL** — requires `type alias + interface` dance | **None** — zero changes needed | Low |
| Dependency cycle creation | High — requires reversing `cmd→internal` direction | **None** | Low |
| `crap_hotspot_test.go` breakage | Likely — test helpers are package-local | None | Negligible |
| Test file renames | ~10 test files must move | `*_test.go` suffix preserves discovery | Subset |
| Exit code mapping changes | Possible if `flagSet` moves | None | Negligible |
| `commandSurfaceUsageDrift()` impact | Low | None | None |

---

## Q2: Risk of Undoing Quality-Gate Work

**Severity: CRITICAL**

### Finding 1 — Cyclomatic/Cognitive Complexity NOT enforced in CI

- **Severity**: CRITICAL
- **Evidence**: `crap_hotspot_test.go` only tests exit-code mapping tables and preview actions. No `go test -coverprofile` gate, no `gocyclo -max=10`, no `crap` tool invocation, no cognitive complexity measurement.
- **Disposition**: **Spec cannot be validated as written.** Before any reorganization, add to CI:
  ```bash
  go test ./... -coverprofile=coverage.out
  which go-jpiry || go install github.com/bk211/fz@latest  # or equivalent
  gocyclo -max=10 ./cmd/sdp-trace/...
  ```
  Or the spec must state "metrics enforced by external tooling not in this repo."

### Finding 2 — CRAP<5 not programmatically enforced

- **Severity**: HIGH
- **Evidence**: No evidence of CRAP score measurement. `crap_hotspot_test.go` tests behavior only.
- **Disposition**: Verify if CRAP<5 is a CI gate or aspirational prose. If CI-gated, add `crap` tool to CI before refactoring; if aspirational, remove from spec or mark `not_assessed`.

### Finding 3 — Metric-gaming risk during reorganization

- **Severity**: HIGH
- **Risk**: Refactoring "packet" helper files into `packet/` subpackage could reduce per-file metrics by distributing complexity. An agent could claim "CRAP improved" while total CRAP for the family is unchanged or worse.
- **Disposition**: Require **aggregate metrics** (not per-file) as the gate. Add `go-jpiry` or equivalent to verify the *total* CRAP of the reorganized package ≤ current total.

### Finding 4 — `command_surface_test.go` drift detection is the real behavioral gate

- **Severity**: MEDIUM (but the gate is fragile)
- **Evidence**: `commandSurfaceUsageDrift()` is the only automated behavior-preservation check. It compares registry vs `usageText`. But it doesn't verify:
  - Exit codes for new commands
  - Flag defaults
  - Subcommand handler registration
- **Disposition**: Strengthen this test before refactoring: add `TestCommandSurfaceExitCodesMatchHandlers`, `TestCommandSurfaceFlagDefaultsStable`.

---

## Q3: Is Behavior Preservation Testable Enough?

**Severity: MEDIUM**

### Finding 5 — Behavioral surface is partially covered

- **Evidence**:
  - `TestCommandSurfaceCoversAllHandlers` ✓ (registry count = handler count)
  - `TestCommandSurfaceUsageDrift` ✓ (usage text drift)
  - `TestCommandSurfaceIncludesKnownProfiles/WitnessKinds/States` ✓ (enumeration)
  - Exit codes tested in `TestStateExitCodeTables`
  - Subcommand failure paths tested in `TestInteractionCommandFailurePaths`, `TestManagedAndRepoObserverCommandFailurePaths`

- **Gaps** (not tested):
  1. **Flag value defaults** — `flagSet.stringValue("missing")` returns `""`, but no test ensures defaults are stable after reorganization
  2. **Subcommand handler registration completeness** — `runInteraction` vs `runInteractionImportTranscript` etc. are tested per-subcommand but the registration map is not snapshot-tested
  3. **JSON output stability** — `TestCommandSurfaceJSONIsDeterministic` only checks round-trip, not schema contract
  4. **Exit code for `exitUsage`** — no test verifies unrecognized flags return `exitUsage`

### Finding 6 — Test file co-location is fragile

- **Evidence**: All test files are in `package main`. If `main_520_flagset.go` moves to a subpackage, `flagSet` becomes `package packet` and test files must `import` or move.
- **Disposition**: Before refactoring, consider moving test helpers (`chdirCLI`, `writeHarnessCLIProfile`, etc.) into a `testutil/` package within `cmd/sdp-trace` to avoid package-boundary test breakage.

### Recommendation:

Add these before refactoring:
```go
// cmd/sdp-trace/command_surface_stability_test.go
func TestCommandSurfaceSchemaContractStability(t *testing.T) {
    // Verify every handler in commandHandlers has an entry in buildCommandSurface()
    // Verify every flag default matches expected values
    // Verify exit codes for all state variants
}
```

---

## Q4: Does Spec Prevent Metric-Gaming Mistaken for Clean Code?

**Severity: HIGH**

### Finding 7 — No spec document exists; claims are unenforceable

- **Severity**: CRITICAL
- **Evidence**: I found no reorganization spec in the worktree. The constraints (`MI 70.1`, `cyclomatic ≤10`, `CRAP<5`) are stated in the request but not formalized in any `docs/` or `M###/` file.
- **Disposition**: **Produce the spec first.** Without a formal spec, "behavior preservation" and metric thresholds are prose claims that cannot be machine-verified.

### Finding 8 — MI 70.1 is undefined in the codebase

- **Severity**: HIGH
- **Evidence**: `crap_hotspot_test.go`, `command_surface_test.go`, and all `*_test.go` files contain no reference to "MI", "Maintainability Index", or a score threshold.
- **Disposition**: Define MI 70.1:
  - Is this Visual Studio's maintainability index?
  - Is there a tool (`gomi`, `lizard`, `oximeter`) that produces it?
  - Is it calculated per-file or aggregate?
  Without this, the claim is unfalsifiable and violates the **Trust Rules** ("Every claim about a gate or verdict must be evidence-backed").

### Finding 9 — "Strict MI 70.1" is prose, not a gate

- **Severity**: MEDIUM
- **Risk**: An agent could claim "MI is preserved" while refactoring files in a way that *increases* system complexity because no measurement exists.
- **Disposition**: Add measurement to CI or mark `not_assessed` with a tracked follow-up.

### Finding 10 — Metric distribution masking

- **Severity**: MEDIUM
- **Risk**: Moving 50 small files into a `packet/` subpackage could make individual file metrics look better (each file is simpler) while the *package-level* complexity is unchanged or increased.
- **Disposition**: If subpackages are chosen, require **package-level** aggregate metrics, not per-file limits.

---

## Summary Table

| # | Finding | Severity | Disposition |
|---|---------|----------|-------------|
| 1 | Cyclomatic/Cognitive NOT enforced in CI | CRITICAL | Add to CI before refactoring |
| 2 | CRAP<5 not programmatically enforced | HIGH | Verify CI gate or mark `not_assessed` |
| 3 | Metric-gaming risk via distribution | HIGH | Require aggregate metrics as gate |
| 4 | `command_surface_test.go` drift detection is fragile | MEDIUM | Strengthen test coverage |
| 5 | Behavioral surface partially covered | MEDIUM | Add stability tests for flags, exit codes |
| 6 | Test file co-location fragile | MEDIUM | Extract test helpers before refactoring |
| 7 | No spec document exists | CRITICAL | Produce spec before work |
| 8 | MI 70.1 undefined | HIGH | Define tool + measurement scope |
| 9 | "Strict MI 70.1" is prose, not gate | MEDIUM | Add measurement or mark `not_assessed` |
| 10 | Metric distribution masking | MEDIUM | Require package-level aggregate metrics |

---

## Recommended Spec Additions

```
## Quality Gates (required before first commit)
- [ ] `go test ./cmd/sdp-trace/...` passes with 0 test changes
- [ ] `gocyclo -max=10 ./cmd/sdp-trace/...` reports 0 files over threshold
- [ ] Aggregate CRAP for reorganized units ≤ current aggregate CRAP
- [ ] `commandSurfaceUsageDrift()` returns zero drift after reorganization
- [ ] All `exitCannotVerify`, `exitUsage` constants map to same values post-refactor

## Behavior Preservation Contract
- [ ] `TestCommandSurfaceCoversAllHandlers` passes unchanged
- [ ] `TestStateExitCodeTables` passes unchanged
- [ ] `sdp-trace command-surface` JSON output is byte-for-byte identical

## Refactoring Strategy
- [ ] Phase 1: Family-prefixed renames only (no package changes)
- [ ] Phase 2: Extract test helpers to `testutil/` package
- [ ] Phase 3: Subpackage extraction only for functions with zero `package main` dependencies
- [ ] No subpackage may contain `flagSet`, `commandHandler`, or `commandSurfaceSchema`
```
