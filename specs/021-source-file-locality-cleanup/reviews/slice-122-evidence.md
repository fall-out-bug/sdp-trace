# Slice 122 Evidence

Date: 2026-06-05T02:42:13Z

Scope:

- Removed `internal/harnessobs/event_type.go`.
- Moved exported `Event` into `internal/harnessobs/run_type.go`.
- Added focused raw JSON assertions for persisted normalized event fields in
  `TestWriteNormalizedEventsWritesJSONL`, including omitted empty optional
  fields and populated `task_ref`, `operation_ref`, and `unavailable_fields`.

Out of scope:

- Event decoding.
- Event identity/ref/content validation.
- Event scanning.
- Event writing behavior.
- Run loading.
- Normalized event generation.
- Schemas, examples, dependencies, package boundary, dependency direction, MI
  baselines, and CRAP threshold.

Plan review:

- Review state: pass from Beauvoir, Peirce, and Halley.

Focused verification:

```sh
gofmt -w internal/harnessobs/run_type.go internal/harnessobs/harnessobs_test.go &&
test "$(go test ./internal/harnessobs -list 'Test(WriteNormalizedEventsWritesJSONL|CollectSessionWritesObservedRun)$' | grep -Ec '^Test(WriteNormalizedEventsWritesJSONL|CollectSessionWritesObservedRun)$')" -eq 2 &&
go test ./internal/harnessobs -run 'Test(WriteNormalizedEventsWritesJSONL|CollectSessionWritesObservedRun)$' -count=1 -v
```

Result: pass after aligning the new raw JSON assertion with the existing
normalized event contract: `content_state` remains `digest_only`, adding a
persisted event with populated optional fields, and asserting empty optional
fields stay omitted.

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

- Spec drift: pass. The implementation matches the reviewed Slice 122 target
  and preserves `Event` fields and JSON tags.
- Constitution drift: pass. No harness/runtime dependency or opaque trust score
  was added.
- Product drift: pass. The change remains a source-file locality cleanup.
- CRAP < 5: pass.
- MI > 70: pass without baseline changes.
- CleanArch hex: pass. Package boundary and dependency direction unchanged.
- CleanCode/SOLID/DRY/YAGNI: pass. One type micro-file was removed without
  adding a new abstraction or dependency.

Implementation review round 1 state: fail

| Reviewer | Harness | Agent ID | Model/provider | Prompt class | Result |
|---|---|---:|---|---|---|
| Beauvoir | Codex subagent | 019e9406-f078-7fd2-b8d0-e22ac17a1e3a | not_assessed | implementation review | major: raw JSON assertion did not cover populated optional `Event` tags or omitempty behavior |
| Peirce | Codex subagent | 019e9406-f40c-79f1-904e-54d0f0b73866 | not_assessed | implementation review | major: evidence overclaimed preservation of `task_ref`, `operation_ref`, and `unavailable_fields` JSON tags |
| Halley | Codex subagent | 019e9406-f7c2-7f80-80d9-86f7cf7e0c22 | not_assessed | implementation review | major: test did not persist an event with non-empty optional fields |

Round 1 fix:

- Added a third persisted JSONL event with `TaskRef`, `OperationRef`, and
  `UnavailableFields`.
- Asserted exact raw JSON keys for `task_ref`, `operation_ref`, and nested
  `unavailable_fields` values.
- Asserted the previous event omits `task_ref`, `operation_ref`, and
  `unavailable_fields` when empty.

Implementation re-review state: pass

| Reviewer | Harness | Agent ID | Model/provider | Prompt class | Result |
|---|---|---:|---|---|---|
| Beauvoir | Codex subagent | 019e9406-f078-7fd2-b8d0-e22ac17a1e3a | not_assessed | implementation re-review | LGTM |
| Peirce | Codex subagent | 019e9406-f40c-79f1-904e-54d0f0b73866 | not_assessed | implementation re-review | LGTM |
| Halley | Codex subagent | 019e9406-f7c2-7f80-80d9-86f7cf7e0c22 | not_assessed | implementation re-review | LGTM |
