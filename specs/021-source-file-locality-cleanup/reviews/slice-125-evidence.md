# Slice 125 Evidence

Date: 2026-06-05T03:14:09Z

Scope:

- Removed `internal/harnessobs/unavailable_field_type.go`.
- Moved exported `UnavailableField` into `internal/harnessobs/run_type.go`.
- Added focused raw JSON assertions for `UnavailableField` inside persisted
  `Event` JSONL, including zero-value emission of `field`, `state`, and
  `reason_code`, plus populated exact key emission.

Out of scope:

- Unavailable field validation.
- Event content validation.
- Event decoding.
- Event writing behavior.
- Run loading.
- Schemas, examples, dependencies, package boundary, dependency direction, MI
  baselines, and CRAP threshold.

Plan review:

- Review state: pass after Halley re-reviewed the strengthened zero-value raw
  JSON assertion requirement; Beauvoir and Peirce returned `LGTM` in round 1.

Focused verification:

```sh
gofmt -w internal/harnessobs/run_type.go internal/harnessobs/harnessobs_test.go &&
test "$(go test ./internal/harnessobs -list 'Test(WriteNormalizedEventsWritesJSONL|ValidateEventRejectsInvalidFields)$' | grep -Ec '^Test(WriteNormalizedEventsWritesJSONL|ValidateEventRejectsInvalidFields)$')" -eq 2 &&
go test ./internal/harnessobs -run 'Test(WriteNormalizedEventsWritesJSONL|ValidateEventRejectsInvalidFields)$' -count=1 -v
```

Result: pass.

Repository verification:

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

Result: pass.

Drift checks:

- Spec drift: pass. The implementation matches the reviewed Slice 125 target
  and preserves `UnavailableField` fields and JSON tags.
- Constitution drift: pass. No harness/runtime dependency or opaque trust score
  was added.
- Product drift: pass. The change remains a source-file locality cleanup.
- CRAP < 5: pass.
- MI > 70: pass without baseline changes.
- CleanArch hex: pass. Package boundary and dependency direction unchanged.
- CleanCode/SOLID/DRY/YAGNI: pass. One type micro-file was removed without
  adding a new abstraction or dependency.

Implementation review state: pass

| Reviewer | Harness | Agent ID | Model/provider | Prompt class | Result |
|---|---|---:|---|---|---|
| Beauvoir | Codex subagent | 019e9406-f078-7fd2-b8d0-e22ac17a1e3a | not_assessed | implementation review | LGTM |
| Peirce | Codex subagent | 019e9406-f40c-79f1-904e-54d0f0b73866 | not_assessed | implementation review | LGTM |
| Halley | Codex subagent | 019e9406-f7c2-7f80-80d9-86f7cf7e0c22 | not_assessed | implementation review | LGTM |
