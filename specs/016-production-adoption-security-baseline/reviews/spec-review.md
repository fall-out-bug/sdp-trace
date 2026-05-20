# Spec Review 016: Production Adoption And Security Baseline

Reviewer: self (Socratic review during intake)
Date: 2026-05-20

## Findings

### Scope
- **Status: acceptable**
- Spec correctly separates "adoption facts" from "production trust claims".
- Acceptance criteria do not claim merge/release approval or external production trust.

### Completeness
- **Finding C1 (minor)**: `gosec` policy (blocking vs advisory) listed in Open Questions but not resolved before implementation phase.
  - Disposition: accepted — will be resolved in WS-016-B triage.
- **Finding C2 (minor)**: Security policy path (`SECURITY.md` vs `.github/SECURITY.md`) listed in Open Questions.
  - Disposition: accepted — will use `.github/SECURITY.md` per GitHub community standard.

### Feasibility
- **Status: acceptable**
- 133 `gosec` findings are classified, not all fixed. Matches spec scope.
- 12 `gitleaks` findings are small enough to triage in one slice.

### Security
- **Finding S1 (major)**: Spec does not explicitly require `gitleaks` allowlist/config to be stored in-repo vs in CI secrets.
  - Disposition: accepted — will store narrow `.gitleaks.toml` or inline allowlist rationale in `docs/security-baseline.md`.
- **Finding S2 (major)**: `gosec` G101 (hardcoded credentials, 2 findings) must be classified as blocker until proven false positive.
  - Disposition: accepted — will verify each G101 in WS-016-B.

### DX/UX
- **Finding D1 (minor)**: No explicit requirement that security docs be linked from README/docs index.
  - Disposition: accepted — will add links in WS-016-A integration step.

## Resolved Open Questions

1. `gosec` policy: advisory for this phase; blocker if new G101/G204 appear after baseline.
2. Security policy path: `.github/SECURITY.md`.
3. Ignored local agent run directories: scanned by hygiene tooling but reported as local clutter, not repository proof.

## Disposition

All findings accepted or deferred to implementation slices. Proceed to Phase 1.

## Remaining `not_assessed`

- Live CI status for this block.
- External customer adoption evidence.
- External security audit.
