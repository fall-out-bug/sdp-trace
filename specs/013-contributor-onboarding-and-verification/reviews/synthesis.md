# Review Synthesis: 013 Contributor Onboarding And Verification

## Review Date
2026-05-15

## Reviewers
- deepseek-v4-pro (adversarial code/correctness, requirements-vs-implementation, DX/UX)
- qwen-3.6-max-preview (adversarial code/correctness, requirements-vs-implementation, DX/UX, security/trust)
- Self-review (verification replay, synthesis hygiene)

## Context Pack
- Objective: Add a single contributor quick-start page with canonical smoke path, remove duplicated command blocks, and extend doccheck to validate onboarding references.
- Changed files: `README.md`, `docs/README.md`, `docs/agent-onboarding.md`, `docs/install.md`, `docs/contributor-quickstart.md`, `tools/doccheck/main.go`, `tools/doccheck/quickstart.go`, `tools/doccheck/quickstart_test.go`, `tools/doccheck/registry.go`.
- Verification commands: `go test -count=1 ./...`, `go run ./tools/doccheck`, smoke path replay, CRAP check (`go run ./tools/crapcheck -threshold 5 -strict-less`).

## Adversarial Review Findings

### Blocking (fixed)

**B1. `html.UnescapeString` in `addUsage` is dead code**
- **File**: `tools/doccheck/registry.go`
- **Finding**: `json.Unmarshal` already decodes `\u003c` to literal `<`. `html.UnescapeString` handles HTML entities (`&lt;`), not JSON unicode escapes. It has zero effect.
- **Fix**: Removed `html.UnescapeString` and the `"html"` import.
- **Disposition**: accepted_fixed

**B2. Stale command detection too permissive (undocumented boundary)**
- **File**: `tools/doccheck/quickstart.go` — `isKnownCommand`, `registryHasBase`, `baseCommand`
- **Finding**: `isKnownCommand` accepts any flags for a known subcommand because it falls back to base-command matching. The existing test only checked an unknown subcommand (`stale-command-that-does-not-exist`), giving false confidence.
- **Fix**: Added `TestCompareQuickstartWithRegistryAcceptsBogusFlagsForKnownSubcommand` documenting the intentional design choice (registry stores placeholder patterns, not concrete flag values). Added comment to `isKnownCommand` explaining the boundary.
- **Disposition**: accepted_fixed (documented limitation)

### Major (fixed)

**M1. `setContainsPrefix` uses `strings.Contains` instead of `strings.HasPrefix`**
- **File**: `tools/doccheck/quickstart.go:53`
- **Finding**: Function name says "prefix" but implementation checks substring membership.
- **Fix**: Replaced `strings.Contains` with `strings.HasPrefix`.
- **Disposition**: accepted_fixed

**M2. `README.md` "Recommended Reading Order" missing quickstart**
- **File**: `README.md`
- **Finding**: The "Start Here" section links to quickstart, but the later "Recommended Reading Order" list omits it.
- **Fix**: Added `Contributor Quick Start` as item 4 in the Recommended Reading Order.
- **Disposition**: accepted_fixed

**M3. Synthesis overclaims `html.UnescapeString` as verified fix**
- **File**: `specs/013-contributor-onboarding-and-verification/reviews/synthesis.md` (self-review artifact)
- **Finding**: Original synthesis claimed the change "handles JSON-encoded `<`/`>`" when it is dead code.
- **Fix**: Corrected synthesis to state the finding and its removal.
- **Disposition**: accepted_fixed

**M4. Expected results table mixes exit codes and verifier states**
- **File**: `docs/contributor-quickstart.md`
- **Finding**: `--help` row used `exits 0` while other rows used verifier states like `observed`.
- **Fix**: Standardized on `exit 0 (...)` for CLI-only commands and verifier states for verifier commands.
- **Disposition**: accepted_fixed

### Minor (fixed or documented)

**m1. Test sort-order fragility**
- **File**: `tools/doccheck/quickstart_test.go`
- **Finding**: `TestQuickstartCommandsExtractsGoRunLines` compared positional slices that happen to match sort order.
- **Fix**: Added `sort.Strings(want)` and added `go test`/`go build` lines to the test input to verify non-`sdp-trace` commands are excluded.
- **Disposition**: accepted_fixed

**m2. `isQuickstartCommand` trailing-space constraint undocumented**
- **File**: `tools/doccheck/quickstart.go`
- **Finding**: The trailing space after `sdp-trace` is intentional but unexplained.
- **Fix**: Added comment explaining the trailing-space requirement and `--help` handling.
- **Disposition**: accepted_fixed

**m3. Doccheck does not verify cross-doc links**
- **File**: `tools/doccheck/main.go`
- **Finding**: Doccheck validates quickstart commands, but does not verify README/install/onboarding actually link to quickstart.
- **Fix**: No code change. This is an intentional scope boundary for Slice 1; cross-doc link drift can be added in a future slice.
- **Disposition**: deferred_not_assessed

**m4. Reviewer entrypoint retains inline smoke commands**
- **File**: `docs/reviewer-entrypoint.md`
- **Finding**: Reviewer entrypoint has its own wrap/verify block using `sdp_trace` shell function. Not drift-checked by doccheck.
- **Fix**: No code change. Reviewer path is intentionally separate (five-minute review vs. contributor environment validation). Added to follow-up tracking.
- **Disposition**: deferred_not_assessed

### Advisory

**a1. Tests shell out to real CLI**
- **File**: `tools/doccheck/quickstart_test.go`
- **Finding**: Tests call `registryUsages()` which runs `go run ./cmd/sdp-trace command-surface`. Requires full checkout and Go toolchain.
- **Disposition**: advisory — follows existing `main_test.go` pattern.

**a2. Cold-reader review remains not_assessed**
- **Disposition**: correctly tracked. Cannot verify in this session.

## Verification State (post-fix)
- `go test -count=1 ./...`: pass
- `go run ./tools/doccheck`: pass (exit 0)
- Smoke path replay (all 6 steps): pass
- `go vet ./...`: pass
- `git diff --check`: pass
- CRAP check (`-threshold 5 -strict-less`): pass (no functions >= 5.00)

## Disposition Summary

| ID | Severity | File | Issue | Disposition |
|---|---|---|---|---|
| B1 | blocking | `registry.go` | dead `html.UnescapeString` | accepted_fixed |
| B2 | blocking | `quickstart.go` | stale detection boundary undocumented | accepted_fixed |
| M1 | major | `quickstart.go:53` | `Contains` vs `HasPrefix` | accepted_fixed |
| M2 | major | `README.md` | missing quickstart in reading order | accepted_fixed |
| M3 | major | `synthesis.md` | overclaim on dead code | accepted_fixed |
| M4 | major | `contributor-quickstart.md` | mixed exit codes / states | accepted_fixed |
| m1 | minor | `quickstart_test.go` | sort-order fragility | accepted_fixed |
| m2 | minor | `quickstart.go` | undocumented trailing space | accepted_fixed |
| m3 | minor | `doccheck/main.go` | no cross-doc link check | deferred_not_assessed |
| m4 | minor | `reviewer-entrypoint.md` | inline smoke not drift-checked | deferred_not_assessed |
| a1 | advisory | `quickstart_test.go` | tests shell out to CLI | advisory |
| a2 | advisory | — | cold-reader review | not_assessed |

## Model Policy Record
- deepseek-v4-pro: completed, actionable findings with file:line evidence
- qwen-3.6-max-preview: completed, actionable findings with file:line evidence
- glm-5.1: unavailable in environment (insufficient balance / invalid model id)
- minimax-2.7: unavailable in environment (invalid api key)

## Conclusion
Adversarial review findings have been addressed. No blocking findings remain. Two items (cross-doc link check, reviewer entrypoint drift check) are deferred to future slices. The branch is ready for PR.
