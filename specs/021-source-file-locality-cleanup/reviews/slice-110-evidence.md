# Slice 110 Evidence

Evidence date: 2026-06-05

## Scope

Slice 110 consolidates `internal/posture` metric row validation helpers:

- Deleted `internal/posture/posture_validate_metric_counts.go`.
- Deleted `internal/posture/posture_validate_metric_identity.go`.
- Moved those predicates into
  `internal/posture/posture_validate_metric_row.go`.
- Extended focused regression coverage in `internal/posture/posture_test.go`.

Excluded surfaces: movement row validation, movement summary validation, export
behavior outside metric row validation, schemas, examples, fixtures,
dependencies, package boundary, dependency direction, and MI baselines.

## Plan Review

Initial plan review found that focused coverage did not require all moved
`malformedMetricCounts` predicates. T021-7650 now requires negative numerator,
negative denominator, invalid unit, and negative not-assessed count coverage.

The first implementation shape also attempted movement helper consolidation,
but direct movement consolidation failed the file-level MI gate. Slice 110 was
narrowed to metric row validation only; movement row helper shards remain for a
later slice.

Final plan re-review: Beauvoir, Peirce, and Halley returned exactly `LGTM`.

## Focused Verification

Command:

```sh
tests='TestValidateMetricRowShapeRejectsMalformedRows'
listed=$(go test ./internal/posture -list "$tests" | rg "^($tests)$")
test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 1
go test ./internal/posture -run "$tests" -count=1 -v
```

Result: verified pass. The focused test includes negative numerator, negative
denominator, invalid unit, and negative not-assessed count cases.

## Repository Verification

Command:

```sh
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal &&
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools &&
go test ./... &&
go vet ./... &&
golangci-lint run &&
go run ./tools/doccheck &&
go run ./tools/hygienecheck &&
jq empty schema/*.json &&
git diff --check &&
go test -count=1 ./... -coverprofile=coverage.out &&
go tool cover -func=coverage.out > coverage-func.txt &&
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt &&
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less &&
rm -f coverage.out coverage-func.txt gocyclo.txt
```

Result: verified pass.

## Boundary Evidence

Verified:

- `git diff --name-status internal/posture` lists only metric row validation
  files and `posture_test.go`.
- Temporary coverage artifacts are absent after verification cleanup.
- No schema, example, dependency, or MI baseline files are changed.

## Implementation Review

| Lane | Harness | Agent id | Model/provider | Prompt class | Timeout | Retries | Fallback | Result |
|---|---|---|---|---|---|---|---|---|
| Beauvoir | Codex subagent | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | not_assessed | implementation review | 360s | 0 | none | LGTM |
| Peirce | Codex subagent | `019e9406-f40c-79f1-904e-54d0f0b73866` | not_assessed | implementation review | 360s | 0 | none | LGTM |
| Halley | Codex subagent | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | not_assessed | implementation review | 360s | 0 | none | LGTM |

Result: verified pass. All three reviewer lanes returned exactly `LGTM`.
