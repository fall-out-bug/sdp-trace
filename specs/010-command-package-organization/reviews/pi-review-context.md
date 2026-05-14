# PI Review Context: 010-command-package-organization

## Objective
Review whether the spec for reorganizing `cmd/sdp-trace` can reduce CLI navigation debt without undoing quality-gate work.

## Repository Rules (from AGENTS.md)
- Machine proof wins over prose.
- No deferred trust closure.
- Source-bound proof requires clean immutable source commit.
- No TODO/FIXME markers in new Go code.
- No Node.js/TS/JS in active product path.
- Bash only as thin command launcher.
- Keep root router under 100 lines.
- Module with >10 skills is too large.
- Test-first behavior for behavior changes.
- Go is the target product code.

## Current State
- `cmd/sdp-trace/` contains 500+ small Go files (~17.7K lines total).
- Files are numbered (`main_001_const.go` … `main_546_commandsurfacejson.go`) plus named files.
- Quality gates pass: strict MI 70.1, cyclomatic/cognitive <=10, CRAP strict <5, coverage-backed.
- The codebase passed strict metrics through broad same-package file splitting (PR #43).
- Review accepted result but recorded navigation and contributor-DX debt.

## Spec Under Review

### Core Claim
> The CLI implementation is organized by command family or another reviewed structure while preserving command behavior and quality gates.

### Proposed Organization Strategies (FR-001)
1. Command-family subpackages
2. Family-prefixed files with generated index
3. Hybrid grouping for only the largest families

### Candidate Families
- packet
- PR review
- assess
- witness
- observe/harness
- release proof
- query/report/gate

### Required Tests
- `go test -count=1 ./...`
- `go run ./cmd/sdp-trace --help`
- `go run ./tools/doccheck`
- `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd internal tools`
- `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd internal`
- complexity gates (cyclomatic/cognitive <=10)
- `go run ./tools/crapcheck` with coverage
- `go vet ./...`
- `git diff --check`
- `go list ./...` for cycles

### Constraints
- No CLI behavior changes.
- No exit semantics changes.
- No lowering quality gates.
- No non-Go tooling.
- No dependency inversion (core trace/evidence must not depend on CLI).

## Review Questions
1. Is the proposed organization small enough for safe slices?
2. Does it risk dependency cycles or import churn?
3. Is behavior preservation testable?
4. Does the spec prevent metric-gaming from being mistaken for Clean Code?
5. Which organization strategy (subpackages, prefixed files, hybrid) is safest given the 500+ file count?
6. Are there hidden risks in moving numbered files that qualitycheck baselines may bind to?
7. Is the acceptance criteria complete?

## Output Format
- Summary first, ordered by severity.
- File/line evidence where applicable.
- Disposition per finding: accepted, accepted_fixed, rejected_false_positive, deferred_not_assessed, cannot_verify, advisory.
- Note any `not_assessed` areas.
