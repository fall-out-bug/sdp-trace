# Slice 1 Review Context: tools/osscompat

## Files Changed
- `tools/osscompat/main.go`
- `tools/osscompat/main_test.go`
- `tools/osscompat/probe.go`
- `tools/osscompat/probe_test.go`
- `tools/osscompat/runner.go`

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
- Scanner verification commands in docs must be copy-pasteable.

## Verification Commands
```text
go test -count=1 ./tools/osscompat
go vet ./tools/osscompat
go run ./tools/osscompat        # exit 0 when no fail probes; exit 1 when jsonschema-wrap-drift reports fail
go run ./tools/osscompat -json
go run ./tools/osscompat -probe jsonschema-fixtures
go run ./tools/hygienecheck
git diff --check
```

## Reviewer Instructions
- Examine the code for issues across all four axes.
- Report file:line evidence for every actionable finding.
- Categorize each finding by axis (quality/ux/dx/security).
- Output exactly `LGTM` only if zero findings across all axes.
