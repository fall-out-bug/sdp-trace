# Slice 44 Evidence

Status: pass

## Scope

Slice 44 is bounded to `internal/harnessobs/harnessobs_344`.

Implemented consolidation:

- moved `readEventsFromPath` into `internal/harnessobs/event_scan_input.go`
  next to event source scanning setup
- removed `internal/harnessobs/harnessobs_344_readeventsfrompath.go`
- added focused branch coverage for profile load success, profile load failure,
  source read failure, source digest return, and existing collect-session
  missing-source `cannot_verify` behavior

Explicit exclusions:

- profile-relative source/output path safety (`harnessobs_345` through
  `harnessobs_347`) remains unchanged and not reassessed by this slice
- raw-event normalization (`harnessobs_348` onward) remains unchanged and not
  reassessed by this slice

## Focused Verification

- pass: `go test ./internal/harnessobs -run 'Test(ReadEventsFromPathBranches|CollectSessionMarksMissingSourceUnavailable)$'`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 internal/harnessobs/event_scan_input.go`

## Repository Verification

- pass: `go test ./...`
- pass: `go vet ./...`
- pass: `golangci-lint run`
- pass: `go run ./tools/doccheck`
- pass: `go run ./tools/hygienecheck`
- pass: `jq empty schema/*.json`
- pass: `git diff --check`
- pass: `go test -count=1 ./... -coverprofile=coverage.out`
- pass: `go tool cover -func=coverage.out > coverage-func.txt`
- pass: `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- pass: `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- pass: `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- pass: `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`

No MI baseline files were changed. `event_scan_input.go` was measured at file
MI 72.7 while keeping the event-reading trust boundary documented in the
grouped file.

## Review Lanes

- correctness reviewer: `019e87ec-dbff-7a40-8bc8-94f111d74f34`,
  result `LGTM`.
- trust/evidence/spec-drift reviewer:
  `019e87ec-fc5e-74c2-8ec5-aef600884a3f`, initial result `major`;
  `019e87ee-d3c6-7011-b0cf-32b035ee225d`, re-review result `LGTM`.
- maintainability/DX reviewer: `019e87ee-f570-7eb0-b5ae-f595547cf19d`,
  initial result `minor`; `019e87f1-b9a4-7682-9f0a-2a0d3114c46a`,
  re-review result `LGTM`.

## Findings

- fixed: staged Slice 44 review artifacts before trust/evidence re-review so
  task closure matched the reviewable diff.
- fixed: replaced padding-style comments with trust-boundary notes for source
  path resolution responsibility, replay limits, and digest withholding on
  scan failure.
- final: all implementation reviewer lanes returned `LGTM`.
