# Spec 021 Slice 4 Evidence

Date: 2026-06-01.

## Scope

Implemented Slice 4 for `cmd/sdp-trace` command-surface schema and runner
shards.

Changed source shape:

- Added `cmd/sdp-trace/command_surface_schema.go`.
- Added `cmd/sdp-trace/command_surface_runner.go`.
- Removed three numbered shards:
  - `cmd/sdp-trace/main_536_commandsurfaceschema.go`
  - `cmd/sdp-trace/main_544_commandsurfacejson.go`
  - `cmd/sdp-trace/main_545_runcommandsurface.go`

No command behavior, schema contract, package boundary, dependency direction,
production trust, release approval, external attestation, or MI-baseline change
was intended.

## Decision Notes

Rejected alternatives:

- One combined schema/metadata file: local pre-change metric analysis measured
  file MI `59.1`.
- Combined metadata/registry file: local pre-change metric analysis measured
  file MI `56.1`.
- Combined constants/registry file: local pre-change metric analysis measured
  file MI `63.3`.

These alternatives would risk a mixed code/baseline PR, so Slice 4 kept schema
types and runner behavior separate.

## Verification

Verified pass:

- `gofmt -w cmd/sdp-trace/command_surface_schema.go cmd/sdp-trace/command_surface_runner.go`
- `go test ./cmd/sdp-trace`
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/command_surface_schema.go cmd/sdp-trace/command_surface_runner.go`
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

- `command_surface_schema.go` file MI: `100.0`.
- `command_surface_runner.go` file MI: `74.2`.
- Moved runner functions have function MI `100.0`.
- Moved runner functions have CRAP below `5`.
- No MI baseline file was changed for this slice.

## Targeted Reviews

Verified pass:

- Harness: `opencode run`
- Model: `opencode-go/glm-5.1`
- Result: `LGTM` after finding fix.
- Harness: `opencode run`
- Model: `kimi-for-coding/k2p6`
- Result: `LGTM`
- Harness: `opencode run`
- Model: `opencode-go/minimax-m3`
- Result: `LGTM`

Review notes:

- Behavior preservation: moved schema types and runner/JSON functions remain
  equivalent after the move.
- Scope: no unrelated code, schema, package boundary, dependency, or baseline
  change found.
- Evidence: production, release, external attestation, and merge approval remain
  `not_assessed`.
- Reviewer finding fixed: the placeholder review section no longer claims
  `Verified pass` with `pending` content.

Not assessed:

- No production deployment, release approval, external attestation, or merge
  approval was assessed.
