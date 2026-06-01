# Spec 021 Slice 3 Evidence

Date: 2026-06-01.

## Scope

Implemented Slice 3 for `cmd/sdp-trace` command-surface list helpers.

Changed source shape:

- Added `cmd/sdp-trace/command_surface_list_helpers.go`.
- Removed two numbered list helper shards:
  - `cmd/sdp-trace/main_582_commandsurfaceflaglisthelpers.go`
  - `cmd/sdp-trace/main_583_commandsurfaceslicehelpers.go`

No command behavior, schema, package boundary, dependency direction, production
trust, release approval, external attestation, or MI-baseline change was
intended.

## Decision Notes

Rejected alternative: one broader argument-helper file that also included
`isHelp` and `isBoolLiteral`. Local pre-change metric analysis measured that
combined file at file MI `63.5`, below the absolute threshold, so Slice 3 was
kept to list helpers only.

## Verification

Verified pass:

- `gofmt -w cmd/sdp-trace/command_surface_list_helpers.go`
- `go test ./cmd/sdp-trace`
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/command_surface_list_helpers.go`
- `go test ./...`
- `go vet ./...`
- `go run ./tools/doccheck`
- `go run ./tools/hygienecheck`
- `jq empty schema/*.json`
- `git diff --check`
- `go test -count=1 ./... -coverprofile=coverage.out`
- `go tool cover -func=coverage.out > coverage-func.txt`
- `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- `go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- `go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`

Quality notes:

- `command_surface_list_helpers.go` file MI: `100.0`.
- Moved functions have function MI `100.0`.
- Moved functions have CRAP `1.00`.
- No MI baseline file was changed for this slice.

## Targeted Review

Verified pass:

- Harness: `opencode run`
- Model: `opencode-go/glm-5.1`
- Result: `LGTM`

Review notes:

- Behavior preservation: `reqFlags`, `optFlags`, `subs`, and `vars` remain
  equivalent after the move.
- Scope: no unrelated code, schema, package boundary, or baseline change found.
- Evidence: no production, release, external attestation, or merge approval
  overclaim found.

Not assessed:

- No production deployment, release approval, external attestation, or merge
  approval was assessed.
