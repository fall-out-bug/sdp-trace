# Slice 127 Evidence

Date: 2026-06-05T03:31:03Z

Scope: residual filename cleanup for `internal/harnessobs/run_type.go`; rename
to `internal/harnessobs/run.go` while preserving exported `UnavailableField`,
`Event`, and `Run` declarations and JSON tags exactly.

## Implementation Summary

- Renamed `internal/harnessobs/run_type.go` to `internal/harnessobs/run.go`.
- Added raw JSON key assertions to `TestCollectSessionWritesObservedRun` for
  persisted `run.json` keys: `schema_version`, `profile_id`,
  `harness_family`, `event_schema_version`, `source_path`, `source_digest`,
  `event_count`, `event_refs`, and `created_at`.
- No schema, example, dependency, package-boundary, exported-name, or MI
  baseline changes were made.

## Focused Verification

Result: pass.

Command:

```sh
test ! -e internal/harnessobs/run_type.go && test -e internal/harnessobs/run.go && test "$(go test ./internal/harnessobs -list 'Test(WriteNormalizedEventsWritesJSONL|CollectSessionWritesObservedRun|LoadRunRejectsUnsafeEventRefs)$' | grep -Ec '^Test(WriteNormalizedEventsWritesJSONL|CollectSessionWritesObservedRun|LoadRunRejectsUnsafeEventRefs)$')" -eq 3 && go test ./internal/harnessobs -run 'Test(WriteNormalizedEventsWritesJSONL|CollectSessionWritesObservedRun|LoadRunRejectsUnsafeEventRefs)$' -count=1 -v
```

Observed tests:

- `TestWriteNormalizedEventsWritesJSONL`: pass.
- `TestCollectSessionWritesObservedRun`: pass.
- `TestLoadRunRejectsUnsafeEventRefs`: pass.

## Repository Verification

Result: pass.

Command:

```sh
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal && go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools && go test ./... && go vet ./... && golangci-lint run && go run ./tools/doccheck && go run ./tools/hygienecheck && jq empty schema/*.json && git diff --check && go test -count=1 ./... -coverprofile=coverage.out && go tool cover -func=coverage.out > coverage-func.txt && go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt && go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less && rm -f coverage.out coverage-func.txt gocyclo.txt
```

## Drift Checks

- Spec drift: pass; implementation matches Slice 127 rename-only scope.
- Constitution drift: pass; machine verification and reviewer evidence are
  recorded; no unchecked green claims.
- Product drift: pass; portable harness observation contract remains unchanged.
- CRAP < 5: pass by local `crapcheck`.
- MI > 70: pass by local `qualitycheck`; no MI baseline changes.
- CleanArch hex: pass; package boundaries unchanged.
- CleanCode/SOLID/DRY/YAGNI: pass; no new abstraction or dependency added.

## Implementation Review

State: pass.

| Reviewer | Agent ID | Harness | Provider/model | Result | Notes |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | LGTM | Returned exactly `LGTM`. |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | LGTM | Returned exactly `LGTM`. |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | LGTM | Returned exactly `LGTM`. |
