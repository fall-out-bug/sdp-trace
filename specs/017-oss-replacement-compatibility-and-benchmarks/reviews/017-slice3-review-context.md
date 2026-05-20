# Slice 3 Review Context: WS-017-B Wrap Output / Schema Drift

## Files Changed
- `examples/flight-recorder/wrap-output-drift/run.json`
- `examples/flight-recorder/wrap-output-drift/README.md`
- `tools/osscompat/probe_test.go`

## Review Axes
1. **Quality:** Correctness of drift evidence, test accuracy, Go idioms.
2. **UX:** Clarity of documentation, discoverability of the blocker.
3. **DX:** File organization, example structure, test maintainability.
4. **Security:** No secrets in fixtures, no unsafe file operations.

## Verification Commands
```text
go test -count=1 ./tools/osscompat
go vet ./tools/osscompat
git diff --check
```

## Reviewer Instructions
- Examine the drift evidence for honesty and completeness.
- Check that the test actually validates the structural mismatch.
- Report file:line evidence for every actionable finding.
- Categorize each finding by axis (quality/ux/dx/security).
- Output exactly `LGTM` only if zero findings across all axes.
