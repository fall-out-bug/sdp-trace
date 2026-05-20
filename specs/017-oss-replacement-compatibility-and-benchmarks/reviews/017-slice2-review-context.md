# Slice 2 Review Context: tools/ossbench

## Files Changed
- `tools/ossbench/main.go`
- `tools/ossbench/main_test.go`
- `tools/ossbench/bench.go`
- `tools/ossbench/bench_test.go`

## Review Axes
1. **Quality:** Code correctness, test coverage, edge cases, error handling, Go idioms.
2. **UX:** Command-line interface clarity, output formatting, help text, error messages.
3. **DX:** Buildability, testability, readability, file organization, documentation.
4. **Security:** Command injection, path traversal, unsafe exec, data exposure, trust boundaries.

## Repository Rules
- Target product code is Go. No Node.js/npm/JS/TS/.mjs in active product path.
- New Go code must be small, readable, testable, covered by focused tests.
- No TODO/FIXME markers.
- Root router under 100 lines; module over 10 skills is too large.

## Verification Commands
```text
go test -count=1 ./tools/ossbench
go vet ./tools/ossbench
go run ./tools/ossbench -list
go run ./tools/ossbench -n 3 true
go run ./tools/ossbench -json -n 2 true
go run ./tools/hygienecheck
git diff --check
```

## Reviewer Instructions
- Examine the code for issues across all four axes.
- Report file:line evidence for every actionable finding.
- Categorize each finding by axis (quality/ux/dx/security).
- Output exactly `LGTM` only if zero findings across all axes.
