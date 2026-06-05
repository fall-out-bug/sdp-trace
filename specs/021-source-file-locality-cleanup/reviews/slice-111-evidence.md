# Slice 111 Evidence

Evidence date: 2026-06-05

## Scope

Slice 111 consolidates `internal/posture` movement row identity validation:

- Deleted `internal/posture/posture_validate_movement_identity.go`.
- Moved `malformedMovementIdentity` into
  `internal/posture/posture_validate_movement_row.go`.

Excluded surfaces: movement value and comparison predicates, movement summary
validation, export behavior outside movement row validation, schemas, examples,
fixtures, dependencies, package boundary, dependency direction, and MI
baselines.

## Plan Review

Three independent reviewer lanes returned exactly `LGTM` for the narrow
movement identity scope. Full movement validation consolidation was not used
because prior feasibility checks showed it failed the file-level MI gate.

## Focused Verification

Command:

```sh
tests='TestPostureHelperCoverageForMovementAndSafety|TestValidateMovementRowRejectsMalformedRows'
listed=$(go test ./internal/posture -list "$tests" | rg "^($tests)$")
test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 2
go test ./internal/posture -run "$tests" -count=1 -v
```

Result: verified pass. Both focused tests were listed and passed.

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

- `git status --short` lists only Slice 111 code and SpecKit artifacts.
- Temporary coverage artifacts are absent after verification cleanup.
- No schema, example, dependency, movement summary, or MI baseline files are
  changed.

## Implementation Review

| Lane | Harness | Agent id | Model/provider | Prompt class | Timeout | Retries | Fallback | Result |
|---|---|---|---|---|---|---|---|
| Beauvoir | Codex subagent | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | not_assessed | implementation review | 360s | 0 | none | LGTM |
| Peirce | Codex subagent | `019e9406-f40c-79f1-904e-54d0f0b73866` | not_assessed | implementation review | 360s | 0 | none | LGTM |
| Halley | Codex subagent | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | not_assessed | implementation review | 360s | 0 | none | LGTM |

Result: verified pass. All three reviewer lanes returned exactly `LGTM`.
