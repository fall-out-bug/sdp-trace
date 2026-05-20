# Spec Review Context: 017 OSS Replacement Compatibility And Benchmarks

## Objective
Review the SpecKit delta for block 017 before implementation delegation.

## Files Under Review
- `specs/017-oss-replacement-compatibility-and-benchmarks/spec.md`
- `specs/017-oss-replacement-compatibility-and-benchmarks/plan.md`
- `specs/017-oss-replacement-compatibility-and-benchmarks/tasks.md`

## Repository Rules (from AGENTS.md)
- Target product code is Go.
- No Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling in active product path.
- Bash allowed only as thin command launcher.
- New Go code must be small, readable, testable, covered by focused tests.
- No TODO/FIXME markers.
- Root router under 100 lines; module over 10 skills is too large.
- Machine proof wins over prose.
- No deferred trust closure.
- Source-bound proof requires clean immutable source commit.

## Verification Commands
```text
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

## Key Trust Concerns
1. Do requirements FR-017-001 through FR-017-006 create hidden dependencies on Node.js/npm or other disallowed tooling?
2. Does the benchmark harness claim authority it should not claim?
3. Are substitution boundaries explicit enough to prevent implementation accident replacing product code?
4. Is the live `wrap` output/schema drift handled with honest `not_assessed` or `blocked` markers?
5. Does the plan create module size violations (>10 skills per module, router >100 lines)?
6. Are verification commands copy-pasteable and reproducible?

## Reviewer Instructions
- Examine the spec, plan, and tasks for inconsistencies, overclaim, hidden coupling, and trust boundary violations.
- Report file:line evidence for every actionable finding.
- Do not approve implementation if Critical or Important findings remain unresolved.
- Output exactly `LGTM` only if zero findings.
