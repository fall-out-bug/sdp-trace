# PI Review Synthesis: 010-command-package-organization

## Reviewers
- **GLM 5.1** (`zai/glm-5.1`) — full prompt with spec, plan, tasks, AGENTS.md rules
- **MiniMax 2.7** (`minimax/MiniMax-M2.7`) — compressed prompt due to model input constraints

## Verified Findings

### F-1. Prefixed files in flat `package main` is the safest default (GLM advisory → accepted)
- `cmd/sdp-trace/` is `package main` with no internal consumers.
- Subpackages require exported handler types + registration indirection for zero coupling benefit.
- **Disposition**: Accepted. Spec updated to rank strategies by risk; subpackages require explicit justification.

### F-2. ~72 files ≤8 lines are metric-splitting artifacts (GLM accepted → verified)
- Verified: 72 `.go` files in `cmd/sdp-trace/` contain ≤8 lines.
- These exist because PR #43 split single logical units across files to satisfy MI/complexity tooling.
- **Disposition**: Accepted. Design note adds "merge tiny files into family groups" as a prerequisite.

### F-3. Need concrete snapshot-based behavior lock (GLM accepted_fixed → accepted)
- Existing `command_surface_test.go` checks determinism and coverage, but drift tests can be updated alongside code.
- **Disposition**: Accepted. Plan updated: snapshot `sdp-trace command-surface` JSON before first move; compare after each family slice.

### F-4. Phase 2 needs a stop condition (GLM accepted → accepted)
- Spec said "proceed family by family" without bounding.
- **Disposition**: Accepted. Plan updated: families are enumerated upfront in design note.

### F-5. Low test-file ratio is advisory (GLM advisory → noted)
- 16 test files for ~548 source files.
- **Disposition**: Advisory recorded. Families without dedicated CLI tests treated as higher-risk moves.

### F-6. MiniMax claim "No spec document exists" (MiniMax CRITICAL → rejected_false_positive)
- Verified: `specs/010-command-package-organization/spec.md`, `plan.md`, `tasks.md` exist.
- **Disposition**: Rejected false positive.

### F-7. MiniMax claim "MI 70.1 undefined / not enforced in CI" (MiniMax CRITICAL/HIGH → rejected_false_positive)
- Verified: `docs/agent-entrypoint.md` defines the gate commands explicitly.
- Verified: `.github/workflows/ci.yml` runs `tools/qualitycheck` with `-mi-under 70`, `-function-mi-under 70`, `-cyclo-over 10`, `-cognitive-over 10`, and `tools/crapcheck`.
- **Disposition**: Rejected false positive.

### F-8. Metric-gaming risk via distribution (MiniMax HIGH → advisory)
- Moving atomized files into subpackages could mask complexity without reducing it.
- **Disposition**: Advisory. Same-package prefixes avoid this; aggregate complexity is unchanged.

### F-9. commandSurfaceUsageDrift is fragile (MiniMax MEDIUM → advisory)
- Test compares against checked-in fixture that can be updated in the same commit.
- **Disposition**: Advisory. Snapshot lock (F-3) addresses this.

## Not Assessed
- Exact per-family test coverage percentage (not measured).
- Whether merging 72 tiny files into larger family files triggers MI/CRAP regressions (will be tested in first slice).

## Consensus Strategy
Both reviewers agree: **family-prefixed files within the existing flat `package main`** is the lowest-risk path.
Subpackages are higher risk and require explicit justification not present in the current codebase.
