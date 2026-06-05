# Slice 128 Evidence

Date: 2026-06-05T03:39:27Z

Scope: residual filename cleanup for
`internal/harnessobs/validation_type.go`; rename to
`internal/harnessobs/validation.go` while preserving exported `Dimension` and
`Validation` declarations and JSON tags exactly.

## Implementation Summary

- Renamed `internal/harnessobs/validation_type.go` to
  `internal/harnessobs/validation.go`.
- Added raw JSON key assertions for persisted validation artifacts, nested
  dimension artifacts, `validation_digest` emission when populated, and
  `validation_digest` omission when `ValidationDigest` is empty.
- No schema, example, dependency, package-boundary, exported-name, or MI
  baseline changes were made.

## Focused Verification

Result: pass.

Command:

```sh
gofmt -w internal/harnessobs/validation.go internal/harnessobs/harnessobs_test.go && test ! -e internal/harnessobs/validation_type.go && test -e internal/harnessobs/validation.go && test "$(go test ./internal/harnessobs -list 'Test(ValidateWritesOutPathWhenPasses|ValidateCannotVerifyWhenRunFileInvalid|LoadValidationRejectsUnsafePathAndSchemaVersion)$' | grep -Ec '^Test(ValidateWritesOutPathWhenPasses|ValidateCannotVerifyWhenRunFileInvalid|LoadValidationRejectsUnsafePathAndSchemaVersion)$')" -eq 3 && go test ./internal/harnessobs -run 'Test(ValidateWritesOutPathWhenPasses|ValidateCannotVerifyWhenRunFileInvalid|LoadValidationRejectsUnsafePathAndSchemaVersion)$' -count=1 -v
```

Observed tests:

- `TestValidateWritesOutPathWhenPasses`: pass.
- `TestValidateCannotVerifyWhenRunFileInvalid`: pass.
- `TestLoadValidationRejectsUnsafePathAndSchemaVersion`: pass.

## Repository Verification

Result: pass.

Command:

```sh
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal && go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools && go test ./... && go vet ./... && golangci-lint run && go run ./tools/doccheck && go run ./tools/hygienecheck && jq empty schema/*.json && git diff --check && go test -count=1 ./... -coverprofile=coverage.out && go tool cover -func=coverage.out > coverage-func.txt && go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt && go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less && rm -f coverage.out coverage-func.txt gocyclo.txt
```

## Drift Checks

- Spec drift: pass; implementation matches Slice 128 rename-only scope.
- Constitution drift: pass; machine verification and reviewer evidence are
  recorded; no unchecked green claims.
- Product drift: pass; portable validation artifact contract remains unchanged.
- CRAP < 5: pass by local `crapcheck`.
- MI > 70: pass by local `qualitycheck`; no MI baseline changes.
- CleanArch hex: pass; package boundaries unchanged.
- CleanCode/SOLID/DRY/YAGNI: pass; no new production abstraction or dependency
  added.

## Implementation Review

State: pass.

| Reviewer | Agent ID | Harness | Provider/model | Result | Notes |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | LGTM | Returned exactly `LGTM`. |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | LGTM | Returned exactly `LGTM`. |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | LGTM | Returned exactly `LGTM`. |
